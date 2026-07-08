package types

// This file provides real protobuf wire-format support for the strategy Msg
// types so that they can be broadcast as Cosmos SDK transactions.
//
// Background: the rest of x/intent uses hand-rolled "fake proto" structs that
// only implement ProtoMessage/Reset/String/XXX_MessageName. That is enough for
// internal registration but NOT for broadcasting via the tx endpoint, because
// the SDK's unknown-field reject path (codec/unknownproto) requires the Go
// type to implement Descriptor() ([]byte, []int), and the tx decoder must be
// able to Marshal/Unmarshal the message bytes on the wire.
//
// To keep the fix surgical, we only add this plumbing for MsgCreateStrategy
// and MsgCancelStrategy (plus their responses) — the messages the website
// actually broadcasts. The other intent Msgs remain as they were.

import (
	"bytes"
	"compress/gzip"
	"fmt"

	"cosmossdk.io/math"

	gogoproto "github.com/cosmos/gogoproto/proto"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// strategyTxFileDescriptorGzipped holds the gzipped serialized
// FileDescriptorProto for the four strategy messages. It is returned from
// each message's Descriptor() method.
var strategyTxFileDescriptorGzipped []byte

// Message indices inside strategyTxFileDescriptorGzipped (must match the
// order of MessageType below).
const (
	idxMsgCreateStrategy         = 0
	idxMsgCreateStrategyResponse = 1
	idxMsgCancelStrategy         = 2
	idxMsgCancelStrategyResponse = 3
)

func init() {
	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/intent/strategy_tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.intent"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: sp("MsgCreateStrategy"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField(1, "creator"),
					stringField(2, "template_name"),
					uint64Field(3, "pool_id"),
					stringField(4, "input_denom"),
					stringField(5, "output_denom"),
					stringField(6, "total_budget"),
					stringField(7, "risk_level"),
					uint64Field(8, "expiry_blocks"),
					uint64Field(9, "num_executions"),
					uint64Field(10, "interval_blocks"),
					uint64Field(11, "grid_levels"),
				},
			},
			{
				Name: sp("MsgCreateStrategyResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField(1, "strategy_id"),
					stringField(2, "chain_id"),
				},
			},
			{
				Name: sp("MsgCancelStrategy"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField(1, "creator"),
					stringField(2, "strategy_id"),
				},
			},
			{
				Name: sp("MsgCancelStrategyResponse"),
			},
		},
	}

	// NOTE: We intentionally do NOT register this FileDescriptorProto with
	// protoregistry.GlobalFiles because proto_register.go already registers
	// "syreen/intent/tx.proto" containing the same message names, which would
	// cause a fully-qualified-name conflict. The gzipped bytes are only used
	// as the return value of Descriptor(), which the SDK's unknownproto path
	// consumes directly.

	// Register a Go type mapping so that unknownproto.protoMessageForTypeName
	// can find these types by their fully-qualified name.
	gogoproto.RegisterType((*MsgCreateStrategy)(nil), "syreen.intent.MsgCreateStrategy")
	gogoproto.RegisterType((*MsgCreateStrategyResponse)(nil), "syreen.intent.MsgCreateStrategyResponse")
	gogoproto.RegisterType((*MsgCancelStrategy)(nil), "syreen.intent.MsgCancelStrategy")
	gogoproto.RegisterType((*MsgCancelStrategyResponse)(nil), "syreen.intent.MsgCancelStrategyResponse")

	raw, err := proto.Marshal(fd)
	if err != nil {
		panic("intent: failed to marshal strategy file descriptor: " + err.Error())
	}
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(raw); err != nil {
		panic("intent: failed to gzip strategy file descriptor: " + err.Error())
	}
	if err := gz.Close(); err != nil {
		panic("intent: failed to close gzip writer for strategy file descriptor: " + err.Error())
	}
	strategyTxFileDescriptorGzipped = buf.Bytes()
}

func stringField(num int32, name string) *descriptorpb.FieldDescriptorProto {
	t := descriptorpb.FieldDescriptorProto_TYPE_STRING
	l := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	n := num
	return &descriptorpb.FieldDescriptorProto{
		Name:     sp(name),
		Number:   &n,
		Type:     &t,
		Label:    &l,
		JsonName: sp(toLowerCamel(name)),
	}
}

func uint64Field(num int32, name string) *descriptorpb.FieldDescriptorProto {
	t := descriptorpb.FieldDescriptorProto_TYPE_UINT64
	l := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	n := num
	return &descriptorpb.FieldDescriptorProto{
		Name:     sp(name),
		Number:   &n,
		Type:     &t,
		Label:    &l,
		JsonName: sp(toLowerCamel(name)),
	}
}

func bytesField(num int32, name string) *descriptorpb.FieldDescriptorProto {
	t := descriptorpb.FieldDescriptorProto_TYPE_BYTES
	l := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	n := num
	return &descriptorpb.FieldDescriptorProto{
		Name:     sp(name),
		Number:   &n,
		Type:     &t,
		Label:    &l,
		JsonName: sp(toLowerCamel(name)),
	}
}

// toLowerCamel converts snake_case to lowerCamelCase (creator -> creator,
// template_name -> templateName).
func toLowerCamel(s string) string {
	out := make([]byte, 0, len(s))
	upper := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' {
			upper = true
			continue
		}
		if upper {
			if c >= 'a' && c <= 'z' {
				c -= 32
			}
			upper = false
		}
		out = append(out, c)
	}
	return string(out)
}

// ---------------------------------------------------------------------------
// Descriptor() methods — required by codec/unknownproto so that the SDK tx
// decoder can walk the message's field types when validating unknown fields.
// ---------------------------------------------------------------------------

func (*MsgCreateStrategy) Descriptor() ([]byte, []int) {
	return strategyTxFileDescriptorGzipped, []int{idxMsgCreateStrategy}
}

func (*MsgCreateStrategyResponse) Descriptor() ([]byte, []int) {
	return strategyTxFileDescriptorGzipped, []int{idxMsgCreateStrategyResponse}
}

func (*MsgCancelStrategy) Descriptor() ([]byte, []int) {
	return strategyTxFileDescriptorGzipped, []int{idxMsgCancelStrategy}
}

func (*MsgCancelStrategyResponse) Descriptor() ([]byte, []int) {
	return strategyTxFileDescriptorGzipped, []int{idxMsgCancelStrategyResponse}
}

// ---------------------------------------------------------------------------
// Marshal / Unmarshal — raw protobuf wire encoding. Field numbers MUST match
// the FileDescriptorProto above and the TypeScript encoder on the website.
// ---------------------------------------------------------------------------

func appendString(b []byte, num protowire.Number, v string) []byte {
	if v == "" {
		return b
	}
	b = protowire.AppendTag(b, num, protowire.BytesType)
	b = protowire.AppendString(b, v)
	return b
}

func appendUint64(b []byte, num protowire.Number, v uint64) []byte {
	if v == 0 {
		return b
	}
	b = protowire.AppendTag(b, num, protowire.VarintType)
	b = protowire.AppendVarint(b, v)
	return b
}

func skipField(data []byte, num protowire.Number, typ protowire.Type) (int, error) {
	n := protowire.ConsumeFieldValue(num, typ, data)
	if n < 0 {
		return 0, protowire.ParseError(n)
	}
	return n, nil
}

// --- MsgCreateStrategy ---

func (m *MsgCreateStrategy) Marshal() ([]byte, error) {
	var b []byte
	b = appendString(b, 1, m.Creator)
	b = appendString(b, 2, m.TemplateName)
	b = appendUint64(b, 3, m.PoolID)
	b = appendString(b, 4, m.InputDenom)
	b = appendString(b, 5, m.OutputDenom)
	// TotalBudget is encoded as its decimal string representation. Zero/nil
	// budgets still emit the field so the server sees an explicit value.
	if !m.TotalBudget.IsNil() {
		b = appendString(b, 6, m.TotalBudget.String())
	}
	b = appendString(b, 7, string(m.RiskLevel))
	b = appendUint64(b, 8, m.ExpiryBlocks)
	b = appendUint64(b, 9, m.NumExecutions)
	b = appendUint64(b, 10, m.IntervalBlocks)
	b = appendUint64(b, 11, m.GridLevels)
	return b, nil
}

func (m *MsgCreateStrategy) Unmarshal(data []byte) error {
	*m = MsgCreateStrategy{TotalBudget: math.ZeroInt()}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		switch num {
		case 1:
			if typ != protowire.BytesType {
				return fmt.Errorf("MsgCreateStrategy: wrong wire type for creator")
			}
			v, nn := protowire.ConsumeString(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			m.Creator = v
			data = data[nn:]
		case 2:
			if typ != protowire.BytesType {
				return fmt.Errorf("MsgCreateStrategy: wrong wire type for template_name")
			}
			v, nn := protowire.ConsumeString(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			m.TemplateName = v
			data = data[nn:]
		case 3:
			if typ != protowire.VarintType {
				return fmt.Errorf("MsgCreateStrategy: wrong wire type for pool_id")
			}
			v, nn := protowire.ConsumeVarint(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			m.PoolID = v
			data = data[nn:]
		case 4:
			if typ != protowire.BytesType {
				return fmt.Errorf("MsgCreateStrategy: wrong wire type for input_denom")
			}
			v, nn := protowire.ConsumeString(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			m.InputDenom = v
			data = data[nn:]
		case 5:
			if typ != protowire.BytesType {
				return fmt.Errorf("MsgCreateStrategy: wrong wire type for output_denom")
			}
			v, nn := protowire.ConsumeString(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			m.OutputDenom = v
			data = data[nn:]
		case 6:
			if typ != protowire.BytesType {
				return fmt.Errorf("MsgCreateStrategy: wrong wire type for total_budget")
			}
			v, nn := protowire.ConsumeString(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			i, ok := math.NewIntFromString(v)
			if !ok {
				return fmt.Errorf("MsgCreateStrategy: invalid total_budget %q", v)
			}
			m.TotalBudget = i
			data = data[nn:]
		case 7:
			if typ != protowire.BytesType {
				return fmt.Errorf("MsgCreateStrategy: wrong wire type for risk_level")
			}
			v, nn := protowire.ConsumeString(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			m.RiskLevel = RiskLevel(v)
			data = data[nn:]
		case 8:
			if typ != protowire.VarintType {
				return fmt.Errorf("MsgCreateStrategy: wrong wire type for expiry_blocks")
			}
			v, nn := protowire.ConsumeVarint(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			m.ExpiryBlocks = v
			data = data[nn:]
		case 9:
			if typ != protowire.VarintType {
				return fmt.Errorf("MsgCreateStrategy: wrong wire type for num_executions")
			}
			v, nn := protowire.ConsumeVarint(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			m.NumExecutions = v
			data = data[nn:]
		case 10:
			if typ != protowire.VarintType {
				return fmt.Errorf("MsgCreateStrategy: wrong wire type for interval_blocks")
			}
			v, nn := protowire.ConsumeVarint(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			m.IntervalBlocks = v
			data = data[nn:]
		case 11:
			if typ != protowire.VarintType {
				return fmt.Errorf("MsgCreateStrategy: wrong wire type for grid_levels")
			}
			v, nn := protowire.ConsumeVarint(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			m.GridLevels = v
			data = data[nn:]
		default:
			nn, err := skipField(data, num, typ)
			if err != nil {
				return err
			}
			data = data[nn:]
		}
	}
	return nil
}

// --- MsgCreateStrategyResponse ---

func (m *MsgCreateStrategyResponse) Marshal() ([]byte, error) {
	var b []byte
	b = appendString(b, 1, m.StrategyID)
	b = appendString(b, 2, m.ChainID)
	return b, nil
}

func (m *MsgCreateStrategyResponse) Unmarshal(data []byte) error {
	*m = MsgCreateStrategyResponse{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		switch num {
		case 1:
			if typ != protowire.BytesType {
				return fmt.Errorf("MsgCreateStrategyResponse: wrong wire type for strategy_id")
			}
			v, nn := protowire.ConsumeString(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			m.StrategyID = v
			data = data[nn:]
		case 2:
			if typ != protowire.BytesType {
				return fmt.Errorf("MsgCreateStrategyResponse: wrong wire type for chain_id")
			}
			v, nn := protowire.ConsumeString(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			m.ChainID = v
			data = data[nn:]
		default:
			nn, err := skipField(data, num, typ)
			if err != nil {
				return err
			}
			data = data[nn:]
		}
	}
	return nil
}

// --- MsgCancelStrategy ---

func (m *MsgCancelStrategy) Marshal() ([]byte, error) {
	var b []byte
	b = appendString(b, 1, m.Creator)
	b = appendString(b, 2, m.StrategyID)
	return b, nil
}

func (m *MsgCancelStrategy) Unmarshal(data []byte) error {
	*m = MsgCancelStrategy{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		switch num {
		case 1:
			if typ != protowire.BytesType {
				return fmt.Errorf("MsgCancelStrategy: wrong wire type for creator")
			}
			v, nn := protowire.ConsumeString(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			m.Creator = v
			data = data[nn:]
		case 2:
			if typ != protowire.BytesType {
				return fmt.Errorf("MsgCancelStrategy: wrong wire type for strategy_id")
			}
			v, nn := protowire.ConsumeString(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			m.StrategyID = v
			data = data[nn:]
		default:
			nn, err := skipField(data, num, typ)
			if err != nil {
				return err
			}
			data = data[nn:]
		}
	}
	return nil
}

// --- MsgCancelStrategyResponse ---

func (m *MsgCancelStrategyResponse) Marshal() ([]byte, error) {
	return nil, nil
}

func (m *MsgCancelStrategyResponse) Unmarshal(data []byte) error {
	*m = MsgCancelStrategyResponse{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		nn, err := skipField(data, num, typ)
		if err != nil {
			return err
		}
		data = data[nn:]
	}
	return nil
}
