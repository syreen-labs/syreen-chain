package vm

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"testing"
)

// counterContract returns a contract definition for a simple counter
func counterContract() string {
	return `{
		"version": "1.0",
		"state": {
			"count": "0",
			"owner": ""
		},
		"instantiate": {
			"actions": [
				{"op": "set", "key": "owner", "value": "$msg.owner"},
				{"op": "set", "key": "count", "value": "$msg.initial_count"},
				{"op": "response", "data": {"status": "initialized"}}
			]
		},
		"execute": {
			"increment": {
				"actions": [
					{"op": "get", "key": "count", "into": "$count"},
					{"op": "int_add", "left": "$count", "right": "1", "into": "$count"},
					{"op": "set", "key": "count", "value": "$count"},
					{"op": "response", "data": {"new_count": "$count"}}
				]
			},
			"decrement": {
				"actions": [
					{"op": "get", "key": "count", "into": "$count"},
					{"op": "require_gt", "left": "$count", "right": "0", "error": "count already zero"},
					{"op": "int_sub", "left": "$count", "right": "1", "into": "$count"},
					{"op": "set", "key": "count", "value": "$count"},
					{"op": "response", "data": {"new_count": "$count"}}
				]
			},
			"transfer_ownership": {
				"actions": [
					{"op": "get", "key": "owner", "into": "$current_owner"},
					{"op": "require_eq", "left": "$sender", "right": "$current_owner", "error": "unauthorized"},
					{"op": "set", "key": "owner", "value": "$msg.new_owner"},
					{"op": "response", "data": {"new_owner": "$msg.new_owner"}}
				]
			},
			"send_tokens": {
				"actions": [
					{"op": "get", "key": "owner", "into": "$current_owner"},
					{"op": "require_eq", "left": "$sender", "right": "$current_owner", "error": "unauthorized"},
					{"op": "bank_send", "from": "$contract", "to": "$msg.recipient", "amount": "$msg.amount", "denom": "$msg.denom"},
					{"op": "response", "data": {"status": "sent"}}
				]
			}
		},
		"query": {
			"get_count": {
				"actions": [
					{"op": "get", "key": "count", "into": "$count"},
					{"op": "response", "data": {"count": "$count"}}
				]
			},
			"get_owner": {
				"actions": [
					{"op": "get", "key": "owner", "into": "$owner"},
					{"op": "response", "data": {"owner": "$owner"}}
				]
			}
		}
	}`
}

func setupCounter(t *testing.T, initialCount string, owner string) (*ContractDefinition, *ExecutionContext) {
	t.Helper()
	def, err := ParseContractDefinition([]byte(counterContract()))
	if err != nil {
		t.Fatalf("parse contract: %v", err)
	}

	msg := map[string]interface{}{
		"owner":         owner,
		"initial_count": initialCount,
	}
	ctx := NewExecutionContext("syreen1contract", owner, "100", nil, msg, 1_000_000)
	eng := NewEngine()
	if err := eng.ExecuteInstantiate(def, ctx); err != nil {
		t.Fatalf("instantiate: %v", err)
	}
	return def, ctx
}

func TestCounter_Increment(t *testing.T) {
	def, ctx := setupCounter(t, "0", "syreen1owner")
	eng := NewEngine()

	// Increment once
	ctx.Msg = map[string]interface{}{}
	if err := eng.Execute(def, "increment", ctx); err != nil {
		t.Fatalf("increment 1: %v", err)
	}
	if ctx.State["count"] != "1" {
		t.Fatalf("expected count=1, got %s", ctx.State["count"])
	}

	// Increment again
	if err := eng.Execute(def, "increment", ctx); err != nil {
		t.Fatalf("increment 2: %v", err)
	}
	if ctx.State["count"] != "2" {
		t.Fatalf("expected count=2, got %s", ctx.State["count"])
	}

	// Check response
	var resp map[string]string
	if err := json.Unmarshal(ctx.Response, &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp["new_count"] != "2" {
		t.Fatalf("expected response new_count=2, got %s", resp["new_count"])
	}
}

func TestCounter_QueryCount(t *testing.T) {
	def, ctx := setupCounter(t, "42", "syreen1owner")
	eng := NewEngine()

	if err := eng.ExecuteQuery(def, "get_count", ctx); err != nil {
		t.Fatalf("query: %v", err)
	}

	var resp map[string]string
	if err := json.Unmarshal(ctx.Response, &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp["count"] != "42" {
		t.Fatalf("expected count=42, got %s", resp["count"])
	}
}

func TestOwnership_TransferOwnership(t *testing.T) {
	def, ctx := setupCounter(t, "0", "syreen1alice")
	eng := NewEngine()

	ctx.Sender = "syreen1alice"
	ctx.Msg = map[string]interface{}{
		"new_owner": "syreen1bob",
	}

	if err := eng.Execute(def, "transfer_ownership", ctx); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if ctx.State["owner"] != "syreen1bob" {
		t.Fatalf("expected owner=syreen1bob, got %s", ctx.State["owner"])
	}

	var resp map[string]string
	if err := json.Unmarshal(ctx.Response, &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp["new_owner"] != "syreen1bob" {
		t.Fatalf("expected response new_owner=syreen1bob, got %s", resp["new_owner"])
	}
}

func TestOwnership_UnauthorizedFails(t *testing.T) {
	def, ctx := setupCounter(t, "0", "syreen1alice")
	eng := NewEngine()

	// Try to transfer as non-owner
	ctx.Sender = "syreen1mallory"
	ctx.Msg = map[string]interface{}{
		"new_owner": "syreen1mallory",
	}

	err := eng.Execute(def, "transfer_ownership", ctx)
	if err == nil {
		t.Fatal("expected error for unauthorized transfer")
	}
	if err.Error() != "contract error: unauthorized" {
		t.Fatalf("unexpected error: %v", err)
	}

	// Owner should not have changed
	if ctx.State["owner"] != "syreen1alice" {
		t.Fatalf("owner should still be alice, got %s", ctx.State["owner"])
	}
}

func TestBankSend_QueuesBankMsg(t *testing.T) {
	def, ctx := setupCounter(t, "0", "syreen1owner")
	eng := NewEngine()

	ctx.Sender = "syreen1owner"
	ctx.Msg = map[string]interface{}{
		"recipient": "syreen1recipient",
		"amount":    "1000",
		"denom":     "usyreen",
	}

	if err := eng.Execute(def, "send_tokens", ctx); err != nil {
		t.Fatalf("send_tokens: %v", err)
	}

	if len(ctx.BankMsgs) != 1 {
		t.Fatalf("expected 1 bank msg, got %d", len(ctx.BankMsgs))
	}

	bm := ctx.BankMsgs[0]
	if bm.FromAddress != "syreen1contract" {
		t.Fatalf("expected from=syreen1contract, got %s", bm.FromAddress)
	}
	if bm.ToAddress != "syreen1recipient" {
		t.Fatalf("expected to=syreen1recipient, got %s", bm.ToAddress)
	}
	if bm.Amount != "1000" {
		t.Fatalf("expected amount=1000, got %s", bm.Amount)
	}
	if bm.Denom != "usyreen" {
		t.Fatalf("expected denom=usyreen, got %s", bm.Denom)
	}
}

func TestGasExhaustion(t *testing.T) {
	def, err := ParseContractDefinition([]byte(counterContract()))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	// Set a very low gas limit
	ctx := NewExecutionContext("syreen1contract", "syreen1owner", "100", nil, map[string]interface{}{
		"owner":         "syreen1owner",
		"initial_count": "0",
	}, 100) // only 100 gas -- not enough for instantiate (3 set ops = 600 gas)

	eng := NewEngine()
	err = eng.ExecuteInstantiate(def, ctx)
	if err == nil {
		t.Fatal("expected out of gas error")
	}
	if err.Error() == "" || !contains(err.Error(), "out of gas") {
		t.Fatalf("expected out of gas error, got: %v", err)
	}
}

func TestIntMath(t *testing.T) {
	eng := NewEngine()

	tests := []struct {
		name   string
		op     string
		left   string
		right  string
		expect string
	}{
		{"add", "int_add", "10", "20", "30"},
		{"sub", "int_sub", "50", "20", "30"},
		{"mul", "int_mul", "6", "7", "42"},
		{"div", "int_div", "100", "3", "33"},
		{"mod", "int_mod", "100", "3", "1"},
		{"add_neg", "int_add", "-5", "3", "-2"},
		{"sub_neg", "int_sub", "3", "10", "-7"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
			action := Action{
				Op:    tc.op,
				Left:  tc.left,
				Right: tc.right,
				Into:  "$result",
			}
			if err := eng.executeAction(action, ctx); err != nil {
				t.Fatalf("error: %v", err)
			}
			if ctx.Variables["result"] != tc.expect {
				t.Fatalf("expected %s, got %s", tc.expect, ctx.Variables["result"])
			}
		})
	}

	// Division by zero
	t.Run("div_by_zero", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{Op: "int_div", Left: "10", Right: "0", Into: "$r"}
		err := eng.executeAction(action, ctx)
		if err == nil {
			t.Fatal("expected division by zero error")
		}
	})

	// Modulo by zero
	t.Run("mod_by_zero", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{Op: "int_mod", Left: "10", Right: "0", Into: "$r"}
		err := eng.executeAction(action, ctx)
		if err == nil {
			t.Fatal("expected division by zero error")
		}
	})
}

func TestRequireGuards(t *testing.T) {
	eng := NewEngine()

	t.Run("require_eq_pass", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{Op: "require_eq", Left: "hello", Right: "hello"}
		if err := eng.executeAction(action, ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("require_eq_fail", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{Op: "require_eq", Left: "hello", Right: "world", Error: "mismatch"}
		err := eng.executeAction(action, ctx)
		if err == nil {
			t.Fatal("expected error")
		}
		if !contains(err.Error(), "mismatch") {
			t.Fatalf("expected 'mismatch' in error, got: %v", err)
		}
	})

	t.Run("require_neq_pass", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{Op: "require_neq", Left: "a", Right: "b"}
		if err := eng.executeAction(action, ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("require_neq_fail", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{Op: "require_neq", Left: "same", Right: "same"}
		if err := eng.executeAction(action, ctx); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("require_gt_pass", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{Op: "require_gt", Left: "10", Right: "5"}
		if err := eng.executeAction(action, ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("require_gt_fail", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{Op: "require_gt", Left: "5", Right: "10", Error: "too small"}
		err := eng.executeAction(action, ctx)
		if err == nil {
			t.Fatal("expected error")
		}
		if !contains(err.Error(), "too small") {
			t.Fatalf("expected 'too small' in error, got: %v", err)
		}
	})

	t.Run("require_gte_pass_equal", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{Op: "require_gte", Left: "5", Right: "5"}
		if err := eng.executeAction(action, ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("require_gte_fail", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{Op: "require_gte", Left: "4", Right: "5"}
		if err := eng.executeAction(action, ctx); err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestStrOperations(t *testing.T) {
	eng := NewEngine()

	t.Run("str_concat", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{Op: "str_concat", Left: "hello", Right: " world", Into: "$result"}
		if err := eng.executeAction(action, ctx); err != nil {
			t.Fatalf("error: %v", err)
		}
		if ctx.Variables["result"] != "hello world" {
			t.Fatalf("expected 'hello world', got %q", ctx.Variables["result"])
		}
	})

	t.Run("str_len", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{Op: "str_len", Value: "hello", Into: "$len"}
		if err := eng.executeAction(action, ctx); err != nil {
			t.Fatalf("error: %v", err)
		}
		if ctx.Variables["len"] != "5" {
			t.Fatalf("expected '5', got %q", ctx.Variables["len"])
		}
	})
}

func TestEmitEvent(t *testing.T) {
	eng := NewEngine()
	ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
	action := Action{
		Op:    "emit_event",
		Event: "transfer",
		Attrs: map[string]string{
			"from":   "$sender",
			"to":     "syreen1bob",
			"amount": "100",
		},
	}
	if err := eng.executeAction(action, ctx); err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(ctx.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(ctx.Events))
	}
	if ctx.Events[0].Type != "transfer" {
		t.Fatalf("expected event type 'transfer', got %q", ctx.Events[0].Type)
	}
	if ctx.Events[0].Attributes["from"] != "s" {
		t.Fatalf("expected from='s', got %q", ctx.Events[0].Attributes["from"])
	}
}

func TestLog(t *testing.T) {
	eng := NewEngine()
	ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
	action := Action{Op: "log", Message: "hello from contract"}
	if err := eng.executeAction(action, ctx); err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(ctx.Logs) != 1 || ctx.Logs[0] != "hello from contract" {
		t.Fatalf("expected log 'hello from contract', got %v", ctx.Logs)
	}
}

func TestConditional(t *testing.T) {
	eng := NewEngine()

	t.Run("truthy", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", map[string]string{"count": "5"}, nil, 1_000_000)
		ctx.Variables["flag"] = "true"
		action := Action{
			Op:        "conditional",
			Condition: "$flag",
			Then: []Action{
				{Op: "set", Key: "branch", Value: "then"},
			},
			Else: []Action{
				{Op: "set", Key: "branch", Value: "else"},
			},
		}
		if err := eng.executeAction(action, ctx); err != nil {
			t.Fatalf("error: %v", err)
		}
		if ctx.State["branch"] != "then" {
			t.Fatalf("expected then branch, got %q", ctx.State["branch"])
		}
	})

	t.Run("falsy", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", map[string]string{}, nil, 1_000_000)
		ctx.Variables["flag"] = "false"
		action := Action{
			Op:        "conditional",
			Condition: "$flag",
			Then: []Action{
				{Op: "set", Key: "branch", Value: "then"},
			},
			Else: []Action{
				{Op: "set", Key: "branch", Value: "else"},
			},
		}
		if err := eng.executeAction(action, ctx); err != nil {
			t.Fatalf("error: %v", err)
		}
		if ctx.State["branch"] != "else" {
			t.Fatalf("expected else branch, got %q", ctx.State["branch"])
		}
	})
}

func TestDeleteState(t *testing.T) {
	eng := NewEngine()
	ctx := NewExecutionContext("c", "s", "1", map[string]string{"key1": "val1", "key2": "val2"}, nil, 1_000_000)
	action := Action{Op: "delete", Key: "key1"}
	if err := eng.executeAction(action, ctx); err != nil {
		t.Fatalf("error: %v", err)
	}
	if _, ok := ctx.State["key1"]; ok {
		t.Fatal("key1 should have been deleted")
	}
	if ctx.State["key2"] != "val2" {
		t.Fatal("key2 should still exist")
	}
	if !ctx.StateChanged {
		t.Fatal("StateChanged should be true")
	}
}

func TestQueryReadOnly(t *testing.T) {
	def, ctx := setupCounter(t, "10", "syreen1owner")
	eng := NewEngine()

	// Query should not allow state modifications to persist
	if err := eng.ExecuteQuery(def, "get_count", ctx); err != nil {
		t.Fatalf("query: %v", err)
	}

	// State should be unchanged
	if ctx.State["count"] != "10" {
		t.Fatalf("expected count=10, got %s", ctx.State["count"])
	}
	if ctx.StateChanged {
		t.Fatal("StateChanged should be false after query")
	}
}

func TestDetermineHandler(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		handler, args, err := DetermineHandler(json.RawMessage(`{"increment": {"by": 5}}`))
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if handler != "increment" {
			t.Fatalf("expected 'increment', got %q", handler)
		}
		if args["by"] != float64(5) {
			t.Fatalf("expected by=5, got %v", args["by"])
		}
	})

	t.Run("empty_args", func(t *testing.T) {
		handler, _, err := DetermineHandler(json.RawMessage(`{"get_count": {}}`))
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if handler != "get_count" {
			t.Fatalf("expected 'get_count', got %q", handler)
		}
	})

	t.Run("multiple_keys_error", func(t *testing.T) {
		_, _, err := DetermineHandler(json.RawMessage(`{"a": {}, "b": {}}`))
		if err == nil {
			t.Fatal("expected error for multiple keys")
		}
	})
}

func TestUnknownHandler(t *testing.T) {
	def, ctx := setupCounter(t, "0", "syreen1owner")
	eng := NewEngine()
	err := eng.Execute(def, "nonexistent", ctx)
	if err == nil {
		t.Fatal("expected error for unknown handler")
	}
	if !contains(err.Error(), "unknown execute handler") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnknownOp(t *testing.T) {
	eng := NewEngine()
	ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
	action := Action{Op: "foobar"}
	err := eng.executeAction(action, ctx)
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestMsgFieldResolution(t *testing.T) {
	ctx := NewExecutionContext("c", "s", "1", nil, map[string]interface{}{
		"nested": map[string]interface{}{
			"deep": "value",
		},
		"number": float64(42),
		"flag":   true,
	}, 1_000_000)

	if v := ctx.resolveValue("$msg.nested.deep"); v != "value" {
		t.Fatalf("expected 'value', got %q", v)
	}
	if v := ctx.resolveValue("$msg.number"); v != "42" {
		t.Fatalf("expected '42', got %q", v)
	}
	if v := ctx.resolveValue("$msg.flag"); v != "true" {
		t.Fatalf("expected 'true', got %q", v)
	}
	if v := ctx.resolveValue("$msg.missing"); v != "" {
		t.Fatalf("expected '', got %q", v)
	}
	if v := ctx.resolveValue("literal"); v != "literal" {
		t.Fatalf("expected 'literal', got %q", v)
	}
}

func TestBuiltinVars(t *testing.T) {
	ctx := NewExecutionContext("syreen1contract", "syreen1sender", "999", nil, nil, 1_000_000)
	if v := ctx.resolveValue("$sender"); v != "syreen1sender" {
		t.Fatalf("expected sender, got %q", v)
	}
	if v := ctx.resolveValue("$contract"); v != "syreen1contract" {
		t.Fatalf("expected contract, got %q", v)
	}
	if v := ctx.resolveValue("$block_height"); v != "999" {
		t.Fatalf("expected 999, got %q", v)
	}
}

func TestParseContractDefinition_Invalid(t *testing.T) {
	_, err := ParseContractDefinition([]byte(`{not json`))
	if err == nil {
		t.Fatal("expected error")
	}

	_, err = ParseContractDefinition([]byte(`{"state":{}}`))
	if err == nil {
		t.Fatal("expected error for missing version")
	}
}

func TestBankSend_ValidationErrors(t *testing.T) {
	eng := NewEngine()

	t.Run("missing_fields", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{Op: "bank_send", From: "$contract"}
		err := eng.executeAction(action, ctx)
		if err == nil {
			t.Fatal("expected error for missing fields")
		}
	})

	t.Run("negative_amount", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{Op: "bank_send", From: "a", To: "b", Amount: "-5", Denom: "usyreen"}
		err := eng.executeAction(action, ctx)
		if err == nil {
			t.Fatal("expected error for negative amount")
		}
	})
}

// contains checks if s contains substr
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// =============================================================================
// Security fix tests
// =============================================================================

// TestConditional_MaxRecursionDepth verifies C-02: deeply nested conditionals
// are rejected once MaxRecursionDepth is exceeded.
func TestConditional_MaxRecursionDepth(t *testing.T) {
	// Build a deeply nested conditional that exceeds MaxRecursionDepth.
	// Each level is: conditional(condition=true, then=[conditional(...)])
	var innerActions []Action
	innerActions = []Action{
		{Op: "set", Key: "reached", Value: "yes"},
	}
	for i := 0; i < MaxRecursionDepth+5; i++ {
		innerActions = []Action{
			{
				Op:        "conditional",
				Condition: "true",
				Then:      innerActions,
			},
		}
	}

	eng := NewEngine()
	ctx := NewExecutionContext("c", "s", "1", map[string]string{}, nil, 100_000_000)

	// Execute the outermost conditional
	err := eng.executeActions(innerActions, ctx)
	if err == nil {
		t.Fatal("expected recursion depth exceeded error")
	}
	if !errors.Is(err, ErrRecursionDepthExceeded) {
		t.Fatalf("expected ErrRecursionDepthExceeded, got: %v", err)
	}

	// Verify that state was NOT set (execution was aborted)
	if ctx.State["reached"] == "yes" {
		t.Fatal("deeply nested conditional should not have completed execution")
	}
}

// TestStrConcat_MaxSize verifies C-03: concatenation producing a string
// larger than MaxStringSize is rejected.
func TestStrConcat_MaxSize(t *testing.T) {
	eng := NewEngine()
	ctx := NewExecutionContext("c", "s", "1", nil, nil, 100_000_000)

	// Create two strings whose combined length exceeds MaxStringSize
	bigLeft := strings.Repeat("A", MaxStringSize/2+1)
	bigRight := strings.Repeat("B", MaxStringSize/2+1)

	ctx.Variables["left"] = bigLeft
	ctx.Variables["right"] = bigRight

	action := Action{Op: "str_concat", Left: "$left", Right: "$right", Into: "$result"}
	err := eng.executeAction(action, ctx)
	if err == nil {
		t.Fatal("expected string too large error")
	}
	if !errors.Is(err, ErrStringTooLarge) {
		t.Fatalf("expected ErrStringTooLarge, got: %v", err)
	}

	// A concat within limits should still work
	t.Run("within_limit", func(t *testing.T) {
		ctx2 := NewExecutionContext("c", "s", "1", nil, nil, 100_000_000)
		action2 := Action{Op: "str_concat", Left: "hello", Right: " world", Into: "$result"}
		if err := eng.executeAction(action2, ctx2); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ctx2.Variables["result"] != "hello world" {
			t.Fatalf("expected 'hello world', got %q", ctx2.Variables["result"])
		}
	})
}

// TestSet_MaxValueSize verifies C-04: state values exceeding MaxStateValueSize are rejected.
func TestSet_MaxValueSize(t *testing.T) {
	eng := NewEngine()
	ctx := NewExecutionContext("c", "s", "1", map[string]string{}, nil, 100_000_000)

	bigValue := strings.Repeat("X", MaxStateValueSize+1)
	ctx.Variables["big"] = bigValue

	action := Action{Op: "set", Key: "mykey", Value: "$big"}
	err := eng.executeAction(action, ctx)
	if err == nil {
		t.Fatal("expected value too large error")
	}
	if !errors.Is(err, ErrValueTooLarge) {
		t.Fatalf("expected ErrValueTooLarge, got: %v", err)
	}

	// Value should NOT have been stored
	if _, ok := ctx.State["mykey"]; ok {
		t.Fatal("oversized value should not have been stored")
	}
}

// TestSet_MaxEntries verifies C-04: more than MaxStateEntries new keys are rejected.
func TestSet_MaxEntries(t *testing.T) {
	eng := NewEngine()
	state := make(map[string]string, MaxStateEntries)
	for i := 0; i < MaxStateEntries; i++ {
		state[fmt.Sprintf("key_%d", i)] = "v"
	}
	ctx := NewExecutionContext("c", "s", "1", state, nil, 100_000_000)

	// Updating an existing key should work
	action := Action{Op: "set", Key: "key_0", Value: "updated"}
	if err := eng.executeAction(action, ctx); err != nil {
		t.Fatalf("updating existing key should work: %v", err)
	}

	// Adding a new key should fail
	action = Action{Op: "set", Key: "new_key", Value: "val"}
	err := eng.executeAction(action, ctx)
	if err == nil {
		t.Fatal("expected too many state entries error")
	}
	if !errors.Is(err, ErrTooManyStateEntries) {
		t.Fatalf("expected ErrTooManyStateEntries, got: %v", err)
	}
}

// TestIntOverflow_Add verifies C-08: int_add overflow detection.
func TestIntOverflow_Add(t *testing.T) {
	eng := NewEngine()

	t.Run("positive_overflow", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{
			Op:    "int_add",
			Left:  fmt.Sprintf("%d", math.MaxInt64),
			Right: "1",
			Into:  "$r",
		}
		err := eng.executeAction(action, ctx)
		if err == nil {
			t.Fatal("expected overflow error")
		}
		if !errors.Is(err, ErrIntegerOverflow) {
			t.Fatalf("expected ErrIntegerOverflow, got: %v", err)
		}
	})

	t.Run("negative_overflow", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{
			Op:    "int_add",
			Left:  fmt.Sprintf("%d", math.MinInt64),
			Right: "-1",
			Into:  "$r",
		}
		err := eng.executeAction(action, ctx)
		if err == nil {
			t.Fatal("expected overflow error")
		}
		if !errors.Is(err, ErrIntegerOverflow) {
			t.Fatalf("expected ErrIntegerOverflow, got: %v", err)
		}
	})

	t.Run("no_overflow", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{Op: "int_add", Left: "100", Right: "200", Into: "$r"}
		if err := eng.executeAction(action, ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ctx.Variables["r"] != "300" {
			t.Fatalf("expected 300, got %s", ctx.Variables["r"])
		}
	})
}

// TestIntOverflow_Sub verifies C-08: int_sub overflow detection.
func TestIntOverflow_Sub(t *testing.T) {
	eng := NewEngine()

	t.Run("positive_overflow", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{
			Op:    "int_sub",
			Left:  fmt.Sprintf("%d", math.MaxInt64),
			Right: "-1",
			Into:  "$r",
		}
		err := eng.executeAction(action, ctx)
		if err == nil {
			t.Fatal("expected overflow error")
		}
		if !errors.Is(err, ErrIntegerOverflow) {
			t.Fatalf("expected ErrIntegerOverflow, got: %v", err)
		}
	})

	t.Run("negative_overflow", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{
			Op:    "int_sub",
			Left:  fmt.Sprintf("%d", math.MinInt64),
			Right: "1",
			Into:  "$r",
		}
		err := eng.executeAction(action, ctx)
		if err == nil {
			t.Fatal("expected overflow error")
		}
		if !errors.Is(err, ErrIntegerOverflow) {
			t.Fatalf("expected ErrIntegerOverflow, got: %v", err)
		}
	})

	t.Run("no_overflow", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{Op: "int_sub", Left: "100", Right: "200", Into: "$r"}
		if err := eng.executeAction(action, ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ctx.Variables["r"] != "-100" {
			t.Fatalf("expected -100, got %s", ctx.Variables["r"])
		}
	})
}

// TestIntOverflow_Mul verifies C-08: int_mul overflow detection.
func TestIntOverflow_Mul(t *testing.T) {
	eng := NewEngine()

	t.Run("positive_overflow", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{
			Op:    "int_mul",
			Left:  fmt.Sprintf("%d", math.MaxInt64),
			Right: "2",
			Into:  "$r",
		}
		err := eng.executeAction(action, ctx)
		if err == nil {
			t.Fatal("expected overflow error")
		}
		if !errors.Is(err, ErrIntegerOverflow) {
			t.Fatalf("expected ErrIntegerOverflow, got: %v", err)
		}
	})

	t.Run("negative_overflow", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{
			Op:    "int_mul",
			Left:  fmt.Sprintf("%d", math.MinInt64),
			Right: "2",
			Into:  "$r",
		}
		err := eng.executeAction(action, ctx)
		if err == nil {
			t.Fatal("expected overflow error")
		}
		if !errors.Is(err, ErrIntegerOverflow) {
			t.Fatalf("expected ErrIntegerOverflow, got: %v", err)
		}
	})

	t.Run("mixed_sign_overflow", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{
			Op:    "int_mul",
			Left:  fmt.Sprintf("%d", math.MaxInt64),
			Right: "-2",
			Into:  "$r",
		}
		err := eng.executeAction(action, ctx)
		if err == nil {
			t.Fatal("expected overflow error")
		}
		if !errors.Is(err, ErrIntegerOverflow) {
			t.Fatalf("expected ErrIntegerOverflow, got: %v", err)
		}
	})

	t.Run("no_overflow", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{Op: "int_mul", Left: "100", Right: "200", Into: "$r"}
		if err := eng.executeAction(action, ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ctx.Variables["r"] != "20000" {
			t.Fatalf("expected 20000, got %s", ctx.Variables["r"])
		}
	})

	t.Run("zero_no_overflow", func(t *testing.T) {
		ctx := NewExecutionContext("c", "s", "1", nil, nil, 1_000_000)
		action := Action{
			Op:    "int_mul",
			Left:  fmt.Sprintf("%d", math.MaxInt64),
			Right: "0",
			Into:  "$r",
		}
		if err := eng.executeAction(action, ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ctx.Variables["r"] != "0" {
			t.Fatalf("expected 0, got %s", ctx.Variables["r"])
		}
	})
}

// TestBankSend_ForgedFromRejected verifies H-03: bank_send only allows
// sending from the contract address.
func TestBankSend_ForgedFromRejected(t *testing.T) {
	eng := NewEngine()

	t.Run("forged_from_rejected", func(t *testing.T) {
		ctx := NewExecutionContext("syreen1contract", "syreen1sender", "1", nil, nil, 1_000_000)
		action := Action{
			Op:     "bank_send",
			From:   "syreen1attacker",
			To:     "syreen1recipient",
			Amount: "1000",
			Denom:  "usyreen",
		}
		err := eng.executeAction(action, ctx)
		if err == nil {
			t.Fatal("expected error for forged from address")
		}
		if !contains(err.Error(), "can only send from contract address") {
			t.Fatalf("expected 'can only send from contract address' in error, got: %v", err)
		}
	})

	t.Run("sender_as_from_rejected", func(t *testing.T) {
		ctx := NewExecutionContext("syreen1contract", "syreen1sender", "1", nil, nil, 1_000_000)
		action := Action{
			Op:     "bank_send",
			From:   "syreen1sender",
			To:     "syreen1recipient",
			Amount: "1000",
			Denom:  "usyreen",
		}
		err := eng.executeAction(action, ctx)
		if err == nil {
			t.Fatal("expected error: sender is not the contract address")
		}
		if !contains(err.Error(), "can only send from contract address") {
			t.Fatalf("expected 'can only send from contract address' in error, got: %v", err)
		}
	})

	t.Run("contract_addr_allowed", func(t *testing.T) {
		ctx := NewExecutionContext("syreen1contract", "syreen1sender", "1", nil, nil, 1_000_000)
		action := Action{
			Op:     "bank_send",
			From:   "syreen1contract",
			To:     "syreen1recipient",
			Amount: "1000",
			Denom:  "usyreen",
		}
		err := eng.executeAction(action, ctx)
		if err != nil {
			t.Fatalf("expected success for contract address, got: %v", err)
		}
		if len(ctx.BankMsgs) != 1 {
			t.Fatalf("expected 1 bank msg, got %d", len(ctx.BankMsgs))
		}
	})
}

// Ensure unused imports are referenced
var _ = json.RawMessage{}
var _ = errors.New
var _ = fmt.Sprintf
var _ = math.MaxInt64
var _ = strings.Repeat
