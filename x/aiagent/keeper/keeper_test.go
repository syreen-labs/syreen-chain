package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"syreen/x/aiagent/keeper"
	"syreen/x/aiagent/types"
)

// ---- mock bank ----

type mockBank struct{ balances map[string]sdk.Coins }

func newMockBank() *mockBank { return &mockBank{balances: map[string]sdk.Coins{}} }

func (m *mockBank) credit(addr sdk.AccAddress, coins sdk.Coins) {
	m.balances[addr.String()] = m.balances[addr.String()].Add(coins...)
}

func (m *mockBank) GetBalance(_ context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	for _, c := range m.balances[addr.String()] {
		if c.Denom == denom {
			return c
		}
	}
	return sdk.NewCoin(denom, math.ZeroInt())
}

func (m *mockBank) SendCoins(_ context.Context, from, to sdk.AccAddress, amt sdk.Coins) error {
	m.balances[from.String()] = m.balances[from.String()].Sub(amt...)
	m.balances[to.String()] = m.balances[to.String()].Add(amt...)
	return nil
}

func (m *mockBank) SendCoinsFromAccountToModule(ctx context.Context, sender sdk.AccAddress, module string, amt sdk.Coins) error {
	return m.SendCoins(ctx, sender, sdk.AccAddress([]byte(module)), amt)
}

func (m *mockBank) SendCoinsFromModuleToAccount(ctx context.Context, module string, recipient sdk.AccAddress, amt sdk.Coins) error {
	return m.SendCoins(ctx, sdk.AccAddress([]byte(module)), recipient, amt)
}

// ---- mock dex ----

type mockDex struct {
	spotPrice math.LegacyDec
	signal    string
	score     int64
	rsi       math.LegacyDec
	denomA    string
	denomB    string
}

func (m *mockDex) GetSpotPrice(_ context.Context, _ uint64, _, _ string) (math.LegacyDec, error) {
	return m.spotPrice, nil
}
func (m *mockDex) SmartSwap(_ context.Context, _ string, inputDenom, outputDenom string, inputAmount, _ math.Int) (sdk.Coin, error) {
	return sdk.NewCoin(outputDenom, inputAmount), nil
}
func (m *mockDex) GetSignalForPool(_ context.Context, _ uint64) (string, int64, math.LegacyDec) {
	return m.signal, m.score, m.rsi
}
func (m *mockDex) GetPoolDenoms(_ context.Context, _ uint64) (string, string, bool) {
	return m.denomA, m.denomB, true
}

// ---- setup ----

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context, *mockBank, *mockDex) {
	t.Helper()
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{Height: 1000}, false, log.NewNopLogger())

	ir := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(ir)

	bank := newMockBank()
	dex := &mockDex{
		spotPrice: math.LegacyNewDec(100),
		signal:    "NEUTRAL",
		score:     50,
		rsi:       math.LegacyNewDec(50),
		denomA:    "usyreen",
		denomB:    "uusdc",
	}

	k := keeper.NewKeeper(cdc, runtime.NewKVStoreService(storeKey), bank, "authority")
	k.SetDexKeeper(dex)
	return k, ctx, bank, dex
}

func defaultConfig() types.AgentConfig {
	return types.AgentConfig{
		PoolID:          1,
		InputDenom:      "uusdc",
		OutputDenom:     "usyreen",
		TradePercent:    math.LegacyNewDecWithPrec(10, 2),
		MaxSlippage:     math.LegacyNewDecWithPrec(1, 2),
		IntervalBlocks:  10,
		MaxTradesPerDay: 5,
	}
}

var ownerAddr = sdk.AccAddress([]byte{1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}).String()

// ---- tests ----

func TestCreateAgent(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	owner, _ := sdk.AccAddressFromBech32(ownerAddr)
	bank.credit(owner, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1_000_000))))

	id, addr, err := k.ExecuteCreateAgent(ctx, ownerAddr, "TestBot", types.StrategyMomentum, defaultConfig(), math.NewInt(500_000))
	require.NoError(t, err)
	require.Equal(t, uint64(1), id)
	require.NotEmpty(t, addr)

	agent, found := k.GetAgent(ctx, id)
	require.True(t, found)
	require.Equal(t, types.AgentStatusActive, agent.Status)
	require.Equal(t, "TestBot", agent.Name)
}

func TestFundAgent(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	owner, _ := sdk.AccAddressFromBech32(ownerAddr)
	bank.credit(owner, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(2_000_000))))

	id, _, err := k.ExecuteCreateAgent(ctx, ownerAddr, "Bot", types.StrategyDCA, defaultConfig(), math.NewInt(500_000))
	require.NoError(t, err)

	require.NoError(t, k.ExecuteFundAgent(ctx, ownerAddr, id, math.NewInt(500_000)))

	agentAddr := keeper.AgentAddress(id)
	bal := bank.GetBalance(ctx, agentAddr, "uusdc")
	require.Equal(t, math.NewInt(1_000_000), bal.Amount)
}

func TestWithdrawAgentFunds(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	owner, _ := sdk.AccAddressFromBech32(ownerAddr)
	bank.credit(owner, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1_000_000))))

	id, _, err := k.ExecuteCreateAgent(ctx, ownerAddr, "Bot", types.StrategyDCA, defaultConfig(), math.NewInt(1_000_000))
	require.NoError(t, err)

	withdrawn, err := k.ExecuteWithdrawAgentFunds(ctx, ownerAddr, id, math.NewInt(400_000))
	require.NoError(t, err)
	require.Equal(t, math.NewInt(400_000), withdrawn)

	bal := bank.GetBalance(ctx, owner, "uusdc")
	require.Equal(t, math.NewInt(400_000), bal.Amount)
}

func TestPauseAndResume(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	owner, _ := sdk.AccAddressFromBech32(ownerAddr)
	bank.credit(owner, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1_000_000))))

	id, _, _ := k.ExecuteCreateAgent(ctx, ownerAddr, "Bot", types.StrategyMomentum, defaultConfig(), math.NewInt(500_000))

	require.NoError(t, k.ExecutePauseAgent(ctx, ownerAddr, id))
	agent, _ := k.GetAgent(ctx, id)
	require.Equal(t, types.AgentStatusPaused, agent.Status)

	require.NoError(t, k.ExecuteResumeAgent(ctx, ownerAddr, id))
	agent, _ = k.GetAgent(ctx, id)
	require.Equal(t, types.AgentStatusActive, agent.Status)
}

func TestPauseAlreadyPaused(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	owner, _ := sdk.AccAddressFromBech32(ownerAddr)
	bank.credit(owner, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1_000_000))))

	id, _, _ := k.ExecuteCreateAgent(ctx, ownerAddr, "Bot", types.StrategyMomentum, defaultConfig(), math.NewInt(500_000))
	require.NoError(t, k.ExecutePauseAgent(ctx, ownerAddr, id))
	require.Error(t, k.ExecutePauseAgent(ctx, ownerAddr, id))
}

func TestUnauthorized(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	owner, _ := sdk.AccAddressFromBech32(ownerAddr)
	other := sdk.AccAddress([]byte{2, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}).String()
	bank.credit(owner, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1_000_000))))

	id, _, _ := k.ExecuteCreateAgent(ctx, ownerAddr, "Bot", types.StrategyMomentum, defaultConfig(), math.NewInt(500_000))

	require.Error(t, k.ExecutePauseAgent(ctx, other, id))
	require.Error(t, k.ExecuteFundAgent(ctx, other, id, math.NewInt(1)))
	_, err := k.ExecuteWithdrawAgentFunds(ctx, other, id, math.NewInt(1))
	require.Error(t, err)
}

func TestUpdateStrategy(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	owner, _ := sdk.AccAddressFromBech32(ownerAddr)
	bank.credit(owner, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1_000_000))))

	id, _, _ := k.ExecuteCreateAgent(ctx, ownerAddr, "Bot", types.StrategyMomentum, defaultConfig(), math.NewInt(500_000))

	newCfg := defaultConfig()
	newCfg.IntervalBlocks = 20
	require.NoError(t, k.ExecuteUpdateStrategy(ctx, ownerAddr, id, newCfg))

	agent, _ := k.GetAgent(ctx, id)
	require.Equal(t, int64(20), agent.Config.IntervalBlocks)
}

func TestBeginBlock_MomentumBuy(t *testing.T) {
	k, ctx, bank, dex := setupKeeper(t)
	owner, _ := sdk.AccAddressFromBech32(ownerAddr)
	bank.credit(owner, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1_000_000))))

	dex.signal = "STRONG_BUY"
	dex.score = 80

	cfg := defaultConfig()
	cfg.IntervalBlocks = 1
	id, _, _ := k.ExecuteCreateAgent(ctx, ownerAddr, "Bot", types.StrategyMomentum, cfg, math.NewInt(1_000_000))

	k.BeginBlockHandler(ctx)

	agent, _ := k.GetAgent(ctx, id)
	require.Equal(t, int64(1), agent.Stats.TotalTrades)
}

func TestBeginBlock_NeutralNoTrade(t *testing.T) {
	k, ctx, bank, dex := setupKeeper(t)
	owner, _ := sdk.AccAddressFromBech32(ownerAddr)
	bank.credit(owner, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1_000_000))))

	dex.signal = "NEUTRAL"
	dex.score = 50

	cfg := defaultConfig()
	cfg.IntervalBlocks = 1
	id, _, _ := k.ExecuteCreateAgent(ctx, ownerAddr, "Bot", types.StrategyMomentum, cfg, math.NewInt(1_000_000))

	k.BeginBlockHandler(ctx)

	agent, _ := k.GetAgent(ctx, id)
	require.Equal(t, int64(0), agent.Stats.TotalTrades)
}

func TestBeginBlock_DCA(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	owner, _ := sdk.AccAddressFromBech32(ownerAddr)
	bank.credit(owner, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1_000_000))))

	cfg := defaultConfig()
	cfg.IntervalBlocks = 1
	id, _, _ := k.ExecuteCreateAgent(ctx, ownerAddr, "Bot", types.StrategyDCA, cfg, math.NewInt(1_000_000))

	k.BeginBlockHandler(ctx)

	agent, _ := k.GetAgent(ctx, id)
	require.Equal(t, int64(1), agent.Stats.TotalTrades)
}

func TestBeginBlock_PausedSkipped(t *testing.T) {
	k, ctx, bank, dex := setupKeeper(t)
	owner, _ := sdk.AccAddressFromBech32(ownerAddr)
	bank.credit(owner, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1_000_000))))

	dex.signal = "STRONG_BUY"
	dex.score = 90

	cfg := defaultConfig()
	cfg.IntervalBlocks = 1
	id, _, _ := k.ExecuteCreateAgent(ctx, ownerAddr, "Bot", types.StrategyMomentum, cfg, math.NewInt(1_000_000))
	require.NoError(t, k.ExecutePauseAgent(ctx, ownerAddr, id))

	k.BeginBlockHandler(ctx)

	agent, _ := k.GetAgent(ctx, id)
	require.Equal(t, int64(0), agent.Stats.TotalTrades)
}

func TestBeginBlock_IntervalRespected(t *testing.T) {
	k, ctx, bank, dex := setupKeeper(t)
	owner, _ := sdk.AccAddressFromBech32(ownerAddr)
	bank.credit(owner, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(10_000_000))))

	dex.signal = "STRONG_BUY"
	dex.score = 90

	cfg := defaultConfig()
	cfg.IntervalBlocks = 100
	id, _, _ := k.ExecuteCreateAgent(ctx, ownerAddr, "Bot", types.StrategyMomentum, cfg, math.NewInt(5_000_000))

	k.BeginBlockHandler(ctx) // block 1000 — trade
	agent, _ := k.GetAgent(ctx, id)
	require.Equal(t, int64(1), agent.Stats.TotalTrades)

	ctx = sdk.NewContext(ctx.MultiStore(), cmtproto.Header{Height: 1001}, false, log.NewNopLogger())
	k.BeginBlockHandler(ctx) // too soon
	agent, _ = k.GetAgent(ctx, id)
	require.Equal(t, int64(1), agent.Stats.TotalTrades)

	ctx = sdk.NewContext(ctx.MultiStore(), cmtproto.Header{Height: 1101}, false, log.NewNopLogger())
	k.BeginBlockHandler(ctx) // past interval
	agent, _ = k.GetAgent(ctx, id)
	require.Equal(t, int64(2), agent.Stats.TotalTrades)
}

func TestBeginBlock_MeanReversion_OversoldBuys(t *testing.T) {
	k, ctx, bank, dex := setupKeeper(t)
	owner, _ := sdk.AccAddressFromBech32(ownerAddr)
	bank.credit(owner, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1_000_000))))

	dex.rsi = math.LegacyNewDec(25) // oversold

	cfg := defaultConfig()
	cfg.IntervalBlocks = 1
	cfg.RSIBuyThreshold = 30
	cfg.RSISellThreshold = 70
	id, _, _ := k.ExecuteCreateAgent(ctx, ownerAddr, "Bot", types.StrategyMeanReversion, cfg, math.NewInt(1_000_000))

	k.BeginBlockHandler(ctx)

	agent, _ := k.GetAgent(ctx, id)
	require.Equal(t, int64(1), agent.Stats.TotalTrades)
}

func TestGetAllAgents(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	owner, _ := sdk.AccAddressFromBech32(ownerAddr)
	bank.credit(owner, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(10_000_000))))

	for i := 0; i < 4; i++ {
		_, _, err := k.ExecuteCreateAgent(ctx, ownerAddr, "Bot", types.StrategyDCA, defaultConfig(), math.NewInt(100_000))
		require.NoError(t, err)
	}
	require.Len(t, k.GetAllAgents(ctx), 4)
}

func TestAgentAddressDeterministic(t *testing.T) {
	a1 := keeper.AgentAddress(1)
	a2 := keeper.AgentAddress(1)
	a3 := keeper.AgentAddress(2)
	require.Equal(t, a1, a2)
	require.NotEqual(t, a1, a3)
}

func TestWithdrawAll_MarksDrained(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	owner, _ := sdk.AccAddressFromBech32(ownerAddr)
	bank.credit(owner, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1_000_000))))

	id, _, _ := k.ExecuteCreateAgent(ctx, ownerAddr, "Bot", types.StrategyDCA, defaultConfig(), math.NewInt(1_000_000))

	_, err := k.ExecuteWithdrawAgentFunds(ctx, ownerAddr, id, math.NewInt(1_000_000))
	require.NoError(t, err)

	agent, _ := k.GetAgent(ctx, id)
	require.Equal(t, types.AgentStatusDrained, agent.Status)
}
