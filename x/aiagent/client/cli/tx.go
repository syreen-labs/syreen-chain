package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"syreen/x/aiagent/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: types.ModuleName, Short: "AI Agent transaction subcommands",
		DisableFlagParsing: true, SuggestionsMinimumDistance: 2, RunE: client.ValidateCmd,
	}
	cmd.AddCommand(CmdCreateAgent(), CmdFundAgent(), CmdWithdrawAgentFunds(), CmdPauseAgent(), CmdResumeAgent())
	return cmd
}

func CmdCreateAgent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-agent [name] [strategy-type] [pool-id] [input-denom] [output-denom] [interval-blocks] [initial-funds]",
		Short: "Create a new AI trading agent", Args: cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			poolID, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil { return fmt.Errorf("invalid pool-id: %w", err) }
			interval, err := strconv.ParseInt(args[5], 10, 64)
			if err != nil { return fmt.Errorf("invalid interval-blocks: %w", err) }
			funds, ok := math.NewIntFromString(args[6])
			if !ok { return fmt.Errorf("invalid initial-funds: %s", args[6]) }
			msg := &types.MsgCreateAgent{
				Owner: clientCtx.GetFromAddress().String(), Name: args[0],
				StrategyType: types.StrategyType(args[1]),
				Config: types.AgentConfig{
					PoolID: poolID, InputDenom: args[3], OutputDenom: args[4],
					IntervalBlocks: interval, TradePercent: math.LegacyNewDecWithPrec(10, 2),
				},
				InitialFunds: funds,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdFundAgent() *cobra.Command {
	cmd := &cobra.Command{
		Use: "fund-agent [agent-id] [amount]", Short: "Fund an existing agent", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			aid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid agent-id: %w", err) }
			amount, ok := math.NewIntFromString(args[1])
			if !ok { return fmt.Errorf("invalid amount: %s", args[1]) }
			msg := &types.MsgFundAgent{Owner: clientCtx.GetFromAddress().String(), AgentID: aid, Amount: amount}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdWithdrawAgentFunds() *cobra.Command {
	cmd := &cobra.Command{
		Use: "withdraw-agent [agent-id] [amount]", Short: "Withdraw funds from an agent", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			aid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid agent-id: %w", err) }
			amount, ok := math.NewIntFromString(args[1])
			if !ok { return fmt.Errorf("invalid amount: %s", args[1]) }
			msg := &types.MsgWithdrawAgentFunds{Owner: clientCtx.GetFromAddress().String(), AgentID: aid, Amount: amount}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdPauseAgent() *cobra.Command {
	cmd := &cobra.Command{
		Use: "pause-agent [agent-id]", Short: "Pause an active agent", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			aid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid agent-id: %w", err) }
			msg := &types.MsgPauseAgent{Owner: clientCtx.GetFromAddress().String(), AgentID: aid}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdResumeAgent() *cobra.Command {
	cmd := &cobra.Command{
		Use: "resume-agent [agent-id]", Short: "Resume a paused agent", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			aid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid agent-id: %w", err) }
			msg := &types.MsgResumeAgent{Owner: clientCtx.GetFromAddress().String(), AgentID: aid}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
