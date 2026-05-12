package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Message types
const (
	TypeMsgStake       = "stake_lp"
	TypeMsgUnstake     = "unstake_lp"
	TypeMsgClaimReward = "claim_reward"
	TypeMsgCreateFarm  = "create_farm"
	TypeMsgUpdateFarm  = "update_farm"
)

// MsgStake stakes LP tokens into a farming pool
type MsgStake struct {
	Sender string   `json:"sender"`
	PoolID uint64   `json:"pool_id"`
	Amount math.Int `json:"amount"`
}

func (msg MsgStake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return err
	}
	if msg.Amount.IsNil() || !msg.Amount.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}

// MsgUnstake removes LP tokens from a farming pool
type MsgUnstake struct {
	Sender string   `json:"sender"`
	PoolID uint64   `json:"pool_id"`
	Amount math.Int `json:"amount"`
}

func (msg MsgUnstake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return err
	}
	if msg.Amount.IsNil() || !msg.Amount.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}

// MsgClaimReward claims pending farming rewards
type MsgClaimReward struct {
	Sender string `json:"sender"`
	PoolID uint64 `json:"pool_id"`
}

func (msg MsgClaimReward) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return err
	}
	return nil
}

// MsgCreateFarm creates a new farming pool (governance or authority only)
type MsgCreateFarm struct {
	Authority      string   `json:"authority"`
	PoolID         uint64   `json:"pool_id"`
	LPDenom        string   `json:"lp_denom"`         // LP token denom from DEX pool (required)
	RewardPerBlock math.Int `json:"reward_per_block"`
	StartBlock     int64    `json:"start_block"`
	EndBlock       int64    `json:"end_block"` // 0 = no end
}

func (msg MsgCreateFarm) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return err
	}
	if msg.LPDenom == "" {
		return ErrInvalidAmount // LP denom must be set
	}
	if msg.RewardPerBlock.IsNil() || !msg.RewardPerBlock.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}

// MsgUpdateFarm updates a farming pool reward rate
type MsgUpdateFarm struct {
	Authority      string   `json:"authority"`
	PoolID         uint64   `json:"pool_id"`
	RewardPerBlock math.Int `json:"reward_per_block"`
	Active         bool     `json:"active"`
}

func (msg MsgUpdateFarm) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return err
	}
	return nil
}

// Response types
type MsgStakeResponse struct {
	PendingReward math.Int `json:"pending_reward"`
}

type MsgUnstakeResponse struct {
	ClaimedReward math.Int `json:"claimed_reward"`
}

type MsgClaimRewardResponse struct {
	Amount math.Int `json:"amount"`
}

type MsgCreateFarmResponse struct{}
type MsgUpdateFarmResponse struct{}
