package types

import "encoding/binary"

const (
	ModuleName = "aiagent"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	AgentPrefix    = "agent/"        // agent/{id} → Agent
	NextAgentIDKey = "next_agent_id"
)

func AgentKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return append([]byte(AgentPrefix), bz...)
}
