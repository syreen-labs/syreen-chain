package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/abstractaccount/types"
)

// GetTxCmd returns the transaction commands for the abstractaccount module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Abstract account transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdCreateSmartAccount(),
		CmdCreateSessionKey(),
		CmdRevokeSessionKey(),
		CmdInitiateRecovery(),
		CmdApproveRecovery(),
		CmdExecuteRecovery(),
		CmdSponsorGas(),
	)

	return cmd
}

func CmdCreateSmartAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-smart-account",
		Short: "Create a new smart account",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress().String()
			msg := &types.MsgCreateSmartAccount{
				Sender:      sender,
				AccountType: types.AccountTypeSession,
				Owners:      []string{sender},
				Threshold:   1,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCreateSessionKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-session-key [grantee] [expiry-unix] [msg-types-csv] [max-amount] [denom]",
		Short: "Create a session key with scoped permissions",
		Long:  "Create a session key for a grantee with specific message type permissions and spending limits",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			grantee := args[0]

			expiryUnix, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid expiry unix timestamp: %w", err)
			}

			msgTypes := strings.Split(args[2], ",")

			maxAmount, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid max-amount: %w", err)
			}

			denom := args[4]

			// Build permissions from msg types
			var permissions []types.Permission
			for _, msgType := range msgTypes {
				permissions = append(permissions, types.Permission{
					MsgType:   strings.TrimSpace(msgType),
					MaxAmount: sdk.NewCoins(sdk.NewInt64Coin(denom, maxAmount)),
				})
			}

			now := time.Now()
			expiryTime := time.Unix(expiryUnix, 0)
			duration := expiryTime.Sub(now)
			if duration <= 0 {
				return fmt.Errorf("expiry must be in the future")
			}

			msg := &types.MsgCreateSessionKey{
				Granter:     clientCtx.GetFromAddress().String(),
				Grantee:     grantee,
				Permissions: permissions,
				Duration:    duration,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRevokeSessionKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-session-key [grantee]",
		Short: "Revoke a session key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgRevokeSessionKey{
				Granter:        clientCtx.GetFromAddress().String(),
				SessionKeyAddr: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdInitiateRecovery() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "initiate-recovery [account] [new-owners-csv]",
		Short: "Initiate social recovery for a smart account",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			account := args[0]
			newOwners := strings.Split(args[1], ",")
			for i, owner := range newOwners {
				newOwners[i] = strings.TrimSpace(owner)
			}

			msg := &types.MsgInitiateRecovery{
				Guardian:  clientCtx.GetFromAddress().String(),
				Account:   account,
				NewOwners: newOwners,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdApproveRecovery() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve-recovery [account]",
		Short: "Approve an active recovery request as a guardian",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgApproveRecovery{
				Guardian: clientCtx.GetFromAddress().String(),
				Account:  args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdExecuteRecovery() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute-recovery [account]",
		Short: "Execute an approved recovery after the delay period",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgExecuteRecovery{
				Sender:  clientCtx.GetFromAddress().String(),
				Account: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdSponsorGas() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sponsor-gas [sponsored-addr] [gas-limit] [max-amount] [denom] [expiry-unix]",
		Short: "Sponsor gas for another account",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sponsored := args[0]

			gasLimit, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid gas-limit: %w", err)
			}

			_, err = strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid max-amount: %w", err)
			}

			_ = args[3] // denom

			expiryUnix, err := strconv.ParseInt(args[4], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid expiry unix timestamp: %w", err)
			}

			now := time.Now()
			expiryTime := time.Unix(expiryUnix, 0)
			duration := expiryTime.Sub(now)
			if duration <= 0 {
				return fmt.Errorf("expiry must be in the future")
			}

			msg := &types.MsgSponsorGas{
				Sponsor:   clientCtx.GetFromAddress().String(),
				Sponsored: sponsored,
				GasLimit:  gasLimit,
				Duration:  duration,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
