package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// NewStakingQueryCmd returns the query commands for the staking module.
func NewStakingQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "staking",
		Short:                      "Querying commands for the staking module",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewQueryValidatorsCmd(),
		NewQueryValidatorCmd(),
		NewQueryDelegationsCmd(),
		NewQueryDelegationCmd(),
		NewQueryUnbondingDelegationsCmd(),
		NewQueryStakingParamsCmd(),
		NewQueryStakingPoolCmd(),
	)

	return cmd
}

// NewQueryValidatorsCmd returns the command to query all validators.
func NewQueryValidatorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validators",
		Short: "Query for all validators",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := stakingtypes.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			status, _ := cmd.Flags().GetString("status")

			res, err := queryClient.Validators(cmd.Context(), &stakingtypes.QueryValidatorsRequest{
				Status:     status,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("status", "", "Filter validators by status (bonded|unbonding|unbonded)")
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "validators")

	return cmd
}

// NewQueryValidatorCmd returns the command to query a single validator.
func NewQueryValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator [validator-addr]",
		Short: "Query a validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := stakingtypes.NewQueryClient(clientCtx)

			res, err := queryClient.Validator(cmd.Context(), &stakingtypes.QueryValidatorRequest{
				ValidatorAddr: args[0],
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

// NewQueryDelegationsCmd returns the command to query delegations for a delegator.
func NewQueryDelegationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegations [delegator-addr]",
		Short: "Query all delegations for a delegator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := stakingtypes.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.DelegatorDelegations(cmd.Context(), &stakingtypes.QueryDelegatorDelegationsRequest{
				DelegatorAddr: args[0],
				Pagination:    pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "delegations")

	return cmd
}

// NewQueryDelegationCmd returns the command to query a specific delegation.
func NewQueryDelegationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegation [delegator-addr] [validator-addr]",
		Short: "Query a specific delegation",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := stakingtypes.NewQueryClient(clientCtx)

			res, err := queryClient.Delegation(cmd.Context(), &stakingtypes.QueryDelegationRequest{
				DelegatorAddr: args[0],
				ValidatorAddr: args[1],
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

// NewQueryUnbondingDelegationsCmd returns the command to query unbonding delegations.
func NewQueryUnbondingDelegationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbonding-delegations [delegator-addr]",
		Short: "Query all unbonding delegations for a delegator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := stakingtypes.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.DelegatorUnbondingDelegations(cmd.Context(), &stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
				DelegatorAddr: args[0],
				Pagination:    pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "unbonding delegations")

	return cmd
}

// NewQueryStakingParamsCmd returns the command to query staking parameters.
func NewQueryStakingParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current staking parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := stakingtypes.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &stakingtypes.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryStakingPoolCmd returns the command to query the staking pool.
func NewQueryStakingPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool",
		Short: "Query the current staking pool values",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := stakingtypes.NewQueryClient(clientCtx)

			res, err := queryClient.Pool(cmd.Context(), &stakingtypes.QueryPoolRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
