package types

import "fmt"

// Access type constants for contract instantiation permissions
const (
	AccessTypeEverybody  = "Everybody"
	AccessTypeOnlyAddress = "OnlyAddress"
	AccessTypeNobody     = "Nobody"
)

// AccessConfig defines the permission for code instantiation
type AccessConfig struct {
	Permission string `json:"permission"`
	Address    string `json:"address,omitempty"`
}

// CodeInfo stores metadata about uploaded code
type CodeInfo struct {
	CodeID                uint64       `json:"code_id"`
	Creator               string       `json:"creator"`
	CodeHash              []byte       `json:"code_hash"`
	InstantiatePermission AccessConfig `json:"instantiate_permission"`
}

func (m *CodeInfo) ProtoMessage()           {}
func (m *CodeInfo) Reset()                  { *m = CodeInfo{} }
func (m *CodeInfo) String() string          { return fmt.Sprintf("code %d by %s", m.CodeID, m.Creator) }
func (m *CodeInfo) XXX_MessageName() string { return "syreen.compute.CodeInfo" }

// ContractInfo stores metadata about a contract instance
type ContractInfo struct {
	Address   string `json:"address"`
	CodeID    uint64 `json:"code_id"`
	Creator   string `json:"creator"`
	Admin     string `json:"admin,omitempty"`
	Label     string `json:"label"`
	CreatedAt int64  `json:"created_at"`
}

func (m *ContractInfo) ProtoMessage()           {}
func (m *ContractInfo) Reset()                  { *m = ContractInfo{} }
func (m *ContractInfo) String() string          { return fmt.Sprintf("contract %s (code %d)", m.Address, m.CodeID) }
func (m *ContractInfo) XXX_MessageName() string { return "syreen.compute.ContractInfo" }

// ContractState represents a key-value pair in contract storage
type ContractState struct {
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}

func (m *ContractState) ProtoMessage()           {}
func (m *ContractState) Reset()                  { *m = ContractState{} }
func (m *ContractState) String() string          { return fmt.Sprintf("state key=%x", m.Key) }
func (m *ContractState) XXX_MessageName() string { return "syreen.compute.ContractState" }

// Model represents a contract state entry for genesis export
type Model struct {
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}

// ContractCodeHistoryEntry tracks code migration history
type ContractCodeHistoryEntry struct {
	CodeID  uint64 `json:"code_id"`
	Updated int64  `json:"updated"`
}
