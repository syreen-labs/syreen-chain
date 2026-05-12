package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"syreen/x/tokenfactory/types"
)

// GetQueryCmd returns the query commands for the tokenfactory module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the tokenfactory module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryDenomInfo(),
	)

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query tokenfactory module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryStore([]byte("params"), types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query params: %w", err)
			}

			if len(res) == 0 {
				params := types.DefaultParams()
				bz, _ := json.MarshalIndent(params, "", "  ")
				return clientCtx.PrintBytes(bz)
			}

			var params types.Params
			if err := json.Unmarshal(res, &params); err != nil {
				return fmt.Errorf("failed to unmarshal params: %w", err)
			}

			bz, _ := json.MarshalIndent(params, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryDenomInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "denom-info [denom]",
		Short: "Query denom authority metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			denom := args[0]
			key := types.DenomAuthorityKey(denom)

			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query denom info: %w", err)
			}

			if len(res) == 0 {
				return fmt.Errorf("denom %s not found", denom)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(res, &result); err != nil {
				return clientCtx.PrintBytes(res)
			}

			bz, _ := json.MarshalIndent(result, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
