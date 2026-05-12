package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/payments/types"
)

type msgServer struct {
	*Keeper
}

func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = &msgServer{}

func (m msgServer) CreateInvoice(ctx context.Context, msg *types.MsgCreateInvoice) (*types.MsgCreateInvoiceResponse, error) {
	invoiceID, err := m.Keeper.CreateInvoice(ctx, msg)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"create_invoice",
		sdk.NewAttribute("creator", msg.Creator),
		sdk.NewAttribute("invoice_id", invoiceID),
		sdk.NewAttribute("booking_id", msg.BookingID),
		sdk.NewAttribute("payee", msg.Payee),
		sdk.NewAttribute("amount", msg.Amount.String()),
	))
	return &types.MsgCreateInvoiceResponse{InvoiceID: invoiceID}, nil
}

func (m msgServer) PayInvoice(ctx context.Context, msg *types.MsgPayInvoice) (*types.MsgPayInvoiceResponse, error) {
	if err := m.Keeper.PayInvoice(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"pay_invoice",
		sdk.NewAttribute("payer", msg.Payer),
		sdk.NewAttribute("invoice_id", msg.InvoiceID),
		sdk.NewAttribute("payment_method", msg.PaymentMethod),
	))
	return &types.MsgPayInvoiceResponse{}, nil
}

func (m msgServer) RefundPayment(ctx context.Context, msg *types.MsgRefundPayment) (*types.MsgRefundPaymentResponse, error) {
	// H-7 fix: defense-in-depth authority check at the msg_server boundary
	// (the keeper.RefundPayment also performs this check).
	if msg.Authority != m.Keeper.GetAuthority() {
		return nil, types.ErrUnauthorized
	}
	if err := m.Keeper.RefundPayment(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"refund_payment",
		sdk.NewAttribute("authority", msg.Authority),
		sdk.NewAttribute("invoice_id", msg.InvoiceID),
		sdk.NewAttribute("reason", msg.Reason),
	))
	return &types.MsgRefundPaymentResponse{}, nil
}

func (m msgServer) SetExchangeRate(ctx context.Context, msg *types.MsgSetExchangeRate) (*types.MsgSetExchangeRateResponse, error) {
	if err := m.Keeper.SetExchangeRateByMsg(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"set_exchange_rate",
		sdk.NewAttribute("authority", msg.Authority),
		sdk.NewAttribute("from_denom", msg.FromDenom),
		sdk.NewAttribute("to_denom", msg.ToDenom),
		sdk.NewAttribute("rate", msg.Rate),
	))
	return &types.MsgSetExchangeRateResponse{}, nil
}

func (m msgServer) WithdrawEarnings(ctx context.Context, msg *types.MsgWithdrawEarnings) (*types.MsgWithdrawEarningsResponse, error) {
	if err := m.Keeper.WithdrawEarnings(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"withdraw_earnings",
		sdk.NewAttribute("address", msg.Address),
		sdk.NewAttribute("amount", msg.Amount.String()),
	))
	return &types.MsgWithdrawEarningsResponse{}, nil
}

func (m msgServer) CancelInvoice(ctx context.Context, msg *types.MsgCancelInvoice) (*types.MsgCancelInvoiceResponse, error) {
	if err := m.Keeper.CancelInvoice(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"cancel_invoice",
		sdk.NewAttribute("creator", msg.Creator),
		sdk.NewAttribute("invoice_id", msg.InvoiceID),
	))
	return &types.MsgCancelInvoiceResponse{}, nil
}
