package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"syreen/x/compute/types"
)

// GetQueryCmd returns the query commands for the compute module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the compute module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryCodeInfo(),
		CmdQueryContractInfo(),
		CmdQueryContractState(),
	)

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query compute module parameters",
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

func CmdQueryCodeInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "code-info [code-id]",
		Short: "Query code metadata by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			codeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid code-id: %w", err)
			}

			key := types.CodeInfoKey(codeID)
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query code info: %w", err)
			}

			if len(res) == 0 {
				return fmt.Errorf("code %d not found", codeID)
			}

			var codeInfo types.CodeInfo
			if err := json.Unmarshal(res, &codeInfo); err != nil {
				return fmt.Errorf("failed to unmarshal code info: %w", err)
			}

			bz, _ := json.MarshalIndent(codeInfo, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryContractInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract-info [contract-addr]",
		Short: "Query contract metadata by address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractAddr := args[0]
			key := types.ContractInfoKey(contractAddr)

			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query contract info: %w", err)
			}

			if len(res) == 0 {
				return fmt.Errorf("contract %s not found", contractAddr)
			}

			var contractInfo types.ContractInfo
			if err := json.Unmarshal(res, &contractInfo); err != nil {
				return fmt.Errorf("failed to unmarshal contract info: %w", err)
			}

			bz, _ := json.MarshalIndent(contractInfo, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryContractState() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract-state [contract-addr] [key]",
		Short: "Query raw contract state by key (hex-encoded)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractAddr := args[0]
			keyHex := args[1]

			keyBytes, err := hex.DecodeString(keyHex)
			if err != nil {
				// Try treating the key as a plain string
				keyBytes = []byte(keyHex)
			}

			storeKey := types.ContractStateKey(contractAddr, keyBytes)
			res, _, err := clientCtx.QueryStore(storeKey, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query contract state: %w", err)
			}

			if len(res) == 0 {
				return fmt.Errorf("no state found for contract %s at key %s", contractAddr, keyHex)
			}

			// Try to pretty-print as JSON if possible
			var jsonData interface{}
			if json.Unmarshal(res, &jsonData) == nil {
				bz, _ := json.MarshalIndent(jsonData, "", "  ")
				return clientCtx.PrintBytes(bz)
			}

			// Otherwise print as hex
			return clientCtx.PrintString(fmt.Sprintf("%x\n", res))
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
