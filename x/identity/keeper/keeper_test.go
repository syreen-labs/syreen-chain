package keeper_test

import (
	"testing"

	"context"

	"cosmossdk.io/core/store"

	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/log"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	cosmosstore "cosmossdk.io/store"

	"syreen/x/identity/keeper"
	"syreen/x/identity/types"
)

// --- Mock Keepers ---

type mockAccountKeeper struct{}

func (m mockAccountKeeper) GetAccount(_ context.Context, _ sdk.AccAddress) sdk.AccountI { return nil }

// --- Test Helper ---

const authority = "syreen1governance"

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("syreen", "syreenpub")
	config.SetBech32PrefixForValidator("syreenvaloper", "syreenvaloperpub")
	config.SetBech32PrefixForConsensusNode("syreenvalcons", "syreenvalconspub")
	config.Seal()
}

func setupKeeper(t *testing.T) (*keeper.Keeper, sdk.Context) {
	t.Helper()

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	stateStore := cosmosstore.NewCommitMultiStore(db, log.NewNopLogger(), storemetrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	if err := stateStore.LoadLatestVersion(); err != nil {
		t.Fatal(err)
	}

	ctx := sdk.NewContext(stateStore, cmtproto.Header{Height: 1}, false, log.NewNopLogger())

	var storeService store.KVStoreService = runtime.NewKVStoreService(storeKey)

	k := keeper.NewKeeper(
		codec.NewProtoCodec(nil),
		storeService,
		mockAccountKeeper{},
		authority,
	)

	// Init genesis with default params
	k.InitGenesis(ctx, *types.DefaultGenesis())

	return k, ctx
}

// Test addresses (valid bech32 with syreen prefix)
const (
	user1    = "syreen1z4tljj0tqa8wrh5p856e3vu5afj5ayrxrxqrqq"
	user2    = "syreen1mety9pctf2dngpawmwcp6unyw4a6rgkkkazpqz"
	user3    = "syreen1pgzph9rze2j2xxavx4n7pdhxlkgsq7razhqd7n"
	verifier1Addr = "syreen1j5a77099lqrf33jwccdv7lywkuwcvx66gujjez"
	verifier2Addr = "syreen1z4tljj0tqa8wrh5p856e3vu5afj5ayrxrxqrqq"
)

// --- Params Tests ---

func TestDefaultParams(t *testing.T) {
	k, ctx := setupKeeper(t)

	params := k.GetParams(ctx)
	if !params.EnableIdentity {
		t.Fatal("identity should be enabled by default")
	}
	if params.VerificationExpiryBlocks != 63115200 {
		t.Fatalf("expected expiry blocks = 63115200, got %d", params.VerificationExpiryBlocks)
	}
	if params.MinTrustScore != 0 {
		t.Fatalf("expected min trust score = 0, got %d", params.MinTrustScore)
	}
	if !params.BasicVerificationFree {
		t.Fatal("basic verification should be free by default")
	}

	t.Log("PASS: Default params loaded correctly")
}

func TestSetGetParams(t *testing.T) {
	k, ctx := setupKeeper(t)

	custom := types.Params{
		VerificationExpiryBlocks: 100000,
		MinTrustScore:            50,
		EnableIdentity:           true,
		BasicVerificationFree:    false,
	}
	if err := k.SetParams(ctx, custom); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	got := k.GetParams(ctx)
	if got.VerificationExpiryBlocks != 100000 {
		t.Fatalf("expected expiry blocks = 100000, got %d", got.VerificationExpiryBlocks)
	}
	if got.MinTrustScore != 50 {
		t.Fatalf("expected min trust score = 50, got %d", got.MinTrustScore)
	}
	if got.BasicVerificationFree {
		t.Fatal("expected basic verification free = false")
	}

	t.Log("PASS: Custom params set and retrieved correctly")
}

// --- Register Identity Tests ---

func TestRegisterIdentity(t *testing.T) {
	k, ctx := setupKeeper(t)

	msg := &types.MsgRegisterIdentity{
		Address:      user1,
		DocumentHash: "abc123hash",
		Nationality:  "NL",
		Level:        "basic",
	}

	if err := k.CreateIdentity(ctx, msg); err != nil {
		t.Fatalf("failed to register identity: %v", err)
	}

	identity, found := k.GetIdentity(ctx, user1)
	if !found {
		t.Fatal("identity should be found after registration")
	}
	if identity.Address != user1 {
		t.Fatalf("expected address %s, got %s", user1, identity.Address)
	}
	if identity.Status != types.StatusPending {
		t.Fatalf("expected status pending, got %s", identity.Status)
	}
	if identity.Level != types.VerificationBasic {
		t.Fatalf("expected level basic, got %s", identity.Level)
	}
	if identity.DocumentHash != "abc123hash" {
		t.Fatalf("expected doc hash abc123hash, got %s", identity.DocumentHash)
	}
	if identity.Nationality != "NL" {
		t.Fatalf("expected nationality NL, got %s", identity.Nationality)
	}

	t.Log("PASS: Identity registered correctly")
}

func TestRegisterIdentityDuplicate(t *testing.T) {
	k, ctx := setupKeeper(t)

	msg := &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	}

	if err := k.CreateIdentity(ctx, msg); err != nil {
		t.Fatalf("first registration should succeed: %v", err)
	}

	err := k.CreateIdentity(ctx, msg)
	if err == nil {
		t.Fatal("duplicate registration should fail")
	}
	if err != types.ErrIdentityExists {
		t.Fatalf("expected ErrIdentityExists, got: %v", err)
	}

	t.Log("PASS: Duplicate identity registration rejected")
}

func TestRegisterIdentityDisabled(t *testing.T) {
	k, ctx := setupKeeper(t)

	params := k.GetParams(ctx)
	params.EnableIdentity = false
	k.SetParams(ctx, params)

	msg := &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	}

	err := k.CreateIdentity(ctx, msg)
	if err == nil {
		t.Fatal("registration should fail when module is disabled")
	}

	t.Log("PASS: Identity registration rejected when module disabled")
}

// --- Register Verifier Tests ---

func TestRegisterVerifier(t *testing.T) {
	k, ctx := setupKeeper(t)

	msg := &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC Provider Inc",
		MaxLevel:  "enhanced",
	}

	if err := k.RegisterVerifier(ctx, msg); err != nil {
		t.Fatalf("failed to register verifier: %v", err)
	}

	verifier, found := k.GetVerifier(ctx, verifier1Addr)
	if !found {
		t.Fatal("verifier should be found after registration")
	}
	if verifier.Name != "KYC Provider Inc" {
		t.Fatalf("expected name 'KYC Provider Inc', got '%s'", verifier.Name)
	}
	if verifier.MaxLevel != types.VerificationEnhanced {
		t.Fatalf("expected max level enhanced, got %s", verifier.MaxLevel)
	}
	if !verifier.Active {
		t.Fatal("verifier should be active after registration")
	}

	t.Log("PASS: Verifier registered correctly")
}

func TestRegisterVerifierUnauthorized(t *testing.T) {
	k, ctx := setupKeeper(t)

	msg := &types.MsgRegisterVerifier{
		Authority: user1, // not the governance authority
		Verifier:  verifier1Addr,
		Name:      "Fake Verifier",
		MaxLevel:  "basic",
	}

	err := k.RegisterVerifier(ctx, msg)
	if err == nil {
		t.Fatal("non-authority should not be able to register verifier")
	}

	t.Log("PASS: Unauthorized verifier registration rejected")
}

func TestRegisterVerifierDuplicate(t *testing.T) {
	k, ctx := setupKeeper(t)

	msg := &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC Provider",
		MaxLevel:  "standard",
	}

	if err := k.RegisterVerifier(ctx, msg); err != nil {
		t.Fatalf("first registration should succeed: %v", err)
	}

	err := k.RegisterVerifier(ctx, msg)
	if err == nil {
		t.Fatal("duplicate verifier registration should fail")
	}

	t.Log("PASS: Duplicate verifier registration rejected")
}

func TestRegisterVerifierInvalidLevel(t *testing.T) {
	k, ctx := setupKeeper(t)

	msg := &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "Bad Verifier",
		MaxLevel:  "super_kyc",
	}

	err := k.RegisterVerifier(ctx, msg)
	if err == nil {
		t.Fatal("invalid max level should fail")
	}

	t.Log("PASS: Invalid verifier level rejected")
}

// --- Verify Identity Tests ---

func TestVerifyIdentity(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Register verifier
	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC Provider",
		MaxLevel:  "enhanced",
	})

	// Register identity
	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address:      user1,
		DocumentHash: "doc123",
		Level:        "standard",
	})

	// Verify
	err := k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "standard",
	})
	if err != nil {
		t.Fatalf("failed to verify identity: %v", err)
	}

	identity, _ := k.GetIdentity(ctx, user1)
	if identity.Status != types.StatusVerified {
		t.Fatalf("expected status verified, got %s", identity.Status)
	}
	if identity.Level != types.VerificationStandard {
		t.Fatalf("expected level standard, got %s", identity.Level)
	}
	if identity.Verifier != verifier1Addr {
		t.Fatalf("expected verifier %s, got %s", verifier1Addr, identity.Verifier)
	}
	if identity.VerifiedAtBlock != 1 {
		t.Fatalf("expected verified at block 1, got %d", identity.VerifiedAtBlock)
	}
	if identity.ExpiresAtBlock == 0 {
		t.Fatal("expected expires at block to be set")
	}

	// Check verifier stats
	verifier, _ := k.GetVerifier(ctx, verifier1Addr)
	if verifier.TotalVerified != 1 {
		t.Fatalf("expected total verified = 1, got %d", verifier.TotalVerified)
	}

	t.Log("PASS: Identity verified correctly")
}

func TestVerifyIdentityUnauthorizedVerifier(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	})

	// Try to verify without being a registered verifier
	err := k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "basic",
	})
	if err == nil {
		t.Fatal("unregistered verifier should not be able to verify")
	}
	if err != types.ErrNotVerifier {
		t.Fatalf("expected ErrNotVerifier, got: %v", err)
	}

	t.Log("PASS: Unauthorized verifier rejected")
}

func TestVerifyIdentityLevelTooHigh(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Register verifier with basic max level
	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "Basic KYC",
		MaxLevel:  "basic",
	})

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address:      user1,
		DocumentHash: "doc123",
		Level:        "enhanced",
	})

	// Try to verify at enhanced level (above verifier's max)
	err := k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "enhanced",
	})
	if err == nil {
		t.Fatal("verifier should not be able to verify above their max level")
	}
	if err != types.ErrLevelTooHigh {
		t.Fatalf("expected ErrLevelTooHigh, got: %v", err)
	}

	t.Log("PASS: Level too high rejected")
}

func TestVerifyIdentityNotFound(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	err := k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "basic",
	})
	if err == nil {
		t.Fatal("verifying nonexistent identity should fail")
	}
	if err != types.ErrIdentityNotFound {
		t.Fatalf("expected ErrIdentityNotFound, got: %v", err)
	}

	t.Log("PASS: Verify nonexistent identity rejected")
}

func TestVerifyIdentityAlreadyVerifiedSameLevel(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "standard",
	})

	// First verify
	k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "standard",
	})

	// Try to verify again at same level
	err := k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "standard",
	})
	if err == nil {
		t.Fatal("verifying already-verified identity at same level should fail")
	}

	t.Log("PASS: Already verified at same level rejected")
}

func TestVerifyIdentityUpgrade(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	})

	// Verify at basic
	k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "basic",
	})

	// Upgrade to enhanced
	err := k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "enhanced",
	})
	if err != nil {
		t.Fatalf("upgrading verification level should succeed: %v", err)
	}

	identity, _ := k.GetIdentity(ctx, user1)
	if identity.Level != types.VerificationEnhanced {
		t.Fatalf("expected level enhanced after upgrade, got %s", identity.Level)
	}

	t.Log("PASS: Identity verification upgraded")
}

func TestVerifyInactiveVerifier(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	// Deactivate verifier
	k.DeactivateVerifier(ctx, &types.MsgDeactivateVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
	})

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	})

	err := k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "basic",
	})
	if err == nil {
		t.Fatal("inactive verifier should not be able to verify")
	}
	if err != types.ErrVerifierNotActive {
		t.Fatalf("expected ErrVerifierNotActive, got: %v", err)
	}

	t.Log("PASS: Inactive verifier rejected")
}

// --- Reject Identity Tests ---

func TestRejectIdentity(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address:      user1,
		DocumentHash: "bad_doc",
		Level:        "standard",
	})

	err := k.RejectIdentity(ctx, &types.MsgRejectIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Reason:   "Document unclear",
	})
	if err != nil {
		t.Fatalf("failed to reject identity: %v", err)
	}

	identity, _ := k.GetIdentity(ctx, user1)
	if identity.Status != types.StatusRejected {
		t.Fatalf("expected status rejected, got %s", identity.Status)
	}
	if identity.RejectionReason != "Document unclear" {
		t.Fatalf("expected rejection reason 'Document unclear', got '%s'", identity.RejectionReason)
	}

	t.Log("PASS: Identity rejected correctly")
}

func TestRejectNonPendingIdentity(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	})

	// Verify first
	k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "basic",
	})

	// Try to reject verified identity
	err := k.RejectIdentity(ctx, &types.MsgRejectIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Reason:   "Changed my mind",
	})
	if err == nil {
		t.Fatal("should not be able to reject a non-pending identity")
	}

	t.Log("PASS: Non-pending identity rejection rejected")
}

// --- Revoke Identity Tests ---

func TestRevokeIdentity(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	})

	k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "basic",
	})

	// Revoke by authority
	err := k.RevokeIdentity(ctx, &types.MsgRevokeIdentity{
		Authority: authority,
		Address:   user1,
		Reason:    "Fraudulent documents",
	})
	if err != nil {
		t.Fatalf("failed to revoke identity: %v", err)
	}

	identity, _ := k.GetIdentity(ctx, user1)
	if identity.Status != types.StatusRevoked {
		t.Fatalf("expected status revoked, got %s", identity.Status)
	}

	t.Log("PASS: Identity revoked correctly")
}

func TestRevokeByOriginalVerifier(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	})

	k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "basic",
	})

	// Revoke by the original verifier
	err := k.RevokeIdentity(ctx, &types.MsgRevokeIdentity{
		Authority: verifier1Addr,
		Address:   user1,
		Reason:    "Re-verification needed",
	})
	if err != nil {
		t.Fatalf("original verifier should be able to revoke: %v", err)
	}

	identity, _ := k.GetIdentity(ctx, user1)
	if identity.Status != types.StatusRevoked {
		t.Fatalf("expected status revoked, got %s", identity.Status)
	}

	t.Log("PASS: Identity revoked by original verifier")
}

func TestRevokeByUnauthorized(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	})

	k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "basic",
	})

	// Try to revoke by random user
	err := k.RevokeIdentity(ctx, &types.MsgRevokeIdentity{
		Authority: user2,
		Address:   user1,
		Reason:    "I want to revoke",
	})
	if err == nil {
		t.Fatal("unauthorized user should not be able to revoke")
	}

	t.Log("PASS: Unauthorized revocation rejected")
}

func TestVerifyRevokedIdentity(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	})

	// Revoke directly (set via authority)
	k.RevokeIdentity(ctx, &types.MsgRevokeIdentity{
		Authority: authority,
		Address:   user1,
		Reason:    "Bad actor",
	})

	// Try to verify revoked identity
	err := k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "basic",
	})
	if err == nil {
		t.Fatal("should not be able to verify a revoked identity")
	}

	t.Log("PASS: Verifying revoked identity rejected")
}

// --- Update Identity Tests ---

func TestUpdateIdentity(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address:      user1,
		DocumentHash: "old_doc",
		Nationality:  "NL",
		Level:        "basic",
	})

	// Verify first
	k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "basic",
	})

	// Update
	err := k.UpdateIdentity(ctx, &types.MsgUpdateIdentity{
		Address:      user1,
		DocumentHash: "new_doc_hash",
		Nationality:  "DE",
	})
	if err != nil {
		t.Fatalf("failed to update identity: %v", err)
	}

	identity, _ := k.GetIdentity(ctx, user1)
	if identity.DocumentHash != "new_doc_hash" {
		t.Fatalf("expected new doc hash, got %s", identity.DocumentHash)
	}
	if identity.Nationality != "DE" {
		t.Fatalf("expected nationality DE, got %s", identity.Nationality)
	}
	// Should reset to pending
	if identity.Status != types.StatusPending {
		t.Fatalf("expected status pending after update, got %s", identity.Status)
	}
	if identity.Verifier != "" {
		t.Fatalf("expected verifier to be cleared after update, got %s", identity.Verifier)
	}
	if identity.VerifiedAtBlock != 0 {
		t.Fatalf("expected verified_at_block = 0 after update, got %d", identity.VerifiedAtBlock)
	}

	t.Log("PASS: Identity updated and reset to pending")
}

func TestUpdateRevokedIdentity(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	})

	// Revoke
	k.RevokeIdentity(ctx, &types.MsgRevokeIdentity{
		Authority: authority,
		Address:   user1,
		Reason:    "Bad",
	})

	err := k.UpdateIdentity(ctx, &types.MsgUpdateIdentity{
		Address:      user1,
		DocumentHash: "new_doc",
	})
	if err == nil {
		t.Fatal("should not be able to update a revoked identity")
	}

	t.Log("PASS: Revoked identity update rejected")
}

func TestUpdateNonexistentIdentity(t *testing.T) {
	k, ctx := setupKeeper(t)

	err := k.UpdateIdentity(ctx, &types.MsgUpdateIdentity{
		Address:      user1,
		DocumentHash: "doc",
	})
	if err == nil {
		t.Fatal("should not be able to update nonexistent identity")
	}

	t.Log("PASS: Nonexistent identity update rejected")
}

// --- Deactivate Verifier Tests ---

func TestDeactivateVerifier(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	err := k.DeactivateVerifier(ctx, &types.MsgDeactivateVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
	})
	if err != nil {
		t.Fatalf("failed to deactivate verifier: %v", err)
	}

	verifier, found := k.GetVerifier(ctx, verifier1Addr)
	if !found {
		t.Fatal("verifier should still exist after deactivation")
	}
	if verifier.Active {
		t.Fatal("verifier should be inactive after deactivation")
	}

	t.Log("PASS: Verifier deactivated correctly")
}

func TestDeactivateVerifierUnauthorized(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	err := k.DeactivateVerifier(ctx, &types.MsgDeactivateVerifier{
		Authority: user1,
		Verifier:  verifier1Addr,
	})
	if err == nil {
		t.Fatal("unauthorized deactivation should fail")
	}

	t.Log("PASS: Unauthorized deactivation rejected")
}

// --- IncrementBookings Tests ---

func TestIncrementBookings(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	})

	// Increment a few times
	for i := 0; i < 5; i++ {
		err := k.IncrementBookings(ctx, &types.MsgIncrementBookings{
			Authority: authority,
			Address:   user1,
		})
		if err != nil {
			t.Fatalf("failed to increment bookings: %v", err)
		}
	}

	identity, _ := k.GetIdentity(ctx, user1)
	if identity.TotalBookings != 5 {
		t.Fatalf("expected total_bookings = 5, got %d", identity.TotalBookings)
	}
	if identity.TrustScore != 50 {
		t.Fatalf("expected trust_score = 50 (5*10), got %d", identity.TrustScore)
	}

	t.Log("PASS: Bookings incremented and trust score calculated")
}

func TestIncrementBookingsTrustScoreCap(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	})

	// Increment 150 times (150 * 10 = 1500, should cap at 1000)
	for i := 0; i < 150; i++ {
		k.IncrementBookings(ctx, &types.MsgIncrementBookings{
			Authority: authority,
			Address:   user1,
		})
	}

	identity, _ := k.GetIdentity(ctx, user1)
	if identity.TotalBookings != 150 {
		t.Fatalf("expected total_bookings = 150, got %d", identity.TotalBookings)
	}
	if identity.TrustScore != 1000 {
		t.Fatalf("expected trust_score capped at 1000, got %d", identity.TrustScore)
	}

	t.Log("PASS: Trust score capped at 1000")
}

func TestIncrementBookingsNoIdentity(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Increment on nonexistent identity (should create minimal one)
	err := k.IncrementBookings(ctx, &types.MsgIncrementBookings{
		Authority: authority,
		Address:   user1,
	})
	if err != nil {
		t.Fatalf("should succeed even without existing identity: %v", err)
	}

	identity, found := k.GetIdentity(ctx, user1)
	if !found {
		t.Fatal("identity should be auto-created")
	}
	if identity.TotalBookings != 1 {
		t.Fatalf("expected total_bookings = 1, got %d", identity.TotalBookings)
	}
	if identity.TrustScore != 10 {
		t.Fatalf("expected trust_score = 10, got %d", identity.TrustScore)
	}

	t.Log("PASS: Bookings incremented for nonexistent identity (auto-created)")
}

func TestIncrementBookingsUnauthorized(t *testing.T) {
	k, ctx := setupKeeper(t)

	err := k.IncrementBookings(ctx, &types.MsgIncrementBookings{
		Authority: user1,
		Address:   user2,
	})
	if err == nil {
		t.Fatal("unauthorized increment should fail")
	}

	t.Log("PASS: Unauthorized increment rejected")
}

// --- IsVerified Tests ---

func TestIsVerified(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Not found = not verified
	if k.IsVerified(ctx, user1) {
		t.Fatal("nonexistent user should not be verified")
	}

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	})

	// Pending = not verified
	if k.IsVerified(ctx, user1) {
		t.Fatal("pending user should not be verified")
	}

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "basic",
	})

	if !k.IsVerified(ctx, user1) {
		t.Fatal("verified user should return true")
	}

	t.Log("PASS: IsVerified works correctly")
}

// --- Expiry Tests ---

func TestExpireIdentities(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Set short expiry
	params := k.GetParams(ctx)
	params.VerificationExpiryBlocks = 10
	k.SetParams(ctx, params)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	})

	k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "basic",
	})

	// At block 1, should be verified
	if !k.IsVerified(ctx, user1) {
		t.Fatal("should be verified at block 1")
	}

	// Advance to block 12 (past expiry of block 1 + 10 = 11)
	ctx = ctx.WithBlockHeight(12)

	// Run expiry
	k.ExpireIdentities(ctx)

	identity, _ := k.GetIdentity(ctx, user1)
	if identity.Status != types.StatusExpired {
		t.Fatalf("expected status expired, got %s", identity.Status)
	}

	if k.IsVerified(ctx, user1) {
		t.Fatal("expired identity should not be verified")
	}

	t.Log("PASS: Identity expired correctly")
}

func TestExpireIdentitiesNotYetExpired(t *testing.T) {
	k, ctx := setupKeeper(t)

	params := k.GetParams(ctx)
	params.VerificationExpiryBlocks = 100
	k.SetParams(ctx, params)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	})

	k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "basic",
	})

	// At block 50 (before expiry at block 101)
	ctx = ctx.WithBlockHeight(50)
	k.ExpireIdentities(ctx)

	identity, _ := k.GetIdentity(ctx, user1)
	if identity.Status != types.StatusVerified {
		t.Fatalf("expected status still verified, got %s", identity.Status)
	}

	t.Log("PASS: Identity not expired before expiry block")
}

// --- Genesis Tests ---

func TestExportImportGenesis(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address:      user1,
		DocumentHash: "doc1",
		Level:        "basic",
	})

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address:      user2,
		DocumentHash: "doc2",
		Level:        "standard",
	})

	// Export genesis
	gs := k.ExportGenesis(ctx)
	if len(gs.Identities) != 2 {
		t.Fatalf("expected 2 identities in genesis, got %d", len(gs.Identities))
	}
	if len(gs.Verifiers) != 1 {
		t.Fatalf("expected 1 verifier in genesis, got %d", len(gs.Verifiers))
	}

	// Create new keeper and import
	k2, ctx2 := setupKeeper(t)
	k2.InitGenesis(ctx2, *gs)

	id1, found := k2.GetIdentity(ctx2, user1)
	if !found {
		t.Fatal("identity user1 should exist after import")
	}
	if id1.DocumentHash != "doc1" {
		t.Fatalf("expected doc hash 'doc1', got '%s'", id1.DocumentHash)
	}

	v1, found := k2.GetVerifier(ctx2, verifier1Addr)
	if !found {
		t.Fatal("verifier should exist after import")
	}
	if v1.Name != "KYC" {
		t.Fatalf("expected verifier name 'KYC', got '%s'", v1.Name)
	}

	t.Log("PASS: Genesis export/import works correctly")
}

// --- Query Tests ---

func TestQueryAllIdentitiesByVerifier(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	})
	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user2,
		Level:   "basic",
	})

	k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "basic",
	})

	// Only user1 was verified by verifier1Addr
	identities := k.GetIdentitiesByVerifier(ctx, verifier1Addr)
	if len(identities) != 1 {
		t.Fatalf("expected 1 identity by verifier, got %d", len(identities))
	}
	if identities[0].Address != user1 {
		t.Fatalf("expected user1, got %s", identities[0].Address)
	}

	t.Log("PASS: Identities by verifier query works")
}

func TestGetAllVerifiers(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC 1",
		MaxLevel:  "enhanced",
	})

	// Use a different address for second verifier
	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  user3,
		Name:      "KYC 2",
		MaxLevel:  "basic",
	})

	verifiers := k.GetAllVerifiers(ctx)
	if len(verifiers) != 2 {
		t.Fatalf("expected 2 verifiers, got %d", len(verifiers))
	}

	t.Log("PASS: GetAllVerifiers works correctly")
}

func TestInvalidVerificationLevel(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.RegisterVerifier(ctx, &types.MsgRegisterVerifier{
		Authority: authority,
		Verifier:  verifier1Addr,
		Name:      "KYC",
		MaxLevel:  "enhanced",
	})

	k.CreateIdentity(ctx, &types.MsgRegisterIdentity{
		Address: user1,
		Level:   "basic",
	})

	// Try to verify with invalid level
	err := k.VerifyIdentity(ctx, &types.MsgVerifyIdentity{
		Verifier: verifier1Addr,
		Address:  user1,
		Level:    "super_ultra_verified",
	})
	if err == nil {
		t.Fatal("invalid verification level should fail")
	}

	t.Log("PASS: Invalid verification level rejected")
}
