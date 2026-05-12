package vm

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// operation defines a VM operation with its gas cost and executor function
type operation struct {
	gasCost uint64
	exec    func(action Action, ctx *ExecutionContext) error
}

// operations registry maps operation names to their implementations
var operations map[string]operation

func init() {
	operations = map[string]operation{
		"get":                {gasCost: 100, exec: opGet},
		"set":                {gasCost: 200, exec: opSet},
		"delete":             {gasCost: 150, exec: opDelete},
		"int_add":            {gasCost: 50, exec: opIntAdd},
		"int_sub":            {gasCost: 50, exec: opIntSub},
		"int_mul":            {gasCost: 50, exec: opIntMul},
		"int_div":            {gasCost: 50, exec: opIntDiv},
		"int_mod":            {gasCost: 50, exec: opIntMod},
		"str_concat":         {gasCost: 50, exec: opStrConcat},
		"str_len":            {gasCost: 50, exec: opStrLen},
		"require_eq":         {gasCost: 50, exec: opRequireEq},
		"require_neq":        {gasCost: 50, exec: opRequireNeq},
		"require_gt":         {gasCost: 50, exec: opRequireGt},
		"require_gte":        {gasCost: 50, exec: opRequireGte},
		"bank_send":          {gasCost: 500, exec: opBankSend},
		"bank_query_balance": {gasCost: 300, exec: opBankQueryBalance},
		"emit_event":         {gasCost: 100, exec: opEmitEvent},
		"response":           {gasCost: 50, exec: opResponse},
		"log":                {gasCost: 10, exec: opLog},
		"conditional":        {gasCost: 50, exec: opConditional},
	}
}

// opGet reads a value from contract state into a variable
func opGet(action Action, ctx *ExecutionContext) error {
	key := ctx.resolveValue(action.Key)
	into := action.Into
	if into == "" {
		return fmt.Errorf("get: 'into' field is required")
	}
	// Strip leading $ from into variable name
	varName := strings.TrimPrefix(into, "$")
	val, ok := ctx.State[key]
	if !ok {
		val = ""
	}
	ctx.Variables[varName] = val
	return nil
}

// opSet writes a value to contract state
func opSet(action Action, ctx *ExecutionContext) error {
	key := ctx.resolveValue(action.Key)
	value := ctx.resolveValue(action.Value)

	// C-04: Enforce maximum state value size
	if len(value) > MaxStateValueSize {
		return ErrValueTooLarge
	}

	// C-04: Enforce maximum number of state entries (only for new keys)
	if _, exists := ctx.State[key]; !exists {
		if len(ctx.State) >= MaxStateEntries {
			return ErrTooManyStateEntries
		}
	}

	// C-04: Charge proportional gas for large values
	extraGas := uint64(len(value)) / 32
	if extraGas > 0 {
		if err := ctx.chargeGas(extraGas); err != nil {
			return err
		}
	}

	ctx.State[key] = value
	ctx.StateChanged = true
	return nil
}

// opDelete removes a key from contract state
func opDelete(action Action, ctx *ExecutionContext) error {
	key := ctx.resolveValue(action.Key)
	delete(ctx.State, key)
	ctx.StateChanged = true
	return nil
}

// parseInt64 parses a resolved string as int64
func parseInt64(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("cannot parse empty string as integer")
	}
	return strconv.ParseInt(s, 10, 64)
}

// opIntAdd performs integer addition: left + right -> into
func opIntAdd(action Action, ctx *ExecutionContext) error {
	left, err := parseInt64(ctx.resolveValue(action.Left))
	if err != nil {
		return fmt.Errorf("int_add: invalid left operand: %w", err)
	}
	right, err := parseInt64(ctx.resolveValue(action.Right))
	if err != nil {
		return fmt.Errorf("int_add: invalid right operand: %w", err)
	}
	// C-08: Overflow detection
	if (right > 0 && left > math.MaxInt64-right) || (right < 0 && left < math.MinInt64-right) {
		return ErrIntegerOverflow
	}
	varName := strings.TrimPrefix(action.Into, "$")
	ctx.Variables[varName] = strconv.FormatInt(left+right, 10)
	return nil
}

// opIntSub performs integer subtraction: left - right -> into
func opIntSub(action Action, ctx *ExecutionContext) error {
	left, err := parseInt64(ctx.resolveValue(action.Left))
	if err != nil {
		return fmt.Errorf("int_sub: invalid left operand: %w", err)
	}
	right, err := parseInt64(ctx.resolveValue(action.Right))
	if err != nil {
		return fmt.Errorf("int_sub: invalid right operand: %w", err)
	}
	// C-08: Overflow detection
	if (right < 0 && left > math.MaxInt64+right) || (right > 0 && left < math.MinInt64+right) {
		return ErrIntegerOverflow
	}
	varName := strings.TrimPrefix(action.Into, "$")
	ctx.Variables[varName] = strconv.FormatInt(left-right, 10)
	return nil
}

// opIntMul performs integer multiplication: left * right -> into
func opIntMul(action Action, ctx *ExecutionContext) error {
	left, err := parseInt64(ctx.resolveValue(action.Left))
	if err != nil {
		return fmt.Errorf("int_mul: invalid left operand: %w", err)
	}
	right, err := parseInt64(ctx.resolveValue(action.Right))
	if err != nil {
		return fmt.Errorf("int_mul: invalid right operand: %w", err)
	}
	// C-08: Overflow detection
	if left != 0 && right != 0 {
		if (left > 0 && right > 0 && left > math.MaxInt64/right) ||
			(left < 0 && right < 0 && left < math.MaxInt64/right) ||
			(left > 0 && right < 0 && right < math.MinInt64/left) ||
			(left < 0 && right > 0 && left < math.MinInt64/right) {
			return ErrIntegerOverflow
		}
	}
	varName := strings.TrimPrefix(action.Into, "$")
	ctx.Variables[varName] = strconv.FormatInt(left*right, 10)
	return nil
}

// opIntDiv performs integer division: left / right -> into (with zero check)
func opIntDiv(action Action, ctx *ExecutionContext) error {
	left, err := parseInt64(ctx.resolveValue(action.Left))
	if err != nil {
		return fmt.Errorf("int_div: invalid left operand: %w", err)
	}
	right, err := parseInt64(ctx.resolveValue(action.Right))
	if err != nil {
		return fmt.Errorf("int_div: invalid right operand: %w", err)
	}
	if right == 0 {
		return fmt.Errorf("int_div: division by zero")
	}
	varName := strings.TrimPrefix(action.Into, "$")
	ctx.Variables[varName] = strconv.FormatInt(left/right, 10)
	return nil
}

// opIntMod performs integer modulo: left % right -> into
func opIntMod(action Action, ctx *ExecutionContext) error {
	left, err := parseInt64(ctx.resolveValue(action.Left))
	if err != nil {
		return fmt.Errorf("int_mod: invalid left operand: %w", err)
	}
	right, err := parseInt64(ctx.resolveValue(action.Right))
	if err != nil {
		return fmt.Errorf("int_mod: invalid right operand: %w", err)
	}
	if right == 0 {
		return fmt.Errorf("int_mod: division by zero")
	}
	varName := strings.TrimPrefix(action.Into, "$")
	ctx.Variables[varName] = strconv.FormatInt(left%right, 10)
	return nil
}

// opStrConcat concatenates left + right -> into
func opStrConcat(action Action, ctx *ExecutionContext) error {
	left := ctx.resolveValue(action.Left)
	right := ctx.resolveValue(action.Right)

	// C-03: Enforce maximum string size
	if len(left)+len(right) > MaxStringSize {
		return ErrStringTooLarge
	}

	// C-03: Charge proportional gas (1 gas per 64 bytes)
	extraGas := uint64(len(left)+len(right)) / 64
	if extraGas > 0 {
		if err := ctx.chargeGas(extraGas); err != nil {
			return err
		}
	}

	varName := strings.TrimPrefix(action.Into, "$")
	ctx.Variables[varName] = left + right
	return nil
}

// opStrLen computes the length of the value -> into
func opStrLen(action Action, ctx *ExecutionContext) error {
	val := ctx.resolveValue(action.Value)
	varName := strings.TrimPrefix(action.Into, "$")
	ctx.Variables[varName] = strconv.Itoa(len(val))
	return nil
}

// opRequireEq asserts that left == right, or returns an error
func opRequireEq(action Action, ctx *ExecutionContext) error {
	left := ctx.resolveValue(action.Left)
	right := ctx.resolveValue(action.Right)
	if left != right {
		errMsg := action.Error
		if errMsg == "" {
			errMsg = fmt.Sprintf("require_eq failed: %q != %q", left, right)
		}
		return fmt.Errorf("contract error: %s", errMsg)
	}
	return nil
}

// opRequireNeq asserts that left != right
func opRequireNeq(action Action, ctx *ExecutionContext) error {
	left := ctx.resolveValue(action.Left)
	right := ctx.resolveValue(action.Right)
	if left == right {
		errMsg := action.Error
		if errMsg == "" {
			errMsg = fmt.Sprintf("require_neq failed: %q == %q", left, right)
		}
		return fmt.Errorf("contract error: %s", errMsg)
	}
	return nil
}

// opRequireGt asserts that left > right (integer comparison)
func opRequireGt(action Action, ctx *ExecutionContext) error {
	left, err := parseInt64(ctx.resolveValue(action.Left))
	if err != nil {
		return fmt.Errorf("require_gt: invalid left operand: %w", err)
	}
	right, err := parseInt64(ctx.resolveValue(action.Right))
	if err != nil {
		return fmt.Errorf("require_gt: invalid right operand: %w", err)
	}
	if left <= right {
		errMsg := action.Error
		if errMsg == "" {
			errMsg = fmt.Sprintf("require_gt failed: %d <= %d", left, right)
		}
		return fmt.Errorf("contract error: %s", errMsg)
	}
	return nil
}

// opRequireGte asserts that left >= right (integer comparison)
func opRequireGte(action Action, ctx *ExecutionContext) error {
	left, err := parseInt64(ctx.resolveValue(action.Left))
	if err != nil {
		return fmt.Errorf("require_gte: invalid left operand: %w", err)
	}
	right, err := parseInt64(ctx.resolveValue(action.Right))
	if err != nil {
		return fmt.Errorf("require_gte: invalid right operand: %w", err)
	}
	if left < right {
		errMsg := action.Error
		if errMsg == "" {
			errMsg = fmt.Sprintf("require_gte failed: %d < %d", left, right)
		}
		return fmt.Errorf("contract error: %s", errMsg)
	}
	return nil
}

// opBankSend queues a bank send message for post-execution processing
func opBankSend(action Action, ctx *ExecutionContext) error {
	from := ctx.resolveValue(action.From)
	to := ctx.resolveValue(action.To)
	amount := ctx.resolveValue(action.Amount)
	denom := ctx.resolveValue(action.Denom)

	if from == "" || to == "" || amount == "" || denom == "" {
		return fmt.Errorf("bank_send: from, to, amount, and denom are all required")
	}

	// H-03: Enforce that from must be the contract address
	if from != ctx.ContractAddr {
		return fmt.Errorf("bank_send: can only send from contract address %s, got %s", ctx.ContractAddr, from)
	}

	// Validate amount is a positive integer
	amt, err := parseInt64(amount)
	if err != nil {
		return fmt.Errorf("bank_send: invalid amount: %w", err)
	}
	if amt <= 0 {
		return fmt.Errorf("bank_send: amount must be positive")
	}

	ctx.BankMsgs = append(ctx.BankMsgs, BankMsg{
		FromAddress: from,
		ToAddress:   to,
		Amount:      amount,
		Denom:       denom,
	})
	return nil
}

// opBankQueryBalance queries the balance of an account (placeholder -- actual
// balance lookup is performed by the keeper layer; here we just set a variable)
func opBankQueryBalance(action Action, ctx *ExecutionContext) error {
	// The keeper integration layer should populate this.
	// At the VM level we store a placeholder or the value from ctx.Variables if
	// the keeper pre-populated it.
	addr := ctx.resolveValue(action.Key)
	denom := ctx.resolveValue(action.Denom)
	varName := strings.TrimPrefix(action.Into, "$")

	// Look for a pre-populated balance key set by the keeper
	balKey := fmt.Sprintf("_balance/%s/%s", addr, denom)
	if bal, ok := ctx.Variables[balKey]; ok {
		ctx.Variables[varName] = bal
	} else {
		ctx.Variables[varName] = "0"
	}
	return nil
}

// opEmitEvent emits a custom event
func opEmitEvent(action Action, ctx *ExecutionContext) error {
	eventType := ctx.resolveValue(action.Event)
	if eventType == "" {
		return fmt.Errorf("emit_event: event type is required")
	}
	resolvedAttrs := make(map[string]string, len(action.Attrs))
	for k, v := range action.Attrs {
		resolvedAttrs[k] = ctx.resolveValue(v)
	}
	ctx.Events = append(ctx.Events, Event{
		Type:       eventType,
		Attributes: resolvedAttrs,
	})
	return nil
}

// opResponse sets the response data, resolving any variable references in the JSON
func opResponse(action Action, ctx *ExecutionContext) error {
	if len(action.Data) == 0 {
		return fmt.Errorf("response: data is required")
	}
	ctx.Response = ctx.resolveResponseData(action.Data)
	return nil
}

// opLog appends a log message
func opLog(action Action, ctx *ExecutionContext) error {
	msg := ctx.resolveValue(action.Message)
	if msg == "" {
		msg = ctx.resolveValue(action.Value)
	}
	ctx.Logs = append(ctx.Logs, msg)
	return nil
}

// opConditional implements if-then-else branching.
// The condition field is evaluated: "true", non-empty, and non-"0" are truthy.
func opConditional(action Action, ctx *ExecutionContext) error {
	// C-02: Check recursion depth before descending
	ctx.RecursionDepth++
	if ctx.RecursionDepth > MaxRecursionDepth {
		ctx.RecursionDepth--
		return ErrRecursionDepthExceeded
	}
	defer func() { ctx.RecursionDepth-- }()

	cond := ctx.resolveValue(action.Condition)
	truthy := cond != "" && cond != "0" && cond != "false"

	// NOTE: gas for the conditional op itself was already charged.
	// Actions inside then/else branches charge their own gas.
	if truthy {
		return executeActionsRecursive(action.Then, ctx)
	}
	return executeActionsRecursive(action.Else, ctx)
}

// executeActionsRecursive runs a slice of actions, dispatching each through
// the operations map. This is separate from Engine.executeActions to break
// the initialization cycle (operations map -> opConditional -> Engine -> operations).
func executeActionsRecursive(actions []Action, ctx *ExecutionContext) error {
	for _, action := range actions {
		op, ok := operations[action.Op]
		if !ok {
			return fmt.Errorf("unknown operation: %s", action.Op)
		}
		if err := ctx.chargeGas(op.gasCost); err != nil {
			return err
		}
		if err := op.exec(action, ctx); err != nil {
			return err
		}
	}
	return nil
}

// DetermineHandler extracts the handler name from a JSON message.
// The message is expected to be a JSON object with a single top-level key,
// which is the handler name, following the CosmWasm convention.
// For example: {"increment": {}} -> "increment", {}
func DetermineHandler(msg json.RawMessage) (string, map[string]interface{}, error) {
	var outer map[string]json.RawMessage
	if err := json.Unmarshal(msg, &outer); err != nil {
		return "", nil, fmt.Errorf("invalid message format: %w", err)
	}
	if len(outer) != 1 {
		return "", nil, fmt.Errorf("message must have exactly one top-level key (handler name), got %d", len(outer))
	}

	for handlerName, innerRaw := range outer {
		var inner map[string]interface{}
		if err := json.Unmarshal(innerRaw, &inner); err != nil {
			// If the inner value is not an object, treat it as empty args
			inner = make(map[string]interface{})
		}
		return handlerName, inner, nil
	}
	// unreachable
	return "", nil, fmt.Errorf("empty message")
}
