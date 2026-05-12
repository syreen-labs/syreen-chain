package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/aiagent/types"
)

const (
	blocksPerDay = int64(17_280) // ~5s blocks: 86400/5
)

type Keeper struct {
	cdc           codec.Codec
	storeService  store.KVStoreService
	bankKeeper    types.BankKeeper
	accountKeeper types.AccountKeeper
	dexKeeper     types.DexKeeper
	authority     string
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	authority string,
) Keeper {
	return Keeper{
		cdc:           cdc,
		storeService:  storeService,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
		authority:     authority,
	}
}

func (k *Keeper) SetDexKeeper(dk types.DexKeeper) { k.dexKeeper = dk }

// calculateMinOut computes the minimum acceptable output for a swap based on
// the current spot price and the agent's MaxSlippage setting. This prevents
// sandwich attacks from extracting value from agent trades.
func (k Keeper) calculateMinOut(ctx context.Context, cfg types.AgentConfig, denomIn, denomOut string, amountIn math.Int) math.Int {
	maxSlippage := cfg.MaxSlippage
	if maxSlippage.IsNil() || !maxSlippage.IsPositive() {
		maxSlippage = math.LegacyNewDecWithPrec(1, 2) // default 1%
	}
	spotPrice, err := k.dexKeeper.GetSpotPrice(ctx, cfg.PoolID, denomIn, denomOut)
	if err != nil || !spotPrice.IsPositive() {
		// Fallback: 5% slippage from a nominal 1:1 price.
		return amountIn.MulRaw(95).QuoRaw(100)
	}
	expectedOut := spotPrice.MulInt(amountIn).TruncateInt()
	minOut := math.LegacyOneDec().Sub(maxSlippage).MulInt(expectedOut).TruncateInt()
	if !minOut.IsPositive() {
		minOut = math.OneInt()
	}
	return minOut
}

// ============================================================
// Agent address derivation
// ============================================================

// AgentAddress returns a deterministic address for an agent based on its ID.
func AgentAddress(id uint64) sdk.AccAddress {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	hash := sha256.Sum256(append([]byte("aiagent/"), bz...))
	return hash[:20]
}

// ============================================================
// Agent CRUD
// ============================================================

func (k Keeper) GetAgent(ctx context.Context, id uint64) (types.Agent, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.AgentKey(id))
	if err != nil || bz == nil {
		return types.Agent{}, false
	}
	var agent types.Agent
	if err := json.Unmarshal(bz, &agent); err != nil {
		return types.Agent{}, false
	}
	return agent, true
}

func (k Keeper) SetAgent(ctx context.Context, agent types.Agent) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(agent)
	_ = kvStore.Set(types.AgentKey(agent.ID), bz)
}

func (k Keeper) GetAllAgents(ctx context.Context) []types.Agent {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.AgentPrefix)
	iter, err := kvStore.Iterator(prefix, append(prefix, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil {
		return nil
	}
	defer iter.Close()
	var agents []types.Agent
	for ; iter.Valid(); iter.Next() {
		var agent types.Agent
		if err := json.Unmarshal(iter.Value(), &agent); err == nil {
			agents = append(agents, agent)
		}
	}
	return agents
}

func (k Keeper) GetNextAgentID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := kvStore.Get([]byte(types.NextAgentIDKey))
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetNextAgentID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	_ = kvStore.Set([]byte(types.NextAgentIDKey), bz)
}

// ============================================================
// Core Operations
// ============================================================

func (k Keeper) ExecuteCreateAgent(
	ctx context.Context,
	owner string,
	name string,
	strategyType types.StrategyType,
	config types.AgentConfig,
	initialFunds math.Int,
) (uint64, string, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	id := k.GetNextAgentID(ctx)
	agentAddr := AgentAddress(id)
	agentAddrStr := agentAddr.String()

	ownerAddr, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		return 0, "", err
	}

	// Register the agent address as a proper account if it doesn't exist yet.
	// This ensures the address is trackable via auth queries and bank sends work correctly.
	if !k.accountKeeper.HasAccount(ctx, agentAddr) {
		acct := k.accountKeeper.NewAccountWithAddress(ctx, agentAddr)
		k.accountKeeper.SetAccount(ctx, acct)
	}

	// fund the agent account
	funds := sdk.NewCoins(sdk.NewCoin(config.InputDenom, initialFunds))
	if err := k.bankKeeper.SendCoins(ctx, ownerAddr, agentAddr, funds); err != nil {
		return 0, "", types.ErrInsufficientFunds
	}

	agent := types.Agent{
		ID:           id,
		Owner:        owner,
		Name:         name,
		StrategyType: strategyType,
		Status:       types.AgentStatusActive,
		Config:       config,
		Stats:        types.AgentStats{TotalVolumeIn: math.ZeroInt(), TotalVolumeOut: math.ZeroInt()},
		CreatedAt:    sdkCtx.BlockHeight(),
		AgentAddress: agentAddrStr,
	}
	k.SetAgent(ctx, agent)
	k.SetNextAgentID(ctx, id+1)

	return id, agentAddrStr, nil
}

func (k Keeper) ExecuteFundAgent(ctx context.Context, owner string, agentID uint64, amount math.Int) error {
	agent, found := k.GetAgent(ctx, agentID)
	if !found {
		return types.ErrAgentNotFound
	}
	if agent.Owner != owner {
		return types.ErrUnauthorized
	}

	ownerAddr, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		return err
	}
	agentAddr := AgentAddress(agentID)

	funds := sdk.NewCoins(sdk.NewCoin(agent.Config.InputDenom, amount))
	return k.bankKeeper.SendCoins(ctx, ownerAddr, agentAddr, funds)
}

func (k Keeper) ExecuteWithdrawAgentFunds(ctx context.Context, owner string, agentID uint64, amount math.Int) (math.Int, error) {
	agent, found := k.GetAgent(ctx, agentID)
	if !found {
		return math.ZeroInt(), types.ErrAgentNotFound
	}
	if agent.Owner != owner {
		return math.ZeroInt(), types.ErrUnauthorized
	}

	ownerAddr, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		return math.ZeroInt(), err
	}
	agentAddr := AgentAddress(agentID)

	// cap to available balance
	balance := k.bankKeeper.GetBalance(ctx, agentAddr, agent.Config.InputDenom)
	withdraw := amount
	if withdraw.GT(balance.Amount) {
		withdraw = balance.Amount
	}
	if !withdraw.IsPositive() {
		return math.ZeroInt(), types.ErrInsufficientFunds
	}

	funds := sdk.NewCoins(sdk.NewCoin(agent.Config.InputDenom, withdraw))
	if err := k.bankKeeper.SendCoins(ctx, agentAddr, ownerAddr, funds); err != nil {
		return math.ZeroInt(), err
	}

	// mark as drained if balance now zero
	newBalance := k.bankKeeper.GetBalance(ctx, agentAddr, agent.Config.InputDenom)
	if !newBalance.Amount.IsPositive() {
		agent.Status = types.AgentStatusDrained
		k.SetAgent(ctx, agent)
	}

	return withdraw, nil
}

func (k Keeper) ExecutePauseAgent(ctx context.Context, owner string, agentID uint64) error {
	agent, found := k.GetAgent(ctx, agentID)
	if !found {
		return types.ErrAgentNotFound
	}
	if agent.Owner != owner {
		return types.ErrUnauthorized
	}
	if agent.Status != types.AgentStatusActive {
		return types.ErrAgentNotActive
	}
	agent.Status = types.AgentStatusPaused
	k.SetAgent(ctx, agent)
	return nil
}

func (k Keeper) ExecuteResumeAgent(ctx context.Context, owner string, agentID uint64) error {
	agent, found := k.GetAgent(ctx, agentID)
	if !found {
		return types.ErrAgentNotFound
	}
	if agent.Owner != owner {
		return types.ErrUnauthorized
	}
	if agent.Status != types.AgentStatusPaused {
		return types.ErrAgentNotPaused
	}
	agent.Status = types.AgentStatusActive
	k.SetAgent(ctx, agent)
	return nil
}

func (k Keeper) ExecuteUpdateStrategy(ctx context.Context, owner string, agentID uint64, config types.AgentConfig) error {
	agent, found := k.GetAgent(ctx, agentID)
	if !found {
		return types.ErrAgentNotFound
	}
	if agent.Owner != owner {
		return types.ErrUnauthorized
	}
	agent.Config = config
	k.SetAgent(ctx, agent)
	return nil
}

// ============================================================
// BeginBlock — execute all active agents
// ============================================================

func (k Keeper) BeginBlockHandler(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := sdkCtx.BlockHeight()

	agents := k.GetAllAgents(ctx)
	for i := range agents {
		agent := &agents[i]
		if agent.Status != types.AgentStatusActive {
			continue
		}

		// reset daily trade counter
		if height-agent.Stats.DayStartBlock >= blocksPerDay {
			agent.Stats.TradesThisDay = 0
			agent.Stats.DayStartBlock = height
		}

		// check interval
		if height-agent.Stats.LastTradeBlock < agent.Config.IntervalBlocks {
			continue
		}

		// check daily cap
		maxPerDay := agent.Config.MaxTradesPerDay
		if maxPerDay <= 0 {
			maxPerDay = 10
		}
		if agent.Stats.TradesThisDay >= maxPerDay {
			continue
		}

		// execute strategy
		executed := k.executeStrategy(ctx, agent, height)
		if executed {
			agent.Stats.LastTradeBlock = height
			agent.Stats.TradesThisDay++
			agent.Stats.TotalTrades++
			k.SetAgent(ctx, *agent)
		}
	}
}

func (k Keeper) executeStrategy(ctx context.Context, agent *types.Agent, height int64) bool {
	switch agent.StrategyType {
	case types.StrategyMomentum:
		return k.executeMomentum(ctx, agent)
	case types.StrategyMeanReversion:
		return k.executeMeanReversion(ctx, agent)
	case types.StrategyDCA:
		return k.executeDCA(ctx, agent)
	case types.StrategyArbitrage:
		return k.executeArbitrage(ctx, agent)
	case types.StrategyGrid:
		return k.executeGrid(ctx, agent, height)
	}
	return false
}

// ============================================================
// Strategy Implementations
// ============================================================

func (k Keeper) executeMomentum(ctx context.Context, agent *types.Agent) bool {
	signal, score, _ := k.dexKeeper.GetSignalForPool(ctx, agent.Config.PoolID)

	cfg := agent.Config
	minStrength := cfg.MinSignalStrength
	if minStrength == "" {
		minStrength = "BUY"
	}

	agentAddr := AgentAddress(agent.ID)
	balance := k.bankKeeper.GetBalance(ctx, agentAddr, cfg.InputDenom)
	if !balance.Amount.IsPositive() {
		return false
	}

	tradePercent := cfg.TradePercent
	if tradePercent.IsNil() || tradePercent.IsZero() {
		tradePercent = math.LegacyNewDecWithPrec(10, 2) // 10% default
	}

	tradeAmount := tradePercent.MulInt(balance.Amount).TruncateInt()
	if !tradeAmount.IsPositive() {
		return false
	}

	// BUY signal
	if (signal == "STRONG_BUY" || (minStrength == "BUY" && (signal == "BUY" || signal == "STRONG_BUY"))) && score >= 60 {
		minOut := k.calculateMinOut(ctx, cfg, cfg.InputDenom, cfg.OutputDenom, tradeAmount)
		_, err := k.dexKeeper.SmartSwap(ctx, agentAddr.String(), cfg.InputDenom, cfg.OutputDenom, tradeAmount, minOut)
		if err == nil {
			agent.Stats.TotalVolumeIn = agent.Stats.TotalVolumeIn.Add(tradeAmount)
			return true
		}
	}

	// SELL signal — swap output back to input
	if signal == "STRONG_SELL" || signal == "SELL" {
		outBalance := k.bankKeeper.GetBalance(ctx, agentAddr, cfg.OutputDenom)
		if outBalance.Amount.IsPositive() {
			sellAmount := tradePercent.MulInt(outBalance.Amount).TruncateInt()
			if sellAmount.IsPositive() {
				minOut := k.calculateMinOut(ctx, cfg, cfg.OutputDenom, cfg.InputDenom, sellAmount)
				_, err := k.dexKeeper.SmartSwap(ctx, agentAddr.String(), cfg.OutputDenom, cfg.InputDenom, sellAmount, minOut)
				if err == nil {
					agent.Stats.TotalVolumeOut = agent.Stats.TotalVolumeOut.Add(sellAmount)
					return true
				}
			}
		}
	}

	return false
}

func (k Keeper) executeMeanReversion(ctx context.Context, agent *types.Agent) bool {
	_, _, rsi := k.dexKeeper.GetSignalForPool(ctx, agent.Config.PoolID)

	cfg := agent.Config
	buyThreshold := cfg.RSIBuyThreshold
	sellThreshold := cfg.RSISellThreshold
	if buyThreshold == 0 {
		buyThreshold = 30
	}
	if sellThreshold == 0 {
		sellThreshold = 70
	}

	agentAddr := AgentAddress(agent.ID)
	tradePercent := cfg.TradePercent
	if tradePercent.IsNil() || tradePercent.IsZero() {
		tradePercent = math.LegacyNewDecWithPrec(10, 2)
	}

	// RSI oversold → BUY
	if rsi.LT(math.LegacyNewDec(buyThreshold)) {
		balance := k.bankKeeper.GetBalance(ctx, agentAddr, cfg.InputDenom)
		if balance.Amount.IsPositive() {
			tradeAmount := tradePercent.MulInt(balance.Amount).TruncateInt()
			if tradeAmount.IsPositive() {
				minOut := k.calculateMinOut(ctx, cfg, cfg.InputDenom, cfg.OutputDenom, tradeAmount)
				_, err := k.dexKeeper.SmartSwap(ctx, agentAddr.String(), cfg.InputDenom, cfg.OutputDenom, tradeAmount, minOut)
				if err == nil {
					agent.Stats.TotalVolumeIn = agent.Stats.TotalVolumeIn.Add(tradeAmount)
					return true
				}
			}
		}
	}

	// RSI overbought → SELL
	if rsi.GT(math.LegacyNewDec(sellThreshold)) {
		outBalance := k.bankKeeper.GetBalance(ctx, agentAddr, cfg.OutputDenom)
		if outBalance.Amount.IsPositive() {
			sellAmount := tradePercent.MulInt(outBalance.Amount).TruncateInt()
			if sellAmount.IsPositive() {
				minOut := k.calculateMinOut(ctx, cfg, cfg.OutputDenom, cfg.InputDenom, sellAmount)
				_, err := k.dexKeeper.SmartSwap(ctx, agentAddr.String(), cfg.OutputDenom, cfg.InputDenom, sellAmount, minOut)
				if err == nil {
					agent.Stats.TotalVolumeOut = agent.Stats.TotalVolumeOut.Add(sellAmount)
					return true
				}
			}
		}
	}

	return false
}

func (k Keeper) executeDCA(ctx context.Context, agent *types.Agent) bool {
	cfg := agent.Config
	agentAddr := AgentAddress(agent.ID)

	balance := k.bankKeeper.GetBalance(ctx, agentAddr, cfg.InputDenom)
	if !balance.Amount.IsPositive() {
		agent.Status = types.AgentStatusDrained
		return false
	}

	tradePercent := cfg.TradePercent
	if tradePercent.IsNil() || tradePercent.IsZero() {
		tradePercent = math.LegacyNewDecWithPrec(5, 2) // 5% per DCA
	}

	tradeAmount := tradePercent.MulInt(balance.Amount).TruncateInt()
	if !tradeAmount.IsPositive() {
		return false
	}

	minOut := k.calculateMinOut(ctx, cfg, cfg.InputDenom, cfg.OutputDenom, tradeAmount)
	_, err := k.dexKeeper.SmartSwap(ctx, agentAddr.String(), cfg.InputDenom, cfg.OutputDenom, tradeAmount, minOut)
	if err == nil {
		agent.Stats.TotalVolumeIn = agent.Stats.TotalVolumeIn.Add(tradeAmount)
		return true
	}
	return false
}

func (k Keeper) executeArbitrage(ctx context.Context, agent *types.Agent) bool {
	cfg := agent.Config
	if cfg.PoolID == 0 || cfg.PoolIDB == 0 {
		return false
	}

	denomA, denomB, ok := k.dexKeeper.GetPoolDenoms(ctx, cfg.PoolID)
	if !ok {
		return false
	}

	// get price in both pools
	priceA, err := k.dexKeeper.GetSpotPrice(ctx, cfg.PoolID, denomA, denomB)
	if err != nil {
		return false
	}
	priceB, err := k.dexKeeper.GetSpotPrice(ctx, cfg.PoolIDB, denomA, denomB)
	if err != nil {
		return false
	}

	minProfit := cfg.MinProfitPercent
	if minProfit.IsNil() || minProfit.IsZero() {
		minProfit = math.LegacyNewDecWithPrec(5, 3) // 0.5% default
	}

	agentAddr := AgentAddress(agent.ID)
	tradePercent := cfg.TradePercent
	if tradePercent.IsNil() || tradePercent.IsZero() {
		tradePercent = math.LegacyNewDecWithPrec(20, 2) // 20% for arb
	}

	// pool A cheaper → buy from A, sell into B
	var priceDiff math.LegacyDec
	var buyCheapPool, sellExpensivePool uint64
	var buyDenom, sellDenom string

	if priceA.LT(priceB) {
		priceDiff = priceB.Sub(priceA).Quo(priceA)
		buyCheapPool = cfg.PoolID
		sellExpensivePool = cfg.PoolIDB
		buyDenom = denomA
		sellDenom = denomB
	} else {
		priceDiff = priceA.Sub(priceB).Quo(priceB)
		buyCheapPool = cfg.PoolIDB
		sellExpensivePool = cfg.PoolID
		buyDenom = denomB
		sellDenom = denomA
	}
	_ = buyCheapPool
	_ = sellExpensivePool

	if priceDiff.LT(minProfit) {
		return false
	}

	// execute: swap inputDenom → buyDenom (cheap pool), then buyDenom → sellDenom (expensive pool)
	inputBalance := k.bankKeeper.GetBalance(ctx, agentAddr, cfg.InputDenom)
	if !inputBalance.Amount.IsPositive() {
		return false
	}
	tradeAmount := tradePercent.MulInt(inputBalance.Amount).TruncateInt()
	if !tradeAmount.IsPositive() {
		return false
	}

	// Wrap both swaps in CacheContext for atomicity — if step 2 fails, step 1 is rolled back
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cacheCtx, write := sdkCtx.CacheContext()

	// step 1: buy the asset from the cheaper pool
	minOut1 := k.calculateMinOut(cacheCtx, cfg, cfg.InputDenom, buyDenom, tradeAmount)
	got, err := k.dexKeeper.SmartSwap(cacheCtx, agentAddr.String(), cfg.InputDenom, buyDenom, tradeAmount, minOut1)
	if err != nil {
		return false
	}
	_ = sellDenom

	// step 2: sell it into the more expensive pool
	_, err = k.dexKeeper.SmartSwap(cacheCtx, agentAddr.String(), got.Denom, cfg.InputDenom, got.Amount, tradeAmount) // minOut = tradeAmount for profit
	if err != nil {
		// CacheContext not committed — both swaps are rolled back automatically
		return false
	}

	// Both swaps succeeded — commit the cached state
	write()

	agent.Stats.TotalVolumeIn = agent.Stats.TotalVolumeIn.Add(tradeAmount)
	return true
}

func (k Keeper) executeGrid(ctx context.Context, agent *types.Agent, height int64) bool {
	cfg := agent.Config
	if cfg.GridLower.IsNil() || cfg.GridUpper.IsNil() || cfg.GridLevels <= 0 {
		return false
	}

	agentAddr := AgentAddress(agent.ID)

	// get current price
	spotPrice, err := k.dexKeeper.GetSpotPrice(ctx, cfg.PoolID, cfg.InputDenom, cfg.OutputDenom)
	if err != nil || spotPrice.IsZero() {
		return false
	}

	// calculate grid step
	gridRange := cfg.GridUpper.Sub(cfg.GridLower)
	gridStep := gridRange.Quo(math.LegacyNewDec(cfg.GridLevels))

	// find which grid level current price is in
	if spotPrice.LT(cfg.GridLower) || spotPrice.GT(cfg.GridUpper) {
		return false // price out of range
	}

	// find the nearest grid lower bound
	levelFloat := spotPrice.Sub(cfg.GridLower).Quo(gridStep)
	level := levelFloat.TruncateInt64()
	gridBottom := cfg.GridLower.Add(gridStep.MulInt64(level))
	gridTop := gridBottom.Add(gridStep)

	midPoint := gridBottom.Add(gridTop).Quo(math.LegacyNewDec(2))

	tradePercent := cfg.TradePercent
	if tradePercent.IsNil() || tradePercent.IsZero() {
		tradePercent = math.LegacyNewDecWithPrec(5, 2)
	}

	// price in lower half of grid → BUY
	if spotPrice.LT(midPoint) {
		balance := k.bankKeeper.GetBalance(ctx, agentAddr, cfg.InputDenom)
		if balance.Amount.IsPositive() {
			tradeAmount := tradePercent.MulInt(balance.Amount).TruncateInt()
			if tradeAmount.IsPositive() {
				minOut := k.calculateMinOut(ctx, cfg, cfg.InputDenom, cfg.OutputDenom, tradeAmount)
				_, err := k.dexKeeper.SmartSwap(ctx, agentAddr.String(), cfg.InputDenom, cfg.OutputDenom, tradeAmount, minOut)
				if err == nil {
					agent.Stats.TotalVolumeIn = agent.Stats.TotalVolumeIn.Add(tradeAmount)
					return true
				}
			}
		}
	} else {
		// price in upper half → SELL
		outBalance := k.bankKeeper.GetBalance(ctx, agentAddr, cfg.OutputDenom)
		if outBalance.Amount.IsPositive() {
			sellAmount := tradePercent.MulInt(outBalance.Amount).TruncateInt()
			if sellAmount.IsPositive() {
				minOut := k.calculateMinOut(ctx, cfg, cfg.OutputDenom, cfg.InputDenom, sellAmount)
				_, err := k.dexKeeper.SmartSwap(ctx, agentAddr.String(), cfg.OutputDenom, cfg.InputDenom, sellAmount, minOut)
				if err == nil {
					agent.Stats.TotalVolumeOut = agent.Stats.TotalVolumeOut.Add(sellAmount)
					return true
				}
			}
		}
	}

	_ = fmt.Sprintf("grid level %d", level) // suppress unused warning
	_ = height
	return false
}
