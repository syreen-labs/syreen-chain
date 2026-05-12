package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// NewSlashingQueryCmd returns the query commands for the slashing module.
func NewSlashingQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "slashing",
		Short:                      "Querying commands for the slashing module",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewQuerySigningInfosCmd(),
		NewQuerySlashingParamsCmd(),
	)

	return cmd
}

// NewQuerySigningInfosCmd returns a command to query all signing infos.
func NewQuerySigningInfosCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signing-infos",
		Short: "Query signing info of all validators",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := slashingtypes.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.SigningInfos(cmd.Context(), &slashingtypes.QuerySigningInfosRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "signing infos")

	return cmd
}

// NewQuerySlashingParamsCmd returns a command to query slashing parameters.
func NewQuerySlashingParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current slashing parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := slashingtypes.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &slashingtypes.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
