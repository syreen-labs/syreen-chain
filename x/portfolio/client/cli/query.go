package cli

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"syreen/x/portfolio/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the portfolio module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(CmdQueryPortfolio())
	return cmd
}

func CmdQueryPortfolio() *cobra.Command {
	cmd := &cobra.Command{
		Use: "portfolio [address]", Short: "Query a portfolio by address", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil { return err }
			key := types.PortfolioKey(args[0])
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil || len(res) == 0 {
				result := map[string]interface{}{"address": args[0], "assets": []interface{}{}, "total_value": "0"}
				bz, _ := json.MarshalIndent(result, "", "  ")
				return clientCtx.PrintBytes(bz)
			}
			var parsed interface{}
			if err := json.Unmarshal(res, &parsed); err != nil { return clientCtx.PrintBytes(res) }
			bz, _ := json.MarshalIndent(parsed, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

