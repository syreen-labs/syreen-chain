package types

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

var (
	// DefaultEVMDenom is the denom used for EVM gas payments
	DefaultEVMDenom = "usyreen"
	// DefaultEnableCreate allows contract creation
	DefaultEnableCreate = true
	// DefaultEnableCall allows contract calls
	DefaultEnableCall = true
)

// Params defines the EVM module parameters
type Params struct {
	EvmDenom     string `json:"evm_denom"`
	EnableCreate bool   `json:"enable_create"`
	EnableCall   bool   `json:"enable_call"`
}

func DefaultParams() Params {
	return Params{
		EvmDenom:     DefaultEVMDenom,
		EnableCreate: DefaultEnableCreate,
		EnableCall:   DefaultEnableCall,
	}
}

func (p Params) Validate() error {
	if p.EvmDenom == "" {
		return fmt.Errorf("evm denom cannot be empty")
	}
	return nil
}

// DefaultChainConfig returns the Ethereum chain config for Syreen EVM.
// All forks enabled from genesis (block 0), using Merge/PoS rules.
func DefaultChainConfig() *params.ChainConfig {
	return &params.ChainConfig{
		ChainID:                       big.NewInt(79733), // Syreen EVM chain ID
		HomesteadBlock:                big.NewInt(0),
		DAOForkBlock:                  big.NewInt(0),
		DAOForkSupport:                true,
		EIP150Block:                   big.NewInt(0),
		EIP155Block:                   big.NewInt(0),
		EIP158Block:                   big.NewInt(0),
		ByzantiumBlock:                big.NewInt(0),
		ConstantinopleBlock:           big.NewInt(0),
		PetersburgBlock:               big.NewInt(0),
		IstanbulBlock:                 big.NewInt(0),
		MuirGlacierBlock:             big.NewInt(0),
		BerlinBlock:                   big.NewInt(0),
		LondonBlock:                   big.NewInt(0),
		ArrowGlacierBlock:            big.NewInt(0),
		GrayGlacierBlock:             big.NewInt(0),
		MergeNetsplitBlock:           big.NewInt(0),
		TerminalTotalDifficulty:       big.NewInt(0),
	}
}

// GenesisState defines the evm module genesis state
type GenesisState struct {
	Params   Params    `json:"params"`
	Accounts []Account `json:"accounts"`
}

// Account represents an EVM account in genesis
type Account struct {
	Address string            `json:"address"`         // hex address
	Code    string            `json:"code"`            // hex code
	Storage map[string]string `json:"storage"`         // hex key -> hex value
	Nonce   uint64            `json:"nonce,omitempty"` // EVM nonce
}

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:   DefaultParams(),
		Accounts: []Account{},
	}
}

func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
