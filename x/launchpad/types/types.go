package types

import "cosmossdk.io/math"

// Launch status constants
const (
	LaunchStatusPending    = "pending"
	LaunchStatusActive     = "active"
	LaunchStatusSuccessful = "successful"
	LaunchStatusFailed     = "failed"
	LaunchStatusFinalized  = "finalized"
)

// Launch represents a token launch / IDO
type Launch struct {
	ID              uint64   `json:"id"`
	Creator         string   `json:"creator"`
	TokenDenom      string   `json:"token_denom"`       // subdenom for tokenfactory
	TokenSupply     math.Int `json:"token_supply"`       // total tokens to distribute
	PricePerToken   math.Int `json:"price_per_token"`    // quote tokens per 1 token (in base units)
	QuoteDenom      string   `json:"quote_denom"`        // e.g. factory/.../uusdc
	SoftCap         math.Int `json:"soft_cap"`           // minimum raise in quote tokens
	HardCap         math.Int `json:"hard_cap"`           // maximum raise in quote tokens
	MaxPerWallet    math.Int `json:"max_per_wallet"`     // max contribution per wallet
	StartBlock      int64    `json:"start_block"`
	EndBlock        int64    `json:"end_block"`
	Status          string   `json:"status"`
	TotalRaised     math.Int `json:"total_raised"`
	Contributors    uint64   `json:"contributors"`
	VestingBlocks   int64    `json:"vesting_blocks"`     // 0 = no vesting, all at TGE
	TGEPercent      uint64   `json:"tge_percent"`        // % unlocked at TGE (0-100)
	FinalizedBlock  int64    `json:"finalized_block"`
	DexPoolID       uint64   `json:"dex_pool_id"`        // created after finalization
	FullTokenDenom  string   `json:"full_token_denom"`   // full denom from tokenfactory
}

// Contribution represents a user's contribution to a launch
type Contribution struct {
	Address   string   `json:"address"`
	LaunchID  uint64   `json:"launch_id"`
	Amount    math.Int `json:"amount"`     // quote tokens contributed
	Claimed   bool     `json:"claimed"`
	ClaimedAt int64    `json:"claimed_at"`
}

// VestingPosition tracks token vesting for a contributor
type VestingPosition struct {
	Address        string   `json:"address"`
	LaunchID       uint64   `json:"launch_id"`
	TotalTokens    math.Int `json:"total_tokens"`
	Claimed        math.Int `json:"claimed"`
	TGEAmount      math.Int `json:"tge_amount"`
	VestingPerBlock math.Int `json:"vesting_per_block"`
	VestStart      int64    `json:"vest_start"`
}

// GenesisState
type GenesisState struct {
	Launches      []Launch         `json:"launches"`
	Contributions []Contribution   `json:"contributions"`
	Vestings      []VestingPosition `json:"vestings"`
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Launches:      []Launch{},
		Contributions: []Contribution{},
		Vestings:      []VestingPosition{},
	}
}

func (gs GenesisState) Validate() error { return nil }
