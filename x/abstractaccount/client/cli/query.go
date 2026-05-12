package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"syreen/x/abstractaccount/types"
)

// GetQueryCmd returns the query commands for the abstractaccount module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the abstractaccount module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQuerySmartAccount(),
		CmdQuerySessionKeys(),
	)

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query abstractaccount module parameters",
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

func CmdQuerySmartAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "smart-account [address]",
		Short: "Query a smart account by address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			address := args[0]
			key := types.SmartAccountKey(address)

			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query smart account: %w", err)
			}

			if len(res) == 0 {
				return fmt.Errorf("smart account %s not found", address)
			}

			var account types.SmartAccount
			if err := json.Unmarshal(res, &account); err != nil {
				return fmt.Errorf("failed to unmarshal smart account: %w", err)
			}

			bz, _ := json.MarshalIndent(account, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQuerySessionKeys() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session-keys [account]",
		Short: "Query session keys for an account (provide granter address)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			granter := args[0]
			prefix := types.SessionKeysByGranterPrefix(granter)

			// Use ABCI query with a subspace iteration
			res, _, err := clientCtx.QueryStore(prefix, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query session keys: %w", err)
			}

			if len(res) == 0 {
				return clientCtx.PrintBytes([]byte("[]"))
			}

			// Try to unmarshal as a single session key first
			var sessionKey types.SessionKey
			if err := json.Unmarshal(res, &sessionKey); err != nil {
				// If not a single key, print raw
				return clientCtx.PrintBytes(res)
			}

			bz, _ := json.MarshalIndent(sessionKey, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
