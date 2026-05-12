package types_test

import (
	"testing"

	"cosmossdk.io/math"
	gogoproto "github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	"syreen/x/intent/types"
)

// TestMsgCreateStrategy_ProtoRoundTrip asserts that MsgCreateStrategy can be
// marshaled to protobuf wire bytes and unmarshaled back with all fields
// preserved. This is the minimum contract needed for the tx decoder to accept
// the message on broadcast.
func TestMsgCreateStrategy_ProtoRoundTrip(t *testing.T) {
	original := &types.MsgCreateStrategy{
		Creator:        "syreen1examplecreator000000000000000000000",
		TemplateName:   types.StrategySafeAccumulate,
		PoolID:         1,
		InputDenom:     "usyreen",
		OutputDenom:    "uusdc",
		TotalBudget:    math.NewInt(1_500_000),
		RiskLevel:      types.RiskModerate,
		ExpiryBlocks:   10_000,
		NumExecutions:  5,
		IntervalBlocks: 100,
		GridLevels:     0,
	}

	bz, err := gogoproto.Marshal(original)
	require.NoError(t, err)
	require.NotEmpty(t, bz)

	decoded := &types.MsgCreateStrategy{}
	require.NoError(t, gogoproto.Unmarshal(bz, decoded))

	require.Equal(t, original.Creator, decoded.Creator)
	require.Equal(t, original.TemplateName, decoded.TemplateName)
	require.Equal(t, original.PoolID, decoded.PoolID)
	require.Equal(t, original.InputDenom, decoded.InputDenom)
	require.Equal(t, original.OutputDenom, decoded.OutputDenom)
	require.True(t, original.TotalBudget.Equal(decoded.TotalBudget))
	require.Equal(t, original.RiskLevel, decoded.RiskLevel)
	require.Equal(t, original.ExpiryBlocks, decoded.ExpiryBlocks)
	require.Equal(t, original.NumExecutions, decoded.NumExecutions)
	require.Equal(t, original.IntervalBlocks, decoded.IntervalBlocks)
	require.Equal(t, original.GridLevels, decoded.GridLevels)
}

// TestMsgCreateStrategy_Descriptor asserts that the Descriptor() method
// returns a non-empty gzipped FileDescriptorProto — the SDK's unknownproto
// reject path requires this and will fail tx decoding without it.
func TestMsgCreateStrategy_Descriptor(t *testing.T) {
	m := &types.MsgCreateStrategy{}
	bz, idx := m.Descriptor()
	require.NotEmpty(t, bz)
	require.Equal(t, []int{0}, idx)

	c := &types.MsgCancelStrategy{}
	bz2, idx2 := c.Descriptor()
	require.NotEmpty(t, bz2)
	require.Equal(t, []int{2}, idx2)
}

// TestMsgCancelStrategy_ProtoRoundTrip verifies wire-format round-trip for
// MsgCancelStrategy.
func TestMsgCancelStrategy_ProtoRoundTrip(t *testing.T) {
	original := &types.MsgCancelStrategy{
		Creator:    "syreen1examplecreator000000000000000000000",
		StrategyID: "strategy-42",
	}

	bz, err := gogoproto.Marshal(original)
	require.NoError(t, err)
	require.NotEmpty(t, bz)

	decoded := &types.MsgCancelStrategy{}
	require.NoError(t, gogoproto.Unmarshal(bz, decoded))

	require.Equal(t, original.Creator, decoded.Creator)
	require.Equal(t, original.StrategyID, decoded.StrategyID)
}
