package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"syreen/x/launchpad/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: types.ModuleName, Short: "Launchpad transaction subcommands",
		DisableFlagParsing: true, SuggestionsMinimumDistance: 2, RunE: client.ValidateCmd,
	}
	cmd.AddCommand(CmdCreateLaunch(), CmdContribute(), CmdClaimTokens(), CmdClaimRefund(), CmdFinalizeLaunch())
	return cmd
}

func CmdCreateLaunch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-launch [token-denom] [token-supply] [price-per-token] [quote-denom] [soft-cap] [hard-cap] [max-per-wallet] [start-block] [end-block] [vesting-blocks] [tge-percent]",
		Short: "Create a new token launch", Args: cobra.ExactArgs(11),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			ts, ok := math.NewIntFromString(args[1])
			if !ok { return fmt.Errorf("invalid token-supply: %s", args[1]) }
			ppt, ok := math.NewIntFromString(args[2])
			if !ok { return fmt.Errorf("invalid price-per-token: %s", args[2]) }
			sc, ok := math.NewIntFromString(args[4])
			if !ok { return fmt.Errorf("invalid soft-cap: %s", args[4]) }
			hc, ok := math.NewIntFromString(args[5])
			if !ok { return fmt.Errorf("invalid hard-cap: %s", args[5]) }
			mpw, ok := math.NewIntFromString(args[6])
			if !ok { return fmt.Errorf("invalid max-per-wallet: %s", args[6]) }
			sb, err := strconv.ParseInt(args[7], 10, 64)
			if err != nil { return fmt.Errorf("invalid start-block: %w", err) }
			eb, err := strconv.ParseInt(args[8], 10, 64)
			if err != nil { return fmt.Errorf("invalid end-block: %w", err) }
			vb, err := strconv.ParseInt(args[9], 10, 64)
			if err != nil { return fmt.Errorf("invalid vesting-blocks: %w", err) }
			tge, err := strconv.ParseUint(args[10], 10, 64)
			if err != nil { return fmt.Errorf("invalid tge-percent: %w", err) }
			msg := &types.MsgCreateLaunch{
				Creator: clientCtx.GetFromAddress().String(), TokenDenom: args[0], TokenSupply: ts,
				PricePerToken: ppt, QuoteDenom: args[3], SoftCap: sc, HardCap: hc,
				MaxPerWallet: mpw, StartBlock: sb, EndBlock: eb, VestingBlocks: vb, TGEPercent: tge,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdContribute() *cobra.Command {
	cmd := &cobra.Command{
		Use: "contribute [launch-id] [amount]", Short: "Contribute to a token launch", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			lid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid launch-id: %w", err) }
			amount, ok := math.NewIntFromString(args[1])
			if !ok { return fmt.Errorf("invalid amount: %s", args[1]) }
			msg := &types.MsgContribute{Sender: clientCtx.GetFromAddress().String(), LaunchID: lid, Amount: amount}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdClaimTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use: "claim-tokens [launch-id]", Short: "Claim tokens after a successful launch", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			lid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid launch-id: %w", err) }
			msg := &types.MsgClaimTokens{Sender: clientCtx.GetFromAddress().String(), LaunchID: lid}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdClaimRefund() *cobra.Command {
	cmd := &cobra.Command{
		Use: "claim-refund [launch-id]", Short: "Claim refund after a failed launch", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			lid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid launch-id: %w", err) }
			msg := &types.MsgClaimRefund{Sender: clientCtx.GetFromAddress().String(), LaunchID: lid}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdFinalizeLaunch() *cobra.Command {
	cmd := &cobra.Command{
		Use: "finalize [launch-id]", Short: "Finalize a launch (authority only)", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			lid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid launch-id: %w", err) }
			msg := &types.MsgFinalizeLaunch{Authority: clientCtx.GetFromAddress().String(), LaunchID: lid}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
