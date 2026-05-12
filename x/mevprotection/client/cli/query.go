package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"syreen/x/mevprotection/types"
)

// GetQueryCmd returns the query commands for the mevprotection module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the mevprotection module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryCommittedTx(),
		CmdQueryMEVRedistribution(),
	)

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query mevprotection module parameters",
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

func CmdQueryMEVRedistribution() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mev-redistribution",
		Short: "Query MEV redistribution stats (reward pool, total redistributed, split percentages)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Query the reward pool
			rewardPoolBz, _, _ := clientCtx.QueryStore([]byte("mev_reward_pool/usyreen"), types.StoreKey)
			rewardPool := "0"
			if len(rewardPoolBz) > 0 {
				var stored string
				if err := json.Unmarshal(rewardPoolBz, &stored); err == nil {
					rewardPool = stored
				}
			}

			// Query lifetime redistributed
			totalBz, _, _ := clientCtx.QueryStore([]byte("mev_total_redistributed/usyreen"), types.StoreKey)
			totalRedistributed := "0"
			if len(totalBz) > 0 {
				var stored string
				if err := json.Unmarshal(totalBz, &stored); err == nil {
					totalRedistributed = stored
				}
			}

			result := map[string]interface{}{
				"reward_pool":           rewardPool,
				"total_redistributed":   totalRedistributed,
				"lp_share_pct":          "60%",
				"staker_share_pct":      "40%",
				"distribution_interval": 100,
			}

			bz, _ := json.MarshalIndent(result, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryCommittedTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "committed-tx [hash-hex]",
		Short: "Query a committed transaction by its hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			txHash, err := hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("invalid hash hex: %w", err)
			}

			key := types.CommittedTxKey(txHash)
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query committed tx: %w", err)
			}

			if len(res) == 0 {
				return fmt.Errorf("committed tx with hash %s not found", args[0])
			}

			var committed types.CommittedTx
			if err := json.Unmarshal(res, &committed); err != nil {
				return fmt.Errorf("failed to unmarshal committed tx: %w", err)
			}

			bz, _ := json.MarshalIndent(committed, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
