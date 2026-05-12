package types

import "fmt"

const (
	DefaultMaxWasmCodeSize  uint64 = 800 * 1024 // 800 KB
	DefaultMaxContractGas   uint64 = 10_000_000_000
)

// Params defines the parameters for the compute module
type Params struct {
	CodeUploadAccess             AccessConfig `json:"code_upload_access"`
	InstantiateDefaultPermission string       `json:"instantiate_default_permission"`
	MaxWasmCodeSize              uint64       `json:"max_wasm_code_size"`
	MaxContractGas               uint64       `json:"max_contract_gas"`
}

func DefaultParams() Params {
	return Params{
		CodeUploadAccess: AccessConfig{
			Permission: AccessTypeEverybody,
		},
		InstantiateDefaultPermission: AccessTypeEverybody,
		MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
		MaxContractGas:               DefaultMaxContractGas,
	}
}

func (p Params) Validate() error {
	if p.MaxWasmCodeSize == 0 {
		return fmt.Errorf("max wasm code size must be positive")
	}
	if p.MaxContractGas == 0 {
		return fmt.Errorf("max contract gas must be positive")
	}
	switch p.InstantiateDefaultPermission {
	case AccessTypeEverybody, AccessTypeOnlyAddress, AccessTypeNobody:
		// valid
	default:
		return fmt.Errorf("invalid instantiate default permission: %s", p.InstantiateDefaultPermission)
	}
	switch p.CodeUploadAccess.Permission {
	case AccessTypeEverybody, AccessTypeOnlyAddress, AccessTypeNobody:
		// valid
	default:
		return fmt.Errorf("invalid code upload access permission: %s", p.CodeUploadAccess.Permission)
	}
	return nil
}
