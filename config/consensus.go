package config

import "time"

// ConsensusConfig defines Syreen's optimized consensus parameters
// for sub-second block times. These values are applied as defaults
// in initTendermintConfig() in cmd/syreend/cmd/root.go.
var ConsensusConfig = struct {
	TimeoutPropose    time.Duration
	TimeoutPrevote    time.Duration
	TimeoutPrecommit  time.Duration
	TimeoutCommit     time.Duration
	CreateEmptyBlocks bool
}{
	TimeoutPropose:    1 * time.Second,
	TimeoutPrevote:    500 * time.Millisecond,
	TimeoutPrecommit:  500 * time.Millisecond,
	TimeoutCommit:     500 * time.Millisecond,
	CreateEmptyBlocks: true,
}
