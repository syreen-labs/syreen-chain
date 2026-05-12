package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"syreen/x/payments/types"
)

// GetQueryCmd returns the query commands for the payments module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the payments module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryInvoice(),
		CmdQueryInvoicesByBooking(),
		CmdQueryInvoicesByPayer(),
		CmdQueryExchangeRate(),
		CmdQueryEarnings(),
	)

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query payments module parameters",
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

func CmdQueryInvoice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoice [invoice-id]",
		Short: "Query an invoice by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			key := types.InvoiceStoreKey(args[0])
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query invoice: %w", err)
			}

			if len(res) == 0 {
				return fmt.Errorf("invoice not found: %s", args[0])
			}

			var invoice types.Invoice
			if err := json.Unmarshal(res, &invoice); err != nil {
				return fmt.Errorf("failed to unmarshal invoice: %w", err)
			}

			bz, _ := json.MarshalIndent(invoice, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryInvoicesByBooking() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoices-by-booking [booking-id]",
		Short: "Query all invoices for a booking",
		Long:  `Lists all invoices associated with a booking ID. Uses gRPC query in production.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			result := map[string]interface{}{
				"booking_id": args[0],
				"note":       "Use gRPC query endpoint /syreen.payments.Query/InvoicesByBooking for full results",
			}

			bz, _ := json.MarshalIndent(result, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryInvoicesByPayer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoices-by-payer [payer-address]",
		Short: "Query all invoices for a payer",
		Long:  `Lists all invoices paid by a given address. Uses gRPC query in production.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			result := map[string]interface{}{
				"payer": args[0],
				"note":  "Use gRPC query endpoint /syreen.payments.Query/InvoicesByPayer for full results",
			}

			bz, _ := json.MarshalIndent(result, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryExchangeRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exchange-rate [from-denom] [to-denom]",
		Short: "Query exchange rate between two denominations",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			key := types.ExchangeRateStoreKey(args[0], args[1])
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query exchange rate: %w", err)
			}

			if len(res) == 0 {
				return fmt.Errorf("exchange rate not found for %s/%s", args[0], args[1])
			}

			var rate types.ExchangeRate
			if err := json.Unmarshal(res, &rate); err != nil {
				return fmt.Errorf("failed to unmarshal exchange rate: %w", err)
			}

			bz, _ := json.MarshalIndent(rate, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryEarnings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "earnings [address]",
		Short: "Query earnings for an address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			key := types.EarningsStoreKey(args[0])
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query earnings: %w", err)
			}

			if len(res) == 0 {
				result := map[string]interface{}{
					"address":         args[0],
					"available":       []interface{}{},
					"total_earned":    []interface{}{},
					"total_withdrawn": []interface{}{},
				}
				bz, _ := json.MarshalIndent(result, "", "  ")
				return clientCtx.PrintBytes(bz)
			}

			var earnings types.Earnings
			if err := json.Unmarshal(res, &earnings); err != nil {
				return fmt.Errorf("failed to unmarshal earnings: %w", err)
			}

			bz, _ := json.MarshalIndent(earnings, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
