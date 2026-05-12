package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"cosmossdk.io/math"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// NewStakingTxCmd returns custom staking tx commands that work around the SDK
// v0.50 bech32 address codec bug where validator addresses get validated
// against the account address codec during tx generation.
func NewStakingTxCmd() *cobra.Command {
	stakingTxCmd := &cobra.Command{
		Use:                        "staking",
		Short:                      "Staking transaction subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	stakingTxCmd.AddCommand(
		NewCreateValidatorCmd(),
		NewDelegateCmd(),
		NewUndelegateCmd(),
		NewRedelegateCmd(),
	)

	return stakingTxCmd
}

// createValidatorJSON is the file format for create-validator.
type createValidatorJSON struct {
	PubKey            json.RawMessage `json:"pubkey"`
	Amount            string          `json:"amount"`
	Moniker           string          `json:"moniker"`
	Identity          string          `json:"identity"`
	Website           string          `json:"website"`
	Security          string          `json:"security"`
	Details           string          `json:"details"`
	CommissionRate    string          `json:"commission-rate"`
	CommissionMaxRate string          `json:"commission-max-rate"`
	CommissionMaxChange string        `json:"commission-max-change-rate"`
	MinSelfDelegation  string         `json:"min-self-delegation"`
}

func NewCreateValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator [path/to/validator.json]",
		Short: "Create a new validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			data, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read validator file: %w", err)
			}

			var vj createValidatorJSON
			if err := json.Unmarshal(data, &vj); err != nil {
				return fmt.Errorf("failed to parse validator json: %w", err)
			}

			amount, err := sdk.ParseCoinNormalized(vj.Amount)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}

			commRate, err := math.LegacyNewDecFromStr(vj.CommissionRate)
			if err != nil {
				return fmt.Errorf("invalid commission rate: %w", err)
			}
			commMaxRate, err := math.LegacyNewDecFromStr(vj.CommissionMaxRate)
			if err != nil {
				return fmt.Errorf("invalid commission max rate: %w", err)
			}
			commMaxChange, err := math.LegacyNewDecFromStr(vj.CommissionMaxChange)
			if err != nil {
				return fmt.Errorf("invalid commission max change rate: %w", err)
			}
			minSelfDel, ok := math.NewIntFromString(vj.MinSelfDelegation)
			if !ok {
				return fmt.Errorf("invalid min self delegation: %s", vj.MinSelfDelegation)
			}

			var pk cryptotypes.PubKey
			if err := clientCtx.Codec.UnmarshalInterfaceJSON(vj.PubKey, &pk); err != nil {
				return fmt.Errorf("failed to parse pubkey: %w", err)
			}

			pkAny, err := codectypes.NewAnyWithValue(pk)
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()

			msg := &stakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{
					Moniker:         vj.Moniker,
					Identity:        vj.Identity,
					Website:         vj.Website,
					SecurityContact: vj.Security,
					Details:         vj.Details,
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          commRate,
					MaxRate:       commMaxRate,
					MaxChangeRate: commMaxChange,
				},
				MinSelfDelegation: minSelfDel,
				DelegatorAddress:  delAddr.String(),
				ValidatorAddress:  sdk.ValAddress(delAddr).String(),
				Pubkey:            pkAny,
				Value:             amount,
			}

			return broadcastStakingMsg(clientCtx, cmd, msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewDelegateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate [validator-addr] [amount]",
		Short: "Delegate liquid tokens to a validator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			msg := stakingtypes.NewMsgDelegate(delAddr.String(), args[0], amount)

			return broadcastStakingMsg(clientCtx, cmd, msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewUndelegateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbond [validator-addr] [amount]",
		Short: "Unbond shares from a validator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			msg := stakingtypes.NewMsgUndelegate(delAddr.String(), args[0], amount)

			return broadcastStakingMsg(clientCtx, cmd, msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewRedelegateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redelegate [src-validator-addr] [dst-validator-addr] [amount]",
		Short: "Redelegate tokens from one validator to another",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			msg := stakingtypes.NewMsgBeginRedelegate(delAddr.String(), args[0], args[1], amount)

			return broadcastStakingMsg(clientCtx, cmd, msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// broadcastStakingMsg builds, signs and broadcasts a staking message directly
// without going through the SDK's broken address validation.
func broadcastStakingMsg(clientCtx client.Context, cmd *cobra.Command, msg sdk.Msg) error {
	txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
	if err != nil {
		return err
	}

	txf, err = txf.Prepare(clientCtx)
	if err != nil {
		return err
	}

	txBuilder := clientCtx.TxConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msg); err != nil {
		return err
	}

	txBuilder.SetFeeAmount(txf.Fees())
	txBuilder.SetGasLimit(txf.Gas())
	txBuilder.SetMemo(txf.Memo())

	if err = tx.Sign(cmd.Context(), txf, clientCtx.FromName, txBuilder, true); err != nil {
		return err
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return err
	}

	res, err := clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	return clientCtx.PrintProto(res)
}

func init() {
	_ = fmt.Sprintf // avoid unused import
}
