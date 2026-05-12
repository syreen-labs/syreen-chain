package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/payments/types"
)

type Keeper struct {
	cdc           codec.Codec
	storeService  store.KVStoreService
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	authority     string
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeService:  storeService,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		authority:     authority,
	}
}

func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetAuthority returns the module's authority address.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// --- Params ---

func (k Keeper) GetParams(ctx context.Context) types.Params {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.ParamsKey))
	if err != nil || bz == nil {
		return types.DefaultParams()
	}
	var params types.Params
	if err := json.Unmarshal(bz, &params); err != nil {
		return types.DefaultParams()
	}
	return params
}

func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return kvStore.Set([]byte(types.ParamsKey), bz)
}

// --- Invoice Counter ---

func (k Keeper) GetInvoiceCount(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.InvoiceCounterKey))
	if err != nil || bz == nil {
		return 0
	}
	count, err := strconv.ParseUint(string(bz), 10, 64)
	if err != nil {
		return 0
	}
	return count
}

func (k Keeper) SetInvoiceCount(ctx context.Context, count uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Set([]byte(types.InvoiceCounterKey), []byte(strconv.FormatUint(count, 10)))
}

func (k Keeper) GetNextInvoiceID(ctx context.Context) string {
	count := k.GetInvoiceCount(ctx)
	count++
	k.SetInvoiceCount(ctx, count)
	return fmt.Sprintf("INV-%06d", count)
}

// --- Invoice CRUD ---

func (k Keeper) GetInvoice(ctx context.Context, id string) (types.Invoice, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.InvoiceStoreKey(id))
	if err != nil || bz == nil {
		return types.Invoice{}, false
	}
	var invoice types.Invoice
	if err := json.Unmarshal(bz, &invoice); err != nil {
		return types.Invoice{}, false
	}
	return invoice, true
}

func (k Keeper) SetInvoice(ctx context.Context, invoice types.Invoice) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(invoice)
	kvStore.Set(types.InvoiceStoreKey(invoice.ID), bz)
}

func (k Keeper) GetAllInvoices(ctx context.Context) []types.Invoice {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.InvoiceKeyPrefix)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var invoices []types.Invoice
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		var invoice types.Invoice
		if err := json.Unmarshal(iter.Value(), &invoice); err != nil {
			continue
		}
		invoices = append(invoices, invoice)
	}
	return invoices
}

func (k Keeper) GetInvoicesByBooking(ctx context.Context, bookingID string) []types.Invoice {
	all := k.GetAllInvoices(ctx)
	var result []types.Invoice
	for _, inv := range all {
		if inv.BookingID == bookingID {
			result = append(result, inv)
		}
	}
	return result
}

func (k Keeper) GetInvoicesByPayer(ctx context.Context, payer string) []types.Invoice {
	all := k.GetAllInvoices(ctx)
	var result []types.Invoice
	for _, inv := range all {
		if inv.Payer == payer {
			result = append(result, inv)
		}
	}
	return result
}

// --- ExchangeRate CRUD ---

func (k Keeper) GetExchangeRate(ctx context.Context, fromDenom, toDenom string) (types.ExchangeRate, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.ExchangeRateStoreKey(fromDenom, toDenom))
	if err != nil || bz == nil {
		return types.ExchangeRate{}, false
	}
	var rate types.ExchangeRate
	if err := json.Unmarshal(bz, &rate); err != nil {
		return types.ExchangeRate{}, false
	}
	return rate, true
}

func (k Keeper) SetExchangeRate(ctx context.Context, rate types.ExchangeRate) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(rate)
	kvStore.Set(types.ExchangeRateStoreKey(rate.FromDenom, rate.ToDenom), bz)
}

func (k Keeper) GetAllExchangeRates(ctx context.Context) []types.ExchangeRate {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.ExchangeRateKeyPrefix)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var rates []types.ExchangeRate
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		var rate types.ExchangeRate
		if err := json.Unmarshal(iter.Value(), &rate); err != nil {
			continue
		}
		rates = append(rates, rate)
	}
	return rates
}

// --- Earnings CRUD ---

func (k Keeper) GetEarnings(ctx context.Context, address string) (types.Earnings, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.EarningsStoreKey(address))
	if err != nil || bz == nil {
		return types.Earnings{}, false
	}
	var earnings types.Earnings
	if err := json.Unmarshal(bz, &earnings); err != nil {
		return types.Earnings{}, false
	}
	return earnings, true
}

func (k Keeper) SetEarnings(ctx context.Context, earnings types.Earnings) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(earnings)
	kvStore.Set(types.EarningsStoreKey(earnings.Address), bz)
}

func (k Keeper) GetAllEarnings(ctx context.Context) []types.Earnings {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.EarningsKeyPrefix)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var earnings []types.Earnings
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		var e types.Earnings
		if err := json.Unmarshal(iter.Value(), &e); err != nil {
			continue
		}
		earnings = append(earnings, e)
	}
	return earnings
}

// AddEarnings adds coins to an address's available earnings and total earned.
func (k *Keeper) AddEarnings(ctx context.Context, address string, amount sdk.Coins) {
	earnings, found := k.GetEarnings(ctx, address)
	if !found {
		earnings = types.Earnings{
			Address:        address,
			Available:      sdk.NewCoins(),
			TotalEarned:    sdk.NewCoins(),
			TotalWithdrawn: sdk.NewCoins(),
		}
	}
	earnings.Available = earnings.Available.Add(amount...)
	earnings.TotalEarned = earnings.TotalEarned.Add(amount...)
	k.SetEarnings(ctx, earnings)
}

// --- Business Logic ---

// CreateInvoice creates a new invoice with auto-increment ID, fee calculation, and expiry.
func (k *Keeper) CreateInvoice(ctx context.Context, msg *types.MsgCreateInvoice) (string, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(ctx)

	if !params.EnablePayments {
		return "", types.ErrPaymentsDisabled
	}

	// Require that the invoice creator is the payee (prevent arbitrary payee invoices)
	if msg.Creator != msg.Payee {
		return "", fmt.Errorf("creator %s must be the payee %s", msg.Creator, msg.Payee)
	}

	// Check for duplicate invoice per booking ID
	if msg.BookingID != "" {
		existing := k.GetInvoicesByBooking(ctx, msg.BookingID)
		if len(existing) > 0 {
			return "", fmt.Errorf("invoice already exists for booking %s", msg.BookingID)
		}
	}

	// Generate invoice ID
	invoiceID := k.GetNextInvoiceID(ctx)

	// Calculate fees
	hundred := math.LegacyNewDec(100)
	amountDec := math.LegacyNewDecFromInt(msg.Amount.Amount)

	platformFeeAmt := amountDec.Mul(params.PlatformFeePercent).Quo(hundred).TruncateInt()
	insuranceFeeAmt := amountDec.Mul(params.InsuranceFeePercent).Quo(hundred).TruncateInt()
	netAmt := msg.Amount.Amount.Sub(platformFeeAmt).Sub(insuranceFeeAmt)

	denom := msg.Amount.Denom

	invoice := types.Invoice{
		ID:             invoiceID,
		BookingID:      msg.BookingID,
		Payer:          "", // Set when paid
		Payee:          msg.Payee,
		Amount:         msg.Amount,
		PlatformFee:    sdk.NewCoin(denom, platformFeeAmt),
		InsuranceFee:   sdk.NewCoin(denom, insuranceFeeAmt),
		NetAmount:      sdk.NewCoin(denom, netAmt),
		Status:         types.InvoiceStatusPending,
		PaymentMethod:  "",
		PaymentTxHash:  "",
		CreatedAtBlock: sdkCtx.BlockHeight(),
		PaidAtBlock:    0,
		ExpiresAtBlock: sdkCtx.BlockHeight() + int64(params.InvoiceExpiryBlocks),
		Description:    msg.Description,
		Currency:       msg.Currency,
		Creator:        msg.Creator,
	}

	k.SetInvoice(ctx, invoice)

	k.Logger(ctx).Info("invoice created",
		"id", invoiceID, "booking", msg.BookingID, "amount", msg.Amount,
		"platform_fee", invoice.PlatformFee, "insurance_fee", invoice.InsuranceFee, "net", invoice.NetAmount)

	return invoiceID, nil
}

// PayInvoice processes payment for an invoice. Transfers funds from payer to module,
// then splits: net to payee earnings, platform fee to platform address, insurance fee to insurance pool.
func (k *Keeper) PayInvoice(ctx context.Context, msg *types.MsgPayInvoice) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(ctx)

	invoice, found := k.GetInvoice(ctx, msg.InvoiceID)
	if !found {
		return types.ErrInvoiceNotFound
	}

	if invoice.Status == types.InvoiceStatusPaid {
		return types.ErrInvoiceAlreadyPaid
	}
	if invoice.Status == types.InvoiceStatusCancelled {
		return types.ErrInvoiceCancelled
	}
	if invoice.Status == types.InvoiceStatusRefunded {
		return types.ErrInvalidInvoiceStatus
	}
	if invoice.Status == types.InvoiceStatusExpired {
		return types.ErrInvoiceExpired
	}

	// Check if expired by block height
	if sdkCtx.BlockHeight() >= invoice.ExpiresAtBlock {
		invoice.Status = types.InvoiceStatusExpired
		k.SetInvoice(ctx, invoice)
		return types.ErrInvoiceExpired
	}

	payerAddr, err := sdk.AccAddressFromBech32(msg.Payer)
	if err != nil {
		return fmt.Errorf("invalid payer address: %w", err)
	}

	// Transfer full amount from payer to module account
	coins := sdk.NewCoins(invoice.Amount)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, payerAddr, types.ModuleName, coins); err != nil {
		return types.ErrInsufficientFunds
	}

	// Add net amount to payee earnings
	netCoins := sdk.NewCoins(invoice.NetAmount)
	k.AddEarnings(ctx, invoice.Payee, netCoins)

	// Send platform fee to platform address (if configured)
	// BUG-4 fix: return fee transfer errors so the tx rolls back instead of
	// silently losing the fee while marking the invoice as paid.
	if params.PlatformAddress != "" && invoice.PlatformFee.IsPositive() {
		platformAddr, err := sdk.AccAddressFromBech32(params.PlatformAddress)
		if err != nil {
			return fmt.Errorf("invalid platform address: %w", err)
		}
		platformCoins := sdk.NewCoins(invoice.PlatformFee)
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, platformAddr, platformCoins); err != nil {
			return fmt.Errorf("failed to send platform fee: %w", err)
		}
	} else if invoice.PlatformFee.IsPositive() {
		// If no platform address configured, add to earnings of authority
		k.AddEarnings(ctx, k.authority, sdk.NewCoins(invoice.PlatformFee))
	}

	// Send insurance fee to insurance pool address (if configured)
	if params.InsurancePoolAddress != "" && invoice.InsuranceFee.IsPositive() {
		insuranceAddr, err := sdk.AccAddressFromBech32(params.InsurancePoolAddress)
		if err != nil {
			return fmt.Errorf("invalid insurance pool address: %w", err)
		}
		insuranceCoins := sdk.NewCoins(invoice.InsuranceFee)
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, insuranceAddr, insuranceCoins); err != nil {
			return fmt.Errorf("failed to send insurance fee: %w", err)
		}
	} else if invoice.InsuranceFee.IsPositive() {
		// If no insurance pool address configured, add to earnings of authority
		k.AddEarnings(ctx, k.authority, sdk.NewCoins(invoice.InsuranceFee))
	}

	// Update invoice
	invoice.Status = types.InvoiceStatusPaid
	invoice.Payer = msg.Payer
	invoice.PaymentMethod = msg.PaymentMethod
	invoice.PaidAtBlock = sdkCtx.BlockHeight()
	k.SetInvoice(ctx, invoice)

	k.Logger(ctx).Info("invoice paid",
		"id", invoice.ID, "payer", msg.Payer, "amount", invoice.Amount, "method", msg.PaymentMethod)

	return nil
}

// RefundPayment refunds a paid invoice. Only authority can refund.
func (k *Keeper) RefundPayment(ctx context.Context, msg *types.MsgRefundPayment) error {
	if msg.Authority != k.authority {
		return types.ErrUnauthorized
	}

	invoice, found := k.GetInvoice(ctx, msg.InvoiceID)
	if !found {
		return types.ErrInvoiceNotFound
	}

	if invoice.Status != types.InvoiceStatusPaid {
		return types.ErrInvalidInvoiceStatus
	}

	if invoice.Payer == "" {
		return fmt.Errorf("invoice has no payer recorded")
	}

	payerAddr, err := sdk.AccAddressFromBech32(invoice.Payer)
	if err != nil {
		return fmt.Errorf("invalid payer address on invoice: %w", err)
	}

	// Deduct from payee earnings
	// BUG-3 fix: if the payee has already withdrawn some/all earnings, the module
	// no longer holds the full invoice amount.  Refunding the full amount would
	// create tokens from nothing (double-spend).  Block the refund in that case.
	earnings, found := k.GetEarnings(ctx, invoice.Payee)
	if found {
		netCoins := sdk.NewCoins(invoice.NetAmount)
		remaining, hasNeg := earnings.Available.SafeSub(netCoins...)
		if hasNeg {
			return fmt.Errorf("cannot refund: payee has already withdrawn earnings for invoice %s", invoice.ID)
		}
		earnings.Available = remaining
		k.SetEarnings(ctx, earnings)
	}

	// Refund only the NetAmount (what the module actually holds for this invoice).
	// The platform fee and insurance fee have already been sent to their respective
	// destinations and cannot be reclaimed from the module account.
	refundCoins := sdk.NewCoins(invoice.NetAmount)
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, payerAddr, refundCoins); err != nil {
		return fmt.Errorf("refund transfer failed: %w", err)
	}

	// Only mark as refunded AFTER successful transfer
	invoice.Status = types.InvoiceStatusRefunded
	k.SetInvoice(ctx, invoice)

	k.Logger(ctx).Info("payment refunded",
		"id", invoice.ID, "payer", invoice.Payer, "amount", invoice.Amount, "reason", msg.Reason)

	return nil
}

// CancelInvoice cancels a pending invoice. Only the creator can cancel.
func (k *Keeper) CancelInvoice(ctx context.Context, msg *types.MsgCancelInvoice) error {
	invoice, found := k.GetInvoice(ctx, msg.InvoiceID)
	if !found {
		return types.ErrInvoiceNotFound
	}

	if invoice.Creator != msg.Creator {
		return types.ErrNotInvoiceCreator
	}

	if invoice.Status != types.InvoiceStatusPending {
		return types.ErrInvalidInvoiceStatus
	}

	invoice.Status = types.InvoiceStatusCancelled
	k.SetInvoice(ctx, invoice)

	k.Logger(ctx).Info("invoice cancelled", "id", invoice.ID, "creator", msg.Creator)

	return nil
}

// SetExchangeRateByMsg sets an exchange rate. Only authority can set rates.
func (k *Keeper) SetExchangeRateByMsg(ctx context.Context, msg *types.MsgSetExchangeRate) error {
	if msg.Authority != k.authority {
		return types.ErrUnauthorized
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	rate, err := math.LegacyNewDecFromStr(msg.Rate)
	if err != nil {
		return fmt.Errorf("invalid rate: %w", err)
	}

	exchangeRate := types.ExchangeRate{
		FromDenom:      msg.FromDenom,
		ToDenom:        msg.ToDenom,
		Rate:           rate,
		UpdatedAtBlock: sdkCtx.BlockHeight(),
		UpdatedBy:      msg.Authority,
	}

	k.SetExchangeRate(ctx, exchangeRate)

	k.Logger(ctx).Info("exchange rate set",
		"from", msg.FromDenom, "to", msg.ToDenom, "rate", msg.Rate)

	return nil
}

// WithdrawEarnings withdraws available earnings to the address.
func (k *Keeper) WithdrawEarnings(ctx context.Context, msg *types.MsgWithdrawEarnings) error {
	earnings, found := k.GetEarnings(ctx, msg.Address)
	if !found {
		return types.ErrNoEarningsAvailable
	}

	withdrawCoins := sdk.NewCoins(msg.Amount)

	// Check available balance
	remaining, hasNeg := earnings.Available.SafeSub(withdrawCoins...)
	if hasNeg {
		return types.ErrNoEarningsAvailable
	}

	// Transfer from module to address
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, withdrawCoins); err != nil {
		return fmt.Errorf("withdrawal transfer failed: %w", err)
	}

	earnings.Available = remaining
	earnings.TotalWithdrawn = earnings.TotalWithdrawn.Add(withdrawCoins...)
	k.SetEarnings(ctx, earnings)

	k.Logger(ctx).Info("earnings withdrawn",
		"address", msg.Address, "amount", msg.Amount, "remaining", earnings.Available)

	return nil
}

// ExpireInvoices checks all pending invoices and marks expired ones. Called in BeginBlock.
func (k *Keeper) ExpireInvoices(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.InvoiceKeyPrefix)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return
	}
	defer iter.Close()

	var toExpire []types.Invoice
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		var invoice types.Invoice
		if err := json.Unmarshal(iter.Value(), &invoice); err != nil {
			continue
		}

		if invoice.Status == types.InvoiceStatusPending && invoice.ExpiresAtBlock > 0 && sdkCtx.BlockHeight() >= invoice.ExpiresAtBlock {
			toExpire = append(toExpire, invoice)
		}
	}

	for _, invoice := range toExpire {
		invoice.Status = types.InvoiceStatusExpired
		k.SetInvoice(ctx, invoice)
		k.Logger(ctx).Info("invoice expired", "id", invoice.ID, "expired_at_block", sdkCtx.BlockHeight())
	}
}

// --- Helpers ---

func prefixEndBytes(prefix []byte) []byte {
	if len(prefix) == 0 {
		return nil
	}
	end := make([]byte, len(prefix))
	copy(end, prefix)
	for i := len(end) - 1; i >= 0; i-- {
		end[i]++
		if end[i] != 0 {
			return end
		}
	}
	return nil
}
