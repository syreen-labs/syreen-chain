package cmd

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/log"

	tmcfg "github.com/cometbft/cometbft/config"
	tmcli "github.com/cometbft/cometbft/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/snapshot"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	distrcli "github.com/cosmos/cosmos-sdk/x/distribution/client/cli"

	ibccorecli "github.com/cosmos/ibc-go/v10/modules/core/client/cli"
	ibctransfercli "github.com/cosmos/ibc-go/v10/modules/apps/transfer/client/cli"

	"syreen/app"

	abstractaccountcli "syreen/x/abstractaccount/client/cli"
	computecli "syreen/x/compute/client/cli"
	feemarketcli "syreen/x/feemarket/client/cli"
	identitycli "syreen/x/identity/client/cli"
	intentcli "syreen/x/intent/client/cli"
	mevprotectioncli "syreen/x/mevprotection/client/cli"
	paymentscli "syreen/x/payments/client/cli"
	tokenfactorycli "syreen/x/tokenfactory/client/cli"
	evmcli "syreen/x/evm/client/cli"
	dexcli "syreen/x/dex/client/cli"
	clmmcli "syreen/x/clmm/client/cli"
	vaultcli "syreen/x/vault/client/cli"
	portfoliocli "syreen/x/portfolio/client/cli"
	farmingcli "syreen/x/farming/client/cli"
	perpscli "syreen/x/perps/client/cli"
	lendingcli "syreen/x/lending/client/cli"
	flashloancli "syreen/x/flashloan/client/cli"
	predictcli "syreen/x/predict/client/cli"
	launchpadcli "syreen/x/launchpad/client/cli"
	optionscli "syreen/x/options/client/cli"
	aiagentcli "syreen/x/aiagent/client/cli"

	syreenconfig "syreen/config"
)

// NewRootCmd creates a new root command for syreend daemon
func NewRootCmd() *cobra.Command {
	encodingConfig := app.MakeEncodingConfig()

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithHomeDir(app.DefaultNodeHome).
		WithViper("")

	rootCmd := &cobra.Command{
		Use:   "syreend",
		Short: "Syreen Blockchain Daemon",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx = initClientCtx.WithCmdContext(cmd.Context())
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customTMConfig := initTendermintConfig()

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customTMConfig)
		},
	}

	initRootCmd(rootCmd, encodingConfig, app.ModuleBasics)

	return rootCmd
}

// initTendermintConfig helps to override default Tendermint Config values for optimal performance
func initTendermintConfig() *tmcfg.Config {
	cfg := tmcfg.DefaultConfig()

	// Apply Syreen's optimized consensus parameters from config/consensus.go
	cfg.Consensus.TimeoutPropose = syreenconfig.ConsensusConfig.TimeoutPropose
	cfg.Consensus.TimeoutPrevote = syreenconfig.ConsensusConfig.TimeoutPrevote
	cfg.Consensus.TimeoutPrecommit = syreenconfig.ConsensusConfig.TimeoutPrecommit
	cfg.Consensus.TimeoutCommit = syreenconfig.ConsensusConfig.TimeoutCommit

	// Produce empty blocks to maintain IBC client updates, time-based logic, and consistent block times
	cfg.Consensus.CreateEmptyBlocks = syreenconfig.ConsensusConfig.CreateEmptyBlocks
	cfg.Consensus.CreateEmptyBlocksInterval = 0

	// P2P optimizations
	cfg.P2P.MaxNumInboundPeers = 100
	cfg.P2P.MaxNumOutboundPeers = 40
	cfg.P2P.SendRate = 20480000 // 20 MB/s
	cfg.P2P.RecvRate = 20480000 // 20 MB/s
	cfg.P2P.FlushThrottleTimeout = 10 * time.Millisecond

	// Mempool settings for high throughput
	cfg.Mempool.Size = 10000
	cfg.Mempool.MaxTxsBytes = 134217728 // 128MB — reasonable for high throughput
	cfg.Mempool.CacheSize = 20000

	return cfg
}

func initAppConfig() (string, interface{}) {
	type CustomAppConfig struct {
		serverconfig.Config
	}

	srvCfg := serverconfig.DefaultConfig()
	srvCfg.MinGasPrices = "0.001usyreen"

	// API configuration — default to localhost for safety; operators override for external access
	srvCfg.API.Enable = true
	srvCfg.API.Swagger = false
	srvCfg.API.Address = "tcp://127.0.0.1:1317"

	// gRPC configuration — default to localhost for safety
	srvCfg.GRPC.Enable = true
	srvCfg.GRPC.Address = "127.0.0.1:9090"

	// State sync
	srvCfg.StateSync.SnapshotInterval = 1000
	srvCfg.StateSync.SnapshotKeepRecent = 2

	customAppConfig := CustomAppConfig{
		Config: *srvCfg,
	}

	customAppTemplate := serverconfig.DefaultConfigTemplate

	return customAppTemplate, customAppConfig
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig app.EncodingConfig, basicManager module.BasicManager) {
	cfg := sdk.GetConfig()
	cfg.Seal()

	rootCmd.AddCommand(
		InitCmd(app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome, genutiltypes.DefaultMessageValidator, encodingConfig.TxConfig.SigningContext().ValidatorAddressCodec()),
		genutilcli.GenTxCmd(basicManager, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome, encodingConfig.TxConfig.SigningContext().ValidatorAddressCodec()),
		genutilcli.ValidateGenesisCmd(basicManager),
		AddGenesisAccountCmd(app.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		debug.Cmd(),
		pruning.Cmd(newApp, app.DefaultNodeHome),
		snapshot.Cmd(newApp),
	)

	server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, appExport, addModuleInitFlags)

	// Add keyring commands
	rootCmd.AddCommand(
		genesisCommand(encodingConfig, basicManager),
		queryCommand(basicManager),
		txCommand(basicManager),
		keys.Commands(),
	)
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

// genesisCommand builds genesis-related cobra.Command
func genesisCommand(encodingConfig app.EncodingConfig, basicManager module.BasicManager, cmds ...*cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "genesis",
		Short:                      "Application's genesis-related subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	gentxModule := basicManager[genutiltypes.ModuleName].(genutil.AppModuleBasic)

	cmd.AddCommand(
		InitCmd(app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome, gentxModule.GenTxValidator, encodingConfig.TxConfig.SigningContext().ValidatorAddressCodec()),
		genutilcli.GenTxCmd(basicManager, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome, encodingConfig.TxConfig.SigningContext().ValidatorAddressCodec()),
		genutilcli.ValidateGenesisCmd(basicManager),
		AddGenesisAccountCmd(app.DefaultNodeHome),
	)

	for _, subCmd := range cmds {
		cmd.AddCommand(subCmd)
	}
	return cmd
}

func queryCommand(_ module.BasicManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.ValidatorCommand(),
		server.QueryBlockCmd(),
		authcmd.QueryTxsByEventsCmd(),
		server.QueryBlocksCmd(),
		authcmd.QueryTxCmd(),
		server.QueryBlockResultsCmd(),
		// Core SDK query commands (SDK v0.50 moved these to AutoCLI, so we provide manual gRPC wrappers)
		NewBankQueryCmd(),
		NewStakingQueryCmd(),
		NewDistributionQueryCmd(),
		NewGovQueryCmd(),
		NewSlashingQueryCmd(),
		NewAuthQueryCmd(),
		NewEvidenceQueryCmd(),
		// IBC query commands
		ibccorecli.GetQueryCmd(),
		ibctransfercli.GetQueryCmd(),
		// Custom module query commands
		tokenfactorycli.GetQueryCmd(),
		feemarketcli.GetQueryCmd(),
		mevprotectioncli.GetQueryCmd(),
		abstractaccountcli.GetQueryCmd(),
		identitycli.GetQueryCmd(),
		intentcli.GetQueryCmd(),
		paymentscli.GetQueryCmd(),
		computecli.GetQueryCmd(),
		evmcli.GetQueryCmd(),
		dexcli.GetQueryCmd(),
		clmmcli.GetQueryCmd(),
		vaultcli.GetQueryCmd(),
		portfoliocli.GetQueryCmd(),
		farmingcli.GetQueryCmd(),
		perpscli.GetQueryCmd(),
		lendingcli.GetQueryCmd(),
		flashloancli.GetQueryCmd(),
		predictcli.GetQueryCmd(),
		launchpadcli.GetQueryCmd(),
		optionscli.GetQueryCmd(),
		aiagentcli.GetQueryCmd(),
	)

	return cmd
}

func txCommand(_ module.BasicManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		authcmd.GetSimulateCmd(),
		// Core SDK tx commands
		bankcli.NewTxCmd(address.NewBech32Codec("syreen")),
		NewStakingTxCmd(), // Custom staking cmds — workaround for SDK v0.50 bech32 address codec bug
		govcli.NewTxCmd(nil),
		distrcli.NewTxCmd(address.NewBech32Codec("syreenvaloper"), address.NewBech32Codec("syreen")),
		NewSlashingTxCmd(),
		// Custom module tx commands
		tokenfactorycli.GetTxCmd(),
		mevprotectioncli.GetTxCmd(),
		abstractaccountcli.GetTxCmd(),
		intentcli.GetTxCmd(),
		computecli.GetTxCmd(),
		evmcli.GetTxCmd(),
		dexcli.GetTxCmd(),
		clmmcli.GetTxCmd(),
		vaultcli.GetTxCmd(),
		portfoliocli.GetTxCmd(),
		farmingcli.GetTxCmd(),
		perpscli.GetTxCmd(),
		lendingcli.GetTxCmd(),
		flashloancli.GetTxCmd(),
		predictcli.GetTxCmd(),
		launchpadcli.GetTxCmd(),
		optionscli.GetTxCmd(),
		aiagentcli.GetTxCmd(),
		// IBC
		ibctransfercli.NewTxCmd(),
	)

	return cmd
}

// newApp creates the application
func newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	baseappOptions := server.DefaultBaseappOptions(appOpts)

	return app.NewSyreenApp(
		logger,
		db,
		traceStore,
		true,
		appOpts,
		baseappOptions...,
	)
}

// appExport creates a new app (optionally at a given height) and exports state.
func appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	var syreenApp *app.SyreenApp

	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	viperAppOpts, ok := appOpts.(*viper.Viper)
	if !ok {
		return servertypes.ExportedApp{}, errors.New("appOpts is not viper.Viper")
	}

	viperAppOpts.Set(server.FlagInvCheckPeriod, 1)
	appOpts = viperAppOpts

	if height != -1 {
		syreenApp = app.NewSyreenApp(logger, db, traceStore, false, appOpts)

		if err := syreenApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		syreenApp = app.NewSyreenApp(logger, db, traceStore, true, appOpts)
	}

	return syreenApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}
