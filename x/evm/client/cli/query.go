package cli

import (
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
			_ = clientCtx
			fmt.Println("EVM Module Parameters:")
			fmt.Printf("  evm_denom: %s\n", types.DefaultEVMDenom)
			fmt.Printf("  enable_create: %v\n", types.DefaultEnableCreate)
			fmt.Printf("  enable_call: %v\n", types.DefaultEnableCall)
			fmt.Printf("  chain_id: 79733 (Syreen EVM)\n")
			return nil
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
			if !common.IsHexAddress(args[0]) {
				return fmt.Errorf("invalid hex address: %s", args[0])
			}
			addr := common.HexToAddress(args[0])
			fmt.Printf("Code at %s: (query requires node connection)\n", addr.Hex())
			return nil
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
			if !common.IsHexAddress(args[0]) {
				return fmt.Errorf("invalid hex address: %s", args[0])
			}
			addr := common.HexToAddress(args[0])
			slot := common.HexToHash(args[1])
			fmt.Printf("Storage at %s slot %s: (query requires node connection)\n", addr.Hex(), slot.Hex())
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
