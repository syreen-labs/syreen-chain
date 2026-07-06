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

// CreatorDenomsPrefix returns the iteration prefix under which every denom for a
// given creator is indexed. Each denom is stored under its own key so that
// creation is O(1) instead of rewriting a single JSON slice on every create.
// A "|" separator is safe because bech32 creator addresses never contain it.
func CreatorDenomsPrefix(creator string) []byte {
	return append([]byte(CreatorPrefixKey), []byte(creator+"|")...)
}

// CreatorDenomIndexKey returns the store key for a single (creator, denom) index entry.
func CreatorDenomIndexKey(creator, denom string) []byte {
	return append(CreatorDenomsPrefix(creator), []byte(denom)...)
}
