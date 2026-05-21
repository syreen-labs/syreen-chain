package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"syreen/x/perps/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: types.ModuleName, Short: "Perps transaction subcommands",
		DisableFlagParsing: true, SuggestionsMinimumDistance: 2, RunE: client.ValidateCmd,
	}
	cmd.AddCommand(CmdOpenPosition(), CmdClosePosition(), CmdAddMargin(), CmdRemoveMargin(), CmdCreateMarket())
	return cmd
}

func CmdOpenPosition() *cobra.Command {
	cmd := &cobra.Command{
		Use: "open-position [market-id] [side] [margin] [leverage]", Short: "Open a perpetual position", Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			marketID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid market-id: %w", err) }
			margin, ok := math.NewIntFromString(args[2])
			if !ok { return fmt.Errorf("invalid margin: %s", args[2]) }
			leverage, err := math.LegacyNewDecFromStr(args[3])
			if err != nil { return fmt.Errorf("invalid leverage: %w", err) }
			msg := &types.MsgOpenPosition{Sender: clientCtx.GetFromAddress().String(), MarketID: marketID, Side: types.Side(args[1]), Margin: margin, Leverage: leverage}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdClosePosition() *cobra.Command {
	cmd := &cobra.Command{
		Use: "close-position [market-id]", Short: "Close a perpetual position", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			marketID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid market-id: %w", err) }
			msg := &types.MsgClosePosition{Sender: clientCtx.GetFromAddress().String(), MarketID: marketID}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdAddMargin() *cobra.Command {
	cmd := &cobra.Command{
		Use: "add-margin [market-id] [amount]", Short: "Add margin to a position", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			marketID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid market-id: %w", err) }
			amount, ok := math.NewIntFromString(args[1])
			if !ok { return fmt.Errorf("invalid amount: %s", args[1]) }
			msg := &types.MsgAddMargin{Sender: clientCtx.GetFromAddress().String(), MarketID: marketID, Amount: amount}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRemoveMargin() *cobra.Command {
	cmd := &cobra.Command{
		Use: "remove-margin [market-id] [amount]", Short: "Remove margin from a position", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			marketID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid market-id: %w", err) }
			amount, ok := math.NewIntFromString(args[1])
			if !ok { return fmt.Errorf("invalid amount: %s", args[1]) }
			msg := &types.MsgRemoveMargin{Sender: clientCtx.GetFromAddress().String(), MarketID: marketID, Amount: amount}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCreateMarket() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create-market [base-denom] [quote-denom] [pool-id] [max-leverage]", Short: "Create a perps market (governance only)", Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			poolID, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil { return fmt.Errorf("invalid pool-id: %w", err) }
			maxLev, err := math.LegacyNewDecFromStr(args[3])
			if err != nil { return fmt.Errorf("invalid max-leverage: %w", err) }
			msg := &types.MsgCreateMarket{Authority: clientCtx.GetFromAddress().String(), BaseDenom: args[0], QuoteDenom: args[1], PoolID: poolID, MaxLeverage: maxLev}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
