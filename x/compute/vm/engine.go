package vm

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Security limits
const (
	MaxRecursionDepth = 10
	MaxStringSize     = 65536  // 64KB
	MaxStateValueSize = 262144 // 256KB
	MaxStateEntries   = 1024
)

// Sentinel errors for security violations
var (
	ErrRecursionDepthExceeded = errors.New("compute VM: maximum recursion depth exceeded")
	ErrStringTooLarge         = errors.New("compute VM: string exceeds maximum size")
	ErrValueTooLarge          = errors.New("compute VM: state value exceeds maximum size")
	ErrTooManyStateEntries    = errors.New("compute VM: too many state entries")
	ErrIntegerOverflow        = errors.New("compute VM: integer overflow")
)

// ContractDefinition is the parsed contract code
type ContractDefinition struct {
	Version     string                 `json:"version"`
	State       map[string]string      `json:"state"`
	Instantiate Handler                `json:"instantiate"`
	Execute     map[string]Handler     `json:"execute"`
	Query       map[string]Handler     `json:"query"`
}

// Handler contains a list of actions to execute
type Handler struct {
	Actions []Action `json:"actions"`
}

// Action represents a single operation in a handler
type Action struct {
	Op    string          `json:"op"`
	Key   string          `json:"key,omitempty"`
	Value string          `json:"value,omitempty"`
	Into  string          `json:"into,omitempty"`
	Left  string          `json:"left,omitempty"`
	Right string          `json:"right,omitempty"`
	Error string          `json:"error,omitempty"`
	Data  json.RawMessage `json:"data,omitempty"`
	From  string          `json:"from,omitempty"`
	To    string          `json:"to,omitempty"`

	Amount string            `json:"amount,omitempty"`
	Denom  string            `json:"denom,omitempty"`
	Event  string            `json:"event,omitempty"`
	Attrs  map[string]string `json:"attrs,omitempty"`

	// For log
	Message string `json:"message,omitempty"`

	// For conditional
	Condition string   `json:"condition,omitempty"`
	Then      []Action `json:"then,omitempty"`
	Else      []Action `json:"else,omitempty"`
}

// Event is a custom event emitted by a contract
type Event struct {
	Type       string            `json:"type"`
	Attributes map[string]string `json:"attributes"`
}

// BankMsg is a bank send message queued during execution
type BankMsg struct {
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	Amount      string `json:"amount"`
	Denom       string `json:"denom"`
}

// ExecutionContext holds the runtime state during VM execution
type ExecutionContext struct {
	ContractAddr   string
	Sender         string
	BlockHeight    string
	State          map[string]string
	Variables      map[string]string
	Msg            map[string]interface{}
	GasUsed        uint64
	GasLimit       uint64
	Response       json.RawMessage
	Events         []Event
	BankMsgs       []BankMsg
	StateChanged   bool
	Logs           []string
	RecursionDepth int
}

// Engine executes contract logic
type Engine struct{}

// NewEngine creates a new VM engine
func NewEngine() *Engine {
	return &Engine{}
}

// ParseContractDefinition parses JSON contract code into a ContractDefinition
func ParseContractDefinition(code []byte) (*ContractDefinition, error) {
	var def ContractDefinition
	if err := json.Unmarshal(code, &def); err != nil {
		return nil, fmt.Errorf("invalid contract definition: %w", err)
	}
	if def.Version == "" {
		return nil, fmt.Errorf("contract definition missing version")
	}
	return &def, nil
}

// NewExecutionContext creates a new execution context with defaults
func NewExecutionContext(contractAddr, sender, blockHeight string, state map[string]string, msg map[string]interface{}, gasLimit uint64) *ExecutionContext {
	stateCopy := make(map[string]string, len(state))
	for k, v := range state {
		stateCopy[k] = v
	}
	return &ExecutionContext{
		ContractAddr: contractAddr,
		Sender:       sender,
		BlockHeight:  blockHeight,
		State:        stateCopy,
		Variables:    make(map[string]string),
		Msg:          msg,
		GasUsed:      0,
		GasLimit:     gasLimit,
		Events:       nil,
		BankMsgs:     nil,
		StateChanged: false,
		Logs:         nil,
	}
}

// ExecuteInstantiate runs the instantiate handler, initializing contract state
func (e *Engine) ExecuteInstantiate(def *ContractDefinition, ctx *ExecutionContext) error {
	// Initialize state from the contract definition defaults
	if def.State != nil {
		for k, v := range def.State {
			if _, exists := ctx.State[k]; !exists {
				ctx.State[k] = v
			}
		}
	}

	return e.executeActions(def.Instantiate.Actions, ctx)
}

// Execute runs a named execute handler
func (e *Engine) Execute(def *ContractDefinition, handler string, ctx *ExecutionContext) error {
	h, ok := def.Execute[handler]
	if !ok {
		return fmt.Errorf("unknown execute handler: %s", handler)
	}
	return e.executeActions(h.Actions, ctx)
}

// ExecuteQuery runs a named query handler (read-only: state changes are discarded)
func (e *Engine) ExecuteQuery(def *ContractDefinition, handler string, ctx *ExecutionContext) error {
	h, ok := def.Query[handler]
	if !ok {
		return fmt.Errorf("unknown query handler: %s", handler)
	}

	// Snapshot state so we can detect and reject mutations
	origState := make(map[string]string, len(ctx.State))
	for k, v := range ctx.State {
		origState[k] = v
	}

	if err := e.executeActions(h.Actions, ctx); err != nil {
		return err
	}

	// Restore original state -- queries are read-only
	ctx.State = origState
	ctx.StateChanged = false
	ctx.BankMsgs = nil
	return nil
}

// executeActions runs a slice of actions sequentially
func (e *Engine) executeActions(actions []Action, ctx *ExecutionContext) error {
	for _, action := range actions {
		if err := e.executeAction(action, ctx); err != nil {
			return err
		}
	}
	return nil
}

// executeAction dispatches a single action to its operation handler
func (e *Engine) executeAction(action Action, ctx *ExecutionContext) error {
	op, ok := operations[action.Op]
	if !ok {
		return fmt.Errorf("unknown operation: %s", action.Op)
	}

	// Charge gas
	if err := ctx.chargeGas(op.gasCost); err != nil {
		return err
	}

	return op.exec(action, ctx)
}

// chargeGas deducts gas and returns an error if limit exceeded
func (ctx *ExecutionContext) chargeGas(amount uint64) error {
	ctx.GasUsed += amount
	if ctx.GasLimit > 0 && ctx.GasUsed > ctx.GasLimit {
		return fmt.Errorf("out of gas: used %d, limit %d", ctx.GasUsed, ctx.GasLimit)
	}
	return nil
}

// resolveValue resolves a string that may contain variable references.
// Variables start with $. Supported forms:
//   - $msg.field       -- reads from the message
//   - $sender          -- tx sender address
//   - $contract        -- contract address
//   - $block_height    -- current block height
//   - $varname         -- runtime variable
//
// If the string does not start with $, it is returned as-is (literal).
func (ctx *ExecutionContext) resolveValue(s string) string {
	if !strings.HasPrefix(s, "$") {
		return s
	}

	name := s[1:]

	// Built-in variables
	switch name {
	case "sender":
		return ctx.Sender
	case "contract":
		return ctx.ContractAddr
	case "block_height":
		return ctx.BlockHeight
	}

	// $msg.field -- read from the message
	if strings.HasPrefix(name, "msg.") {
		fieldPath := name[4:]
		return ctx.resolveMsg(fieldPath)
	}

	// Runtime variable
	if v, ok := ctx.Variables[name]; ok {
		return v
	}

	return ""
}

// resolveMsg reads a dotted path from the message map.
func (ctx *ExecutionContext) resolveMsg(path string) string {
	parts := strings.Split(path, ".")
	var current interface{} = ctx.Msg
	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return ""
		}
		current, ok = m[part]
		if !ok {
			return ""
		}
	}
	switch v := current.(type) {
	case string:
		return v
	case float64:
		// JSON numbers -- format without decimal if integer
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case nil:
		return ""
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

// resolveResponseData resolves variable references inside a json.RawMessage.
// It walks the JSON structure and replaces any string value starting with "$".
func (ctx *ExecutionContext) resolveResponseData(data json.RawMessage) json.RawMessage {
	if len(data) == 0 {
		return data
	}

	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return data
	}

	resolved := ctx.resolveInterface(raw)
	out, err := json.Marshal(resolved)
	if err != nil {
		return data
	}
	return out
}

func (ctx *ExecutionContext) resolveInterface(v interface{}) interface{} {
	return ctx.resolveInterfaceDepth(v, 0)
}

func (ctx *ExecutionContext) resolveInterfaceDepth(v interface{}, depth int) interface{} {
	if depth > MaxRecursionDepth {
		return v // stop resolving at max depth to prevent stack overflow
	}
	switch val := v.(type) {
	case string:
		if strings.HasPrefix(val, "$") {
			return ctx.resolveValue(val)
		}
		return val
	case map[string]interface{}:
		result := make(map[string]interface{}, len(val))
		for k, v2 := range val {
			result[k] = ctx.resolveInterfaceDepth(v2, depth+1)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, v2 := range val {
			result[i] = ctx.resolveInterfaceDepth(v2, depth+1)
		}
		return result
	default:
		return val
	}
}
