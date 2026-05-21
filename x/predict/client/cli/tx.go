package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"syreen/x/predict/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: types.ModuleName, Short: "Predict transaction subcommands",
		DisableFlagParsing: true, SuggestionsMinimumDistance: 2, RunE: client.ValidateCmd,
	}
	cmd.AddCommand(CmdCreateMarket(), CmdBuyShares(), CmdSellShares(), CmdResolveMarket(), CmdClaimWinnings())
	return cmd
}

func CmdCreateMarket() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create-market [question] [resolver] [quote-denom] [resolution-block] [initial-liquidity]", Short: "Create a prediction market", Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			rb, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil { return fmt.Errorf("invalid resolution-block: %w", err) }
			il, ok := math.NewIntFromString(args[4])
			if !ok { return fmt.Errorf("invalid initial-liquidity: %s", args[4]) }
			msg := &types.MsgCreateMarket{Creator: clientCtx.GetFromAddress().String(), Question: args[0], Resolver: args[1], QuoteDenom: args[2], ResolutionBlock: rb, InitialLiquidity: il}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdBuyShares() *cobra.Command {
	cmd := &cobra.Command{
		Use: "buy-shares [market-id] [outcome] [amount]", Short: "Buy YES or NO shares", Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			mid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid market-id: %w", err) }
			amount, ok := math.NewIntFromString(args[2])
			if !ok { return fmt.Errorf("invalid amount: %s", args[2]) }
			msg := &types.MsgBuyShares{Sender: clientCtx.GetFromAddress().String(), MarketID: mid, Outcome: args[1], Amount: amount}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdSellShares() *cobra.Command {
	cmd := &cobra.Command{
		Use: "sell-shares [market-id] [outcome] [shares]", Short: "Sell YES or NO shares", Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			mid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid market-id: %w", err) }
			shares, ok := math.NewIntFromString(args[2])
			if !ok { return fmt.Errorf("invalid shares: %s", args[2]) }
			msg := &types.MsgSellShares{Sender: clientCtx.GetFromAddress().String(), MarketID: mid, Outcome: args[1], Shares: shares}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdResolveMarket() *cobra.Command {
	cmd := &cobra.Command{
		Use: "resolve-market [market-id] [outcome]", Short: "Resolve a prediction market (yes/no/void)", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			mid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid market-id: %w", err) }
			msg := &types.MsgResolveMarket{Resolver: clientCtx.GetFromAddress().String(), MarketID: mid, Outcome: args[1]}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdClaimWinnings() *cobra.Command {
	cmd := &cobra.Command{
		Use: "claim-winnings [market-id]", Short: "Claim payout from a resolved market", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			mid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid market-id: %w", err) }
			msg := &types.MsgClaimWinnings{Sender: clientCtx.GetFromAddress().String(), MarketID: mid}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
