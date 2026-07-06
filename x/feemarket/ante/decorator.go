package ante

import (
	"context"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"syreen/x/feemarket/types"
)

// FeeMarketKeeper defines the interface needed by the fee market ante decorator.
type FeeMarketKeeper interface {
	GetBaseFee(ctx context.Context) math.LegacyDec
	GetFeeLane(ctx context.Context, name string) (types.FeeLane, bool)
	GetParams(ctx context.Context) types.Params
	BurnBaseFee(ctx context.Context, fees sdk.Coins) error
}

// FeeMarketDecorator enforces the EIP-1559 dynamic base fee on all transactions.
// It rejects transactions whose fee per gas is below the current base fee,
// and routes the base fee portion to the fee market module for burning.
type FeeMarketDecorator struct {
	feeMarketKeeper FeeMarketKeeper
	bankKeeper      types.BankKeeper
}

// NewFeeMarketDecorator creates a new FeeMarketDecorator.
func NewFeeMarketDecorator(fmk FeeMarketKeeper, bk types.BankKeeper) FeeMarketDecorator {
	return FeeMarketDecorator{
		feeMarketKeeper: fmk,
		bankKeeper:      bk,
	}
}

// AnteHandle checks that the transaction fee meets the dynamic base fee requirement.
func (fmd FeeMarketDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Skip fee check during simulation and genesis (block height 0)
	if simulate || ctx.BlockHeight() == 0 {
		return next(ctx, tx, simulate)
	}

	// Gasless DEX transactions skip fee validation but still proceed through
	// the handler chain. The DEX swap fee (0.3%) handles economic contribution.
	// FIX: Previously gasless TXs bypassed ALL fee logic including burns.
	if ctx.Value(GaslessTxKey{}) != nil {
		return next(ctx, tx, simulate)
	}

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.ErrTxDecode.Wrap("tx must implement sdk.FeeTx")
	}

	feeCoins := feeTx.GetFee()
	gasWanted := feeTx.GetGas()

	// If gas is 0, let standard validation handle it
	if gasWanted == 0 {
		return next(ctx, tx, simulate)
	}

	baseFee := fmd.feeMarketKeeper.GetBaseFee(ctx)

	// Determine the lane for this transaction (default lane for now)
	laneName := determineLane(tx)
	lane, found := fmd.feeMarketKeeper.GetFeeLane(ctx, laneName)
	if !found {
		// Fallback to default lane
		lane, found = fmd.feeMarketKeeper.GetFeeLane(ctx, "default")
		if !found {
			// If no lanes configured, use base fee with 1x multiplier
			lane = types.FeeLane{
				Name:              "default",
				BaseFeeMultiplier: math.LegacyOneDec(),
			}
		}
	}

	// Calculate required fee per gas: baseFee * laneMultiplier
	requiredFeePerGas := baseFee.Mul(lane.BaseFeeMultiplier)

	// Calculate minimum total fee: requiredFeePerGas * gasWanted
	requiredTotalFee := requiredFeePerGas.Mul(math.LegacyNewDecFromInt(math.NewIntFromUint64(gasWanted)))

	// Enforce minimum fee of 1000 usyreen (~0.001 SYREEN) to prevent zero-fee spam (H-07).
	// When the base fee drops very low (e.g., near-zero during low activity),
	// the calculated required fee could round to zero, allowing cheap/free transactions.
	minFee := math.LegacyNewDec(1000)
	if requiredTotalFee.LT(minFee) {
		requiredTotalFee = minFee
	}

	// Check if the provided fee meets the minimum
	// We check against the native denom (usyreen)
	const bondDenom = "usyreen"
	providedFee := math.LegacyZeroDec()
	for _, coin := range feeCoins {
		if coin.Denom == bondDenom {
			providedFee = math.LegacyNewDecFromInt(coin.Amount)
			break
		}
	}

	if providedFee.LT(requiredTotalFee) {
		return ctx, types.ErrFeeTooLow.Wrapf(
			"insufficient fee: provided %s%s, required minimum %s%s (base_fee=%s, lane=%s, gas=%d)",
			providedFee.TruncateInt().String(), bondDenom,
			requiredTotalFee.Ceil().TruncateInt().String(), bondDenom,
			baseFee.String(), laneName, gasWanted,
		)
	}

	// Burn the base fee portion (BurnRatio applied HERE, not in BurnBaseFee).
	// FIX: Previously BurnRatio was applied twice. Now applied only here.
	params := fmd.feeMarketKeeper.GetParams(ctx)
	burnAmount := requiredFeePerGas.Mul(math.LegacyNewDecFromInt(math.NewIntFromUint64(gasWanted))).Mul(params.BurnRatio)
	burnInt := burnAmount.TruncateInt()

	if burnInt.IsPositive() {
		// I-16: Verify fee_collector has sufficient balance before burning.
		// Cap the burn at the available balance to prevent failures.
		feeCollectorAddr := authtypes.NewModuleAddress("fee_collector")
		feeCollectorBal := fmd.bankKeeper.GetBalance(ctx, feeCollectorAddr, bondDenom)
		if burnInt.GT(feeCollectorBal.Amount) {
			burnInt = feeCollectorBal.Amount // cap at available balance
		}

		if burnInt.IsPositive() {
			burnCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, burnInt))

			// Send the burn portion from fee_collector to feemarket module
			err := fmd.bankKeeper.SendCoinsFromModuleToModule(
				ctx,
				"fee_collector",
				types.ModuleName,
				burnCoins,
			)
			if err != nil {
				// M-07: Reject tx if burn transfer fails to preserve economic model
				return ctx, fmt.Errorf("feemarket: failed to transfer burn portion: %w", err)
			}
			// Burn the coins in the feemarket module account
			if err := fmd.feeMarketKeeper.BurnBaseFee(ctx, burnCoins); err != nil {
				return ctx, fmt.Errorf("feemarket: failed to burn base fee: %w", err)
			}
		}
	}

	return next(ctx, tx, simulate)
}

// determineLane inspects the transaction messages and returns the appropriate fee lane.
// For mixed-message transactions, the lane with the HIGHEST fee multiplier is used
// to prevent fee underpayment by mixing expensive messages with cheap-lane messages (H-08).
// The governance lane (0.8x) is only applied when ALL messages are governance messages.
func determineLane(tx sdk.Tx) string {
	msgs := tx.GetMsgs()
	if len(msgs) == 0 {
		return "default"
	}

	// Find the highest-cost (most expensive) lane across all messages.
	highestLane := "default"
	highestMultiplier := math.LegacyOneDec() // default multiplier 1.0

	allGovernance := true

	defiMultiplier := math.LegacyNewDecWithPrec(15, 1)  // 1.5
	ibcMultiplier := math.LegacyNewDecWithPrec(12, 1)    // 1.2

	for _, msg := range msgs {
		typeURL := sdk.MsgTypeURL(msg)

		if !isGovernanceMsg(typeURL) {
			allGovernance = false
		}

		switch {
		case isDeFiMsg(typeURL):
			if defiMultiplier.GT(highestMultiplier) {
				highestLane = "defi"
				highestMultiplier = defiMultiplier
			}
		case isIBCMsg(typeURL):
			if ibcMultiplier.GT(highestMultiplier) {
				highestLane = "ibc"
				highestMultiplier = ibcMultiplier
			}
			// Governance is a cheaper lane (0.8x), so we never "upgrade" to it
			// from default. It is handled below.
		}
	}

	// Only use the governance lane (lower fee) if ALL messages are governance.
	// This prevents a mixed tx like [MsgVote, MsgSend] from getting the
	// discounted governance rate.
	if allGovernance && highestLane == "default" {
		return "governance"
	}

	return highestLane
}

func isGovernanceMsg(typeURL string) bool {
	switch typeURL {
	case "/cosmos.gov.v1.MsgVote",
		"/cosmos.gov.v1.MsgVoteWeighted",
		"/cosmos.gov.v1.MsgSubmitProposal",
		"/cosmos.gov.v1.MsgDeposit",
		"/cosmos.gov.v1beta1.MsgVote",
		"/cosmos.gov.v1beta1.MsgVoteWeighted",
		"/cosmos.gov.v1beta1.MsgSubmitProposal",
		"/cosmos.gov.v1beta1.MsgDeposit":
		return true
	}
	return false
}

func isIBCMsg(typeURL string) bool {
	// IBC messages all start with /ibc.
	if len(typeURL) > 5 && typeURL[:5] == "/ibc." {
		return true
	}
	return false
}

func isDeFiMsg(typeURL string) bool {
	switch typeURL {
	case "/cosmos.staking.v1beta1.MsgDelegate",
		"/cosmos.staking.v1beta1.MsgUndelegate",
		"/cosmos.staking.v1beta1.MsgBeginRedelegate",
		"/syreen.tokenfactory.MsgMint",
		"/syreen.tokenfactory.MsgBurn",
		"/syreen.intent.MsgSubmitIntent",
		"/syreen.intent.MsgSubmitSolution":
		return true
	}
	return false
}
