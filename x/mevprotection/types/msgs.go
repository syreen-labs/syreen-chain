package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgCommitTx = "commit_tx"
	TypeMsgRevealTx = "reveal_tx"
)

// MsgCommitTx submits an encrypted transaction commitment
type MsgCommitTx struct {
	Sender      string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	TxHash      []byte `protobuf:"bytes,2,opt,name=tx_hash,proto3" json:"tx_hash"`
	EncryptedTx []byte `protobuf:"bytes,3,opt,name=encrypted_tx,proto3" json:"encrypted_tx"`
}

func (m *MsgCommitTx) ProtoMessage()           {}
func (m *MsgCommitTx) Reset()                  { *m = MsgCommitTx{} }
func (m *MsgCommitTx) String() string          { return fmt.Sprintf("commit_tx: %s hash=%x", m.Sender, m.TxHash) }
func (m *MsgCommitTx) XXX_MessageName() string { return "syreen.mevprotection.MsgCommitTx" }

func (m *MsgCommitTx) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return fmt.Errorf("invalid sender address: %w", err)
	}
	// TxHash must be exactly 32 bytes (SHA-256 digest)
	if len(m.TxHash) != 32 {
		return fmt.Errorf("tx_hash must be exactly 32 bytes (SHA-256 digest), got %d", len(m.TxHash))
	}
	if len(m.EncryptedTx) == 0 {
		return fmt.Errorf("encrypted tx body cannot be empty")
	}
	// Max encrypted tx size: 256KB
	if len(m.EncryptedTx) > 256*1024 {
		return fmt.Errorf("encrypted tx body too large: %d bytes (max 262144)", len(m.EncryptedTx))
	}
	return nil
}

// MsgRevealTx reveals a previously committed transaction
type MsgRevealTx struct {
	Sender     string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	CommitHash []byte `protobuf:"bytes,2,opt,name=commit_hash,proto3" json:"commit_hash"`
	TxBody     []byte `protobuf:"bytes,3,opt,name=tx_body,proto3" json:"tx_body"`
	Nonce      []byte `protobuf:"bytes,4,opt,name=nonce,proto3" json:"nonce"`
}

func (m *MsgRevealTx) ProtoMessage()           {}
func (m *MsgRevealTx) Reset()                  { *m = MsgRevealTx{} }
func (m *MsgRevealTx) String() string          { return fmt.Sprintf("reveal_tx: %s commit=%x", m.Sender, m.CommitHash) }
func (m *MsgRevealTx) XXX_MessageName() string { return "syreen.mevprotection.MsgRevealTx" }

func (m *MsgRevealTx) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return fmt.Errorf("invalid sender address: %w", err)
	}
	// CommitHash must be exactly 32 bytes (SHA-256 digest)
	if len(m.CommitHash) != 32 {
		return fmt.Errorf("commit_hash must be exactly 32 bytes (SHA-256 digest), got %d", len(m.CommitHash))
	}
	if len(m.TxBody) == 0 {
		return fmt.Errorf("tx body cannot be empty")
	}
	if len(m.TxBody) > 256*1024 {
		return fmt.Errorf("tx body too large: %d bytes (max 262144)", len(m.TxBody))
	}
	if len(m.Nonce) == 0 {
		return fmt.Errorf("nonce cannot be empty")
	}
	if len(m.Nonce) > 64 {
		return fmt.Errorf("nonce too large: %d bytes (max 64)", len(m.Nonce))
	}
	return nil
}

