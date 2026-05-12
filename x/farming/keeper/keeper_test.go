package keeper_test

import (
	"testing"

	"context"

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

	"syreen/x/farming/keeper"
	"syreen/x/farming/types"

	dbm "github.com/cosmos/cosmos-db"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

// Mock keepers
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

func (m *mockBankKeeper) SpendableCoins(_ context.Context, addr sdk.AccAddress) sdk.Coins {
	coins := sdk.NewCoins()
	if bals, ok := m.balances[addr.String()]; ok {
		for denom, amt := range bals {
			coins = coins.Add(sdk.NewCoin(denom, amt))
		}
	}
	return coins
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

func (m *mockBankKeeper) SendCoinsFromModuleToModule(_ context.Context, senderModule, recipientModule string, amt sdk.Coins) error {
	fromAddr := authtypes.NewModuleAddress(senderModule)
	toAddr := authtypes.NewModuleAddress(recipientModule)
	for _, coin := range amt {
		fromBal := m.getBalance(fromAddr.String(), coin.Denom)
		if fromBal.LT(coin.Amount) {
			return types.ErrInsufficientFunds
		}
		m.setBalance(fromAddr.String(), coin.Denom, fromBal.Sub(coin.Amount))
		toBal := m.getBalance(toAddr.String(), coin.Denom)
		m.setBalance(toAddr.String(), coin.Denom, toBal.Add(coin.Amount))
	}
	return nil
}

// Test helpers
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

const lpDenom = "factory/syreen1abc/lp/1"
const rewardDenom = "usyreen"

func createTestFarm(t *testing.T, k keeper.Keeper, ctx sdk.Context, bk *mockBankKeeper) {
	authority := authtypes.NewModuleAddress("gov").String()
	err := k.CreateFarmPool(ctx, authority, 1, lpDenom, math.NewInt(10_000_000), 0, 0)
	require.NoError(t, err)
	// Fund farming module with rewards
	farmModuleAddr := authtypes.NewModuleAddress(types.ModuleName)
	bk.setBalance(farmModuleAddr.String(), rewardDenom, math.NewInt(1_000_000_000_000)) // 1M SYR
}

// ============================================================
// Tests
// ============================================================

func TestCreateFarm(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)

	farm, ok := k.GetFarm(ctx, 1)
	require.True(t, ok)
	require.Equal(t, uint64(1), farm.PoolID)
	require.Equal(t, lpDenom, farm.LPDenom)
	require.Equal(t, "usyreen", farm.RewardDenom)
	require.Equal(t, math.NewInt(10_000_000), farm.RewardPerBlock)
	require.True(t, farm.Active)
	require.True(t, farm.TotalStaked.IsZero())
}

func TestCreateFarmUnauthorized(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	err := k.CreateFarmPool(ctx, testAddr(1), 1, lpDenom, math.NewInt(10_000_000), 0, 0)
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestCreateFarmDuplicate(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)
	authority := authtypes.NewModuleAddress("gov").String()
	err := k.CreateFarmPool(ctx, authority, 1, lpDenom, math.NewInt(5_000_000), 0, 0)
	require.ErrorIs(t, err, types.ErrFarmAlreadyExists)
}

func TestStakeLP(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)

	user := testAddr(1)
	bk.setBalance(user, lpDenom, math.NewInt(1_000_000))

	_, err := k.StakeLP(ctx, user, 1, math.NewInt(500_000))
	require.NoError(t, err)

	// Check position
	pos, ok := k.GetPosition(ctx, 1, user)
	require.True(t, ok)
	require.Equal(t, math.NewInt(500_000), pos.Amount)

	// Check farm total
	farm, _ := k.GetFarm(ctx, 1)
	require.Equal(t, math.NewInt(500_000), farm.TotalStaked)

	// Check LP tokens moved to module
	require.Equal(t, math.NewInt(500_000), bk.getBalance(user, lpDenom))
}

func TestStakeLPNoFarm(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	_, err := k.StakeLP(ctx, testAddr(1), 99, math.NewInt(100))
	require.ErrorIs(t, err, types.ErrFarmNotFound)
}

func TestStakeLPInactiveFarm(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)
	authority := authtypes.NewModuleAddress("gov").String()
	err := k.UpdateFarmConfig(ctx, authority, 1, math.NewInt(10_000_000), false)
	require.NoError(t, err)

	bk.setBalance(testAddr(1), lpDenom, math.NewInt(1_000_000))
	_, err = k.StakeLP(ctx, testAddr(1), 1, math.NewInt(100))
	require.ErrorIs(t, err, types.ErrFarmNotActive)
}

func TestStakeLPInsufficientBalance(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)

	bk.setBalance(testAddr(1), lpDenom, math.NewInt(100))
	_, err := k.StakeLP(ctx, testAddr(1), 1, math.NewInt(1_000))
	require.Error(t, err) // insufficient funds
}

func TestUnstakeLP(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)

	user := testAddr(1)
	bk.setBalance(user, lpDenom, math.NewInt(1_000_000))

	_, err := k.StakeLP(ctx, user, 1, math.NewInt(500_000))
	require.NoError(t, err)

	// Advance blocks to accumulate rewards
	ctx = ctx.WithBlockHeight(11)
	claimed, err := k.UnstakeLP(ctx, user, 1, math.NewInt(200_000))
	require.NoError(t, err)

	// Should have received rewards
	require.True(t, claimed.IsPositive(), "should have claimed rewards, got: %s", claimed.String())

	// Check position reduced
	pos, ok := k.GetPosition(ctx, 1, user)
	require.True(t, ok)
	require.Equal(t, math.NewInt(300_000), pos.Amount)

	// Check LP returned
	require.Equal(t, math.NewInt(700_000), bk.getBalance(user, lpDenom))
}

func TestUnstakeAll(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)

	user := testAddr(1)
	bk.setBalance(user, lpDenom, math.NewInt(1_000_000))

	_, err := k.StakeLP(ctx, user, 1, math.NewInt(500_000))
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(11)
	_, err = k.UnstakeLP(ctx, user, 1, math.NewInt(500_000))
	require.NoError(t, err)

	// Position should be deleted
	_, ok := k.GetPosition(ctx, 1, user)
	require.False(t, ok)

	// All LP returned
	require.Equal(t, math.NewInt(1_000_000), bk.getBalance(user, lpDenom))
}

func TestUnstakeTooMuch(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)

	user := testAddr(1)
	bk.setBalance(user, lpDenom, math.NewInt(1_000_000))
	_, _ = k.StakeLP(ctx, user, 1, math.NewInt(500_000))

	_, err := k.UnstakeLP(ctx, user, 1, math.NewInt(600_000))
	require.ErrorIs(t, err, types.ErrInsufficientStake)
}

func TestUnstakeNoPosition(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)

	_, err := k.UnstakeLP(ctx, testAddr(1), 1, math.NewInt(100))
	require.ErrorIs(t, err, types.ErrNoPosition)
}

func TestClaimReward(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)

	user := testAddr(1)
	bk.setBalance(user, lpDenom, math.NewInt(1_000_000))

	_, err := k.StakeLP(ctx, user, 1, math.NewInt(1_000_000))
	require.NoError(t, err)

	// Advance 100 blocks
	ctx = ctx.WithBlockHeight(101)
	claimed, err := k.ClaimFarmReward(ctx, user, 1)
	require.NoError(t, err)

	// 100 blocks × 10_000_000 per block = 1_000_000_000 rewards
	require.Equal(t, math.NewInt(1_000_000_000), claimed)

	// Check user received SYR
	require.Equal(t, math.NewInt(1_000_000_000), bk.getBalance(user, rewardDenom))

	// Position should still exist
	pos, ok := k.GetPosition(ctx, 1, user)
	require.True(t, ok)
	require.Equal(t, math.NewInt(1_000_000), pos.Amount)
}

func TestClaimRewardNoPosition(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)

	_, err := k.ClaimFarmReward(ctx, testAddr(1), 1)
	require.ErrorIs(t, err, types.ErrNoPosition)
}

func TestClaimRewardNoPending(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)

	user := testAddr(1)
	bk.setBalance(user, lpDenom, math.NewInt(1_000_000))
	_, _ = k.StakeLP(ctx, user, 1, math.NewInt(1_000_000))

	// Same block, no rewards yet
	_, err := k.ClaimFarmReward(ctx, user, 1)
	require.ErrorIs(t, err, types.ErrNoPendingReward)
}

func TestRewardDistributionProportional(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)

	user1 := testAddr(1)
	user2 := testAddr(2)
	bk.setBalance(user1, lpDenom, math.NewInt(3_000_000))
	bk.setBalance(user2, lpDenom, math.NewInt(1_000_000))

	// User1 stakes 3x more than User2
	_, _ = k.StakeLP(ctx, user1, 1, math.NewInt(3_000_000))
	_, _ = k.StakeLP(ctx, user2, 1, math.NewInt(1_000_000))

	// Advance 100 blocks
	ctx = ctx.WithBlockHeight(101)

	// Total rewards: 100 × 10_000_000 = 1_000_000_000
	// User1 should get 75% = 750_000_000
	// User2 should get 25% = 250_000_000

	claimed1, err := k.ClaimFarmReward(ctx, user1, 1)
	require.NoError(t, err)

	claimed2, err := k.ClaimFarmReward(ctx, user2, 1)
	require.NoError(t, err)

	require.Equal(t, math.NewInt(750_000_000), claimed1)
	require.Equal(t, math.NewInt(250_000_000), claimed2)
}

func TestRewardAfterPartialUnstake(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)

	user := testAddr(1)
	bk.setBalance(user, lpDenom, math.NewInt(1_000_000))
	_, _ = k.StakeLP(ctx, user, 1, math.NewInt(1_000_000))

	// Advance 50 blocks
	ctx = ctx.WithBlockHeight(51)
	// Unstake half — should claim 50 blocks of rewards
	claimed, _ := k.UnstakeLP(ctx, user, 1, math.NewInt(500_000))
	// 50 blocks × 10_000_000 = 500_000_000
	require.Equal(t, math.NewInt(500_000_000), claimed)

	// Advance another 50 blocks
	ctx = ctx.WithBlockHeight(101)
	// Claim with half the stake — 50 blocks × 10_000_000 = 500_000_000 (only staker)
	claimed2, _ := k.ClaimFarmReward(ctx, user, 1)
	require.Equal(t, math.NewInt(500_000_000), claimed2)
}

func TestMultipleStakes(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)

	user := testAddr(1)
	bk.setBalance(user, lpDenom, math.NewInt(2_000_000))

	_, _ = k.StakeLP(ctx, user, 1, math.NewInt(500_000))

	// Advance 10 blocks then stake more
	ctx = ctx.WithBlockHeight(11)
	pending, err := k.StakeLP(ctx, user, 1, math.NewInt(500_000))
	require.NoError(t, err)
	// Should report pending rewards from first 10 blocks
	require.Equal(t, math.NewInt(100_000_000), pending)

	// Check position total
	pos, _ := k.GetPosition(ctx, 1, user)
	require.Equal(t, math.NewInt(1_000_000), pos.Amount)
}

func TestUpdateFarmConfig(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)

	authority := authtypes.NewModuleAddress("gov").String()

	// Double the reward rate
	err := k.UpdateFarmConfig(ctx, authority, 1, math.NewInt(20_000_000), true)
	require.NoError(t, err)

	farm, _ := k.GetFarm(ctx, 1)
	require.Equal(t, math.NewInt(20_000_000), farm.RewardPerBlock)
	require.True(t, farm.Active)

	// Deactivate
	err = k.UpdateFarmConfig(ctx, authority, 1, math.NewInt(20_000_000), false)
	require.NoError(t, err)
	farm, _ = k.GetFarm(ctx, 1)
	require.False(t, farm.Active)
}

func TestUpdateFarmUnauthorized(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)
	err := k.UpdateFarmConfig(ctx, testAddr(1), 1, math.NewInt(5_000_000), true)
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestGetAllFarms(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	authority := authtypes.NewModuleAddress("gov").String()

	_ = k.CreateFarmPool(ctx, authority, 1, "lp1", math.NewInt(10_000_000), 0, 0)
	_ = k.CreateFarmPool(ctx, authority, 2, "lp2", math.NewInt(5_000_000), 0, 0)

	farmModuleAddr := authtypes.NewModuleAddress(types.ModuleName)
	bk.setBalance(farmModuleAddr.String(), rewardDenom, math.NewInt(1_000_000_000_000))

	farms := k.GetAllFarms(ctx)
	require.Len(t, farms, 2)
}

func TestGetPositionsByAddress(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	authority := authtypes.NewModuleAddress("gov").String()

	_ = k.CreateFarmPool(ctx, authority, 1, "lp1", math.NewInt(10_000_000), 0, 0)
	_ = k.CreateFarmPool(ctx, authority, 2, "lp2", math.NewInt(5_000_000), 0, 0)

	farmModuleAddr := authtypes.NewModuleAddress(types.ModuleName)
	bk.setBalance(farmModuleAddr.String(), rewardDenom, math.NewInt(1_000_000_000_000))

	user := testAddr(1)
	bk.setBalance(user, "lp1", math.NewInt(1_000_000))
	bk.setBalance(user, "lp2", math.NewInt(500_000))

	_, _ = k.StakeLP(ctx, user, 1, math.NewInt(1_000_000))
	_, _ = k.StakeLP(ctx, user, 2, math.NewInt(500_000))

	positions := k.GetPositionsByAddress(ctx, user)
	require.Len(t, positions, 2)
}

func TestDistributeRewards(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)

	user := testAddr(1)
	bk.setBalance(user, lpDenom, math.NewInt(1_000_000))
	_, _ = k.StakeLP(ctx, user, 1, math.NewInt(1_000_000))

	// Simulate BeginBlock advancing
	ctx = ctx.WithBlockHeight(11)
	k.DistributeRewards(ctx)

	// Check farm updated
	farm, _ := k.GetFarm(ctx, 1)
	require.Equal(t, int64(11), farm.LastRewardBlock)
	require.True(t, farm.AccRewardPerShare.IsPositive())
}

func TestRewardCapsAtModuleBalance(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)

	// Set module balance to only 50 SYR
	farmModuleAddr := authtypes.NewModuleAddress(types.ModuleName)
	bk.setBalance(farmModuleAddr.String(), rewardDenom, math.NewInt(50_000_000))

	user := testAddr(1)
	bk.setBalance(user, lpDenom, math.NewInt(1_000_000))
	_, _ = k.StakeLP(ctx, user, 1, math.NewInt(1_000_000))

	// Advance 100 blocks — would be 1000 SYR but only 50 available
	ctx = ctx.WithBlockHeight(101)
	claimed, err := k.ClaimFarmReward(ctx, user, 1)
	require.NoError(t, err)
	// Should get at most 50 SYR (50_000_000 usyreen)
	require.True(t, claimed.LTE(math.NewInt(50_000_000)), "claimed %s should be <= 50M", claimed)
}

func TestFarmAPR(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	createTestFarm(t, k, ctx, bk)

	// No stakers — infinite APR
	farm, _ := k.GetFarm(ctx, 1)
	apr := k.GetFarmAPR(ctx, farm)
	require.Equal(t, "∞", apr)

	// Stake some LP
	user := testAddr(1)
	bk.setBalance(user, lpDenom, math.NewInt(1_000_000))
	_, _ = k.StakeLP(ctx, user, 1, math.NewInt(1_000_000))

	farm, _ = k.GetFarm(ctx, 1)
	apr = k.GetFarmAPR(ctx, farm)
	// 10_000_000 per block × 63_072_000 blocks/year / 1_000_000 staked × 100
	// = 630_720_000_000_000 / 1_000_000 * 100 = very high %
	require.NotEqual(t, "∞", apr)
}

func TestFarmEndBlock(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	authority := authtypes.NewModuleAddress("gov").String()

	// Create farm that ends at block 50
	err := k.CreateFarmPool(ctx, authority, 1, lpDenom, math.NewInt(10_000_000), 1, 50)
	require.NoError(t, err)

	// Try to stake after end
	ctx = ctx.WithBlockHeight(51)
	_, err = k.StakeLP(ctx, testAddr(1), 1, math.NewInt(100))
	require.ErrorIs(t, err, types.ErrFarmEnded)
}
