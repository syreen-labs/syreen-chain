package types

import "encoding/binary"

const (
	ModuleName = "options"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	OptionPrefix   = "option/"   // option/{id} → Option
	NextOptionIDKey = "next_option_id"

	// Index: active options by expiry block for fast BeginBlock sweep
	ExpiryIndexPrefix = "expiry/" // expiry/{block}/{id} → exists
)

func OptionKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return append([]byte(OptionPrefix), bz...)
}

func ExpiryIndexKey(block int64, id uint64) []byte {
	blockBz := make([]byte, 8)
	binary.BigEndian.PutUint64(blockBz, uint64(block))
	idBz := make([]byte, 8)
	binary.BigEndian.PutUint64(idBz, id)
	return append(append([]byte(ExpiryIndexPrefix), blockBz...), idBz...)
}

func ExpiryIndexPrefixForBlock(block int64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(block))
	return append([]byte(ExpiryIndexPrefix), bz...)
}
