package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// NewBankQueryCmd returns the query commands for the bank module (SDK v0.50 moved these to AutoCLI).
func NewBankQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "bank",
		Short:                      "Querying commands for the bank module",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewBankBalancesCmd(),
		NewBankTotalSupplyCmd(),
		NewBankDenomsMetadataCmd(),
	)

	return cmd
}

// NewBankBalancesCmd returns a command to query account balances.
func NewBankBalancesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balances [address]",
		Short: "Query for account balances by address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := banktypes.NewQueryClient(clientCtx)

			denom, _ := cmd.Flags().GetString("denom")
			if denom != "" {
				res, err := queryClient.Balance(cmd.Context(), &banktypes.QueryBalanceRequest{
					Address: args[0],
					Denom:   denom,
				})
				if err != nil {
					return err
				}
				return clientCtx.PrintProto(res)
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.AllBalances(cmd.Context(), &banktypes.QueryAllBalancesRequest{
				Address:    args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("denom", "", "Query a specific denomination")
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all balances")

	return cmd
}

// NewBankTotalSupplyCmd returns a command to query total supply.
func NewBankTotalSupplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-supply",
		Short: "Query the total supply of coins on the chain",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := banktypes.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.TotalSupply(cmd.Context(), &banktypes.QueryTotalSupplyRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "total supply")

	return cmd
}

// NewBankDenomsMetadataCmd returns a command to query denomination metadata.
func NewBankDenomsMetadataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "denoms-metadata",
		Short: "Query the client metadata for coin denominations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := banktypes.NewQueryClient(clientCtx)

			denom, _ := cmd.Flags().GetString("denom")
			if denom != "" {
				res, err := queryClient.DenomMetadata(cmd.Context(), &banktypes.QueryDenomMetadataRequest{
					Denom: denom,
				})
				if err != nil {
					return err
				}
				return clientCtx.PrintProto(res)
			}

			res, err := queryClient.DenomsMetadata(cmd.Context(), &banktypes.QueryDenomsMetadataRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("denom", "", "Query metadata for a specific denomination")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

