package cli

import (
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"syreen/x/predict/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: types.ModuleName, Short: "Querying commands for the predict module",
		DisableFlagParsing: true, SuggestionsMinimumDistance: 2, RunE: client.ValidateCmd,
	}
	cmd.AddCommand(CmdQueryMarkets())
	return cmd
}

func CmdQueryMarkets() *cobra.Command {
	cmd := &cobra.Command{
		Use: "markets", Short: "Query all prediction markets", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil { return err }
			return clientCtx.PrintString("Use the gRPC/REST endpoint: /syreen/predict/v1/markets\n")
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
