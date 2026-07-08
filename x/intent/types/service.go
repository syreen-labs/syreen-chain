package types

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protowire"
)

// Msg response types are defined in intent_proto.go (with Descriptor/Marshal/Unmarshal).

// --- Hand-rolled proto codec helpers ---
//
// The query request/response types below carry fields the reflection-based proto
// codec cannot represent: plain string fields (no protobuf struct tag) and domain
// structs/slices holding json.RawMessage (Intent.Body, Solution fields, etc.).
// Because these types have NO protobuf tags at all, reflection marshaling emits
// EMPTY output — so over gRPC/REST every field is silently dropped and clients see
// blank responses (this is the same class of bug as the July QuerySolutions fix).
// Each domain object is encoded as its JSON bytes in a length-delimited field,
// mirroring exactly how the keeper stores them (json.Marshal). Helpers below keep
// the per-type Marshal/Unmarshal small and uniform.

// appendString (proto3 string field, omitted when empty) is defined in
// strategy_proto.go and reused here.

// appendBool appends a proto3 bool field, omitting it when false (proto3 default).
func appendBool(b []byte, num protowire.Number, v bool) []byte {
	if !v {
		return b
	}
	b = protowire.AppendTag(b, num, protowire.VarintType)
	return protowire.AppendVarint(b, 1)
}

// appendJSON encodes v as JSON and appends it as a length-delimited bytes field.
func appendJSON(b []byte, num protowire.Number, v interface{}) ([]byte, error) {
	jb, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	b = protowire.AppendTag(b, num, protowire.BytesType)
	return protowire.AppendBytes(b, jb), nil
}

// sizedBuffer is the common MarshalToSizedBuffer body: marshal, then copy into the
// tail of data (gogoproto's contract).
func sizedBuffer(data []byte, marshal func() ([]byte, error)) (int, error) {
	b, err := marshal()
	if err != nil {
		return 0, err
	}
	n := len(b)
	copy(data[len(data)-n:], b)
	return n, nil
}

// marshalIntentSlice / unmarshalIntentSlice encode a []Intent as repeated
// length-delimited JSON bytes in field #1 — shared by the list-style responses.
func marshalIntentSlice(intents []Intent) ([]byte, error) {
	var b []byte
	var err error
	for i := range intents {
		if b, err = appendJSON(b, 1, &intents[i]); err != nil {
			return nil, err
		}
	}
	return b, nil
}

func unmarshalIntentSlice(data []byte) ([]Intent, error) {
	var out []Intent
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return nil, protowire.ParseError(n)
		}
		data = data[n:]
		if num == 1 && typ == protowire.BytesType {
			v, vn := protowire.ConsumeBytes(data)
			if vn < 0 {
				return nil, protowire.ParseError(vn)
			}
			data = data[vn:]
			var it Intent
			if err := json.Unmarshal(v, &it); err != nil {
				return nil, err
			}
			out = append(out, it)
		} else {
			vn := protowire.ConsumeFieldValue(num, typ, data)
			if vn < 0 {
				return nil, protowire.ParseError(vn)
			}
			data = data[vn:]
		}
	}
	return out, nil
}

// --- MsgServer Interface ---

type MsgServer interface {
	SubmitIntent(context.Context, *MsgSubmitIntent) (*MsgSubmitIntentResponse, error)
	RegisterSolver(context.Context, *MsgRegisterSolver) (*MsgRegisterSolverResponse, error)
	DeregisterSolver(context.Context, *MsgDeregisterSolver) (*MsgDeregisterSolverResponse, error)
	SubmitSolution(context.Context, *MsgSubmitSolution) (*MsgSubmitSolutionResponse, error)
	FulfillIntent(context.Context, *MsgFulfillIntent) (*MsgFulfillIntentResponse, error)
	CancelIntent(context.Context, *MsgCancelIntent) (*MsgCancelIntentResponse, error)
	SubmitChain(context.Context, *MsgSubmitChain) (*MsgSubmitChainResponse, error)
	CancelChain(context.Context, *MsgCancelChain) (*MsgCancelChainResponse, error)
	CreateStrategy(context.Context, *MsgCreateStrategy) (*MsgCreateStrategyResponse, error)
	CancelStrategy(context.Context, *MsgCancelStrategy) (*MsgCancelStrategyResponse, error)
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
}

// --- Query Request/Response Types ---

type QueryParamsRequest struct{}

func (m *QueryParamsRequest) ProtoMessage()           {}
func (m *QueryParamsRequest) Reset()                  { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string          { return "query_params_request" }
func (m *QueryParamsRequest) XXX_MessageName() string { return "syreen.intent.QueryParamsRequest" }

type QueryParamsResponse struct {
	Params Params `json:"params"`
}

func (m *QueryParamsResponse) ProtoMessage()           {}
func (m *QueryParamsResponse) Reset()                  { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string          { return fmt.Sprintf("params: %+v", m.Params) }
func (m *QueryParamsResponse) XXX_MessageName() string { return "syreen.intent.QueryParamsResponse" }

func (m *QueryParamsResponse) Marshal() ([]byte, error)              { return appendJSON(nil, 1, &m.Params) }
func (m *QueryParamsResponse) MarshalToSizedBuffer(d []byte) (int, error) { return sizedBuffer(d, m.Marshal) }
func (m *QueryParamsResponse) Size() int                            { b, _ := m.Marshal(); return len(b) }
func (m *QueryParamsResponse) Unmarshal(data []byte) error {
	*m = QueryParamsResponse{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		if num == 1 && typ == protowire.BytesType {
			v, vn := protowire.ConsumeBytes(data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
			if err := json.Unmarshal(v, &m.Params); err != nil {
				return err
			}
		} else {
			vn := protowire.ConsumeFieldValue(num, typ, data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
		}
	}
	return nil
}

type QueryIntentRequest struct {
	IntentID string `json:"intent_id"`
}

func (m *QueryIntentRequest) ProtoMessage()           {}
func (m *QueryIntentRequest) Reset()                  { *m = QueryIntentRequest{} }
func (m *QueryIntentRequest) String() string          { return fmt.Sprintf("query_intent: id=%s", m.IntentID) }
func (m *QueryIntentRequest) XXX_MessageName() string { return "syreen.intent.QueryIntentRequest" }

func (m *QueryIntentRequest) Marshal() ([]byte, error)              { return appendString(nil, 1, m.IntentID), nil }
func (m *QueryIntentRequest) MarshalToSizedBuffer(d []byte) (int, error) { return sizedBuffer(d, m.Marshal) }
func (m *QueryIntentRequest) Size() int                            { b, _ := m.Marshal(); return len(b) }
func (m *QueryIntentRequest) Unmarshal(data []byte) error {
	*m = QueryIntentRequest{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		if num == 1 && typ == protowire.BytesType {
			v, vn := protowire.ConsumeString(data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
			m.IntentID = v
		} else {
			vn := protowire.ConsumeFieldValue(num, typ, data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
		}
	}
	return nil
}

type QueryIntentResponse struct {
	Intent Intent `json:"intent"`
	Found  bool   `json:"found"`
}

func (m *QueryIntentResponse) ProtoMessage()           {}
func (m *QueryIntentResponse) Reset()                  { *m = QueryIntentResponse{} }
func (m *QueryIntentResponse) String() string          { return fmt.Sprintf("intent: %+v", m.Intent) }
func (m *QueryIntentResponse) XXX_MessageName() string { return "syreen.intent.QueryIntentResponse" }

func (m *QueryIntentResponse) Marshal() ([]byte, error) {
	b, err := appendJSON(nil, 1, &m.Intent)
	if err != nil {
		return nil, err
	}
	return appendBool(b, 2, m.Found), nil
}
func (m *QueryIntentResponse) MarshalToSizedBuffer(d []byte) (int, error) { return sizedBuffer(d, m.Marshal) }
func (m *QueryIntentResponse) Size() int                            { b, _ := m.Marshal(); return len(b) }
func (m *QueryIntentResponse) Unmarshal(data []byte) error {
	*m = QueryIntentResponse{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		switch {
		case num == 1 && typ == protowire.BytesType:
			v, vn := protowire.ConsumeBytes(data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
			if err := json.Unmarshal(v, &m.Intent); err != nil {
				return err
			}
		case num == 2 && typ == protowire.VarintType:
			v, vn := protowire.ConsumeVarint(data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
			m.Found = v != 0
		default:
			vn := protowire.ConsumeFieldValue(num, typ, data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
		}
	}
	return nil
}

type QuerySolverRequest struct {
	Address string `json:"address"`
}

func (m *QuerySolverRequest) ProtoMessage()           {}
func (m *QuerySolverRequest) Reset()                  { *m = QuerySolverRequest{} }
func (m *QuerySolverRequest) String() string          { return fmt.Sprintf("query_solver: addr=%s", m.Address) }
func (m *QuerySolverRequest) XXX_MessageName() string { return "syreen.intent.QuerySolverRequest" }

func (m *QuerySolverRequest) Marshal() ([]byte, error)              { return appendString(nil, 1, m.Address), nil }
func (m *QuerySolverRequest) MarshalToSizedBuffer(d []byte) (int, error) { return sizedBuffer(d, m.Marshal) }
func (m *QuerySolverRequest) Size() int                            { b, _ := m.Marshal(); return len(b) }
func (m *QuerySolverRequest) Unmarshal(data []byte) error {
	*m = QuerySolverRequest{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		if num == 1 && typ == protowire.BytesType {
			v, vn := protowire.ConsumeString(data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
			m.Address = v
		} else {
			vn := protowire.ConsumeFieldValue(num, typ, data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
		}
	}
	return nil
}

type QuerySolverResponse struct {
	Solver Solver `json:"solver"`
	Found  bool   `json:"found"`
}

func (m *QuerySolverResponse) ProtoMessage()           {}
func (m *QuerySolverResponse) Reset()                  { *m = QuerySolverResponse{} }
func (m *QuerySolverResponse) String() string          { return fmt.Sprintf("solver: %+v", m.Solver) }
func (m *QuerySolverResponse) XXX_MessageName() string { return "syreen.intent.QuerySolverResponse" }

func (m *QuerySolverResponse) Marshal() ([]byte, error) {
	b, err := appendJSON(nil, 1, &m.Solver)
	if err != nil {
		return nil, err
	}
	return appendBool(b, 2, m.Found), nil
}
func (m *QuerySolverResponse) MarshalToSizedBuffer(d []byte) (int, error) { return sizedBuffer(d, m.Marshal) }
func (m *QuerySolverResponse) Size() int                            { b, _ := m.Marshal(); return len(b) }
func (m *QuerySolverResponse) Unmarshal(data []byte) error {
	*m = QuerySolverResponse{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		switch {
		case num == 1 && typ == protowire.BytesType:
			v, vn := protowire.ConsumeBytes(data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
			if err := json.Unmarshal(v, &m.Solver); err != nil {
				return err
			}
		case num == 2 && typ == protowire.VarintType:
			v, vn := protowire.ConsumeVarint(data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
			m.Found = v != 0
		default:
			vn := protowire.ConsumeFieldValue(num, typ, data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
		}
	}
	return nil
}

type QuerySolutionsRequest struct {
	IntentID string `json:"intent_id"`
}

func (m *QuerySolutionsRequest) ProtoMessage()           {}
func (m *QuerySolutionsRequest) Reset()                  { *m = QuerySolutionsRequest{} }
func (m *QuerySolutionsRequest) String() string          { return fmt.Sprintf("query_solutions: intent=%s", m.IntentID) }
func (m *QuerySolutionsRequest) XXX_MessageName() string { return "syreen.intent.QuerySolutionsRequest" }

// Marshal/Unmarshal/Size: hand-rolled proto codec so IntentID (which has no
// protobuf struct tag) survives the wire. Without this, reflection marshaling
// drops IntentID and the gRPC query runs against an empty intent id.
func (m *QuerySolutionsRequest) Marshal() ([]byte, error) {
	var b []byte
	if m.IntentID != "" {
		b = protowire.AppendTag(b, 1, protowire.BytesType)
		b = protowire.AppendString(b, m.IntentID)
	}
	return b, nil
}

func (m *QuerySolutionsRequest) MarshalToSizedBuffer(data []byte) (int, error) {
	b, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	n := len(b)
	copy(data[len(data)-n:], b)
	return n, nil
}

func (m *QuerySolutionsRequest) Size() int {
	b, _ := m.Marshal()
	return len(b)
}

func (m *QuerySolutionsRequest) Unmarshal(data []byte) error {
	m.IntentID = ""
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		if num == 1 && typ == protowire.BytesType {
			v, vn := protowire.ConsumeString(data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
			m.IntentID = v
		} else {
			vn := protowire.ConsumeFieldValue(num, typ, data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
		}
	}
	return nil
}

type QuerySolutionsResponse struct {
	Solutions []Solution `json:"solutions"`
}

func (m *QuerySolutionsResponse) ProtoMessage()           {}
func (m *QuerySolutionsResponse) Reset()                  { *m = QuerySolutionsResponse{} }
func (m *QuerySolutionsResponse) String() string          { return fmt.Sprintf("solutions: %d", len(m.Solutions)) }
func (m *QuerySolutionsResponse) XXX_MessageName() string { return "syreen.intent.QuerySolutionsResponse" }

// Marshal/Unmarshal/Size: hand-rolled proto codec for the response.
// Without these, the codec falls back to reflection encoding, which cannot
// represent Solution's []json.RawMessage fields — so the Solutions slice was
// silently DROPPED and the query returned empty even though the server found
// the solutions. Each Solution is encoded as a repeated length-delimited
// bytes field (#1) holding its JSON, matching how solutions are stored.
func (m *QuerySolutionsResponse) Marshal() ([]byte, error) {
	var b []byte
	for i := range m.Solutions {
		sb, err := json.Marshal(&m.Solutions[i])
		if err != nil {
			return nil, err
		}
		b = protowire.AppendTag(b, 1, protowire.BytesType)
		b = protowire.AppendBytes(b, sb)
	}
	return b, nil
}

func (m *QuerySolutionsResponse) MarshalToSizedBuffer(data []byte) (int, error) {
	b, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	n := len(b)
	copy(data[len(data)-n:], b)
	return n, nil
}

func (m *QuerySolutionsResponse) Size() int {
	b, _ := m.Marshal()
	return len(b)
}

func (m *QuerySolutionsResponse) Unmarshal(data []byte) error {
	m.Solutions = nil
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		if num == 1 && typ == protowire.BytesType {
			v, vn := protowire.ConsumeBytes(data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
			var s Solution
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
			m.Solutions = append(m.Solutions, s)
		} else {
			vn := protowire.ConsumeFieldValue(num, typ, data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
		}
	}
	return nil
}

// --- Intents (list all) Query ---

type QueryIntentsRequest struct{}

func (m *QueryIntentsRequest) ProtoMessage()           {}
func (m *QueryIntentsRequest) Reset()                  { *m = QueryIntentsRequest{} }
func (m *QueryIntentsRequest) String() string          { return "query_intents_request" }
func (m *QueryIntentsRequest) XXX_MessageName() string { return "syreen.intent.QueryIntentsRequest" }

type QueryIntentsResponse struct {
	Intents []Intent `json:"intents"`
}

func (m *QueryIntentsResponse) ProtoMessage()           {}
func (m *QueryIntentsResponse) Reset()                  { *m = QueryIntentsResponse{} }
func (m *QueryIntentsResponse) String() string          { return fmt.Sprintf("intents: %d", len(m.Intents)) }
func (m *QueryIntentsResponse) XXX_MessageName() string { return "syreen.intent.QueryIntentsResponse" }

func (m *QueryIntentsResponse) Marshal() ([]byte, error)              { return marshalIntentSlice(m.Intents) }
func (m *QueryIntentsResponse) MarshalToSizedBuffer(d []byte) (int, error) { return sizedBuffer(d, m.Marshal) }
func (m *QueryIntentsResponse) Size() int                            { b, _ := m.Marshal(); return len(b) }
func (m *QueryIntentsResponse) Unmarshal(data []byte) error {
	out, err := unmarshalIntentSlice(data)
	if err != nil {
		return err
	}
	m.Intents = out
	return nil
}

// --- IntentsByType Query ---

type QueryIntentsByTypeRequest struct {
	IntentType string `json:"intent_type"`
}

func (m *QueryIntentsByTypeRequest) ProtoMessage()           {}
func (m *QueryIntentsByTypeRequest) Reset()                  { *m = QueryIntentsByTypeRequest{} }
func (m *QueryIntentsByTypeRequest) String() string          { return fmt.Sprintf("query_intents_by_type: type=%s", m.IntentType) }
func (m *QueryIntentsByTypeRequest) XXX_MessageName() string { return "syreen.intent.QueryIntentsByTypeRequest" }

func (m *QueryIntentsByTypeRequest) Marshal() ([]byte, error)              { return appendString(nil, 1, m.IntentType), nil }
func (m *QueryIntentsByTypeRequest) MarshalToSizedBuffer(d []byte) (int, error) { return sizedBuffer(d, m.Marshal) }
func (m *QueryIntentsByTypeRequest) Size() int                            { b, _ := m.Marshal(); return len(b) }
func (m *QueryIntentsByTypeRequest) Unmarshal(data []byte) error {
	*m = QueryIntentsByTypeRequest{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		if num == 1 && typ == protowire.BytesType {
			v, vn := protowire.ConsumeString(data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
			m.IntentType = v
		} else {
			vn := protowire.ConsumeFieldValue(num, typ, data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
		}
	}
	return nil
}

type QueryIntentsByTypeResponse struct {
	Intents []Intent `json:"intents"`
}

func (m *QueryIntentsByTypeResponse) ProtoMessage()           {}
func (m *QueryIntentsByTypeResponse) Reset()                  { *m = QueryIntentsByTypeResponse{} }
func (m *QueryIntentsByTypeResponse) String() string          { return fmt.Sprintf("intents_by_type: %d", len(m.Intents)) }
func (m *QueryIntentsByTypeResponse) XXX_MessageName() string { return "syreen.intent.QueryIntentsByTypeResponse" }

func (m *QueryIntentsByTypeResponse) Marshal() ([]byte, error)              { return marshalIntentSlice(m.Intents) }
func (m *QueryIntentsByTypeResponse) MarshalToSizedBuffer(d []byte) (int, error) { return sizedBuffer(d, m.Marshal) }
func (m *QueryIntentsByTypeResponse) Size() int                            { b, _ := m.Marshal(); return len(b) }
func (m *QueryIntentsByTypeResponse) Unmarshal(data []byte) error {
	out, err := unmarshalIntentSlice(data)
	if err != nil {
		return err
	}
	m.Intents = out
	return nil
}

// --- QueryServer Interface ---

type QueryServer interface {
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	Intent(context.Context, *QueryIntentRequest) (*QueryIntentResponse, error)
	Intents(context.Context, *QueryIntentsRequest) (*QueryIntentsResponse, error)
	Solver(context.Context, *QuerySolverRequest) (*QuerySolverResponse, error)
	Solutions(context.Context, *QuerySolutionsRequest) (*QuerySolutionsResponse, error)
	IntentsByType(context.Context, *QueryIntentsByTypeRequest) (*QueryIntentsByTypeResponse, error)
	StrategyTemplates(context.Context, *QueryStrategyTemplatesRequest) (*QueryStrategyTemplatesResponse, error)
	Strategies(context.Context, *QueryStrategiesRequest) (*QueryStrategiesResponse, error)
	Strategy(context.Context, *QueryStrategyRequest) (*QueryStrategyResponse, error)
}

// --- Msg Handler Functions ---

func _Msg_SubmitIntent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSubmitIntent)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SubmitIntent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/SubmitIntent"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SubmitIntent(ctx, req.(*MsgSubmitIntent))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RegisterSolver_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterSolver)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RegisterSolver(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/RegisterSolver"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RegisterSolver(ctx, req.(*MsgRegisterSolver))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_DeregisterSolver_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDeregisterSolver)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).DeregisterSolver(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/DeregisterSolver"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).DeregisterSolver(ctx, req.(*MsgDeregisterSolver))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SubmitSolution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSubmitSolution)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SubmitSolution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/SubmitSolution"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SubmitSolution(ctx, req.(*MsgSubmitSolution))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_FulfillIntent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgFulfillIntent)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).FulfillIntent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/FulfillIntent"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).FulfillIntent(ctx, req.(*MsgFulfillIntent))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CancelIntent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCancelIntent)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CancelIntent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/CancelIntent"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CancelIntent(ctx, req.(*MsgCancelIntent))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SubmitChain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSubmitChain)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SubmitChain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/SubmitChain"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SubmitChain(ctx, req.(*MsgSubmitChain))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CancelChain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCancelChain)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CancelChain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/CancelChain"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CancelChain(ctx, req.(*MsgCancelChain))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateStrategy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateStrategy)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateStrategy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/CreateStrategy"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateStrategy(ctx, req.(*MsgCreateStrategy))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CancelStrategy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCancelStrategy)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CancelStrategy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/CancelStrategy"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CancelStrategy(ctx, req.(*MsgCancelStrategy))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/UpdateParams"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

// --- Query Handler Functions ---

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/Params"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Intent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryIntentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Intent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/Intent"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Intent(ctx, req.(*QueryIntentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Solver_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySolverRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Solver(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/Solver"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Solver(ctx, req.(*QuerySolverRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Solutions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySolutionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Solutions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/Solutions"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Solutions(ctx, req.(*QuerySolutionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Intents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryIntentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Intents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/Intents"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Intents(ctx, req.(*QueryIntentsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_IntentsByType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryIntentsByTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).IntentsByType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/IntentsByType"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).IntentsByType(ctx, req.(*QueryIntentsByTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_StrategyTemplates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryStrategyTemplatesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).StrategyTemplates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/StrategyTemplates"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).StrategyTemplates(ctx, req.(*QueryStrategyTemplatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Strategies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryStrategiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Strategies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/Strategies"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Strategies(ctx, req.(*QueryStrategiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Strategy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryStrategyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Strategy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/Strategy"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Strategy(ctx, req.(*QueryStrategyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// --- Service Descriptors ---

var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.intent.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitIntent",
			Handler:    _Msg_SubmitIntent_Handler,
		},
		{
			MethodName: "RegisterSolver",
			Handler:    _Msg_RegisterSolver_Handler,
		},
		{
			MethodName: "DeregisterSolver",
			Handler:    _Msg_DeregisterSolver_Handler,
		},
		{
			MethodName: "SubmitSolution",
			Handler:    _Msg_SubmitSolution_Handler,
		},
		{
			MethodName: "FulfillIntent",
			Handler:    _Msg_FulfillIntent_Handler,
		},
		{
			MethodName: "CancelIntent",
			Handler:    _Msg_CancelIntent_Handler,
		},
		{
			MethodName: "SubmitChain",
			Handler:    _Msg_SubmitChain_Handler,
		},
		{
			MethodName: "CancelChain",
			Handler:    _Msg_CancelChain_Handler,
		},
		{
			MethodName: "CreateStrategy",
			Handler:    _Msg_CreateStrategy_Handler,
		},
		{
			MethodName: "CancelStrategy",
			Handler:    _Msg_CancelStrategy_Handler,
		},
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/intent/tx.proto",
}

var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.intent.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "Intent",
			Handler:    _Query_Intent_Handler,
		},
		{
			MethodName: "Intents",
			Handler:    _Query_Intents_Handler,
		},
		{
			MethodName: "Solver",
			Handler:    _Query_Solver_Handler,
		},
		{
			MethodName: "Solutions",
			Handler:    _Query_Solutions_Handler,
		},
		{
			MethodName: "IntentsByType",
			Handler:    _Query_IntentsByType_Handler,
		},
		{
			MethodName: "StrategyTemplates",
			Handler:    _Query_StrategyTemplates_Handler,
		},
		{
			MethodName: "Strategies",
			Handler:    _Query_Strategies_Handler,
		},
		{
			MethodName: "Strategy",
			Handler:    _Query_Strategy_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/intent/query.proto",
}

// RegisterMsgServer registers the MsgServer implementation with the gRPC server
func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

// RegisterQueryServer registers the QueryServer implementation with the gRPC server
func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}
