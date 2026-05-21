package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"syreen/x/options/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: types.ModuleName, Short: "Options transaction subcommands",
		DisableFlagParsing: true, SuggestionsMinimumDistance: 2, RunE: client.ValidateCmd,
	}
	cmd.AddCommand(CmdWriteOption(), CmdBuyOption(), CmdExerciseOption(), CmdCancelOption())
	return cmd
}

func CmdWriteOption() *cobra.Command {
	cmd := &cobra.Command{
		Use: "write [pool-id] [option-type] [strike-price] [amount] [expiry-block] [custom-premium]", Short: "Write (create) an option", Args: cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid pool-id: %w", err) }
			strike, err := math.LegacyNewDecFromStr(args[2])
			if err != nil { return fmt.Errorf("invalid strike-price: %w", err) }
			amount, ok := math.NewIntFromString(args[3])
			if !ok { return fmt.Errorf("invalid amount: %s", args[3]) }
			expiry, err := strconv.ParseInt(args[4], 10, 64)
			if err != nil { return fmt.Errorf("invalid expiry-block: %w", err) }
			premium, ok := math.NewIntFromString(args[5])
			if !ok { return fmt.Errorf("invalid custom-premium: %s", args[5]) }
			msg := &types.MsgWriteOption{
				Writer: clientCtx.GetFromAddress().String(), PoolID: poolID,
				OptionType: types.OptionType(args[1]), StrikePrice: strike,
				Amount: amount, ExpiryBlock: expiry, CustomPremium: premium,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdBuyOption() *cobra.Command {
	cmd := &cobra.Command{
		Use: "buy [option-id]", Short: "Buy an option", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			oid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid option-id: %w", err) }
			msg := &types.MsgBuyOption{Buyer: clientCtx.GetFromAddress().String(), OptionID: oid}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdExerciseOption() *cobra.Command {
	cmd := &cobra.Command{
		Use: "exercise [option-id]", Short: "Exercise an option", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			oid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid option-id: %w", err) }
			msg := &types.MsgExerciseOption{Buyer: clientCtx.GetFromAddress().String(), OptionID: oid}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCancelOption() *cobra.Command {
	cmd := &cobra.Command{
		Use: "cancel [option-id]", Short: "Cancel an unsold option", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			oid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid option-id: %w", err) }
			msg := &types.MsgCancelOption{Writer: clientCtx.GetFromAddress().String(), OptionID: oid}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
