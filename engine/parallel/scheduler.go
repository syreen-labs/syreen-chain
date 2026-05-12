package parallel

import (
	"crypto/sha256"
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	computetypes "syreen/x/compute/types"
	intenttypes "syreen/x/intent/types"
	tokenfactorytypes "syreen/x/tokenfactory/types"
)

// Scheduler analyzes transactions and groups them for parallel execution.
type Scheduler struct {
	maxParallel int
}

// NewScheduler creates a new dependency-aware scheduler.
func NewScheduler(maxParallel int) *Scheduler {
	if maxParallel <= 0 {
		maxParallel = 4
	}
	return &Scheduler{maxParallel: maxParallel}
}

// PredictAccessKeys returns predicted store keys that a transaction will access.
// This is a heuristic -- actual access is tracked during execution for validation.
func PredictAccessKeys(tx sdk.Tx) (reads []string, writes []string) {
	msgs := tx.GetMsgs()

	// Extract signer addresses from the tx signatures
	if sigTx, ok := tx.(interface{ GetSigners() [][]byte }); ok {
		for _, signer := range sigTx.GetSigners() {
			addrKey := fmt.Sprintf("acc:%x", signer)
			reads = append(reads, addrKey)
			writes = append(writes, addrKey) // sequence number update
		}
	}

	for _, msg := range msgs {
		r, w := predictMsgAccessKeys(msg)
		reads = append(reads, r...)
		writes = append(writes, w...)
	}

	// Deduplicate and sort for deterministic output across all nodes.
	reads = uniqueSorted(reads)
	writes = uniqueSorted(writes)
	return
}

// uniqueSorted deduplicates and lexicographically sorts a string slice.
func uniqueSorted(ss []string) []string {
	if len(ss) <= 1 {
		return ss
	}
	seen := make(map[string]struct{}, len(ss))
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}

// predictMsgAccessKeys returns predicted read/write keys for a specific message type.
func predictMsgAccessKeys(msg sdk.Msg) (reads []string, writes []string) {
	switch m := msg.(type) {

	// ---- Bank ----
	case *banktypes.MsgSend:
		reads = append(reads, fmt.Sprintf("bank:%s", m.FromAddress))
		reads = append(reads, fmt.Sprintf("bank:%s", m.ToAddress))
		writes = append(writes, fmt.Sprintf("bank:%s", m.FromAddress))
		writes = append(writes, fmt.Sprintf("bank:%s", m.ToAddress))

	case *banktypes.MsgMultiSend:
		for _, input := range m.Inputs {
			reads = append(reads, fmt.Sprintf("bank:%s", input.Address))
			writes = append(writes, fmt.Sprintf("bank:%s", input.Address))
		}
		for _, output := range m.Outputs {
			reads = append(reads, fmt.Sprintf("bank:%s", output.Address))
			writes = append(writes, fmt.Sprintf("bank:%s", output.Address))
		}

	// ---- Staking ----
	case *stakingtypes.MsgDelegate:
		reads = append(reads, fmt.Sprintf("staking:val:%s", m.ValidatorAddress))
		reads = append(reads, fmt.Sprintf("staking:del:%s", m.DelegatorAddress))
		writes = append(writes, fmt.Sprintf("staking:val:%s", m.ValidatorAddress))
		writes = append(writes, fmt.Sprintf("staking:del:%s", m.DelegatorAddress))
		reads = append(reads, fmt.Sprintf("bank:%s", m.DelegatorAddress))
		writes = append(writes, fmt.Sprintf("bank:%s", m.DelegatorAddress))

	case *stakingtypes.MsgUndelegate:
		reads = append(reads, fmt.Sprintf("staking:val:%s", m.ValidatorAddress))
		reads = append(reads, fmt.Sprintf("staking:del:%s", m.DelegatorAddress))
		writes = append(writes, fmt.Sprintf("staking:val:%s", m.ValidatorAddress))
		writes = append(writes, fmt.Sprintf("staking:del:%s", m.DelegatorAddress))
		reads = append(reads, fmt.Sprintf("bank:%s", m.DelegatorAddress))
		writes = append(writes, fmt.Sprintf("bank:%s", m.DelegatorAddress))

	case *stakingtypes.MsgBeginRedelegate:
		reads = append(reads, fmt.Sprintf("staking:val:%s", m.ValidatorSrcAddress))
		reads = append(reads, fmt.Sprintf("staking:val:%s", m.ValidatorDstAddress))
		reads = append(reads, fmt.Sprintf("staking:del:%s", m.DelegatorAddress))
		writes = append(writes, fmt.Sprintf("staking:val:%s", m.ValidatorSrcAddress))
		writes = append(writes, fmt.Sprintf("staking:val:%s", m.ValidatorDstAddress))
		writes = append(writes, fmt.Sprintf("staking:del:%s", m.DelegatorAddress))

	// ---- Distribution ----
	case *distrtypes.MsgWithdrawDelegatorReward:
		reads = append(reads, fmt.Sprintf("distr:val:%s", m.ValidatorAddress))
		reads = append(reads, fmt.Sprintf("distr:del:%s", m.DelegatorAddress))
		writes = append(writes, fmt.Sprintf("distr:val:%s", m.ValidatorAddress))
		writes = append(writes, fmt.Sprintf("distr:del:%s", m.DelegatorAddress))
		writes = append(writes, fmt.Sprintf("bank:%s", m.DelegatorAddress))

	// ---- Governance ----
	case *govv1beta1.MsgVote:
		reads = append(reads, fmt.Sprintf("gov:proposal:%d", m.ProposalId))
		writes = append(writes, fmt.Sprintf("gov:vote:%d:%s", m.ProposalId, m.Voter))

	case *govv1.MsgVote:
		reads = append(reads, fmt.Sprintf("gov:proposal:%d", m.ProposalId))
		writes = append(writes, fmt.Sprintf("gov:vote:%d:%s", m.ProposalId, m.Voter))

	case *govv1beta1.MsgDeposit:
		reads = append(reads, fmt.Sprintf("gov:proposal:%d", m.ProposalId))
		writes = append(writes, fmt.Sprintf("gov:proposal:%d", m.ProposalId))
		reads = append(reads, fmt.Sprintf("bank:%s", m.Depositor))
		writes = append(writes, fmt.Sprintf("bank:%s", m.Depositor))

	case *govv1.MsgDeposit:
		reads = append(reads, fmt.Sprintf("gov:proposal:%d", m.ProposalId))
		writes = append(writes, fmt.Sprintf("gov:proposal:%d", m.ProposalId))
		reads = append(reads, fmt.Sprintf("bank:%s", m.Depositor))
		writes = append(writes, fmt.Sprintf("bank:%s", m.Depositor))

	case *govv1beta1.MsgSubmitProposal:
		// New proposal writes a new proposal ID; conflicts with any other proposal submission
		writes = append(writes, "gov:next_proposal_id")
		if signers, err := extractMsgSigners(m); err == nil {
			for _, signer := range signers {
				reads = append(reads, fmt.Sprintf("bank:%s", signer))
				writes = append(writes, fmt.Sprintf("bank:%s", signer))
			}
		}

	case *govv1.MsgSubmitProposal:
		writes = append(writes, "gov:next_proposal_id")
		reads = append(reads, fmt.Sprintf("bank:%s", m.Proposer))
		writes = append(writes, fmt.Sprintf("bank:%s", m.Proposer))

	// ---- Token Factory ----
	case *tokenfactorytypes.MsgCreateDenom:
		writes = append(writes, fmt.Sprintf("tf:denom:%s/%s", m.Sender, m.Subdenom))
		reads = append(reads, fmt.Sprintf("bank:%s", m.Sender))
		writes = append(writes, fmt.Sprintf("bank:%s", m.Sender))

	case *tokenfactorytypes.MsgMint:
		denom := m.Amount.Denom
		reads = append(reads, fmt.Sprintf("tf:admin:%s", denom))
		writes = append(writes, fmt.Sprintf("tf:supply:%s", denom))
		if m.MintTo != "" {
			writes = append(writes, fmt.Sprintf("bank:%s", m.MintTo))
		} else {
			writes = append(writes, fmt.Sprintf("bank:%s", m.Sender))
		}

	case *tokenfactorytypes.MsgBurn:
		denom := m.Amount.Denom
		reads = append(reads, fmt.Sprintf("tf:admin:%s", denom))
		writes = append(writes, fmt.Sprintf("tf:supply:%s", denom))
		if m.BurnFrom != "" {
			reads = append(reads, fmt.Sprintf("bank:%s", m.BurnFrom))
			writes = append(writes, fmt.Sprintf("bank:%s", m.BurnFrom))
		} else {
			reads = append(reads, fmt.Sprintf("bank:%s", m.Sender))
			writes = append(writes, fmt.Sprintf("bank:%s", m.Sender))
		}

	case *tokenfactorytypes.MsgChangeAdmin:
		reads = append(reads, fmt.Sprintf("tf:admin:%s", m.Denom))
		writes = append(writes, fmt.Sprintf("tf:admin:%s", m.Denom))

	// ---- Intent ----
	case *intenttypes.MsgSubmitIntent:
		writes = append(writes, fmt.Sprintf("intent:next_id"))
		reads = append(reads, fmt.Sprintf("bank:%s", m.Creator))
		writes = append(writes, fmt.Sprintf("bank:%s", m.Creator))

	case *intenttypes.MsgSubmitSolution:
		reads = append(reads, fmt.Sprintf("intent:%s", m.IntentID))
		writes = append(writes, fmt.Sprintf("intent:%s", m.IntentID))
		reads = append(reads, fmt.Sprintf("intent:solver:%s", m.SolverAddr))
		writes = append(writes, fmt.Sprintf("intent:solver:%s", m.SolverAddr))

	// ---- Compute (Smart Contracts) ----
	case *computetypes.MsgStoreCode:
		// Store code writes a new code ID; conflicts with other code stores
		writes = append(writes, "compute:next_code_id")

	case *computetypes.MsgInstantiateContract:
		// Instantiation writes a new contract address; conflicts with other instantiations
		writes = append(writes, "compute:next_contract_seq")
		reads = append(reads, fmt.Sprintf("compute:code:%d", m.CodeID))
		reads = append(reads, fmt.Sprintf("bank:%s", m.Sender))
		writes = append(writes, fmt.Sprintf("bank:%s", m.Sender))

	case *computetypes.MsgExecuteContract:
		// Contract execution: conservative -- assume contract state conflicts
		reads = append(reads, fmt.Sprintf("compute:contract:%s", m.Contract))
		writes = append(writes, fmt.Sprintf("compute:contract:%s", m.Contract))
		reads = append(reads, fmt.Sprintf("bank:%s", m.Sender))
		writes = append(writes, fmt.Sprintf("bank:%s", m.Sender))
		writes = append(writes, fmt.Sprintf("bank:%s", m.Contract))

	case *computetypes.MsgMigrateContract:
		reads = append(reads, fmt.Sprintf("compute:contract:%s", m.Contract))
		writes = append(writes, fmt.Sprintf("compute:contract:%s", m.Contract))
		reads = append(reads, fmt.Sprintf("compute:code:%d", m.CodeID))

	default:
		// For unknown msg types, extract signer addresses as read/write keys
		// so that txs from the same signer are never parallelized.
		// Also use msg type URL as a conservative module-level conflict key.
		typeURL := sdk.MsgTypeURL(msg)
		writes = append(writes, fmt.Sprintf("msg:%s", typeURL))
		if signers, err := extractMsgSigners(msg); err == nil {
			for _, signer := range signers {
				addrKey := fmt.Sprintf("acc:%s", signer)
				reads = append(reads, addrKey)
				writes = append(writes, addrKey)
			}
		}
	}
	return
}

// extractMsgSigners attempts to get signer addresses from a message.
func extractMsgSigners(msg sdk.Msg) ([]string, error) {
	type legacySignerMsg interface {
		GetSigners() []sdk.AccAddress
	}
	if lm, ok := msg.(legacySignerMsg); ok {
		signers := lm.GetSigners()
		addrs := make([]string, len(signers))
		for i, s := range signers {
			addrs[i] = s.String()
		}
		return addrs, nil
	}
	return nil, fmt.Errorf("cannot extract signers")
}

// AnalyzeDependencies builds a dependency graph and returns execution groups
// that can be run in parallel. Within each group, all transactions are independent.
//
// All internal access-key tracking uses sorted slices (not maps) to guarantee
// deterministic grouping across all nodes, preventing consensus forks.
func (s *Scheduler) AnalyzeDependencies(tasks []*TxTask) []ExecutionGroup {
	if len(tasks) == 0 {
		return nil
	}

	// For each task, predict what keys it will read/write.
	// PredictAccessKeys already returns deduplicated, sorted slices.
	type txAccess struct {
		reads  []string
		writes []string
	}
	accesses := make([]txAccess, len(tasks))

	for i, task := range tasks {
		if task.Tx == nil {
			// Failed to decode -- must run sequentially
			accesses[i] = txAccess{
				writes: []string{fmt.Sprintf("decode_err:%d", i)},
			}
			continue
		}
		reads, writes := PredictAccessKeys(task.Tx)
		accesses[i] = txAccess{reads: reads, writes: writes}
	}

	// Greedy grouping: assign each tx to the first group where it doesn't conflict.
	// Group access sets use maps for O(1) lookups, but we never iterate them --
	// we only iterate the tx's sorted slices for the conflict check.
	var groups []ExecutionGroup
	type groupAccess struct {
		writes map[string]struct{}
		reads  map[string]struct{}
	}
	groupAccesses := []groupAccess{}

	for i, task := range tasks {
		placed := false
		for g := 0; g < len(groups); g++ {
			if conflictsSorted(accesses[i], groupAccesses[g]) {
				continue
			}
			groups[g].Tasks = append(groups[g].Tasks, task)
			for _, w := range accesses[i].writes {
				groupAccesses[g].writes[w] = struct{}{}
			}
			for _, r := range accesses[i].reads {
				groupAccesses[g].reads[r] = struct{}{}
			}
			placed = true
			break
		}
		if !placed {
			ga := groupAccess{
				writes: make(map[string]struct{}),
				reads:  make(map[string]struct{}),
			}
			for _, w := range accesses[i].writes {
				ga.writes[w] = struct{}{}
			}
			for _, r := range accesses[i].reads {
				ga.reads[r] = struct{}{}
			}
			groups = append(groups, ExecutionGroup{Tasks: []*TxTask{task}})
			groupAccesses = append(groupAccesses, ga)
		}
	}

	return groups
}

// conflictsSorted checks if a tx's access set conflicts with a group's access set.
// The tx's reads and writes are sorted slices, iterated in deterministic order.
// The group's reads and writes are maps used only for O(1) membership lookups.
func conflictsSorted(tx struct {
	reads  []string
	writes []string
}, group struct {
	writes map[string]struct{}
	reads  map[string]struct{}
}) bool {
	// WAW: tx writes something group writes
	// WAR: tx writes something group reads
	for _, w := range tx.writes {
		if _, ok := group.writes[w]; ok {
			return true
		}
		if _, ok := group.reads[w]; ok {
			return true
		}
	}
	// RAW: tx reads something group writes
	for _, r := range tx.reads {
		if _, ok := group.writes[r]; ok {
			return true
		}
	}
	return false
}

// TxFingerprint creates a short hash of a tx for logging.
func TxFingerprint(txBytes []byte) string {
	h := sha256.Sum256(txBytes)
	return fmt.Sprintf("%x", h[:4])
}
