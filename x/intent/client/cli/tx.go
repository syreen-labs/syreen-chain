package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/intent/types"
)

// GetTxCmd returns the transaction commands for the intent module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Intent module transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdSubmitIntent(),
		CmdRegisterSolver(),
		CmdDeregisterSolver(),
		CmdSubmitSolution(),
		CmdFulfillIntent(),
		CmdLimitBuy(),
		CmdLimitSell(),
		CmdStopLoss(),
		CmdTakeProfit(),
		CmdDCA(),
	)

	return cmd
}

func CmdSubmitIntent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-intent [type] [body-json] [max-fee-amount] [max-fee-denom] [tip-amount] [tip-denom]",
		Short: "Submit a new intent declaration",
		Long: `Submit a new intent. Types: swap, transfer, defi, custom.
Body is a JSON string describing the intent details.`,
		Args: cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			intentType := args[0]
			body := json.RawMessage(args[1])
			if !json.Valid(body) {
				return fmt.Errorf("invalid JSON body: %s", args[1])
			}

			maxFeeAmount, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid max-fee-amount: %w", err)
			}
			maxFeeDenom := args[3]

			tipAmount, err := strconv.ParseInt(args[4], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid tip-amount: %w", err)
			}
			tipDenom := args[5]

			msg := &types.MsgSubmitIntent{
				Creator:      clientCtx.GetFromAddress().String(),
				IntentType:   intentType,
				Body:         body,
				MaxFee:       sdk.NewCoins(sdk.NewCoin(maxFeeDenom, math.NewInt(maxFeeAmount))),
				Tip:          sdk.NewCoins(sdk.NewCoin(tipDenom, math.NewInt(tipAmount))),
				ExpiryBlocks: 100, // Default expiry
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRegisterSolver() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-solver [moniker] [stake-amount] [stake-denom]",
		Short: "Register as an intent solver",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			moniker := args[0]

			stakeAmount, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid stake-amount: %w", err)
			}
			stakeDenom := args[2]

			msg := &types.MsgRegisterSolver{
				Address:     clientCtx.GetFromAddress().String(),
				Moniker:     moniker,
				StakeAmount: sdk.NewCoin(stakeDenom, math.NewInt(stakeAmount)),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdDeregisterSolver() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deregister-solver",
		Short: "Deregister as an intent solver",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgDeregisterSolver{
				Address: clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdSubmitSolution() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-solution [intent-id] [solution-msgs-json] [expected-outcome] [gas-estimate]",
		Short: "Submit a solution for an intent",
		Long: `Submit a solution for a pending intent.
solution-msgs-json: JSON array of execution messages.
expected-outcome: JSON describing expected outcome.
gas-estimate: estimated gas for execution.`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			intentID := args[0]

			// Parse execution msgs as JSON array
			var execMsgs []json.RawMessage
			if err := json.Unmarshal([]byte(args[1]), &execMsgs); err != nil {
				return fmt.Errorf("invalid solution-msgs-json: %w", err)
			}

			expectedOutcome := json.RawMessage(args[2])
			if !json.Valid(expectedOutcome) {
				return fmt.Errorf("invalid expected-outcome JSON: %s", args[2])
			}

			_, err = strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid gas-estimate: %w", err)
			}

			msg := &types.MsgSubmitSolution{
				SolverAddr:      clientCtx.GetFromAddress().String(),
				IntentID:        intentID,
				ExecutionMsgs:   execMsgs,
				ExpectedOutcome: expectedOutcome,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdLimitBuy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "limit-buy [pool-id] [input-denom] [input-amount] [output-denom] [target-price] [min-output] [fee-amount] [expiry-blocks]",
		Short: "Place a limit buy order (buy output when price drops to target)",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool-id: %w", err)
			}

			inputAmount, ok := math.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid input-amount: %s", args[2])
			}

			targetPrice, err := math.LegacyNewDecFromStr(args[4])
			if err != nil {
				return fmt.Errorf("invalid target-price: %w", err)
			}

			minOutput, ok := math.NewIntFromString(args[5])
			if !ok {
				return fmt.Errorf("invalid min-output: %s", args[5])
			}

			feeAmount, err := strconv.ParseInt(args[6], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid fee-amount: %w", err)
			}

			expiryBlocks, err := strconv.ParseUint(args[7], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid expiry-blocks: %w", err)
			}

			body := types.LimitBuyIntent{
				InputDenom:      args[1],
				InputAmount:     inputAmount,
				OutputDenom:     args[3],
				TargetPrice:     targetPrice,
				PoolID:          poolID,
				MinOutputAmount: minOutput,
			}
			bodyBz, _ := json.Marshal(body)

			msg := &types.MsgSubmitIntent{
				Creator:      clientCtx.GetFromAddress().String(),
				IntentType:   types.IntentTypeLimitBuy,
				Body:         bodyBz,
				MaxFee:       sdk.NewCoins(sdk.NewCoin(args[1], math.NewInt(feeAmount))),
				Tip:          sdk.NewCoins(),
				ExpiryBlocks: expiryBlocks,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdLimitSell() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "limit-sell [pool-id] [input-denom] [input-amount] [output-denom] [target-price] [min-output] [fee-amount] [expiry-blocks]",
		Short: "Place a limit sell order (sell input when price rises to target)",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool-id: %w", err)
			}

			inputAmount, ok := math.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid input-amount: %s", args[2])
			}

			targetPrice, err := math.LegacyNewDecFromStr(args[4])
			if err != nil {
				return fmt.Errorf("invalid target-price: %w", err)
			}

			minOutput, ok := math.NewIntFromString(args[5])
			if !ok {
				return fmt.Errorf("invalid min-output: %s", args[5])
			}

			feeAmount, err := strconv.ParseInt(args[6], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid fee-amount: %w", err)
			}

			expiryBlocks, err := strconv.ParseUint(args[7], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid expiry-blocks: %w", err)
			}

			body := types.LimitSellIntent{
				InputDenom:      args[1],
				InputAmount:     inputAmount,
				OutputDenom:     args[3],
				TargetPrice:     targetPrice,
				PoolID:          poolID,
				MinOutputAmount: minOutput,
			}
			bodyBz, _ := json.Marshal(body)

			msg := &types.MsgSubmitIntent{
				Creator:      clientCtx.GetFromAddress().String(),
				IntentType:   types.IntentTypeLimitSell,
				Body:         bodyBz,
				MaxFee:       sdk.NewCoins(sdk.NewCoin(args[1], math.NewInt(feeAmount))),
				Tip:          sdk.NewCoins(),
				ExpiryBlocks: expiryBlocks,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdStopLoss() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop-loss [pool-id] [input-denom] [input-amount] [output-denom] [stop-price] [min-output] [fee-amount] [expiry-blocks]",
		Short: "Place a stop-loss order (sell input when price drops to stop price)",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool-id: %w", err)
			}

			inputAmount, ok := math.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid input-amount: %s", args[2])
			}

			stopPrice, err := math.LegacyNewDecFromStr(args[4])
			if err != nil {
				return fmt.Errorf("invalid stop-price: %w", err)
			}

			minOutput, ok := math.NewIntFromString(args[5])
			if !ok {
				return fmt.Errorf("invalid min-output: %s", args[5])
			}

			feeAmount, err := strconv.ParseInt(args[6], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid fee-amount: %w", err)
			}

			expiryBlocks, err := strconv.ParseUint(args[7], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid expiry-blocks: %w", err)
			}

			body := types.StopLossIntent{
				InputDenom:      args[1],
				InputAmount:     inputAmount,
				OutputDenom:     args[3],
				StopPrice:       stopPrice,
				PoolID:          poolID,
				MinOutputAmount: minOutput,
			}
			bodyBz, _ := json.Marshal(body)

			msg := &types.MsgSubmitIntent{
				Creator:      clientCtx.GetFromAddress().String(),
				IntentType:   types.IntentTypeStopLoss,
				Body:         bodyBz,
				MaxFee:       sdk.NewCoins(sdk.NewCoin(args[1], math.NewInt(feeAmount))),
				Tip:          sdk.NewCoins(),
				ExpiryBlocks: expiryBlocks,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdTakeProfit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "take-profit [pool-id] [input-denom] [input-amount] [output-denom] [target-price] [min-output] [fee-amount] [expiry-blocks]",
		Short: "Place a take-profit order (sell input when price rises to target)",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool-id: %w", err)
			}

			inputAmount, ok := math.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid input-amount: %s", args[2])
			}

			targetPrice, err := math.LegacyNewDecFromStr(args[4])
			if err != nil {
				return fmt.Errorf("invalid target-price: %w", err)
			}

			minOutput, ok := math.NewIntFromString(args[5])
			if !ok {
				return fmt.Errorf("invalid min-output: %s", args[5])
			}

			feeAmount, err := strconv.ParseInt(args[6], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid fee-amount: %w", err)
			}

			expiryBlocks, err := strconv.ParseUint(args[7], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid expiry-blocks: %w", err)
			}

			body := types.TakeProfitIntent{
				InputDenom:      args[1],
				InputAmount:     inputAmount,
				OutputDenom:     args[3],
				TargetPrice:     targetPrice,
				PoolID:          poolID,
				MinOutputAmount: minOutput,
			}
			bodyBz, _ := json.Marshal(body)

			msg := &types.MsgSubmitIntent{
				Creator:      clientCtx.GetFromAddress().String(),
				IntentType:   types.IntentTypeTakeProfit,
				Body:         bodyBz,
				MaxFee:       sdk.NewCoins(sdk.NewCoin(args[1], math.NewInt(feeAmount))),
				Tip:          sdk.NewCoins(),
				ExpiryBlocks: expiryBlocks,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdDCA() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dca [pool-id] [input-denom] [total-amount] [output-denom] [num-executions] [interval-blocks] [fee-amount] [expiry-blocks]",
		Short: "Place a DCA order (split buy across multiple blocks)",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool-id: %w", err)
			}

			totalAmount, ok := math.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid total-amount: %s", args[2])
			}

			numExec, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid num-executions: %w", err)
			}

			interval, err := strconv.ParseUint(args[5], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid interval-blocks: %w", err)
			}

			feeAmount, err := strconv.ParseInt(args[6], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid fee-amount: %w", err)
			}

			expiryBlocks, err := strconv.ParseUint(args[7], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid expiry-blocks: %w", err)
			}

			body := types.DCAIntent{
				InputDenom:     args[1],
				TotalAmount:    totalAmount,
				OutputDenom:    args[3],
				NumExecutions:  numExec,
				IntervalBlocks: interval,
				PoolID:         poolID,
			}
			bodyBz, _ := json.Marshal(body)

			msg := &types.MsgSubmitIntent{
				Creator:      clientCtx.GetFromAddress().String(),
				IntentType:   types.IntentTypeDCA,
				Body:         bodyBz,
				MaxFee:       sdk.NewCoins(sdk.NewCoin(args[1], math.NewInt(feeAmount))),
				Tip:          sdk.NewCoins(),
				ExpiryBlocks: expiryBlocks,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdFulfillIntent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fulfill-intent [intent-id] [optional:solver-addr]",
		Short: "Fulfill an intent (empty solver triggers auction selection)",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			intentID := args[0]
			solverAddr := ""
			if len(args) == 2 {
				solverAddr = args[1]
			}

			// If no solver addr provided, use the sender address
			if solverAddr == "" {
				solverAddr = clientCtx.GetFromAddress().String()
			}

			msg := &types.MsgFulfillIntent{
				SolverAddr: solverAddr,
				IntentID:   intentID,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
