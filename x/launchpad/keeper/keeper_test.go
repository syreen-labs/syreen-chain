package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	"syreen/x/launchpad/keeper"
	"syreen/x/launchpad/types"

	dbm "github.com/cosmos/cosmos-db"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

// ============================================================
// Mock Keepers
// ============================================================

type mockAccountKeeper struct {
	moduleAddrs map[string]sdk.AccAddress
}

func (m mockAccountKeeper) GetModuleAddress(name string) sdk.AccAddress {
	if addr, ok := m.moduleAddrs[name]; ok {
		return addr
	}
	return authtypes.NewModuleAddress(name)
}

func (m mockAccountKeeper) GetModuleAccount(_ context.Context, name string) sdk.ModuleAccountI {
	return nil
}

type mockBankKeeper struct {
	balances map[string]map[string]math.Int // addr -> denom -> amount
}

func newMockBankKeeper() *mockBankKeeper {
	return &mockBankKeeper{balances: make(map[string]map[string]math.Int)}
}

func (m *mockBankKeeper) setBalance(addr, denom string, amount math.Int) {
	if m.balances[addr] == nil {
		m.balances[addr] = make(map[string]math.Int)
	}
	m.balances[addr][denom] = amount
}

func (m *mockBankKeeper) getBalance(addr, denom string) math.Int {
	if m.balances[addr] == nil {
		return math.ZeroInt()
	}
	if amt, ok := m.balances[addr][denom]; ok {
		return amt
	}
	return math.ZeroInt()
}

func (m *mockBankKeeper) GetBalance(_ context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	amt := m.getBalance(addr.String(), denom)
	return sdk.NewCoin(denom, amt)
}

func (m *mockBankKeeper) SendCoins(_ context.Context, from, to sdk.AccAddress, amt sdk.Coins) error {
	for _, coin := range amt {
		fromBal := m.getBalance(from.String(), coin.Denom)
		if fromBal.LT(coin.Amount) {
			return types.ErrInsufficientFunds
		}
		m.setBalance(from.String(), coin.Denom, fromBal.Sub(coin.Amount))
		toBal := m.getBalance(to.String(), coin.Denom)
		m.setBalance(to.String(), coin.Denom, toBal.Add(coin.Amount))
	}
	return nil
}

func (m *mockBankKeeper) SendCoinsFromAccountToModule(_ context.Context, sender sdk.AccAddress, module string, amt sdk.Coins) error {
	moduleAddr := authtypes.NewModuleAddress(module)
	for _, coin := range amt {
		fromBal := m.getBalance(sender.String(), coin.Denom)
		if fromBal.LT(coin.Amount) {
			return types.ErrInsufficientFunds
		}
		m.setBalance(sender.String(), coin.Denom, fromBal.Sub(coin.Amount))
		toBal := m.getBalance(moduleAddr.String(), coin.Denom)
		m.setBalance(moduleAddr.String(), coin.Denom, toBal.Add(coin.Amount))
	}
	return nil
}

func (m *mockBankKeeper) SendCoinsFromModuleToAccount(_ context.Context, module string, recipient sdk.AccAddress, amt sdk.Coins) error {
	moduleAddr := authtypes.NewModuleAddress(module)
	for _, coin := range amt {
		fromBal := m.getBalance(moduleAddr.String(), coin.Denom)
		if fromBal.LT(coin.Amount) {
			return types.ErrInsufficientFunds
		}
		m.setBalance(moduleAddr.String(), coin.Denom, fromBal.Sub(coin.Amount))
		toBal := m.getBalance(recipient.String(), coin.Denom)
		m.setBalance(recipient.String(), coin.Denom, toBal.Add(coin.Amount))
	}
	return nil
}

func (m *mockBankKeeper) MintCoins(_ context.Context, moduleName string, amt sdk.Coins) error {
	moduleAddr := authtypes.NewModuleAddress(moduleName)
	for _, coin := range amt {
		bal := m.getBalance(moduleAddr.String(), coin.Denom)
		m.setBalance(moduleAddr.String(), coin.Denom, bal.Add(coin.Amount))
	}
	return nil
}

// ============================================================
// Test Setup
// ============================================================

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context, *mockBankKeeper) {
	t.Helper()
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{Height: 1}, false, log.NewNopLogger())

	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

	accountKeeper := mockAccountKeeper{
		moduleAddrs: map[string]sdk.AccAddress{
			types.ModuleName: authtypes.NewModuleAddress(types.ModuleName),
		},
	}
	bankKeeper := newMockBankKeeper()

	k := keeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		accountKeeper,
		bankKeeper,
		authtypes.NewModuleAddress("gov").String(),
	)

	return k, ctx, bankKeeper
}

func testAddr(i int) string {
	return sdk.AccAddress([]byte{byte(i), 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}).String()
}

const quoteDenom = "uusdc"
const tokenDenom = "newtoken"

func fundUser(bk *mockBankKeeper, userIdx int, denom string, amount math.Int) string {
	addr := testAddr(userIdx)
	bk.setBalance(addr, denom, amount)
	return addr
}

func fundModule(bk *mockBankKeeper, denom string, amount math.Int) {
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName).String()
	bk.setBalance(moduleAddr, denom, amount)
}

func defaultCreateLaunchMsg(creator string) *types.MsgCreateLaunch {
	return &types.MsgCreateLaunch{
		Creator:       creator,
		TokenDenom:    tokenDenom,
		TokenSupply:   math.NewInt(1_000_000),
		PricePerToken: math.NewInt(1), // 1 quote token per token
		QuoteDenom:    quoteDenom,
		SoftCap:       math.NewInt(100_000),
		HardCap:       math.NewInt(500_000),
		MaxPerWallet:  math.NewInt(100_000),
		StartBlock:    5,
		EndBlock:      100,
		VestingBlocks: 0,
		TGEPercent:    100,
	}
}

// ============================================================
// Tests
// ============================================================

func TestCreateLaunch(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)

	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, uint64(1), launchID)

	launch, found := k.GetLaunch(ctx, launchID)
	require.True(t, found)
	require.Equal(t, types.LaunchStatusPending, launch.Status)
	require.Equal(t, tokenDenom, launch.TokenDenom)
	require.Equal(t, math.NewInt(1_000_000), launch.TokenSupply)
	require.Equal(t, math.ZeroInt(), launch.TotalRaised)
	require.Equal(t, uint64(0), launch.Contributors)
}

func TestCreateLaunchSequentialIDs(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	creator := testAddr(1)

	id1, err := k.ExecuteCreateLaunch(ctx, defaultCreateLaunchMsg(creator))
	require.NoError(t, err)
	require.Equal(t, uint64(1), id1)

	id2, err := k.ExecuteCreateLaunch(ctx, defaultCreateLaunchMsg(creator))
	require.NoError(t, err)
	require.Equal(t, uint64(2), id2)
}

func TestContributeBasic(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	// Advance to start block
	ctx = ctx.WithBlockHeight(10)

	user := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(50_000))
	require.NoError(t, err)

	launch, _ := k.GetLaunch(ctx, launchID)
	require.Equal(t, types.LaunchStatusActive, launch.Status)
	require.Equal(t, math.NewInt(50_000), launch.TotalRaised)
	require.Equal(t, uint64(1), launch.Contributors)

	contrib, found := k.GetContribution(ctx, launchID, user)
	require.True(t, found)
	require.Equal(t, math.NewInt(50_000), contrib.Amount)
}

func TestContributeMultipleUsers(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	launchID, err := k.ExecuteCreateLaunch(ctx, defaultCreateLaunchMsg(creator))
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(10)

	user1 := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))
	user2 := fundUser(bk, 3, quoteDenom, math.NewInt(200_000))

	err = k.ExecuteContribute(ctx, user1, launchID, math.NewInt(80_000))
	require.NoError(t, err)

	err = k.ExecuteContribute(ctx, user2, launchID, math.NewInt(60_000))
	require.NoError(t, err)

	launch, _ := k.GetLaunch(ctx, launchID)
	require.Equal(t, math.NewInt(140_000), launch.TotalRaised)
	require.Equal(t, uint64(2), launch.Contributors)
}

func TestContributeBeforeStartRejected(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	// Block 1, start is 5
	user := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(50_000))
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrLaunchNotActive)
}

func TestContributeAfterEndRejected(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	// First activate it
	ctx = ctx.WithBlockHeight(10)
	user := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(10_000))
	require.NoError(t, err)

	// Now go past end block
	ctx = ctx.WithBlockHeight(101)
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(10_000))
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrLaunchEnded)
}

func TestContributeExceedHardCap(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	msg.HardCap = math.NewInt(100_000)
	msg.MaxPerWallet = math.NewInt(200_000) // high max per wallet
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(10)

	user := fundUser(bk, 2, quoteDenom, math.NewInt(500_000))
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(80_000))
	require.NoError(t, err)

	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(30_000))
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrHardCapExceeded)
}

func TestContributeMaxPerWallet(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	msg.MaxPerWallet = math.NewInt(50_000)
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(10)

	user := fundUser(bk, 2, quoteDenom, math.NewInt(500_000))
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(40_000))
	require.NoError(t, err)

	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(20_000))
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrMaxPerWalletExceeded)
}

func TestFinalizeSuccess(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(10)

	// Contribute enough to meet soft cap
	user1 := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))
	user2 := fundUser(bk, 3, quoteDenom, math.NewInt(200_000))
	err = k.ExecuteContribute(ctx, user1, launchID, math.NewInt(60_000))
	require.NoError(t, err)
	err = k.ExecuteContribute(ctx, user2, launchID, math.NewInt(60_000))
	require.NoError(t, err)

	// Finalize
	ctx = ctx.WithBlockHeight(101)
	status, err := k.ExecuteFinalizeLaunch(ctx, authtypes.NewModuleAddress("gov").String(), launchID)
	require.NoError(t, err)
	require.Equal(t, types.LaunchStatusSuccessful, status)

	launch, _ := k.GetLaunch(ctx, launchID)
	require.Equal(t, types.LaunchStatusSuccessful, launch.Status)
	require.Equal(t, int64(101), launch.FinalizedBlock)
}

func TestFinalizeFailed(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(10)

	// Contribute less than soft cap
	user := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(50_000))
	require.NoError(t, err)

	// Finalize
	ctx = ctx.WithBlockHeight(101)
	status, err := k.ExecuteFinalizeLaunch(ctx, authtypes.NewModuleAddress("gov").String(), launchID)
	require.NoError(t, err)
	require.Equal(t, types.LaunchStatusFailed, status)
}

func TestClaimTokensNoVesting(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(10)

	user1 := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))
	user2 := fundUser(bk, 3, quoteDenom, math.NewInt(200_000))

	err = k.ExecuteContribute(ctx, user1, launchID, math.NewInt(60_000))
	require.NoError(t, err)
	err = k.ExecuteContribute(ctx, user2, launchID, math.NewInt(40_000))
	require.NoError(t, err)

	// Finalize (no tokenfactory/dex keepers set = tokens tracked via tokenDenom)
	ctx = ctx.WithBlockHeight(101)
	status, err := k.ExecuteFinalizeLaunch(ctx, authtypes.NewModuleAddress("gov").String(), launchID)
	require.NoError(t, err)
	require.Equal(t, types.LaunchStatusSuccessful, status)

	// Fund module with tokens for distribution
	fundModule(bk, tokenDenom, math.NewInt(1_000_000))

	// User1 contributed 60k out of 100k total, so gets 60% of 1M = 600k tokens
	amount1, err := k.ExecuteClaimTokens(ctx, user1, launchID)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(600_000), amount1)

	// User2 contributed 40k out of 100k total, so gets 40% of 1M = 400k tokens
	amount2, err := k.ExecuteClaimTokens(ctx, user2, launchID)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(400_000), amount2)
}

func TestClaimTokensAlreadyClaimed(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(10)
	user := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(100_000))
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(101)
	_, err = k.ExecuteFinalizeLaunch(ctx, authtypes.NewModuleAddress("gov").String(), launchID)
	require.NoError(t, err)

	fundModule(bk, tokenDenom, math.NewInt(1_000_000))

	_, err = k.ExecuteClaimTokens(ctx, user, launchID)
	require.NoError(t, err)

	// Second claim should fail
	_, err = k.ExecuteClaimTokens(ctx, user, launchID)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrAlreadyClaimed)
}

func TestClaimRefund(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(10)
	user := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(50_000))
	require.NoError(t, err)

	// Verify user balance decreased
	userBal := bk.getBalance(user, quoteDenom)
	require.Equal(t, math.NewInt(150_000), userBal)

	// Finalize as failed (below soft cap)
	ctx = ctx.WithBlockHeight(101)
	status, err := k.ExecuteFinalizeLaunch(ctx, authtypes.NewModuleAddress("gov").String(), launchID)
	require.NoError(t, err)
	require.Equal(t, types.LaunchStatusFailed, status)

	// Claim refund
	refundAmount, err := k.ExecuteClaimRefund(ctx, user, launchID)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(50_000), refundAmount)

	// Verify user balance restored
	userBal = bk.getBalance(user, quoteDenom)
	require.Equal(t, math.NewInt(200_000), userBal)
}

func TestClaimRefundAlreadyClaimed(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(10)
	user := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(50_000))
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(101)
	_, err = k.ExecuteFinalizeLaunch(ctx, authtypes.NewModuleAddress("gov").String(), launchID)
	require.NoError(t, err)

	_, err = k.ExecuteClaimRefund(ctx, user, launchID)
	require.NoError(t, err)

	// Second refund should fail
	_, err = k.ExecuteClaimRefund(ctx, user, launchID)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrAlreadyClaimed)
}

func TestClaimRefundOnSuccessfulLaunchFails(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(10)
	user := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(100_000))
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(101)
	_, err = k.ExecuteFinalizeLaunch(ctx, authtypes.NewModuleAddress("gov").String(), launchID)
	require.NoError(t, err)

	// Try to claim refund on successful launch
	_, err = k.ExecuteClaimRefund(ctx, user, launchID)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrLaunchNotFailed)
}

func TestClaimTokensOnFailedLaunchFails(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(10)
	user := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(50_000))
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(101)
	_, err = k.ExecuteFinalizeLaunch(ctx, authtypes.NewModuleAddress("gov").String(), launchID)
	require.NoError(t, err)

	// Try to claim tokens on failed launch
	_, err = k.ExecuteClaimTokens(ctx, user, launchID)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrLaunchNotSuccessful)
}

func TestVestingClaims(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	msg.VestingBlocks = 100
	msg.TGEPercent = 25 // 25% at TGE, 75% vests over 100 blocks
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(10)
	user := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(100_000))
	require.NoError(t, err)

	// Finalize at block 101
	ctx = ctx.WithBlockHeight(101)
	status, err := k.ExecuteFinalizeLaunch(ctx, authtypes.NewModuleAddress("gov").String(), launchID)
	require.NoError(t, err)
	require.Equal(t, types.LaunchStatusSuccessful, status)

	// Fund module with tokens
	fundModule(bk, tokenDenom, math.NewInt(1_000_000))

	// First claim should give TGE amount (25% of 1M = 250k)
	tgeAmount, err := k.ExecuteClaimTokens(ctx, user, launchID)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(250_000), tgeAmount)

	// After 50 blocks, should have 50% of vesting (75% * 50/100 = 37.5% = 375k additional)
	ctx = ctx.WithBlockHeight(151)
	vestedAmount, err := k.ExecuteClaimTokens(ctx, user, launchID)
	require.NoError(t, err)
	// Vesting: 750k over 100 blocks = 7500 per block. 50 blocks = 375000
	require.Equal(t, math.NewInt(375_000), vestedAmount)

	// After all 100 blocks, should get remaining
	ctx = ctx.WithBlockHeight(201)
	finalAmount, err := k.ExecuteClaimTokens(ctx, user, launchID)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(375_000), finalAmount)

	// Nothing more to claim
	_, err = k.ExecuteClaimTokens(ctx, user, launchID)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrNothingToClaim)
}

func TestBeginBlockAutoFinalize(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	msg.StartBlock = 5
	msg.EndBlock = 50
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	// Contribute
	ctx = ctx.WithBlockHeight(10)
	user := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(100_000))
	require.NoError(t, err)

	// Check it's active
	launch, _ := k.GetLaunch(ctx, launchID)
	require.Equal(t, types.LaunchStatusActive, launch.Status)

	// Trigger BeginBlock after end
	ctx = ctx.WithBlockHeight(51)
	k.CheckAndFinalizeExpiredLaunches(ctx)

	// Should be finalized as successful (100k >= 100k soft cap)
	launch, _ = k.GetLaunch(ctx, launchID)
	require.Equal(t, types.LaunchStatusSuccessful, launch.Status)
}

func TestBeginBlockAutoFinalizeFailed(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	msg.StartBlock = 5
	msg.EndBlock = 50
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	// Small contribution
	ctx = ctx.WithBlockHeight(10)
	user := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(10_000))
	require.NoError(t, err)

	// Trigger BeginBlock after end
	ctx = ctx.WithBlockHeight(51)
	k.CheckAndFinalizeExpiredLaunches(ctx)

	// Should be finalized as failed (10k < 100k soft cap)
	launch, _ := k.GetLaunch(ctx, launchID)
	require.Equal(t, types.LaunchStatusFailed, launch.Status)
}

func TestDoubleFinalizeFails(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(10)
	user := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(100_000))
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(101)
	_, err = k.ExecuteFinalizeLaunch(ctx, authtypes.NewModuleAddress("gov").String(), launchID)
	require.NoError(t, err)

	// Try to finalize again
	_, err = k.ExecuteFinalizeLaunch(ctx, authtypes.NewModuleAddress("gov").String(), launchID)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrLaunchAlreadyFinalized)
}

func TestContributeInsufficientFunds(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(10)
	// Fund user with only 10k but try to contribute 50k
	user := fundUser(bk, 2, quoteDenom, math.NewInt(10_000))
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(50_000))
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrInsufficientFunds)
}

func TestNoContributionClaim(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(10)
	user := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(100_000))
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(101)
	_, err = k.ExecuteFinalizeLaunch(ctx, authtypes.NewModuleAddress("gov").String(), launchID)
	require.NoError(t, err)

	// User who didn't contribute tries to claim
	nonContributor := testAddr(5)
	_, err = k.ExecuteClaimTokens(ctx, nonContributor, launchID)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrNoContribution)
}

func TestMsgServer(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)

	// Test via msg server interface
	resp, err := k.CreateLaunch(ctx, defaultCreateLaunchMsg(creator))
	require.NoError(t, err)
	require.Equal(t, uint64(1), resp.LaunchID)

	ctx = ctx.WithBlockHeight(10)
	user := fundUser(bk, 2, quoteDenom, math.NewInt(200_000))

	_, err = k.Contribute(ctx, &types.MsgContribute{Sender: user, LaunchID: 1, Amount: math.NewInt(100_000)})
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(101)
	finalizeResp, err := k.FinalizeLaunch(ctx, &types.MsgFinalizeLaunch{Authority: creator, LaunchID: 1})
	require.NoError(t, err)
	require.Equal(t, types.LaunchStatusSuccessful, finalizeResp.Status)

	fundModule(bk, tokenDenom, math.NewInt(1_000_000))

	claimResp, err := k.ClaimTokens(ctx, &types.MsgClaimTokens{Sender: user, LaunchID: 1})
	require.NoError(t, err)
	require.Equal(t, math.NewInt(1_000_000), claimResp.Amount)
}

func TestIncrementalContributions(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	creator := testAddr(1)
	msg := defaultCreateLaunchMsg(creator)
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(10)
	user := fundUser(bk, 2, quoteDenom, math.NewInt(500_000))

	// Three incremental contributions
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(30_000))
	require.NoError(t, err)
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(30_000))
	require.NoError(t, err)
	err = k.ExecuteContribute(ctx, user, launchID, math.NewInt(30_000))
	require.NoError(t, err)

	// Total should be 90k
	contrib, found := k.GetContribution(ctx, launchID, user)
	require.True(t, found)
	require.Equal(t, math.NewInt(90_000), contrib.Amount)

	// Only 1 contributor
	launch, _ := k.GetLaunch(ctx, launchID)
	require.Equal(t, uint64(1), launch.Contributors)
}

func TestLaunchNotFound(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	ctx = ctx.WithBlockHeight(10)
	err := k.ExecuteContribute(ctx, testAddr(1), 999, math.NewInt(1000))
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrLaunchNotFound)
}
