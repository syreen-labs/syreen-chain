package ante

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/abstractaccount/types"
)

// SessionKeyKeeper defines the keeper methods needed by the session key decorator.
// Only read-only methods are required for the ante decorator; spend tracking is
// done in the PostHandler via RecordSessionKeyUsage after tx execution.
type SessionKeyKeeper interface {
	GetSessionKeyByGrantee(ctx context.Context, grantee string) (types.SessionKey, bool)
	ValidateSessionKey(ctx context.Context, key string, msgType string, spendAmount sdk.Coins) error
	RecordSessionKeyUsage(ctx context.Context, key string, spendAmount sdk.Coins) error
}

// SessionKeyDecorator checks if the tx signer is a session key grantee.
// It performs READ-ONLY validation: verifies expiry, message type permissions,
// and that spend amounts would be within limits, but does NOT persist any
// state changes. This is critical because the decorator runs before signature
// verification -- an unsigned tx must not be able to alter session key state (H-01).
//
// Spend tracking (updating the Used field) is handled by SessionKeyPostDecorator
// after the transaction has been fully verified and executed.
type SessionKeyDecorator struct {
	aaKeeper SessionKeyKeeper
}

// NewSessionKeyDecorator returns a new SessionKeyDecorator.
func NewSessionKeyDecorator(aaKeeper SessionKeyKeeper) SessionKeyDecorator {
	return SessionKeyDecorator{
		aaKeeper: aaKeeper,
	}
}

func (skd SessionKeyDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	sigTx, ok := tx.(sdk.Tx)
	if !ok {
		return next(ctx, tx, simulate)
	}

	msgs := sigTx.GetMsgs()
	if len(msgs) == 0 {
		return next(ctx, tx, simulate)
	}

	// Get the first signer address from the first message's signer field.
	firstMsg := msgs[0]
	signerAddr := extractMsgSigner(firstMsg)
	if signerAddr == "" {
		return next(ctx, tx, simulate)
	}

	// Check if this signer is a session key grantee.
	sessionKey, found := skd.aaKeeper.GetSessionKeyByGrantee(ctx, signerAddr)
	if !found {
		// Not a session key; pass through to normal signature verification.
		return next(ctx, tx, simulate)
	}

	// --- READ-ONLY validation (H-01 fix) ---
	// Validate inline without calling keeper.ValidateSessionKey() to ensure
	// no state is mutated before signature verification completes.

	// Check time-based expiry
	if ctx.BlockTime().After(sessionKey.Expiry) {
		return ctx, types.ErrSessionExpired
	}

	// Check block-height-based expiry
	if sessionKey.ExpiresAt > 0 && ctx.BlockHeight() > sessionKey.ExpiresAt {
		return ctx, types.ErrSessionExpired
	}

	// Accumulate total spend across all messages for limit checking.
	var totalSpend sdk.Coins

	// Validate each message against session key permissions.
	for _, msg := range msgs {
		msgType := sdk.MsgTypeURL(msg)

		// Check that at least one permission allows this message type.
		hasPermission := false
		for _, perm := range sessionKey.Permissions {
			if perm.MsgType == msgType || perm.MsgType == "*" {
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			return ctx, types.ErrPermissionDenied
		}

		// Accumulate spend amount if the message carries value.
		if sendMsg, ok := msg.(interface{ GetAmount() sdk.Coins }); ok {
			totalSpend = totalSpend.Add(sendMsg.GetAmount()...)
		}
	}

	// Check spend limits (read-only: verify the tx would stay within limits
	// but do NOT update the Used field -- that happens in PostHandler).
	if sessionKey.SpendLimit != nil && !sessionKey.SpendLimit.IsZero() &&
		totalSpend != nil && !totalSpend.IsZero() {
		projectedUsed := sessionKey.Used.Add(totalSpend...)
		for _, used := range projectedUsed {
			limit := sessionKey.SpendLimit.AmountOf(used.Denom)
			if !limit.IsZero() && used.Amount.GT(limit) {
				return ctx, types.ErrSessionLimitExceeded
			}
		}
	}

	// Session key is valid for all messages. Emit an event and continue.
	// NOTE: Actual spend tracking (updating Used) happens in SessionKeyPostDecorator
	// via keeper.RecordSessionKeyUsage() after the tx is fully executed.
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"session_key_auth",
		sdk.NewAttribute("session_key", signerAddr),
		sdk.NewAttribute("granter", sessionKey.Granter),
	))

	return next(ctx, tx, simulate)
}

// SessionKeyPostDecorator records session key spend usage after a transaction
// has been successfully executed, ensuring atomicity (H-13).
type SessionKeyPostDecorator struct {
	aaKeeper SessionKeyKeeper
}

// NewSessionKeyPostDecorator returns a new SessionKeyPostDecorator.
func NewSessionKeyPostDecorator(aaKeeper SessionKeyKeeper) SessionKeyPostDecorator {
	return SessionKeyPostDecorator{aaKeeper: aaKeeper}
}

// PostHandle records spend usage for session keys after successful tx execution.
func (spd SessionKeyPostDecorator) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, success bool, next sdk.PostHandler) (sdk.Context, error) {
	if !success {
		return next(ctx, tx, simulate, success)
	}

	msgs := tx.GetMsgs()
	if len(msgs) == 0 {
		return next(ctx, tx, simulate, success)
	}

	// Determine the signer from the first message's signer field.
	firstMsg := msgs[0]
	signerAddr := extractMsgSigner(firstMsg)
	if signerAddr == "" {
		return next(ctx, tx, simulate, success)
	}

	// Check if this signer is a session key grantee
	sessionKey, found := spd.aaKeeper.GetSessionKeyByGrantee(ctx, signerAddr)
	if !found {
		return next(ctx, tx, simulate, success)
	}

	// Record spend usage for each message
	for _, msg := range msgs {
		var spendAmount sdk.Coins
		if sendMsg, ok := msg.(interface{ GetAmount() sdk.Coins }); ok {
			spendAmount = sendMsg.GetAmount()
		}
		if spendAmount != nil && !spendAmount.IsZero() {
			if err := spd.aaKeeper.RecordSessionKeyUsage(ctx, sessionKey.Key, spendAmount); err != nil {
				return ctx, err
			}
		}
	}

	return next(ctx, tx, simulate, success)
}

// extractMsgSigner extracts the signer address string from a message by
// checking for known abstractaccount message types. This replaces the
// deprecated GetSigners() method removed in SDK v0.53.
func extractMsgSigner(msg sdk.Msg) string {
	switch m := msg.(type) {
	case *types.MsgCreateSmartAccount:
		return m.Sender
	case *types.MsgCreateSessionKey:
		return m.Granter
	case *types.MsgRevokeSessionKey:
		return m.Granter
	case *types.MsgInitiateRecovery:
		return m.Guardian
	case *types.MsgApproveRecovery:
		return m.Guardian
	case *types.MsgExecuteRecovery:
		return m.Sender
	case *types.MsgSponsorGas:
		return m.Sponsor
	case *types.MsgBatchExecute:
		return m.Sender
	default:
		// For non-abstractaccount messages, try common signer field interfaces.
		if m, ok := msg.(interface{ GetSender() string }); ok {
			return m.GetSender()
		}
		if m, ok := msg.(interface{ GetCreator() string }); ok {
			return m.GetCreator()
		}
		return ""
	}
}
