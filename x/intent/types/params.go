package types

import (
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Default parameter values
var (
	DefaultMinSolverStake        = sdk.NewCoin("usyreen", math.NewInt(1000000000)) // 1000 SYR
	DefaultSolvingWindow         = uint64(100)
	DefaultMaxSolutionsPerIntent = uint64(20)
	DefaultSolverSlashFraction   = math.LegacyNewDecWithPrec(10, 2) // 10%
	// DefaultIntentExpiryBlocks was originally 200 which was far too short
	// for multi-tranche DCA chains (a 10-tranche DCA needs 6000+ blocks and
	// SubmitChain caps chain expiry at IntentExpiryBlocks * numSteps). Raised
	// to 20000 (~28 hours at 5s blocks) so long-running strategy chains work
	// out of the box without the cap slamming them down to 400 blocks.
	DefaultIntentExpiryBlocks    = uint64(20000)
	DefaultSolverUnbondingBlocks = uint64(100) // Stake locked for 100 blocks after deregistration
	DefaultEnableIntents         = true
	DefaultMaxChainSteps         = uint64(10)
	DefaultEnableChains          = true
	DefaultAllowedMsgTypes       = []string{
		"/cosmos.bank.v1beta1.MsgSend",
		"/cosmos.bank.v1beta1.MsgMultiSend",
		"/ibc.applications.transfer.v1.MsgTransfer",
		"/syreen.compute.MsgExecuteContract",
		// DEX swaps are the legitimate best-execution settlement primitive: a
		// solver crafts a creator-signed MsgSwap that consumes the creator's own
		// input (bounded by the executeSolutionMsgs balance invariants) and
		// delivers output back to the creator. This is safe because it cannot
		// touch pooled escrow or spend a creator denom beyond the declared input.
		"/syreen.dex.MsgSwap",
		"/syreen.dex.MsgMultiHopSwap",
	}
)

// Params defines the parameters for the intent module
type Params struct {
	MinSolverStake        sdk.Coin       `json:"min_solver_stake"`
	SolvingWindow         uint64         `json:"solving_window"`
	MaxSolutionsPerIntent uint64         `json:"max_solutions_per_intent"`
	SolverSlashFraction   math.LegacyDec `json:"solver_slash_fraction"`
	IntentExpiryBlocks    uint64         `json:"intent_expiry_blocks"`
	SolverUnbondingBlocks uint64         `json:"solver_unbonding_blocks"`
	EnableIntents         bool           `json:"enable_intents"`
	AllowedMsgTypes       []string       `json:"allowed_msg_types"`
	MaxChainSteps         uint64         `json:"max_chain_steps"`
	EnableChains          bool           `json:"enable_chains"`
}

// DefaultParams returns the default set of parameters
func DefaultParams() Params {
	return Params{
		MinSolverStake:        DefaultMinSolverStake,
		SolvingWindow:         DefaultSolvingWindow,
		MaxSolutionsPerIntent: DefaultMaxSolutionsPerIntent,
		SolverSlashFraction:   DefaultSolverSlashFraction,
		IntentExpiryBlocks:    DefaultIntentExpiryBlocks,
		SolverUnbondingBlocks: DefaultSolverUnbondingBlocks,
		EnableIntents:         DefaultEnableIntents,
		AllowedMsgTypes:       DefaultAllowedMsgTypes,
		MaxChainSteps:         DefaultMaxChainSteps,
		EnableChains:          DefaultEnableChains,
	}
}

// Validate checks that the parameters are valid
func (p Params) Validate() error {
	if !p.MinSolverStake.IsPositive() {
		return fmt.Errorf("min solver stake must be positive: %s", p.MinSolverStake)
	}
	if p.SolvingWindow == 0 {
		return fmt.Errorf("solving window must be greater than 0")
	}
	if p.MaxSolutionsPerIntent == 0 {
		return fmt.Errorf("max solutions per intent must be greater than 0")
	}
	if p.SolverSlashFraction.IsNegative() || p.SolverSlashFraction.GT(math.LegacyOneDec()) {
		return fmt.Errorf("solver slash fraction must be between 0 and 1: %s", p.SolverSlashFraction)
	}
	if p.IntentExpiryBlocks == 0 {
		return fmt.Errorf("intent expiry blocks must be greater than 0")
	}
	if p.SolverUnbondingBlocks == 0 {
		return fmt.Errorf("solver unbonding blocks must be greater than 0")
	}
	if p.EnableChains && p.MaxChainSteps == 0 {
		return fmt.Errorf("max chain steps must be greater than 0 when chains are enabled")
	}
	for _, t := range p.AllowedMsgTypes {
		if t == "" || t[0] != '/' {
			return fmt.Errorf("allowed msg type must be a non-empty type URL beginning with '/': %q", t)
		}
	}
	return nil
}
