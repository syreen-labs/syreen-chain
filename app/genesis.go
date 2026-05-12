package app

import (
	"encoding/json"
	"time"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"

	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState defines the genesis state of the application.
type GenesisState map[string]json.RawMessage

// BondDenom is the native staking/fee denomination for the Syreen chain.
const BondDenom = "usyreen"

// DefaultGenesisState returns the default genesis state with usyreen as the
// bond denom instead of the SDK default "stake". This eliminates the need for
// post-init sed fixups on every fresh chain initialisation.
func DefaultGenesisState(cdc codec.JSONCodec) GenesisState {
	genState := ModuleBasics.DefaultGenesis(cdc)

	// Fix staking: denom, 150 validators, 21-day unbonding, 5% min commission
	var stakingGenesis stakingtypes.GenesisState
	cdc.MustUnmarshalJSON(genState[stakingtypes.ModuleName], &stakingGenesis)
	stakingGenesis.Params.BondDenom = BondDenom
	stakingGenesis.Params.MaxValidators = 150
	stakingGenesis.Params.UnbondingTime = 21 * 24 * time.Hour // 21 days
	stakingGenesis.Params.MinCommissionRate = math.LegacyNewDecWithPrec(5, 2) // 5%
	genState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(&stakingGenesis)

	// Fix mint: denom, inflation 5-15% range, target 10%, 63M blocks/year (~500ms blocks)
	var mintGenesis minttypes.GenesisState
	cdc.MustUnmarshalJSON(genState[minttypes.ModuleName], &mintGenesis)
	mintGenesis.Params.MintDenom = BondDenom
	mintGenesis.Params.InflationMin = math.LegacyNewDecWithPrec(5, 2)  // 5%
	mintGenesis.Params.InflationMax = math.LegacyNewDecWithPrec(15, 2) // 15%
	mintGenesis.Params.InflationRateChange = math.LegacyNewDecWithPrec(10, 2) // 10% rate of change
	mintGenesis.Params.GoalBonded = math.LegacyNewDecWithPrec(67, 2)  // 67% target bonded
	mintGenesis.Params.BlocksPerYear = 63115200                        // ~500ms blocks
	mintGenesis.Minter.Inflation = math.LegacyNewDecWithPrec(10, 2)   // start at 10%
	genState[minttypes.ModuleName] = cdc.MustMarshalJSON(&mintGenesis)

	// Fix crisis denom
	var crisisGenesis crisistypes.GenesisState
	cdc.MustUnmarshalJSON(genState[crisistypes.ModuleName], &crisisGenesis)
	crisisGenesis.ConstantFee = sdk.NewCoin(BondDenom, crisisGenesis.ConstantFee.Amount)
	genState[crisistypes.ModuleName] = cdc.MustMarshalJSON(&crisisGenesis)

	// Fix gov: denom, 100 SYR min deposit, 5-day voting period
	var govGenesis govv1.GenesisState
	cdc.MustUnmarshalJSON(genState["gov"], &govGenesis)
	if govGenesis.Params != nil {
		votingPeriod := 5 * 24 * time.Hour // 5 days
		depositPeriod := 5 * 24 * time.Hour
		govGenesis.Params.VotingPeriod = &votingPeriod
		govGenesis.Params.MaxDepositPeriod = &depositPeriod
		govGenesis.Params.MinDeposit = sdk.NewCoins(sdk.NewCoin(BondDenom, math.NewInt(100_000_000))) // 100 SYR
		govGenesis.Params.ExpeditedMinDeposit = sdk.NewCoins(sdk.NewCoin(BondDenom, math.NewInt(500_000_000))) // 500 SYR
		expeditedVoting := 2 * 24 * time.Hour // 2 days
		govGenesis.Params.ExpeditedVotingPeriod = &expeditedVoting
	}
	genState["gov"] = cdc.MustMarshalJSON(&govGenesis)

	return genState
}
