package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/dex/types"
)

// ---------------------------------------------------------------------------
// IBC Order Packet Data Types
// ---------------------------------------------------------------------------

// IBCOrderPacketData defines the IBC packet payload for cross-chain order submission.
type IBCOrderPacketData struct {
	OrderType   string `json:"order_type"`   // "limit" or "market"
	Side        string `json:"side"`         // "buy" or "sell"
	PoolID      uint64 `json:"pool_id"`
	Price       string `json:"price"`
	Quantity    string `json:"quantity"`
	TimeInForce string `json:"time_in_force"`
	Sender      string `json:"sender"`       // sender on source chain
	SourceChain string `json:"source_chain"` // chain-id
}

// IBCOrderAckData defines the acknowledgement returned after processing an IBC order.
type IBCOrderAckData struct {
	OrderID uint64 `json:"order_id"`
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
}

// IBCOrderInfo stores metadata about a cross-chain order for tracking.
type IBCOrderInfo struct {
	OrderID     uint64 `json:"order_id"`
	PoolID      uint64 `json:"pool_id"`
	OrderType   string `json:"order_type"`
	Side        string `json:"side"`
	Price       string `json:"price"`
	Quantity    string `json:"quantity"`
	TimeInForce string `json:"time_in_force"`
	Sender      string `json:"sender"`
	SourceChain string `json:"source_chain"`
	LocalAddr   string `json:"local_addr"`   // escrow address on Syreen
	Status      string `json:"status"`
	CreatedAt   int64  `json:"created_at"`
}

// IBCReturnInfo stores pending IBC return transfer information.
type IBCReturnInfo struct {
	OrderID     uint64    `json:"order_id"`
	ReturnCoins sdk.Coins `json:"return_coins"`
	DestChannel string    `json:"dest_channel"`
	DestPort    string    `json:"dest_port"`
	Sender      string    `json:"sender"`       // original sender on source chain
	SourceChain string    `json:"source_chain"`
	Queued      bool      `json:"queued"`
}

// ---------------------------------------------------------------------------
// Store Key Prefixes
// ---------------------------------------------------------------------------

const (
	IBCOrderPrefix            = "ibc_order/"
	IBCOrderReturnPrefix      = "ibc_order_return/"
	IBCOrderCounterKey        = "ibc_order_counter"
	IBCAllowedChainsKey       = "ibc_allowed_chains"
	IBCPendingReturnPrefix    = "ibc_pending_return/" // orderID -> PendingReturnInfo
	IBCPendingReturnStatusKey = "pending_return"      // status value for pending returns
)

// PendingReturnInfo stores information about tokens waiting for IBC return.
type PendingReturnInfo struct {
	OrderID     uint64    `json:"order_id"`
	ReturnCoins sdk.Coins `json:"return_coins"`
	Sender      string    `json:"sender"`       // original sender on source chain
	SourceChain string    `json:"source_chain"`
	EscrowAddr  string    `json:"escrow_addr"`  // local escrow address holding the funds
	Status      string    `json:"status"`        // "pending_return" or "claimed"
	CreatedAt   int64     `json:"created_at"`
}

// ibcOrderKey returns the store key for a cross-chain order.
// Format: ibc_order/<sourceChain>/<senderAddr>/<orderID>
func ibcOrderKey(sourceChain, sender string, orderID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, orderID)
	return append([]byte(IBCOrderPrefix+sourceChain+"/"+sender+"/"), bz...)
}

// ibcOrderChainPrefix returns the prefix for all orders from a source chain.
func ibcOrderChainPrefix(sourceChain string) []byte {
	return []byte(IBCOrderPrefix + sourceChain + "/")
}

// ibcOrderChainSenderPrefix returns the prefix for all orders from a specific sender on a chain.
func ibcOrderChainSenderPrefix(sourceChain, sender string) []byte {
	return []byte(IBCOrderPrefix + sourceChain + "/" + sender + "/")
}

// ibcReturnKey returns the store key for a pending IBC return.
// Format: ibc_order_return/<orderID>
func ibcReturnKey(orderID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, orderID)
	return append([]byte(IBCOrderReturnPrefix), bz...)
}

// ---------------------------------------------------------------------------
// Escrow Address Derivation
// ---------------------------------------------------------------------------

// DeriveIBCEscrowAddress deterministically generates a local escrow address
// for a remote sender so that on-chain orders can be attributed to them.
// Uses SHA-256(sourceChain + "/" + sender) truncated to 20 bytes.
func DeriveIBCEscrowAddress(sourceChain, sender string) sdk.AccAddress {
	hash := sha256.Sum256([]byte("ibc-order-escrow/" + sourceChain + "/" + sender))
	return sdk.AccAddress(hash[:20])
}

// ---------------------------------------------------------------------------
// ProcessIBCOrder handles an incoming IBC order packet.
// It validates the packet, maps the remote sender to a local escrow address,
// verifies the escrow has sufficient funds, and places the order.
// Returns the order ID or an error.
// ---------------------------------------------------------------------------

func (k Keeper) ProcessIBCOrder(ctx context.Context, packet IBCOrderPacketData) (uint64, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// TODO: Source chain verification MUST come from IBC channel metadata (the
	// channel handshake binds a channel to a specific counterparty chain-id),
	// NOT from the packet payload which is attacker-controlled. Until full IBC
	// module integration, we enforce a whitelist of allowed source chains.
	if !k.IsChainAllowed(ctx, packet.SourceChain) {
		return 0, fmt.Errorf("source chain %q is not in the allowed chains whitelist; IBC orders are disabled until explicitly configured", packet.SourceChain)
	}

	// Validate order type
	if packet.OrderType != types.OrderTypeLimit && packet.OrderType != types.OrderTypeMarket {
		return 0, fmt.Errorf("invalid IBC order type: %s (must be 'limit' or 'market')", packet.OrderType)
	}

	// Validate side
	if packet.Side != types.OrderSideBuy && packet.Side != types.OrderSideSell {
		return 0, fmt.Errorf("invalid IBC order side: %s (must be 'buy' or 'sell')", packet.Side)
	}

	// Validate pool exists
	_, found := k.GetPool(ctx, packet.PoolID)
	if !found {
		return 0, types.ErrPoolNotFound
	}

	// Parse price
	if packet.OrderType == types.OrderTypeLimit && (packet.Price == "" || packet.Price == "0") {
		return 0, fmt.Errorf("limit orders require a non-zero price")
	}

	// Parse quantity
	parsedQuantity, ok := math.NewIntFromString(packet.Quantity)
	if !ok || !parsedQuantity.IsPositive() {
		return 0, fmt.Errorf("invalid quantity: %s", packet.Quantity)
	}

	// Validate source chain and sender
	if packet.SourceChain == "" {
		return 0, fmt.Errorf("source_chain is required")
	}
	if packet.Sender == "" {
		return 0, fmt.Errorf("sender is required")
	}

	// Derive local escrow address for the remote sender
	escrowAddr := DeriveIBCEscrowAddress(packet.SourceChain, packet.Sender)
	localAddr := escrowAddr.String()

	// Determine time-in-force, default to GTC
	tif := packet.TimeInForce
	if tif == "" {
		tif = types.TimeInForceGTC
	}

	// Verify the escrow address has been funded (IBC transfer must arrive before order).
	// PlaceOrder will also fail on insufficient funds during escrow, but we check
	// upfront for a clearer error message.
	pool, _ := k.GetPool(ctx, packet.PoolID)
	var requiredDenom string
	if packet.Side == types.OrderSideBuy {
		requiredDenom = pool.DenomB // buying base, need quote
	} else {
		requiredDenom = pool.DenomA // selling base, need base
	}
	escrowBal := k.bankKeeper.GetBalance(sdkCtx, escrowAddr, requiredDenom)
	if escrowBal.IsZero() {
		return 0, fmt.Errorf("escrow address %s has no %s balance; IBC transfer must arrive before order", localAddr, requiredDenom)
	}

	// Build the MsgPlaceOrder using the escrow address as the creator
	msg := &types.MsgPlaceOrder{
		Creator:      localAddr,
		PoolID:       packet.PoolID,
		Side:         packet.Side,
		OrderType:    packet.OrderType,
		Price:        packet.Price,
		Quantity:     packet.Quantity,
		TimeInForce:  tif,
		TriggerPrice: "0",
	}

	// Place the order through the existing order book logic
	orderID, status, err := k.PlaceOrder(ctx, msg)
	if err != nil {
		return 0, fmt.Errorf("failed to place IBC order: %w", err)
	}

	// Store the IBC order tracking info
	info := IBCOrderInfo{
		OrderID:     orderID,
		PoolID:      packet.PoolID,
		OrderType:   packet.OrderType,
		Side:        packet.Side,
		Price:       packet.Price,
		Quantity:    packet.Quantity,
		TimeInForce: tif,
		Sender:      packet.Sender,
		SourceChain: packet.SourceChain,
		LocalAddr:   localAddr,
		Status:      status,
		CreatedAt:   sdkCtx.BlockHeight(),
	}
	k.SetIBCOrder(ctx, info)

	k.Logger(ctx).Info("processed IBC order",
		"order_id", orderID,
		"pool_id", packet.PoolID,
		"source_chain", packet.SourceChain,
		"sender", packet.Sender,
		"local_addr", localAddr,
		"status", status,
	)

	return orderID, nil
}

// ---------------------------------------------------------------------------
// QueueIBCReturn queues output tokens for IBC transfer back to the source chain.
// ---------------------------------------------------------------------------

func (k Keeper) QueueIBCReturn(ctx context.Context, orderID uint64, returnCoins sdk.Coins, destChannel, destPort string) {
	// Look up the IBC order to get sender/chain info
	info, found := k.GetIBCOrderByID(ctx, orderID)
	if !found {
		k.Logger(ctx).Error("cannot queue IBC return: order not found", "order_id", orderID)
		return
	}

	ret := IBCReturnInfo{
		OrderID:     orderID,
		ReturnCoins: returnCoins,
		DestChannel: destChannel,
		DestPort:    destPort,
		Sender:      info.Sender,
		SourceChain: info.SourceChain,
		Queued:      true,
	}

	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(ret)
	if err != nil {
		k.Logger(ctx).Error("failed to marshal IBC return", "error", err)
		return
	}
	kvStore.Set(ibcReturnKey(orderID), bz)

	k.Logger(ctx).Info("queued IBC return",
		"order_id", orderID,
		"return_coins", returnCoins.String(),
		"dest_channel", destChannel,
		"source_chain", info.SourceChain,
		"sender", info.Sender,
	)
}

// ---------------------------------------------------------------------------
// ProcessIBCReturns processes all pending IBC returns.
// Called in EndBlock. For now, it marks returns as processed and emits events.
// Actual IBC transfer execution requires the IBC transfer keeper, which will
// be wired when the full IBC module integration is completed in app.go.
// ---------------------------------------------------------------------------

func (k Keeper) ProcessIBCReturns(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	kvStore := k.storeService.OpenKVStore(ctx)

	prefix := []byte(IBCOrderReturnPrefix)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return
	}

	// Collect entries to process — do NOT mutate KV store during iteration
	type returnAction struct {
		key []byte
		ret IBCReturnInfo
	}
	var actions []returnAction

	for ; iter.Valid(); iter.Next() {
		var ret IBCReturnInfo
		if err := json.Unmarshal(iter.Value(), &ret); err != nil {
			continue
		}
		if !ret.Queued {
			continue
		}
		actions = append(actions, returnAction{
			key: append([]byte(nil), iter.Key()...),
			ret: ret,
		})
	}
	iter.Close()

	// Apply modifications after iterator is closed.
	// Since the IBC transfer keeper isn't wired yet, we send the return tokens
	// to the escrow address so they are recoverable, and store a PendingReturn
	// record that users can query. Once IBC is fully wired, ClaimIBCReturn or
	// automatic IBC transfer will handle the actual cross-chain return.
	for _, action := range actions {
		// Look up the IBC order to get the escrow address
		orderInfo, found := k.GetIBCOrderByID(ctx, action.ret.OrderID)
		escrowAddr := ""
		if found {
			escrowAddr = orderInfo.LocalAddr
		}

		// Store as PendingReturn for user query and manual claim
		pendingReturn := PendingReturnInfo{
			OrderID:     action.ret.OrderID,
			ReturnCoins: action.ret.ReturnCoins,
			Sender:      action.ret.Sender,
			SourceChain: action.ret.SourceChain,
			EscrowAddr:  escrowAddr,
			Status:      IBCPendingReturnStatusKey,
			CreatedAt:   sdkCtx.BlockHeight(),
		}
		pendingBz, _ := json.Marshal(pendingReturn)
		kvStore.Set(ibcPendingReturnKey(action.ret.OrderID), pendingBz)

		sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
			"ibc_order_return_pending",
			sdk.NewAttribute("order_id", fmt.Sprintf("%d", action.ret.OrderID)),
			sdk.NewAttribute("return_coins", action.ret.ReturnCoins.String()),
			sdk.NewAttribute("sender", action.ret.Sender),
			sdk.NewAttribute("source_chain", action.ret.SourceChain),
			sdk.NewAttribute("escrow_addr", escrowAddr),
			sdk.NewAttribute("status", "pending_return"),
		))

		// Mark original return entry as processed
		action.ret.Queued = false
		bz, _ := json.Marshal(action.ret)
		kvStore.Set(action.key, bz)

		k.Logger(ctx).Info("IBC return held as pending (IBC transfer not wired yet)",
			"order_id", action.ret.OrderID,
			"return_coins", action.ret.ReturnCoins.String(),
			"escrow_addr", escrowAddr,
		)
	}
}

// ---------------------------------------------------------------------------
// IBC Order Storage CRUD
// ---------------------------------------------------------------------------

// SetIBCOrder stores an IBC order tracking record.
func (k Keeper) SetIBCOrder(ctx context.Context, info IBCOrderInfo) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(info)
	if err != nil {
		return
	}
	key := ibcOrderKey(info.SourceChain, info.Sender, info.OrderID)
	kvStore.Set(key, bz)
}

// GetIBCOrderByID searches for an IBC order by order ID across all chains.
// This is a convenience method that scans the ibc_order prefix.
func (k Keeper) GetIBCOrderByID(ctx context.Context, orderID uint64) (IBCOrderInfo, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)

	prefix := []byte(IBCOrderPrefix)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return IBCOrderInfo{}, false
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var info IBCOrderInfo
		if err := json.Unmarshal(iter.Value(), &info); err != nil {
			continue
		}
		if info.OrderID == orderID {
			return info, true
		}
	}
	return IBCOrderInfo{}, false
}

// GetIBCOrdersByChain returns all IBC orders from a given source chain.
func (k Keeper) GetIBCOrdersByChain(ctx context.Context, sourceChain string) []IBCOrderInfo {
	kvStore := k.storeService.OpenKVStore(ctx)

	prefix := ibcOrderChainPrefix(sourceChain)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var orders []IBCOrderInfo
	for ; iter.Valid(); iter.Next() {
		var info IBCOrderInfo
		if err := json.Unmarshal(iter.Value(), &info); err != nil {
			continue
		}
		orders = append(orders, info)
	}

	// Sort by order ID for deterministic output
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].OrderID < orders[j].OrderID
	})
	return orders
}

// GetIBCOrdersByChainAndSender returns all IBC orders from a specific sender on a chain.
func (k Keeper) GetIBCOrdersByChainAndSender(ctx context.Context, sourceChain, sender string) []IBCOrderInfo {
	kvStore := k.storeService.OpenKVStore(ctx)

	prefix := ibcOrderChainSenderPrefix(sourceChain, sender)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var orders []IBCOrderInfo
	for ; iter.Valid(); iter.Next() {
		var info IBCOrderInfo
		if err := json.Unmarshal(iter.Value(), &info); err != nil {
			continue
		}
		orders = append(orders, info)
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].OrderID < orders[j].OrderID
	})
	return orders
}

// ---------------------------------------------------------------------------
// OnRecvIBCOrderPacket handles an incoming IBC order packet.
// This is the entry point called from the IBC module's OnRecvPacket handler.
// ---------------------------------------------------------------------------

func (k Keeper) OnRecvIBCOrderPacket(ctx context.Context, packetData []byte) IBCOrderAckData {
	var packet IBCOrderPacketData
	if err := json.Unmarshal(packetData, &packet); err != nil {
		return IBCOrderAckData{
			Status: "error",
			Error:  fmt.Sprintf("failed to unmarshal IBC order packet: %v", err),
		}
	}

	orderID, err := k.ProcessIBCOrder(ctx, packet)
	if err != nil {
		return IBCOrderAckData{
			Status: "error",
			Error:  err.Error(),
		}
	}

	return IBCOrderAckData{
		OrderID: orderID,
		Status:  "accepted",
	}
}

// ---------------------------------------------------------------------------
// Query: IBCOrders
// ---------------------------------------------------------------------------

// QueryIBCOrders returns IBC orders filtered by source chain and optionally sender.
func (k Keeper) QueryIBCOrders(ctx context.Context, sourceChain, sender string) []IBCOrderInfo {
	if sender != "" {
		return k.GetIBCOrdersByChainAndSender(ctx, sourceChain, sender)
	}
	return k.GetIBCOrdersByChain(ctx, sourceChain)
}

// ---------------------------------------------------------------------------
// IBC Allowed Chains Whitelist
// ---------------------------------------------------------------------------

// ibcPendingReturnKey returns the store key for a pending IBC return.
func ibcPendingReturnKey(orderID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, orderID)
	return append([]byte(IBCPendingReturnPrefix), bz...)
}

// IsChainAllowed checks if a source chain is in the allowed chains whitelist.
// Default: empty whitelist = all chains rejected until explicitly configured.
func (k Keeper) IsChainAllowed(ctx context.Context, chain string) bool {
	chains := k.GetAllowedChains(ctx)
	for _, c := range chains {
		if c == chain {
			return true
		}
	}
	return false
}

// GetAllowedChains returns the list of allowed IBC source chains.
func (k Keeper) GetAllowedChains(ctx context.Context) []string {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(IBCAllowedChainsKey))
	if err != nil || bz == nil {
		return nil
	}
	var chains []string
	if err := json.Unmarshal(bz, &chains); err != nil {
		return nil
	}
	return chains
}

// SetAllowedChains sets the list of allowed IBC source chains.
func (k Keeper) SetAllowedChains(ctx context.Context, chains []string) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(chains)
	kvStore.Set([]byte(IBCAllowedChainsKey), bz)
}

// ---------------------------------------------------------------------------
// Pending Return Queries
// ---------------------------------------------------------------------------

// GetPendingReturn returns a pending IBC return by order ID.
func (k Keeper) GetPendingReturn(ctx context.Context, orderID uint64) (PendingReturnInfo, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(ibcPendingReturnKey(orderID))
	if err != nil || bz == nil {
		return PendingReturnInfo{}, false
	}
	var info PendingReturnInfo
	if err := json.Unmarshal(bz, &info); err != nil {
		return PendingReturnInfo{}, false
	}
	return info, true
}

// GetAllPendingReturns returns all pending IBC returns.
func (k Keeper) GetAllPendingReturns(ctx context.Context) []PendingReturnInfo {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(IBCPendingReturnPrefix)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var returns []PendingReturnInfo
	for ; iter.Valid(); iter.Next() {
		var info PendingReturnInfo
		if err := json.Unmarshal(iter.Value(), &info); err != nil {
			continue
		}
		if info.Status == IBCPendingReturnStatusKey {
			returns = append(returns, info)
		}
	}
	sort.Slice(returns, func(i, j int) bool {
		return returns[i].OrderID < returns[j].OrderID
	})
	return returns
}

// suppress unused import warning
var _ = hex.EncodeToString
