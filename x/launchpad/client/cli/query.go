package cli

import (
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"syreen/x/launchpad/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: types.ModuleName, Short: "Querying commands for the launchpad module",
		DisableFlagParsing: true, SuggestionsMinimumDistance: 2, RunE: client.ValidateCmd,
	}
	cmd.AddCommand(CmdQueryLaunches())
	return cmd
}

func CmdQueryLaunches() *cobra.Command {
	cmd := &cobra.Command{
		Use: "launches", Short: "Query all token launches", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil { return err }
			return clientCtx.PrintString("Use the gRPC/REST endpoint: /syreen/launchpad/v1/launches\n")
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
