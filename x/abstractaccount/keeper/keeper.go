package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/abstractaccount/types"
)

type Keeper struct {
	cdc           codec.Codec
	storeService  store.KVStoreService
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	msgRouter     types.MsgRouter
	authority     string
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	msgRouter types.MsgRouter,
	authority string,
) Keeper {
	return Keeper{
		cdc:           cdc,
		storeService:  storeService,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		msgRouter:     msgRouter,
		authority:     authority,
	}
}

func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ---------------------------------------------------------------------------
// Smart Account operations
// ---------------------------------------------------------------------------

// CreateSmartAccount creates a new smart wallet account
func (k Keeper) CreateSmartAccount(ctx context.Context, msg *types.MsgCreateSmartAccount) (string, error) {
	// Derive a deterministic address for the smart account
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return "", types.ErrInvalidAddress
	}

	// Get a monotonic nonce to ensure uniqueness
	nonce := k.nextAccountNonce(ctx)

	// Create a 20-byte address via SHA256 hash of (sender || accountType || nonce)
	h := sha256.New()
	h.Write([]byte("syreen/abstractaccount/"))
	h.Write(senderAddr)
	h.Write([]byte(msg.AccountType))
	nonceBz := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBz, nonce)
	h.Write(nonceBz)
	smartAddr := sdk.AccAddress(h.Sum(nil)[:20]) // Truncate to 20 bytes
	smartAddrStr := smartAddr.String()

	// Check if account already exists
	if _, found := k.GetSmartAccount(ctx, smartAddrStr); found {
		return "", types.ErrAccountExists
	}

	account := types.SmartAccount{
		Address:     smartAddrStr,
		AccountType: msg.AccountType,
		Owners:      msg.Owners,
		Threshold:   msg.Threshold,
	}

	// Ensure the account exists in the auth module
	if k.accountKeeper.GetAccount(ctx, smartAddr) == nil {
		newAcc := k.accountKeeper.NewAccountWithAddress(ctx, smartAddr)
		k.accountKeeper.SetAccount(ctx, newAcc)
	}

	k.SetSmartAccount(ctx, account)

	// If it's a social account type, auto-create a recovery config with sender as initial guardian
	if msg.AccountType == types.AccountTypeSocial {
		params := k.GetParams(ctx)
		config := types.RecoveryConfig{
			Account:     smartAddrStr,
			Guardians:   []string{senderAddr.String()},
			Threshold:   1,
			DelayPeriod: params.RecoveryDelayPeriod,
		}
		k.SetRecoveryConfig(ctx, config)
	}

	k.Logger(ctx).Info("created smart account",
		"address", smartAddrStr,
		"type", msg.AccountType,
		"owners", len(msg.Owners),
		"threshold", msg.Threshold,
	)

	return smartAddrStr, nil
}

// ---------------------------------------------------------------------------
// Session Key operations
// ---------------------------------------------------------------------------

// ValidateSessionKey performs a READ-ONLY check of whether a session key is
// authorized for a given message type and whether the spend would exceed limits.
// H-13: This does NOT persist any Used update. The spend tracking is recorded
// separately via RecordSessionKeyUsage after the transaction succeeds, ensuring
// atomicity — a failed tx does not consume spend budget.
func (k Keeper) ValidateSessionKey(ctx context.Context, key string, msgType string, spendAmount sdk.Coins) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// We need to find the session key - iterate by looking up the key address
	// The key is stored by granter/grantee, so we need to check by grantee
	sessionKey, found := k.GetSessionKeyByGrantee(ctx, key)
	if !found {
		return types.ErrSessionKeyNotFound
	}

	// Check time-based expiry
	if sdkCtx.BlockTime().After(sessionKey.Expiry) {
		return types.ErrSessionExpired
	}

	// Check block-height-based expiry
	if sessionKey.ExpiresAt > 0 && sdkCtx.BlockHeight() > sessionKey.ExpiresAt {
		return fmt.Errorf("session key expired at block %d, current block %d", sessionKey.ExpiresAt, sdkCtx.BlockHeight())
	}

	// Check permissions
	hasPermission := false
	for _, perm := range sessionKey.Permissions {
		if perm.MsgType == msgType || perm.MsgType == "*" {
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		return types.ErrPermissionDenied
	}

	// Enforce spend limits (read-only: check only, do NOT write Used)
	if sessionKey.SpendLimit != nil && !sessionKey.SpendLimit.IsZero() && spendAmount != nil && !spendAmount.IsZero() {
		newUsed := sessionKey.Used.Add(spendAmount...)
		// Check each denom against the limit
		for _, used := range newUsed {
			limit := sessionKey.SpendLimit.AmountOf(used.Denom)
			if !limit.IsZero() && used.Amount.GT(limit) {
				return types.ErrSessionLimitExceeded
			}
		}
	}

	return nil
}

// RecordSessionKeyUsage persists the spend amount against a session key's Used
// budget. This should be called AFTER the transaction has successfully executed
// (e.g., in a PostHandler) to ensure atomicity with tx execution.
func (k Keeper) RecordSessionKeyUsage(ctx context.Context, key string, spendAmount sdk.Coins) error {
	sessionKey, found := k.GetSessionKeyByGrantee(ctx, key)
	if !found {
		return types.ErrSessionKeyNotFound
	}

	if sessionKey.SpendLimit != nil && !sessionKey.SpendLimit.IsZero() && spendAmount != nil && !spendAmount.IsZero() {
		sessionKey.Used = sessionKey.Used.Add(spendAmount...)
		k.SetSessionKey(ctx, sessionKey)
	}

	return nil
}

// CreateSessionKey creates a new session key with scoped permissions.
// The granter is the smart account address; it must list itself (or the tx signer)
// as an owner. The semantic is "create a session key on MY smart account" —
// only the account itself (via its owners) can grant session keys.
// C-10 note: This is NOT an authorization bypass. The granter address IS the
// smart account address, and we verify granter is in that account's Owners list.
// The msg is signed by the granter, so only an owner of the smart account can
// submit this transaction. There is no separate "target account" parameter
// because the granter always operates on its own account.
func (k Keeper) CreateSessionKey(ctx context.Context, granter, grantee string, perms []types.Permission, duration time.Duration) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Validate granter is an owner of a smart account
	account, found := k.GetSmartAccount(ctx, granter)
	if !found {
		return types.ErrAccountNotFound
	}

	// Verify that the granter is actually an owner of this account
	isOwner := false
	for _, owner := range account.Owners {
		if owner == granter {
			isOwner = true
			break
		}
	}
	if !isOwner {
		return fmt.Errorf("granter %s is not an owner of account %s", granter, account.Address)
	}

	// Restrict wildcard permissions — require spend limits when using wildcard
	for _, perm := range perms {
		if perm.MsgType == "*" && (perm.MaxAmount == nil || perm.MaxAmount.IsZero()) {
			return fmt.Errorf("wildcard (*) permissions require a spend limit (MaxAmount)")
		}
	}

	// Check duration against params
	params := k.GetParams(ctx)
	if duration > params.MaxSessionKeyDuration {
		return types.ErrSessionKeyDuration
	}

	expiry := sdkCtx.BlockTime().Add(duration)

	// Calculate total spend limit from permissions
	var spendLimit sdk.Coins
	for _, perm := range perms {
		if perm.MaxAmount != nil {
			spendLimit = spendLimit.Add(perm.MaxAmount...)
		}
	}

	sessionKey := types.SessionKey{
		Key:         grantee,
		Granter:     granter,
		Permissions: perms,
		Expiry:      expiry,
		SpendLimit:  spendLimit,
		Used:        sdk.NewCoins(),
	}

	k.SetSessionKey(ctx, sessionKey)

	k.Logger(ctx).Info("created session key",
		"granter", granter,
		"grantee", grantee,
		"expiry", expiry,
		"permissions", len(perms),
	)

	return nil
}

// RevokeSessionKey revokes a session key
func (k Keeper) RevokeSessionKey(ctx context.Context, granter, keyAddr string) error {
	// Verify granter is an owner of the smart account
	account, acctFound := k.GetSmartAccount(ctx, granter)
	if !acctFound {
		return types.ErrAccountNotFound
	}
	isOwner := false
	for _, owner := range account.Owners {
		if owner == granter {
			isOwner = true
			break
		}
	}
	if !isOwner {
		return fmt.Errorf("granter %s is not an owner of account %s", granter, account.Address)
	}

	_, found := k.GetSessionKey(ctx, granter, keyAddr)
	if !found {
		return types.ErrSessionKeyNotFound
	}

	k.DeleteSessionKey(ctx, granter, keyAddr)

	k.Logger(ctx).Info("revoked session key", "granter", granter, "key", keyAddr)
	return nil
}

// ---------------------------------------------------------------------------
// Social Recovery operations
// ---------------------------------------------------------------------------

// InitiateRecovery starts a recovery process for an account
func (k Keeper) InitiateRecovery(ctx context.Context, guardian, account string, newOwners []string) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Check recovery config exists
	config, found := k.GetRecoveryConfig(ctx, account)
	if !found {
		return types.ErrRecoveryConfigNotFound
	}

	// Verify guardian is in the guardians list
	isGuardian := false
	for _, g := range config.Guardians {
		if g == guardian {
			isGuardian = true
			break
		}
	}
	if !isGuardian {
		return types.ErrNotGuardian
	}

	// Check no existing recovery in progress
	if _, found := k.GetRecoveryRequest(ctx, account); found {
		return types.ErrRecoveryInProgress
	}

	request := types.RecoveryRequest{
		Account:     account,
		NewOwners:   newOwners,
		Approvals:   []string{guardian},
		InitiatedAt: sdkCtx.BlockTime(),
	}

	k.SetRecoveryRequest(ctx, request)

	k.Logger(ctx).Info("initiated recovery",
		"guardian", guardian,
		"account", account,
		"new_owners", len(newOwners),
	)

	return nil
}

// ApproveRecovery adds a guardian approval to an active recovery request
func (k Keeper) ApproveRecovery(ctx context.Context, guardian, account string) error {
	// Check recovery config exists
	config, found := k.GetRecoveryConfig(ctx, account)
	if !found {
		return types.ErrRecoveryConfigNotFound
	}

	// Verify guardian
	isGuardian := false
	for _, g := range config.Guardians {
		if g == guardian {
			isGuardian = true
			break
		}
	}
	if !isGuardian {
		return types.ErrNotGuardian
	}

	// Get active recovery request
	request, found := k.GetRecoveryRequest(ctx, account)
	if !found {
		return types.ErrRecoveryNotFound
	}

	// Check for duplicate approval
	for _, approval := range request.Approvals {
		if approval == guardian {
			return types.ErrDuplicateApproval
		}
	}

	request.Approvals = append(request.Approvals, guardian)
	k.SetRecoveryRequest(ctx, request)

	k.Logger(ctx).Info("approved recovery",
		"guardian", guardian,
		"account", account,
		"approvals", len(request.Approvals),
	)

	return nil
}

// ExecuteRecovery executes a recovery after threshold approvals and delay period.
// The sender must be one of the guardians who approved the recovery.
func (k Keeper) ExecuteRecovery(ctx context.Context, sender string, account string) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	config, found := k.GetRecoveryConfig(ctx, account)
	if !found {
		return types.ErrRecoveryConfigNotFound
	}

	request, found := k.GetRecoveryRequest(ctx, account)
	if !found {
		return types.ErrRecoveryNotFound
	}

	// H-09: Validate that EVERY approval is from a registered guardian.
	// This prevents attackers from stuffing fake approvals into the list.
	for _, approval := range request.Approvals {
		isGuardian := false
		for _, guardian := range config.Guardians {
			if approval == guardian {
				isGuardian = true
				break
			}
		}
		if !isGuardian {
			return fmt.Errorf("approval from %s is not a registered guardian", approval)
		}
	}

	// Verify the sender is one of the approving guardians
	senderIsGuardian := false
	for _, approval := range request.Approvals {
		if approval == sender {
			senderIsGuardian = true
			break
		}
	}
	if !senderIsGuardian {
		return fmt.Errorf("sender %s is not an approving guardian for this recovery", sender)
	}

	// Check threshold
	if uint32(len(request.Approvals)) < config.Threshold {
		return types.ErrInsufficientGuardians
	}

	// Check delay period
	if sdkCtx.BlockTime().Before(request.InitiatedAt.Add(config.DelayPeriod)) {
		return types.ErrRecoveryDelayNotMet
	}

	// Execute: swap owners on the smart account
	smartAccount, found := k.GetSmartAccount(ctx, account)
	if !found {
		return types.ErrAccountNotFound
	}

	// Validate that the new owner set can still meet the threshold requirement.
	// If fewer new owners than the current threshold, reset threshold to avoid bricking the account.
	if uint32(len(request.NewOwners)) < smartAccount.Threshold {
		smartAccount.Threshold = uint32(len(request.NewOwners))
		if smartAccount.Threshold == 0 {
			return fmt.Errorf("recovery would result in zero owners, which bricks the account")
		}
	}

	smartAccount.Owners = request.NewOwners
	k.SetSmartAccount(ctx, smartAccount)

	// Clean up the recovery request
	k.DeleteRecoveryRequest(ctx, account)

	k.Logger(ctx).Info("executed recovery",
		"sender", sender,
		"account", account,
		"new_owners", len(request.NewOwners),
	)

	return nil
}

// CancelRecovery allows any current account owner to cancel an active recovery request.
// This provides a defense mechanism if a recovery is initiated maliciously.
func (k Keeper) CancelRecovery(ctx context.Context, sender string, account string) error {
	// Verify the smart account exists
	smartAccount, found := k.GetSmartAccount(ctx, account)
	if !found {
		return types.ErrAccountNotFound
	}

	// Verify the sender is a current owner of the account
	isOwner := false
	for _, owner := range smartAccount.Owners {
		if owner == sender {
			isOwner = true
			break
		}
	}
	if !isOwner {
		return fmt.Errorf("sender %s is not an owner of account %s", sender, account)
	}

	// Verify there is an active recovery request
	_, found = k.GetRecoveryRequest(ctx, account)
	if !found {
		return types.ErrRecoveryNotFound
	}

	// Delete the recovery request
	k.DeleteRecoveryRequest(ctx, account)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"recovery_cancelled",
		sdk.NewAttribute("account", account),
		sdk.NewAttribute("cancelled_by", sender),
	))

	k.Logger(ctx).Info("recovery cancelled",
		"sender", sender,
		"account", account,
	)

	return nil
}

// ---------------------------------------------------------------------------
// Gas Sponsorship operations
// ---------------------------------------------------------------------------

// SponsorGas creates a gas sponsorship record.
// I-04: Enforces a per-sponsor limit of 5 active sponsorships and validates
// that the sponsor has a minimum balance to cover the gas limit.
func (k Keeper) SponsorGas(ctx context.Context, sponsor, sponsored string, gasLimit uint64, duration time.Duration) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	params := k.GetParams(ctx)
	if !params.EnableGasSponsorship {
		return types.ErrGasSponsorshipDisabled
	}

	// I-04: Check per-sponsor active sponsorship count (max 5)
	const maxActiveSponsorships = 5
	activeCount := k.countActiveSponsorships(ctx, sponsor)
	if activeCount >= maxActiveSponsorships {
		return fmt.Errorf("sponsor has reached the maximum of %d active sponsorships", maxActiveSponsorships)
	}

	// I-04: Check sponsor has at least enough balance to cover this sponsorship
	// (gasLimit * 50 usyreen per gas unit minimum)
	sponsorAddr, err := sdk.AccAddressFromBech32(sponsor)
	if err != nil {
		return fmt.Errorf("invalid sponsor address: %w", err)
	}
	gasLimitInt := math.NewIntFromUint64(gasLimit)
	minBalanceAmount := gasLimitInt.MulRaw(50)
	minBalance := sdk.NewCoin("usyreen", minBalanceAmount) // 50 usyreen per gas unit minimum
	sponsorBal := k.bankKeeper.GetBalance(ctx, sponsorAddr, "usyreen")
	if sponsorBal.Amount.LT(minBalance.Amount) {
		return fmt.Errorf("sponsor has insufficient balance: need at least %s, have %s", minBalance, sponsorBal)
	}

	expiry := sdkCtx.BlockTime().Add(duration)

	gasSponsor := types.GasSponsor{
		Sponsor:        sponsor,
		Sponsored:      sponsored,
		GasLimit:       gasLimit,
		Expiry:         expiry,
		TotalSponsored: sdk.NewCoins(),
	}

	k.SetGasSponsor(ctx, gasSponsor)

	k.Logger(ctx).Info("created gas sponsorship",
		"sponsor", sponsor,
		"sponsored", sponsored,
		"gas_limit", gasLimit,
		"expiry", expiry,
	)

	return nil
}

// countActiveSponsorships counts the number of active (non-expired) sponsorships
// for a given sponsor address.
func (k Keeper) countActiveSponsorships(ctx context.Context, sponsor string) int {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.GasSponsorPrefix + sponsor + "/")

	iter, err := kvStore.Iterator(prefix, append(prefix, 0xFF))
	if err != nil {
		return 0
	}
	defer iter.Close()

	count := 0
	for ; iter.Valid(); iter.Next() {
		var gs types.GasSponsor
		if err := json.Unmarshal(iter.Value(), &gs); err != nil {
			continue
		}
		if sdkCtx.BlockTime().Before(gs.Expiry) {
			count++
		}
	}
	return count
}

// CheckGasSponsor returns a sponsor for the given account if one exists, is active,
// and has remaining gas budget.
func (k Keeper) CheckGasSponsor(ctx context.Context, sender string, gasWanted uint64) (*types.GasSponsor, bool) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Look up sponsors for this account by iterating the store
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.GasSponsorsBySponsored(sender)

	// Check the reverse index
	iter, err := kvStore.Iterator(prefix, append(prefix, 0xFF))
	if err != nil {
		return nil, false
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		sponsorAddr := string(iter.Value())
		gasSponsor, found := k.GetGasSponsor(ctx, sponsorAddr, sender)
		if !found {
			continue
		}
		// Check if still active
		if sdkCtx.BlockTime().After(gasSponsor.Expiry) {
			continue
		}
		// I-14: Check gas limit using uint64 tracking instead of fake "gas" denom
		totalUsed := gasSponsor.TotalGasUsed
		if totalUsed+gasWanted > gasSponsor.GasLimit {
			continue // This sponsor is exhausted
		}
		return &gasSponsor, true
	}

	return nil, false
}

// DeductGasSponsorship records gas usage against a sponsor's budget.
// I-14: Uses uint64 counter instead of fake "gas" denom.
func (k Keeper) DeductGasSponsorship(ctx context.Context, sponsor *types.GasSponsor, gasUsed uint64) {
	sponsor.TotalGasUsed += gasUsed
	k.SetGasSponsor(ctx, *sponsor)
}

// ---------------------------------------------------------------------------
// Batch Execution
// ---------------------------------------------------------------------------

// ExecuteBatch executes multiple messages atomically using CacheContext.
// All messages succeed or all are rolled back.
func (k Keeper) ExecuteBatch(ctx context.Context, sender string, messages []json.RawMessage) ([][]byte, error) {
	params := k.GetParams(ctx)
	if uint32(len(messages)) > params.MaxBatchSize {
		return nil, types.ErrBatchTooLarge
	}

	// Verify sender is a smart account owner
	_, found := k.GetSmartAccount(ctx, sender)
	if !found {
		return nil, types.ErrAccountNotFound
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	k.Logger(ctx).Info("executing batch",
		"sender", sender,
		"messages", len(messages),
	)

	// I-13: Allowlist approach — only permit known-safe message types in batch execution.
	// Any message type NOT in this list is rejected.
	allowedBatchMsgTypes := map[string]bool{
		"/cosmos.bank.v1beta1.MsgSend":                  true,
		"/cosmos.bank.v1beta1.MsgMultiSend":             true,
		"/cosmos.staking.v1beta1.MsgDelegate":           true,
		"/cosmos.staking.v1beta1.MsgUndelegate":         true,
		"/cosmos.staking.v1beta1.MsgBeginRedelegate":    true,
		"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward": true,
		"/syreen.dex.MsgSwap":                           true,
		"/syreen.dex.MsgPlaceOrder":                     true,
		"/syreen.dex.MsgCancelOrder":                    true,
		"/syreen.dex.MsgModifyOrder":                    true,
		"/syreen.dex.MsgAddLiquidity":                   true,
		"/syreen.dex.MsgRemoveLiquidity":                true,
	}

	// Use CacheContext for atomic execution: all-or-nothing.
	cacheCtx, writeFn := sdkCtx.CacheContext()

	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return nil, types.ErrInvalidAddress
	}

	results := make([][]byte, 0, len(messages))

	for i, rawMsg := range messages {
		// Decode the message using the codec's JSON unmarshaler.
		// The raw JSON must contain a "@type" field for proper deserialization.
		var anyMsg struct {
			Type string `json:"@type"`
		}
		if err := json.Unmarshal(rawMsg, &anyMsg); err != nil {
			return nil, fmt.Errorf("batch message %d: failed to parse @type: %w", i, err)
		}
		if anyMsg.Type == "" {
			return nil, fmt.Errorf("batch message %d: missing @type field", i)
		}

		// I-13: Check message type against allowlist — reject anything not explicitly permitted
		if !allowedBatchMsgTypes[anyMsg.Type] {
			return nil, fmt.Errorf("batch message %d: message type %s not allowed in batch execution", i, anyMsg.Type)
		}

		// Resolve the message type from the codec's interface registry.
		msgIface, err := k.cdc.InterfaceRegistry().Resolve(anyMsg.Type)
		if err != nil {
			return nil, fmt.Errorf("batch message %d: unknown message type %s: %w", i, anyMsg.Type, err)
		}

		sdkMsg, ok := msgIface.(sdk.Msg)
		if !ok {
			return nil, fmt.Errorf("batch message %d: resolved type %s does not implement sdk.Msg", i, anyMsg.Type)
		}

		// Unmarshal the full message using encoding/json since the custom Msg
		// types in this chain use JSON struct tags rather than protobuf encoding.
		if err := json.Unmarshal(rawMsg, sdkMsg); err != nil {
			return nil, fmt.Errorf("batch message %d: failed to unmarshal: %w", i, err)
		}

		// Validate basic message constraints.
		if hasValidate, ok := sdkMsg.(interface{ ValidateBasic() error }); ok {
			if err := hasValidate.ValidateBasic(); err != nil {
				return nil, fmt.Errorf("batch message %d: validation failed: %w", i, err)
			}
		}

		// Verify the sender is a signer of this message (authorization check).
		if hasSigners, ok := sdkMsg.(interface{ GetSigners() []sdk.AccAddress }); ok {
			signers := hasSigners.GetSigners()
			senderIsSigner := false
			for _, signer := range signers {
				if signer.Equals(senderAddr) {
					senderIsSigner = true
					break
				}
			}
			if !senderIsSigner {
				return nil, fmt.Errorf("batch message %d: sender %s is not a signer of this message", i, sender)
			}
		}

		// Route the message to its handler.
		handler := k.msgRouter.Handler(sdkMsg)
		if handler == nil {
			return nil, fmt.Errorf("batch message %d: no handler found for %s", i, anyMsg.Type)
		}

		// Execute the message within the cached context.
		msgResp, err := handler(cacheCtx, sdkMsg)
		if err != nil {
			return nil, fmt.Errorf("batch message %d: execution failed: %w", i, err)
		}

		// Collect the response bytes if available.
		if msgResp != nil {
			bz, err := k.cdc.MarshalJSON(msgResp)
			if err != nil {
				results = append(results, nil)
			} else {
				results = append(results, bz)
			}
		} else {
			results = append(results, nil)
		}
	}

	// All messages succeeded; commit the cached writes.
	writeFn()

	return results, nil
}

// ---------------------------------------------------------------------------
// Store helpers
// ---------------------------------------------------------------------------

// GetSmartAccount retrieves a smart account from the store
func (k Keeper) GetSmartAccount(ctx context.Context, address string) (types.SmartAccount, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.SmartAccountKey(address))
	if err != nil || bz == nil {
		return types.SmartAccount{}, false
	}
	var account types.SmartAccount
	if err := json.Unmarshal(bz, &account); err != nil {
		return types.SmartAccount{}, false
	}
	return account, true
}

// SetSmartAccount stores a smart account
func (k Keeper) SetSmartAccount(ctx context.Context, account types.SmartAccount) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(account)
	kvStore.Set(types.SmartAccountKey(account.Address), bz)
}

// GetSessionKey retrieves a session key by granter and grantee
func (k Keeper) GetSessionKey(ctx context.Context, granter, grantee string) (types.SessionKey, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.SessionKeyKey(granter, grantee))
	if err != nil || bz == nil {
		return types.SessionKey{}, false
	}
	var sessionKey types.SessionKey
	if err := json.Unmarshal(bz, &sessionKey); err != nil {
		return types.SessionKey{}, false
	}
	return sessionKey, true
}

// GetSessionKeyByGrantee searches for a session key by the grantee address.
// FIX: Only returns session keys where the grantee matches AND the key is not expired.
// If multiple matches exist, returns the one with the fewest permissions (most restrictive)
// to prevent hijacking via permissive session keys from other granters.
func (k Keeper) GetSessionKeyByGrantee(ctx context.Context, grantee string) (types.SessionKey, bool) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.SessionKeyPrefix)

	iter, err := kvStore.Iterator(prefix, append(prefix, 0xFF))
	if err != nil {
		return types.SessionKey{}, false
	}
	defer iter.Close()

	var bestKey types.SessionKey
	found := false
	for ; iter.Valid(); iter.Next() {
		var sessionKey types.SessionKey
		if err := json.Unmarshal(iter.Value(), &sessionKey); err != nil {
			continue
		}
		if sessionKey.Key != grantee {
			continue
		}
		// Skip expired keys
		if !sessionKey.Expiry.IsZero() && sdkCtx.BlockTime().After(sessionKey.Expiry) {
			continue
		}
		if sessionKey.ExpiresAt > 0 && sdkCtx.BlockHeight() > sessionKey.ExpiresAt {
			continue
		}
		if !found {
			bestKey = sessionKey
			found = true
		} else if len(sessionKey.Permissions) < len(bestKey.Permissions) {
			// Prefer the most restrictive key (fewest permissions)
			bestKey = sessionKey
		}
	}
	return bestKey, found
}

// SetSessionKey stores a session key
func (k Keeper) SetSessionKey(ctx context.Context, sessionKey types.SessionKey) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(sessionKey)
	kvStore.Set(types.SessionKeyKey(sessionKey.Granter, sessionKey.Key), bz)
}

// DeleteSessionKey removes a session key from the store
func (k Keeper) DeleteSessionKey(ctx context.Context, granter, grantee string) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Delete(types.SessionKeyKey(granter, grantee))
}

// GetRecoveryConfig retrieves the recovery configuration for an account
func (k Keeper) GetRecoveryConfig(ctx context.Context, account string) (types.RecoveryConfig, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.RecoveryConfigKey(account))
	if err != nil || bz == nil {
		return types.RecoveryConfig{}, false
	}
	var config types.RecoveryConfig
	if err := json.Unmarshal(bz, &config); err != nil {
		return types.RecoveryConfig{}, false
	}
	return config, true
}

// SetRecoveryConfig stores a recovery configuration
func (k Keeper) SetRecoveryConfig(ctx context.Context, config types.RecoveryConfig) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(config)
	kvStore.Set(types.RecoveryConfigKey(config.Account), bz)
}

// GetRecoveryRequest retrieves an active recovery request
func (k Keeper) GetRecoveryRequest(ctx context.Context, account string) (types.RecoveryRequest, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.RecoveryRequestKey(account))
	if err != nil || bz == nil {
		return types.RecoveryRequest{}, false
	}
	var request types.RecoveryRequest
	if err := json.Unmarshal(bz, &request); err != nil {
		return types.RecoveryRequest{}, false
	}
	return request, true
}

// SetRecoveryRequest stores a recovery request
func (k Keeper) SetRecoveryRequest(ctx context.Context, request types.RecoveryRequest) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(request)
	kvStore.Set(types.RecoveryRequestKey(request.Account), bz)
}

// DeleteRecoveryRequest removes a recovery request from the store
func (k Keeper) DeleteRecoveryRequest(ctx context.Context, account string) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Delete(types.RecoveryRequestKey(account))
}

// GetGasSponsor retrieves a gas sponsorship record
func (k Keeper) GetGasSponsor(ctx context.Context, sponsor, sponsored string) (types.GasSponsor, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.GasSponsorKey(sponsor, sponsored))
	if err != nil || bz == nil {
		return types.GasSponsor{}, false
	}
	var gasSponsor types.GasSponsor
	if err := json.Unmarshal(bz, &gasSponsor); err != nil {
		return types.GasSponsor{}, false
	}
	return gasSponsor, true
}

// SetGasSponsor stores a gas sponsorship record and maintains a reverse index
func (k Keeper) SetGasSponsor(ctx context.Context, gasSponsor types.GasSponsor) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(gasSponsor)
	kvStore.Set(types.GasSponsorKey(gasSponsor.Sponsor, gasSponsor.Sponsored), bz)

	// Maintain reverse index: sponsored -> sponsor
	reverseKey := append(types.GasSponsorsBySponsored(gasSponsor.Sponsored), []byte(gasSponsor.Sponsor)...)
	kvStore.Set(reverseKey, []byte(gasSponsor.Sponsor))
}

// GetParams returns the module parameters
func (k Keeper) GetParams(ctx context.Context) types.Params {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte("params"))
	if err != nil || bz == nil {
		return types.DefaultParams()
	}
	var params types.Params
	if err := json.Unmarshal(bz, &params); err != nil {
		return types.DefaultParams()
	}
	return params
}

// SetParams stores the module parameters
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return kvStore.Set([]byte("params"), bz)
}

// nextAccountNonce returns and increments a monotonic nonce for address derivation.
func (k Keeper) nextAccountNonce(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.AccountNonceKey))
	var nonce uint64
	if err == nil && bz != nil && len(bz) == 8 {
		nonce = binary.BigEndian.Uint64(bz)
	}
	nonce++
	out := make([]byte, 8)
	binary.BigEndian.PutUint64(out, nonce)
	kvStore.Set([]byte(types.AccountNonceKey), out)
	return nonce
}
