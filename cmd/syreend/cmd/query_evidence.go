package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	evidencetypes "cosmossdk.io/x/evidence/types"
)

// NewEvidenceQueryCmd returns the query commands for the evidence module (SDK v0.50 moved these to AutoCLI).
func NewEvidenceQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "evidence",
		Short:                      "Querying commands for the evidence module",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewQueryEvidenceCmd(),
		NewQueryAllEvidenceCmd(),
	)

	return cmd
}

// NewQueryEvidenceCmd returns a command to query evidence by hash.
func NewQueryEvidenceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [hash]",
		Short: "Query evidence by hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := evidencetypes.NewQueryClient(clientCtx)

			res, err := queryClient.Evidence(cmd.Context(), &evidencetypes.QueryEvidenceRequest{
				Hash: args[0],
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

// NewQueryAllEvidenceCmd returns a command to query all evidence.
func NewQueryAllEvidenceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Query all submitted evidence (paginated)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := evidencetypes.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.AllEvidence(cmd.Context(), &evidencetypes.QueryAllEvidenceRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "evidence")

	return cmd
}
