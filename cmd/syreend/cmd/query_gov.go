package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// NewGovQueryCmd returns the query commands for the gov module.
func NewGovQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "gov",
		Short:                      "Querying commands for the governance module",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewQueryProposalCmd(),
		NewQueryProposalsCmd(),
		NewQueryGovParamsCmd(),
	)

	return cmd
}

// NewQueryProposalCmd returns the command to query a single proposal.
func NewQueryProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal [proposal-id]",
		Short: "Query details of a single proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := govtypes.NewQueryClient(clientCtx)

			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid proposal id: %w", err)
			}

			res, err := queryClient.Proposal(cmd.Context(), &govtypes.QueryProposalRequest{
				ProposalId: proposalID,
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

// NewQueryProposalsCmd returns the command to query all proposals.
func NewQueryProposalsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposals",
		Short: "Query all proposals with optional filters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := govtypes.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Proposals(cmd.Context(), &govtypes.QueryProposalsRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "proposals")

	return cmd
}

// NewQueryGovParamsCmd returns the command to query governance parameters.
func NewQueryGovParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the governance parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := govtypes.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &govtypes.QueryParamsRequest{
				ParamsType: "tallying",
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
