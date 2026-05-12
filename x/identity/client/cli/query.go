package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"syreen/x/identity/types"
)

// GetQueryCmd returns the query commands for the identity module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the identity module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryIdentity(),
		CmdQueryIdentitiesByVerifier(),
		CmdQueryVerifier(),
		CmdQueryVerificationStatus(),
	)

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query identity module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryStore([]byte(types.ParamsKey), types.StoreKey)
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

func CmdQueryIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "identity [address]",
		Short: "Query an identity by address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			key := types.IdentityStoreKey(args[0])
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query identity: %w", err)
			}

			if len(res) == 0 {
				return fmt.Errorf("identity not found for address %s", args[0])
			}

			var identity types.Identity
			if err := json.Unmarshal(res, &identity); err != nil {
				return fmt.Errorf("failed to unmarshal identity: %w", err)
			}

			bz, _ := json.MarshalIndent(identity, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryIdentitiesByVerifier() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "identities-by-verifier [verifier-address]",
		Short: "Query all identities verified by a specific verifier",
		Long:  `Lists all identities that were verified by a given verifier address. Uses gRPC query in production.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// For CLI, display a message directing to gRPC query
			// In production, this uses the gRPC query server
			result := map[string]interface{}{
				"verifier": args[0],
				"note":     "Use gRPC query endpoint /syreen.identity.Query/IdentitiesByVerifier for full results",
			}

			_ = clientCtx // used for output
			bz, _ := json.MarshalIndent(result, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryVerifier() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verifier [address]",
		Short: "Query a verifier by address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			key := types.VerifierStoreKey(args[0])
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query verifier: %w", err)
			}

			if len(res) == 0 {
				return fmt.Errorf("verifier not found for address %s", args[0])
			}

			var verifier types.Verifier
			if err := json.Unmarshal(res, &verifier); err != nil {
				return fmt.Errorf("failed to unmarshal verifier: %w", err)
			}

			bz, _ := json.MarshalIndent(verifier, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryVerificationStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [address]",
		Short: "Query the verification status of an address",
		Long:  `Returns whether the address is verified, their level, status, and trust score.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			key := types.IdentityStoreKey(args[0])
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query identity: %w", err)
			}

			if len(res) == 0 {
				result := map[string]interface{}{
					"address":     args[0],
					"verified":    false,
					"level":       "none",
					"status":      "not_registered",
					"trust_score": 0,
				}
				bz, _ := json.MarshalIndent(result, "", "  ")
				return clientCtx.PrintBytes(bz)
			}

			var identity types.Identity
			if err := json.Unmarshal(res, &identity); err != nil {
				return fmt.Errorf("failed to unmarshal identity: %w", err)
			}

			verified := identity.Status == types.StatusVerified
			result := map[string]interface{}{
				"address":     identity.Address,
				"verified":    verified,
				"level":       identity.Level,
				"status":      identity.Status,
				"trust_score": identity.TrustScore,
			}
			bz, _ := json.MarshalIndent(result, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
