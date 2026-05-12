package types

import "encoding/binary"

const (
	ModuleName = "launchpad"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	LaunchPrefix         = "launch/"        // launch/{launchID} -> Launch
	ContributionPrefix   = "contrib/"       // contrib/{launchID}/{address} -> Contribution
	ContribByAddrPrefix  = "caddr/"         // caddr/{address}/{launchID} -> exists
	VestingPrefix        = "vest/"          // vest/{launchID}/{address} -> VestingPosition
	NextLaunchIDKey      = "next_launch_id"
	ActiveLaunchPrefix   = "active/"        // active/{launchID} -> exists (index for BeginBlock)
)

func LaunchKey(launchID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, launchID)
	return append([]byte(LaunchPrefix), bz...)
}

func ContributionKey(launchID uint64, address string) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, launchID)
	return append(append([]byte(ContributionPrefix), bz...), []byte("/"+address)...)
}

func ContributionLaunchPrefix(launchID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, launchID)
	return append([]byte(ContributionPrefix), bz...)
}

func ContribByAddrKey(address string, launchID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, launchID)
	return append(append([]byte(ContribByAddrPrefix), []byte(address+"/")...), bz...)
}

func VestingKey(launchID uint64, address string) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, launchID)
	return append(append([]byte(VestingPrefix), bz...), []byte("/"+address)...)
}

func ActiveLaunchKey(launchID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, launchID)
	return append([]byte(ActiveLaunchPrefix), bz...)
}
