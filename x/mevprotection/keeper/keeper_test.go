package keeper_test

import (
	"context"
	"crypto/sha256"
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"

	"syreen/x/mevprotection/keeper"
	"syreen/x/mevprotection/types"
)

// ---------------------------------------------------------------------------
// mock keepers
// ---------------------------------------------------------------------------

type mockStakingKeeper struct {
	// bondedCount controls how many bonded validators GetBondedValidatorsByPower reports.
	bondedCount int
	// slashCalls / jailCalls track enforcement actions for assertions.
	slashCalls *int
}

func (m *mockStakingKeeper) GetValidator(_ context.Context, _ sdk.ValAddress) (stakingtypes.Validator, error) {
	return stakingtypes.Validator{}, nil
}

func (m *mockStakingKeeper) GetValidatorByConsAddr(_ context.Context, _ sdk.ConsAddress) (stakingtypes.Validator, error) {
	// Return a validator with non-zero tokens so GetConsensusPower > 0.
	return stakingtypes.Validator{
		Tokens: sdk.DefaultPowerReduction.MulRaw(10),
		Status: stakingtypes.Bonded,
	}, nil
}

func (m *mockStakingKeeper) Slash(_ context.Context, _ sdk.ConsAddress, _ int64, _ int64, _ math.LegacyDec) (math.Int, error) {
	if m.slashCalls != nil {
		*m.slashCalls++
	}
	return math.ZeroInt(), nil
}

func (m *mockStakingKeeper) GetBondedValidatorsByPower(_ context.Context) ([]stakingtypes.Validator, error) {
	vals := make([]stakingtypes.Validator, m.bondedCount)
	return vals, nil
}

type mockBankKeeper struct{}

func (m *mockBankKeeper) MintCoins(_ context.Context, _ string, _ sdk.Coins) error {
	return nil
}

func (m *mockBankKeeper) SendCoinsFromModuleToModule(_ context.Context, _, _ string, _ sdk.Coins) error {
	return nil
}

func (m *mockBankKeeper) SendCoinsFromModuleToAccount(_ context.Context, _ string, _ sdk.AccAddress, _ sdk.Coins) error {
	return nil
}

type mockSlashingKeeper struct {
	jailCalls *int
}

func (m *mockSlashingKeeper) Slash(_ context.Context, _ sdk.ConsAddress, _ math.LegacyDec, _ int64, _ int64) error {
	return nil
}

func (m *mockSlashingKeeper) Jail(_ context.Context, _ sdk.ConsAddress) error {
	if m.jailCalls != nil {
		*m.jailCalls++
	}
	return nil
}

// ---------------------------------------------------------------------------
// setup
// ---------------------------------------------------------------------------

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context) {
	t.Helper()
	k, ctx, _, _ := setupKeeperWithMocks(t, &mockStakingKeeper{}, &mockSlashingKeeper{})
	return k, ctx
}

// setupKeeperWithMocks builds a keeper with caller-provided staking/slashing
// mocks so tests can inspect enforcement behaviour.
func setupKeeperWithMocks(t *testing.T, sk *mockStakingKeeper, slk *mockSlashingKeeper) (keeper.Keeper, sdk.Context, *mockStakingKeeper, *mockSlashingKeeper) {
	t.Helper()
	storeKey := storetypes.NewKVStoreKey("mevprotection")
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{
		Height: 10,
		Time:   time.Now(),
	}, false, log.NewNopLogger())
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	storeService := runtime.NewKVStoreService(storeKey)

	k := keeper.NewKeeper(cdc, storeService, sk, slk, &mockBankKeeper{}, "authority")
	return k, ctx, sk, slk
}

// commitRevealReorderedBlock sets up commit/reveal records with ascending
// timestamps and returns a fully-reversed block ordering that CheckFairOrdering
// flags as MEV (drift exceeds MaxTxDelay=1). Shared by MEV detection tests.
func commitRevealReorderedBlock(t *testing.T, k keeper.Keeper, ctx sdk.Context) [][]byte {
	t.Helper()
	params := types.DefaultParams()
	params.FairOrderConfig.MaxTxDelay = 1
	require.NoError(t, k.SetParams(ctx, params))

	bodies := make([][]byte, 3)
	nonces := make([][]byte, 3)
	hashes := make([][]byte, 3)
	for i := 0; i < 3; i++ {
		bodies[i] = []byte{byte('A' + i), byte(i), byte(i + 10)}
		nonces[i] = []byte{byte(i + 100)}
		hashes[i] = makeCommitHash(bodies[i], nonces[i])
	}
	for i := 0; i < 3; i++ {
		ctxT := ctx.WithBlockTime(time.Unix(int64(1000+i*100), 0))
		require.NoError(t, k.CommitTx(ctxT, validAddr(), hashes[i], []byte("enc")))
	}
	revealCtx := ctx.WithBlockHeight(20) // within reveal window [20,31]
	for i := 0; i < 3; i++ {
		require.NoError(t, k.RevealTx(revealCtx, validAddr(), hashes[i], bodies[i], nonces[i]))
	}
	return [][]byte{bodies[2], bodies[1], bodies[0]}
}

func validAddr() string {
	addr := sdk.AccAddress([]byte("test_address_padded_"))
	return addr.String()
}

func makeCommitHash(body, nonce []byte) []byte {
	h := sha256.Sum256(append(body, nonce...))
	return h[:]
}

// ---------------------------------------------------------------------------
// CommitTx
// ---------------------------------------------------------------------------

func TestCommitTx_Valid(t *testing.T) {
	k, ctx := setupKeeper(t)
	body := []byte("tx-body")
	nonce := []byte("nonce")
	hash := makeCommitHash(body, nonce)

	err := k.CommitTx(ctx, validAddr(), hash, []byte("encrypted"))
	require.NoError(t, err)

	committed, found := k.GetCommittedTx(ctx, hash)
	require.True(t, found)
	require.Equal(t, validAddr(), committed.Sender)
	require.Equal(t, int64(10), committed.BlockHeight)
}

func TestCommitTx_DuplicateCommit(t *testing.T) {
	k, ctx := setupKeeper(t)
	hash := makeCommitHash([]byte("body"), []byte("nonce"))

	err := k.CommitTx(ctx, validAddr(), hash, []byte("enc"))
	require.NoError(t, err)

	err = k.CommitTx(ctx, validAddr(), hash, []byte("enc2"))
	require.ErrorIs(t, err, types.ErrDuplicateCommit)
}

func TestCommitTx_InvalidHashLength(t *testing.T) {
	k, ctx := setupKeeper(t)
	err := k.CommitTx(ctx, validAddr(), []byte("short"), []byte("enc"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "32 bytes")
}

// ---------------------------------------------------------------------------
// RevealTx
// ---------------------------------------------------------------------------

func TestRevealTx_Valid(t *testing.T) {
	k, ctx := setupKeeper(t)
	body := []byte("tx-body-reveal")
	nonce := []byte("nonce-reveal")
	hash := makeCommitHash(body, nonce)

	// Commit first
	require.NoError(t, k.CommitTx(ctx, validAddr(), hash, []byte("encrypted")))

	// Default params: CommitWindow=20, RevealWindow=1 → reveal window is
	// [10+20/2, 10+20+1] = [20, 31]. Reveal at height 20 (start of window).
	err := k.RevealTx(ctx.WithBlockHeight(20), validAddr(), hash, body, nonce)
	require.NoError(t, err)
}

func TestRevealTx_Mismatch(t *testing.T) {
	k, ctx := setupKeeper(t)
	body := []byte("correct-body")
	nonce := []byte("correct-nonce")
	hash := makeCommitHash(body, nonce)

	require.NoError(t, k.CommitTx(ctx, validAddr(), hash, []byte("enc")))

	// Reveal with wrong body
	err := k.RevealTx(ctx, validAddr(), hash, []byte("wrong-body"), nonce)
	require.ErrorIs(t, err, types.ErrRevealMismatch)
}

func TestRevealTx_WrongSender(t *testing.T) {
	k, ctx := setupKeeper(t)
	body := []byte("body-sender")
	nonce := []byte("nonce-sender")
	hash := makeCommitHash(body, nonce)

	require.NoError(t, k.CommitTx(ctx, validAddr(), hash, []byte("enc")))

	otherAddr := sdk.AccAddress([]byte("other_address_paddd_")).String()
	err := k.RevealTx(ctx, otherAddr, hash, body, nonce)
	require.ErrorIs(t, err, types.ErrRevealMismatch)
}

func TestRevealTx_Expired(t *testing.T) {
	k, ctx := setupKeeper(t)
	body := []byte("body-expire")
	nonce := []byte("nonce-expire")
	hash := makeCommitHash(body, nonce)

	// Commit at height 10
	require.NoError(t, k.CommitTx(ctx, validAddr(), hash, []byte("enc")))

	// Default params: CommitWindow=20, RevealWindow=1, so maxRevealHeight = 10+20+1 = 31
	// Move to height 32 (past the window)
	futureCtx := ctx.WithBlockHeight(32)
	err := k.RevealTx(futureCtx, validAddr(), hash, body, nonce)
	require.ErrorIs(t, err, types.ErrRevealExpired)
}

func TestRevealTx_NotFound(t *testing.T) {
	k, ctx := setupKeeper(t)
	hash := makeCommitHash([]byte("nonexistent"), []byte("nonce"))
	err := k.RevealTx(ctx, validAddr(), hash, []byte("nonexistent"), []byte("nonce"))
	require.ErrorIs(t, err, types.ErrRevealMismatch)
}

// ---------------------------------------------------------------------------
// GetCommittedTx / SetCommittedTx roundtrip
// ---------------------------------------------------------------------------

func TestSetGetCommittedTx_Roundtrip(t *testing.T) {
	k, ctx := setupKeeper(t)
	hash := makeCommitHash([]byte("rt-body"), []byte("rt-nonce"))

	committed := types.CommittedTx{
		Sender:      validAddr(),
		TxHash:      hash,
		EncryptedTx: []byte("encrypted-rt"),
		BlockHeight: 42,
		Timestamp:   1700000000,
	}
	k.SetCommittedTx(ctx, hash, committed)

	got, found := k.GetCommittedTx(ctx, hash)
	require.True(t, found)
	require.Equal(t, committed.Sender, got.Sender)
	require.Equal(t, committed.BlockHeight, got.BlockHeight)
	require.Equal(t, committed.Timestamp, got.Timestamp)
	require.Equal(t, committed.EncryptedTx, got.EncryptedTx)
}

func TestGetCommittedTx_NotFound(t *testing.T) {
	k, ctx := setupKeeper(t)
	_, found := k.GetCommittedTx(ctx, []byte("nonexistent-hash-padded-to-32bb"))
	require.False(t, found)
}

// ---------------------------------------------------------------------------
// GetPendingReveals
// ---------------------------------------------------------------------------

func TestGetPendingReveals(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Commit two txs at height 10
	body1 := []byte("pending-body-1")
	nonce1 := []byte("pending-nonce-1")
	hash1 := makeCommitHash(body1, nonce1)
	require.NoError(t, k.CommitTx(ctx, validAddr(), hash1, []byte("enc1")))

	body2 := []byte("pending-body-2")
	nonce2 := []byte("pending-nonce-2")
	hash2 := makeCommitHash(body2, nonce2)
	require.NoError(t, k.CommitTx(ctx, validAddr(), hash2, []byte("enc2")))

	// Reveal only the first one (at height 20, within the reveal window [20,31])
	require.NoError(t, k.RevealTx(ctx.WithBlockHeight(20), validAddr(), hash1, body1, nonce1))

	// GetPendingReveals should return only the second (unrevealed) one
	pending := k.GetPendingReveals(ctx)
	require.Len(t, pending, 1)
	require.Equal(t, hash2, pending[0].TxHash)
}

func TestGetPendingReveals_ExpiredNotReturned(t *testing.T) {
	k, ctx := setupKeeper(t)

	body := []byte("expire-pending")
	nonce := []byte("expire-nonce")
	hash := makeCommitHash(body, nonce)
	require.NoError(t, k.CommitTx(ctx, validAddr(), hash, []byte("enc")))

	// Move past the reveal window
	futureCtx := ctx.WithBlockHeight(100)
	pending := k.GetPendingReveals(futureCtx)
	require.Empty(t, pending)
}

// ---------------------------------------------------------------------------
// PruneExpiredCommits
// ---------------------------------------------------------------------------

func TestPruneExpiredCommits(t *testing.T) {
	k, ctx := setupKeeper(t)

	body := []byte("prune-body")
	nonce := []byte("prune-nonce")
	hash := makeCommitHash(body, nonce)
	require.NoError(t, k.CommitTx(ctx, validAddr(), hash, []byte("enc")))

	// Verify it exists
	_, found := k.GetCommittedTx(ctx, hash)
	require.True(t, found)

	// Prune at a height past the reveal window (committed at 10, window = 20+1 = 21, so 32 is past)
	futureCtx := ctx.WithBlockHeight(32)
	k.PruneExpiredCommits(futureCtx)

	// Should be gone
	_, found = k.GetCommittedTx(futureCtx, hash)
	require.False(t, found)
}

func TestPruneExpiredCommits_KeepsActiveCommits(t *testing.T) {
	k, ctx := setupKeeper(t)

	body := []byte("active-body")
	nonce := []byte("active-nonce")
	hash := makeCommitHash(body, nonce)
	require.NoError(t, k.CommitTx(ctx, validAddr(), hash, []byte("enc")))

	// Prune at a height within the window (height 11, max = 14)
	sameCtx := ctx.WithBlockHeight(11)
	k.PruneExpiredCommits(sameCtx)

	// Should still exist
	_, found := k.GetCommittedTx(sameCtx, hash)
	require.True(t, found)
}

// ---------------------------------------------------------------------------
// SetParams / GetParams roundtrip
// ---------------------------------------------------------------------------

func TestSetGetParams_Roundtrip(t *testing.T) {
	k, ctx := setupKeeper(t)

	params := types.Params{
		FairOrderConfig: types.FairOrderConfig{
			EnableCommitReveal: false,
			CommitWindow:       5,
			RevealWindow:       10,
			MaxTxDelay:         20,
		},
	}
	require.NoError(t, k.SetParams(ctx, params))

	got := k.GetParams(ctx)
	require.Equal(t, params.FairOrderConfig.EnableCommitReveal, got.FairOrderConfig.EnableCommitReveal)
	require.Equal(t, params.FairOrderConfig.CommitWindow, got.FairOrderConfig.CommitWindow)
	require.Equal(t, params.FairOrderConfig.RevealWindow, got.FairOrderConfig.RevealWindow)
	require.Equal(t, params.FairOrderConfig.MaxTxDelay, got.FairOrderConfig.MaxTxDelay)
}

func TestGetParams_DefaultWhenNotSet(t *testing.T) {
	k, ctx := setupKeeper(t)
	got := k.GetParams(ctx)
	require.Equal(t, types.DefaultParams(), got)
}

// ---------------------------------------------------------------------------
// Genesis InitGenesis / ExportGenesis roundtrip
// ---------------------------------------------------------------------------

func TestGenesis_InitExport_Roundtrip(t *testing.T) {
	k, ctx := setupKeeper(t)

	genesisIn := types.GenesisState{
		Params: types.Params{
			FairOrderConfig: types.FairOrderConfig{
				EnableCommitReveal: true,
				CommitWindow:       3,
				RevealWindow:       5,
				MaxTxDelay:         7,
			},
		},
		Penalties: []types.MEVPenalty{
			{
				ValidatorAddr: "val1",
				PenaltyType:   "reordering",
				Amount:        10,
				Evidence:      []byte("proof"),
				BlockHeight:   50,
			},
		},
	}

	k.InitGenesis(ctx, genesisIn)

	exported := k.ExportGenesis(ctx)
	require.NotNil(t, exported)
	// Params should roundtrip
	require.Equal(t, genesisIn.Params.FairOrderConfig.CommitWindow, exported.Params.FairOrderConfig.CommitWindow)
	require.Equal(t, genesisIn.Params.FairOrderConfig.RevealWindow, exported.Params.FairOrderConfig.RevealWindow)
	require.Equal(t, genesisIn.Params.FairOrderConfig.MaxTxDelay, exported.Params.FairOrderConfig.MaxTxDelay)
	require.Equal(t, genesisIn.Params.FairOrderConfig.EnableCommitReveal, exported.Params.FairOrderConfig.EnableCommitReveal)
}

func TestGenesis_DefaultInit(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.InitGenesis(ctx, *types.DefaultGenesis())

	exported := k.ExportGenesis(ctx)
	require.Equal(t, types.DefaultParams(), exported.Params)
	require.Empty(t, exported.Penalties)
}

// ---------------------------------------------------------------------------
// CheckFairOrdering
// ---------------------------------------------------------------------------

func TestCheckFairOrdering_DisabledWhenCommitRevealOff(t *testing.T) {
	k, ctx := setupKeeper(t)

	params := types.DefaultParams()
	params.FairOrderConfig.EnableCommitReveal = false
	require.NoError(t, k.SetParams(ctx, params))

	// Even with multiple txs, should return nil when disabled
	err := k.CheckFairOrdering(ctx, [][]byte{[]byte("tx1"), []byte("tx2")})
	require.NoError(t, err)
}

func TestCheckFairOrdering_SingleTx(t *testing.T) {
	k, ctx := setupKeeper(t)
	// Single tx should always pass
	err := k.CheckFairOrdering(ctx, [][]byte{[]byte("only-tx")})
	require.NoError(t, err)
}

func TestCheckFairOrdering_ValidOrdering(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Create two commits with distinct timestamps (earlier commit first in block)
	body1 := []byte("fair-tx-1")
	nonce1 := []byte("fair-nonce-1")
	hash1 := makeCommitHash(body1, nonce1)

	body2 := []byte("fair-tx-2")
	nonce2 := []byte("fair-nonce-2")
	hash2 := makeCommitHash(body2, nonce2)

	// Commit first tx at an earlier time
	ctx1 := ctx.WithBlockTime(time.Unix(1000, 0))
	require.NoError(t, k.CommitTx(ctx1, validAddr(), hash1, []byte("enc1")))

	// Commit second tx at a later time
	ctx2 := ctx.WithBlockTime(time.Unix(2000, 0))
	require.NoError(t, k.CommitTx(ctx2, validAddr(), hash2, []byte("enc2")))

	// Reveal both (at height 20, within the reveal window [20,31])
	revealCtx := ctx.WithBlockHeight(20)
	require.NoError(t, k.RevealTx(revealCtx, validAddr(), hash1, body1, nonce1))
	require.NoError(t, k.RevealTx(revealCtx, validAddr(), hash2, body2, nonce2))

	// Block with correct ordering (body1 first, body2 second)
	err := k.CheckFairOrdering(ctx, [][]byte{body1, body2})
	require.NoError(t, err)
}

func TestCheckFairOrdering_MEVDetectedOnReordering(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Use MaxTxDelay = 1 to make MEV detection more sensitive
	params := types.DefaultParams()
	params.FairOrderConfig.MaxTxDelay = 1
	require.NoError(t, k.SetParams(ctx, params))

	// We need enough txs that swapping them exceeds MaxTxDelay
	// Create 3 commits with ascending timestamps
	bodies := make([][]byte, 3)
	nonces := make([][]byte, 3)
	hashes := make([][]byte, 3)

	for i := 0; i < 3; i++ {
		bodies[i] = []byte{byte('A' + i), byte(i), byte(i + 10)}
		nonces[i] = []byte{byte(i + 100)}
		hashes[i] = makeCommitHash(bodies[i], nonces[i])
	}

	for i := 0; i < 3; i++ {
		ctxT := ctx.WithBlockTime(time.Unix(int64(1000+i*100), 0))
		require.NoError(t, k.CommitTx(ctxT, validAddr(), hashes[i], []byte("enc")))
	}

	revealCtx := ctx.WithBlockHeight(20)
	for i := 0; i < 3; i++ {
		require.NoError(t, k.RevealTx(revealCtx, validAddr(), hashes[i], bodies[i], nonces[i]))
	}

	// Reverse the order completely: [body2, body1, body0] -- drift = 2, exceeds MaxTxDelay=1
	reordered := [][]byte{bodies[2], bodies[1], bodies[0]}
	err := k.CheckFairOrdering(ctx, reordered)
	require.ErrorIs(t, err, types.ErrMEVDetected)
}

func TestCheckFairOrdering_NoReveals(t *testing.T) {
	k, ctx := setupKeeper(t)
	// Block txs with no matching reveals -- should pass (no evidence of MEV)
	err := k.CheckFairOrdering(ctx, [][]byte{[]byte("random1"), []byte("random2")})
	require.NoError(t, err)
}

// ---------------------------------------------------------------------------
// C-07: Deterministic ordering tests
// ---------------------------------------------------------------------------

// TestGetAllRevealedTxs_DeterministicOrder inserts revealed txs and verifies
// that GetAllRevealedTxs always returns them sorted by commit hash,
// guaranteeing determinism regardless of DB backend iteration order.
func TestGetAllRevealedTxs_DeterministicOrder(t *testing.T) {
	k, ctx := setupKeeper(t)

	type testCase struct {
		body  []byte
		nonce []byte
	}
	// Deliberately unordered test data
	cases := []testCase{
		{body: []byte("tx_zzz_padding_____"), nonce: []byte("nonce1")},
		{body: []byte("tx_aaa_padding_____"), nonce: []byte("nonce2")},
		{body: []byte("tx_mmm_padding_____"), nonce: []byte("nonce3")},
		{body: []byte("tx_bbb_padding_____"), nonce: []byte("nonce4")},
		{body: []byte("tx_yyy_padding_____"), nonce: []byte("nonce5")},
	}

	sender := validAddr()
	for _, tc := range cases {
		h := sha256.New()
		h.Write(tc.body)
		h.Write(tc.nonce)
		commitHash := h.Sum(nil)

		require.NoError(t, k.CommitTx(ctx, sender, commitHash, []byte("encrypted")))
		// Reveal at height 20, within the reveal window [20,31].
		require.NoError(t, k.RevealTx(ctx.WithBlockHeight(20), sender, commitHash, tc.body, tc.nonce))
	}

	results := k.GetAllRevealedTxs(ctx)
	require.Len(t, results, 5)

	// Verify sorted by CommittedHash
	for i := 1; i < len(results); i++ {
		require.LessOrEqual(t,
			string(results[i-1].Revealed.CommittedHash),
			string(results[i].Revealed.CommittedHash),
			"results[%d] should be <= results[%d] by CommittedHash", i-1, i)
	}

	// Determinism: multiple calls yield identical ordering
	for iter := 0; iter < 20; iter++ {
		results2 := k.GetAllRevealedTxs(ctx)
		require.Len(t, results2, 5)
		for i := range results {
			require.Equal(t, results[i].Revealed.CommittedHash, results2[i].Revealed.CommittedHash,
				"iteration %d: result[%d] commit hash differs", iter, i)
		}
	}
}

// TestGetPendingReveals_DeterministicOrder verifies that pending reveals are
// returned sorted by tx hash for determinism.
func TestGetPendingReveals_DeterministicOrder(t *testing.T) {
	k, ctx := setupKeeper(t)

	sender := validAddr()
	for i := 0; i < 5; i++ {
		h := sha256.Sum256([]byte{byte(i), byte(i + 100), byte(i + 200)})
		txHash := h[:]
		require.NoError(t, k.CommitTx(ctx, sender, txHash, []byte("encrypted")))
	}

	pending := k.GetPendingReveals(ctx)
	require.Len(t, pending, 5)

	// Verify sorted by TxHash
	for i := 1; i < len(pending); i++ {
		require.LessOrEqual(t,
			string(pending[i-1].TxHash),
			string(pending[i].TxHash),
			"pending[%d] should be <= pending[%d] by TxHash", i-1, i)
	}

	// Determinism across calls
	for iter := 0; iter < 20; iter++ {
		pending2 := k.GetPendingReveals(ctx)
		require.Len(t, pending2, 5)
		for i := range pending {
			require.Equal(t, pending[i].TxHash, pending2[i].TxHash,
				"iteration %d: pending[%d] hash differs", iter, i)
		}
	}
}

// ---------------------------------------------------------------------------
// DetectMEV solo-validator safety guard
// ---------------------------------------------------------------------------

// TestDetectMEV_SkipsSlashingWhenTooFewValidators proves that when the bonded
// validator count is below the safety threshold, DetectMEV records the penalty
// but does NOT slash or jail the proposer — protecting a low-validator-count
// chain from being halted by its own MEV enforcement.
func TestDetectMEV_SkipsSlashingWhenTooFewValidators(t *testing.T) {
	slashCalls := 0
	jailCalls := 0
	sk := &mockStakingKeeper{bondedCount: 1, slashCalls: &slashCalls}
	slk := &mockSlashingKeeper{jailCalls: &jailCalls}
	k, ctx, _, _ := setupKeeperWithMocks(t, sk, slk)

	reordered := commitRevealReorderedBlock(t, k, ctx)

	penalty, err := k.DetectMEV(ctx, reordered)
	require.NoError(t, err)
	// A penalty is still recorded as evidence...
	require.NotNil(t, penalty)
	require.Equal(t, "reordering", penalty.PenaltyType)
	// ...but the validator set must NOT be touched.
	require.Equal(t, 0, slashCalls, "must not slash when validator count below threshold")
	require.Equal(t, 0, jailCalls, "must not jail when validator count below threshold")
}

// TestDetectMEV_SkipsSlashingWhenJailingDropsBelowThreshold proves the guard
// also fires when the bonded count equals the threshold, since jailing the
// proposer would drop the active set below the safe minimum.
func TestDetectMEV_SkipsSlashingWhenJailingDropsBelowThreshold(t *testing.T) {
	slashCalls := 0
	jailCalls := 0
	sk := &mockStakingKeeper{bondedCount: 4, slashCalls: &slashCalls}
	slk := &mockSlashingKeeper{jailCalls: &jailCalls}
	k, ctx, _, _ := setupKeeperWithMocks(t, sk, slk)

	reordered := commitRevealReorderedBlock(t, k, ctx)

	penalty, err := k.DetectMEV(ctx, reordered)
	require.NoError(t, err)
	require.NotNil(t, penalty)
	require.Equal(t, 0, slashCalls, "must not slash when jailing would drop below threshold")
	require.Equal(t, 0, jailCalls, "must not jail when jailing would drop below threshold")
}

// TestDetectMEV_SlashesWhenEnoughValidators confirms enforcement still runs
// when the validator set is comfortably above the safety threshold.
func TestDetectMEV_SlashesWhenEnoughValidators(t *testing.T) {
	slashCalls := 0
	jailCalls := 0
	sk := &mockStakingKeeper{bondedCount: 10, slashCalls: &slashCalls}
	slk := &mockSlashingKeeper{jailCalls: &jailCalls}
	k, ctx, _, _ := setupKeeperWithMocks(t, sk, slk)

	reordered := commitRevealReorderedBlock(t, k, ctx)

	penalty, err := k.DetectMEV(ctx, reordered)
	require.NoError(t, err)
	require.NotNil(t, penalty)
	require.Equal(t, 1, slashCalls, "should slash when validator count is safely above threshold")
	require.Equal(t, 1, jailCalls, "should jail when validator count is safely above threshold")
}

// TestDetectMEV_NoViolationReturnsNil confirms DetectMEV is a no-op when there
// is no MEV violation.
func TestDetectMEV_NoViolationReturnsNil(t *testing.T) {
	slashCalls := 0
	jailCalls := 0
	sk := &mockStakingKeeper{bondedCount: 10, slashCalls: &slashCalls}
	slk := &mockSlashingKeeper{jailCalls: &jailCalls}
	k, ctx, _, _ := setupKeeperWithMocks(t, sk, slk)

	penalty, err := k.DetectMEV(ctx, [][]byte{[]byte("only-one-tx")})
	require.NoError(t, err)
	require.Nil(t, penalty)
	require.Equal(t, 0, slashCalls)
	require.Equal(t, 0, jailCalls)
}
