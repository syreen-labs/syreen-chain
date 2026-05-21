package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"syreen/x/clmm/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the CLMM module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryPool(),
		CmdQueryPools(),
		CmdQueryPosition(),
	)

	return cmd
}

func CmdQueryPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool [pool-id]",
		Short: "Query a CL pool by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool-id: %w", err)
			}
			key := types.CLPoolKey(poolID)
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil || len(res) == 0 {
				return fmt.Errorf("CL pool %d not found", poolID)
			}
			var parsed interface{}
			if err := json.Unmarshal(res, &parsed); err != nil {
				return clientCtx.PrintBytes(res)
			}
			bz, _ := json.MarshalIndent(parsed, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryPools() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pools",
		Short: "Query all CL pools",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			var pools []interface{}
			for id := uint64(1); id <= 100; id++ {
				key := types.CLPoolKey(id)
				res, _, err := clientCtx.QueryStore(key, types.StoreKey)
				if err != nil || len(res) == 0 {
					break
				}
				var pool interface{}
				if err := json.Unmarshal(res, &pool); err == nil {
					pools = append(pools, pool)
				}
			}
			if pools == nil {
				pools = []interface{}{}
			}
			bz, _ := json.MarshalIndent(pools, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryPosition() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "position [position-id]",
		Short: "Query a CL position by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			posID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid position-id: %w", err)
			}
			key := types.CLPositionKey(posID)
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil || len(res) == 0 {
				return fmt.Errorf("position %d not found", posID)
			}
			var parsed interface{}
			if err := json.Unmarshal(res, &parsed); err != nil {
				return clientCtx.PrintBytes(res)
			}
			bz, _ := json.MarshalIndent(parsed, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
