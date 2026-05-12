package types

const (
	ModuleName = "identity"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
	QuerierRoute = ModuleName

	// Key prefixes for store
	IdentityKeyPrefix    = "identity/"
	VerifierKeyPrefix    = "verifier/"
	UserIndexKeyPrefix   = "user_index/"
	ParamsKey            = "params"
)

func IdentityStoreKey(address string) []byte {
	return []byte(IdentityKeyPrefix + address)
}

func VerifierStoreKey(address string) []byte {
	return []byte(VerifierKeyPrefix + address)
}

func UserIndexKey(address string) []byte {
	return []byte(UserIndexKeyPrefix + address)
}
