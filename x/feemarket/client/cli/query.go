package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"syreen/x/feemarket/types"
)

// GetQueryCmd returns the query commands for the feemarket module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the feemarket module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryBaseFee(),
		CmdQueryFeeLanes(),
	)

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query feemarket module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryStore([]byte(types.ParamsKey), types.StoreKey)
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

func CmdQueryBaseFee() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "base-fee",
		Short: "Query the current base fee",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryStore([]byte(types.BaseFeeKey), types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query base fee: %w", err)
			}

			if len(res) == 0 {
				state := types.DefaultFeeState()
				result := map[string]string{"base_fee": state.BaseFee.String()}
				bz, _ := json.MarshalIndent(result, "", "  ")
				return clientCtx.PrintBytes(bz)
			}

			var state types.FeeState
			if err := json.Unmarshal(res, &state); err != nil {
				return fmt.Errorf("failed to unmarshal fee state: %w", err)
			}

			result := map[string]string{"base_fee": state.BaseFee.String()}
			bz, _ := json.MarshalIndent(result, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryFeeLanes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fee-lanes",
		Short: "Query all fee lanes and their multipliers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Fee lanes are stored in params
			res, _, err := clientCtx.QueryStore([]byte(types.ParamsKey), types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query params: %w", err)
			}

			var params types.Params
			if len(res) == 0 {
				params = types.DefaultParams()
			} else {
				if err := json.Unmarshal(res, &params); err != nil {
					return fmt.Errorf("failed to unmarshal params: %w", err)
				}
			}

			bz, _ := json.MarshalIndent(params.FeeLanes, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
