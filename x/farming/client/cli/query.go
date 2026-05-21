package cli

import (
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"syreen/x/farming/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: types.ModuleName, Short: "Querying commands for the farming module",
		DisableFlagParsing: true, SuggestionsMinimumDistance: 2, RunE: client.ValidateCmd,
	}
	cmd.AddCommand(CmdQueryFarms())
	return cmd
}

func CmdQueryFarms() *cobra.Command {
	cmd := &cobra.Command{
		Use: "farms", Short: "Query all farming pools", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil { return err }
			return clientCtx.PrintString("Use the gRPC/REST endpoint: /syreen/farming/v1/farms\n")
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
