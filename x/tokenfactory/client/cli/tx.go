package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/tokenfactory/types"
)

// GetTxCmd returns the transaction commands for the tokenfactory module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Token factory transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdCreateDenom(),
		CmdMint(),
		CmdBurn(),
		CmdChangeAdmin(),
	)

	return cmd
}

func CmdCreateDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-denom [subdenom]",
		Short: "Create a new token denomination (factory/{creator}/{subdenom})",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgCreateDenom{
				Sender:   clientCtx.GetFromAddress().String(),
				Subdenom: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdMint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [denom] [amount] [optional:recipient]",
		Short: "Mint tokens of a factory denomination",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			denom := args[0]
			amount, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}

			mintTo := ""
			if len(args) == 3 {
				mintTo = args[2]
			}

			msg := &types.MsgMint{
				Sender: clientCtx.GetFromAddress().String(),
				Amount: sdk.NewCoin(denom, math.NewInt(amount)),
				MintTo: mintTo,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdBurn() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [denom] [amount]",
		Short: "Burn tokens of a factory denomination",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			denom := args[0]
			amount, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}

			msg := &types.MsgBurn{
				Sender: clientCtx.GetFromAddress().String(),
				Amount: sdk.NewCoin(denom, math.NewInt(amount)),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdChangeAdmin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-admin [denom] [new-admin]",
		Short: "Transfer admin of a factory denomination to a new address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgChangeAdmin{
				Sender:   clientCtx.GetFromAddress().String(),
				Denom:    args[0],
				NewAdmin: args[1],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
