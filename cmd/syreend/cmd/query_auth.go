package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// NewAuthQueryCmd returns the query commands for the auth module (SDK v0.50 moved these to AutoCLI).
func NewAuthQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "auth",
		Short:                      "Querying commands for the auth module",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewQueryAccountCmd(),
		NewQueryAccountsCmd(),
		NewQueryAuthParamsCmd(),
		NewQueryModuleAccountsCmd(),
	)

	return cmd
}

// NewQueryAccountCmd returns a command to query an account by address.
func NewQueryAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account [address]",
		Short: "Query account details by address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := authtypes.NewQueryClient(clientCtx)

			res, err := queryClient.Account(cmd.Context(), &authtypes.QueryAccountRequest{
				Address: args[0],
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

// NewQueryAccountsCmd returns a command to query all accounts.
func NewQueryAccountsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "Query all accounts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := authtypes.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Accounts(cmd.Context(), &authtypes.QueryAccountsRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "accounts")

	return cmd
}

// NewQueryAuthParamsCmd returns a command to query auth parameters.
func NewQueryAuthParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current auth parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := authtypes.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &authtypes.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryModuleAccountsCmd returns a command to query all module accounts.
func NewQueryModuleAccountsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module-accounts",
		Short: "Query all module accounts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := authtypes.NewQueryClient(clientCtx)

			res, err := queryClient.ModuleAccounts(cmd.Context(), &authtypes.QueryModuleAccountsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
