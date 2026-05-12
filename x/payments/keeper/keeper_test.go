package keeper_test

import (
	"fmt"
	"testing"

	"context"

	"cosmossdk.io/core/store"
	"cosmossdk.io/math"

	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/log"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	cosmosstore "cosmossdk.io/store"

	"syreen/x/payments/keeper"
	"syreen/x/payments/types"
)

// --- Mock Keepers ---

type mockAccountKeeper struct{}

func (m mockAccountKeeper) GetAccount(_ context.Context, _ sdk.AccAddress) sdk.AccountI { return nil }
func (m mockAccountKeeper) GetModuleAddress(_ string) sdk.AccAddress                    { return nil }

type mockBankKeeper struct {
	balances      map[string]sdk.Coins
	moduleBalance sdk.Coins
	sendErr       error
}

func newMockBankKeeper() *mockBankKeeper {
	return &mockBankKeeper{
		balances:      make(map[string]sdk.Coins),
		moduleBalance: sdk.NewCoins(),
	}
}

func (m *mockBankKeeper) SendCoins(_ context.Context, from, to sdk.AccAddress, amt sdk.Coins) error {
	if m.sendErr != nil {
		return m.sendErr
	}
	return nil
}

func (m *mockBankKeeper) SendCoinsFromAccountToModule(_ context.Context, sender sdk.AccAddress, _ string, amt sdk.Coins) error {
	if m.sendErr != nil {
		return m.sendErr
	}
	// Deduct from sender
	key := sender.String()
	remaining, hasNeg := m.balances[key].SafeSub(amt...)
	if hasNeg {
		return fmt.Errorf("insufficient funds")
	}
	m.balances[key] = remaining
	m.moduleBalance = m.moduleBalance.Add(amt...)
	return nil
}

func (m *mockBankKeeper) SendCoinsFromModuleToAccount(_ context.Context, _ string, recipient sdk.AccAddress, amt sdk.Coins) error {
	if m.sendErr != nil {
		return m.sendErr
	}
	remaining, hasNeg := m.moduleBalance.SafeSub(amt...)
	if hasNeg {
		return fmt.Errorf("insufficient module funds")
	}
	m.moduleBalance = remaining
	key := recipient.String()
	m.balances[key] = m.balances[key].Add(amt...)
	return nil
}

func (m *mockBankKeeper) GetBalance(_ context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	coins, ok := m.balances[addr.String()]
	if !ok {
		return sdk.NewCoin(denom, math.ZeroInt())
	}
	for _, c := range coins {
		if c.Denom == denom {
			return c
		}
	}
	return sdk.NewCoin(denom, math.ZeroInt())
}

func (m *mockBankKeeper) SpendableCoins(_ context.Context, addr sdk.AccAddress) sdk.Coins {
	return m.balances[addr.String()]
}

// --- Test Helper ---

const authority = "syreen1governance"

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("syreen", "syreenpub")
	config.SetBech32PrefixForValidator("syreenvaloper", "syreenvaloperpub")
	config.SetBech32PrefixForConsensusNode("syreenvalcons", "syreenvalconspub")
	config.Seal()
}

func setupKeeper(t *testing.T) (*keeper.Keeper, sdk.Context, *mockBankKeeper) {
	t.Helper()

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	stateStore := cosmosstore.NewCommitMultiStore(db, log.NewNopLogger(), storemetrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	if err := stateStore.LoadLatestVersion(); err != nil {
		t.Fatal(err)
	}

	ctx := sdk.NewContext(stateStore, cmtproto.Header{Height: 1}, false, log.NewNopLogger())

	var storeService store.KVStoreService = runtime.NewKVStoreService(storeKey)

	bk := newMockBankKeeper()
	k := keeper.NewKeeper(
		codec.NewProtoCodec(nil),
		storeService,
		mockAccountKeeper{},
		bk,
		authority,
	)

	// Init genesis with default params
	k.InitGenesis(ctx, *types.DefaultGenesis())

	return k, ctx, bk
}

// Test addresses (valid bech32 with syreen prefix)
const (
	creator1 = "syreen1z4tljj0tqa8wrh5p856e3vu5afj5ayrxrxqrqq"
	payer1   = "syreen1mety9pctf2dngpawmwcp6unyw4a6rgkkkazpqz"
	payee1   = "syreen1pgzph9rze2j2xxavx4n7pdhxlkgsq7razhqd7n"
	payee2   = "syreen1j5a77099lqrf33jwccdv7lywkuwcvx66gujjez"
)

// ============================================================
// TESTS
// ============================================================

func TestCreateInvoice(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	msg := &types.MsgCreateInvoice{
		Creator:     creator1,
		BookingID:   "BK-001",
		Payee:       payee1,
		Amount:      sdk.NewCoin("usyreen", math.NewInt(1000000)),
		Description: "Hotel booking",
		Currency:    "USD",
	}

	id, err := k.CreateInvoice(ctx, msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "INV-000001" {
		t.Fatalf("expected INV-000001, got %s", id)
	}

	invoice, found := k.GetInvoice(ctx, id)
	if !found {
		t.Fatal("invoice not found after creation")
	}
	if invoice.BookingID != "BK-001" {
		t.Fatalf("expected booking BK-001, got %s", invoice.BookingID)
	}
	if invoice.Status != types.InvoiceStatusPending {
		t.Fatalf("expected pending status, got %s", invoice.Status)
	}
	if invoice.Creator != creator1 {
		t.Fatalf("expected creator %s, got %s", creator1, invoice.Creator)
	}
}

func TestCreateInvoiceFeeCalculation(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	// Amount = 1,000,000 usyreen
	// Platform fee = 2.5% = 25,000
	// Insurance fee = 0.5% = 5,000
	// Net = 97% = 970,000
	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-002",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(1000000)),
	}

	id, err := k.CreateInvoice(ctx, msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	invoice, _ := k.GetInvoice(ctx, id)

	if !invoice.PlatformFee.Amount.Equal(math.NewInt(25000)) {
		t.Fatalf("expected platform fee 25000, got %s", invoice.PlatformFee.Amount)
	}
	if !invoice.InsuranceFee.Amount.Equal(math.NewInt(5000)) {
		t.Fatalf("expected insurance fee 5000, got %s", invoice.InsuranceFee.Amount)
	}
	if !invoice.NetAmount.Amount.Equal(math.NewInt(970000)) {
		t.Fatalf("expected net amount 970000, got %s", invoice.NetAmount.Amount)
	}
}

func TestCreateInvoiceAutoIncrement(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	for i := 1; i <= 5; i++ {
		msg := &types.MsgCreateInvoice{
			Creator:   creator1,
			BookingID: fmt.Sprintf("BK-%03d", i),
			Payee:     payee1,
			Amount:    sdk.NewCoin("usyreen", math.NewInt(100000)),
		}
		id, err := k.CreateInvoice(ctx, msg)
		if err != nil {
			t.Fatalf("unexpected error on invoice %d: %v", i, err)
		}
		expected := fmt.Sprintf("INV-%06d", i)
		if id != expected {
			t.Fatalf("expected %s, got %s", expected, id)
		}
	}
}

func TestPayInvoice(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	// Fund payer
	payerAddr, _ := sdk.AccAddressFromBech32(payer1)
	bk.balances[payerAddr.String()] = sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10000000)))

	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-PAY",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(1000000)),
	}
	id, _ := k.CreateInvoice(ctx, msg)

	payMsg := &types.MsgPayInvoice{
		Payer:         payer1,
		InvoiceID:     id,
		PaymentMethod: "syr",
	}

	err := k.PayInvoice(ctx, payMsg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	invoice, _ := k.GetInvoice(ctx, id)
	if invoice.Status != types.InvoiceStatusPaid {
		t.Fatalf("expected paid status, got %s", invoice.Status)
	}
	if invoice.Payer != payer1 {
		t.Fatalf("expected payer %s, got %s", payer1, invoice.Payer)
	}
	if invoice.PaymentMethod != "syr" {
		t.Fatalf("expected payment method syr, got %s", invoice.PaymentMethod)
	}
}

func TestPayInvoiceEarnings(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	payerAddr, _ := sdk.AccAddressFromBech32(payer1)
	bk.balances[payerAddr.String()] = sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10000000)))

	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-EARN",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(1000000)),
	}
	id, _ := k.CreateInvoice(ctx, msg)

	payMsg := &types.MsgPayInvoice{
		Payer:         payer1,
		InvoiceID:     id,
		PaymentMethod: "usdc",
	}
	k.PayInvoice(ctx, payMsg)

	// Check payee earnings: net = 970000
	earnings, found := k.GetEarnings(ctx, payee1)
	if !found {
		t.Fatal("payee earnings not found")
	}
	expected := math.NewInt(970000)
	if !earnings.Available.AmountOf("usyreen").Equal(expected) {
		t.Fatalf("expected available earnings 970000, got %s", earnings.Available.AmountOf("usyreen"))
	}
	if !earnings.TotalEarned.AmountOf("usyreen").Equal(expected) {
		t.Fatalf("expected total earned 970000, got %s", earnings.TotalEarned.AmountOf("usyreen"))
	}
}

func TestDoublePay(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	payerAddr, _ := sdk.AccAddressFromBech32(payer1)
	bk.balances[payerAddr.String()] = sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10000000)))

	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-DBL",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(100000)),
	}
	id, _ := k.CreateInvoice(ctx, msg)

	payMsg := &types.MsgPayInvoice{
		Payer:         payer1,
		InvoiceID:     id,
		PaymentMethod: "syr",
	}
	k.PayInvoice(ctx, payMsg)

	// Try to pay again
	err := k.PayInvoice(ctx, payMsg)
	if err == nil {
		t.Fatal("expected error on double pay")
	}
	if err != types.ErrInvoiceAlreadyPaid {
		t.Fatalf("expected ErrInvoiceAlreadyPaid, got %v", err)
	}
}

func TestPayNonExistentInvoice(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	payMsg := &types.MsgPayInvoice{
		Payer:         payer1,
		InvoiceID:     "INV-NONEXISTENT",
		PaymentMethod: "syr",
	}

	err := k.PayInvoice(ctx, payMsg)
	if err == nil {
		t.Fatal("expected error for non-existent invoice")
	}
	if err != types.ErrInvoiceNotFound {
		t.Fatalf("expected ErrInvoiceNotFound, got %v", err)
	}
}

func TestPayInsufficientFunds(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	// Payer has no funds
	payerAddr, _ := sdk.AccAddressFromBech32(payer1)
	bk.balances[payerAddr.String()] = sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(100))) // way too little

	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-BROKE",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(1000000)),
	}
	id, _ := k.CreateInvoice(ctx, msg)

	payMsg := &types.MsgPayInvoice{
		Payer:         payer1,
		InvoiceID:     id,
		PaymentMethod: "syr",
	}

	err := k.PayInvoice(ctx, payMsg)
	if err == nil {
		t.Fatal("expected error for insufficient funds")
	}
}

func TestRefundPayment(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	payerAddr, _ := sdk.AccAddressFromBech32(payer1)
	bk.balances[payerAddr.String()] = sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10000000)))

	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-REFUND",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(1000000)),
	}
	id, _ := k.CreateInvoice(ctx, msg)

	// Pay it
	payMsg := &types.MsgPayInvoice{Payer: payer1, InvoiceID: id, PaymentMethod: "syr"}
	k.PayInvoice(ctx, payMsg)

	// Refund it
	refundMsg := &types.MsgRefundPayment{
		Authority: authority,
		InvoiceID: id,
		Reason:    "customer requested refund",
	}
	err := k.RefundPayment(ctx, refundMsg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	invoice, _ := k.GetInvoice(ctx, id)
	if invoice.Status != types.InvoiceStatusRefunded {
		t.Fatalf("expected refunded status, got %s", invoice.Status)
	}
}

func TestRefundUnauthorized(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	payerAddr, _ := sdk.AccAddressFromBech32(payer1)
	bk.balances[payerAddr.String()] = sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10000000)))

	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-UNAUTH",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(100000)),
	}
	id, _ := k.CreateInvoice(ctx, msg)

	payMsg := &types.MsgPayInvoice{Payer: payer1, InvoiceID: id, PaymentMethod: "syr"}
	k.PayInvoice(ctx, payMsg)

	// Try refund with non-authority
	refundMsg := &types.MsgRefundPayment{
		Authority: payer1, // not the authority
		InvoiceID: id,
		Reason:    "unauthorized attempt",
	}
	err := k.RefundPayment(ctx, refundMsg)
	if err == nil {
		t.Fatal("expected unauthorized error")
	}
	if err != types.ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestRefundPendingInvoice(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-REFPEND",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(100000)),
	}
	id, _ := k.CreateInvoice(ctx, msg)

	// Try to refund a pending (not paid) invoice
	refundMsg := &types.MsgRefundPayment{
		Authority: authority,
		InvoiceID: id,
		Reason:    "trying to refund pending",
	}
	err := k.RefundPayment(ctx, refundMsg)
	if err == nil {
		t.Fatal("expected error refunding pending invoice")
	}
	if err != types.ErrInvalidInvoiceStatus {
		t.Fatalf("expected ErrInvalidInvoiceStatus, got %v", err)
	}
}

func TestCancelInvoice(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-CANCEL",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(100000)),
	}
	id, _ := k.CreateInvoice(ctx, msg)

	cancelMsg := &types.MsgCancelInvoice{
		Creator:   creator1,
		InvoiceID: id,
	}
	err := k.CancelInvoice(ctx, cancelMsg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	invoice, _ := k.GetInvoice(ctx, id)
	if invoice.Status != types.InvoiceStatusCancelled {
		t.Fatalf("expected cancelled status, got %s", invoice.Status)
	}
}

func TestCancelInvoiceNotCreator(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-CANCNC",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(100000)),
	}
	id, _ := k.CreateInvoice(ctx, msg)

	cancelMsg := &types.MsgCancelInvoice{
		Creator:   payer1, // not the creator
		InvoiceID: id,
	}
	err := k.CancelInvoice(ctx, cancelMsg)
	if err == nil {
		t.Fatal("expected error for non-creator cancel")
	}
	if err != types.ErrNotInvoiceCreator {
		t.Fatalf("expected ErrNotInvoiceCreator, got %v", err)
	}
}

func TestCancelPaidInvoice(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	payerAddr, _ := sdk.AccAddressFromBech32(payer1)
	bk.balances[payerAddr.String()] = sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10000000)))

	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-CANPAID",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(100000)),
	}
	id, _ := k.CreateInvoice(ctx, msg)

	payMsg := &types.MsgPayInvoice{Payer: payer1, InvoiceID: id, PaymentMethod: "syr"}
	k.PayInvoice(ctx, payMsg)

	cancelMsg := &types.MsgCancelInvoice{Creator: creator1, InvoiceID: id}
	err := k.CancelInvoice(ctx, cancelMsg)
	if err == nil {
		t.Fatal("expected error cancelling paid invoice")
	}
	if err != types.ErrInvalidInvoiceStatus {
		t.Fatalf("expected ErrInvalidInvoiceStatus, got %v", err)
	}
}

func TestPayCancelledInvoice(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	payerAddr, _ := sdk.AccAddressFromBech32(payer1)
	bk.balances[payerAddr.String()] = sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10000000)))

	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-PAYCANC",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(100000)),
	}
	id, _ := k.CreateInvoice(ctx, msg)

	// Cancel
	cancelMsg := &types.MsgCancelInvoice{Creator: creator1, InvoiceID: id}
	k.CancelInvoice(ctx, cancelMsg)

	// Try to pay
	payMsg := &types.MsgPayInvoice{Payer: payer1, InvoiceID: id, PaymentMethod: "syr"}
	err := k.PayInvoice(ctx, payMsg)
	if err == nil {
		t.Fatal("expected error paying cancelled invoice")
	}
	if err != types.ErrInvoiceCancelled {
		t.Fatalf("expected ErrInvoiceCancelled, got %v", err)
	}
}

func TestExpireInvoices(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	// Set short expiry for testing
	params := k.GetParams(ctx)
	params.InvoiceExpiryBlocks = 10
	k.SetParams(ctx, params)

	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-EXPIRE",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(100000)),
	}
	id, _ := k.CreateInvoice(ctx, msg)

	invoice, _ := k.GetInvoice(ctx, id)
	if invoice.Status != types.InvoiceStatusPending {
		t.Fatalf("expected pending, got %s", invoice.Status)
	}

	// Advance to past expiry
	ctx = ctx.WithBlockHeight(invoice.ExpiresAtBlock + 1)
	k.ExpireInvoices(ctx)

	invoice, _ = k.GetInvoice(ctx, id)
	if invoice.Status != types.InvoiceStatusExpired {
		t.Fatalf("expected expired, got %s", invoice.Status)
	}
}

func TestPayExpiredInvoice(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	payerAddr, _ := sdk.AccAddressFromBech32(payer1)
	bk.balances[payerAddr.String()] = sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10000000)))

	params := k.GetParams(ctx)
	params.InvoiceExpiryBlocks = 5
	k.SetParams(ctx, params)

	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-PAYEXP",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(100000)),
	}
	id, _ := k.CreateInvoice(ctx, msg)

	// Advance past expiry
	invoice, _ := k.GetInvoice(ctx, id)
	ctx = ctx.WithBlockHeight(invoice.ExpiresAtBlock + 1)

	payMsg := &types.MsgPayInvoice{Payer: payer1, InvoiceID: id, PaymentMethod: "syr"}
	err := k.PayInvoice(ctx, payMsg)
	if err == nil {
		t.Fatal("expected error paying expired invoice")
	}
	if err != types.ErrInvoiceExpired {
		t.Fatalf("expected ErrInvoiceExpired, got %v", err)
	}
}

func TestSetExchangeRate(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	msg := &types.MsgSetExchangeRate{
		Authority: authority,
		FromDenom: "usyreen",
		ToDenom:   "uusdc",
		Rate:      "0.001",
	}

	err := k.SetExchangeRateByMsg(ctx, msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rate, found := k.GetExchangeRate(ctx, "usyreen", "uusdc")
	if !found {
		t.Fatal("exchange rate not found")
	}
	expectedRate, _ := math.LegacyNewDecFromStr("0.001")
	if !rate.Rate.Equal(expectedRate) {
		t.Fatalf("expected rate 0.001, got %s", rate.Rate)
	}
	if rate.UpdatedBy != authority {
		t.Fatalf("expected updatedBy %s, got %s", authority, rate.UpdatedBy)
	}
}

func TestSetExchangeRateUnauthorized(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	msg := &types.MsgSetExchangeRate{
		Authority: payer1, // not authority
		FromDenom: "usyreen",
		ToDenom:   "uusdc",
		Rate:      "0.001",
	}

	err := k.SetExchangeRateByMsg(ctx, msg)
	if err == nil {
		t.Fatal("expected unauthorized error")
	}
	if err != types.ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestWithdrawEarnings(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	payerAddr, _ := sdk.AccAddressFromBech32(payer1)
	bk.balances[payerAddr.String()] = sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10000000)))

	// Create and pay invoice to generate earnings
	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-WITHDRAW",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(1000000)),
	}
	id, _ := k.CreateInvoice(ctx, msg)

	payMsg := &types.MsgPayInvoice{Payer: payer1, InvoiceID: id, PaymentMethod: "syr"}
	k.PayInvoice(ctx, payMsg)

	// Payee should have 970000 in earnings
	earnings, found := k.GetEarnings(ctx, payee1)
	if !found {
		t.Fatal("earnings not found")
	}
	if !earnings.Available.AmountOf("usyreen").Equal(math.NewInt(970000)) {
		t.Fatalf("expected 970000 available, got %s", earnings.Available.AmountOf("usyreen"))
	}

	// Withdraw 500000
	withdrawMsg := &types.MsgWithdrawEarnings{
		Address: payee1,
		Amount:  sdk.NewCoin("usyreen", math.NewInt(500000)),
	}
	err := k.WithdrawEarnings(ctx, withdrawMsg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	earnings, _ = k.GetEarnings(ctx, payee1)
	if !earnings.Available.AmountOf("usyreen").Equal(math.NewInt(470000)) {
		t.Fatalf("expected 470000 remaining, got %s", earnings.Available.AmountOf("usyreen"))
	}
	if !earnings.TotalWithdrawn.AmountOf("usyreen").Equal(math.NewInt(500000)) {
		t.Fatalf("expected 500000 withdrawn, got %s", earnings.TotalWithdrawn.AmountOf("usyreen"))
	}
}

func TestWithdrawEarningsNoEarnings(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	withdrawMsg := &types.MsgWithdrawEarnings{
		Address: payee2,
		Amount:  sdk.NewCoin("usyreen", math.NewInt(100)),
	}
	err := k.WithdrawEarnings(ctx, withdrawMsg)
	if err == nil {
		t.Fatal("expected error for no earnings")
	}
	if err != types.ErrNoEarningsAvailable {
		t.Fatalf("expected ErrNoEarningsAvailable, got %v", err)
	}
}

func TestWithdrawExceedsAvailable(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	payerAddr, _ := sdk.AccAddressFromBech32(payer1)
	bk.balances[payerAddr.String()] = sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10000000)))

	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-EXCEED",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(100000)),
	}
	id, _ := k.CreateInvoice(ctx, msg)

	payMsg := &types.MsgPayInvoice{Payer: payer1, InvoiceID: id, PaymentMethod: "syr"}
	k.PayInvoice(ctx, payMsg)

	// Try to withdraw more than available
	withdrawMsg := &types.MsgWithdrawEarnings{
		Address: payee1,
		Amount:  sdk.NewCoin("usyreen", math.NewInt(99999999)),
	}
	err := k.WithdrawEarnings(ctx, withdrawMsg)
	if err == nil {
		t.Fatal("expected error for exceeding available")
	}
	if err != types.ErrNoEarningsAvailable {
		t.Fatalf("expected ErrNoEarningsAvailable, got %v", err)
	}
}

func TestPaymentsDisabled(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	params := k.GetParams(ctx)
	params.EnablePayments = false
	k.SetParams(ctx, params)

	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-DISABLED",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(100000)),
	}
	_, err := k.CreateInvoice(ctx, msg)
	if err == nil {
		t.Fatal("expected error when payments disabled")
	}
	if err != types.ErrPaymentsDisabled {
		t.Fatalf("expected ErrPaymentsDisabled, got %v", err)
	}
}

func TestInvoicesByBooking(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	for i := 0; i < 3; i++ {
		msg := &types.MsgCreateInvoice{
			Creator:   creator1,
			BookingID: "BK-MULTI",
			Payee:     payee1,
			Amount:    sdk.NewCoin("usyreen", math.NewInt(100000)),
		}
		k.CreateInvoice(ctx, msg)
	}
	// Different booking
	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-OTHER",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(100000)),
	}
	k.CreateInvoice(ctx, msg)

	invoices := k.GetInvoicesByBooking(ctx, "BK-MULTI")
	if len(invoices) != 3 {
		t.Fatalf("expected 3 invoices, got %d", len(invoices))
	}
}

func TestInvoicesByPayer(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	payerAddr, _ := sdk.AccAddressFromBech32(payer1)
	bk.balances[payerAddr.String()] = sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(100000000)))

	for i := 0; i < 2; i++ {
		msg := &types.MsgCreateInvoice{
			Creator:   creator1,
			BookingID: fmt.Sprintf("BK-PAYER-%d", i),
			Payee:     payee1,
			Amount:    sdk.NewCoin("usyreen", math.NewInt(100000)),
		}
		id, _ := k.CreateInvoice(ctx, msg)
		payMsg := &types.MsgPayInvoice{Payer: payer1, InvoiceID: id, PaymentMethod: "syr"}
		k.PayInvoice(ctx, payMsg)
	}

	invoices := k.GetInvoicesByPayer(ctx, payer1)
	if len(invoices) != 2 {
		t.Fatalf("expected 2 invoices, got %d", len(invoices))
	}
}

func TestFeeCalculationSmallAmount(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	// Small amount: 100 usyreen
	// Platform fee = 2.5% = 2 (truncated from 2.5)
	// Insurance fee = 0.5% = 0 (truncated from 0.5)
	// Net = 100 - 2 - 0 = 98
	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-SMALL",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(100)),
	}

	id, _ := k.CreateInvoice(ctx, msg)
	invoice, _ := k.GetInvoice(ctx, id)

	if !invoice.PlatformFee.Amount.Equal(math.NewInt(2)) {
		t.Fatalf("expected platform fee 2, got %s", invoice.PlatformFee.Amount)
	}
	if !invoice.InsuranceFee.Amount.Equal(math.NewInt(0)) {
		t.Fatalf("expected insurance fee 0, got %s", invoice.InsuranceFee.Amount)
	}
	if !invoice.NetAmount.Amount.Equal(math.NewInt(98)) {
		t.Fatalf("expected net 98, got %s", invoice.NetAmount.Amount)
	}
}

func TestFeeCalculationLargeAmount(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	// 100,000,000 usyreen (100 SYR)
	// Platform fee = 2.5% = 2,500,000
	// Insurance fee = 0.5% = 500,000
	// Net = 97,000,000
	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-LARGE",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(100000000)),
	}

	id, _ := k.CreateInvoice(ctx, msg)
	invoice, _ := k.GetInvoice(ctx, id)

	if !invoice.PlatformFee.Amount.Equal(math.NewInt(2500000)) {
		t.Fatalf("expected platform fee 2500000, got %s", invoice.PlatformFee.Amount)
	}
	if !invoice.InsuranceFee.Amount.Equal(math.NewInt(500000)) {
		t.Fatalf("expected insurance fee 500000, got %s", invoice.InsuranceFee.Amount)
	}
	if !invoice.NetAmount.Amount.Equal(math.NewInt(97000000)) {
		t.Fatalf("expected net 97000000, got %s", invoice.NetAmount.Amount)
	}
}

func TestGenesisExportImport(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	// Create some data
	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-GEN",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(100000)),
	}
	k.CreateInvoice(ctx, msg)

	rateMsg := &types.MsgSetExchangeRate{
		Authority: authority,
		FromDenom: "usyreen",
		ToDenom:   "uusdc",
		Rate:      "0.5",
	}
	k.SetExchangeRateByMsg(ctx, rateMsg)

	// Export
	gs := k.ExportGenesis(ctx)
	if len(gs.Invoices) != 1 {
		t.Fatalf("expected 1 invoice in genesis, got %d", len(gs.Invoices))
	}
	if len(gs.ExchangeRates) != 1 {
		t.Fatalf("expected 1 exchange rate in genesis, got %d", len(gs.ExchangeRates))
	}
	if gs.InvoiceCount != 1 {
		t.Fatalf("expected invoice count 1, got %d", gs.InvoiceCount)
	}

	// Re-import into fresh keeper
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	stateStore := cosmosstore.NewCommitMultiStore(db, log.NewNopLogger(), storemetrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.LoadLatestVersion()

	ctx2 := sdk.NewContext(stateStore, cmtproto.Header{Height: 100}, false, log.NewNopLogger())
	var storeService store.KVStoreService = runtime.NewKVStoreService(storeKey)
	k2 := keeper.NewKeeper(codec.NewProtoCodec(nil), storeService, mockAccountKeeper{}, newMockBankKeeper(), authority)
	k2.InitGenesis(ctx2, *gs)

	// Verify data survived
	inv, found := k2.GetInvoice(ctx2, "INV-000001")
	if !found {
		t.Fatal("invoice not found after re-import")
	}
	if inv.BookingID != "BK-GEN" {
		t.Fatalf("expected BK-GEN, got %s", inv.BookingID)
	}

	rate, found := k2.GetExchangeRate(ctx2, "usyreen", "uusdc")
	if !found {
		t.Fatal("exchange rate not found after re-import")
	}
	expectedRate, _ := math.LegacyNewDecFromStr("0.5")
	if !rate.Rate.Equal(expectedRate) {
		t.Fatalf("expected rate 0.5, got %s", rate.Rate)
	}
}

func TestParamsValidation(t *testing.T) {
	params := types.DefaultParams()
	if err := params.Validate(); err != nil {
		t.Fatalf("default params should be valid: %v", err)
	}

	// Negative fee
	badParams := params
	badParams.PlatformFeePercent = math.LegacyNewDec(-1)
	if err := badParams.Validate(); err == nil {
		t.Fatal("expected error for negative fee")
	}

	// Fee over 100
	badParams2 := params
	badParams2.PlatformFeePercent = math.LegacyNewDec(101)
	if err := badParams2.Validate(); err == nil {
		t.Fatal("expected error for fee > 100")
	}

	// Combined fees over 100
	badParams3 := params
	badParams3.PlatformFeePercent = math.LegacyNewDec(60)
	badParams3.InsuranceFeePercent = math.LegacyNewDec(50)
	if err := badParams3.Validate(); err == nil {
		t.Fatal("expected error for combined fees > 100")
	}

	// Zero expiry
	badParams4 := params
	badParams4.InvoiceExpiryBlocks = 0
	if err := badParams4.Validate(); err == nil {
		t.Fatal("expected error for zero expiry blocks")
	}
}

func TestMsgValidateBasic(t *testing.T) {
	// Valid create invoice
	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-001",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(1000000)),
	}
	if err := msg.ValidateBasic(); err != nil {
		t.Fatalf("valid msg should pass: %v", err)
	}

	// Invalid creator address
	badMsg := &types.MsgCreateInvoice{
		Creator:   "invalid",
		BookingID: "BK-001",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(1000000)),
	}
	if err := badMsg.ValidateBasic(); err == nil {
		t.Fatal("expected error for invalid creator")
	}

	// Empty booking ID
	badMsg2 := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(1000000)),
	}
	if err := badMsg2.ValidateBasic(); err == nil {
		t.Fatal("expected error for empty booking ID")
	}

	// Zero amount
	badMsg3 := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-001",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(0)),
	}
	if err := badMsg3.ValidateBasic(); err == nil {
		t.Fatal("expected error for zero amount")
	}
}

func TestMsgPayInvoiceValidateBasic(t *testing.T) {
	// Valid
	msg := &types.MsgPayInvoice{
		Payer:         payer1,
		InvoiceID:     "INV-000001",
		PaymentMethod: "syr",
	}
	if err := msg.ValidateBasic(); err != nil {
		t.Fatalf("valid msg should pass: %v", err)
	}

	// Invalid payment method
	badMsg := &types.MsgPayInvoice{
		Payer:         payer1,
		InvoiceID:     "INV-000001",
		PaymentMethod: "invalid",
	}
	if err := badMsg.ValidateBasic(); err == nil {
		t.Fatal("expected error for invalid payment method")
	}
}

func TestMsgSetExchangeRateValidateBasic(t *testing.T) {
	// Valid
	msg := &types.MsgSetExchangeRate{
		Authority: creator1,
		FromDenom: "usyreen",
		ToDenom:   "uusdc",
		Rate:      "1.5",
	}
	if err := msg.ValidateBasic(); err != nil {
		t.Fatalf("valid msg should pass: %v", err)
	}

	// Same denom
	badMsg := &types.MsgSetExchangeRate{
		Authority: creator1,
		FromDenom: "usyreen",
		ToDenom:   "usyreen",
		Rate:      "1.0",
	}
	if err := badMsg.ValidateBasic(); err == nil {
		t.Fatal("expected error for same denom")
	}

	// Negative rate
	badMsg2 := &types.MsgSetExchangeRate{
		Authority: creator1,
		FromDenom: "usyreen",
		ToDenom:   "uusdc",
		Rate:      "-1.0",
	}
	if err := badMsg2.ValidateBasic(); err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestMultipleEarnings(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	payerAddr, _ := sdk.AccAddressFromBech32(payer1)
	bk.balances[payerAddr.String()] = sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(100000000)))

	// Create and pay 3 invoices to same payee
	for i := 0; i < 3; i++ {
		msg := &types.MsgCreateInvoice{
			Creator:   creator1,
			BookingID: fmt.Sprintf("BK-ME-%d", i),
			Payee:     payee1,
			Amount:    sdk.NewCoin("usyreen", math.NewInt(1000000)),
		}
		id, _ := k.CreateInvoice(ctx, msg)
		payMsg := &types.MsgPayInvoice{Payer: payer1, InvoiceID: id, PaymentMethod: "syr"}
		k.PayInvoice(ctx, payMsg)
	}

	// Should have 3 * 970000 = 2910000
	earnings, found := k.GetEarnings(ctx, payee1)
	if !found {
		t.Fatal("earnings not found")
	}
	expected := math.NewInt(2910000)
	if !earnings.Available.AmountOf("usyreen").Equal(expected) {
		t.Fatalf("expected 2910000 available, got %s", earnings.Available.AmountOf("usyreen"))
	}
}

func TestExpiryDoesNotAffectPaidInvoices(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	payerAddr, _ := sdk.AccAddressFromBech32(payer1)
	bk.balances[payerAddr.String()] = sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10000000)))

	params := k.GetParams(ctx)
	params.InvoiceExpiryBlocks = 5
	k.SetParams(ctx, params)

	msg := &types.MsgCreateInvoice{
		Creator:   creator1,
		BookingID: "BK-NEXP",
		Payee:     payee1,
		Amount:    sdk.NewCoin("usyreen", math.NewInt(100000)),
	}
	id, _ := k.CreateInvoice(ctx, msg)

	// Pay it
	payMsg := &types.MsgPayInvoice{Payer: payer1, InvoiceID: id, PaymentMethod: "syr"}
	k.PayInvoice(ctx, payMsg)

	// Advance past expiry
	invoice, _ := k.GetInvoice(ctx, id)
	ctx = ctx.WithBlockHeight(invoice.ExpiresAtBlock + 10)
	k.ExpireInvoices(ctx)

	// Should still be paid, not expired
	invoice, _ = k.GetInvoice(ctx, id)
	if invoice.Status != types.InvoiceStatusPaid {
		t.Fatalf("expected paid status, got %s", invoice.Status)
	}
}

func TestGetAllInvoicesEmpty(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	invoices := k.GetAllInvoices(ctx)
	if invoices != nil && len(invoices) != 0 {
		t.Fatalf("expected empty invoices, got %d", len(invoices))
	}
}

func TestGetAllExchangeRatesEmpty(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	rates := k.GetAllExchangeRates(ctx)
	if rates != nil && len(rates) != 0 {
		t.Fatalf("expected empty rates, got %d", len(rates))
	}
}

func TestGetAllEarningsEmpty(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	earnings := k.GetAllEarnings(ctx)
	if earnings != nil && len(earnings) != 0 {
		t.Fatalf("expected empty earnings, got %d", len(earnings))
	}
}
