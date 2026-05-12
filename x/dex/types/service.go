package types

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// Query request/response types
// ---------------------------------------------------------------------------

type QueryParamsRequest struct{}

func (m *QueryParamsRequest) ProtoMessage()           {}
func (m *QueryParamsRequest) Reset()                  { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string          { return "query_params_request" }
func (m *QueryParamsRequest) XXX_MessageName() string { return "syreen.dex.QueryParamsRequest" }

type QueryParamsResponse struct {
	Params Params `json:"params"`
}

func (m *QueryParamsResponse) ProtoMessage()           {}
func (m *QueryParamsResponse) Reset()                  { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string          { return fmt.Sprintf("params: %+v", m.Params) }
func (m *QueryParamsResponse) XXX_MessageName() string { return "syreen.dex.QueryParamsResponse" }

type QueryPoolRequest struct {
	PoolID uint64 `json:"pool_id"`
}

func (m *QueryPoolRequest) ProtoMessage()           {}
func (m *QueryPoolRequest) Reset()                  { *m = QueryPoolRequest{} }
func (m *QueryPoolRequest) String() string          { return fmt.Sprintf("query_pool: %d", m.PoolID) }
func (m *QueryPoolRequest) XXX_MessageName() string { return "syreen.dex.QueryPoolRequest" }

type QueryPoolResponse struct {
	Pool Pool `json:"pool"`
}

func (m *QueryPoolResponse) ProtoMessage()           {}
func (m *QueryPoolResponse) Reset()                  { *m = QueryPoolResponse{} }
func (m *QueryPoolResponse) String() string          { return fmt.Sprintf("pool: %+v", m.Pool) }
func (m *QueryPoolResponse) XXX_MessageName() string { return "syreen.dex.QueryPoolResponse" }

type QueryPoolsRequest struct{}

func (m *QueryPoolsRequest) ProtoMessage()           {}
func (m *QueryPoolsRequest) Reset()                  { *m = QueryPoolsRequest{} }
func (m *QueryPoolsRequest) String() string          { return "query_pools_request" }
func (m *QueryPoolsRequest) XXX_MessageName() string { return "syreen.dex.QueryPoolsRequest" }

type QueryPoolsResponse struct {
	Pools []Pool `json:"pools"`
}

func (m *QueryPoolsResponse) ProtoMessage()           {}
func (m *QueryPoolsResponse) Reset()                  { *m = QueryPoolsResponse{} }
func (m *QueryPoolsResponse) String() string          { return fmt.Sprintf("pools: %d", len(m.Pools)) }
func (m *QueryPoolsResponse) XXX_MessageName() string { return "syreen.dex.QueryPoolsResponse" }

type QueryQuoteRequest struct {
	PoolID  uint64   `json:"pool_id"`
	TokenIn sdk.Coin `json:"token_in"`
}

func (m *QueryQuoteRequest) ProtoMessage()           {}
func (m *QueryQuoteRequest) Reset()                  { *m = QueryQuoteRequest{} }
func (m *QueryQuoteRequest) String() string          { return fmt.Sprintf("query_quote: pool %d %s", m.PoolID, m.TokenIn) }
func (m *QueryQuoteRequest) XXX_MessageName() string { return "syreen.dex.QueryQuoteRequest" }

type QueryQuoteResponse struct {
	TokenOut    sdk.Coin       `json:"token_out"`
	PriceImpact math.LegacyDec `json:"price_impact"`
}

func (m *QueryQuoteResponse) ProtoMessage()           {}
func (m *QueryQuoteResponse) Reset()                  { *m = QueryQuoteResponse{} }
func (m *QueryQuoteResponse) String() string          { return fmt.Sprintf("quote: %s impact: %s", m.TokenOut, m.PriceImpact) }
func (m *QueryQuoteResponse) XXX_MessageName() string { return "syreen.dex.QueryQuoteResponse" }

type QuerySpotPriceRequest struct {
	PoolID uint64 `json:"pool_id"`
}

func (m *QuerySpotPriceRequest) ProtoMessage()           {}
func (m *QuerySpotPriceRequest) Reset()                  { *m = QuerySpotPriceRequest{} }
func (m *QuerySpotPriceRequest) String() string          { return fmt.Sprintf("query_spot_price: %d", m.PoolID) }
func (m *QuerySpotPriceRequest) XXX_MessageName() string { return "syreen.dex.QuerySpotPriceRequest" }

type QuerySpotPriceResponse struct {
	PoolID  uint64         `json:"pool_id"`
	PriceAB math.LegacyDec `json:"price_a_b"`
	PriceBA math.LegacyDec `json:"price_b_a"`
	DenomA  string         `json:"denom_a"`
	DenomB  string         `json:"denom_b"`
}

func (m *QuerySpotPriceResponse) ProtoMessage()           {}
func (m *QuerySpotPriceResponse) Reset()                  { *m = QuerySpotPriceResponse{} }
func (m *QuerySpotPriceResponse) String() string          { return fmt.Sprintf("spot_price: pool %d", m.PoolID) }
func (m *QuerySpotPriceResponse) XXX_MessageName() string { return "syreen.dex.QuerySpotPriceResponse" }

// ---------------------------------------------------------------------------
// Query request/response types for Smart Order Router
// ---------------------------------------------------------------------------

type QueryOptimalRouteRequest struct {
	InputDenom  string   `json:"input_denom"`
	OutputDenom string   `json:"output_denom"`
	Amount      math.Int `json:"amount"`
}

func (m *QueryOptimalRouteRequest) ProtoMessage()           {}
func (m *QueryOptimalRouteRequest) Reset()                  { *m = QueryOptimalRouteRequest{} }
func (m *QueryOptimalRouteRequest) String() string          { return fmt.Sprintf("query_optimal_route: %s->%s %s", m.InputDenom, m.OutputDenom, m.Amount) }
func (m *QueryOptimalRouteRequest) XXX_MessageName() string { return "syreen.dex.QueryOptimalRouteRequest" }

type RouteHopInfo struct {
	PoolID   uint64 `json:"pool_id"`
	DenomIn  string `json:"denom_in"`
	DenomOut string `json:"denom_out"`
}

type RouteInfo struct {
	Hops        []RouteHopInfo `json:"hops"`
	ExpectedOut math.Int       `json:"expected_out"`
	PriceImpact math.LegacyDec `json:"price_impact"`
	TotalFees   math.LegacyDec `json:"total_fees"`
	Score       math.LegacyDec `json:"score"`
}

type SplitRouteInfo struct {
	Route       RouteInfo `json:"route"`
	InputAmount math.Int  `json:"input_amount"`
	OutputAmount math.Int `json:"output_amount"`
}

type QueryOptimalRouteResponse struct {
	InputDenom   string           `json:"input_denom"`
	OutputDenom  string           `json:"output_denom"`
	InputAmount  math.Int         `json:"input_amount"`
	OutputAmount math.Int         `json:"output_amount"`
	PriceImpact  math.LegacyDec  `json:"price_impact"`
	Routes       []SplitRouteInfo `json:"routes"`
	IsSplit      bool             `json:"is_split"`
	RouteType    string           `json:"route_type"`
}

func (m *QueryOptimalRouteResponse) ProtoMessage()           {}
func (m *QueryOptimalRouteResponse) Reset()                  { *m = QueryOptimalRouteResponse{} }
func (m *QueryOptimalRouteResponse) String() string          { return fmt.Sprintf("optimal_route: %s->%s out=%s", m.InputDenom, m.OutputDenom, m.OutputAmount) }
func (m *QueryOptimalRouteResponse) XXX_MessageName() string { return "syreen.dex.QueryOptimalRouteResponse" }

type QueryRoutesRequest struct {
	InputDenom  string   `json:"input_denom"`
	OutputDenom string   `json:"output_denom"`
	Amount      math.Int `json:"amount"`
}

func (m *QueryRoutesRequest) ProtoMessage()           {}
func (m *QueryRoutesRequest) Reset()                  { *m = QueryRoutesRequest{} }
func (m *QueryRoutesRequest) String() string          { return fmt.Sprintf("query_routes: %s->%s", m.InputDenom, m.OutputDenom) }
func (m *QueryRoutesRequest) XXX_MessageName() string { return "syreen.dex.QueryRoutesRequest" }

type QueryRoutesResponse struct {
	Routes []RouteInfo `json:"routes"`
}

func (m *QueryRoutesResponse) ProtoMessage()           {}
func (m *QueryRoutesResponse) Reset()                  { *m = QueryRoutesResponse{} }
func (m *QueryRoutesResponse) String() string          { return fmt.Sprintf("routes: %d", len(m.Routes)) }
func (m *QueryRoutesResponse) XXX_MessageName() string { return "syreen.dex.QueryRoutesResponse" }

// ---------------------------------------------------------------------------
// Signals query request/response types
// ---------------------------------------------------------------------------

type QuerySignalsRequest struct {
	PoolID uint64 `json:"pool_id"`
}

func (m *QuerySignalsRequest) ProtoMessage()           {}
func (m *QuerySignalsRequest) Reset()                  { *m = QuerySignalsRequest{} }
func (m *QuerySignalsRequest) String() string          { return fmt.Sprintf("query_signals: %d", m.PoolID) }
func (m *QuerySignalsRequest) XXX_MessageName() string { return "syreen.dex.QuerySignalsRequest" }

type QuerySignalsResponse struct {
	Signals json.RawMessage `json:"signals"`
}

func (m *QuerySignalsResponse) ProtoMessage()           {}
func (m *QuerySignalsResponse) Reset()                  { *m = QuerySignalsResponse{} }
func (m *QuerySignalsResponse) String() string          { return "query_signals_response" }
func (m *QuerySignalsResponse) XXX_MessageName() string { return "syreen.dex.QuerySignalsResponse" }

type QuerySignalHistoryRequest struct {
	PoolID uint64 `json:"pool_id"`
}

func (m *QuerySignalHistoryRequest) ProtoMessage()           {}
func (m *QuerySignalHistoryRequest) Reset()                  { *m = QuerySignalHistoryRequest{} }
func (m *QuerySignalHistoryRequest) String() string          { return fmt.Sprintf("query_signal_history: %d", m.PoolID) }
func (m *QuerySignalHistoryRequest) XXX_MessageName() string { return "syreen.dex.QuerySignalHistoryRequest" }

type QuerySignalHistoryResponse struct {
	History json.RawMessage `json:"history"`
}

func (m *QuerySignalHistoryResponse) ProtoMessage()           {}
func (m *QuerySignalHistoryResponse) Reset()                  { *m = QuerySignalHistoryResponse{} }
func (m *QuerySignalHistoryResponse) String() string          { return "query_signal_history_response" }
func (m *QuerySignalHistoryResponse) XXX_MessageName() string { return "syreen.dex.QuerySignalHistoryResponse" }

// ---------------------------------------------------------------------------
// Risk Score query types
// ---------------------------------------------------------------------------

type QueryRiskScoreRequest struct {
	PoolID uint64 `json:"pool_id"`
}

func (m *QueryRiskScoreRequest) ProtoMessage()           {}
func (m *QueryRiskScoreRequest) Reset()                  { *m = QueryRiskScoreRequest{} }
func (m *QueryRiskScoreRequest) String() string          { return fmt.Sprintf("query_risk_score: %d", m.PoolID) }
func (m *QueryRiskScoreRequest) XXX_MessageName() string { return "syreen.dex.QueryRiskScoreRequest" }

type QueryRiskScoreResponse struct {
	RiskScore json.RawMessage `json:"risk_score"`
}

func (m *QueryRiskScoreResponse) ProtoMessage()           {}
func (m *QueryRiskScoreResponse) Reset()                  { *m = QueryRiskScoreResponse{} }
func (m *QueryRiskScoreResponse) String() string          { return "risk_score_response" }
func (m *QueryRiskScoreResponse) XXX_MessageName() string { return "syreen.dex.QueryRiskScoreResponse" }

type QueryTradeRiskRequest struct {
	PoolID     uint64   `json:"pool_id"`
	InputDenom string   `json:"input_denom"`
	Amount     math.Int `json:"amount"`
}

func (m *QueryTradeRiskRequest) ProtoMessage()           {}
func (m *QueryTradeRiskRequest) Reset()                  { *m = QueryTradeRiskRequest{} }
func (m *QueryTradeRiskRequest) String() string          { return fmt.Sprintf("query_trade_risk: pool %d", m.PoolID) }
func (m *QueryTradeRiskRequest) XXX_MessageName() string { return "syreen.dex.QueryTradeRiskRequest" }

type QueryTradeRiskResponse struct {
	TradeRisk json.RawMessage `json:"trade_risk"`
}

func (m *QueryTradeRiskResponse) ProtoMessage()           {}
func (m *QueryTradeRiskResponse) Reset()                  { *m = QueryTradeRiskResponse{} }
func (m *QueryTradeRiskResponse) String() string          { return "trade_risk_response" }
func (m *QueryTradeRiskResponse) XXX_MessageName() string { return "syreen.dex.QueryTradeRiskResponse" }

type QueryRiskHistoryRequest struct {
	PoolID uint64 `json:"pool_id"`
}

func (m *QueryRiskHistoryRequest) ProtoMessage()           {}
func (m *QueryRiskHistoryRequest) Reset()                  { *m = QueryRiskHistoryRequest{} }
func (m *QueryRiskHistoryRequest) String() string          { return fmt.Sprintf("query_risk_history: %d", m.PoolID) }
func (m *QueryRiskHistoryRequest) XXX_MessageName() string { return "syreen.dex.QueryRiskHistoryRequest" }

type QueryRiskHistoryResponse struct {
	RiskHistory json.RawMessage `json:"risk_history"`
}

func (m *QueryRiskHistoryResponse) ProtoMessage()           {}
func (m *QueryRiskHistoryResponse) Reset()                  { *m = QueryRiskHistoryResponse{} }
func (m *QueryRiskHistoryResponse) String() string          { return "risk_history_response" }
func (m *QueryRiskHistoryResponse) XXX_MessageName() string { return "syreen.dex.QueryRiskHistoryResponse" }

// ---------------------------------------------------------------------------
// Sentiment Oracle query request/response types
// ---------------------------------------------------------------------------

type QuerySentimentRequest struct {
	PoolID uint64 `json:"pool_id"`
}

func (m *QuerySentimentRequest) ProtoMessage()           {}
func (m *QuerySentimentRequest) Reset()                  { *m = QuerySentimentRequest{} }
func (m *QuerySentimentRequest) String() string          { return fmt.Sprintf("query_sentiment: %d", m.PoolID) }
func (m *QuerySentimentRequest) XXX_MessageName() string { return "syreen.dex.QuerySentimentRequest" }

type QuerySentimentResponse struct {
	PoolID         uint64      `json:"pool_id"`
	FearGreedIndex int64       `json:"fear_greed_index"`
	FearGreedLabel string      `json:"fear_greed_label"`
	Mood           string      `json:"mood"`
	Indicators     interface{} `json:"indicators"`
	UpdatedAt      int64       `json:"updated_at_height"`
	BlocksAnalyzed int64       `json:"blocks_analyzed"`
}

func (m *QuerySentimentResponse) ProtoMessage()           {}
func (m *QuerySentimentResponse) Reset()                  { *m = QuerySentimentResponse{} }
func (m *QuerySentimentResponse) String() string          { return fmt.Sprintf("sentiment: pool %d F&G %d", m.PoolID, m.FearGreedIndex) }
func (m *QuerySentimentResponse) XXX_MessageName() string { return "syreen.dex.QuerySentimentResponse" }

type QuerySentimentHistoryRequest struct {
	PoolID uint64 `json:"pool_id"`
	Limit  int64  `json:"limit"`
}

func (m *QuerySentimentHistoryRequest) ProtoMessage()           {}
func (m *QuerySentimentHistoryRequest) Reset()                  { *m = QuerySentimentHistoryRequest{} }
func (m *QuerySentimentHistoryRequest) String() string          { return fmt.Sprintf("query_sentiment_history: %d", m.PoolID) }
func (m *QuerySentimentHistoryRequest) XXX_MessageName() string { return "syreen.dex.QuerySentimentHistoryRequest" }

type QuerySentimentHistoryResponse struct {
	PoolID  uint64      `json:"pool_id"`
	History interface{} `json:"history"`
}

func (m *QuerySentimentHistoryResponse) ProtoMessage()           {}
func (m *QuerySentimentHistoryResponse) Reset()                  { *m = QuerySentimentHistoryResponse{} }
func (m *QuerySentimentHistoryResponse) String() string          { return fmt.Sprintf("sentiment_history: pool %d", m.PoolID) }
func (m *QuerySentimentHistoryResponse) XXX_MessageName() string { return "syreen.dex.QuerySentimentHistoryResponse" }

type QuerySentimentAlertsRequest struct {
	PoolID uint64 `json:"pool_id"`
}

func (m *QuerySentimentAlertsRequest) ProtoMessage()           {}
func (m *QuerySentimentAlertsRequest) Reset()                  { *m = QuerySentimentAlertsRequest{} }
func (m *QuerySentimentAlertsRequest) String() string          { return fmt.Sprintf("query_sentiment_alerts: %d", m.PoolID) }
func (m *QuerySentimentAlertsRequest) XXX_MessageName() string { return "syreen.dex.QuerySentimentAlertsRequest" }

type QuerySentimentAlertsResponse struct {
	PoolID uint64      `json:"pool_id"`
	Alerts interface{} `json:"alerts"`
}

func (m *QuerySentimentAlertsResponse) ProtoMessage()           {}
func (m *QuerySentimentAlertsResponse) Reset()                  { *m = QuerySentimentAlertsResponse{} }
func (m *QuerySentimentAlertsResponse) String() string          { return fmt.Sprintf("sentiment_alerts: pool %d", m.PoolID) }
func (m *QuerySentimentAlertsResponse) XXX_MessageName() string { return "syreen.dex.QuerySentimentAlertsResponse" }

// ---------------------------------------------------------------------------
// Oracle Composite query request/response types
// ---------------------------------------------------------------------------

type QueryOracleCompositeRequest struct {
	PoolID uint64 `json:"pool_id"`
}

func (m *QueryOracleCompositeRequest) ProtoMessage()           {}
func (m *QueryOracleCompositeRequest) Reset()                  { *m = QueryOracleCompositeRequest{} }
func (m *QueryOracleCompositeRequest) String() string          { return fmt.Sprintf("query_oracle_composite: %d", m.PoolID) }
func (m *QueryOracleCompositeRequest) XXX_MessageName() string { return "syreen.dex.QueryOracleCompositeRequest" }

type QueryOracleCompositeResponse struct {
	Oracle json.RawMessage `json:"oracle"`
}

func (m *QueryOracleCompositeResponse) ProtoMessage()           {}
func (m *QueryOracleCompositeResponse) Reset()                  { *m = QueryOracleCompositeResponse{} }
func (m *QueryOracleCompositeResponse) String() string          { return "oracle_composite" }
func (m *QueryOracleCompositeResponse) XXX_MessageName() string { return "syreen.dex.QueryOracleCompositeResponse" }

type QueryWhaleAlertsRequest struct {
	PoolID uint64 `json:"pool_id"`
	Limit  int64  `json:"limit"`
}

func (m *QueryWhaleAlertsRequest) ProtoMessage()           {}
func (m *QueryWhaleAlertsRequest) Reset()                  { *m = QueryWhaleAlertsRequest{} }
func (m *QueryWhaleAlertsRequest) String() string          { return fmt.Sprintf("query_whale_alerts: %d", m.PoolID) }
func (m *QueryWhaleAlertsRequest) XXX_MessageName() string { return "syreen.dex.QueryWhaleAlertsRequest" }

type QueryWhaleAlertsResponse struct {
	Alerts json.RawMessage `json:"alerts"`
}

func (m *QueryWhaleAlertsResponse) ProtoMessage()           {}
func (m *QueryWhaleAlertsResponse) Reset()                  { *m = QueryWhaleAlertsResponse{} }
func (m *QueryWhaleAlertsResponse) String() string          { return "whale_alerts" }
func (m *QueryWhaleAlertsResponse) XXX_MessageName() string { return "syreen.dex.QueryWhaleAlertsResponse" }

// ---------------------------------------------------------------------------
// Referral query request/response types
// ---------------------------------------------------------------------------

type QueryReferralInfoRequest struct {
	Address string `json:"address"`
}

func (m *QueryReferralInfoRequest) ProtoMessage()           {}
func (m *QueryReferralInfoRequest) Reset()                  { *m = QueryReferralInfoRequest{} }
func (m *QueryReferralInfoRequest) String() string          { return fmt.Sprintf("query_referral_info: %s", m.Address) }
func (m *QueryReferralInfoRequest) XXX_MessageName() string { return "syreen.dex.QueryReferralInfoRequest" }

type QueryReferralInfoResponse struct {
	Address       string `json:"address"`
	ReferralCode  string `json:"referral_code"`
	TotalReferrals int64  `json:"total_referrals"`
	TotalVolume   string `json:"total_volume"`
	TotalEarned   string `json:"total_earned"`
	Tier          string `json:"tier"`
	ReferrerPct   int64  `json:"referrer_pct"`
	DiscountPct   int64  `json:"discount_pct"`
}

func (m *QueryReferralInfoResponse) ProtoMessage()           {}
func (m *QueryReferralInfoResponse) Reset()                  { *m = QueryReferralInfoResponse{} }
func (m *QueryReferralInfoResponse) String() string          { return fmt.Sprintf("referral_info: %s tier=%s", m.Address, m.Tier) }
func (m *QueryReferralInfoResponse) XXX_MessageName() string { return "syreen.dex.QueryReferralInfoResponse" }

type QueryReferralListRequest struct {
	Address string `json:"address"`
}

func (m *QueryReferralListRequest) ProtoMessage()           {}
func (m *QueryReferralListRequest) Reset()                  { *m = QueryReferralListRequest{} }
func (m *QueryReferralListRequest) String() string          { return fmt.Sprintf("query_referral_list: %s", m.Address) }
func (m *QueryReferralListRequest) XXX_MessageName() string { return "syreen.dex.QueryReferralListRequest" }

type ReferredUserInfo struct {
	Address     string `json:"address"`
	TotalVolume string `json:"total_volume"`
	TotalSaved  string `json:"total_saved"`
}

type QueryReferralListResponse struct {
	Referrer  string             `json:"referrer"`
	Referrals []ReferredUserInfo `json:"referrals"`
}

func (m *QueryReferralListResponse) ProtoMessage()           {}
func (m *QueryReferralListResponse) Reset()                  { *m = QueryReferralListResponse{} }
func (m *QueryReferralListResponse) String() string          { return fmt.Sprintf("referral_list: %s count=%d", m.Referrer, len(m.Referrals)) }
func (m *QueryReferralListResponse) XXX_MessageName() string { return "syreen.dex.QueryReferralListResponse" }

type QueryReferralCodeLookupRequest struct {
	Code string `json:"code"`
}

func (m *QueryReferralCodeLookupRequest) ProtoMessage()           {}
func (m *QueryReferralCodeLookupRequest) Reset()                  { *m = QueryReferralCodeLookupRequest{} }
func (m *QueryReferralCodeLookupRequest) String() string          { return fmt.Sprintf("query_referral_code: %s", m.Code) }
func (m *QueryReferralCodeLookupRequest) XXX_MessageName() string { return "syreen.dex.QueryReferralCodeLookupRequest" }

type QueryReferralCodeLookupResponse struct {
	Code          string `json:"code"`
	ReferrerAddr  string `json:"referrer_address"`
	TotalReferrals int64  `json:"total_referrals"`
	Tier          string `json:"tier"`
}

func (m *QueryReferralCodeLookupResponse) ProtoMessage()           {}
func (m *QueryReferralCodeLookupResponse) Reset()                  { *m = QueryReferralCodeLookupResponse{} }
func (m *QueryReferralCodeLookupResponse) String() string          { return fmt.Sprintf("referral_code: %s -> %s", m.Code, m.ReferrerAddr) }
func (m *QueryReferralCodeLookupResponse) XXX_MessageName() string { return "syreen.dex.QueryReferralCodeLookupResponse" }

type QueryReferralEarningsRequest struct {
	Address string `json:"address"`
}

func (m *QueryReferralEarningsRequest) ProtoMessage()           {}
func (m *QueryReferralEarningsRequest) Reset()                  { *m = QueryReferralEarningsRequest{} }
func (m *QueryReferralEarningsRequest) String() string          { return fmt.Sprintf("query_referral_earnings: %s", m.Address) }
func (m *QueryReferralEarningsRequest) XXX_MessageName() string { return "syreen.dex.QueryReferralEarningsRequest" }

type ReferralEarningEntry struct {
	Height   int64  `json:"height"`
	Trader   string `json:"trader"`
	Denom    string `json:"denom"`
	Amount   string `json:"amount"`
	PoolID   uint64 `json:"pool_id"`
}

type QueryReferralEarningsResponse struct {
	Address       string                 `json:"address"`
	TotalEarned   string                 `json:"total_earned"`
	Entries       []ReferralEarningEntry `json:"entries"`
}

func (m *QueryReferralEarningsResponse) ProtoMessage()           {}
func (m *QueryReferralEarningsResponse) Reset()                  { *m = QueryReferralEarningsResponse{} }
func (m *QueryReferralEarningsResponse) String() string          { return fmt.Sprintf("referral_earnings: %s total=%s", m.Address, m.TotalEarned) }
func (m *QueryReferralEarningsResponse) XXX_MessageName() string { return "syreen.dex.QueryReferralEarningsResponse" }

// ---------------------------------------------------------------------------
// Copy Trading query types
// ---------------------------------------------------------------------------

type QueryLeaderboardRequest struct {
	Limit int64 `json:"limit"`
}

func (m *QueryLeaderboardRequest) ProtoMessage()           {}
func (m *QueryLeaderboardRequest) Reset()                  { *m = QueryLeaderboardRequest{} }
func (m *QueryLeaderboardRequest) String() string          { return "query_leaderboard" }
func (m *QueryLeaderboardRequest) XXX_MessageName() string { return "syreen.dex.QueryLeaderboardRequest" }

type QueryLeaderboardResponse struct {
	Traders json.RawMessage `json:"traders"`
}

func (m *QueryLeaderboardResponse) ProtoMessage()           {}
func (m *QueryLeaderboardResponse) Reset()                  { *m = QueryLeaderboardResponse{} }
func (m *QueryLeaderboardResponse) String() string          { return "leaderboard_response" }
func (m *QueryLeaderboardResponse) XXX_MessageName() string { return "syreen.dex.QueryLeaderboardResponse" }

type QueryTraderStatsRequest struct {
	Address string `json:"address"`
}

func (m *QueryTraderStatsRequest) ProtoMessage()           {}
func (m *QueryTraderStatsRequest) Reset()                  { *m = QueryTraderStatsRequest{} }
func (m *QueryTraderStatsRequest) String() string          { return fmt.Sprintf("query_trader_stats: %s", m.Address) }
func (m *QueryTraderStatsRequest) XXX_MessageName() string { return "syreen.dex.QueryTraderStatsRequest" }

type QueryTraderStatsResponse struct {
	Stats json.RawMessage `json:"stats"`
}

func (m *QueryTraderStatsResponse) ProtoMessage()           {}
func (m *QueryTraderStatsResponse) Reset()                  { *m = QueryTraderStatsResponse{} }
func (m *QueryTraderStatsResponse) String() string          { return "trader_stats_response" }
func (m *QueryTraderStatsResponse) XXX_MessageName() string { return "syreen.dex.QueryTraderStatsResponse" }

type QueryTraderFollowersRequest struct {
	Address string `json:"address"`
}

func (m *QueryTraderFollowersRequest) ProtoMessage()           {}
func (m *QueryTraderFollowersRequest) Reset()                  { *m = QueryTraderFollowersRequest{} }
func (m *QueryTraderFollowersRequest) String() string          { return fmt.Sprintf("query_trader_followers: %s", m.Address) }
func (m *QueryTraderFollowersRequest) XXX_MessageName() string { return "syreen.dex.QueryTraderFollowersRequest" }

type QueryTraderFollowersResponse struct {
	Followers json.RawMessage `json:"followers"`
}

func (m *QueryTraderFollowersResponse) ProtoMessage()           {}
func (m *QueryTraderFollowersResponse) Reset()                  { *m = QueryTraderFollowersResponse{} }
func (m *QueryTraderFollowersResponse) String() string          { return "trader_followers_response" }
func (m *QueryTraderFollowersResponse) XXX_MessageName() string { return "syreen.dex.QueryTraderFollowersResponse" }

type QueryCopyFollowingRequest struct {
	Address string `json:"address"`
}

func (m *QueryCopyFollowingRequest) ProtoMessage()           {}
func (m *QueryCopyFollowingRequest) Reset()                  { *m = QueryCopyFollowingRequest{} }
func (m *QueryCopyFollowingRequest) String() string          { return fmt.Sprintf("query_copy_following: %s", m.Address) }
func (m *QueryCopyFollowingRequest) XXX_MessageName() string { return "syreen.dex.QueryCopyFollowingRequest" }

type QueryCopyFollowingResponse struct {
	Following json.RawMessage `json:"following"`
}

func (m *QueryCopyFollowingResponse) ProtoMessage()           {}
func (m *QueryCopyFollowingResponse) Reset()                  { *m = QueryCopyFollowingResponse{} }
func (m *QueryCopyFollowingResponse) String() string          { return "copy_following_response" }
func (m *QueryCopyFollowingResponse) XXX_MessageName() string { return "syreen.dex.QueryCopyFollowingResponse" }

type QueryCopyHistoryRequest struct {
	Address string `json:"address"`
}

func (m *QueryCopyHistoryRequest) ProtoMessage()           {}
func (m *QueryCopyHistoryRequest) Reset()                  { *m = QueryCopyHistoryRequest{} }
func (m *QueryCopyHistoryRequest) String() string          { return fmt.Sprintf("query_copy_history: %s", m.Address) }
func (m *QueryCopyHistoryRequest) XXX_MessageName() string { return "syreen.dex.QueryCopyHistoryRequest" }

type QueryCopyHistoryResponse struct {
	History json.RawMessage `json:"history"`
}

func (m *QueryCopyHistoryResponse) ProtoMessage()           {}
func (m *QueryCopyHistoryResponse) Reset()                  { *m = QueryCopyHistoryResponse{} }
func (m *QueryCopyHistoryResponse) String() string          { return "copy_history_response" }
func (m *QueryCopyHistoryResponse) XXX_MessageName() string { return "syreen.dex.QueryCopyHistoryResponse" }

// ---------------------------------------------------------------------------
// Time-Weighted Governance query request/response types
// ---------------------------------------------------------------------------

type QueryTimeWeightedPowerRequest struct {
	Address string `json:"address"`
}

func (m *QueryTimeWeightedPowerRequest) ProtoMessage()           {}
func (m *QueryTimeWeightedPowerRequest) Reset()                  { *m = QueryTimeWeightedPowerRequest{} }
func (m *QueryTimeWeightedPowerRequest) String() string          { return fmt.Sprintf("query_time_weighted_power: %s", m.Address) }
func (m *QueryTimeWeightedPowerRequest) XXX_MessageName() string { return "syreen.dex.QueryTimeWeightedPowerRequest" }

type QueryTimeWeightedPowerResponse struct {
	PowerInfo json.RawMessage `json:"power_info"`
}

func (m *QueryTimeWeightedPowerResponse) ProtoMessage()           {}
func (m *QueryTimeWeightedPowerResponse) Reset()                  { *m = QueryTimeWeightedPowerResponse{} }
func (m *QueryTimeWeightedPowerResponse) String() string          { return "time_weighted_power_response" }
func (m *QueryTimeWeightedPowerResponse) XXX_MessageName() string { return "syreen.dex.QueryTimeWeightedPowerResponse" }

// ---------------------------------------------------------------------------
// Trader Record query request/response types
// ---------------------------------------------------------------------------

type QueryTraderRecordRequest struct {
	Address string `json:"address"`
}

func (m *QueryTraderRecordRequest) ProtoMessage()           {}
func (m *QueryTraderRecordRequest) Reset()                  { *m = QueryTraderRecordRequest{} }
func (m *QueryTraderRecordRequest) String() string          { return fmt.Sprintf("query_trader_record: %s", m.Address) }
func (m *QueryTraderRecordRequest) XXX_MessageName() string { return "syreen.dex.QueryTraderRecordRequest" }

type QueryTraderRecordResponse struct {
	Record json.RawMessage `json:"record"`
}

func (m *QueryTraderRecordResponse) ProtoMessage()           {}
func (m *QueryTraderRecordResponse) Reset()                  { *m = QueryTraderRecordResponse{} }
func (m *QueryTraderRecordResponse) String() string          { return "trader_record_response" }
func (m *QueryTraderRecordResponse) XXX_MessageName() string { return "syreen.dex.QueryTraderRecordResponse" }

type QueryTraderLeaderboardRequest struct {
	SortBy string `json:"sort_by"`
	Limit  int64  `json:"limit"`
}

func (m *QueryTraderLeaderboardRequest) ProtoMessage()           {}
func (m *QueryTraderLeaderboardRequest) Reset()                  { *m = QueryTraderLeaderboardRequest{} }
func (m *QueryTraderLeaderboardRequest) String() string          { return "query_trader_leaderboard" }
func (m *QueryTraderLeaderboardRequest) XXX_MessageName() string { return "syreen.dex.QueryTraderLeaderboardRequest" }

type QueryTraderLeaderboardResponse struct {
	Traders json.RawMessage `json:"traders"`
}

func (m *QueryTraderLeaderboardResponse) ProtoMessage()           {}
func (m *QueryTraderLeaderboardResponse) Reset()                  { *m = QueryTraderLeaderboardResponse{} }
func (m *QueryTraderLeaderboardResponse) String() string          { return "trader_leaderboard_response" }
func (m *QueryTraderLeaderboardResponse) XXX_MessageName() string { return "syreen.dex.QueryTraderLeaderboardResponse" }

// suppress unused import warnings
var _ math.Int

// ---------------------------------------------------------------------------
// Order Book query request/response types
// ---------------------------------------------------------------------------

type QueryOrderBookRequest struct {
	PoolID uint64 `json:"pool_id"`
}

func (m *QueryOrderBookRequest) ProtoMessage()           {}
func (m *QueryOrderBookRequest) Reset()                  { *m = QueryOrderBookRequest{} }
func (m *QueryOrderBookRequest) String() string          { return fmt.Sprintf("query_orderbook: %d", m.PoolID) }
func (m *QueryOrderBookRequest) XXX_MessageName() string { return "syreen.dex.QueryOrderBookRequest" }

type QueryOrderBookResponse struct {
	PoolID uint64           `json:"pool_id"`
	Bids   []OrderBookLevel `json:"bids"`
	Asks   []OrderBookLevel `json:"asks"`
}

func (m *QueryOrderBookResponse) ProtoMessage()           {}
func (m *QueryOrderBookResponse) Reset()                  { *m = QueryOrderBookResponse{} }
func (m *QueryOrderBookResponse) String() string          { return fmt.Sprintf("orderbook: pool %d", m.PoolID) }
func (m *QueryOrderBookResponse) XXX_MessageName() string { return "syreen.dex.QueryOrderBookResponse" }

type QueryOrderBookDepthRequest struct {
	PoolID uint64 `json:"pool_id"`
}

func (m *QueryOrderBookDepthRequest) ProtoMessage()           {}
func (m *QueryOrderBookDepthRequest) Reset()                  { *m = QueryOrderBookDepthRequest{} }
func (m *QueryOrderBookDepthRequest) String() string          { return fmt.Sprintf("query_orderbook_depth: %d", m.PoolID) }
func (m *QueryOrderBookDepthRequest) XXX_MessageName() string { return "syreen.dex.QueryOrderBookDepthRequest" }

type QueryOrderBookDepthResponse struct {
	PoolID uint64           `json:"pool_id"`
	Bids   []OrderBookLevel `json:"bids"`
	Asks   []OrderBookLevel `json:"asks"`
}

func (m *QueryOrderBookDepthResponse) ProtoMessage()           {}
func (m *QueryOrderBookDepthResponse) Reset()                  { *m = QueryOrderBookDepthResponse{} }
func (m *QueryOrderBookDepthResponse) String() string          { return fmt.Sprintf("orderbook_depth: pool %d", m.PoolID) }
func (m *QueryOrderBookDepthResponse) XXX_MessageName() string { return "syreen.dex.QueryOrderBookDepthResponse" }

type QueryUserOrdersRequest struct {
	Address string `json:"address"`
}

func (m *QueryUserOrdersRequest) ProtoMessage()           {}
func (m *QueryUserOrdersRequest) Reset()                  { *m = QueryUserOrdersRequest{} }
func (m *QueryUserOrdersRequest) String() string          { return fmt.Sprintf("query_user_orders: %s", m.Address) }
func (m *QueryUserOrdersRequest) XXX_MessageName() string { return "syreen.dex.QueryUserOrdersRequest" }

type QueryUserOrdersResponse struct {
	Orders []Order `json:"orders"`
}

func (m *QueryUserOrdersResponse) ProtoMessage()           {}
func (m *QueryUserOrdersResponse) Reset()                  { *m = QueryUserOrdersResponse{} }
func (m *QueryUserOrdersResponse) String() string          { return fmt.Sprintf("user_orders: %d", len(m.Orders)) }
func (m *QueryUserOrdersResponse) XXX_MessageName() string { return "syreen.dex.QueryUserOrdersResponse" }

type QueryOrderDetailRequest struct {
	OrderID uint64 `json:"order_id"`
}

func (m *QueryOrderDetailRequest) ProtoMessage()           {}
func (m *QueryOrderDetailRequest) Reset()                  { *m = QueryOrderDetailRequest{} }
func (m *QueryOrderDetailRequest) String() string          { return fmt.Sprintf("query_order_detail: %d", m.OrderID) }
func (m *QueryOrderDetailRequest) XXX_MessageName() string { return "syreen.dex.QueryOrderDetailRequest" }

type QueryOrderDetailResponse struct {
	Order Order `json:"order"`
}

func (m *QueryOrderDetailResponse) ProtoMessage()           {}
func (m *QueryOrderDetailResponse) Reset()                  { *m = QueryOrderDetailResponse{} }
func (m *QueryOrderDetailResponse) String() string          { return "order_detail_response" }
func (m *QueryOrderDetailResponse) XXX_MessageName() string { return "syreen.dex.QueryOrderDetailResponse" }

type QueryTradesRequest struct {
	PoolID uint64 `json:"pool_id"`
	Limit  int64  `json:"limit"`
}

func (m *QueryTradesRequest) ProtoMessage()           {}
func (m *QueryTradesRequest) Reset()                  { *m = QueryTradesRequest{} }
func (m *QueryTradesRequest) String() string          { return fmt.Sprintf("query_trades: %d", m.PoolID) }
func (m *QueryTradesRequest) XXX_MessageName() string { return "syreen.dex.QueryTradesRequest" }

type QueryTradesResponse struct {
	Trades []Trade `json:"trades"`
}

func (m *QueryTradesResponse) ProtoMessage()           {}
func (m *QueryTradesResponse) Reset()                  { *m = QueryTradesResponse{} }
func (m *QueryTradesResponse) String() string          { return fmt.Sprintf("trades: %d", len(m.Trades)) }
func (m *QueryTradesResponse) XXX_MessageName() string { return "syreen.dex.QueryTradesResponse" }

// ---------------------------------------------------------------------------
// IBC Order Queries
// ---------------------------------------------------------------------------

type QueryIBCOrdersRequest struct {
	SourceChain string `json:"source_chain"`
	Sender      string `json:"sender,omitempty"`
}

func (m *QueryIBCOrdersRequest) ProtoMessage()           {}
func (m *QueryIBCOrdersRequest) Reset()                  { *m = QueryIBCOrdersRequest{} }
func (m *QueryIBCOrdersRequest) String() string          { return fmt.Sprintf("query_ibc_orders: %s", m.SourceChain) }
func (m *QueryIBCOrdersRequest) XXX_MessageName() string { return "syreen.dex.QueryIBCOrdersRequest" }

type QueryIBCOrdersResponse struct {
	Orders json.RawMessage `json:"orders"`
}

func (m *QueryIBCOrdersResponse) ProtoMessage()           {}
func (m *QueryIBCOrdersResponse) Reset()                  { *m = QueryIBCOrdersResponse{} }
func (m *QueryIBCOrdersResponse) String() string          { return "ibc_orders_response" }
func (m *QueryIBCOrdersResponse) XXX_MessageName() string { return "syreen.dex.QueryIBCOrdersResponse" }

// ---------------------------------------------------------------------------
// Dynamic Fee Query Types
// ---------------------------------------------------------------------------

type QueryPoolFeeRequest struct {
	PoolID uint64 `json:"pool_id"`
}

func (m *QueryPoolFeeRequest) ProtoMessage()           {}
func (m *QueryPoolFeeRequest) Reset()                  { *m = QueryPoolFeeRequest{} }
func (m *QueryPoolFeeRequest) String() string          { return fmt.Sprintf("query_pool_fee: %d", m.PoolID) }
func (m *QueryPoolFeeRequest) XXX_MessageName() string { return "syreen.dex.QueryPoolFeeRequest" }

type QueryPoolFeeResponse struct {
	PoolID               uint64         `json:"pool_id"`
	EffectiveFee         math.LegacyDec `json:"effective_fee"`
	SwapFee              math.LegacyDec `json:"swap_fee"`
	VolatilityFeeEnabled bool           `json:"volatility_fee_enabled"`
	BaseFee              math.LegacyDec `json:"base_fee,omitempty"`
	MaxFee               math.LegacyDec `json:"max_fee,omitempty"`
	VolatilityMultiplier math.LegacyDec `json:"volatility_multiplier,omitempty"`
	CurrentVolatility    math.LegacyDec `json:"current_volatility"`
}

func (m *QueryPoolFeeResponse) ProtoMessage()           {}
func (m *QueryPoolFeeResponse) Reset()                  { *m = QueryPoolFeeResponse{} }
func (m *QueryPoolFeeResponse) String() string          { return fmt.Sprintf("pool_fee: %d effective=%s", m.PoolID, m.EffectiveFee) }
func (m *QueryPoolFeeResponse) XXX_MessageName() string { return "syreen.dex.QueryPoolFeeResponse" }

// ---------------------------------------------------------------------------
// QueryServer interface
// ---------------------------------------------------------------------------

type QueryServer interface {
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	Pool(context.Context, *QueryPoolRequest) (*QueryPoolResponse, error)
	Pools(context.Context, *QueryPoolsRequest) (*QueryPoolsResponse, error)
	Quote(context.Context, *QueryQuoteRequest) (*QueryQuoteResponse, error)
	SpotPrice(context.Context, *QuerySpotPriceRequest) (*QuerySpotPriceResponse, error)
	OptimalRoute(context.Context, *QueryOptimalRouteRequest) (*QueryOptimalRouteResponse, error)
	Routes(context.Context, *QueryRoutesRequest) (*QueryRoutesResponse, error)
	Signals(context.Context, *QuerySignalsRequest) (*QuerySignalsResponse, error)
	SignalHistory(context.Context, *QuerySignalHistoryRequest) (*QuerySignalHistoryResponse, error)
	RiskScore(context.Context, *QueryRiskScoreRequest) (*QueryRiskScoreResponse, error)
	TradeRisk(context.Context, *QueryTradeRiskRequest) (*QueryTradeRiskResponse, error)
	RiskHistory(context.Context, *QueryRiskHistoryRequest) (*QueryRiskHistoryResponse, error)
	Sentiment(context.Context, *QuerySentimentRequest) (*QuerySentimentResponse, error)
	SentimentHistory(context.Context, *QuerySentimentHistoryRequest) (*QuerySentimentHistoryResponse, error)
	SentimentAlerts(context.Context, *QuerySentimentAlertsRequest) (*QuerySentimentAlertsResponse, error)
	OracleComposite(context.Context, *QueryOracleCompositeRequest) (*QueryOracleCompositeResponse, error)
	WhaleAlerts(context.Context, *QueryWhaleAlertsRequest) (*QueryWhaleAlertsResponse, error)
	ReferralInfo(context.Context, *QueryReferralInfoRequest) (*QueryReferralInfoResponse, error)
	ReferralList(context.Context, *QueryReferralListRequest) (*QueryReferralListResponse, error)
	ReferralCodeLookup(context.Context, *QueryReferralCodeLookupRequest) (*QueryReferralCodeLookupResponse, error)
	ReferralEarnings(context.Context, *QueryReferralEarningsRequest) (*QueryReferralEarningsResponse, error)
	Leaderboard(context.Context, *QueryLeaderboardRequest) (*QueryLeaderboardResponse, error)
	TraderStats(context.Context, *QueryTraderStatsRequest) (*QueryTraderStatsResponse, error)
	TraderFollowers(context.Context, *QueryTraderFollowersRequest) (*QueryTraderFollowersResponse, error)
	CopyFollowing(context.Context, *QueryCopyFollowingRequest) (*QueryCopyFollowingResponse, error)
	CopyHistory(context.Context, *QueryCopyHistoryRequest) (*QueryCopyHistoryResponse, error)
	// Order Book queries
	OrderBook(context.Context, *QueryOrderBookRequest) (*QueryOrderBookResponse, error)
	OrderBookDepth(context.Context, *QueryOrderBookDepthRequest) (*QueryOrderBookDepthResponse, error)
	UserOrders(context.Context, *QueryUserOrdersRequest) (*QueryUserOrdersResponse, error)
	OrderDetail(context.Context, *QueryOrderDetailRequest) (*QueryOrderDetailResponse, error)
	Trades(context.Context, *QueryTradesRequest) (*QueryTradesResponse, error)
	// Time-Weighted Governance
	TimeWeightedPower(context.Context, *QueryTimeWeightedPowerRequest) (*QueryTimeWeightedPowerResponse, error)
	// Trader Track Records
	TraderRecord(context.Context, *QueryTraderRecordRequest) (*QueryTraderRecordResponse, error)
	TraderLeaderboard(context.Context, *QueryTraderLeaderboardRequest) (*QueryTraderLeaderboardResponse, error)
	// IBC Cross-Chain Order queries
	IBCOrders(context.Context, *QueryIBCOrdersRequest) (*QueryIBCOrdersResponse, error)
	// Dynamic Fee queries
	PoolFee(context.Context, *QueryPoolFeeRequest) (*QueryPoolFeeResponse, error)
}

// ---------------------------------------------------------------------------
// RegisterMsgServer
// ---------------------------------------------------------------------------

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

// ---------------------------------------------------------------------------
// RegisterQueryServer
// ---------------------------------------------------------------------------

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

// ---------------------------------------------------------------------------
// Msg service handlers
// ---------------------------------------------------------------------------

func _Msg_CreatePool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreatePool)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreatePool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Msg/CreatePool"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreatePool(ctx, req.(*MsgCreatePool))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_AddLiquidity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAddLiquidity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AddLiquidity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Msg/AddLiquidity"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AddLiquidity(ctx, req.(*MsgAddLiquidity))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RemoveLiquidity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRemoveLiquidity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RemoveLiquidity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Msg/RemoveLiquidity"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RemoveLiquidity(ctx, req.(*MsgRemoveLiquidity))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Swap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSwap)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Swap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Msg/Swap"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Swap(ctx, req.(*MsgSwap))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateReferralCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateReferralCode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateReferralCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Msg/CreateReferralCode"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateReferralCode(ctx, req.(*MsgCreateReferralCode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RegisterReferral_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterReferral)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RegisterReferral(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Msg/RegisterReferral"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RegisterReferral(ctx, req.(*MsgRegisterReferral))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_PlaceOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgPlaceOrder)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).PlaceOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Msg/PlaceOrder"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).PlaceOrder(ctx, req.(*MsgPlaceOrder))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CancelOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCancelOrder)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CancelOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Msg/CancelOrder"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CancelOrder(ctx, req.(*MsgCancelOrder))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ModifyOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgModifyOrder)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ModifyOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Msg/ModifyOrder"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ModifyOrder(ctx, req.(*MsgModifyOrder))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_MultiHopSwap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgMultiHopSwap)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).MultiHopSwap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Msg/MultiHopSwap"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).MultiHopSwap(ctx, req.(*MsgMultiHopSwap))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.dex.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreatePool", Handler: _Msg_CreatePool_Handler},
		{MethodName: "AddLiquidity", Handler: _Msg_AddLiquidity_Handler},
		{MethodName: "RemoveLiquidity", Handler: _Msg_RemoveLiquidity_Handler},
		{MethodName: "Swap", Handler: _Msg_Swap_Handler},
		{MethodName: "CreateReferralCode", Handler: _Msg_CreateReferralCode_Handler},
		{MethodName: "RegisterReferral", Handler: _Msg_RegisterReferral_Handler},
		{MethodName: "PlaceOrder", Handler: _Msg_PlaceOrder_Handler},
		{MethodName: "CancelOrder", Handler: _Msg_CancelOrder_Handler},
		{MethodName: "ModifyOrder", Handler: _Msg_ModifyOrder_Handler},
		{MethodName: "FollowTrader", Handler: _Msg_FollowTrader_Handler},
		{MethodName: "UnfollowTrader", Handler: _Msg_UnfollowTrader_Handler},
		{MethodName: "UpdateCopySettings", Handler: _Msg_UpdateCopySettings_Handler},
		{MethodName: "ClaimReferralRewards", Handler: _Msg_ClaimReferralRewards_Handler},
		{MethodName: "MultiHopSwap", Handler: _Msg_MultiHopSwap_Handler},
		{MethodName: "SetPoolFeeConfig", Handler: _Msg_SetPoolFeeConfig_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/dex/tx.proto",
}

func _Msg_SetPoolFeeConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetPoolFeeConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetPoolFeeConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Msg/SetPoolFeeConfig"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetPoolFeeConfig(ctx, req.(*MsgSetPoolFeeConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ClaimReferralRewards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClaimReferralRewards)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ClaimReferralRewards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Msg/ClaimReferralRewards"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClaimReferralRewards(ctx, req.(*MsgClaimReferralRewards))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_FollowTrader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgFollowTrader)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).FollowTrader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Msg/FollowTrader"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).FollowTrader(ctx, req.(*MsgFollowTrader))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UnfollowTrader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUnfollowTrader)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UnfollowTrader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Msg/UnfollowTrader"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UnfollowTrader(ctx, req.(*MsgUnfollowTrader))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateCopySettings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateCopySettings)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateCopySettings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Msg/UpdateCopySettings"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateCopySettings(ctx, req.(*MsgUpdateCopySettings))
	}
	return interceptor(ctx, in, info, handler)
}

// ---------------------------------------------------------------------------
// Query service handlers
// ---------------------------------------------------------------------------

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/Params"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Pool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPoolRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Pool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/Pool"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Pool(ctx, req.(*QueryPoolRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Pools_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPoolsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Pools(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/Pools"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Pools(ctx, req.(*QueryPoolsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Quote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryQuoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Quote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/Quote"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Quote(ctx, req.(*QueryQuoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_SpotPrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySpotPriceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).SpotPrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/SpotPrice"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).SpotPrice(ctx, req.(*QuerySpotPriceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_OptimalRoute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryOptimalRouteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).OptimalRoute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/OptimalRoute"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).OptimalRoute(ctx, req.(*QueryOptimalRouteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Routes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRoutesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Routes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/Routes"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Routes(ctx, req.(*QueryRoutesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Signals_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySignalsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Signals(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/Signals"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Signals(ctx, req.(*QuerySignalsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_SignalHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySignalHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).SignalHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/SignalHistory"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).SignalHistory(ctx, req.(*QuerySignalHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_RiskScore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRiskScoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).RiskScore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/RiskScore"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).RiskScore(ctx, req.(*QueryRiskScoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_TradeRisk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTradeRiskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).TradeRisk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/TradeRisk"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).TradeRisk(ctx, req.(*QueryTradeRiskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_RiskHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRiskHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).RiskHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/RiskHistory"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).RiskHistory(ctx, req.(*QueryRiskHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Sentiment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySentimentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Sentiment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/Sentiment"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Sentiment(ctx, req.(*QuerySentimentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_SentimentHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySentimentHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).SentimentHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/SentimentHistory"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).SentimentHistory(ctx, req.(*QuerySentimentHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_SentimentAlerts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySentimentAlertsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).SentimentAlerts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/SentimentAlerts"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).SentimentAlerts(ctx, req.(*QuerySentimentAlertsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_OracleComposite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryOracleCompositeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).OracleComposite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/OracleComposite"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).OracleComposite(ctx, req.(*QueryOracleCompositeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_WhaleAlerts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryWhaleAlertsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).WhaleAlerts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/WhaleAlerts"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).WhaleAlerts(ctx, req.(*QueryWhaleAlertsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ReferralInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryReferralInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ReferralInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/ReferralInfo"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ReferralInfo(ctx, req.(*QueryReferralInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ReferralList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryReferralListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ReferralList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/ReferralList"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ReferralList(ctx, req.(*QueryReferralListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ReferralCodeLookup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryReferralCodeLookupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ReferralCodeLookup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/ReferralCodeLookup"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ReferralCodeLookup(ctx, req.(*QueryReferralCodeLookupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ReferralEarnings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryReferralEarningsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ReferralEarnings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/ReferralEarnings"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ReferralEarnings(ctx, req.(*QueryReferralEarningsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_OrderBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryOrderBookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).OrderBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/OrderBook"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).OrderBook(ctx, req.(*QueryOrderBookRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_OrderBookDepth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryOrderBookDepthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).OrderBookDepth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/OrderBookDepth"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).OrderBookDepth(ctx, req.(*QueryOrderBookDepthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_UserOrders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryUserOrdersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).UserOrders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/UserOrders"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).UserOrders(ctx, req.(*QueryUserOrdersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_OrderDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryOrderDetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).OrderDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/OrderDetail"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).OrderDetail(ctx, req.(*QueryOrderDetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Trades_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTradesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Trades(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/Trades"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Trades(ctx, req.(*QueryTradesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.dex.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Params", Handler: _Query_Params_Handler},
		{MethodName: "Pool", Handler: _Query_Pool_Handler},
		{MethodName: "Pools", Handler: _Query_Pools_Handler},
		{MethodName: "Quote", Handler: _Query_Quote_Handler},
		{MethodName: "SpotPrice", Handler: _Query_SpotPrice_Handler},
		{MethodName: "OptimalRoute", Handler: _Query_OptimalRoute_Handler},
		{MethodName: "Routes", Handler: _Query_Routes_Handler},
		{MethodName: "Signals", Handler: _Query_Signals_Handler},
		{MethodName: "SignalHistory", Handler: _Query_SignalHistory_Handler},
		{MethodName: "RiskScore", Handler: _Query_RiskScore_Handler},
		{MethodName: "TradeRisk", Handler: _Query_TradeRisk_Handler},
		{MethodName: "RiskHistory", Handler: _Query_RiskHistory_Handler},
		{MethodName: "Sentiment", Handler: _Query_Sentiment_Handler},
		{MethodName: "SentimentHistory", Handler: _Query_SentimentHistory_Handler},
		{MethodName: "SentimentAlerts", Handler: _Query_SentimentAlerts_Handler},
		{MethodName: "OracleComposite", Handler: _Query_OracleComposite_Handler},
		{MethodName: "WhaleAlerts", Handler: _Query_WhaleAlerts_Handler},
		{MethodName: "ReferralInfo", Handler: _Query_ReferralInfo_Handler},
		{MethodName: "ReferralList", Handler: _Query_ReferralList_Handler},
		{MethodName: "ReferralCodeLookup", Handler: _Query_ReferralCodeLookup_Handler},
		{MethodName: "ReferralEarnings", Handler: _Query_ReferralEarnings_Handler},
		{MethodName: "OrderBook", Handler: _Query_OrderBook_Handler},
		{MethodName: "OrderBookDepth", Handler: _Query_OrderBookDepth_Handler},
		{MethodName: "UserOrders", Handler: _Query_UserOrders_Handler},
		{MethodName: "OrderDetail", Handler: _Query_OrderDetail_Handler},
		{MethodName: "Trades", Handler: _Query_Trades_Handler},
		{MethodName: "Leaderboard", Handler: _Query_Leaderboard_Handler},
		{MethodName: "TraderStats", Handler: _Query_TraderStats_Handler},
		{MethodName: "TraderFollowers", Handler: _Query_TraderFollowers_Handler},
		{MethodName: "CopyFollowing", Handler: _Query_CopyFollowing_Handler},
		{MethodName: "CopyHistory", Handler: _Query_CopyHistory_Handler},
		{MethodName: "TimeWeightedPower", Handler: _Query_TimeWeightedPower_Handler},
		{MethodName: "TraderRecord", Handler: _Query_TraderRecord_Handler},
		{MethodName: "TraderLeaderboard", Handler: _Query_TraderLeaderboard_Handler},
		{MethodName: "IBCOrders", Handler: _Query_IBCOrders_Handler},
		{MethodName: "PoolFee", Handler: _Query_PoolFee_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/dex/query.proto",
}

func _Query_PoolFee_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPoolFeeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).PoolFee(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/PoolFee"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).PoolFee(ctx, req.(*QueryPoolFeeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Leaderboard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryLeaderboardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Leaderboard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/Leaderboard"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Leaderboard(ctx, req.(*QueryLeaderboardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_TraderStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTraderStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).TraderStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/TraderStats"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).TraderStats(ctx, req.(*QueryTraderStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_TraderFollowers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTraderFollowersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).TraderFollowers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/TraderFollowers"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).TraderFollowers(ctx, req.(*QueryTraderFollowersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_CopyFollowing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryCopyFollowingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).CopyFollowing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/CopyFollowing"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).CopyFollowing(ctx, req.(*QueryCopyFollowingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_CopyHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryCopyHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).CopyHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/CopyHistory"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).CopyHistory(ctx, req.(*QueryCopyHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_TimeWeightedPower_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTimeWeightedPowerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).TimeWeightedPower(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/TimeWeightedPower"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).TimeWeightedPower(ctx, req.(*QueryTimeWeightedPowerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_TraderRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTraderRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).TraderRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/TraderRecord"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).TraderRecord(ctx, req.(*QueryTraderRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_TraderLeaderboard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTraderLeaderboardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).TraderLeaderboard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/TraderLeaderboard"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).TraderLeaderboard(ctx, req.(*QueryTraderLeaderboardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_IBCOrders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryIBCOrdersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).IBCOrders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.dex.Query/IBCOrders"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).IBCOrders(ctx, req.(*QueryIBCOrdersRequest))
	}
	return interceptor(ctx, in, info, handler)
}
