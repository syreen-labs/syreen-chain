package types

import "context"

// MsgServer defines the tokenfactory msg server interface
type MsgServer interface {
	CreateDenom(context.Context, *MsgCreateDenom) (*MsgCreateDenomResponse, error)
	Mint(context.Context, *MsgMint) (*MsgMintResponse, error)
	Burn(context.Context, *MsgBurn) (*MsgBurnResponse, error)
	ChangeAdmin(context.Context, *MsgChangeAdmin) (*MsgChangeAdminResponse, error)
}

// Response types
type MsgCreateDenomResponse struct {
	NewTokenDenom string `protobuf:"bytes,1,opt,name=new_token_denom,json=newTokenDenom,proto3" json:"new_token_denom"`
}

func (m *MsgCreateDenomResponse) ProtoMessage()          {}
func (m *MsgCreateDenomResponse) Reset()                 { *m = MsgCreateDenomResponse{} }
func (m *MsgCreateDenomResponse) String() string         { return m.NewTokenDenom }
func (m *MsgCreateDenomResponse) XXX_MessageName() string { return "syreen.tokenfactory.MsgCreateDenomResponse" }

type MsgMintResponse struct{}
type MsgBurnResponse struct{}
type MsgChangeAdminResponse struct{}
