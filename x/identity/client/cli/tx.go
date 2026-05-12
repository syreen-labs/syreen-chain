package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"syreen/x/identity/types"
)

// GetTxCmd returns the transaction commands for the identity module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Identity module transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdRegisterIdentity(),
		CmdVerifyIdentity(),
		CmdRejectIdentity(),
		CmdRevokeIdentity(),
		CmdUpdateIdentity(),
		CmdRegisterVerifier(),
		CmdDeactivateVerifier(),
		CmdIncrementBookings(),
	)

	return cmd
}

func CmdRegisterIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [level] [document-hash] [nationality]",
		Short: "Register your identity for verification",
		Long: `Register an on-chain identity. Levels: basic, standard, enhanced.
document-hash is required for standard and enhanced levels.
nationality is optional (ISO 3166-1 alpha-2 code).`,
		Args: cobra.RangeArgs(1, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			docHash := ""
			nationality := ""
			if len(args) > 1 {
				docHash = args[1]
			}
			if len(args) > 2 {
				nationality = args[2]
			}

			msg := &types.MsgRegisterIdentity{
				Address:      clientCtx.GetFromAddress().String(),
				Level:        args[0],
				DocumentHash: docHash,
				Nationality:  nationality,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdVerifyIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify [address] [level]",
		Short: "Verify an identity (verifier only)",
		Long:  `Approve a pending identity at a given verification level. Only authorized verifiers can call this.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgVerifyIdentity{
				Verifier: clientCtx.GetFromAddress().String(),
				Address:  args[0],
				Level:    args[1],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRejectIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reject [address] [reason]",
		Short: "Reject an identity verification request (verifier only)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgRejectIdentity{
				Verifier: clientCtx.GetFromAddress().String(),
				Address:  args[0],
				Reason:   args[1],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRevokeIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke [address] [reason]",
		Short: "Revoke a verified identity (governance or original verifier only)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgRevokeIdentity{
				Authority: clientCtx.GetFromAddress().String(),
				Address:   args[0],
				Reason:    args[1],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdUpdateIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [document-hash] [nationality]",
		Short: "Update your identity (resets to pending for re-verification)",
		Long:  `Update document hash and/or nationality. Both are optional but at least one should be provided.`,
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			docHash := ""
			nationality := ""
			if len(args) > 0 {
				docHash = args[0]
			}
			if len(args) > 1 {
				nationality = args[1]
			}

			msg := &types.MsgUpdateIdentity{
				Address:      clientCtx.GetFromAddress().String(),
				DocumentHash: docHash,
				Nationality:  nationality,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRegisterVerifier() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-verifier [verifier-address] [name] [max-level]",
		Short: "Register a new identity verifier (governance only)",
		Long:  `Register an authorized identity verifier. max-level: basic, standard, or enhanced.`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgRegisterVerifier{
				Authority: clientCtx.GetFromAddress().String(),
				Verifier:  args[0],
				Name:      args[1],
				MaxLevel:  args[2],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdDeactivateVerifier() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deactivate-verifier [verifier-address]",
		Short: "Deactivate a verifier (governance only)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgDeactivateVerifier{
				Authority: clientCtx.GetFromAddress().String(),
				Verifier:  args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdIncrementBookings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "increment-bookings [address]",
		Short: "Increment booking count for an identity (authority only)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgIncrementBookings{
				Authority: clientCtx.GetFromAddress().String(),
				Address:   args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
