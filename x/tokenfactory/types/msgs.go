package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MaxSubdenomLength  = 44
	TypeMsgCreateDenom = "create_denom"
	TypeMsgMint        = "tf_mint"
	TypeMsgBurn        = "tf_burn"
	TypeMsgChangeAdmin = "change_admin"
)

// GetDenom returns the full denom string for a factory token
func GetDenom(creator, subdenom string) string {
	return fmt.Sprintf("factory/%s/%s", creator, subdenom)
}

// DeconstructDenom splits a factory denom into creator and subdenom
func DeconstructDenom(denom string) (creator string, subdenom string, err error) {
	if !strings.HasPrefix(denom, "factory/") {
		return "", "", fmt.Errorf("denom %s is not a factory denom", denom)
	}
	parts := strings.SplitN(denom, "/", 3)
	if len(parts) != 3 {
		return "", "", fmt.Errorf("invalid factory denom: %s", denom)
	}
	return parts[1], parts[2], nil
}

// MsgCreateDenom creates a new factory denom
type MsgCreateDenom struct {
	Sender   string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	Subdenom string `protobuf:"bytes,2,opt,name=subdenom,proto3" json:"subdenom"`
}

func (m *MsgCreateDenom) ProtoMessage()             {}
func (m *MsgCreateDenom) Reset()                    { *m = MsgCreateDenom{} }
func (m *MsgCreateDenom) String() string             { return fmt.Sprintf("create_denom: %s/%s", m.Sender, m.Subdenom) }
func (m *MsgCreateDenom) XXX_MessageName() string    { return "syreen.tokenfactory.MsgCreateDenom" }

func (m *MsgCreateDenom) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return ErrInvalidCreator
	}
	if len(m.Subdenom) == 0 {
		return ErrSubdenomTooShort
	}
	if len(m.Subdenom) > MaxSubdenomLength {
		return ErrSubdenomTooLong
	}
	// Validate subdenom characters: alphanumeric, underscore, hyphen, period only
	for _, c := range m.Subdenom {
		if !isValidSubdenomChar(c) {
			return fmt.Errorf("invalid character '%c' in subdenom: only alphanumeric, underscore, hyphen, and period allowed", c)
		}
	}
	return nil
}

func isValidSubdenomChar(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-' || c == '.'
}

// MsgMint mints tokens from a factory denom
type MsgMint struct {
	Sender string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	Amount sdk.Coin `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount"`
	MintTo string   `protobuf:"bytes,3,opt,name=mint_to_address,json=mintToAddress,proto3" json:"mint_to_address"`
}

func (m *MsgMint) ProtoMessage()             {}
func (m *MsgMint) Reset()                    { *m = MsgMint{} }
func (m *MsgMint) String() string             { return fmt.Sprintf("tf_mint: %s %s", m.Sender, m.Amount) }
func (m *MsgMint) XXX_MessageName() string    { return "syreen.tokenfactory.MsgMint" }

func (m *MsgMint) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return ErrInvalidCreator
	}
	if !m.Amount.IsValid() || m.Amount.IsZero() {
		return fmt.Errorf("invalid mint amount")
	}
	if m.MintTo != "" {
		_, err := sdk.AccAddressFromBech32(m.MintTo)
		if err != nil {
			return fmt.Errorf("invalid mint_to address")
		}
	}
	return nil
}

// MsgBurn burns tokens of a factory denom
type MsgBurn struct {
	Sender   string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	Amount   sdk.Coin `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount"`
	BurnFrom string   `protobuf:"bytes,3,opt,name=burn_from_address,json=burnFromAddress,proto3" json:"burn_from_address"`
}

func (m *MsgBurn) ProtoMessage()             {}
func (m *MsgBurn) Reset()                    { *m = MsgBurn{} }
func (m *MsgBurn) String() string             { return fmt.Sprintf("tf_burn: %s %s", m.Sender, m.Amount) }
func (m *MsgBurn) XXX_MessageName() string    { return "syreen.tokenfactory.MsgBurn" }

func (m *MsgBurn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return ErrInvalidCreator
	}
	if !m.Amount.IsValid() || m.Amount.IsZero() {
		return fmt.Errorf("invalid burn amount")
	}
	// BurnFrom must be empty (burn from sender) or a valid address
	if m.BurnFrom != "" {
		if _, err := sdk.AccAddressFromBech32(m.BurnFrom); err != nil {
			return fmt.Errorf("invalid burn_from address: %w", err)
		}
	}
	return nil
}

// MsgChangeAdmin changes the admin of a factory denom
type MsgChangeAdmin struct {
	Sender   string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	Denom    string `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom"`
	NewAdmin string `protobuf:"bytes,3,opt,name=new_admin,json=newAdmin,proto3" json:"new_admin"`
}

func (m *MsgChangeAdmin) ProtoMessage()             {}
func (m *MsgChangeAdmin) Reset()                    { *m = MsgChangeAdmin{} }
func (m *MsgChangeAdmin) String() string             { return fmt.Sprintf("change_admin: %s %s", m.Denom, m.NewAdmin) }
func (m *MsgChangeAdmin) XXX_MessageName() string    { return "syreen.tokenfactory.MsgChangeAdmin" }

func (m *MsgChangeAdmin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return ErrInvalidCreator
	}
	_, err = sdk.AccAddressFromBech32(m.NewAdmin)
	if err != nil {
		return ErrInvalidAdmin
	}
	return nil
}

