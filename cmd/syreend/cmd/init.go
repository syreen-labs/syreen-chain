package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	cfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	"cosmossdk.io/math/unsafe"

	"syreen/app"
)

const (
	flagOverwrite = "overwrite"
	flagRecover   = "recover"
)

type initPrintInfo struct {
	Moniker    string          `json:"moniker" yaml:"moniker"`
	ChainID    string          `json:"chain_id" yaml:"chain_id"`
	NodeID     string          `json:"node_id" yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir" yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
}

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application, using usyreen as the bond denom.
func InitCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			switch {
			case chainID != "":
			case clientCtx.ChainID != "":
				chainID = clientCtx.ChainID
			default:
				chainID = fmt.Sprintf("test-chain-%v", unsafe.Str(6))
			}

			// Get bip39 mnemonic
			var mnemonic string
			recover, _ := cmd.Flags().GetBool(flagRecover)
			if recover {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				value, err := input.GetString("Enter your bip39 mnemonic", inBuf)
				if err != nil {
					return err
				}

				mnemonic = value
				if !bip39.IsMnemonicValid(mnemonic) {
					return errors.New("invalid mnemonic")
				}
			}

			// Get initial height
			initHeight, _ := cmd.Flags().GetInt64(flags.FlagInitHeight)
			if initHeight < 1 {
				initHeight = 1
			}

			nodeID, _, err := genutil.InitializeNodeValidatorFilesFromMnemonic(config, mnemonic)
			if err != nil {
				return err
			}

			config.Moniker = args[0]

			genFile := config.GenesisFile()
			overwrite, _ := cmd.Flags().GetBool(flagOverwrite)

			// use os.Stat to check if the file exists
			_, err = os.Stat(genFile)
			if !overwrite && !os.IsNotExist(err) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}

			// Use our custom DefaultGenesisState which patches denoms to usyreen
			appGenState := app.DefaultGenesisState(cdc)

			appState, err := json.MarshalIndent(appGenState, "", " ")
			if err != nil {
				return fmt.Errorf("failed to marshal default genesis state: %w", err)
			}

			appGenesis := &genutiltypes.AppGenesis{}
			if _, err := os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else {
				appGenesis, err = genutiltypes.AppGenesisFromFile(genFile)
				if err != nil {
					return fmt.Errorf("failed to read genesis doc from file: %w", err)
				}
			}

			appGenesis.AppName = version.AppName
			appGenesis.AppVersion = version.Version
			appGenesis.ChainID = chainID
			appGenesis.AppState = appState
			appGenesis.InitialHeight = initHeight
			appGenesis.Consensus = &genutiltypes.ConsensusGenesis{
				Validators: nil,
			}

			if err = genutil.ExportGenesisFile(appGenesis, genFile); err != nil {
				return fmt.Errorf("failed to export genesis file: %w", err)
			}

			toPrint := initPrintInfo{
				Moniker:    config.Moniker,
				ChainID:    chainID,
				NodeID:     nodeID,
				AppMessage: appState,
			}

			out, err := json.MarshalIndent(toPrint, "", " ")
			if err != nil {
				return err
			}

			_, err = fmt.Fprintf(os.Stderr, "%s\n", out)

			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
			return err
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().Bool(flagRecover, false, "provide seed phrase to recover existing key instead of creating")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().Int64(flags.FlagInitHeight, 1, "specify the initial block height at genesis")

	return cmd
}
