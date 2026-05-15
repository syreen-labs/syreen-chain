package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// NewSlashingTxCmd returns the tx commands for the slashing module.
func NewSlashingTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "slashing",
		Short:                      "Slashing transaction subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewUnjailCmd(),
	)

	return cmd
}

// NewUnjailCmd returns a command to unjail a validator.
func NewUnjailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unjail",
		Short: "Unjail a jailed validator",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valAddr := clientCtx.GetFromAddress()

			msg := &slashingtypes.MsgUnjail{
				ValidatorAddr: sdk.ValAddress(valAddr).String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
