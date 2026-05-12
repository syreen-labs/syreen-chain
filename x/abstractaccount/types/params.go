package types

import (
	"fmt"
	"time"
)

var (
	DefaultMaxSessionKeyDuration = 24 * time.Hour
	DefaultMaxGuardians          = uint32(7)
	DefaultRecoveryDelayPeriod   = 48 * time.Hour
	DefaultMaxBatchSize          = uint32(10)
	DefaultEnableGasSponsorship  = true
)

// Params defines the parameters for the abstractaccount module
type Params struct {
	MaxSessionKeyDuration time.Duration `json:"max_session_key_duration"`
	MaxGuardians          uint32        `json:"max_guardians"`
	RecoveryDelayPeriod   time.Duration `json:"recovery_delay_period"`
	MaxBatchSize          uint32        `json:"max_batch_size"`
	EnableGasSponsorship  bool          `json:"enable_gas_sponsorship"`
}

func DefaultParams() Params {
	return Params{
		MaxSessionKeyDuration: DefaultMaxSessionKeyDuration,
		MaxGuardians:          DefaultMaxGuardians,
		RecoveryDelayPeriod:   DefaultRecoveryDelayPeriod,
		MaxBatchSize:          DefaultMaxBatchSize,
		EnableGasSponsorship:  DefaultEnableGasSponsorship,
	}
}

func (p Params) Validate() error {
	if p.MaxSessionKeyDuration <= 0 {
		return fmt.Errorf("max session key duration must be positive: %s", p.MaxSessionKeyDuration)
	}
	if p.MaxGuardians == 0 {
		return fmt.Errorf("max guardians must be positive")
	}
	if p.RecoveryDelayPeriod <= 0 {
		return fmt.Errorf("recovery delay period must be positive: %s", p.RecoveryDelayPeriod)
	}
	if p.MaxBatchSize == 0 {
		return fmt.Errorf("max batch size must be positive")
	}
	return nil
}
