package types

import (
	"encoding/binary"
	"fmt"
	"time"

	"cosmossdk.io/math"
)

// ---------------------------------------------------------------------------
// Order Book Constants
// ---------------------------------------------------------------------------

const (
	// Order sides
	OrderSideBuy  = "buy"
	OrderSideSell = "sell"

	// Order types
	OrderTypeLimit          = "limit"
	OrderTypeMarket         = "market"
	OrderTypeStopLoss       = "stop_loss"
	OrderTypeTakeProfit     = "take_profit"
	OrderTypeStopLossLimit  = "stop_loss_limit"
	OrderTypeTakeProfitLimit = "take_profit_limit"

	// Time in force
	TimeInForceGTC = "GTC" // Good-Till-Cancelled
	TimeInForceIOC = "IOC" // Immediate-Or-Cancel
	TimeInForceFOK = "FOK" // Fill-Or-Kill

	// Order statuses
	OrderStatusOpen      = "open"
	OrderStatusFilled    = "filled"
	OrderStatusPartial   = "partial"
	OrderStatusCancelled = "cancelled"
	OrderStatusExpired   = "expired"
	OrderStatusErrored   = "errored"    // H3: refund/transfer failure during fill
	OrderStatusTriggered = "triggered"  // conditional order has been triggered
	OrderStatusPending   = "pending"    // conditional order waiting for trigger

	// Store prefixes
	OrderPrefix              = "order/"
	OrderBookBidPrefix       = "ob_bid/"       // orderbook bids: ob_bid/<poolID>/<price_inverted>/<timestamp>/<orderID>
	OrderBookAskPrefix       = "ob_ask/"       // orderbook asks: ob_ask/<poolID>/<price>/<timestamp>/<orderID>
	OrdersByAddressPrefix    = "ob_addr/"       // user's orders: ob_addr/<address>/<orderID>
	TradeHistoryPrefix       = "ob_trade/"      // trade history: ob_trade/<poolID>/<tradeID>
	PoolTradeCountPrefix     = "pool_trade_count/" // per-pool trade count: pool_trade_count/<poolID>
	ConditionalOrderPrefix   = "ob_cond/"       // conditional orders: ob_cond/<poolID>/<orderID>
	NextOrderIDKey           = "next_order_id"
	NextTradeIDKey           = "next_trade_id"

	// Limits
	MaxOrdersPerSide      = 500
	DefaultOrderExpiry     = int64(100000) // blocks
	MaxTradesPerPool       = 1000
	PricePrecisionDecimals = 18 // precision for price encoding in keys
)

// ---------------------------------------------------------------------------
// Order represents a limit or market order on the order book
// ---------------------------------------------------------------------------

type Order struct {
	ID            uint64         `json:"id"`
	PoolID        uint64         `json:"pool_id"`
	Creator       string         `json:"creator"`
	Side          string         `json:"side"`           // "buy" or "sell"
	OrderType     string         `json:"order_type"`     // "limit", "market", "stop_loss", "take_profit", "stop_loss_limit", "take_profit_limit"
	Price         math.LegacyDec `json:"price"`          // price in quote per base (DenomB per DenomA)
	Quantity      math.Int       `json:"quantity"`        // total quantity in base denom (DenomA)
	FilledQty     math.Int       `json:"filled_qty"`     // how much has been filled
	TimeInForce   string         `json:"time_in_force"`  // GTC, IOC, FOK
	Status        string         `json:"status"`          // open, filled, partial, cancelled, expired, pending, triggered
	CreatedAt     int64          `json:"created_at"`      // block height
	ExpiresAt     int64          `json:"expires_at"`      // block height when order expires
	LastUpdatedAt int64          `json:"last_updated_at"` // block height of last fill
	TriggerPrice  math.LegacyDec `json:"trigger_price"`  // trigger price for stop-loss/take-profit orders
}

// RemainingQty returns the unfilled quantity.
func (o *Order) RemainingQty() math.Int {
	return o.Quantity.Sub(o.FilledQty)
}

// IsFilled returns true if the order is fully filled.
func (o *Order) IsFilled() bool {
	return o.FilledQty.GTE(o.Quantity)
}

// ---------------------------------------------------------------------------
// Trade represents a matched trade between two orders or an order and the AMM
// ---------------------------------------------------------------------------

type Trade struct {
	ID         uint64         `json:"id"`
	PoolID     uint64         `json:"pool_id"`
	Price      math.LegacyDec `json:"price"`
	Quantity   math.Int       `json:"quantity"` // amount in base denom
	MakerAddr  string         `json:"maker_addr"`
	TakerAddr  string         `json:"taker_addr"`
	MakerOrder uint64         `json:"maker_order_id"`
	TakerOrder uint64         `json:"taker_order_id"`
	Side       string         `json:"side"`       // taker's side
	Source     string         `json:"source"`      // "orderbook" or "amm" or "hybrid"
	Timestamp  time.Time      `json:"timestamp"`
	BlockHeight int64         `json:"block_height"`
}

// ---------------------------------------------------------------------------
// OrderBookSummary is the response for the order book query
// ---------------------------------------------------------------------------

type OrderBookSummary struct {
	PoolID uint64           `json:"pool_id"`
	Bids   []OrderBookLevel `json:"bids"`
	Asks   []OrderBookLevel `json:"asks"`
}

type OrderBookLevel struct {
	Price    math.LegacyDec `json:"price"`
	Quantity math.Int       `json:"quantity"`
	Orders   int            `json:"orders"`
}

// ---------------------------------------------------------------------------
// KV Store Key Functions
// ---------------------------------------------------------------------------

// OrderKey returns the store key for an order by its ID.
func OrderKey(orderID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, orderID)
	return append([]byte(OrderPrefix), bz...)
}

// OrderBookBidKey returns the key for a bid order in price-time priority.
// Bids are sorted by price DESCENDING (highest first), so we invert the price.
// Key format: ob_bid/<poolID_8bytes><inverted_price_32bytes><height_8bytes><orderID_8bytes>
func OrderBookBidKey(poolID uint64, price math.LegacyDec, height int64, orderID uint64) []byte {
	poolBz := make([]byte, 8)
	binary.BigEndian.PutUint64(poolBz, poolID)

	priceBz := encodePriceInverted(price)

	heightBz := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBz, uint64(height))

	orderBz := make([]byte, 8)
	binary.BigEndian.PutUint64(orderBz, orderID)

	key := append([]byte(OrderBookBidPrefix), poolBz...)
	key = append(key, priceBz...)
	key = append(key, heightBz...)
	key = append(key, orderBz...)
	return key
}

// OrderBookAskKey returns the key for an ask order in price-time priority.
// Asks are sorted by price ASCENDING (lowest first).
// Key format: ob_ask/<poolID_8bytes><price_32bytes><height_8bytes><orderID_8bytes>
func OrderBookAskKey(poolID uint64, price math.LegacyDec, height int64, orderID uint64) []byte {
	poolBz := make([]byte, 8)
	binary.BigEndian.PutUint64(poolBz, poolID)

	priceBz := encodePrice(price)

	heightBz := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBz, uint64(height))

	orderBz := make([]byte, 8)
	binary.BigEndian.PutUint64(orderBz, orderID)

	key := append([]byte(OrderBookAskPrefix), poolBz...)
	key = append(key, priceBz...)
	key = append(key, heightBz...)
	key = append(key, orderBz...)
	return key
}

// OrderBookBidPoolPrefix returns the prefix for all bids for a pool.
func OrderBookBidPoolPrefix(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(OrderBookBidPrefix), bz...)
}

// OrderBookAskPoolPrefix returns the prefix for all asks for a pool.
func OrderBookAskPoolPrefix(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(OrderBookAskPrefix), bz...)
}

// OrdersByAddressKey returns the key for an order indexed by user address.
func OrdersByAddressKey(address string, orderID uint64) []byte {
	orderBz := make([]byte, 8)
	binary.BigEndian.PutUint64(orderBz, orderID)
	key := append([]byte(OrdersByAddressPrefix), []byte(address+"/")...)
	return append(key, orderBz...)
}

// OrdersByAddressAddrPrefix returns the prefix for all orders for an address.
func OrdersByAddressAddrPrefix(address string) []byte {
	return append([]byte(OrdersByAddressPrefix), []byte(address+"/")...)
}

// TradeHistoryKey returns the key for a trade.
func TradeHistoryKey(poolID uint64, tradeID uint64) []byte {
	poolBz := make([]byte, 8)
	binary.BigEndian.PutUint64(poolBz, poolID)
	tradeBz := make([]byte, 8)
	binary.BigEndian.PutUint64(tradeBz, tradeID)
	key := append([]byte(TradeHistoryPrefix), poolBz...)
	return append(key, tradeBz...)
}

// TradeHistoryPoolPrefix returns the prefix for all trades for a pool.
func TradeHistoryPoolPrefix(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(TradeHistoryPrefix), bz...)
}

// PoolTradeCountKey returns the key for the per-pool trade count counter.
func PoolTradeCountKey(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(PoolTradeCountPrefix), bz...)
}

// ConditionalOrderKey returns the key for a conditional order indexed by pool.
// Key format: ob_cond/<poolID_8bytes><orderID_8bytes>
func ConditionalOrderKey(poolID uint64, orderID uint64) []byte {
	poolBz := make([]byte, 8)
	binary.BigEndian.PutUint64(poolBz, poolID)
	orderBz := make([]byte, 8)
	binary.BigEndian.PutUint64(orderBz, orderID)
	key := append([]byte(ConditionalOrderPrefix), poolBz...)
	return append(key, orderBz...)
}

// ConditionalOrderPoolPrefix returns the prefix for all conditional orders for a pool.
func ConditionalOrderPoolPrefix(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(ConditionalOrderPrefix), bz...)
}

// IsConditionalOrderType returns true if the order type is a conditional type.
func IsConditionalOrderType(orderType string) bool {
	return orderType == OrderTypeStopLoss || orderType == OrderTypeTakeProfit ||
		orderType == OrderTypeStopLossLimit || orderType == OrderTypeTakeProfitLimit
}

// ---------------------------------------------------------------------------
// Price encoding: encode price as a 32-byte sortable representation
// We use a fixed-point encoding: integer part (14 bytes) + decimal part (18 bytes)
// ---------------------------------------------------------------------------

func encodePrice(price math.LegacyDec) []byte {
	// price is stored as an Int with 18 decimals of precision
	// e.g., price 1.5 is stored as 1500000000000000000
	raw := price.MulInt64(1e18).TruncateInt()

	// Encode as 32-byte big-endian (positive numbers only for prices)
	bz := make([]byte, 32)
	if raw.IsNegative() {
		return bz // zero-filled for negative prices (shouldn't happen)
	}

	// Convert to big.Int bytes and right-align in 32 bytes
	bigBytes := raw.BigInt().Bytes()
	if len(bigBytes) > 32 {
		bigBytes = bigBytes[:32]
	}
	copy(bz[32-len(bigBytes):], bigBytes)
	return bz
}

func encodePriceInverted(price math.LegacyDec) []byte {
	bz := encodePrice(price)
	// Invert all bytes for descending sort
	for i := range bz {
		bz[i] = ^bz[i]
	}
	return bz
}

// suppress unused variable warnings
var _ = fmt.Sprintf
