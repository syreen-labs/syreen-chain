package keeper_test

import (
	"context"
	"fmt"
	"testing"

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
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"

	"syreen/x/dex/keeper"
	"syreen/x/dex/types"
)

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("syreen", "syreenpub")
	config.Seal()
}

// ---------------------------------------------------------------------------
// Mock keepers
// ---------------------------------------------------------------------------

type mockAccountKeeper struct{}

func (m *mockAccountKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	return sdk.AccAddress([]byte("dex_module_addr_pad_")) // 20 bytes
}

type mockBankKeeper struct {
	balances map[string]sdk.Coins // address string -> coins
}

func newMockBankKeeper() *mockBankKeeper {
	return &mockBankKeeper{
		balances: make(map[string]sdk.Coins),
	}
}

func (m *mockBankKeeper) fundAccount(addr string, coins sdk.Coins) {
	existing := m.balances[addr]
	m.balances[addr] = existing.Add(coins...)
}

func (m *mockBankKeeper) getBalance(addr string, denom string) math.Int {
	for _, c := range m.balances[addr] {
		if c.Denom == denom {
			return c.Amount
		}
	}
	return math.ZeroInt()
}

func (m *mockBankKeeper) SendCoins(_ context.Context, from, to sdk.AccAddress, amt sdk.Coins) error {
	fromStr := from.String()
	toStr := to.String()
	for _, coin := range amt {
		bal := m.getBalance(fromStr, coin.Denom)
		if bal.LT(coin.Amount) {
			return fmt.Errorf("insufficient funds: %s < %s for %s", bal, coin.Amount, coin.Denom)
		}
	}
	m.balances[fromStr] = m.balances[fromStr].Sub(amt...)
	m.balances[toStr] = m.balances[toStr].Add(amt...)
	return nil
}

func (m *mockBankKeeper) SendCoinsFromAccountToModule(_ context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	senderStr := senderAddr.String()
	for _, coin := range amt {
		bal := m.getBalance(senderStr, coin.Denom)
		if bal.LT(coin.Amount) {
			return fmt.Errorf("insufficient funds: %s < %s for %s", bal, coin.Amount, coin.Denom)
		}
	}
	m.balances[senderStr] = m.balances[senderStr].Sub(amt...)
	m.balances[recipientModule] = m.balances[recipientModule].Add(amt...)
	return nil
}

func (m *mockBankKeeper) SendCoinsFromModuleToAccount(_ context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	recipientStr := recipientAddr.String()
	for _, coin := range amt {
		bal := m.getBalance(senderModule, coin.Denom)
		if bal.LT(coin.Amount) {
			return fmt.Errorf("insufficient module funds: %s < %s for %s", bal, coin.Amount, coin.Denom)
		}
	}
	m.balances[senderModule] = m.balances[senderModule].Sub(amt...)
	m.balances[recipientStr] = m.balances[recipientStr].Add(amt...)
	return nil
}

func (m *mockBankKeeper) MintCoins(_ context.Context, moduleName string, amt sdk.Coins) error {
	m.balances[moduleName] = m.balances[moduleName].Add(amt...)
	return nil
}

func (m *mockBankKeeper) GetBalance(_ context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	bal := m.getBalance(addr.String(), denom)
	return sdk.NewCoin(denom, bal)
}

func (m *mockBankKeeper) BurnCoins(_ context.Context, moduleName string, amt sdk.Coins) error {
	for _, coin := range amt {
		bal := m.getBalance(moduleName, coin.Denom)
		if bal.LT(coin.Amount) {
			return fmt.Errorf("insufficient module funds to burn")
		}
	}
	m.balances[moduleName] = m.balances[moduleName].Sub(amt...)
	return nil
}

type mockTokenFactoryKeeper struct {
	denoms map[string]bool
	bk     *mockBankKeeper
}

func newMockTokenFactoryKeeper(bk *mockBankKeeper) *mockTokenFactoryKeeper {
	return &mockTokenFactoryKeeper{
		denoms: make(map[string]bool),
		bk:     bk,
	}
}

func (m *mockTokenFactoryKeeper) CreateDenom(_ context.Context, creator, subdenom string) (string, error) {
	denom := fmt.Sprintf("factory/%s/%s", creator, subdenom)
	if m.denoms[denom] {
		return "", fmt.Errorf("denom already exists: %s", denom)
	}
	m.denoms[denom] = true
	return denom, nil
}

func (m *mockTokenFactoryKeeper) Mint(_ context.Context, sender string, amount sdk.Coin, mintTo string) error {
	if !m.denoms[amount.Denom] {
		return fmt.Errorf("denom does not exist: %s", amount.Denom)
	}
	// Mint to module first, then to recipient
	target := sender
	if mintTo != "" {
		target = mintTo
	}
	m.bk.balances[target] = m.bk.balances[target].Add(amount)
	return nil
}

func (m *mockTokenFactoryKeeper) Burn(_ context.Context, sender string, amount sdk.Coin, _ string) error {
	if !m.denoms[amount.Denom] {
		return fmt.Errorf("denom does not exist: %s", amount.Denom)
	}
	// Burn from the module account (tokens should already be there)
	bal := m.bk.getBalance(types.ModuleName, amount.Denom)
	if bal.LT(amount.Amount) {
		return fmt.Errorf("insufficient funds to burn")
	}
	m.bk.balances[types.ModuleName] = m.bk.balances[types.ModuleName].Sub(amount)
	return nil
}

// ---------------------------------------------------------------------------
// Setup helper
// ---------------------------------------------------------------------------

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context, *mockBankKeeper, *mockTokenFactoryKeeper) {
	t.Helper()

	storeKey := storetypes.NewKVStoreKey("dex")
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{Height: 100}, false, log.NewNopLogger())
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	storeService := runtime.NewKVStoreService(storeKey)

	bk := newMockBankKeeper()
	ak := &mockAccountKeeper{}
	tfk := newMockTokenFactoryKeeper(bk)

	k := keeper.NewKeeper(cdc, storeService, ak, bk, tfk, "authority")
	return k, ctx, bk, tfk
}

func testAddr() string {
	return sdk.AccAddress([]byte("test_address_padded_")).String()
}

func testAddr2() string {
	return sdk.AccAddress([]byte("second_addr_padding_")).String()
}

func moduleAddr() string {
	return sdk.AccAddress([]byte("dex_module_addr_pad_")).String()
}

// ---------------------------------------------------------------------------
// TestCreatePool
// ---------------------------------------------------------------------------

func TestCreatePool(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	// Fund creator with tokens
	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 10_000_000),
		sdk.NewInt64Coin("uusdc", 10_000_000),
	))

	poolID, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(5_000_000), math.NewInt(5_000_000))
	require.NoError(t, err)
	require.Equal(t, uint64(1), poolID)

	// Verify pool stored
	pool, found := k.GetPool(ctx, poolID)
	require.True(t, found)
	require.Equal(t, "usyreen", pool.DenomA) // alphabetically sorted: usyreen < uusdc
	require.Equal(t, "uusdc", pool.DenomB)
	require.Equal(t, math.NewInt(5_000_000), pool.ReserveA)
	require.Equal(t, math.NewInt(5_000_000), pool.ReserveB)
	require.True(t, pool.TotalShares.IsPositive())
	require.Equal(t, creator, pool.Creator)

	// Verify LP tokens minted to creator
	lpDenom := fmt.Sprintf("factory/%s/pool-1", moduleAddr())
	lpBal := bk.getBalance(creator, lpDenom)
	require.True(t, lpBal.IsPositive())
	require.Equal(t, pool.TotalShares, lpBal)

	// Verify denom pair index
	poolByDenom, found := k.GetPoolByDenomPair(ctx, "usyreen", "uusdc")
	require.True(t, found)
	require.Equal(t, poolID, poolByDenom.ID)

	// Verify next pool ID incremented
	require.Equal(t, uint64(2), k.GetNextPoolID(ctx))
}

// ---------------------------------------------------------------------------
// TestCreatePoolDuplicateFails
// ---------------------------------------------------------------------------

func TestCreatePoolDuplicateFails(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 20_000_000),
		sdk.NewInt64Coin("uusdc", 20_000_000),
	))

	_, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(5_000_000), math.NewInt(5_000_000))
	require.NoError(t, err)

	// Try creating same pair again
	_, err = k.CreatePool(ctx, creator, "uusdc", "usyreen", math.NewInt(5_000_000), math.NewInt(5_000_000))
	require.ErrorIs(t, err, types.ErrPoolAlreadyExists)
}

// ---------------------------------------------------------------------------
// TestAddLiquidity
// ---------------------------------------------------------------------------

func TestAddLiquidity(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	provider := testAddr2()

	// Create pool
	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 10_000_000),
		sdk.NewInt64Coin("uusdc", 10_000_000),
	))
	poolID, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(10_000_000), math.NewInt(10_000_000))
	require.NoError(t, err)

	poolBefore, _ := k.GetPool(ctx, poolID)

	// Provider adds liquidity
	bk.fundAccount(provider, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 5_000_000),
		sdk.NewInt64Coin("uusdc", 5_000_000),
	))
	shares, err := k.AddLiquidity(ctx, provider, poolID, math.NewInt(5_000_000), math.NewInt(5_000_000), math.ZeroInt())
	require.NoError(t, err)
	require.True(t, shares.IsPositive())

	// Verify pool updated
	poolAfter, _ := k.GetPool(ctx, poolID)
	require.Equal(t, poolBefore.ReserveA.Add(math.NewInt(5_000_000)), poolAfter.ReserveA)
	require.Equal(t, poolBefore.ReserveB.Add(math.NewInt(5_000_000)), poolAfter.ReserveB)
	require.Equal(t, poolBefore.TotalShares.Add(shares), poolAfter.TotalShares)

	// LP tokens minted to provider
	lpDenom := fmt.Sprintf("factory/%s/pool-%d", moduleAddr(), poolID)
	lpBal := bk.getBalance(provider, lpDenom)
	require.Equal(t, shares, lpBal)
}

// ---------------------------------------------------------------------------
// TestRemoveLiquidity
// ---------------------------------------------------------------------------

func TestRemoveLiquidity(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 10_000_000),
		sdk.NewInt64Coin("uusdc", 10_000_000),
	))
	poolID, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(10_000_000), math.NewInt(10_000_000))
	require.NoError(t, err)

	pool, _ := k.GetPool(ctx, poolID)
	totalShares := pool.TotalShares

	// Remove half the liquidity
	halfShares := totalShares.QuoRaw(2)

	// Transfer LP tokens to module for burning (mock needs them there)
	lpDenom := fmt.Sprintf("factory/%s/pool-%d", moduleAddr(), poolID)
	bk.balances[types.ModuleName] = bk.balances[types.ModuleName].Add(sdk.NewCoin(lpDenom, halfShares))
	// Deduct from creator
	bk.balances[creator] = bk.balances[creator].Sub(sdk.NewCoin(lpDenom, halfShares))
	// Fund module with tokens to return
	bk.balances[types.ModuleName] = bk.balances[types.ModuleName].Add(
		sdk.NewCoin("usyreen", math.NewInt(5_000_000)),
		sdk.NewCoin("uusdc", math.NewInt(5_000_000)),
	)
	// Give LP tokens back to creator for the test
	bk.balances[creator] = bk.balances[creator].Add(sdk.NewCoin(lpDenom, halfShares))
	// Reset module LP balance for proper flow
	bk.balances[types.ModuleName] = bk.balances[types.ModuleName].Sub(sdk.NewCoin(lpDenom, halfShares))

	amountA, amountB, err := k.RemoveLiquidity(ctx, creator, poolID, halfShares, math.ZeroInt(), math.ZeroInt())
	require.NoError(t, err)
	require.True(t, amountA.IsPositive())
	require.True(t, amountB.IsPositive())

	// Verify pool updated
	poolAfter, _ := k.GetPool(ctx, poolID)
	require.Equal(t, pool.ReserveA.Sub(amountA), poolAfter.ReserveA)
	require.Equal(t, pool.ReserveB.Sub(amountB), poolAfter.ReserveB)
	require.Equal(t, pool.TotalShares.Sub(halfShares), poolAfter.TotalShares)
}

// ---------------------------------------------------------------------------
// TestSwap
// ---------------------------------------------------------------------------

func TestSwap(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	trader := testAddr2()

	// Create pool: 10M usyreen / 10M uusdc (1:1)
	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 10_000_000),
		sdk.NewInt64Coin("uusdc", 10_000_000),
	))
	poolID, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(10_000_000), math.NewInt(10_000_000))
	require.NoError(t, err)

	// Fund module with tokens for the output side
	bk.balances[types.ModuleName] = bk.balances[types.ModuleName].Add(
		sdk.NewCoin("usyreen", math.NewInt(10_000_000)),
		sdk.NewCoin("uusdc", math.NewInt(10_000_000)),
	)

	// Trader swaps 1M usyreen for uusdc
	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))

	tokenIn := sdk.NewInt64Coin("usyreen", 1_000_000)
	tokenOut, err := k.Swap(ctx, trader, poolID, tokenIn, math.ZeroInt())
	require.NoError(t, err)
	require.Equal(t, "uusdc", tokenOut.Denom)
	require.True(t, tokenOut.Amount.IsPositive())

	// Output should be less than 1M due to price impact and fee
	require.True(t, tokenOut.Amount.LT(math.NewInt(1_000_000)))

	// Verify pool reserves updated
	pool, _ := k.GetPool(ctx, poolID)
	require.Equal(t, math.NewInt(10_000_000).Add(math.NewInt(1_000_000)), pool.ReserveA) // usyreen increased
	require.Equal(t, math.NewInt(10_000_000).Sub(tokenOut.Amount), pool.ReserveB)        // uusdc decreased

	// Verify x*y=k holds (approximately, with fee it actually increases)
	kBefore := math.NewInt(10_000_000).Mul(math.NewInt(10_000_000))
	kAfter := pool.ReserveA.Mul(pool.ReserveB)
	require.True(t, kAfter.GTE(kBefore), "k should increase or stay same due to fees")
}

// ---------------------------------------------------------------------------
// TestSwapSlippageExceeded
// ---------------------------------------------------------------------------

func TestSwapSlippageExceeded(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	trader := testAddr2()

	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 10_000_000),
		sdk.NewInt64Coin("uusdc", 10_000_000),
	))
	poolID, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(10_000_000), math.NewInt(10_000_000))
	require.NoError(t, err)

	bk.balances[types.ModuleName] = bk.balances[types.ModuleName].Add(
		sdk.NewCoin("usyreen", math.NewInt(10_000_000)),
		sdk.NewCoin("uusdc", math.NewInt(10_000_000)),
	)

	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))

	// Set unrealistically high minimum output
	tokenIn := sdk.NewInt64Coin("usyreen", 1_000_000)
	_, err = k.Swap(ctx, trader, poolID, tokenIn, math.NewInt(999_999))
	require.ErrorIs(t, err, types.ErrSlippageExceeded)
}

// ---------------------------------------------------------------------------
// TestSwapInsufficientLiquidity
// ---------------------------------------------------------------------------

func TestSwapInsufficientLiquidity(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	trader := testAddr2()

	// Create pool with minimal liquidity
	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 1_000_000),
		sdk.NewInt64Coin("uusdc", 1_000_000),
	))
	poolID, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(1_000_000), math.NewInt(1_000_000))
	require.NoError(t, err)

	bk.balances[types.ModuleName] = bk.balances[types.ModuleName].Add(
		sdk.NewCoin("usyreen", math.NewInt(1_000_000)),
		sdk.NewCoin("uusdc", math.NewInt(1_000_000)),
	)

	// Try to swap a tiny amount (1 unit) - output will be ~0 after fees
	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1)))
	tokenIn := sdk.NewInt64Coin("usyreen", 1)
	_, err = k.Swap(ctx, trader, poolID, tokenIn, math.ZeroInt())
	require.ErrorIs(t, err, types.ErrInsufficientLiquidity)
}

// ---------------------------------------------------------------------------
// TestGetQuote
// ---------------------------------------------------------------------------

func TestGetQuote(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 10_000_000),
		sdk.NewInt64Coin("uusdc", 10_000_000),
	))
	poolID, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(10_000_000), math.NewInt(10_000_000))
	require.NoError(t, err)

	tokenIn := sdk.NewInt64Coin("usyreen", 1_000_000)
	tokenOut, priceImpact, err := k.GetQuote(ctx, poolID, tokenIn)
	require.NoError(t, err)
	require.Equal(t, "uusdc", tokenOut.Denom)
	require.True(t, tokenOut.Amount.IsPositive())
	require.True(t, tokenOut.Amount.LT(math.NewInt(1_000_000)))
	require.True(t, priceImpact.IsPositive())
}

// ---------------------------------------------------------------------------
// TestMultipleSwaps
// ---------------------------------------------------------------------------

func TestMultipleSwaps(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	trader := testAddr2()

	// Create pool
	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 100_000_000),
		sdk.NewInt64Coin("uusdc", 100_000_000),
	))
	poolID, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(100_000_000), math.NewInt(100_000_000))
	require.NoError(t, err)

	bk.balances[types.ModuleName] = bk.balances[types.ModuleName].Add(
		sdk.NewCoin("usyreen", math.NewInt(100_000_000)),
		sdk.NewCoin("uusdc", math.NewInt(100_000_000)),
	)

	kInitial := math.NewInt(100_000_000).Mul(math.NewInt(100_000_000))

	// Perform multiple swaps in both directions
	for i := 0; i < 5; i++ {
		// Swap usyreen -> uusdc
		bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))
		tokenOut, err := k.Swap(ctx, trader, poolID, sdk.NewInt64Coin("usyreen", 1_000_000), math.ZeroInt())
		require.NoError(t, err)
		require.True(t, tokenOut.Amount.IsPositive())

		// Swap uusdc -> usyreen
		bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("uusdc", 1_000_000)))
		bk.balances[types.ModuleName] = bk.balances[types.ModuleName].Add(sdk.NewCoin("usyreen", math.NewInt(1_000_000))) // ensure module has tokens
		tokenOut2, err := k.Swap(ctx, trader, poolID, sdk.NewInt64Coin("uusdc", 1_000_000), math.ZeroInt())
		require.NoError(t, err)
		require.True(t, tokenOut2.Amount.IsPositive())
	}

	// After swaps, k should be >= initial k (fees accumulate)
	pool, _ := k.GetPool(ctx, poolID)
	kFinal := pool.ReserveA.Mul(pool.ReserveB)
	require.True(t, kFinal.GTE(kInitial), "k should grow due to accumulated fees")
}

// ---------------------------------------------------------------------------
// TestParams
// ---------------------------------------------------------------------------

func TestSetGetParams(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	p := k.GetParams(ctx)
	require.Equal(t, types.DefaultParams(), p)

	custom := types.Params{
		DefaultSwapFee:      math.LegacyNewDecWithPrec(5, 3), // 0.5%
		ProtocolFeeShare:    math.LegacyNewDecWithPrec(5, 1),  // 50%
		MinInitialLiquidity: math.NewInt(2_000_000),
		PoolCreationFee:     sdk.NewCoins(),
	}
	require.NoError(t, k.SetParams(ctx, custom))
	got := k.GetParams(ctx)
	require.Equal(t, custom.DefaultSwapFee, got.DefaultSwapFee)
	require.Equal(t, custom.MinInitialLiquidity, got.MinInitialLiquidity)
}

// ---------------------------------------------------------------------------
// TestGenesisRoundtrip
// ---------------------------------------------------------------------------

func TestGenesisRoundtrip(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 10_000_000),
		sdk.NewInt64Coin("uusdc", 10_000_000),
	))
	_, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(10_000_000), math.NewInt(10_000_000))
	require.NoError(t, err)

	// Export
	exported := k.ExportGenesis(ctx)
	require.NotNil(t, exported)
	require.Len(t, exported.Pools, 1)
	require.Equal(t, uint64(2), exported.NextPoolID)

	// Import into fresh keeper
	k2, ctx2, _, _ := setupKeeper(t)
	k2.InitGenesis(ctx2, *exported)

	// Verify
	exported2 := k2.ExportGenesis(ctx2)
	require.Len(t, exported2.Pools, 1)
	require.Equal(t, exported.NextPoolID, exported2.NextPoolID)
	require.Equal(t, exported.Pools[0].ID, exported2.Pools[0].ID)
	require.Equal(t, exported.Pools[0].DenomA, exported2.Pools[0].DenomA)
	require.Equal(t, exported.Pools[0].ReserveA, exported2.Pools[0].ReserveA)

	// Verify denom pair index rebuilt
	pool, found := k2.GetPoolByDenomPair(ctx2, "usyreen", "uusdc")
	require.True(t, found)
	require.Equal(t, uint64(1), pool.ID)
}

func TestInitGenesis_DefaultState(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	defaultGS := types.DefaultGenesis()
	k.InitGenesis(ctx, *defaultGS)

	params := k.GetParams(ctx)
	require.Equal(t, types.DefaultParams(), params)

	pools := k.GetAllPools(ctx)
	require.Empty(t, pools)
	require.Equal(t, uint64(1), k.GetNextPoolID(ctx))
}

// ---------------------------------------------------------------------------
// TestDenomSorting
// ---------------------------------------------------------------------------

func TestCreatePoolDenomsSorted(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	// Pass denoms in reverse order
	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("uusdc", 5_000_000),
		sdk.NewInt64Coin("usyreen", 3_000_000),
	))
	poolID, err := k.CreatePool(ctx, creator, "uusdc", "usyreen", math.NewInt(5_000_000), math.NewInt(3_000_000))
	require.NoError(t, err)

	pool, found := k.GetPool(ctx, poolID)
	require.True(t, found)
	// DenomA should be alphabetically first: usyreen < uusdc
	require.Equal(t, "usyreen", pool.DenomA)
	require.Equal(t, "uusdc", pool.DenomB)
	// Amounts should be swapped to match sorted denoms
	require.Equal(t, math.NewInt(3_000_000), pool.ReserveA) // usyreen amount
	require.Equal(t, math.NewInt(5_000_000), pool.ReserveB) // uusdc amount
}
