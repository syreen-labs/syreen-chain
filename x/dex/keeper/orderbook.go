package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/dex/types"
)

// ---------------------------------------------------------------------------
// Order ID management
// ---------------------------------------------------------------------------

func (k Keeper) GetNextOrderID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.NextOrderIDKey))
	if err != nil || bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetNextOrderID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	kvStore.Set([]byte(types.NextOrderIDKey), bz)
}

func (k Keeper) GetNextTradeID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.NextTradeIDKey))
	if err != nil || bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetNextTradeID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	kvStore.Set([]byte(types.NextTradeIDKey), bz)
}

// ---------------------------------------------------------------------------
// Order CRUD
// ---------------------------------------------------------------------------

func (k Keeper) GetOrder(ctx context.Context, orderID uint64) (types.Order, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.OrderKey(orderID))
	if err != nil || bz == nil {
		return types.Order{}, false
	}
	var order types.Order
	if err := json.Unmarshal(bz, &order); err != nil {
		return types.Order{}, false
	}
	return order, true
}

func (k Keeper) SetOrder(ctx context.Context, order types.Order) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(order)
	kvStore.Set(types.OrderKey(order.ID), bz)
}

func (k Keeper) DeleteOrder(ctx context.Context, order types.Order) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Delete(types.OrderKey(order.ID))
}

// ---------------------------------------------------------------------------
// Order Book Index Management
// ---------------------------------------------------------------------------

func (k Keeper) addOrderToBook(ctx context.Context, order types.Order) {
	kvStore := k.storeService.OpenKVStore(ctx)

	// Store orderID as value
	orderIDBz := make([]byte, 8)
	binary.BigEndian.PutUint64(orderIDBz, order.ID)

	if order.Side == types.OrderSideBuy {
		key := types.OrderBookBidKey(order.PoolID, order.Price, order.CreatedAt, order.ID)
		kvStore.Set(key, orderIDBz)
	} else {
		key := types.OrderBookAskKey(order.PoolID, order.Price, order.CreatedAt, order.ID)
		kvStore.Set(key, orderIDBz)
	}

	// Add to address index
	kvStore.Set(types.OrdersByAddressKey(order.Creator, order.ID), orderIDBz)
}

func (k Keeper) removeOrderFromBook(ctx context.Context, order types.Order) {
	kvStore := k.storeService.OpenKVStore(ctx)

	if order.Side == types.OrderSideBuy {
		key := types.OrderBookBidKey(order.PoolID, order.Price, order.CreatedAt, order.ID)
		kvStore.Delete(key)
	} else {
		key := types.OrderBookAskKey(order.PoolID, order.Price, order.CreatedAt, order.ID)
		kvStore.Delete(key)
	}

	// Remove from address index
	kvStore.Delete(types.OrdersByAddressKey(order.Creator, order.ID))
}

// ---------------------------------------------------------------------------
// Count orders per side per pool
// ---------------------------------------------------------------------------

func (k Keeper) countOrdersOnSide(ctx context.Context, poolID uint64, side string) int {
	kvStore := k.storeService.OpenKVStore(ctx)

	var prefix []byte
	if side == types.OrderSideBuy {
		prefix = types.OrderBookBidPoolPrefix(poolID)
	} else {
		prefix = types.OrderBookAskPoolPrefix(poolID)
	}

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return 0
	}
	defer iter.Close()

	count := 0
	for ; iter.Valid(); iter.Next() {
		count++
	}
	return count
}

// ---------------------------------------------------------------------------
// Get orders for a side (sorted by price-time priority)
// ---------------------------------------------------------------------------

func (k Keeper) getBidOrders(ctx context.Context, poolID uint64, limit int) []types.Order {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.OrderBookBidPoolPrefix(poolID)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var orders []types.Order
	for ; iter.Valid() && len(orders) < limit; iter.Next() {
		orderID := binary.BigEndian.Uint64(iter.Value())
		order, found := k.GetOrder(ctx, orderID)
		if found && order.Status == types.OrderStatusOpen {
			orders = append(orders, order)
		}
	}
	return orders
}

func (k Keeper) getAskOrders(ctx context.Context, poolID uint64, limit int) []types.Order {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.OrderBookAskPoolPrefix(poolID)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var orders []types.Order
	for ; iter.Valid() && len(orders) < limit; iter.Next() {
		orderID := binary.BigEndian.Uint64(iter.Value())
		order, found := k.GetOrder(ctx, orderID)
		if found && order.Status == types.OrderStatusOpen {
			orders = append(orders, order)
		}
	}
	return orders
}

// ---------------------------------------------------------------------------
// PlaceOrder: main entry point for placing orders
// ---------------------------------------------------------------------------

func (k Keeper) PlaceOrder(ctx context.Context, msg *types.MsgPlaceOrder) (uint64, string, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockHeight := sdkCtx.BlockHeight()

	pool, found := k.GetPool(ctx, msg.PoolID)
	if !found {
		return 0, "", types.ErrPoolNotFound
	}

	creatorAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return 0, "", types.ErrInvalidSender
	}

	isConditional := types.IsConditionalOrderType(msg.OrderType)

	// Parse price from string
	var parsedPrice math.LegacyDec
	if msg.Price != "" && msg.Price != "0" {
		var parseErr error
		parsedPrice, parseErr = math.LegacyNewDecFromStr(msg.Price)
		if parseErr != nil {
			return 0, "", fmt.Errorf("invalid price: %w", parseErr)
		}
	} else {
		parsedPrice = math.LegacyZeroDec()
	}

	// Parse trigger price from string
	// Parse Quantity
	parsedQuantity, qtyOk := math.NewIntFromString(msg.Quantity)
	if !qtyOk {
		return 0, "", fmt.Errorf("invalid quantity: %s", msg.Quantity)
	}

	var parsedTriggerPrice math.LegacyDec
	if msg.TriggerPrice != "" && msg.TriggerPrice != "0" {
		var parseErr error
		parsedTriggerPrice, parseErr = math.LegacyNewDecFromStr(msg.TriggerPrice)
		if parseErr != nil {
			return 0, "", fmt.Errorf("invalid trigger price: %w", parseErr)
		}
	} else {
		parsedTriggerPrice = math.LegacyZeroDec()
	}

	// Determine the price for market orders (use AMM spot price)
	price := parsedPrice
	if msg.OrderType == types.OrderTypeMarket {
		ammPrice, err := k.getAMMPrice(ctx, pool, msg.Side)
		if err != nil {
			return 0, "", types.ErrMarketNoLiquidity
		}
		price = ammPrice
	}

	// For conditional market-style orders (stop_loss, take_profit), use trigger price as escrow price
	// with a 10% buffer for buy orders to cover market price movement at execution time
	if msg.OrderType == types.OrderTypeStopLoss || msg.OrderType == types.OrderTypeTakeProfit {
		price = parsedTriggerPrice
		if msg.Side == types.OrderSideBuy {
			// Add 10% buffer: escrow = triggerPrice * 1.1
			price = parsedTriggerPrice.Mul(math.LegacyNewDecWithPrec(110, 2))
		}
	}
	// For conditional limit-style orders (stop_loss_limit, take_profit_limit), Price is already set from msg.Price

	// Check order book capacity for limit orders that will rest on the book
	if msg.OrderType == types.OrderTypeLimit && msg.TimeInForce == types.TimeInForceGTC {
		count := k.countOrdersOnSide(ctx, msg.PoolID, msg.Side)
		if count >= types.MaxOrdersPerSide {
			return 0, "", types.ErrOrderBookFull
		}
	}

	// Calculate escrow amount
	escrowCoin := k.calculateEscrow(pool, msg.Side, price, parsedQuantity)

	// Escrow funds from creator to module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, sdk.NewCoins(escrowCoin)); err != nil {
		return 0, "", err
	}

	// Create order
	orderID := k.GetNextOrderID(ctx)
	k.SetNextOrderID(ctx, orderID+1)

	// Set trigger price (zero for non-conditional orders)
	triggerPrice := math.LegacyZeroDec()
	if isConditional {
		triggerPrice = parsedTriggerPrice
	}

	order := types.Order{
		ID:            orderID,
		PoolID:        msg.PoolID,
		Creator:       msg.Creator,
		Side:          msg.Side,
		OrderType:     msg.OrderType,
		Price:         price,
		Quantity:      parsedQuantity,
		FilledQty:     math.ZeroInt(),
		TimeInForce:   msg.TimeInForce,
		Status:        types.OrderStatusOpen,
		CreatedAt:     blockHeight,
		ExpiresAt:     blockHeight + types.DefaultOrderExpiry,
		LastUpdatedAt: blockHeight,
		TriggerPrice:  triggerPrice,
	}

	// Conditional orders go to pending state — they sit in the conditional
	// index until the trigger price is hit.
	if isConditional {
		order.Status = types.OrderStatusPending
		k.SetOrder(ctx, order)
		k.addConditionalOrder(ctx, order)
		// Also add to address index for user queries
		kvStore := k.storeService.OpenKVStore(ctx)
		orderIDBz := make([]byte, 8)
		binary.BigEndian.PutUint64(orderIDBz, order.ID)
		kvStore.Set(types.OrdersByAddressKey(order.Creator, order.ID), orderIDBz)

		k.Logger(ctx).Info("conditional order placed",
			"order_id", orderID,
			"pool_id", msg.PoolID,
			"side", msg.Side,
			"type", msg.OrderType,
			"trigger_price", triggerPrice,
			"price", price,
			"quantity", msg.Quantity,
			"status", order.Status,
		)
		return orderID, order.Status, nil
	}

	// For IOC and FOK orders, try to match immediately
	if msg.TimeInForce == types.TimeInForceIOC || msg.TimeInForce == types.TimeInForceFOK || msg.OrderType == types.OrderTypeMarket {
		filledQty := k.tryImmediateMatch(ctx, &order, pool)

		if msg.TimeInForce == types.TimeInForceFOK && filledQty.LT(order.Quantity) {
			// FOK: must fill entirely, refund (H3: propagate refund failures).
			if err := k.refundEscrow(ctx, creatorAddr, escrowCoin); err != nil {
				return 0, "", fmt.Errorf("FOK refund failed: %w", err)
			}
			return 0, "", types.ErrFOKNotFillable
		}

		order.FilledQty = filledQty
		if order.IsFilled() {
			order.Status = types.OrderStatusFilled
		} else if filledQty.IsPositive() {
			order.Status = types.OrderStatusPartial
		}

		// IOC and market orders: cancel any unfilled portion
		if (msg.TimeInForce == types.TimeInForceIOC || msg.OrderType == types.OrderTypeMarket) && !order.IsFilled() {
			// Refund unfilled portion (H3: propagate refund failures so the
			// order is not silently marked done with funds stuck in the module).
			if filledQty.IsPositive() {
				unfilledQty := order.Quantity.Sub(filledQty)
				refundCoin := k.calculateEscrow(pool, msg.Side, price, unfilledQty)
				if err := k.refundEscrow(ctx, creatorAddr, refundCoin); err != nil {
					return 0, "", fmt.Errorf("IOC unfilled refund failed: %w", err)
				}
			} else {
				if err := k.refundEscrow(ctx, creatorAddr, escrowCoin); err != nil {
					return 0, "", fmt.Errorf("IOC refund failed: %w", err)
				}
			}
			if !order.IsFilled() {
				order.Status = types.OrderStatusCancelled
			}
		}
	}

	// Save order
	k.SetOrder(ctx, order)

	// If order is still open (GTC limit orders that weren't fully filled), add to book
	if order.Status == types.OrderStatusOpen {
		k.addOrderToBook(ctx, order)
	}

	k.Logger(ctx).Info("order placed",
		"order_id", orderID,
		"pool_id", msg.PoolID,
		"side", msg.Side,
		"type", msg.OrderType,
		"price", price,
		"quantity", msg.Quantity,
		"status", order.Status,
	)

	return orderID, order.Status, nil
}

// ---------------------------------------------------------------------------
// CancelOrder
// ---------------------------------------------------------------------------

func (k Keeper) CancelOrder(ctx context.Context, creator string, orderID uint64) error {
	order, found := k.GetOrder(ctx, orderID)
	if !found {
		return types.ErrOrderNotFound
	}

	if order.Creator != creator {
		return types.ErrOrderNotOwned
	}

	if order.Status != types.OrderStatusOpen && order.Status != types.OrderStatusPartial && order.Status != types.OrderStatusPending {
		return types.ErrOrderAlreadyCancelled
	}

	pool, found := k.GetPool(ctx, order.PoolID)
	if !found {
		return types.ErrPoolNotFound
	}

	creatorAddr, _ := sdk.AccAddressFromBech32(creator)

	// Refund unfilled portion (H3: must succeed before we mark cancelled).
	unfilledQty := order.RemainingQty()
	if unfilledQty.IsPositive() {
		refundCoin := k.calculateEscrow(pool, order.Side, order.Price, unfilledQty)
		if err := k.refundEscrow(ctx, creatorAddr, refundCoin); err != nil {
			return fmt.Errorf("cancel refund failed: %w", err)
		}
	}

	// Remove from book or conditional index and update status
	if order.Status == types.OrderStatusPending {
		k.removeConditionalOrder(ctx, order)
	} else {
		k.removeOrderFromBook(ctx, order)
	}
	order.Status = types.OrderStatusCancelled
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	order.LastUpdatedAt = sdkCtx.BlockHeight()
	k.SetOrder(ctx, order)

	k.Logger(ctx).Info("order cancelled", "order_id", orderID, "creator", creator)
	return nil
}

// ---------------------------------------------------------------------------
// ModifyOrder: cancel + replace
// ---------------------------------------------------------------------------

func (k Keeper) ModifyOrder(ctx context.Context, creator string, orderID uint64, newPriceStr string, newQuantity math.Int) (uint64, error) {
	order, found := k.GetOrder(ctx, orderID)
	if !found {
		return 0, types.ErrOrderNotFound
	}

	if order.Creator != creator {
		return 0, types.ErrOrderNotOwned
	}

	if order.Status != types.OrderStatusOpen && order.Status != types.OrderStatusPartial {
		return 0, types.ErrOrderAlreadyCancelled
	}

	// Cancel existing order
	if err := k.CancelOrder(ctx, creator, orderID); err != nil {
		return 0, err
	}

	// Determine new values
	price := order.Price
	if newPriceStr != "" && newPriceStr != "0" {
		parsedNewPrice, parseErr := math.LegacyNewDecFromStr(newPriceStr)
		if parseErr != nil {
			return 0, fmt.Errorf("invalid new price: %w", parseErr)
		}
		if parsedNewPrice.IsPositive() {
			price = parsedNewPrice
		}
	}
	quantity := order.RemainingQty()
	if !newQuantity.IsNil() && newQuantity.IsPositive() {
		quantity = newQuantity
	}

	// Place new order
	newMsg := &types.MsgPlaceOrder{
		Creator:      creator,
		PoolID:       order.PoolID,
		Side:         order.Side,
		OrderType:    order.OrderType,
		Price:        price.String(),
		Quantity:     quantity.String(),
		TimeInForce:  order.TimeInForce,
		TriggerPrice: order.TriggerPrice.String(),
	}

	newID, _, err := k.PlaceOrder(ctx, newMsg)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

// ---------------------------------------------------------------------------
// Matching Engine: runs in BeginBlock
// ---------------------------------------------------------------------------

func (k Keeper) RunOrderBookMatching(ctx sdk.Context) {
	pools := k.GetAllPools(ctx)
	for _, pool := range pools {
		k.matchOrdersForPool(ctx, pool)
		k.cleanupExpiredOrders(ctx, pool.ID)
	}
}

func (k Keeper) matchOrdersForPool(ctx context.Context, pool types.Pool) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockHeight := sdkCtx.BlockHeight()
	blockTime := sdkCtx.BlockTime()

	// Get best bid and ask orders
	bids := k.getBidOrders(ctx, pool.ID, types.MaxOrdersPerSide)
	asks := k.getAskOrders(ctx, pool.ID, types.MaxOrdersPerSide)

	bidIdx := 0
	askIdx := 0

	for bidIdx < len(bids) && askIdx < len(asks) {
		bid := &bids[bidIdx]
		ask := &asks[askIdx]

		// Check if bid price >= ask price (orders can match)
		if bid.Price.LT(ask.Price) {
			break // No more matches possible
		}

		// Match at the maker's price (the order that was placed first)
		matchPrice := ask.Price
		if bid.CreatedAt < ask.CreatedAt {
			matchPrice = bid.Price
		}

		// Determine match quantity
		bidRemaining := bid.RemainingQty()
		askRemaining := ask.RemainingQty()
		matchQty := bidRemaining
		if askRemaining.LT(matchQty) {
			matchQty = askRemaining
		}

		if matchQty.IsZero() {
			break
		}

		// Execute the match
		if err := k.executeMatch(ctx, bid, ask, matchPrice, matchQty, pool, blockHeight, blockTime, "orderbook"); err != nil {
			// D-03: Don't update order status on failed match; log and skip
			if errors.Is(err, types.ErrSelfTrade) {
				// D-16: Skip self-trade, advance both indices
				bidIdx++
				continue
			}
			k.Logger(ctx).Error("executeMatch failed in BeginBlock", "error", err, "buy_order", bid.ID, "sell_order", ask.ID)
			// Skip this pair — try next ask for this bid
			askIdx++
			continue
		}

		// Update bid
		bid.FilledQty = bid.FilledQty.Add(matchQty)
		bid.LastUpdatedAt = blockHeight
		if bid.IsFilled() {
			bid.Status = types.OrderStatusFilled
			k.removeOrderFromBook(ctx, *bid)
			bidIdx++
		} else {
			bid.Status = types.OrderStatusPartial
		}
		k.SetOrder(ctx, *bid)

		// Update ask
		ask.FilledQty = ask.FilledQty.Add(matchQty)
		ask.LastUpdatedAt = blockHeight
		if ask.IsFilled() {
			ask.Status = types.OrderStatusFilled
			k.removeOrderFromBook(ctx, *ask)
			askIdx++
		} else {
			ask.Status = types.OrderStatusPartial
		}
		k.SetOrder(ctx, *ask)
	}

	// After order book matching, check remaining orders against AMM for better prices
	k.matchRemainingAgainstAMM(ctx, pool)
}

// matchRemainingAgainstAMM checks if any remaining open orders can get a better deal from the AMM.
func (k Keeper) matchRemainingAgainstAMM(ctx context.Context, pool types.Pool) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockHeight := sdkCtx.BlockHeight()
	blockTime := sdkCtx.BlockTime()

	// C1-FIX: Per-block volume cap — max 30% of each reserve consumed via orderbook AMM fills
	maxDenomBConsumed := pool.ReserveB.MulRaw(30).QuoRaw(100)
	maxDenomAConsumed := pool.ReserveA.MulRaw(30).QuoRaw(100)
	totalDenomBConsumed := math.ZeroInt()
	totalDenomAConsumed := math.ZeroInt()

	// Check buy orders against AMM
	bids := k.getBidOrders(ctx, pool.ID, 50) // check top 50
	for _, bid := range bids {
		remaining := bid.RemainingQty()
		if remaining.IsZero() {
			continue
		}

		// C1-FIX: Stop if we've consumed 30% of reserves this block
		if totalDenomBConsumed.GTE(maxDenomBConsumed) {
			break
		}

		// For a buy order: user wants to buy DenomA, paying DenomB
		// AMM price: how much DenomB per DenomA
		ammPrice, err := k.getAMMPrice(ctx, pool, types.OrderSideBuy)
		if err != nil {
			continue
		}

		// If AMM price is better (lower) than the limit price, fill via AMM
		if ammPrice.LTE(bid.Price) {
			// H2-FIX: Cap tokenIn at actual escrowed amount for this order
			escrowRemaining := bid.Price.MulInt(remaining).Ceil().TruncateInt()
			tokenIn := sdk.NewCoin(pool.DenomB, escrowRemaining)

			// C1-FIX: Further cap tokenIn to stay within per-block limit
			budgetLeft := maxDenomBConsumed.Sub(totalDenomBConsumed)
			if tokenIn.Amount.GT(budgetLeft) {
				tokenIn = sdk.NewCoin(pool.DenomB, budgetLeft)
			}

			quote, _, qErr := k.GetQuote(ctx, pool.ID, tokenIn)
			if qErr != nil || quote.Amount.IsZero() {
				continue
			}

			// Only fill what we can get from AMM
			fillQty := quote.Amount
			if fillQty.GT(remaining) {
				fillQty = remaining
			}

			if fillQty.IsZero() {
				continue
			}

			// Execute the AMM fill
			creatorAddr, _ := sdk.AccAddressFromBech32(bid.Creator)

			// H2-FIX: actualTokenIn capped at escrow remaining
			actualTokenIn := sdk.NewCoin(pool.DenomB, ammPrice.MulInt(fillQty).Ceil().TruncateInt())
			if actualTokenIn.Amount.GT(escrowRemaining) {
				actualTokenIn = sdk.NewCoin(pool.DenomB, escrowRemaining)
			}

			tokenOut, swapErr := k.executeAMMSwapInternal(ctx, pool.ID, actualTokenIn)
			if swapErr != nil {
				continue
			}

			// C1-FIX: Track cumulative volume
			totalDenomBConsumed = totalDenomBConsumed.Add(actualTokenIn.Amount)

			actualFillQty := tokenOut.Amount
			if actualFillQty.GT(remaining) {
				// H5-FIX: Return excess to pool reserves, not to order creator
				excess := actualFillQty.Sub(remaining)
				if excess.IsPositive() {
					// Add excess back to pool reserves instead of giving free tokens
					refreshedPool, ok := k.GetPool(ctx, pool.ID)
					if ok {
						refreshedPool.ReserveA = refreshedPool.ReserveA.Add(excess)
						k.SetPool(ctx, refreshedPool)
						pool = refreshedPool
					}
				}
				actualFillQty = remaining
			}

			// Send the output tokens to the order creator
			if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, sdk.NewCoins(sdk.NewCoin(pool.DenomA, actualFillQty))); err != nil {
				continue
			}

			// Record trade
			k.recordTrade(ctx, pool.ID, ammPrice, actualFillQty, bid.Creator, "AMM", bid.ID, 0, types.OrderSideBuy, "amm", blockHeight, blockTime)

			// BUY-FIX: refund price-improvement escrow on the filled portion.
			// The order escrowed bid.Price per unit but the AMM fill only spent
			// actualTokenIn (at the lower ammPrice). Release the difference so it
			// isn't stranded when the order is later removed from the book. Mirrors
			// the (escrow - payment) accounting in executeMatch. The unfilled
			// portion keeps its escrow, so partial fills are not over-refunded.
			escrowForFill := bid.Price.MulInt(actualFillQty).Ceil().TruncateInt()
			if priceImprovement := escrowForFill.Sub(actualTokenIn.Amount); priceImprovement.IsPositive() {
				_ = k.refundEscrow(ctx, creatorAddr, sdk.NewCoin(pool.DenomB, priceImprovement))
			}

			// Update order
			bid.FilledQty = bid.FilledQty.Add(actualFillQty)
			bid.LastUpdatedAt = blockHeight
			if bid.IsFilled() {
				bid.Status = types.OrderStatusFilled
				k.removeOrderFromBook(ctx, bid)
			} else {
				bid.Status = types.OrderStatusPartial
			}
			k.SetOrder(ctx, bid)

			// Refresh pool. H4: if the pool is gone (deleted, errored read),
			// stop matching this side rather than continuing with stale data.
			refreshed, ok := k.GetPool(ctx, pool.ID)
			if !ok {
				k.Logger(ctx).Error("matchRemainingAgainstAMM: pool refresh failed (bid loop), aborting matching for pool",
					"pool_id", pool.ID)
				return
			}
			pool = refreshed
		}
	}

	// Check sell orders against AMM
	asks := k.getAskOrders(ctx, pool.ID, 50)
	for _, ask := range asks {
		remaining := ask.RemainingQty()
		if remaining.IsZero() {
			continue
		}

		// C1-FIX: Stop if we've consumed 30% of DenomA reserves this block
		if totalDenomAConsumed.GTE(maxDenomAConsumed) {
			break
		}

		// For a sell order: user wants to sell DenomA for DenomB
		ammPrice, err := k.getAMMPrice(ctx, pool, types.OrderSideSell)
		if err != nil {
			continue
		}

		// If AMM price is better (higher) than the limit price, fill via AMM
		if ammPrice.GTE(ask.Price) {
			// C1-FIX: Cap sell amount to stay within per-block limit
			sellAmount := remaining
			budgetLeft := maxDenomAConsumed.Sub(totalDenomAConsumed)
			budgetCapped := false
			if sellAmount.GT(budgetLeft) {
				sellAmount = budgetLeft
				budgetCapped = true
			}
			if sellAmount.IsZero() {
				break
			}
			tokenIn := sdk.NewCoin(pool.DenomA, sellAmount)
			quote, _, qErr := k.GetQuote(ctx, pool.ID, tokenIn)
			if qErr != nil || quote.Amount.IsZero() {
				continue
			}

			// Execute AMM swap from escrowed funds
			tokenOut, swapErr := k.executeAMMSwapInternal(ctx, pool.ID, tokenIn)
			if swapErr != nil {
				continue
			}

			// C1-FIX: Track cumulative volume
			totalDenomAConsumed = totalDenomAConsumed.Add(tokenIn.Amount)

			creatorAddr, _ := sdk.AccAddressFromBech32(ask.Creator)

			// Send the output tokens (DenomB) to the order creator
			if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, sdk.NewCoins(tokenOut)); err != nil {
				continue
			}

			// SELL-FIX: record and advance FilledQty by the ACTUAL amount swapped
			// (sellAmount), not the uncapped remaining. Using remaining here would
			// strand (remaining - sellAmount) of the seller's DenomA escrow and
			// could push FilledQty above Quantity when the budget cap fires.
			k.recordTrade(ctx, pool.ID, ammPrice, sellAmount, ask.Creator, "AMM", ask.ID, 0, types.OrderSideSell, "amm", blockHeight, blockTime)

			// Update order by the actual amount filled
			ask.FilledQty = ask.FilledQty.Add(sellAmount)
			ask.LastUpdatedAt = blockHeight
			if ask.IsFilled() {
				ask.Status = types.OrderStatusFilled
				k.removeOrderFromBook(ctx, ask)
			} else {
				ask.Status = types.OrderStatusPartial
			}
			k.SetOrder(ctx, ask)

			// Refresh pool. H4: if pool is gone, stop matching this side.
			refreshed, ok := k.GetPool(ctx, pool.ID)
			if !ok {
				k.Logger(ctx).Error("matchRemainingAgainstAMM: pool refresh failed (ask loop), aborting matching for pool",
					"pool_id", pool.ID)
				return
			}
			pool = refreshed

			// C1-FIX: per-block budget exhausted — stop filling rather than
			// leaving the remainder marked as anything but genuinely filled.
			if budgetCapped {
				break
			}
		}
	}
}

// ---------------------------------------------------------------------------
// tryImmediateMatch: for IOC/FOK/market orders, try to match immediately
// ---------------------------------------------------------------------------

func (k Keeper) tryImmediateMatch(ctx context.Context, order *types.Order, pool types.Pool) math.Int {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockHeight := sdkCtx.BlockHeight()
	blockTime := sdkCtx.BlockTime()

	filledQty := math.ZeroInt()
	remaining := order.Quantity

	if order.Side == types.OrderSideBuy {
		// Match against asks
		asks := k.getAskOrders(ctx, pool.ID, types.MaxOrdersPerSide)
		for i := range asks {
			ask := &asks[i]
			if remaining.IsZero() {
				break
			}
			if ask.Price.GT(order.Price) && order.OrderType != types.OrderTypeMarket {
				break // No more asks at acceptable price
			}

			matchQty := ask.RemainingQty()
			if remaining.LT(matchQty) {
				matchQty = remaining
			}

			matchPrice := ask.Price

			// Execute match
			if err := k.executeMatch(ctx, order, ask, matchPrice, matchQty, pool, blockHeight, blockTime, "orderbook"); err != nil {
				if errors.Is(err, types.ErrSelfTrade) {
					continue // D-16: skip self-trade, try next ask
				}
				k.Logger(ctx).Error("executeMatch failed in tryImmediateMatch (buy)", "error", err)
				continue
			}

			ask.FilledQty = ask.FilledQty.Add(matchQty)
			ask.LastUpdatedAt = blockHeight
			if ask.IsFilled() {
				ask.Status = types.OrderStatusFilled
				k.removeOrderFromBook(ctx, *ask)
			} else {
				ask.Status = types.OrderStatusPartial
			}
			k.SetOrder(ctx, *ask)

			filledQty = filledQty.Add(matchQty)
			remaining = remaining.Sub(matchQty)
		}

		// Try AMM for remaining quantity
		if remaining.IsPositive() {
			ammFill := k.tryAMMFill(ctx, pool, order.Side, order.Price, remaining, order.Creator, order.ID, blockHeight, blockTime)
			filledQty = filledQty.Add(ammFill)
		}
	} else {
		// Match against bids
		bids := k.getBidOrders(ctx, pool.ID, types.MaxOrdersPerSide)
		for i := range bids {
			bid := &bids[i]
			if remaining.IsZero() {
				break
			}
			if bid.Price.LT(order.Price) && order.OrderType != types.OrderTypeMarket {
				break
			}

			matchQty := bid.RemainingQty()
			if remaining.LT(matchQty) {
				matchQty = remaining
			}

			matchPrice := bid.Price

			if err := k.executeMatch(ctx, bid, order, matchPrice, matchQty, pool, blockHeight, blockTime, "orderbook"); err != nil {
				if errors.Is(err, types.ErrSelfTrade) {
					continue // D-16: skip self-trade, try next bid
				}
				k.Logger(ctx).Error("executeMatch failed in tryImmediateMatch (sell)", "error", err)
				continue
			}

			bid.FilledQty = bid.FilledQty.Add(matchQty)
			bid.LastUpdatedAt = blockHeight
			if bid.IsFilled() {
				bid.Status = types.OrderStatusFilled
				k.removeOrderFromBook(ctx, *bid)
			} else {
				bid.Status = types.OrderStatusPartial
			}
			k.SetOrder(ctx, *bid)

			filledQty = filledQty.Add(matchQty)
			remaining = remaining.Sub(matchQty)
		}

		// Try AMM for remaining quantity
		if remaining.IsPositive() {
			ammFill := k.tryAMMFill(ctx, pool, order.Side, order.Price, remaining, order.Creator, order.ID, blockHeight, blockTime)
			filledQty = filledQty.Add(ammFill)
		}
	}

	return filledQty
}

// tryAMMFill tries to fill remaining order quantity via AMM.
func (k Keeper) tryAMMFill(ctx context.Context, pool types.Pool, side string, limitPrice math.LegacyDec, qty math.Int, creator string, orderID uint64, blockHeight int64, blockTime time.Time) math.Int {
	if side == types.OrderSideBuy {
		// Buy DenomA using DenomB from escrow
		maxTokenInAmount := limitPrice.MulInt(qty).Ceil().TruncateInt()
		maxTokenIn := sdk.NewCoin(pool.DenomB, maxTokenInAmount)

		quote, _, err := k.GetQuote(ctx, pool.ID, maxTokenIn)
		if err != nil || quote.Amount.IsZero() {
			return math.ZeroInt()
		}

		// D-07: Determine the actual DenomB to swap. If the quote gives us more
		// DenomA than requested (qty), reduce the DenomB input proportionally to
		// avoid overpaying and leaking excess DenomB to the module account.
		actualTokenInAmount := maxTokenInAmount
		if quote.Amount.GT(qty) {
			// Scale down: we only need qty out of quote.Amount
			// actualIn = maxIn * qty / quote.Amount (+ 1 for ceiling)
			actualTokenInAmount = maxTokenInAmount.Mul(qty).Quo(quote.Amount).Add(math.OneInt())
			if actualTokenInAmount.GT(maxTokenInAmount) {
				actualTokenInAmount = maxTokenInAmount
			}
		}
		tokenIn := sdk.NewCoin(pool.DenomB, actualTokenInAmount)

		tokenOut, err := k.executeAMMSwapInternal(ctx, pool.ID, tokenIn)
		if err != nil {
			return math.ZeroInt()
		}

		fillQty := tokenOut.Amount
		if fillQty.GT(qty) {
			fillQty = qty
		}

		creatorAddr, _ := sdk.AccAddressFromBech32(creator)
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, sdk.NewCoins(sdk.NewCoin(pool.DenomA, fillQty))); err != nil {
			return math.ZeroInt()
		}

		// D-07: Refund any excess DenomA that the AMM returned beyond fillQty.
		// H3: Previously this swallowed the error. We now check and, on
		// failure, mark the originating order as errored and emit a
		// "refund_failed" event. We still return the executed fillQty so the
		// main fill (which already left the module account) is accounted for.
		excessDenomA := tokenOut.Amount.Sub(fillQty)
		if excessDenomA.IsPositive() {
			excessCoin := sdk.NewCoin(pool.DenomA, excessDenomA)
			if refundErr := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, sdk.NewCoins(excessCoin)); refundErr != nil {
				k.Logger(ctx).Error("tryAMMFill: excess refund failed, marking order errored",
					"order_id", orderID,
					"creator", creator,
					"excess", excessCoin.String(),
					"error", refundErr,
				)
				sdkCtx := sdk.UnwrapSDKContext(ctx)
				sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
					"refund_failed",
					sdk.NewAttribute("order_id", fmt.Sprintf("%d", orderID)),
					sdk.NewAttribute("address", creatorAddr.String()),
					sdk.NewAttribute("coin", excessCoin.String()),
					sdk.NewAttribute("error", refundErr.Error()),
				))
				if errOrder, found := k.GetOrder(ctx, orderID); found {
					errOrder.Status = types.OrderStatusErrored
					errOrder.LastUpdatedAt = blockHeight
					k.SetOrder(ctx, errOrder)
					k.removeOrderFromBook(ctx, errOrder)
				}
				return math.ZeroInt()
			}
		}

		// Record trade
		ammPrice, _ := k.getAMMPrice(ctx, pool, side)
		k.recordTrade(ctx, pool.ID, ammPrice, fillQty, creator, "AMM", orderID, 0, side, "amm", blockHeight, blockTime)

		return fillQty
	}

	// Sell DenomA for DenomB
	tokenIn := sdk.NewCoin(pool.DenomA, qty)
	quote, _, err := k.GetQuote(ctx, pool.ID, tokenIn)
	if err != nil || quote.Amount.IsZero() {
		return math.ZeroInt()
	}

	// Check if AMM price meets limit
	effectivePrice := math.LegacyNewDecFromInt(quote.Amount).Quo(math.LegacyNewDecFromInt(qty))
	if effectivePrice.LT(limitPrice) {
		return math.ZeroInt()
	}

	tokenOut, err := k.executeAMMSwapInternal(ctx, pool.ID, tokenIn)
	if err != nil {
		return math.ZeroInt()
	}

	creatorAddr, _ := sdk.AccAddressFromBech32(creator)
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, sdk.NewCoins(tokenOut)); err != nil {
		return math.ZeroInt()
	}

	ammPrice, _ := k.getAMMPrice(ctx, pool, side)
	k.recordTrade(ctx, pool.ID, ammPrice, qty, creator, "AMM", orderID, 0, side, "amm", blockHeight, blockTime)

	return qty
}

// ---------------------------------------------------------------------------
// executeMatch: transfer funds between maker and taker
// ---------------------------------------------------------------------------

func (k Keeper) executeMatch(ctx context.Context, buyOrder, sellOrder *types.Order, matchPrice math.LegacyDec, matchQty math.Int, pool types.Pool, blockHeight int64, blockTime time.Time, source string) error {
	// D-16: Prevent wash trading — reject self-trades
	if buyOrder.Creator == sellOrder.Creator {
		return types.ErrSelfTrade
	}

	// The buyer's escrow (DenomB) pays the seller
	// The seller's escrow (DenomA) goes to the buyer

	buyerAddr, _ := sdk.AccAddressFromBech32(buyOrder.Creator)
	sellerAddr, _ := sdk.AccAddressFromBech32(sellOrder.Creator)

	// Transfer DenomA from module (seller's escrow) to buyer
	denomAAmount := matchQty
	denomACoin := sdk.NewCoin(pool.DenomA, denomAAmount)
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, buyerAddr, sdk.NewCoins(denomACoin)); err != nil {
		return fmt.Errorf("failed to transfer DenomA to buyer: %w", err)
	}

	// Transfer DenomB from module (buyer's escrow) to seller
	denomBAmount := matchPrice.MulInt(matchQty).TruncateInt()
	denomBCoin := sdk.NewCoin(pool.DenomB, denomBAmount)
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sellerAddr, sdk.NewCoins(denomBCoin)); err != nil {
		return fmt.Errorf("failed to transfer DenomB to seller: %w", err)
	}

	// Refund buyer's excess escrow: calculate as (escrowed - payment) to avoid
	// rounding dust from independent truncation of escrow and payment amounts.
	escrowAmount := buyOrder.Price.MulInt(matchQty).Ceil().TruncateInt()
	refundAmount := escrowAmount.Sub(denomBAmount)
	if refundAmount.IsPositive() {
		refundCoin := sdk.NewCoin(pool.DenomB, refundAmount)
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, buyerAddr, sdk.NewCoins(refundCoin)); err != nil {
			return fmt.Errorf("failed to refund excess DenomB to buyer: %w", err)
		}
	}

	// Record trade
	k.recordTrade(ctx, pool.ID, matchPrice, matchQty, sellOrder.Creator, buyOrder.Creator, sellOrder.ID, buyOrder.ID, types.OrderSideBuy, source, blockHeight, blockTime)
	return nil
}

// ---------------------------------------------------------------------------
// executeAMMSwapInternal: execute swap from module account funds (already escrowed)
// ---------------------------------------------------------------------------

func (k Keeper) executeAMMSwapInternal(ctx context.Context, poolID uint64, tokenIn sdk.Coin) (sdk.Coin, error) {
	// Risk engine: check if pool is halted or volume limit exceeded
	if err := k.CheckSwapRisk(ctx, poolID, tokenIn.Amount); err != nil {
		return sdk.Coin{}, err
	}

	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return sdk.Coin{}, types.ErrPoolNotFound
	}

	var reserveIn, reserveOut *math.Int
	var denomOut string

	if tokenIn.Denom == pool.DenomA {
		reserveIn = &pool.ReserveA
		reserveOut = &pool.ReserveB
		denomOut = pool.DenomB
	} else if tokenIn.Denom == pool.DenomB {
		reserveIn = &pool.ReserveB
		reserveOut = &pool.ReserveA
		denomOut = pool.DenomA
	} else {
		return sdk.Coin{}, types.ErrInvalidDenom
	}

	// Calculate output using constant-product formula with dynamic fee
	effectiveFee := k.GetEffectiveFee(ctx, poolID)
	feeBps := effectiveFee.MulInt64(10000).TruncateInt64()
	amountInAfterFee := tokenIn.Amount.MulRaw(10000 - feeBps).QuoRaw(10000)
	outputAmount := reserveOut.Mul(amountInAfterFee).Quo(reserveIn.Add(amountInAfterFee))

	if outputAmount.IsZero() || outputAmount.GT(*reserveOut) {
		return sdk.Coin{}, types.ErrInsufficientLiquidity
	}

	// Update reserves — funds are already in module account
	if tokenIn.Denom == pool.DenomA {
		pool.ReserveA = pool.ReserveA.Add(tokenIn.Amount)
		pool.ReserveB = pool.ReserveB.Sub(outputAmount)
	} else {
		pool.ReserveB = pool.ReserveB.Add(tokenIn.Amount)
		pool.ReserveA = pool.ReserveA.Sub(outputAmount)
	}
	k.SetPool(ctx, pool)

	// Record swap volume in risk engine AFTER successful swap
	k.RecordSwapVolume(ctx, poolID, tokenIn.Amount)

	// Record trade volume for fee tier tracking
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	k.RecordTradeVolume(sdkCtx, "", tokenIn.Amount)

	return sdk.NewCoin(denomOut, outputAmount), nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func (k Keeper) getAMMPrice(ctx context.Context, pool types.Pool, side string) (math.LegacyDec, error) {
	// Price is DenomB per DenomA
	if pool.ReserveA.IsZero() || pool.ReserveB.IsZero() {
		return math.LegacyDec{}, fmt.Errorf("zero reserves")
	}
	price := math.LegacyNewDecFromInt(pool.ReserveB).Quo(math.LegacyNewDecFromInt(pool.ReserveA))
	return price, nil
}

func (k Keeper) calculateEscrow(pool types.Pool, side string, price math.LegacyDec, quantity math.Int) sdk.Coin {
	if side == types.OrderSideBuy {
		// Buying DenomA, escrow DenomB (price * quantity)
		amount := price.MulInt(quantity).Ceil().TruncateInt()
		return sdk.NewCoin(pool.DenomB, amount)
	}
	// Selling DenomA, escrow DenomA (the quantity itself)
	return sdk.NewCoin(pool.DenomA, quantity)
}

// refundEscrow attempts to return escrowed funds to `addr`. It returns an
// error if the bank send fails so callers can react (rollback, log, emit a
// "refund_failed" event) instead of silently dropping the user's funds.
//
// H3: Previously this swallowed the error with `_ =`, which could brick orders
// or burn user funds without any audit trail.
func (k Keeper) refundEscrow(ctx context.Context, addr sdk.AccAddress, coin sdk.Coin) error {
	if !coin.Amount.IsPositive() {
		return nil
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(coin)); err != nil {
		k.Logger(ctx).Error("refundEscrow failed",
			"address", addr.String(),
			"coin", coin.String(),
			"error", err,
		)
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
			"refund_failed",
			sdk.NewAttribute("address", addr.String()),
			sdk.NewAttribute("coin", coin.String()),
			sdk.NewAttribute("error", err.Error()),
		))
		return err
	}
	return nil
}

// ---------------------------------------------------------------------------
// Trade recording
// ---------------------------------------------------------------------------

func (k Keeper) recordTrade(ctx context.Context, poolID uint64, price math.LegacyDec, quantity math.Int, makerAddr, takerAddr string, makerOrderID, takerOrderID uint64, side, source string, blockHeight int64, blockTime time.Time) {
	tradeID := k.GetNextTradeID(ctx)
	k.SetNextTradeID(ctx, tradeID+1)

	trade := types.Trade{
		ID:          tradeID,
		PoolID:      poolID,
		Price:       price,
		Quantity:    quantity,
		MakerAddr:   makerAddr,
		TakerAddr:   takerAddr,
		MakerOrder:  makerOrderID,
		TakerOrder:  takerOrderID,
		Side:        side,
		Source:      source,
		Timestamp:   blockTime,
		BlockHeight: blockHeight,
	}

	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(trade)
	kvStore.Set(types.TradeHistoryKey(poolID, tradeID), bz)

	// Update OHLCV candles with this trade's price and volume
	k.UpdateCandles(ctx, poolID, price, quantity)

	// Bump per-pool counter (O(1))
	count := k.getPoolTradeCount(ctx, poolID) + 1
	k.setPoolTradeCount(ctx, poolID, count)

	// Only trim when over cap — uses the counter to avoid iterating every swap.
	if count > uint64(types.MaxTradesPerPool) {
		k.trimOldestTrade(ctx, poolID)
	}
}

// getPoolTradeCount returns the per-pool trade count (0 if unset).
func (k Keeper) getPoolTradeCount(ctx context.Context, poolID uint64) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.PoolTradeCountKey(poolID))
	if err != nil || bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

// setPoolTradeCount persists the per-pool trade count.
func (k Keeper) setPoolTradeCount(ctx context.Context, poolID uint64, count uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	kvStore.Set(types.PoolTradeCountKey(poolID), bz)
}

// trimOldestTrade deletes exactly one oldest trade for the given pool and
// decrements the per-pool counter. Called only when count > MaxTradesPerPool,
// so it runs at most once per recordTrade. O(1) amortized.
func (k Keeper) trimOldestTrade(ctx context.Context, poolID uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.TradeHistoryPoolPrefix(poolID)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return
	}
	defer iter.Close()

	if !iter.Valid() {
		return
	}
	oldestKey := append([]byte{}, iter.Key()...)
	kvStore.Delete(oldestKey)

	count := k.getPoolTradeCount(ctx, poolID)
	if count > 0 {
		k.setPoolTradeCount(ctx, poolID, count-1)
	}
}

// GetTradesByPool returns the persisted AMM + orderbook trade history for a pool,
// newest first, capped at `limit` (or all if limit == 0).
func (k Keeper) GetTradesByPool(ctx context.Context, poolID uint64, limit int) []types.Trade {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.TradeHistoryPoolPrefix(poolID)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var trades []types.Trade
	for ; iter.Valid(); iter.Next() {
		var t types.Trade
		if err := json.Unmarshal(iter.Value(), &t); err != nil {
			continue
		}
		trades = append(trades, t)
	}

	// Reverse so newest is first (keys are ascending tradeID = ascending time).
	for i, j := 0, len(trades)-1; i < j; i, j = i+1, j-1 {
		trades[i], trades[j] = trades[j], trades[i]
	}

	if limit > 0 && len(trades) > limit {
		trades = trades[:limit]
	}
	return trades
}


// ---------------------------------------------------------------------------
// Conditional Order Index Management
// ---------------------------------------------------------------------------

func (k Keeper) addConditionalOrder(ctx context.Context, order types.Order) {
	kvStore := k.storeService.OpenKVStore(ctx)
	orderIDBz := make([]byte, 8)
	binary.BigEndian.PutUint64(orderIDBz, order.ID)
	kvStore.Set(types.ConditionalOrderKey(order.PoolID, order.ID), orderIDBz)
}

func (k Keeper) removeConditionalOrder(ctx context.Context, order types.Order) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Delete(types.ConditionalOrderKey(order.PoolID, order.ID))
	// Also clean address index
	kvStore.Delete(types.OrdersByAddressKey(order.Creator, order.ID))
}

// ---------------------------------------------------------------------------
// ProcessConditionalOrders: checks all pending stop-loss/take-profit orders
// against the current pool price and triggers them. Called in BeginBlock
// BEFORE RunOrderBookMatching so triggered orders get matched in the same block.
// ---------------------------------------------------------------------------

func (k Keeper) ProcessConditionalOrders(ctx sdk.Context) {
	pools := k.GetAllPools(ctx)
	for _, pool := range pools {
		k.processConditionalOrdersForPool(ctx, pool)
	}
}

func (k Keeper) processConditionalOrdersForPool(ctx context.Context, pool types.Pool) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockHeight := sdkCtx.BlockHeight()

	// Get current pool price (DenomB per DenomA)
	currentPrice, err := k.getAMMPrice(ctx, pool, types.OrderSideSell)
	if err != nil {
		return // no liquidity, skip
	}

	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.ConditionalOrderPoolPrefix(pool.ID)

	iter, iterErr := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if iterErr != nil {
		return
	}

	// Collect order IDs to process (avoid modifying store during iteration)
	var orderIDs []uint64
	for ; iter.Valid(); iter.Next() {
		orderID := binary.BigEndian.Uint64(iter.Value())
		orderIDs = append(orderIDs, orderID)
	}
	iter.Close()

	for _, orderID := range orderIDs {
		order, found := k.GetOrder(ctx, orderID)
		if !found || order.Status != types.OrderStatusPending {
			// Order was cancelled or already processed — clean up stale index entry
			kvStore.Delete(types.ConditionalOrderKey(pool.ID, orderID))
			continue
		}

		triggered := false
		switch order.OrderType {
		case types.OrderTypeStopLoss, types.OrderTypeStopLossLimit:
			if order.Side == types.OrderSideSell {
				// Sell stop-loss: "sell when price DROPS to X" — exit a long to limit losses
				triggered = currentPrice.LTE(order.TriggerPrice)
			} else {
				// Buy stop-loss: "buy when price RISES to X" — exit a short to limit losses
				triggered = currentPrice.GTE(order.TriggerPrice)
			}
		case types.OrderTypeTakeProfit, types.OrderTypeTakeProfitLimit:
			if order.Side == types.OrderSideSell {
				// Sell take-profit: "sell when price RISES to X" — take profit on a long
				triggered = currentPrice.GTE(order.TriggerPrice)
			} else {
				// Buy take-profit: "buy when price DROPS to X" — take profit on a short
				triggered = currentPrice.LTE(order.TriggerPrice)
			}
		}

		if !triggered {
			continue
		}

		// Mark as triggered
		order.Status = types.OrderStatusTriggered
		order.LastUpdatedAt = blockHeight
		k.SetOrder(ctx, order)

		// Remove from conditional index
		k.removeConditionalOrder(ctx, order)

		k.Logger(ctx).Info("conditional order triggered",
			"order_id", order.ID,
			"type", order.OrderType,
			"trigger_price", order.TriggerPrice,
			"current_price", currentPrice,
		)

		// Convert to executable order
		if order.OrderType == types.OrderTypeStopLoss || order.OrderType == types.OrderTypeTakeProfit {
			// Market-style: convert to market order and add to book for immediate matching.
			// Keep order.Price at TriggerPrice (the escrow price) so refund math stays correct.
			order.OrderType = types.OrderTypeMarket
			order.Status = types.OrderStatusOpen
			k.SetOrder(ctx, order)
			k.addOrderToBook(ctx, order)
		} else {
			// Limit-style (stop_loss_limit, take_profit_limit): convert to limit order
			// Price was already set from msg.Price at placement time
			order.OrderType = types.OrderTypeLimit
			order.Status = types.OrderStatusOpen
			k.SetOrder(ctx, order)
			k.addOrderToBook(ctx, order)
		}
	}
}

// ---------------------------------------------------------------------------
// Cleanup expired orders
// ---------------------------------------------------------------------------

func (k Keeper) cleanupExpiredOrders(ctx context.Context, poolID uint64) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockHeight := sdkCtx.BlockHeight()

	// Check bids
	bids := k.getBidOrders(ctx, poolID, types.MaxOrdersPerSide)
	for _, bid := range bids {
		if bid.ExpiresAt > 0 && blockHeight >= bid.ExpiresAt {
			k.expireOrder(ctx, bid)
		}
	}

	// Check asks
	asks := k.getAskOrders(ctx, poolID, types.MaxOrdersPerSide)
	for _, ask := range asks {
		if ask.ExpiresAt > 0 && blockHeight >= ask.ExpiresAt {
			k.expireOrder(ctx, ask)
		}
	}

	// Check conditional (pending) orders
	k.cleanupExpiredConditionalOrders(ctx, poolID)
}

func (k Keeper) cleanupExpiredConditionalOrders(ctx context.Context, poolID uint64) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockHeight := sdkCtx.BlockHeight()
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.ConditionalOrderPoolPrefix(poolID)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return
	}

	type staleEntry struct {
		id    uint64
		order types.Order
		found bool
	}
	var toExpire []staleEntry
	for ; iter.Valid(); iter.Next() {
		orderID := binary.BigEndian.Uint64(iter.Value())
		order, found := k.GetOrder(ctx, orderID)
		if !found {
			toExpire = append(toExpire, staleEntry{id: orderID, found: false})
			continue
		}
		if order.Status == types.OrderStatusPending && order.ExpiresAt > 0 && blockHeight >= order.ExpiresAt {
			toExpire = append(toExpire, staleEntry{id: orderID, order: order, found: true})
		}
	}
	iter.Close()

	for _, entry := range toExpire {
		if !entry.found {
			kvStore.Delete(types.ConditionalOrderKey(poolID, entry.id))
			continue
		}
		k.expireConditionalOrder(ctx, entry.order)
	}
}

func (k Keeper) expireConditionalOrder(ctx context.Context, order types.Order) {
	pool, found := k.GetPool(ctx, order.PoolID)
	if !found {
		return
	}

	creatorAddr, _ := sdk.AccAddressFromBech32(order.Creator)

	// Refund escrowed funds
	unfilledQty := order.RemainingQty()
	if unfilledQty.IsPositive() {
		refundCoin := k.calculateEscrow(pool, order.Side, order.Price, unfilledQty)
		if err := k.refundEscrow(ctx, creatorAddr, refundCoin); err != nil {
			k.Logger(ctx).Error("expireConditionalOrder: refund failed, leaving order pending for retry",
				"order_id", order.ID, "error", err)
			return
		}
	}

	k.removeConditionalOrder(ctx, order)
	order.Status = types.OrderStatusExpired
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	order.LastUpdatedAt = sdkCtx.BlockHeight()
	k.SetOrder(ctx, order)

	k.Logger(ctx).Info("conditional order expired", "order_id", order.ID)
}

func (k Keeper) expireOrder(ctx context.Context, order types.Order) {
	pool, found := k.GetPool(ctx, order.PoolID)
	if !found {
		return
	}

	creatorAddr, _ := sdk.AccAddressFromBech32(order.Creator)

	// Refund unfilled portion. expireOrder is called from BeginBlock, so we
	// cannot propagate an error to a tx. H3: log + emit event inside
	// refundEscrow. If the refund fails, leave the order in place rather than
	// silently expiring it with funds stuck in the module account — the
	// keeper will retry on the next block.
	unfilledQty := order.RemainingQty()
	if unfilledQty.IsPositive() {
		refundCoin := k.calculateEscrow(pool, order.Side, order.Price, unfilledQty)
		if err := k.refundEscrow(ctx, creatorAddr, refundCoin); err != nil {
			k.Logger(ctx).Error("expireOrder: refund failed, leaving order open for retry",
				"order_id", order.ID, "error", err)
			return
		}
	}

	k.removeOrderFromBook(ctx, order)
	order.Status = types.OrderStatusExpired
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	order.LastUpdatedAt = sdkCtx.BlockHeight()
	k.SetOrder(ctx, order)

	k.Logger(ctx).Info("order expired", "order_id", order.ID)
}

// ---------------------------------------------------------------------------
// Query helpers
// ---------------------------------------------------------------------------

func (k Keeper) GetOrderBook(ctx context.Context, poolID uint64) types.OrderBookSummary {
	bids := k.getBidOrders(ctx, poolID, types.MaxOrdersPerSide)
	asks := k.getAskOrders(ctx, poolID, types.MaxOrdersPerSide)

	// Aggregate bids by price level
	bidLevels := aggregateOrders(bids)
	askLevels := aggregateOrders(asks)

	return types.OrderBookSummary{
		PoolID: poolID,
		Bids:   bidLevels,
		Asks:   askLevels,
	}
}

func aggregateOrders(orders []types.Order) []types.OrderBookLevel {
	levelMap := make(map[string]*types.OrderBookLevel)
	var levelOrder []string

	for _, o := range orders {
		key := o.Price.String()
		if lvl, exists := levelMap[key]; exists {
			lvl.Quantity = lvl.Quantity.Add(o.RemainingQty())
			lvl.Orders++
		} else {
			levelMap[key] = &types.OrderBookLevel{
				Price:    o.Price,
				Quantity: o.RemainingQty(),
				Orders:   1,
			}
			levelOrder = append(levelOrder, key)
		}
	}

	var levels []types.OrderBookLevel
	for _, key := range levelOrder {
		levels = append(levels, *levelMap[key])
	}
	return levels
}

func (k Keeper) GetOrdersByAddress(ctx context.Context, address string) []types.Order {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.OrdersByAddressAddrPrefix(address)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var orders []types.Order
	for ; iter.Valid(); iter.Next() {
		orderID := binary.BigEndian.Uint64(iter.Value())
		order, found := k.GetOrder(ctx, orderID)
		if found {
			orders = append(orders, order)
		}
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].ID < orders[j].ID
	})
	return orders
}

func (k Keeper) GetTradeHistory(ctx context.Context, poolID uint64, limit int) []types.Trade {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.TradeHistoryPoolPrefix(poolID)

	// We want most recent first, so iterate in reverse
	iter, err := kvStore.ReverseIterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var trades []types.Trade
	for ; iter.Valid() && len(trades) < limit; iter.Next() {
		var trade types.Trade
		if err := json.Unmarshal(iter.Value(), &trade); err != nil {
			continue
		}
		trades = append(trades, trade)
	}
	return trades
}

// suppress unused variable warnings
var _ = time.Now
