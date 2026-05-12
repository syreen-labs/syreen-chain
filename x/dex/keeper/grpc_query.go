package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/math"

	"syreen/x/dex/types"
)

type queryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface.
func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

var _ types.QueryServer = queryServer{}

func (q queryServer) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params := q.Keeper.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

func (q queryServer) Pool(ctx context.Context, req *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	if req.PoolID == 0 {
		return nil, types.ErrPoolNotFound
	}

	pool, found := q.Keeper.GetPool(ctx, req.PoolID)
	if !found {
		return nil, types.ErrPoolNotFound
	}

	return &types.QueryPoolResponse{Pool: pool}, nil
}

func (q queryServer) Pools(ctx context.Context, _ *types.QueryPoolsRequest) (*types.QueryPoolsResponse, error) {
	pools := q.Keeper.GetAllPools(ctx)
	if pools == nil {
		pools = []types.Pool{}
	}
	return &types.QueryPoolsResponse{Pools: pools}, nil
}

func (q queryServer) SpotPrice(ctx context.Context, req *types.QuerySpotPriceRequest) (*types.QuerySpotPriceResponse, error) {
	if req.PoolID == 0 {
		return nil, types.ErrPoolNotFound
	}

	pool, found := q.Keeper.GetPool(ctx, req.PoolID)
	if !found {
		return nil, types.ErrPoolNotFound
	}

	priceAB, err := q.Keeper.GetSpotPrice(ctx, req.PoolID, pool.DenomA, pool.DenomB)
	if err != nil {
		return nil, err
	}

	priceBA, err := q.Keeper.GetSpotPrice(ctx, req.PoolID, pool.DenomB, pool.DenomA)
	if err != nil {
		return nil, err
	}

	return &types.QuerySpotPriceResponse{
		PoolID:  req.PoolID,
		PriceAB: priceAB,
		PriceBA: priceBA,
		DenomA:  pool.DenomA,
		DenomB:  pool.DenomB,
	}, nil
}

func (q queryServer) Quote(ctx context.Context, req *types.QueryQuoteRequest) (*types.QueryQuoteResponse, error) {
	if req.PoolID == 0 {
		return nil, types.ErrPoolNotFound
	}

	tokenOut, priceImpact, err := q.Keeper.GetQuote(ctx, req.PoolID, req.TokenIn)
	if err != nil {
		return nil, err
	}

	return &types.QueryQuoteResponse{
		TokenOut:    tokenOut,
		PriceImpact: priceImpact,
	}, nil
}

func (q queryServer) OptimalRoute(ctx context.Context, req *types.QueryOptimalRouteRequest) (*types.QueryOptimalRouteResponse, error) {
	if req.InputDenom == "" || req.OutputDenom == "" {
		return nil, types.ErrInvalidDenom
	}

	result, err := q.Keeper.FindOptimalRoute(ctx, req.InputDenom, req.OutputDenom, req.Amount)
	if err != nil {
		return nil, err
	}

	// Convert internal types to response types
	resp := &types.QueryOptimalRouteResponse{
		InputDenom:   result.InputDenom,
		OutputDenom:  result.OutputDenom,
		InputAmount:  result.InputAmount,
		OutputAmount: result.OutputAmount,
		PriceImpact:  result.PriceImpact,
		IsSplit:      result.IsSplit,
		RouteType:    result.RouteType,
	}

	for _, sr := range result.Routes {
		var hops []types.RouteHopInfo
		for _, h := range sr.Route.Hops {
			hops = append(hops, types.RouteHopInfo{
				PoolID:   h.PoolID,
				DenomIn:  h.DenomIn,
				DenomOut: h.DenomOut,
			})
		}
		resp.Routes = append(resp.Routes, types.SplitRouteInfo{
			Route: types.RouteInfo{
				Hops:        hops,
				ExpectedOut: sr.Route.ExpectedOut,
				PriceImpact: sr.Route.PriceImpact,
				TotalFees:   sr.Route.TotalFees,
				Score:       sr.Route.Score,
			},
			InputAmount:  sr.InputAmount,
			OutputAmount: sr.OutputAmount,
		})
	}

	return resp, nil
}

func (q queryServer) Routes(ctx context.Context, req *types.QueryRoutesRequest) (*types.QueryRoutesResponse, error) {
	if req.InputDenom == "" || req.OutputDenom == "" {
		return nil, types.ErrInvalidDenom
	}

	// Use a default amount for scoring if not provided
	amount := req.Amount
	if amount.IsNil() || amount.IsZero() {
		amount = math.NewInt(1_000_000) // default 1M base units
	}

	routes, err := q.Keeper.GetAllRoutesWithScores(ctx, req.InputDenom, req.OutputDenom, amount)
	if err != nil {
		return nil, err
	}

	var routeInfos []types.RouteInfo
	for _, r := range routes {
		var hops []types.RouteHopInfo
		for _, h := range r.Hops {
			hops = append(hops, types.RouteHopInfo{
				PoolID:   h.PoolID,
				DenomIn:  h.DenomIn,
				DenomOut: h.DenomOut,
			})
		}
		routeInfos = append(routeInfos, types.RouteInfo{
			Hops:        hops,
			ExpectedOut: r.ExpectedOut,
			PriceImpact: r.PriceImpact,
			TotalFees:   r.TotalFees,
			Score:       r.Score,
		})
	}

	if routeInfos == nil {
		routeInfos = []types.RouteInfo{}
	}

	return &types.QueryRoutesResponse{Routes: routeInfos}, nil
}

func (q queryServer) Signals(ctx context.Context, req *types.QuerySignalsRequest) (*types.QuerySignalsResponse, error) {
	if req.PoolID == 0 {
		return nil, fmt.Errorf("pool_id is required")
	}

	signals := q.Keeper.ComputePoolSignals(ctx, req.PoolID)
	if signals == nil {
		return nil, fmt.Errorf("no signal data available for pool %d", req.PoolID)
	}

	bz, err := json.Marshal(signals)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signals: %w", err)
	}

	return &types.QuerySignalsResponse{Signals: bz}, nil
}

func (q queryServer) SignalHistory(ctx context.Context, req *types.QuerySignalHistoryRequest) (*types.QuerySignalHistoryResponse, error) {
	if req.PoolID == 0 {
		return nil, fmt.Errorf("pool_id is required")
	}

	history, found := q.Keeper.GetSignalHistory(ctx, req.PoolID)
	if !found {
		return nil, fmt.Errorf("no signal history available for pool %d", req.PoolID)
	}

	bz, err := json.Marshal(history)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signal history: %w", err)
	}

	return &types.QuerySignalHistoryResponse{History: bz}, nil
}

func (q queryServer) RiskScore(ctx context.Context, req *types.QueryRiskScoreRequest) (*types.QueryRiskScoreResponse, error) {
	if req.PoolID == 0 {
		return nil, fmt.Errorf("pool_id is required")
	}

	score, found := q.Keeper.GetPoolRiskScore(ctx, req.PoolID)
	if !found {
		return nil, fmt.Errorf("no risk score available for pool %d", req.PoolID)
	}

	bz, err := json.Marshal(score)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal risk score: %w", err)
	}

	return &types.QueryRiskScoreResponse{RiskScore: bz}, nil
}

func (q queryServer) TradeRisk(ctx context.Context, req *types.QueryTradeRiskRequest) (*types.QueryTradeRiskResponse, error) {
	if req.PoolID == 0 {
		return nil, fmt.Errorf("pool_id is required")
	}

	tradeRisk, err := q.Keeper.CalculateTradeRisk(ctx, req.PoolID, req.InputDenom, req.Amount)
	if err != nil {
		return nil, err
	}

	bz, err := json.Marshal(tradeRisk)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal trade risk: %w", err)
	}

	return &types.QueryTradeRiskResponse{TradeRisk: bz}, nil
}

func (q queryServer) RiskHistory(ctx context.Context, req *types.QueryRiskHistoryRequest) (*types.QueryRiskHistoryResponse, error) {
	if req.PoolID == 0 {
		return nil, fmt.Errorf("pool_id is required")
	}

	history := q.Keeper.GetRiskHistory(ctx, req.PoolID)

	bz, err := json.Marshal(history)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal risk history: %w", err)
	}

	return &types.QueryRiskHistoryResponse{RiskHistory: bz}, nil
}

func (q queryServer) Sentiment(ctx context.Context, req *types.QuerySentimentRequest) (*types.QuerySentimentResponse, error) {
	if req.PoolID == 0 {
		return nil, fmt.Errorf("pool_id is required")
	}

	state, found := q.Keeper.GetSentimentState(ctx, req.PoolID)
	if !found {
		return &types.QuerySentimentResponse{
			PoolID:         req.PoolID,
			FearGreedIndex: 50,
			FearGreedLabel: "Neutral",
			Mood:           string(MoodNeutral),
			Indicators:     SentimentIndicators{BuySellRatio: 50, TradeSizeTrend: 50, IntentSentiment: 50, LiquidityFlow: 50, WhaleActivity: 50, VolumeTrend: 50},
			BlocksAnalyzed: 0,
		}, nil
	}

	return &types.QuerySentimentResponse{
		PoolID:         state.PoolID,
		FearGreedIndex: state.FearGreedIndex,
		FearGreedLabel: state.FearGreedLabel,
		Mood:           string(state.Mood),
		Indicators:     state.Indicators,
		UpdatedAt:      state.UpdatedAtHeight,
		BlocksAnalyzed: state.BlocksAnalyzed,
	}, nil
}

func (q queryServer) SentimentHistory(ctx context.Context, req *types.QuerySentimentHistoryRequest) (*types.QuerySentimentHistoryResponse, error) {
	if req.PoolID == 0 {
		return nil, fmt.Errorf("pool_id is required")
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 200
	}

	history := q.Keeper.GetSentimentHistory(ctx, req.PoolID, limit)
	if history == nil {
		history = []BlockTradingData{}
	}

	return &types.QuerySentimentHistoryResponse{
		PoolID:  req.PoolID,
		History: history,
	}, nil
}

func (q queryServer) OracleComposite(ctx context.Context, req *types.QueryOracleCompositeRequest) (*types.QueryOracleCompositeResponse, error) {
	if req.PoolID == 0 {
		return nil, fmt.Errorf("pool_id is required")
	}

	composite := q.Keeper.GetOracleComposite(ctx, req.PoolID)
	bz, err := json.Marshal(composite)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal oracle composite: %w", err)
	}

	return &types.QueryOracleCompositeResponse{Oracle: bz}, nil
}

func (q queryServer) WhaleAlerts(ctx context.Context, req *types.QueryWhaleAlertsRequest) (*types.QueryWhaleAlertsResponse, error) {
	if req.PoolID == 0 {
		return nil, fmt.Errorf("pool_id is required")
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 50
	}

	alerts := q.Keeper.GetWhaleAlerts(ctx, req.PoolID, limit)
	if alerts == nil {
		alerts = []WhaleAlert{}
	}

	bz, err := json.Marshal(alerts)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal whale alerts: %w", err)
	}

	return &types.QueryWhaleAlertsResponse{Alerts: bz}, nil
}

func (q queryServer) SentimentAlerts(ctx context.Context, req *types.QuerySentimentAlertsRequest) (*types.QuerySentimentAlertsResponse, error) {
	if req.PoolID == 0 {
		return nil, fmt.Errorf("pool_id is required")
	}

	alerts := q.Keeper.GetSentimentAlerts(ctx, req.PoolID)

	return &types.QuerySentimentAlertsResponse{
		PoolID: req.PoolID,
		Alerts: alerts.Alerts,
	}, nil
}

// ---------------------------------------------------------------------------
// Order Book Queries
// ---------------------------------------------------------------------------

func (q queryServer) OrderBook(ctx context.Context, req *types.QueryOrderBookRequest) (*types.QueryOrderBookResponse, error) {
	if req.PoolID == 0 {
		return nil, fmt.Errorf("pool_id is required")
	}

	ob := q.Keeper.GetOrderBook(ctx, req.PoolID)
	if ob.Bids == nil {
		ob.Bids = []types.OrderBookLevel{}
	}
	if ob.Asks == nil {
		ob.Asks = []types.OrderBookLevel{}
	}

	return &types.QueryOrderBookResponse{
		PoolID: ob.PoolID,
		Bids:   ob.Bids,
		Asks:   ob.Asks,
	}, nil
}

func (q queryServer) OrderBookDepth(ctx context.Context, req *types.QueryOrderBookDepthRequest) (*types.QueryOrderBookDepthResponse, error) {
	if req.PoolID == 0 {
		return nil, fmt.Errorf("pool_id is required")
	}

	ob := q.Keeper.GetOrderBook(ctx, req.PoolID)
	if ob.Bids == nil {
		ob.Bids = []types.OrderBookLevel{}
	}
	if ob.Asks == nil {
		ob.Asks = []types.OrderBookLevel{}
	}

	return &types.QueryOrderBookDepthResponse{
		PoolID: ob.PoolID,
		Bids:   ob.Bids,
		Asks:   ob.Asks,
	}, nil
}

func (q queryServer) UserOrders(ctx context.Context, req *types.QueryUserOrdersRequest) (*types.QueryUserOrdersResponse, error) {
	if req.Address == "" {
		return nil, fmt.Errorf("address is required")
	}

	orders := q.Keeper.GetOrdersByAddress(ctx, req.Address)
	if orders == nil {
		orders = []types.Order{}
	}

	return &types.QueryUserOrdersResponse{Orders: orders}, nil
}

func (q queryServer) OrderDetail(ctx context.Context, req *types.QueryOrderDetailRequest) (*types.QueryOrderDetailResponse, error) {
	if req.OrderID == 0 {
		return nil, types.ErrOrderNotFound
	}

	order, found := q.Keeper.GetOrder(ctx, req.OrderID)
	if !found {
		return nil, types.ErrOrderNotFound
	}

	return &types.QueryOrderDetailResponse{Order: order}, nil
}

func (q queryServer) Trades(ctx context.Context, req *types.QueryTradesRequest) (*types.QueryTradesResponse, error) {
	if req.PoolID == 0 {
		return nil, fmt.Errorf("pool_id is required")
	}

	limit := 100
	if req.Limit > 0 && req.Limit < 1000 {
		limit = int(req.Limit)
	}

	trades := q.Keeper.GetTradeHistory(ctx, req.PoolID, limit)
	if trades == nil {
		trades = []types.Trade{}
	}

	return &types.QueryTradesResponse{Trades: trades}, nil
}

// ---------------------------------------------------------------------------
// Referral Queries
// ---------------------------------------------------------------------------

func (q queryServer) ReferralInfo(ctx context.Context, req *types.QueryReferralInfoRequest) (*types.QueryReferralInfoResponse, error) {
	if req.Address == "" {
		return nil, fmt.Errorf("address is required")
	}

	resp := &types.QueryReferralInfoResponse{
		Address: req.Address,
	}

	code, hasCode := q.Keeper.GetReferralCodeByAddress(ctx, req.Address)
	if hasCode {
		resp.ReferralCode = code
	}

	stats, found := q.Keeper.GetReferrerStats(ctx, req.Address)
	if found {
		resp.TotalReferrals = stats.TotalReferrals
		resp.TotalVolume = stats.TotalVolume.String()
		resp.TotalEarned = stats.TotalEarned.String()
		resp.Tier = stats.Tier
		resp.ReferrerPct = GetTierReferrerPctBps(stats.Tier)
		resp.DiscountPct = GetTierDiscountPctBps(stats.Tier)
	} else {
		resp.TotalVolume = "0"
		resp.TotalEarned = "0"
		resp.Tier = TierBronze
		resp.ReferrerPct = GetTierReferrerPctBps(TierBronze)
		resp.DiscountPct = GetTierDiscountPctBps(TierBronze)
	}

	return resp, nil
}

func (q queryServer) ReferralList(ctx context.Context, req *types.QueryReferralListRequest) (*types.QueryReferralListResponse, error) {
	if req.Address == "" {
		return nil, fmt.Errorf("address is required")
	}

	users := q.Keeper.GetAllReferredUsers(ctx, req.Address)
	var referrals []types.ReferredUserInfo
	for _, u := range users {
		referrals = append(referrals, types.ReferredUserInfo{
			Address:     u.Address,
			TotalVolume: u.TotalVolume.String(),
			TotalSaved:  u.TotalSaved.String(),
		})
	}
	if referrals == nil {
		referrals = []types.ReferredUserInfo{}
	}

	return &types.QueryReferralListResponse{
		Referrer:  req.Address,
		Referrals: referrals,
	}, nil
}

func (q queryServer) ReferralCodeLookup(ctx context.Context, req *types.QueryReferralCodeLookupRequest) (*types.QueryReferralCodeLookupResponse, error) {
	if req.Code == "" {
		return nil, fmt.Errorf("code is required")
	}

	referrerAddr, found := q.Keeper.GetReferrerByCode(ctx, req.Code)
	if !found {
		return nil, types.ErrReferralCodeNotFound
	}

	stats, _ := q.Keeper.GetReferrerStats(ctx, referrerAddr)

	return &types.QueryReferralCodeLookupResponse{
		Code:           req.Code,
		ReferrerAddr:   referrerAddr,
		TotalReferrals: stats.TotalReferrals,
		Tier:           stats.Tier,
	}, nil
}

func (q queryServer) ReferralEarnings(ctx context.Context, req *types.QueryReferralEarningsRequest) (*types.QueryReferralEarningsResponse, error) {
	if req.Address == "" {
		return nil, fmt.Errorf("address is required")
	}

	earnings := q.Keeper.GetReferralEarnings(ctx, req.Address)
	var entries []types.ReferralEarningEntry
	for _, e := range earnings.Entries {
		entries = append(entries, types.ReferralEarningEntry{
			Height: e.Height,
			Trader: e.Trader,
			Denom:  e.Denom,
			Amount: e.Amount.String(),
			PoolID: e.PoolID,
		})
	}
	if entries == nil {
		entries = []types.ReferralEarningEntry{}
	}

	stats, found := q.Keeper.GetReferrerStats(ctx, req.Address)
	totalEarned := "0"
	if found {
		totalEarned = stats.TotalEarned.String()
	}

	return &types.QueryReferralEarningsResponse{
		Address:     req.Address,
		TotalEarned: totalEarned,
		Entries:     entries,
	}, nil
}

// ---------------------------------------------------------------------------
// Copy Trading Queries
// ---------------------------------------------------------------------------

func (q queryServer) Leaderboard(ctx context.Context, req *types.QueryLeaderboardRequest) (*types.QueryLeaderboardResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	traders := q.Keeper.GetLeaderboard(ctx, limit)
	if traders == nil {
		traders = []TraderStats{}
	}

	bz, err := json.Marshal(traders)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal leaderboard: %w", err)
	}

	return &types.QueryLeaderboardResponse{Traders: bz}, nil
}

func (q queryServer) TraderStats(ctx context.Context, req *types.QueryTraderStatsRequest) (*types.QueryTraderStatsResponse, error) {
	if req.Address == "" {
		return nil, fmt.Errorf("address is required")
	}

	stats, found := q.Keeper.GetTraderStats(ctx, req.Address)
	if !found {
		stats = q.Keeper.GetOrCreateTraderStats(ctx, req.Address)
	}

	bz, err := json.Marshal(stats)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal trader stats: %w", err)
	}

	return &types.QueryTraderStatsResponse{Stats: bz}, nil
}

func (q queryServer) TraderFollowers(ctx context.Context, req *types.QueryTraderFollowersRequest) (*types.QueryTraderFollowersResponse, error) {
	if req.Address == "" {
		return nil, fmt.Errorf("address is required")
	}

	followers := q.Keeper.GetFollowers(ctx, req.Address)
	if followers == nil {
		followers = []string{}
	}

	bz, err := json.Marshal(followers)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal followers: %w", err)
	}

	return &types.QueryTraderFollowersResponse{Followers: bz}, nil
}

func (q queryServer) CopyFollowing(ctx context.Context, req *types.QueryCopyFollowingRequest) (*types.QueryCopyFollowingResponse, error) {
	if req.Address == "" {
		return nil, fmt.Errorf("address is required")
	}

	following := q.Keeper.GetFollowing(ctx, req.Address)
	if following == nil {
		following = []CopySettings{}
	}

	bz, err := json.Marshal(following)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal following: %w", err)
	}

	return &types.QueryCopyFollowingResponse{Following: bz}, nil
}

func (q queryServer) CopyHistory(ctx context.Context, req *types.QueryCopyHistoryRequest) (*types.QueryCopyHistoryResponse, error) {
	if req.Address == "" {
		return nil, fmt.Errorf("address is required")
	}

	logs := q.Keeper.GetCopyTradeLog(ctx, req.Address)
	if logs == nil {
		logs = []CopyTradeLog{}
	}

	bz, err := json.Marshal(logs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal copy history: %w", err)
	}

	return &types.QueryCopyHistoryResponse{History: bz}, nil
}

// ---------------------------------------------------------------------------
// Time-Weighted Governance Queries
// ---------------------------------------------------------------------------

func (q queryServer) TimeWeightedPower(ctx context.Context, req *types.QueryTimeWeightedPowerRequest) (*types.QueryTimeWeightedPowerResponse, error) {
	if req.Address == "" {
		return nil, fmt.Errorf("address is required")
	}

	info, err := q.Keeper.GetTimeWeightedPowerInfo(ctx, req.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get time-weighted power: %w", err)
	}

	bz, err := json.Marshal(info)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal power info: %w", err)
	}

	return &types.QueryTimeWeightedPowerResponse{PowerInfo: bz}, nil
}

// ---------------------------------------------------------------------------
// Trader Track Record Queries
// ---------------------------------------------------------------------------

func (q queryServer) TraderRecord(ctx context.Context, req *types.QueryTraderRecordRequest) (*types.QueryTraderRecordResponse, error) {
	if req.Address == "" {
		return nil, fmt.Errorf("address is required")
	}

	record := q.Keeper.GetTraderRecord(ctx, req.Address)

	bz, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal trader record: %w", err)
	}

	return &types.QueryTraderRecordResponse{Record: bz}, nil
}

func (q queryServer) TraderLeaderboard(ctx context.Context, req *types.QueryTraderLeaderboardRequest) (*types.QueryTraderLeaderboardResponse, error) {
	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "pnl"
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	traders := q.Keeper.GetTraderLeaderboard(ctx, sortBy, limit)
	if traders == nil {
		traders = []TraderRecord{}
	}

	bz, err := json.Marshal(traders)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal leaderboard: %w", err)
	}

	return &types.QueryTraderLeaderboardResponse{Traders: bz}, nil
}

// ---------------------------------------------------------------------------
// IBC Cross-Chain Order Queries
// ---------------------------------------------------------------------------

func (q queryServer) IBCOrders(ctx context.Context, req *types.QueryIBCOrdersRequest) (*types.QueryIBCOrdersResponse, error) {
	if req.SourceChain == "" {
		return nil, fmt.Errorf("source_chain is required")
	}

	orders := q.Keeper.QueryIBCOrders(ctx, req.SourceChain, req.Sender)
	if orders == nil {
		orders = []IBCOrderInfo{}
	}

	bz, err := json.Marshal(orders)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal IBC orders: %w", err)
	}

	return &types.QueryIBCOrdersResponse{Orders: bz}, nil
}

// ---------------------------------------------------------------------------
// Dynamic Pool Fee Query
// ---------------------------------------------------------------------------

func (q queryServer) PoolFee(ctx context.Context, req *types.QueryPoolFeeRequest) (*types.QueryPoolFeeResponse, error) {
	if req.PoolID == 0 {
		return nil, fmt.Errorf("pool_id is required")
	}

	pool, found := q.Keeper.GetPool(ctx, req.PoolID)
	if !found {
		return nil, types.ErrPoolNotFound
	}

	effectiveFee := q.Keeper.GetEffectiveFee(ctx, req.PoolID)
	volatility := q.Keeper.CalculateVolatility(ctx, req.PoolID)

	return &types.QueryPoolFeeResponse{
		PoolID:               pool.ID,
		EffectiveFee:         effectiveFee,
		SwapFee:              pool.SwapFee,
		VolatilityFeeEnabled: pool.VolatilityFeeEnabled,
		BaseFee:              pool.BaseFee,
		MaxFee:               pool.MaxFee,
		VolatilityMultiplier: pool.VolatilityMultiplier,
		CurrentVolatility:    volatility,
	}, nil
}
