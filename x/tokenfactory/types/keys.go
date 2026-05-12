package types

const (
	ModuleName = "tokenfactory"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	// DenomsPrefixKey is the prefix for storing denom authority metadata
	DenomsPrefixKey = "denoms/"

	// CreatorPrefixKey is the prefix for indexing denoms by creator
	CreatorPrefixKey = "creator/"
)

// DenomAuthorityKey returns the store key for a denom's authority metadata
func DenomAuthorityKey(denom string) []byte {
	return append([]byte(DenomsPrefixKey), []byte(denom)...)
}

// CreatorDenomsKey returns the prefix key for a creator's denoms
func CreatorDenomsKey(creator string) []byte {
	return append([]byte(CreatorPrefixKey), []byte(creator)...)
}
