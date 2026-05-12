package types

import "cosmossdk.io/errors"

var (
	ErrPoolNotFound          = errors.Register(ModuleName, 2, "pool not found")
	ErrPoolAlreadyExists     = errors.Register(ModuleName, 3, "pool already exists for this denom pair")
	ErrInsufficientLiquidity = errors.Register(ModuleName, 4, "insufficient liquidity in pool")
	ErrSlippageExceeded      = errors.Register(ModuleName, 5, "slippage tolerance exceeded")
	ErrInvalidDenom          = errors.Register(ModuleName, 6, "invalid denomination")
	ErrInvalidAmount         = errors.Register(ModuleName, 7, "invalid amount")
	ErrInvalidSender         = errors.Register(ModuleName, 8, "invalid sender address")
	ErrSameDenom             = errors.Register(ModuleName, 9, "cannot create pool with same denomination on both sides")
	ErrZeroLiquidity         = errors.Register(ModuleName, 10, "liquidity amount must be greater than zero")
	ErrMinLiquidityNotMet    = errors.Register(ModuleName, 11, "minimum initial liquidity requirement not met")
	ErrNoRouteFound          = errors.Register(ModuleName, 12, "no route found between the specified denominations")

	// Referral errors
	ErrReferralCodeExists    = errors.Register(ModuleName, 20, "referral code already exists for this address")
	ErrReferralCodeNotFound  = errors.Register(ModuleName, 21, "referral code not found")
	ErrAlreadyReferred       = errors.Register(ModuleName, 22, "user already has a referrer")
	ErrSelfReferral          = errors.Register(ModuleName, 23, "cannot refer yourself")
	ErrCircularReferral      = errors.Register(ModuleName, 24, "circular referral not allowed")

	// Order book errors
	ErrOrderNotFound         = errors.Register(ModuleName, 30, "order not found")
	ErrOrderNotOwned         = errors.Register(ModuleName, 31, "order not owned by sender")
	ErrInvalidOrderSide      = errors.Register(ModuleName, 32, "invalid order side: must be buy or sell")
	ErrInvalidOrderType      = errors.Register(ModuleName, 33, "invalid order type: must be limit or market")
	ErrInvalidOrderPrice     = errors.Register(ModuleName, 34, "invalid order price: must be positive")
	ErrInvalidTimeInForce    = errors.Register(ModuleName, 35, "invalid time in force: must be GTC, IOC, or FOK")
	ErrOrderBookFull         = errors.Register(ModuleName, 36, "order book is full for this side")
	ErrOrderAlreadyCancelled = errors.Register(ModuleName, 37, "order is already cancelled or filled")
	ErrFOKNotFillable        = errors.Register(ModuleName, 38, "fill-or-kill order cannot be fully filled")
	ErrMarketNoLiquidity     = errors.Register(ModuleName, 39, "no liquidity available for market order")
	ErrSelfTrade             = errors.Register(ModuleName, 40, "self-trading is not allowed")
	ErrInvalidTriggerPrice   = errors.Register(ModuleName, 41, "invalid trigger price: must be positive for conditional orders")
	ErrTriggerPriceNotAllowed = errors.Register(ModuleName, 42, "trigger price not allowed for regular limit/market orders")

	// AMM protection errors
	ErrSwapTooLarge = errors.Register(ModuleName, 50, "swap input exceeds 30% of pool reserve")
	ErrPoolDrained  = errors.Register(ModuleName, 51, "pool reserves fully drained; must recreate pool")

	// Multi-hop errors
	ErrTooManyHops = errors.Register(ModuleName, 60, "route exceeds maximum of 4 hops")

	// Dynamic fee errors
	ErrNotPoolCreator   = errors.Register(ModuleName, 70, "only the pool creator can configure fees")
	ErrInvalidFeeConfig = errors.Register(ModuleName, 71, "invalid fee configuration: max_fee must be >= base_fee")
)
