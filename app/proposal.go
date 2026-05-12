package app

import (
	"bytes"
	"crypto/sha256"
	"sort"

	abci "github.com/cometbft/cometbft/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	mevkeeper "syreen/x/mevprotection/keeper"
)

// revealedTxInfo holds ordering metadata for a revealed transaction in a proposed block.
type revealedTxInfo struct {
	txIndex           int    // position of the tx in the block
	commitBlockHeight int64  // height at which the commitment was originally made
	commitTimestamp   int64  // timestamp of the original commitment
	commitHash        []byte // hash of the original commitment (for deterministic tiebreaking)
}

// PrepareProposalHandler returns an ABCI PrepareProposal handler that enforces
// MEV protection by fair-ordering revealed transactions.
//
// Ordering rules:
//  1. Revealed transactions (those that went through commit-reveal) come first,
//     sorted by their original commit block height (first committed = first included).
//     Ties are broken by commit timestamp, then by commit hash for determinism.
//  2. Normal transactions (not part of commit-reveal) come after all revealed txs,
//     preserving their original mempool order.
//  3. The total block size respects the max bytes limit.
func (app *SyreenApp) PrepareProposalHandler() sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
		params := app.MEVProtectionKeeper.GetParams(ctx)
		if !params.FairOrderConfig.EnableCommitReveal {
			// Commit-reveal disabled: pass through unchanged
			return &abci.ResponsePrepareProposal{Txs: req.Txs}, nil
		}

		// Build a lookup of revealed tx body hashes -> commit metadata
		revealedIndex := buildRevealedIndex(ctx, app.MEVProtectionKeeper)

		var revealed []indexedTx
		var normal []indexedTx

		for i, txBytes := range req.Txs {
			bodyHash := sha256.Sum256(txBytes)
			if info, ok := revealedIndex[string(bodyHash[:])]; ok {
				revealed = append(revealed, indexedTx{
					txBytes:           txBytes,
					origIndex:         i,
					commitBlockHeight: info.commitBlockHeight,
					commitTimestamp:   info.commitTimestamp,
					commitHash:        info.commitHash,
				})
			} else {
				normal = append(normal, indexedTx{
					txBytes:   txBytes,
					origIndex: i,
				})
			}
		}

		// Sort revealed txs by commit block height (fair ordering).
		sort.SliceStable(revealed, func(i, j int) bool {
			if revealed[i].commitBlockHeight != revealed[j].commitBlockHeight {
				return revealed[i].commitBlockHeight < revealed[j].commitBlockHeight
			}
			if revealed[i].commitTimestamp != revealed[j].commitTimestamp {
				return revealed[i].commitTimestamp < revealed[j].commitTimestamp
			}
			return bytes.Compare(revealed[i].commitHash, revealed[j].commitHash) < 0
		})

		// Assemble final tx list respecting max block size
		maxBytes := req.MaxTxBytes
		var finalTxs [][]byte
		var totalSize int64

		for _, tx := range revealed {
			txSize := int64(len(tx.txBytes))
			if totalSize+txSize > maxBytes {
				break
			}
			finalTxs = append(finalTxs, tx.txBytes)
			totalSize += txSize
		}

		for _, tx := range normal {
			txSize := int64(len(tx.txBytes))
			if totalSize+txSize > maxBytes {
				break
			}
			finalTxs = append(finalTxs, tx.txBytes)
			totalSize += txSize
		}

		return &abci.ResponsePrepareProposal{Txs: finalTxs}, nil
	}
}

// ProcessProposalHandler returns an ABCI ProcessProposal handler that validates
// proposed blocks for MEV protection compliance.
//
// It verifies that revealed transactions appear in commit-height order within
// the block. A tolerance of MaxTxDelay positions of drift is allowed.
func (app *SyreenApp) ProcessProposalHandler() sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		params := app.MEVProtectionKeeper.GetParams(ctx)
		if !params.FairOrderConfig.EnableCommitReveal {
			return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
		}

		revealedIndex := buildRevealedIndex(ctx, app.MEVProtectionKeeper)

		// Collect revealed txs in the order they appear in the proposed block
		var blockRevealed []revealedTxInfo
		for i, txBytes := range req.Txs {
			bodyHash := sha256.Sum256(txBytes)
			if info, ok := revealedIndex[string(bodyHash[:])]; ok {
				blockRevealed = append(blockRevealed, revealedTxInfo{
					txIndex:           i,
					commitBlockHeight: info.commitBlockHeight,
					commitTimestamp:   info.commitTimestamp,
					commitHash:        info.commitHash,
				})
			}
		}

		// Fewer than 2 revealed txs — nothing to check
		if len(blockRevealed) < 2 {
			return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
		}

		// Build the expected (sorted) ordering
		sorted := make([]revealedTxInfo, len(blockRevealed))
		copy(sorted, blockRevealed)
		sort.SliceStable(sorted, func(i, j int) bool {
			if sorted[i].commitBlockHeight != sorted[j].commitBlockHeight {
				return sorted[i].commitBlockHeight < sorted[j].commitBlockHeight
			}
			if sorted[i].commitTimestamp != sorted[j].commitTimestamp {
				return sorted[i].commitTimestamp < sorted[j].commitTimestamp
			}
			return bytes.Compare(sorted[i].commitHash, sorted[j].commitHash) < 0
		})

		maxDelay := int(params.FairOrderConfig.MaxTxDelay)

		for i := range blockRevealed {
			if !bytes.Equal(blockRevealed[i].commitHash, sorted[i].commitHash) {
				// Compute position drift between actual and expected
				drift := blockRevealed[i].txIndex - sorted[i].txIndex
				if drift < 0 {
					drift = -drift
				}
				if drift > maxDelay {
					return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
				}
			}
		}

		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}
}

// indexedTx is a transaction annotated with ordering metadata for PrepareProposal.
type indexedTx struct {
	txBytes           []byte
	origIndex         int
	commitBlockHeight int64
	commitTimestamp   int64
	commitHash        []byte
}

// commitInfo holds commit metadata for a revealed transaction body hash.
type commitInfo struct {
	commitBlockHeight int64
	commitTimestamp   int64
	commitHash        []byte
}

// buildRevealedIndex creates a map from SHA256(revealed tx body) -> commit metadata.
// This allows O(1) lookup when classifying block transactions.
func buildRevealedIndex(ctx sdk.Context, keeper mevkeeper.Keeper) map[string]commitInfo {
	allRevealed := keeper.GetAllRevealedTxs(ctx)
	index := make(map[string]commitInfo, len(allRevealed))
	for _, r := range allRevealed {
		bodyHash := sha256.Sum256(r.Revealed.ActualTx)
		index[string(bodyHash[:])] = commitInfo{
			commitBlockHeight: r.CommitBlockHeight,
			commitTimestamp:   r.CommitTimestamp,
			commitHash:        r.Revealed.CommittedHash,
		}
	}
	return index
}
