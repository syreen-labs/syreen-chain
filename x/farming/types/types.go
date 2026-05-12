package types

import (
	"cosmossdk.io/math"
)

// FarmPool represents a farming pool that distributes rewards to LP stakers
type FarmPool struct {
	PoolID          uint64        `json:"pool_id"`
	LPDenom         string        `json:"lp_denom"`          // LP token denom from DEX pool
	RewardDenom     string        `json:"reward_denom"`       // reward token (usyreen)
	RewardPerBlock  math.Int      `json:"reward_per_block"`   // rewards distributed per block
	TotalStaked     math.Int      `json:"total_staked"`       // total LP tokens staked
	AccRewardPerShare math.LegacyDec `json:"acc_reward_per_share"` // accumulated reward per share (scaled by 1e12)
	LastRewardBlock int64         `json:"last_reward_block"`  // last block rewards were calculated
	StartBlock      int64         `json:"start_block"`        // farming start block
	EndBlock        int64         `json:"end_block"`          // farming end block (0 = no end)
	Active          bool          `json:"active"`
	Creator         string        `json:"creator"`            // who created the farm (governance or authority)
}

// Position represents a user's staked LP tokens in a farm
type Position struct {
	Address     string        `json:"address"`
	PoolID      uint64        `json:"pool_id"`
	Amount      math.Int      `json:"amount"`           // LP tokens staked
	RewardDebt  math.LegacyDec `json:"reward_debt"`     // reward debt for accurate reward calc
	PendingReward math.Int    `json:"pending_reward"`   // unclaimed rewards
	StakedAt    int64         `json:"staked_at"`        // block height when staked
}

// FarmingParams stores global farming parameters
type FarmingParams struct {
	DefaultRewardDenom string   `json:"default_reward_denom"` // usyreen
	MaxRewardPerBlock  math.Int `json:"max_reward_per_block"` // max reward per block per farm
	EcosystemAddress   string   `json:"ecosystem_address"`    // address that funds rewards
}

// GenesisState defines the farming module's genesis state
type GenesisState struct {
	Farms     []FarmPool     `json:"farms"`
	Positions []Position     `json:"positions"`
	Params    FarmingParams  `json:"params"`
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Farms:     []FarmPool{},
		Positions: []Position{},
		Params: FarmingParams{
			DefaultRewardDenom: "usyreen",
			MaxRewardPerBlock:  math.NewInt(100_000_000), // 100 SYR per block max
			EcosystemAddress:   "",
		},
	}
}

// Validate performs genesis state validation
func (gs GenesisState) Validate() error {
	return nil
}

// Precision factor for reward calculations (1e12)
var RewardPrecision = math.LegacyNewDec(1_000_000_000_000)
