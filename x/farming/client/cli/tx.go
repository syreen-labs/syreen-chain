package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"syreen/x/farming/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Farming transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(CmdStake(), CmdUnstake(), CmdClaimReward(), CmdCreateFarm(), CmdUpdateFarm())
	return cmd
}

func CmdStake() *cobra.Command {
	cmd := &cobra.Command{
		Use: "stake [pool-id] [amount]", Short: "Stake LP tokens into a farming pool", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid pool-id: %w", err) }
			amount, ok := math.NewIntFromString(args[1])
			if !ok { return fmt.Errorf("invalid amount: %s", args[1]) }
			msg := &types.MsgStake{Sender: clientCtx.GetFromAddress().String(), PoolID: poolID, Amount: amount}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdUnstake() *cobra.Command {
	cmd := &cobra.Command{
		Use: "unstake [pool-id] [amount]", Short: "Unstake LP tokens from a farming pool", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid pool-id: %w", err) }
			amount, ok := math.NewIntFromString(args[1])
			if !ok { return fmt.Errorf("invalid amount: %s", args[1]) }
			msg := &types.MsgUnstake{Sender: clientCtx.GetFromAddress().String(), PoolID: poolID, Amount: amount}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdClaimReward() *cobra.Command {
	cmd := &cobra.Command{
		Use: "claim-reward [pool-id]", Short: "Claim pending farming rewards", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid pool-id: %w", err) }
			msg := &types.MsgClaimReward{Sender: clientCtx.GetFromAddress().String(), PoolID: poolID}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCreateFarm() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create-farm [pool-id] [lp-denom] [reward-per-block] [start-block] [end-block]", Short: "Create a farming pool (authority only)", Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid pool-id: %w", err) }
			rpb, ok := math.NewIntFromString(args[2])
			if !ok { return fmt.Errorf("invalid reward-per-block: %s", args[2]) }
			sb, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil { return fmt.Errorf("invalid start-block: %w", err) }
			eb, err := strconv.ParseInt(args[4], 10, 64)
			if err != nil { return fmt.Errorf("invalid end-block: %w", err) }
			msg := &types.MsgCreateFarm{Authority: clientCtx.GetFromAddress().String(), PoolID: poolID, LPDenom: args[1], RewardPerBlock: rpb, StartBlock: sb, EndBlock: eb}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdUpdateFarm() *cobra.Command {
	cmd := &cobra.Command{
		Use: "update-farm [pool-id] [reward-per-block] [active]", Short: "Update a farming pool", Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid pool-id: %w", err) }
			rpb, ok := math.NewIntFromString(args[1])
			if !ok { return fmt.Errorf("invalid reward-per-block: %s", args[1]) }
			active := args[2] == "true"
			msg := &types.MsgUpdateFarm{Authority: clientCtx.GetFromAddress().String(), PoolID: poolID, RewardPerBlock: rpb, Active: active}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
