package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/payments/types"
)

// GetTxCmd returns the transaction commands for the payments module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Payments module transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdCreateInvoice(),
		CmdPayInvoice(),
		CmdRefundPayment(),
		CmdSetExchangeRate(),
		CmdWithdrawEarnings(),
		CmdCancelInvoice(),
	)

	return cmd
}

func CmdCreateInvoice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-invoice [booking-id] [payee] [amount] [description] [currency]",
		Short: "Create a new payment invoice for a booking",
		Long:  `Create a payment invoice with automatic fee calculation. Amount format: 1000000usyreen`,
		Args:  cobra.RangeArgs(3, 5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			description := ""
			currency := ""
			if len(args) > 3 {
				description = args[3]
			}
			if len(args) > 4 {
				currency = args[4]
			}

			msg := &types.MsgCreateInvoice{
				Creator:     clientCtx.GetFromAddress().String(),
				BookingID:   args[0],
				Payee:       args[1],
				Amount:      amount,
				Description: description,
				Currency:    currency,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdPayInvoice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pay-invoice [invoice-id] [payment-method]",
		Short: "Pay an existing invoice",
		Long:  `Pay an invoice. Payment methods: crypto, usdc, syr`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgPayInvoice{
				Payer:         clientCtx.GetFromAddress().String(),
				InvoiceID:     args[0],
				PaymentMethod: args[1],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRefundPayment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refund [invoice-id] [reason]",
		Short: "Refund a paid invoice (authority only)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgRefundPayment{
				Authority: clientCtx.GetFromAddress().String(),
				InvoiceID: args[0],
				Reason:    args[1],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdSetExchangeRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-rate [from-denom] [to-denom] [rate]",
		Short: "Set exchange rate between two denominations (authority only)",
		Long:  `Set the exchange rate. Example: set-rate usyreen uusdc 0.001 (1 SYR = 0.001 USDC)`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgSetExchangeRate{
				Authority: clientCtx.GetFromAddress().String(),
				FromDenom: args[0],
				ToDenom:   args[1],
				Rate:      args[2],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdWithdrawEarnings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [amount]",
		Short: "Withdraw your available earnings",
		Long:  `Withdraw accumulated earnings. Amount format: 1000000usyreen`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := &types.MsgWithdrawEarnings{
				Address: clientCtx.GetFromAddress().String(),
				Amount:  amount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCancelInvoice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-invoice [invoice-id]",
		Short: "Cancel a pending invoice (creator only)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgCancelInvoice{
				Creator:   clientCtx.GetFromAddress().String(),
				InvoiceID: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
