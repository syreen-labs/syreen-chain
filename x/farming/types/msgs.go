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
	Sender string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	PoolID uint64   `protobuf:"varint,2,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
	Amount math.Int `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
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
	Sender string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	PoolID uint64   `protobuf:"varint,2,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
	Amount math.Int `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
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
	Sender string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	PoolID uint64 `protobuf:"varint,2,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
}

func (msg MsgClaimReward) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return err
	}
	return nil
}

// MsgCreateFarm creates a new farming pool (governance or authority only)
type MsgCreateFarm struct {
	Authority      string   `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority"`
	PoolID         uint64   `protobuf:"varint,2,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
	LPDenom        string   `protobuf:"bytes,3,opt,name=lp_denom,json=lpDenom,proto3" json:"lp_denom"`
	RewardPerBlock math.Int `protobuf:"bytes,4,opt,name=reward_per_block,json=rewardPerBlock,proto3" json:"reward_per_block"`
	StartBlock     int64    `protobuf:"varint,5,opt,name=start_block,json=startBlock,proto3" json:"start_block"`
	EndBlock       int64    `protobuf:"varint,6,opt,name=end_block,json=endBlock,proto3" json:"end_block"`
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
	Authority      string   `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority"`
	PoolID         uint64   `protobuf:"varint,2,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
	RewardPerBlock math.Int `protobuf:"bytes,3,opt,name=reward_per_block,json=rewardPerBlock,proto3" json:"reward_per_block"`
	Active         bool     `protobuf:"varint,4,opt,name=active,proto3" json:"active"`
}

func (msg MsgUpdateFarm) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return err
	}
	return nil
}

// Response types
type MsgStakeResponse struct {
	PendingReward math.Int `protobuf:"bytes,1,opt,name=pending_reward,json=pendingReward,proto3" json:"pending_reward"`
}

type MsgUnstakeResponse struct {
	ClaimedReward math.Int `protobuf:"bytes,1,opt,name=claimed_reward,json=claimedReward,proto3" json:"claimed_reward"`
}

type MsgClaimRewardResponse struct {
	Amount math.Int `protobuf:"bytes,1,opt,name=amount,proto3" json:"amount"`
}

type MsgCreateFarmResponse struct{}
type MsgUpdateFarmResponse struct{}
