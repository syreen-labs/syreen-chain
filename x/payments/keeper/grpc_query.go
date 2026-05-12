package keeper

import (
	"context"
	"fmt"

	"syreen/x/payments/types"
)

type queryServer struct {
	Keeper
}

func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

var _ types.QueryServer = queryServer{}

func (q queryServer) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params := q.Keeper.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

func (q queryServer) Invoice(ctx context.Context, req *types.QueryInvoiceRequest) (*types.QueryInvoiceResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: request cannot be nil")
	}
	if req.InvoiceID == "" {
		return nil, fmt.Errorf("invalid request: invoice_id cannot be empty")
	}
	invoice, found := q.Keeper.GetInvoice(ctx, req.InvoiceID)
	return &types.QueryInvoiceResponse{
		Invoice: invoice,
		Found:   found,
	}, nil
}

func (q queryServer) InvoicesByBooking(ctx context.Context, req *types.QueryInvoicesByBookingRequest) (*types.QueryInvoicesByBookingResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: request cannot be nil")
	}
	if req.BookingID == "" {
		return nil, fmt.Errorf("invalid request: booking_id cannot be empty")
	}
	invoices := q.Keeper.GetInvoicesByBooking(ctx, req.BookingID)
	return &types.QueryInvoicesByBookingResponse{Invoices: invoices}, nil
}

func (q queryServer) InvoicesByPayer(ctx context.Context, req *types.QueryInvoicesByPayerRequest) (*types.QueryInvoicesByPayerResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: request cannot be nil")
	}
	if req.Payer == "" {
		return nil, fmt.Errorf("invalid request: payer cannot be empty")
	}
	invoices := q.Keeper.GetInvoicesByPayer(ctx, req.Payer)
	return &types.QueryInvoicesByPayerResponse{Invoices: invoices}, nil
}

func (q queryServer) ExchangeRate(ctx context.Context, req *types.QueryExchangeRateRequest) (*types.QueryExchangeRateResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: request cannot be nil")
	}
	if req.FromDenom == "" || req.ToDenom == "" {
		return nil, fmt.Errorf("invalid request: from_denom and to_denom are required")
	}
	rate, found := q.Keeper.GetExchangeRate(ctx, req.FromDenom, req.ToDenom)
	return &types.QueryExchangeRateResponse{
		ExchangeRate: rate,
		Found:        found,
	}, nil
}

func (q queryServer) Earnings(ctx context.Context, req *types.QueryEarningsRequest) (*types.QueryEarningsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: request cannot be nil")
	}
	if req.Address == "" {
		return nil, fmt.Errorf("invalid request: address cannot be empty")
	}
	earnings, found := q.Keeper.GetEarnings(ctx, req.Address)
	return &types.QueryEarningsResponse{
		Earnings: earnings,
		Found:    found,
	}, nil
}
