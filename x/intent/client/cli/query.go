package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"syreen/x/intent/types"
)

// GetQueryCmd returns the query commands for the intent module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the intent module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryIntent(),
		CmdQuerySolver(),
		CmdQuerySolutions(),
	)

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query intent module parameters",
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

func CmdQueryIntent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "intent [id]",
		Short: "Query an intent by its ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			intentID := args[0]
			key := types.IntentKey(intentID)

			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query intent: %w", err)
			}

			if len(res) == 0 {
				return fmt.Errorf("intent %s not found", intentID)
			}

			var intent types.Intent
			if err := json.Unmarshal(res, &intent); err != nil {
				return fmt.Errorf("failed to unmarshal intent: %w", err)
			}

			bz, _ := json.MarshalIndent(intent, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQuerySolver() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "solver [address]",
		Short: "Query a solver by address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			address := args[0]
			key := types.SolverKey(address)

			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query solver: %w", err)
			}

			if len(res) == 0 {
				return fmt.Errorf("solver %s not found", address)
			}

			var solver types.Solver
			if err := json.Unmarshal(res, &solver); err != nil {
				return fmt.Errorf("failed to unmarshal solver: %w", err)
			}

			bz, _ := json.MarshalIndent(solver, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQuerySolutions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "solutions [intent-id]",
		Short: "Query all solutions for an intent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			intentID := args[0]
			prefix := types.SolutionsByIntentPrefix(intentID)

			res, _, err := clientCtx.QueryStore(prefix, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query solutions: %w", err)
			}

			if len(res) == 0 {
				return clientCtx.PrintBytes([]byte("[]"))
			}

			// Try to unmarshal as a single solution first
			var solution types.Solution
			if err := json.Unmarshal(res, &solution); err != nil {
				return clientCtx.PrintBytes(res)
			}

			bz, _ := json.MarshalIndent(solution, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
