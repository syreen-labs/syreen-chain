package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// NewDistributionQueryCmd returns the query commands for the distribution module.
func NewDistributionQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "distribution",
		Aliases:                    []string{"distr"},
		Short:                      "Querying commands for the distribution module",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewQueryDelegatorRewardsCmd(),
		NewQueryDistributionParamsCmd(),
		NewQueryCommunityPoolCmd(),
	)

	return cmd
}

// NewQueryDelegatorRewardsCmd returns the command to query delegator rewards.
func NewQueryDelegatorRewardsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards [delegator-addr] [validator-addr]",
		Short: "Query all distribution delegator rewards or rewards from a particular validator",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := distrtypes.NewQueryClient(clientCtx)

			if len(args) == 2 {
				res, err := queryClient.DelegationRewards(cmd.Context(), &distrtypes.QueryDelegationRewardsRequest{
					DelegatorAddress: args[0],
					ValidatorAddress: args[1],
				})
				if err != nil {
					return err
				}
				return clientCtx.PrintProto(res)
			}

			res, err := queryClient.DelegationTotalRewards(cmd.Context(), &distrtypes.QueryDelegationTotalRewardsRequest{
				DelegatorAddress: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryDistributionParamsCmd returns the command to query distribution parameters.
func NewQueryDistributionParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query distribution params",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := distrtypes.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &distrtypes.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryCommunityPoolCmd returns the command to query the community pool.
func NewQueryCommunityPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "community-pool",
		Short: "Query the community pool coins",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := distrtypes.NewQueryClient(clientCtx)

			res, err := queryClient.CommunityPool(cmd.Context(), &distrtypes.QueryCommunityPoolRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
