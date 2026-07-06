package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/identity/types"
)

type Keeper struct {
	cdc           codec.Codec
	storeService  store.KVStoreService
	accountKeeper types.AccountKeeper
	authority     string
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	accountKeeper types.AccountKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeService:  storeService,
		accountKeeper: accountKeeper,
		authority:     authority,
	}
}

func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetAuthority returns the module's authority address.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// --- Params ---

func (k Keeper) GetParams(ctx context.Context) types.Params {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.ParamsKey))
	if err != nil || bz == nil {
		return types.DefaultParams()
	}
	var params types.Params
	if err := json.Unmarshal(bz, &params); err != nil {
		return types.DefaultParams()
	}
	return params
}

func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return kvStore.Set([]byte(types.ParamsKey), bz)
}

// --- Identity Store ---

func (k Keeper) GetIdentity(ctx context.Context, address string) (types.Identity, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.IdentityStoreKey(address))
	if err != nil || bz == nil {
		return types.Identity{}, false
	}
	var identity types.Identity
	if err := json.Unmarshal(bz, &identity); err != nil {
		return types.Identity{}, false
	}
	return identity, true
}

func (k Keeper) SetIdentity(ctx context.Context, identity types.Identity) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(identity)
	kvStore.Set(types.IdentityStoreKey(identity.Address), bz)
}

func (k Keeper) GetAllIdentities(ctx context.Context) []types.Identity {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.IdentityKeyPrefix)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var identities []types.Identity
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		var identity types.Identity
		if err := json.Unmarshal(iter.Value(), &identity); err != nil {
			continue
		}
		identities = append(identities, identity)
	}
	return identities
}

func (k Keeper) GetIdentitiesByVerifier(ctx context.Context, verifier string) []types.Identity {
	all := k.GetAllIdentities(ctx)
	var result []types.Identity
	for _, id := range all {
		if id.Verifier == verifier {
			result = append(result, id)
		}
	}
	return result
}

// --- Verifier Store ---

func (k Keeper) GetVerifier(ctx context.Context, address string) (types.Verifier, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.VerifierStoreKey(address))
	if err != nil || bz == nil {
		return types.Verifier{}, false
	}
	var verifier types.Verifier
	if err := json.Unmarshal(bz, &verifier); err != nil {
		return types.Verifier{}, false
	}
	return verifier, true
}

func (k Keeper) SetVerifier(ctx context.Context, verifier types.Verifier) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(verifier)
	kvStore.Set(types.VerifierStoreKey(verifier.Address), bz)
}

func (k Keeper) GetAllVerifiers(ctx context.Context) []types.Verifier {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.VerifierKeyPrefix)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var verifiers []types.Verifier
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		var verifier types.Verifier
		if err := json.Unmarshal(iter.Value(), &verifier); err != nil {
			continue
		}
		verifiers = append(verifiers, verifier)
	}
	return verifiers
}

// --- Business Logic ---

// levelRank returns a numeric rank for comparison. Higher rank = more verified.
func levelRank(level types.VerificationLevel) int {
	switch level {
	case types.VerificationBasic:
		return 1
	case types.VerificationStandard:
		return 2
	case types.VerificationEnhanced:
		return 3
	default:
		return 0
	}
}

// CreateIdentity registers a new identity for verification.
func (k *Keeper) CreateIdentity(ctx context.Context, msg *types.MsgRegisterIdentity) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(ctx)

	if !params.EnableIdentity {
		return fmt.Errorf("identity module is disabled")
	}

	// Check if identity already exists
	_, exists := k.GetIdentity(ctx, msg.Address)
	if exists {
		return types.ErrIdentityExists
	}

	identity := types.Identity{
		Address:        msg.Address,
		Level:          types.VerificationLevel(msg.Level),
		Status:         types.StatusPending,
		DocumentHash:   msg.DocumentHash,
		Nationality:    msg.Nationality,
		CreatedAtBlock: sdkCtx.BlockHeight(),
	}

	k.SetIdentity(ctx, identity)
	k.Logger(ctx).Info("identity registered",
		"address", msg.Address, "level", msg.Level, "nationality", msg.Nationality)

	return nil
}

// VerifyIdentity allows an authorized verifier to approve an identity.
func (k *Keeper) VerifyIdentity(ctx context.Context, msg *types.MsgVerifyIdentity) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(ctx)

	// Check verifier exists and is active
	verifier, found := k.GetVerifier(ctx, msg.Verifier)
	if !found {
		return types.ErrNotVerifier
	}
	if !verifier.Active {
		return types.ErrVerifierNotActive
	}

	// Validate the requested level
	requestedLevel := types.VerificationLevel(msg.Level)
	switch requestedLevel {
	case types.VerificationBasic, types.VerificationStandard, types.VerificationEnhanced:
	default:
		return types.ErrInvalidLevel
	}

	// Check verifier is authorized for this level
	if levelRank(requestedLevel) > levelRank(verifier.MaxLevel) {
		return types.ErrLevelTooHigh
	}

	// Get identity
	identity, found := k.GetIdentity(ctx, msg.Address)
	if !found {
		return types.ErrIdentityNotFound
	}

	// Cannot verify a revoked identity
	if identity.Status == types.StatusRevoked {
		return types.ErrIdentityRevoked
	}

	// Cannot downgrade already verified identity
	if identity.Status == types.StatusVerified && levelRank(identity.Level) >= levelRank(requestedLevel) {
		return types.ErrAlreadyVerified
	}

	// Update identity
	identity.Level = requestedLevel
	identity.Status = types.StatusVerified
	identity.Verifier = msg.Verifier
	identity.VerifiedAtBlock = sdkCtx.BlockHeight()
	identity.ExpiresAtBlock = sdkCtx.BlockHeight() + int64(params.VerificationExpiryBlocks)
	identity.RejectionReason = ""

	k.SetIdentity(ctx, identity)

	// Update verifier stats
	verifier.TotalVerified++
	k.SetVerifier(ctx, verifier)

	k.Logger(ctx).Info("identity verified",
		"address", msg.Address, "level", msg.Level, "verifier", msg.Verifier,
		"expires_at", identity.ExpiresAtBlock)

	return nil
}

// RejectIdentity allows a verifier to reject an identity verification request.
func (k *Keeper) RejectIdentity(ctx context.Context, msg *types.MsgRejectIdentity) error {
	// Check verifier exists and is active
	verifier, found := k.GetVerifier(ctx, msg.Verifier)
	if !found {
		return types.ErrNotVerifier
	}
	if !verifier.Active {
		return types.ErrVerifierNotActive
	}

	// Get identity
	identity, found := k.GetIdentity(ctx, msg.Address)
	if !found {
		return types.ErrIdentityNotFound
	}

	// Can only reject pending identities
	if identity.Status != types.StatusPending {
		return fmt.Errorf("can only reject pending identities, current status: %s", identity.Status)
	}

	identity.Status = types.StatusRejected
	identity.RejectionReason = msg.Reason
	identity.Verifier = msg.Verifier

	k.SetIdentity(ctx, identity)

	k.Logger(ctx).Info("identity rejected",
		"address", msg.Address, "verifier", msg.Verifier, "reason", msg.Reason)

	return nil
}

// RevokeIdentity allows the authority (governance) or the original verifier to revoke an identity.
func (k *Keeper) RevokeIdentity(ctx context.Context, msg *types.MsgRevokeIdentity) error {
	// Must be governance authority or the original verifier
	identity, found := k.GetIdentity(ctx, msg.Address)
	if !found {
		return types.ErrIdentityNotFound
	}

	isAuthority := msg.Authority == k.authority
	isVerifier := msg.Authority == identity.Verifier

	if !isAuthority && !isVerifier {
		return fmt.Errorf("only governance authority or original verifier can revoke an identity")
	}

	identity.Status = types.StatusRevoked
	identity.RejectionReason = msg.Reason

	k.SetIdentity(ctx, identity)

	k.Logger(ctx).Info("identity revoked",
		"address", msg.Address, "authority", msg.Authority, "reason", msg.Reason)

	return nil
}

// UpdateIdentity allows the identity owner to update their document hash and nationality.
// Resets status to pending for re-verification.
func (k *Keeper) UpdateIdentity(ctx context.Context, msg *types.MsgUpdateIdentity) error {
	identity, found := k.GetIdentity(ctx, msg.Address)
	if !found {
		return types.ErrIdentityNotFound
	}

	// Only the owner can update their identity
	if identity.Address != msg.Address {
		return types.ErrNotIdentityOwner
	}

	// Cannot update a revoked identity
	if identity.Status == types.StatusRevoked {
		return types.ErrIdentityRevoked
	}

	// Update fields if provided
	if msg.DocumentHash != "" {
		identity.DocumentHash = msg.DocumentHash
	}
	if msg.Nationality != "" {
		identity.Nationality = msg.Nationality
	}

	// Reset to pending for re-verification
	identity.Status = types.StatusPending
	identity.VerifiedAtBlock = 0
	identity.ExpiresAtBlock = 0
	identity.Verifier = ""
	identity.RejectionReason = ""

	k.SetIdentity(ctx, identity)

	k.Logger(ctx).Info("identity updated, reset to pending",
		"address", msg.Address)

	return nil
}

// RegisterVerifier registers a new authorized verifier. Only governance authority can do this.
func (k *Keeper) RegisterVerifier(ctx context.Context, msg *types.MsgRegisterVerifier) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if msg.Authority != k.authority {
		return fmt.Errorf("unauthorized: only governance authority can register verifiers")
	}

	// Check if verifier already exists
	_, exists := k.GetVerifier(ctx, msg.Verifier)
	if exists {
		return types.ErrVerifierExists
	}

	maxLevel := types.VerificationLevel(msg.MaxLevel)
	switch maxLevel {
	case types.VerificationBasic, types.VerificationStandard, types.VerificationEnhanced:
	default:
		return types.ErrInvalidLevel
	}

	verifier := types.Verifier{
		Address:      msg.Verifier,
		Name:         msg.Name,
		MaxLevel:     maxLevel,
		Active:       true,
		RegisteredAt: sdkCtx.BlockHeight(),
	}

	k.SetVerifier(ctx, verifier)

	k.Logger(ctx).Info("verifier registered",
		"address", msg.Verifier, "name", msg.Name, "max_level", msg.MaxLevel)

	return nil
}

// DeactivateVerifier deactivates a verifier. Only governance authority can do this.
func (k *Keeper) DeactivateVerifier(ctx context.Context, msg *types.MsgDeactivateVerifier) error {
	if msg.Authority != k.authority {
		return fmt.Errorf("unauthorized: only governance authority can deactivate verifiers")
	}

	verifier, found := k.GetVerifier(ctx, msg.Verifier)
	if !found {
		return types.ErrNotVerifier
	}

	verifier.Active = false
	k.SetVerifier(ctx, verifier)

	k.Logger(ctx).Info("verifier deactivated", "address", msg.Verifier)

	return nil
}

// IncrementBookings increments the total_bookings count and updates trust_score.
// Called by travelescrow module when a booking is completed.
func (k *Keeper) IncrementBookings(ctx context.Context, msg *types.MsgIncrementBookings) error {
	if msg.Authority != k.authority {
		return fmt.Errorf("unauthorized: only authority can increment bookings")
	}

	identity, found := k.GetIdentity(ctx, msg.Address)
	if !found {
		// If identity doesn't exist, create a minimal one
		identity = types.Identity{
			Address: msg.Address,
			Level:   types.VerificationNone,
			Status:  types.StatusPending,
		}
	}

	identity.TotalBookings++

	// Trust score = total_bookings * 10, capped at 1000
	score := identity.TotalBookings * 10
	if score > 1000 {
		score = 1000
	}
	identity.TrustScore = score

	k.SetIdentity(ctx, identity)

	k.Logger(ctx).Info("bookings incremented",
		"address", msg.Address, "total_bookings", identity.TotalBookings, "trust_score", identity.TrustScore)

	return nil
}

// IsVerified returns true if the address has a verified identity that is not expired or revoked.
func (k Keeper) IsVerified(ctx context.Context, address string) bool {
	identity, found := k.GetIdentity(ctx, address)
	if !found {
		return false
	}
	if identity.Status != types.StatusVerified {
		return false
	}
	// Check expiry
	if identity.ExpiresAtBlock > 0 {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		if sdkCtx.BlockHeight() >= identity.ExpiresAtBlock {
			return false
		}
	}
	return true
}

// maxExpiryScanPerBlock bounds how many identities are expired in a single block.
// ExpireIdentities does a full scan of the identity keyspace every block with no
// expiry index; that keyspace is attacker-growable, so an unbounded scan lets
// per-block cost grow forever until the chain misses timeout_commit and stalls.
// Capping the number of expirations per block bounds worst-case work; identities
// beyond the cap are expired over subsequent blocks (expiry a few blocks late is
// acceptable — once expired they no longer match the filter). Deterministic:
// store iteration order is identical on all validators. Long-term fix: add a
// height-bucketed expiry index so only identities due at/before the current
// height are visited.
const maxExpiryScanPerBlock = 500

// ExpireIdentities checks identities and marks expired ones (BeginBlocker).
func (k *Keeper) ExpireIdentities(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.IdentityKeyPrefix)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return
	}
	defer iter.Close()

	var toExpire []types.Identity
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		var identity types.Identity
		if err := json.Unmarshal(iter.Value(), &identity); err != nil {
			continue
		}

		// Expire verified identities past their expiry block
		if identity.Status == types.StatusVerified && identity.ExpiresAtBlock > 0 && sdkCtx.BlockHeight() >= identity.ExpiresAtBlock {
			toExpire = append(toExpire, identity)
			// Bound per-block work: stop collecting once the cap is reached.
			// Remaining expired identities are handled in subsequent blocks.
			if len(toExpire) >= maxExpiryScanPerBlock {
				break
			}
		}
	}

	for _, identity := range toExpire {
		identity.Status = types.StatusExpired
		k.SetIdentity(ctx, identity)
		k.Logger(ctx).Info("identity expired", "address", identity.Address, "expired_at_block", sdkCtx.BlockHeight())
	}
}

// --- Helpers ---

func prefixEndBytes(prefix []byte) []byte {
	if len(prefix) == 0 {
		return nil
	}
	end := make([]byte, len(prefix))
	copy(end, prefix)
	for i := len(end) - 1; i >= 0; i-- {
		end[i]++
		if end[i] != 0 {
			return end
		}
	}
	return nil
}
