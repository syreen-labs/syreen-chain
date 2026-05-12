package types

import "encoding/binary"

const (
	ModuleName = "clmm"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	CLPoolPrefix      = "clpool/"     // clpool/{poolID} -> CLPool
	CLPositionPrefix  = "clpos/"      // clpos/{positionID} -> CLPosition
	TickInfoPrefix    = "tick/"       // tick/{poolID}/{tickIndex} -> TickInfo
	PosByOwnerPrefix  = "clpos_own/"  // clpos_own/{owner}/{positionID} -> exists
	PosByPoolPrefix   = "clpos_pool/" // clpos_pool/{poolID}/{positionID} -> exists
	NextCLPoolIDKey   = "next_clpool_id"
	NextCLPositionIDKey = "next_clpos_id"
)

func CLPoolKey(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(CLPoolPrefix), bz...)
}

func CLPositionKey(positionID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, positionID)
	return append([]byte(CLPositionPrefix), bz...)
}

func TickInfoKey(poolID uint64, tickIndex int64) []byte {
	poolBz := make([]byte, 8)
	binary.BigEndian.PutUint64(poolBz, poolID)
	// Encode tick as uint64 with offset so negative ticks sort correctly
	tickBz := make([]byte, 8)
	binary.BigEndian.PutUint64(tickBz, uint64(tickIndex+1<<62))
	return append(append([]byte(TickInfoPrefix), poolBz...), tickBz...)
}

func TickInfoPoolPrefix(poolID uint64) []byte {
	poolBz := make([]byte, 8)
	binary.BigEndian.PutUint64(poolBz, poolID)
	return append([]byte(TickInfoPrefix), poolBz...)
}

func PosByOwnerKey(owner string, positionID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, positionID)
	return append(append([]byte(PosByOwnerPrefix), []byte(owner+"/")...), bz...)
}

func PosByOwnerPrefixKey(owner string) []byte {
	return []byte(PosByOwnerPrefix + owner + "/")
}

func PosByPoolKey(poolID uint64, positionID uint64) []byte {
	poolBz := make([]byte, 8)
	binary.BigEndian.PutUint64(poolBz, poolID)
	posBz := make([]byte, 8)
	binary.BigEndian.PutUint64(posBz, positionID)
	return append(append([]byte(PosByPoolPrefix), poolBz...), posBz...)
}
