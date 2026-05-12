package types

const (
	ModuleName = "abstractaccount"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	// SmartAccountPrefix is the prefix for storing smart account data
	SmartAccountPrefix = "smart_account/"

	// SessionKeyPrefix is the prefix for storing session keys
	SessionKeyPrefix = "session_key/"

	// RecoveryConfigPrefix is the prefix for storing recovery configurations
	RecoveryConfigPrefix = "recovery_config/"

	// RecoveryRequestPrefix is the prefix for storing active recovery requests
	RecoveryRequestPrefix = "recovery_request/"

	// GasSponsorPrefix is the prefix for storing gas sponsorship records
	GasSponsorPrefix = "gas_sponsor/"

	// AccountNonceKey is the store key for the monotonic account creation nonce
	AccountNonceKey = "account_nonce"
)

// SmartAccountKey returns the store key for a smart account
func SmartAccountKey(address string) []byte {
	return append([]byte(SmartAccountPrefix), []byte(address)...)
}

// SessionKeyKey returns the store key for a session key
func SessionKeyKey(granter, grantee string) []byte {
	return append([]byte(SessionKeyPrefix), []byte(granter+"/"+grantee)...)
}

// SessionKeysByGranterPrefix returns the prefix for all session keys from a granter
func SessionKeysByGranterPrefix(granter string) []byte {
	return append([]byte(SessionKeyPrefix), []byte(granter+"/")...)
}

// RecoveryConfigKey returns the store key for a recovery config
func RecoveryConfigKey(account string) []byte {
	return append([]byte(RecoveryConfigPrefix), []byte(account)...)
}

// RecoveryRequestKey returns the store key for an active recovery request
func RecoveryRequestKey(account string) []byte {
	return append([]byte(RecoveryRequestPrefix), []byte(account)...)
}

// GasSponsorKey returns the store key for a gas sponsorship record
func GasSponsorKey(sponsor, sponsored string) []byte {
	return append([]byte(GasSponsorPrefix), []byte(sponsor+"/"+sponsored)...)
}

// GasSponsorsBySponsored returns the prefix for looking up sponsors by sponsored account
func GasSponsorsBySponsored(sponsored string) []byte {
	return append([]byte(GasSponsorPrefix+"by_sponsored/"), []byte(sponsored+"/")...)
}
