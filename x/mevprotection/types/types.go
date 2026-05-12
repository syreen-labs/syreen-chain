package types

import "fmt"

// CommittedTx represents an encrypted transaction submitted in the commit phase
type CommittedTx struct {
	Sender      string `json:"sender"`
	TxHash      []byte `json:"tx_hash"`
	EncryptedTx []byte `json:"encrypted_tx"` // Stored encrypted payload, revealed later
	BlockHeight int64  `json:"block_height"`
	Timestamp   int64  `json:"timestamp"`
}

func (m *CommittedTx) ProtoMessage()           {}
func (m *CommittedTx) Reset()                  { *m = CommittedTx{} }
func (m *CommittedTx) String() string          { return fmt.Sprintf("committed_tx: %s at height %d", m.Sender, m.BlockHeight) }
func (m *CommittedTx) XXX_MessageName() string { return "syreen.mevprotection.CommittedTx" }

// RevealedTx represents a revealed transaction matched against a prior commitment
type RevealedTx struct {
	CommittedHash []byte `json:"committed_hash"`
	ActualTx      []byte `json:"actual_tx"`
	Proof         []byte `json:"proof"`
}

func (m *RevealedTx) ProtoMessage()           {}
func (m *RevealedTx) Reset()                  { *m = RevealedTx{} }
func (m *RevealedTx) String() string          { return fmt.Sprintf("revealed_tx: hash=%x", m.CommittedHash) }
func (m *RevealedTx) XXX_MessageName() string { return "syreen.mevprotection.RevealedTx" }

// FairOrderConfig defines the configuration for fair ordering enforcement
type FairOrderConfig struct {
	EnableCommitReveal bool   `json:"enable_commit_reveal"`
	CommitWindow       uint64 `json:"commit_window"`
	RevealWindow       uint64 `json:"reveal_window"`
	MaxTxDelay         uint64 `json:"max_tx_delay"`
}

func (m *FairOrderConfig) ProtoMessage()           {}
func (m *FairOrderConfig) Reset()                  { *m = FairOrderConfig{} }
func (m *FairOrderConfig) String() string          { return fmt.Sprintf("fair_order_config: commit=%v", m.EnableCommitReveal) }
func (m *FairOrderConfig) XXX_MessageName() string { return "syreen.mevprotection.FairOrderConfig" }

// MEVPenalty records a penalty assessed against a validator for MEV extraction
type MEVPenalty struct {
	ValidatorAddr string `json:"validator_addr"`
	PenaltyType   string `json:"penalty_type"`
	Amount        uint64 `json:"amount"`
	Evidence      []byte `json:"evidence"`
	BlockHeight   int64  `json:"block_height"`
}

func (m *MEVPenalty) ProtoMessage()           {}
func (m *MEVPenalty) Reset()                  { *m = MEVPenalty{} }
func (m *MEVPenalty) String() string          { return fmt.Sprintf("mev_penalty: %s type=%s at height %d", m.ValidatorAddr, m.PenaltyType, m.BlockHeight) }
func (m *MEVPenalty) XXX_MessageName() string { return "syreen.mevprotection.MEVPenalty" }
