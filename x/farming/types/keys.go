package types

import "encoding/binary"

const (
	ModuleName = "farming"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	// KV prefixes
	FarmPrefix         = "farm/"          // farm/{poolID} → FarmPool
	PositionPrefix     = "pos/"           // pos/{poolID}/{address} → Position
	PositionByAddr     = "pos_addr/"      // pos_addr/{address}/{poolID} → exists
	RewardAccumPrefix  = "raccum/"        // raccum/{poolID} → accumulated reward per share
	GlobalRewardPrefix = "greward/"       // greward → total rewards distributed
	ParamsKey          = "params"
)

func FarmKey(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(FarmPrefix), bz...)
}

func PositionKey(poolID uint64, address string) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append(append([]byte(PositionPrefix), bz...), []byte("/"+address)...)
}

func PositionByAddrKey(address string, poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append(append([]byte(PositionByAddr), []byte(address+"/")...), bz...)
}

func PositionByAddrPrefix(address string) []byte {
	return []byte(PositionByAddr + address + "/")
}

func PositionPoolPrefix(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(PositionPrefix), bz...)
}

func RewardAccumKey(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(RewardAccumPrefix), bz...)
}
