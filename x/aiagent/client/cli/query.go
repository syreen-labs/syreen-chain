package cli

import (
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"syreen/x/aiagent/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: types.ModuleName, Short: "Querying commands for the AI agent module",
		DisableFlagParsing: true, SuggestionsMinimumDistance: 2, RunE: client.ValidateCmd,
	}
	cmd.AddCommand(CmdQueryAgents())
	return cmd
}

func CmdQueryAgents() *cobra.Command {
	cmd := &cobra.Command{
		Use: "agents", Short: "Query all AI agents", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil { return err }
			return clientCtx.PrintString("Use the gRPC/REST endpoint: /syreen/aiagent/v1/agents\n")
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
