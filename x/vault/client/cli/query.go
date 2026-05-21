package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"syreen/x/vault/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the vault module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(CmdQueryVault(), CmdQueryVaults())
	return cmd
}

func CmdQueryVault() *cobra.Command {
	cmd := &cobra.Command{
		Use: "vault [vault-id]", Short: "Query a vault by ID", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil { return err }
			vaultID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid vault-id: %w", err) }
			key := types.VaultKey(vaultID)
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil || len(res) == 0 { return fmt.Errorf("vault %d not found", vaultID) }
			var parsed interface{}
			if err := json.Unmarshal(res, &parsed); err != nil { return clientCtx.PrintBytes(res) }
			bz, _ := json.MarshalIndent(parsed, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryVaults() *cobra.Command {
	cmd := &cobra.Command{
		Use: "vaults", Short: "Query all vaults", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil { return err }
			var vaults []interface{}
			for id := uint64(1); id <= 100; id++ {
				key := types.VaultKey(id)
				res, _, err := clientCtx.QueryStore(key, types.StoreKey)
				if err != nil || len(res) == 0 { break }
				var v interface{}
				if err := json.Unmarshal(res, &v); err == nil { vaults = append(vaults, v) }
			}
			if vaults == nil { vaults = []interface{}{} }
			bz, _ := json.MarshalIndent(vaults, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
