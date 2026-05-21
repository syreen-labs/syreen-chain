package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"syreen/x/portfolio/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Portfolio transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		CmdRecordTrade(),
		CmdCreateCompetition(),
		CmdJoinCompetition(),
		CmdEndCompetition(),
		CmdUpdatePortfolio(),
	)
	return cmd
}

func CmdRecordTrade() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record-trade [trade-type] [denom] [amount] [price] [pnl]",
		Short: "Record a trade event for portfolio tracking",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			amount, ok := math.NewIntFromString(args[2])
			if !ok { return fmt.Errorf("invalid amount: %s", args[2]) }
			price, ok := math.NewIntFromString(args[3])
			if !ok { return fmt.Errorf("invalid price: %s", args[3]) }
			pnl, ok := math.NewIntFromString(args[4])
			if !ok { return fmt.Errorf("invalid pnl: %s", args[4]) }
			msg := &types.MsgRecordTrade{
				Sender: clientCtx.GetFromAddress().String(), TradeType: args[0],
				Denom: args[1], Amount: amount, Price: price, PnL: pnl,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCreateCompetition() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-competition [name] [start-block] [end-block] [prize-denom] [prize-pool] [entry-fee] [max-participants]",
		Short: "Create a trading competition",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			startBlock, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil { return fmt.Errorf("invalid start-block: %w", err) }
			endBlock, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil { return fmt.Errorf("invalid end-block: %w", err) }
			prizePool, ok := math.NewIntFromString(args[4])
			if !ok { return fmt.Errorf("invalid prize-pool: %s", args[4]) }
			entryFee, ok := math.NewIntFromString(args[5])
			if !ok { return fmt.Errorf("invalid entry-fee: %s", args[5]) }
			maxPart, err := strconv.ParseUint(args[6], 10, 64)
			if err != nil { return fmt.Errorf("invalid max-participants: %w", err) }
			msg := &types.MsgCreateCompetition{
				Creator: clientCtx.GetFromAddress().String(), Name: args[0],
				StartBlock: startBlock, EndBlock: endBlock, PrizeDenom: args[3],
				PrizePool: prizePool, EntryFee: entryFee, MaxParticipants: maxPart,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdJoinCompetition() *cobra.Command {
	cmd := &cobra.Command{
		Use: "join-competition [competition-id]", Short: "Join a trading competition", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			compID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid competition-id: %w", err) }
			msg := &types.MsgJoinCompetition{Sender: clientCtx.GetFromAddress().String(), CompetitionID: compID}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdEndCompetition() *cobra.Command {
	cmd := &cobra.Command{
		Use: "end-competition [competition-id]", Short: "End a competition and distribute prizes", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			compID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid competition-id: %w", err) }
			msg := &types.MsgEndCompetition{Authority: clientCtx.GetFromAddress().String(), CompetitionID: compID}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdUpdatePortfolio() *cobra.Command {
	cmd := &cobra.Command{
		Use: "update", Short: "Trigger a portfolio recalculation", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			msg := &types.MsgUpdatePortfolio{Sender: clientCtx.GetFromAddress().String()}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
