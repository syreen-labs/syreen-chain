package keeper

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/mevprotection/types"
)

type Keeper struct {
	cdc            codec.Codec
	storeService   store.KVStoreService
	stakingKeeper  types.StakingKeeper
	slashingKeeper types.SlashingKeeper
	bankKeeper     types.BankKeeper
	authority      string
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	stakingKeeper types.StakingKeeper,
	slashingKeeper types.SlashingKeeper,
	bankKeeper types.BankKeeper,
	authority string,
) Keeper {
	return Keeper{
		cdc:            cdc,
		storeService:   storeService,
		stakingKeeper:  stakingKeeper,
		slashingKeeper: slashingKeeper,
		bankKeeper:     bankKeeper,
		authority:      authority,
	}
}

// GetAuthority returns the module's governance authority address
func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// CommitTx stores a new transaction commitment in the commit-reveal scheme.
// txHash must be SHA256(txBody || nonce) — a 32-byte digest.
// encryptedTx is stored alongside the commitment for later reference.
func (k Keeper) CommitTx(ctx context.Context, sender string, txHash []byte, encryptedTx []byte) error {
	// Validate hash length (defense in depth — ValidateBasic also checks)
	if len(txHash) != 32 {
		return fmt.Errorf("tx_hash must be 32 bytes, got %d", len(txHash))
	}

	// Check for duplicate commitment
	_, found := k.GetCommittedTx(ctx, txHash)
	if found {
		return types.ErrDuplicateCommit
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	committed := types.CommittedTx{
		Sender:      sender,
		TxHash:      txHash,
		EncryptedTx: encryptedTx,
		BlockHeight: sdkCtx.BlockHeight(),
		Timestamp:   sdkCtx.BlockTime().Unix(),
	}

	k.SetCommittedTx(ctx, txHash, committed)

	k.Logger(ctx).Info("committed tx", "sender", sender, "hash", fmt.Sprintf("%x", txHash), "height", sdkCtx.BlockHeight())
	return nil
}

// RevealTx validates and stores a reveal against a prior commitment
func (k Keeper) RevealTx(ctx context.Context, sender string, commitHash []byte, txBody []byte, nonce []byte) error {
	committed, found := k.GetCommittedTx(ctx, commitHash)
	if !found {
		return types.ErrRevealMismatch
	}

	// Verify the sender matches the original committer
	if committed.Sender != sender {
		return types.ErrRevealMismatch
	}

	// Verify the reveal matches the commitment by hashing txBody + nonce
	h := sha256.New()
	h.Write(txBody)
	h.Write(nonce)
	computedHash := h.Sum(nil)

	if !bytes.Equal(computedHash, commitHash) {
		return types.ErrRevealMismatch
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(ctx)

	// Minimum reveal delay: must wait at least half the commit window before revealing
	minRevealHeight := committed.BlockHeight + int64(params.FairOrderConfig.CommitWindow/2)
	if sdkCtx.BlockHeight() < minRevealHeight {
		return types.ErrRevealTooEarly
	}

	// Check reveal window hasn't expired
	maxRevealHeight := committed.BlockHeight + int64(params.FairOrderConfig.CommitWindow) + int64(params.FairOrderConfig.RevealWindow)
	if sdkCtx.BlockHeight() > maxRevealHeight {
		return types.ErrRevealExpired
	}

	revealed := types.RevealedTx{
		CommittedHash: commitHash,
		ActualTx:      txBody,
		Proof:         nonce,
	}

	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(revealed)
	if err != nil {
		return err
	}
	if err := kvStore.Set(types.RevealedTxKey(commitHash), bz); err != nil {
		return err
	}

	k.Logger(ctx).Info("revealed tx", "sender", sender, "commit_hash", fmt.Sprintf("%x", commitHash), "height", sdkCtx.BlockHeight())
	return nil
}

// orderedTx is used internally for fair ordering checks
type orderedTx struct {
	Hash      []byte `json:"hash"`
	Timestamp int64  `json:"timestamp"`
	Index     int    `json:"index"`
}

// CheckFairOrdering validates that revealed transactions within a block are ordered
// fairly by their original commit timestamp.
// It cross-references block transactions against revealed tx records.
func (k Keeper) CheckFairOrdering(ctx context.Context, blockTxs [][]byte) error {
	if len(blockTxs) <= 1 {
		return nil
	}

	params := k.GetParams(ctx)
	if !params.FairOrderConfig.EnableCommitReveal {
		return nil
	}

	// Build an index of raw tx hash -> block position for fast lookup
	txIndex := make(map[string]int, len(blockTxs))
	for i, tx := range blockTxs {
		h := sha256.Sum256(tx)
		txIndex[string(h[:])] = i
	}

	// Iterate all revealed txs and match them against block transactions
	var ordered []orderedTx
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.RevealedTxPrefixKey)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var revealed types.RevealedTx
		if err := json.Unmarshal(iter.Value(), &revealed); err != nil {
			continue
		}

		// Hash the revealed tx body to see if it matches any block tx
		bodyHash := sha256.Sum256(revealed.ActualTx)
		if idx, found := txIndex[string(bodyHash[:])]; found {
			// Look up the original commit to get the timestamp
			committed, commitFound := k.GetCommittedTx(ctx, revealed.CommittedHash)
			if commitFound {
				ordered = append(ordered, orderedTx{
					Hash:      revealed.CommittedHash,
					Timestamp: committed.Timestamp,
					Index:     idx,
				})
			}
		}
	}

	if len(ordered) <= 1 {
		return nil
	}

	// Verify ordering: committed txs should be ordered by timestamp
	// Use stable sort so ties (same-block commits) maintain original order
	sorted := make([]orderedTx, len(ordered))
	copy(sorted, ordered)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Timestamp != sorted[j].Timestamp {
			return sorted[i].Timestamp < sorted[j].Timestamp
		}
		// Deterministic tiebreaker: by commit hash
		return bytes.Compare(sorted[i].Hash, sorted[j].Hash) < 0
	})

	maxDelay := int(params.FairOrderConfig.MaxTxDelay)
	for i := range ordered {
		if !bytes.Equal(ordered[i].Hash, sorted[i].Hash) {
			// Allow some tolerance based on MaxTxDelay
			positionDrift := ordered[i].Index - sorted[i].Index
			if positionDrift < 0 {
				positionDrift = -positionDrift
			}
			if positionDrift > maxDelay {
				return types.ErrMEVDetected
			}
		}
	}

	return nil
}

// DetectMEV checks for suspicious reordering patterns in a block.
// When a violation is found it records a penalty AND slashes the proposer validator.
func (k Keeper) DetectMEV(ctx context.Context, blockTxs [][]byte) (*types.MEVPenalty, error) {
	err := k.CheckFairOrdering(ctx, blockTxs)
	if err == nil {
		return nil, nil
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	proposerAddr := sdkCtx.BlockHeader().ProposerAddress
	consAddr := sdk.ConsAddress(proposerAddr)

	// Build evidence from the misordered transactions
	evidence, _ := json.Marshal(map[string]interface{}{
		"block_height": sdkCtx.BlockHeight(),
		"tx_count":     len(blockTxs),
		"violation":    "transaction ordering does not match arrival timestamps",
	})

	penalty := &types.MEVPenalty{
		ValidatorAddr: consAddr.String(),
		PenaltyType:   "reordering",
		Amount:        1,
		Evidence:      evidence,
		BlockHeight:   sdkCtx.BlockHeight(),
	}

	// Store the penalty record
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, marshalErr := json.Marshal(penalty)
	if marshalErr != nil {
		return nil, marshalErr
	}
	if setErr := kvStore.Set(types.PenaltyKey(penalty.ValidatorAddr, penalty.BlockHeight), bz); setErr != nil {
		return nil, setErr
	}

	// Determine slash fraction from params (default 5%)
	params := k.GetParams(ctx)
	slashFractionStr := params.SlashFraction
	if slashFractionStr == "" {
		slashFractionStr = types.DefaultSlashFraction
	}
	slashFraction, parseErr := math.LegacyNewDecFromStr(slashFractionStr)
	if parseErr != nil {
		k.Logger(ctx).Error("invalid slash fraction in params, using default", "err", parseErr)
		slashFraction = math.LegacyNewDecWithPrec(5, 2) // 0.05
	}

	// Look up validator by consensus address to get actual voting power.
	// If lookup fails, skip slashing entirely to avoid over-slashing.
	validator, valErr := k.stakingKeeper.GetValidatorByConsAddr(ctx, consAddr)
	if valErr != nil {
		k.Logger(ctx).Error("cannot slash: validator not found", "consAddr", consAddr.String(), "error", valErr)
		return penalty, nil
	}
	validatorPower := validator.GetConsensusPower(sdk.DefaultPowerReduction)
	if validatorPower <= 0 {
		k.Logger(ctx).Error("cannot slash: validator has zero power", "consAddr", consAddr.String())
		return penalty, nil
	}
	k.Logger(ctx).Info("MEV slash using actual validator power",
		"consAddr", consAddr.String(),
		"power", validatorPower,
	)

	slashedAmt, slashErr := k.stakingKeeper.Slash(ctx, consAddr, sdkCtx.BlockHeight(), validatorPower, slashFraction)
	if slashErr != nil {
		k.Logger(ctx).Error("failed to slash validator for MEV", "validator", consAddr.String(), "err", slashErr)
	} else {
		k.Logger(ctx).Info("slashed validator for MEV", "validator", consAddr.String(), "fraction", slashFraction.String(), "height", sdkCtx.BlockHeight())

		// Accumulate slashed amount for MEV redistribution to LPs and stakers
		if slashedAmt.IsPositive() {
			bondDenom := "usyreen" // default bond denom
			rewardCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, slashedAmt))
			k.AccumulateMEVReward(ctx, rewardCoins)
		}
	}

	// Jail the validator via the slashing keeper to prevent further proposals
	if jailErr := k.slashingKeeper.Jail(ctx, consAddr); jailErr != nil {
		k.Logger(ctx).Error("failed to jail validator for MEV", "validator", consAddr.String(), "err", jailErr)
	}

	k.Logger(ctx).Info("MEV detected", "validator", consAddr.String(), "height", penalty.BlockHeight)
	return penalty, nil
}

// GetRevealedTx retrieves a revealed transaction by its commit hash
func (k Keeper) GetRevealedTx(ctx context.Context, commitHash []byte) (types.RevealedTx, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.RevealedTxKey(commitHash))
	if err != nil || bz == nil {
		return types.RevealedTx{}, false
	}
	var revealed types.RevealedTx
	if err := json.Unmarshal(bz, &revealed); err != nil {
		return types.RevealedTx{}, false
	}
	return revealed, true
}

// GetAllRevealedTxs returns all currently stored revealed transactions along with their original commit data.
// Each entry contains the revealed tx and the commit block height for ordering purposes.
func (k Keeper) GetAllRevealedTxs(ctx context.Context) []RevealedWithCommit {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.RevealedTxPrefixKey)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var results []RevealedWithCommit
	for ; iter.Valid(); iter.Next() {
		var revealed types.RevealedTx
		if err := json.Unmarshal(iter.Value(), &revealed); err != nil {
			continue
		}
		committed, found := k.GetCommittedTx(ctx, revealed.CommittedHash)
		if !found {
			continue
		}
		results = append(results, RevealedWithCommit{
			Revealed:         revealed,
			CommitBlockHeight: committed.BlockHeight,
			CommitTimestamp:   committed.Timestamp,
		})
	}

	// Explicit deterministic sort by commit hash, regardless of DB backend
	// iteration order. This prevents consensus forks if different backends
	// return keys in different orders.
	sort.Slice(results, func(i, j int) bool {
		return bytes.Compare(results[i].Revealed.CommittedHash, results[j].Revealed.CommittedHash) < 0
	})

	return results
}

// RevealedWithCommit pairs a revealed transaction with its original commit metadata.
type RevealedWithCommit struct {
	Revealed          types.RevealedTx
	CommitBlockHeight int64
	CommitTimestamp    int64
}

// prefixEndBytes returns the end key for a prefix scan (increment last byte).
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
	return nil // overflow: prefix is all 0xFF
}

// GetCommittedTx retrieves a committed transaction by its hash
func (k Keeper) GetCommittedTx(ctx context.Context, txHash []byte) (types.CommittedTx, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.CommittedTxKey(txHash))
	if err != nil || bz == nil {
		return types.CommittedTx{}, false
	}
	var committed types.CommittedTx
	if err := json.Unmarshal(bz, &committed); err != nil {
		return types.CommittedTx{}, false
	}
	return committed, true
}

// SetCommittedTx stores a committed transaction
func (k Keeper) SetCommittedTx(ctx context.Context, txHash []byte, committed types.CommittedTx) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(committed)
	kvStore.Set(types.CommittedTxKey(txHash), bz)
}

// GetPendingReveals returns committed txs that have not yet been revealed and are still within the reveal window
func (k Keeper) GetPendingReveals(ctx context.Context) []types.CommittedTx {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(ctx)
	kvStore := k.storeService.OpenKVStore(ctx)

	var pending []types.CommittedTx

	// Iterate over committed tx prefix with bounded end key
	prefix := []byte(types.CommittedTxPrefixKey)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {

		var committed types.CommittedTx
		if err := json.Unmarshal(iter.Value(), &committed); err != nil {
			continue
		}

		// Check if still within reveal window
		maxRevealHeight := committed.BlockHeight + int64(params.FairOrderConfig.CommitWindow) + int64(params.FairOrderConfig.RevealWindow)
		if sdkCtx.BlockHeight() <= maxRevealHeight {
			// Check if not yet revealed
			revealKey := types.RevealedTxKey(committed.TxHash)
			revealBz, err := kvStore.Get(revealKey)
			if err != nil || revealBz == nil {
				pending = append(pending, committed)
			}
		}
	}

	// Explicit deterministic sort by tx hash to guarantee identical results
	// regardless of KV store backend iteration order.
	sort.Slice(pending, func(i, j int) bool {
		return bytes.Compare(pending[i].TxHash, pending[j].TxHash) < 0
	})

	return pending
}

// BeginBlockMEVCheck is called during BeginBlock as a secondary safety net.
// It records an event if any MEV penalty was already assessed for this height,
// serving as an audit trail. The primary enforcement is in ProcessProposal.
func (k Keeper) BeginBlockMEVCheck(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(ctx)
	if !params.FairOrderConfig.EnableCommitReveal {
		return
	}

	// Check if the previous block's proposer has any penalty recorded.
	// This acts as an audit log entry in the block events.
	prevHeight := sdkCtx.BlockHeight() - 1
	if prevHeight <= 0 {
		return
	}

	proposerAddr := sdkCtx.BlockHeader().ProposerAddress
	if len(proposerAddr) == 0 {
		return
	}

	consAddr := sdk.ConsAddress(proposerAddr)
	kvStore := k.storeService.OpenKVStore(ctx)
	penaltyKey := types.PenaltyKey(consAddr.String(), prevHeight)
	bz, err := kvStore.Get(penaltyKey)
	if err != nil || bz == nil {
		return
	}

	k.Logger(ctx).Info("MEV penalty record found for previous block proposer",
		"proposer", consAddr.String(),
		"penalty_height", prevHeight,
	)

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"mev_penalty_detected",
		sdk.NewAttribute("proposer", consAddr.String()),
		sdk.NewAttribute("infraction_height", fmt.Sprintf("%d", prevHeight)),
	))
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

// SetParams stores the module parameters. Authority must be validated by the caller.
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return kvStore.Set([]byte("params"), bz)
}

// maxPrunePerBlock bounds how many expired commits are pruned in a single block.
// PruneExpiredCommits does a full scan of the committed-tx keyspace every block
// with no expiry index; that keyspace is attacker-growable, so an unbounded scan
// lets per-block cost grow forever until the chain misses timeout_commit and
// stalls. Capping the number of prunes per block bounds worst-case work; commits
// beyond the cap are pruned over subsequent blocks (pruning a few blocks late is
// acceptable — once deleted they leave the candidate set). Deterministic: store
// iteration order is identical on all validators. Long-term fix: add a
// height-bucketed expiry index so only commits due at/before the current height
// are visited.
const maxPrunePerBlock = 500

// PruneExpiredCommits removes committed transactions whose reveal window has expired.
// Should be called in BeginBlocker to prevent unbounded state growth.
func (k Keeper) PruneExpiredCommits(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(ctx)
	kvStore := k.storeService.OpenKVStore(ctx)

	prefix := []byte(types.CommittedTxPrefixKey)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return
	}
	defer iter.Close()

	// prunedCommits counts committed-tx entries selected for deletion (each may
	// also enqueue a corresponding reveal key); it is bounded by maxPrunePerBlock.
	var keysToDelete [][]byte
	prunedCommits := 0
	for ; iter.Valid(); iter.Next() {
		var committed types.CommittedTx
		if err := json.Unmarshal(iter.Value(), &committed); err != nil {
			// Corrupted entry — mark for deletion
			keysToDelete = append(keysToDelete, append([]byte(nil), iter.Key()...))
			prunedCommits++
			if prunedCommits >= maxPrunePerBlock {
				break
			}
			continue
		}

		maxRevealHeight := committed.BlockHeight + int64(params.FairOrderConfig.CommitWindow) + int64(params.FairOrderConfig.RevealWindow)
		if sdkCtx.BlockHeight() > maxRevealHeight {
			keysToDelete = append(keysToDelete, append([]byte(nil), iter.Key()...))
			// Also delete corresponding reveal if it exists
			revealKey := types.RevealedTxKey(committed.TxHash)
			keysToDelete = append(keysToDelete, revealKey)
			// Bound per-block work: stop collecting once the cap is reached.
			// Remaining expired commits are pruned in subsequent blocks.
			prunedCommits++
			if prunedCommits >= maxPrunePerBlock {
				break
			}
		}
	}

	for _, key := range keysToDelete {
		kvStore.Delete(key)
	}

	if len(keysToDelete) > 0 {
		k.Logger(ctx).Debug("pruned expired commits", "count", len(keysToDelete)/2)
	}
}
