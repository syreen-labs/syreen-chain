package types

import "encoding/binary"

const (
	ModuleName = "vault"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	VaultPrefix          = "vault/"        // vault/{vaultID} -> Vault
	VaultDepositPrefix   = "vdep/"         // vdep/{vaultID}/{address} -> VaultDeposit
	DepositByAddrPrefix  = "vdep_addr/"    // vdep_addr/{address}/{vaultID} -> exists
	NextVaultIDKey       = "next_vault_id"
)

func VaultKey(vaultID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, vaultID)
	return append([]byte(VaultPrefix), bz...)
}

func VaultDepositKey(vaultID uint64, address string) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, vaultID)
	return append(append([]byte(VaultDepositPrefix), bz...), []byte("/"+address)...)
}

func DepositByAddrKey(address string, vaultID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, vaultID)
	return append(append([]byte(DepositByAddrPrefix), []byte(address+"/")...), bz...)
}

func DepositByAddrPrefixKey(address string) []byte {
	return []byte(DepositByAddrPrefix + address + "/")
}

func VaultDepositPrefixKey(vaultID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, vaultID)
	return append([]byte(VaultDepositPrefix), bz...)
}
