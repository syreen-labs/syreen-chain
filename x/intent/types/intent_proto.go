package types

// This file provides real protobuf wire-format support for the core intent Msg
// types (MsgSubmitIntent, MsgRegisterSolver, MsgDeregisterSolver,
// MsgSubmitSolution, MsgFulfillIntent, MsgSubmitChain, MsgCancelChain) so
// that they can be broadcast as Cosmos SDK transactions.
//
// The SDK's unknown-field reject path (codec/unknownproto) requires the Go
// type to implement Descriptor() ([]byte, []int), and the tx decoder must be
// able to Marshal/Unmarshal the message bytes on the wire.

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gogoproto "github.com/cosmos/gogoproto/proto"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

var intentTxFileDescriptorGzipped []byte

const (
	idxMsgSubmitIntent            = 0
	idxMsgSubmitIntentResponse    = 1
	idxMsgRegisterSolver          = 2
	idxMsgRegisterSolverResponse  = 3
	idxMsgDeregisterSolver        = 4
	idxMsgDeregisterSolverResponse = 5
	idxMsgSubmitSolution          = 6
	idxMsgSubmitSolutionResponse  = 7
	idxMsgFulfillIntent           = 8
	idxMsgFulfillIntentResponse   = 9
	idxMsgSubmitChain             = 10
	idxMsgSubmitChainResponse     = 11
	idxMsgCancelChain             = 12
	idxMsgCancelChainResponse     = 13
	idxMsgCancelIntent            = 14
	idxMsgCancelIntentResponse    = 15
)

func init() {
	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/intent/intent_tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.intent"),
		MessageType: []*descriptorpb.DescriptorProto{
			{ // 0: MsgSubmitIntent
				Name: sp("MsgSubmitIntent"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField(1, "creator"),
					stringField(2, "intent_type"),
					stringField(3, "body"),
					stringField(4, "max_fee"),
					stringField(5, "tip"),
					uint64Field(6, "expiry_blocks"),
				},
			},
			{Name: sp("MsgSubmitIntentResponse")}, // 1
			{ // 2: MsgRegisterSolver
				Name: sp("MsgRegisterSolver"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField(1, "address"),
					stringField(2, "moniker"),
					stringField(3, "stake_amount"),
				},
			},
			{Name: sp("MsgRegisterSolverResponse")}, // 3
			{ // 4: MsgDeregisterSolver
				Name: sp("MsgDeregisterSolver"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField(1, "address"),
				},
			},
			{Name: sp("MsgDeregisterSolverResponse")}, // 5
			{ // 6: MsgSubmitSolution
				Name: sp("MsgSubmitSolution"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField(1, "solver_addr"),
					stringField(2, "intent_id"),
					stringField(3, "execution_msgs"),
					stringField(4, "expected_outcome"),
				},
			},
			{Name: sp("MsgSubmitSolutionResponse")}, // 7
			{ // 8: MsgFulfillIntent
				Name: sp("MsgFulfillIntent"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField(1, "sender"),
					stringField(2, "solver_addr"),
					stringField(3, "intent_id"),
				},
			},
			{Name: sp("MsgFulfillIntentResponse")}, // 9
			{ // 10: MsgSubmitChain
				Name: sp("MsgSubmitChain"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField(1, "creator"),
					stringField(2, "steps"),
					stringField(3, "max_fee"),
					stringField(4, "tip"),
					uint64Field(5, "expiry_blocks"),
				},
			},
			{ // 11: MsgSubmitChainResponse
				Name: sp("MsgSubmitChainResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField(1, "chain_id"),
				},
			},
			{ // 12: MsgCancelChain
				Name: sp("MsgCancelChain"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField(1, "creator"),
					stringField(2, "chain_id"),
				},
			},
			{Name: sp("MsgCancelChainResponse")}, // 13
			{ // 14: MsgCancelIntent
				Name: sp("MsgCancelIntent"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField(1, "creator"),
					stringField(2, "intent_id"),
				},
			},
			{Name: sp("MsgCancelIntentResponse")}, // 15
		},
	}

	gogoproto.RegisterType((*MsgSubmitIntent)(nil), "syreen.intent.MsgSubmitIntent")
	gogoproto.RegisterType((*MsgSubmitIntentResponse)(nil), "syreen.intent.MsgSubmitIntentResponse")
	gogoproto.RegisterType((*MsgRegisterSolver)(nil), "syreen.intent.MsgRegisterSolver")
	gogoproto.RegisterType((*MsgRegisterSolverResponse)(nil), "syreen.intent.MsgRegisterSolverResponse")
	gogoproto.RegisterType((*MsgDeregisterSolver)(nil), "syreen.intent.MsgDeregisterSolver")
	gogoproto.RegisterType((*MsgDeregisterSolverResponse)(nil), "syreen.intent.MsgDeregisterSolverResponse")
	gogoproto.RegisterType((*MsgSubmitSolution)(nil), "syreen.intent.MsgSubmitSolution")
	gogoproto.RegisterType((*MsgSubmitSolutionResponse)(nil), "syreen.intent.MsgSubmitSolutionResponse")
	gogoproto.RegisterType((*MsgFulfillIntent)(nil), "syreen.intent.MsgFulfillIntent")
	gogoproto.RegisterType((*MsgFulfillIntentResponse)(nil), "syreen.intent.MsgFulfillIntentResponse")
	gogoproto.RegisterType((*MsgSubmitChain)(nil), "syreen.intent.MsgSubmitChain")
	gogoproto.RegisterType((*MsgSubmitChainResponse)(nil), "syreen.intent.MsgSubmitChainResponse")
	gogoproto.RegisterType((*MsgCancelChain)(nil), "syreen.intent.MsgCancelChain")
	gogoproto.RegisterType((*MsgCancelChainResponse)(nil), "syreen.intent.MsgCancelChainResponse")
	gogoproto.RegisterType((*MsgCancelIntent)(nil), "syreen.intent.MsgCancelIntent")
	gogoproto.RegisterType((*MsgCancelIntentResponse)(nil), "syreen.intent.MsgCancelIntentResponse")

	raw, err := proto.Marshal(fd)
	if err != nil {
		panic("intent: failed to marshal intent file descriptor: " + err.Error())
	}
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(raw); err != nil {
		panic("intent: failed to gzip intent file descriptor: " + err.Error())
	}
	if err := gz.Close(); err != nil {
		panic("intent: failed to close gzip writer for intent file descriptor: " + err.Error())
	}
	intentTxFileDescriptorGzipped = buf.Bytes()
}

// ---------------------------------------------------------------------------
// Descriptor() methods
// ---------------------------------------------------------------------------

func (*MsgSubmitIntent) Descriptor() ([]byte, []int) {
	return intentTxFileDescriptorGzipped, []int{idxMsgSubmitIntent}
}

func (*MsgRegisterSolver) Descriptor() ([]byte, []int) {
	return intentTxFileDescriptorGzipped, []int{idxMsgRegisterSolver}
}

func (*MsgDeregisterSolver) Descriptor() ([]byte, []int) {
	return intentTxFileDescriptorGzipped, []int{idxMsgDeregisterSolver}
}

func (*MsgSubmitSolution) Descriptor() ([]byte, []int) {
	return intentTxFileDescriptorGzipped, []int{idxMsgSubmitSolution}
}

func (*MsgFulfillIntent) Descriptor() ([]byte, []int) {
	return intentTxFileDescriptorGzipped, []int{idxMsgFulfillIntent}
}

func (*MsgSubmitChain) Descriptor() ([]byte, []int) {
	return intentTxFileDescriptorGzipped, []int{idxMsgSubmitChain}
}

func (*MsgSubmitChainResponse) Descriptor() ([]byte, []int) {
	return intentTxFileDescriptorGzipped, []int{idxMsgSubmitChainResponse}
}

func (*MsgCancelChain) Descriptor() ([]byte, []int) {
	return intentTxFileDescriptorGzipped, []int{idxMsgCancelChain}
}

func (*MsgCancelIntent) Descriptor() ([]byte, []int) {
	return intentTxFileDescriptorGzipped, []int{idxMsgCancelIntent}
}

func (*MsgCancelIntentResponse) Descriptor() ([]byte, []int) {
	return intentTxFileDescriptorGzipped, []int{idxMsgCancelIntentResponse}
}

// Response types
func (*MsgSubmitIntentResponse) Descriptor() ([]byte, []int) {
	return intentTxFileDescriptorGzipped, []int{idxMsgSubmitIntentResponse}
}

func (*MsgRegisterSolverResponse) Descriptor() ([]byte, []int) {
	return intentTxFileDescriptorGzipped, []int{idxMsgRegisterSolverResponse}
}

func (*MsgDeregisterSolverResponse) Descriptor() ([]byte, []int) {
	return intentTxFileDescriptorGzipped, []int{idxMsgDeregisterSolverResponse}
}

func (*MsgSubmitSolutionResponse) Descriptor() ([]byte, []int) {
	return intentTxFileDescriptorGzipped, []int{idxMsgSubmitSolutionResponse}
}

func (*MsgFulfillIntentResponse) Descriptor() ([]byte, []int) {
	return intentTxFileDescriptorGzipped, []int{idxMsgFulfillIntentResponse}
}

func (*MsgCancelChainResponse) Descriptor() ([]byte, []int) {
	return intentTxFileDescriptorGzipped, []int{idxMsgCancelChainResponse}
}

// ---------------------------------------------------------------------------
// Marshal / Unmarshal — MsgRegisterSolver
// ---------------------------------------------------------------------------

func (m *MsgRegisterSolver) Marshal() ([]byte, error) {
	var b []byte
	b = appendString(b, 1, m.Address)
	b = appendString(b, 2, m.Moniker)
	if m.StakeAmount.IsPositive() {
		b = appendString(b, 3, m.StakeAmount.String())
	}
	return b, nil
}

func (m *MsgRegisterSolver) MarshalTo(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	copy(dAtA, bz)
	return len(bz), nil
}

func (m *MsgRegisterSolver) Size() int {
	bz, _ := m.Marshal()
	return len(bz)
}

func (m *MsgRegisterSolver) Unmarshal(data []byte) error {
	*m = MsgRegisterSolver{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		switch num {
		case 1:
			if typ != protowire.BytesType {
				return fmt.Errorf("MsgRegisterSolver: wrong wire type for address")
			}
			v, nn := protowire.ConsumeString(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			m.Address = v
			data = data[nn:]
		case 2:
			if typ != protowire.BytesType {
				return fmt.Errorf("MsgRegisterSolver: wrong wire type for moniker")
			}
			v, nn := protowire.ConsumeString(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			m.Moniker = v
			data = data[nn:]
		case 3:
			if typ != protowire.BytesType {
				return fmt.Errorf("MsgRegisterSolver: wrong wire type for stake_amount")
			}
			v, nn := protowire.ConsumeString(data)
			if nn < 0 {
				return protowire.ParseError(nn)
			}
			coin, err := sdk.ParseCoinNormalized(v)
			if err != nil {
				return fmt.Errorf("MsgRegisterSolver: invalid stake_amount %q: %w", v, err)
			}
			m.StakeAmount = coin
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

// ---------------------------------------------------------------------------
// Marshal / Unmarshal — MsgRegisterSolverResponse
// ---------------------------------------------------------------------------

type MsgRegisterSolverResponse struct{}

func (m *MsgRegisterSolverResponse) ProtoMessage()           {}
func (m *MsgRegisterSolverResponse) Reset()                  { *m = MsgRegisterSolverResponse{} }
func (m *MsgRegisterSolverResponse) String() string          { return "register_solver_response" }
func (m *MsgRegisterSolverResponse) XXX_MessageName() string { return "syreen.intent.MsgRegisterSolverResponse" }

func (m *MsgRegisterSolverResponse) Marshal() ([]byte, error) { return nil, nil }
func (m *MsgRegisterSolverResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgRegisterSolverResponse) Size() int { return 0 }
func (m *MsgRegisterSolverResponse) Unmarshal(data []byte) error {
	*m = MsgRegisterSolverResponse{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 { return protowire.ParseError(n) }
		data = data[n:]
		nn, err := skipField(data, num, typ)
		if err != nil { return err }
		data = data[nn:]
	}
	return nil
}

// ---------------------------------------------------------------------------
// Marshal / Unmarshal — MsgDeregisterSolver
// ---------------------------------------------------------------------------

func (m *MsgDeregisterSolver) Marshal() ([]byte, error) {
	var b []byte
	b = appendString(b, 1, m.Address)
	return b, nil
}

func (m *MsgDeregisterSolver) MarshalTo(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil { return 0, err }
	copy(dAtA, bz)
	return len(bz), nil
}

func (m *MsgDeregisterSolver) Size() int {
	bz, _ := m.Marshal()
	return len(bz)
}

func (m *MsgDeregisterSolver) Unmarshal(data []byte) error {
	*m = MsgDeregisterSolver{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 { return protowire.ParseError(n) }
		data = data[n:]
		switch num {
		case 1:
			if typ != protowire.BytesType {
				return fmt.Errorf("MsgDeregisterSolver: wrong wire type for address")
			}
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.Address = v
			data = data[nn:]
		default:
			nn, err := skipField(data, num, typ)
			if err != nil { return err }
			data = data[nn:]
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Marshal / Unmarshal — MsgDeregisterSolverResponse
// ---------------------------------------------------------------------------

type MsgDeregisterSolverResponse struct{}

func (m *MsgDeregisterSolverResponse) ProtoMessage()           {}
func (m *MsgDeregisterSolverResponse) Reset()                  { *m = MsgDeregisterSolverResponse{} }
func (m *MsgDeregisterSolverResponse) String() string          { return "deregister_solver_response" }
func (m *MsgDeregisterSolverResponse) XXX_MessageName() string { return "syreen.intent.MsgDeregisterSolverResponse" }

func (m *MsgDeregisterSolverResponse) Marshal() ([]byte, error) { return nil, nil }
func (m *MsgDeregisterSolverResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgDeregisterSolverResponse) Size() int { return 0 }
func (m *MsgDeregisterSolverResponse) Unmarshal(data []byte) error {
	*m = MsgDeregisterSolverResponse{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 { return protowire.ParseError(n) }
		data = data[n:]
		nn, err := skipField(data, num, typ)
		if err != nil { return err }
		data = data[nn:]
	}
	return nil
}

// ---------------------------------------------------------------------------
// Marshal / Unmarshal — MsgSubmitIntent
// ---------------------------------------------------------------------------

func (m *MsgSubmitIntent) Marshal() ([]byte, error) {
	var b []byte
	b = appendString(b, 1, m.Creator)
	b = appendString(b, 2, m.IntentType)
	if len(m.Body) > 0 {
		b = appendString(b, 3, string(m.Body))
	}
	if len(m.MaxFee) > 0 {
		b = appendString(b, 4, m.MaxFee.String())
	}
	if len(m.Tip) > 0 {
		b = appendString(b, 5, m.Tip.String())
	}
	b = appendUint64(b, 6, m.ExpiryBlocks)
	return b, nil
}

func (m *MsgSubmitIntent) MarshalTo(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil { return 0, err }
	copy(dAtA, bz)
	return len(bz), nil
}

func (m *MsgSubmitIntent) Size() int {
	bz, _ := m.Marshal()
	return len(bz)
}

func (m *MsgSubmitIntent) Unmarshal(data []byte) error {
	*m = MsgSubmitIntent{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 { return protowire.ParseError(n) }
		data = data[n:]
		switch num {
		case 1:
			if typ != protowire.BytesType { return fmt.Errorf("MsgSubmitIntent: wrong wire type for creator") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.Creator = v
			data = data[nn:]
		case 2:
			if typ != protowire.BytesType { return fmt.Errorf("MsgSubmitIntent: wrong wire type for intent_type") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.IntentType = v
			data = data[nn:]
		case 3:
			if typ != protowire.BytesType { return fmt.Errorf("MsgSubmitIntent: wrong wire type for body") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.Body = json.RawMessage(v)
			data = data[nn:]
		case 4:
			if typ != protowire.BytesType { return fmt.Errorf("MsgSubmitIntent: wrong wire type for max_fee") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			coins, err := sdk.ParseCoinsNormalized(v)
			if err != nil { return fmt.Errorf("MsgSubmitIntent: invalid max_fee: %w", err) }
			m.MaxFee = coins
			data = data[nn:]
		case 5:
			if typ != protowire.BytesType { return fmt.Errorf("MsgSubmitIntent: wrong wire type for tip") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			coins, err := sdk.ParseCoinsNormalized(v)
			if err != nil { return fmt.Errorf("MsgSubmitIntent: invalid tip: %w", err) }
			m.Tip = coins
			data = data[nn:]
		case 6:
			if typ != protowire.VarintType { return fmt.Errorf("MsgSubmitIntent: wrong wire type for expiry_blocks") }
			v, nn := protowire.ConsumeVarint(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.ExpiryBlocks = v
			data = data[nn:]
		default:
			nn, err := skipField(data, num, typ)
			if err != nil { return err }
			data = data[nn:]
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgSubmitIntentResponse
// ---------------------------------------------------------------------------

type MsgSubmitIntentResponse struct {
	IntentID string `protobuf:"bytes,1,opt,name=intent_id,proto3" json:"intent_id"`
}

func (m *MsgSubmitIntentResponse) ProtoMessage()           {}
func (m *MsgSubmitIntentResponse) Reset()                  { *m = MsgSubmitIntentResponse{} }
func (m *MsgSubmitIntentResponse) String() string          { return m.IntentID }
func (m *MsgSubmitIntentResponse) XXX_MessageName() string { return "syreen.intent.MsgSubmitIntentResponse" }

func (m *MsgSubmitIntentResponse) Marshal() ([]byte, error) {
	var b []byte
	b = appendString(b, 1, m.IntentID)
	return b, nil
}
func (m *MsgSubmitIntentResponse) MarshalTo(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil { return 0, err }
	copy(dAtA, bz)
	return len(bz), nil
}
func (m *MsgSubmitIntentResponse) Size() int { bz, _ := m.Marshal(); return len(bz) }
func (m *MsgSubmitIntentResponse) Unmarshal(data []byte) error {
	*m = MsgSubmitIntentResponse{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 { return protowire.ParseError(n) }
		data = data[n:]
		switch num {
		case 1:
			if typ != protowire.BytesType { return fmt.Errorf("wrong wire type") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.IntentID = v
			data = data[nn:]
		default:
			nn, err := skipField(data, num, typ)
			if err != nil { return err }
			data = data[nn:]
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Marshal / Unmarshal — MsgSubmitSolution
// ---------------------------------------------------------------------------

func (m *MsgSubmitSolution) Marshal() ([]byte, error) {
	var b []byte
	b = appendString(b, 1, m.SolverAddr)
	b = appendString(b, 2, m.IntentID)
	for _, em := range m.ExecutionMsgs {
		b = appendString(b, 3, string(em))
	}
	if len(m.ExpectedOutcome) > 0 {
		b = appendString(b, 4, string(m.ExpectedOutcome))
	}
	return b, nil
}

func (m *MsgSubmitSolution) MarshalTo(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil { return 0, err }
	copy(dAtA, bz)
	return len(bz), nil
}

func (m *MsgSubmitSolution) Size() int { bz, _ := m.Marshal(); return len(bz) }

func (m *MsgSubmitSolution) Unmarshal(data []byte) error {
	*m = MsgSubmitSolution{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 { return protowire.ParseError(n) }
		data = data[n:]
		switch num {
		case 1:
			if typ != protowire.BytesType { return fmt.Errorf("MsgSubmitSolution: wrong wire type for solver_addr") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.SolverAddr = v
			data = data[nn:]
		case 2:
			if typ != protowire.BytesType { return fmt.Errorf("MsgSubmitSolution: wrong wire type for intent_id") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.IntentID = v
			data = data[nn:]
		case 3:
			if typ != protowire.BytesType { return fmt.Errorf("MsgSubmitSolution: wrong wire type for execution_msgs") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.ExecutionMsgs = append(m.ExecutionMsgs, json.RawMessage(v))
			data = data[nn:]
		case 4:
			if typ != protowire.BytesType { return fmt.Errorf("MsgSubmitSolution: wrong wire type for expected_outcome") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.ExpectedOutcome = json.RawMessage(v)
			data = data[nn:]
		default:
			nn, err := skipField(data, num, typ)
			if err != nil { return err }
			data = data[nn:]
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgSubmitSolutionResponse
// ---------------------------------------------------------------------------

type MsgSubmitSolutionResponse struct{}

func (m *MsgSubmitSolutionResponse) ProtoMessage()           {}
func (m *MsgSubmitSolutionResponse) Reset()                  { *m = MsgSubmitSolutionResponse{} }
func (m *MsgSubmitSolutionResponse) String() string          { return "submit_solution_response" }
func (m *MsgSubmitSolutionResponse) XXX_MessageName() string { return "syreen.intent.MsgSubmitSolutionResponse" }

func (m *MsgSubmitSolutionResponse) Marshal() ([]byte, error) { return nil, nil }
func (m *MsgSubmitSolutionResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgSubmitSolutionResponse) Size() int { return 0 }
func (m *MsgSubmitSolutionResponse) Unmarshal(data []byte) error {
	*m = MsgSubmitSolutionResponse{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 { return protowire.ParseError(n) }
		data = data[n:]
		nn, err := skipField(data, num, typ)
		if err != nil { return err }
		data = data[nn:]
	}
	return nil
}

// ---------------------------------------------------------------------------
// Marshal / Unmarshal — MsgFulfillIntent
// ---------------------------------------------------------------------------

func (m *MsgFulfillIntent) Marshal() ([]byte, error) {
	var b []byte
	b = appendString(b, 1, m.Sender)
	b = appendString(b, 2, m.SolverAddr)
	b = appendString(b, 3, m.IntentID)
	return b, nil
}

func (m *MsgFulfillIntent) MarshalTo(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil { return 0, err }
	copy(dAtA, bz)
	return len(bz), nil
}

func (m *MsgFulfillIntent) Size() int { bz, _ := m.Marshal(); return len(bz) }

func (m *MsgFulfillIntent) Unmarshal(data []byte) error {
	*m = MsgFulfillIntent{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 { return protowire.ParseError(n) }
		data = data[n:]
		switch num {
		case 1:
			if typ != protowire.BytesType { return fmt.Errorf("MsgFulfillIntent: wrong wire type for sender") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.Sender = v
			data = data[nn:]
		case 2:
			if typ != protowire.BytesType { return fmt.Errorf("MsgFulfillIntent: wrong wire type for solver_addr") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.SolverAddr = v
			data = data[nn:]
		case 3:
			if typ != protowire.BytesType { return fmt.Errorf("MsgFulfillIntent: wrong wire type for intent_id") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.IntentID = v
			data = data[nn:]
		default:
			nn, err := skipField(data, num, typ)
			if err != nil { return err }
			data = data[nn:]
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgFulfillIntentResponse
// ---------------------------------------------------------------------------

type MsgFulfillIntentResponse struct{}

func (m *MsgFulfillIntentResponse) ProtoMessage()           {}
func (m *MsgFulfillIntentResponse) Reset()                  { *m = MsgFulfillIntentResponse{} }
func (m *MsgFulfillIntentResponse) String() string          { return "fulfill_intent_response" }
func (m *MsgFulfillIntentResponse) XXX_MessageName() string { return "syreen.intent.MsgFulfillIntentResponse" }

func (m *MsgFulfillIntentResponse) Marshal() ([]byte, error) { return nil, nil }
func (m *MsgFulfillIntentResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgFulfillIntentResponse) Size() int { return 0 }
func (m *MsgFulfillIntentResponse) Unmarshal(data []byte) error {
	*m = MsgFulfillIntentResponse{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 { return protowire.ParseError(n) }
		data = data[n:]
		nn, err := skipField(data, num, typ)
		if err != nil { return err }
		data = data[nn:]
	}
	return nil
}

// ---------------------------------------------------------------------------
// Marshal / Unmarshal — MsgSubmitChain
// ---------------------------------------------------------------------------

func (m *MsgSubmitChain) Marshal() ([]byte, error) {
	var b []byte
	b = appendString(b, 1, m.Creator)
	for _, step := range m.Steps {
		stepJSON, err := json.Marshal(step)
		if err != nil { return nil, fmt.Errorf("MsgSubmitChain: failed to marshal step: %w", err) }
		b = appendString(b, 2, string(stepJSON))
	}
	if len(m.MaxFee) > 0 {
		b = appendString(b, 3, m.MaxFee.String())
	}
	if len(m.Tip) > 0 {
		b = appendString(b, 4, m.Tip.String())
	}
	b = appendUint64(b, 5, m.ExpiryBlocks)
	return b, nil
}

func (m *MsgSubmitChain) MarshalTo(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil { return 0, err }
	copy(dAtA, bz)
	return len(bz), nil
}

func (m *MsgSubmitChain) Size() int { bz, _ := m.Marshal(); return len(bz) }

func (m *MsgSubmitChain) Unmarshal(data []byte) error {
	*m = MsgSubmitChain{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 { return protowire.ParseError(n) }
		data = data[n:]
		switch num {
		case 1:
			if typ != protowire.BytesType { return fmt.Errorf("MsgSubmitChain: wrong wire type for creator") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.Creator = v
			data = data[nn:]
		case 2:
			if typ != protowire.BytesType { return fmt.Errorf("MsgSubmitChain: wrong wire type for steps") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			var step ChainStep
			if err := json.Unmarshal([]byte(v), &step); err != nil {
				return fmt.Errorf("MsgSubmitChain: invalid step: %w", err)
			}
			m.Steps = append(m.Steps, step)
			data = data[nn:]
		case 3:
			if typ != protowire.BytesType { return fmt.Errorf("MsgSubmitChain: wrong wire type for max_fee") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			coins, err := sdk.ParseCoinsNormalized(v)
			if err != nil { return fmt.Errorf("MsgSubmitChain: invalid max_fee: %w", err) }
			m.MaxFee = coins
			data = data[nn:]
		case 4:
			if typ != protowire.BytesType { return fmt.Errorf("MsgSubmitChain: wrong wire type for tip") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			coins, err := sdk.ParseCoinsNormalized(v)
			if err != nil { return fmt.Errorf("MsgSubmitChain: invalid tip: %w", err) }
			m.Tip = coins
			data = data[nn:]
		case 5:
			if typ != protowire.VarintType { return fmt.Errorf("MsgSubmitChain: wrong wire type for expiry_blocks") }
			v, nn := protowire.ConsumeVarint(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.ExpiryBlocks = v
			data = data[nn:]
		default:
			nn, err := skipField(data, num, typ)
			if err != nil { return err }
			data = data[nn:]
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgSubmitChainResponse — Marshal/Unmarshal
// ---------------------------------------------------------------------------

func (m *MsgSubmitChainResponse) Marshal() ([]byte, error) {
	var b []byte
	b = appendString(b, 1, m.ChainID)
	return b, nil
}

func (m *MsgSubmitChainResponse) MarshalTo(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil { return 0, err }
	copy(dAtA, bz)
	return len(bz), nil
}

func (m *MsgSubmitChainResponse) Size() int { bz, _ := m.Marshal(); return len(bz) }

func (m *MsgSubmitChainResponse) Unmarshal(data []byte) error {
	*m = MsgSubmitChainResponse{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 { return protowire.ParseError(n) }
		data = data[n:]
		switch num {
		case 1:
			if typ != protowire.BytesType { return fmt.Errorf("wrong wire type") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.ChainID = v
			data = data[nn:]
		default:
			nn, err := skipField(data, num, typ)
			if err != nil { return err }
			data = data[nn:]
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Marshal / Unmarshal — MsgCancelChain
// ---------------------------------------------------------------------------

func (m *MsgCancelChain) Marshal() ([]byte, error) {
	var b []byte
	b = appendString(b, 1, m.Creator)
	b = appendString(b, 2, m.ChainID)
	return b, nil
}

func (m *MsgCancelChain) MarshalTo(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil { return 0, err }
	copy(dAtA, bz)
	return len(bz), nil
}

func (m *MsgCancelChain) Size() int { bz, _ := m.Marshal(); return len(bz) }

func (m *MsgCancelChain) Unmarshal(data []byte) error {
	*m = MsgCancelChain{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 { return protowire.ParseError(n) }
		data = data[n:]
		switch num {
		case 1:
			if typ != protowire.BytesType { return fmt.Errorf("MsgCancelChain: wrong wire type for creator") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.Creator = v
			data = data[nn:]
		case 2:
			if typ != protowire.BytesType { return fmt.Errorf("MsgCancelChain: wrong wire type for chain_id") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.ChainID = v
			data = data[nn:]
		default:
			nn, err := skipField(data, num, typ)
			if err != nil { return err }
			data = data[nn:]
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgCancelChainResponse — Marshal/Unmarshal
// ---------------------------------------------------------------------------

func (m *MsgCancelChainResponse) Marshal() ([]byte, error) { return nil, nil }
func (m *MsgCancelChainResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgCancelChainResponse) Size() int { return 0 }
func (m *MsgCancelChainResponse) Unmarshal(data []byte) error {
	*m = MsgCancelChainResponse{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 { return protowire.ParseError(n) }
		data = data[n:]
		nn, err := skipField(data, num, typ)
		if err != nil { return err }
		data = data[nn:]
	}
	return nil
}

// ---------------------------------------------------------------------------
// Marshal / Unmarshal — MsgCancelIntent
// ---------------------------------------------------------------------------

func (m *MsgCancelIntent) Marshal() ([]byte, error) {
	var b []byte
	b = appendString(b, 1, m.Creator)
	b = appendString(b, 2, m.IntentID)
	return b, nil
}

func (m *MsgCancelIntent) MarshalTo(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil { return 0, err }
	copy(dAtA, bz)
	return len(bz), nil
}

func (m *MsgCancelIntent) Size() int { bz, _ := m.Marshal(); return len(bz) }

func (m *MsgCancelIntent) Unmarshal(data []byte) error {
	*m = MsgCancelIntent{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 { return protowire.ParseError(n) }
		data = data[n:]
		switch num {
		case 1:
			if typ != protowire.BytesType { return fmt.Errorf("MsgCancelIntent: wrong wire type for creator") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.Creator = v
			data = data[nn:]
		case 2:
			if typ != protowire.BytesType { return fmt.Errorf("MsgCancelIntent: wrong wire type for intent_id") }
			v, nn := protowire.ConsumeString(data)
			if nn < 0 { return protowire.ParseError(nn) }
			m.IntentID = v
			data = data[nn:]
		default:
			nn, err := skipField(data, num, typ)
			if err != nil { return err }
			data = data[nn:]
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgCancelIntentResponse — Marshal/Unmarshal
// ---------------------------------------------------------------------------

func (m *MsgCancelIntentResponse) Marshal() ([]byte, error) { return nil, nil }
func (m *MsgCancelIntentResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgCancelIntentResponse) Size() int { return 0 }
func (m *MsgCancelIntentResponse) Unmarshal(data []byte) error {
	*m = MsgCancelIntentResponse{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 { return protowire.ParseError(n) }
		data = data[n:]
		nn, err := skipField(data, num, typ)
		if err != nil { return err }
		data = data[nn:]
	}
	return nil
}
