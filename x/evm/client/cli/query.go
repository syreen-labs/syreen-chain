package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"syreen/x/evm/types"
)

// GetQueryCmd returns the query commands for the EVM module
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "EVM query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryCode(),
		CmdQueryStorage(),
	)

	return cmd
}

// CmdQueryParams queries the EVM module parameters
func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the EVM module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryStore(types.KeyParams(), types.StoreKey)
			if err != nil || len(res) == 0 {
				result := map[string]interface{}{
					"evm_denom":     types.DefaultEVMDenom,
					"enable_create": types.DefaultEnableCreate,
					"enable_call":   types.DefaultEnableCall,
					"chain_id":      79733,
				}
				bz, _ := json.MarshalIndent(result, "", "  ")
				return clientCtx.PrintBytes(bz)
			}

			var parsed interface{}
			if err := json.Unmarshal(res, &parsed); err != nil {
				return clientCtx.PrintBytes(res)
			}
			bz, _ := json.MarshalIndent(parsed, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryCode queries the bytecode at an EVM address
func CmdQueryCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "code [hex-address]",
		Short: "Query the EVM bytecode at an address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			if !common.IsHexAddress(args[0]) {
				return fmt.Errorf("invalid hex address: %s", args[0])
			}
			addr := common.HexToAddress(args[0])

			// Query EVM code from the store using the code key prefix
			key := types.KeyCode(addr.Bytes())
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil || len(res) == 0 {
				result := map[string]interface{}{
					"address": addr.Hex(),
					"code":    "",
					"size":    0,
				}
				bz, _ := json.MarshalIndent(result, "", "  ")
				return clientCtx.PrintBytes(bz)
			}

			result := map[string]interface{}{
				"address": addr.Hex(),
				"code":    hex.EncodeToString(res),
				"size":    len(res),
			}
			bz, _ := json.MarshalIndent(result, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryStorage queries a storage slot at an EVM address
func CmdQueryStorage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage [hex-address] [slot-hex]",
		Short: "Query a storage slot at an EVM address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			if !common.IsHexAddress(args[0]) {
				return fmt.Errorf("invalid hex address: %s", args[0])
			}
			addr := common.HexToAddress(args[0])
			slot := common.HexToHash(args[1])

			key := types.KeyStorage(addr.Bytes(), slot.Bytes())
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil || len(res) == 0 {
				result := map[string]interface{}{
					"address": addr.Hex(),
					"slot":    slot.Hex(),
					"value":   common.Hash{}.Hex(),
				}
				bz, _ := json.MarshalIndent(result, "", "  ")
				return clientCtx.PrintBytes(bz)
			}

			result := map[string]interface{}{
				"address": addr.Hex(),
				"slot":    slot.Hex(),
				"value":   common.BytesToHash(res).Hex(),
			}
			bz, _ := json.MarshalIndent(result, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
