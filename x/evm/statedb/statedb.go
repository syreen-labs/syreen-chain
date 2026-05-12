package statedb

import (
	"math/big"
	"sort"

	"cosmossdk.io/core/store"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/stateless"
	"github.com/ethereum/go-ethereum/core/tracing"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	ethparams "github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie/utils"
	"github.com/holiman/uint256"

	evmtypes "syreen/x/evm/types"
)

var _ interface {
	// Verify we implement the full vm.StateDB interface at compile time.
	// We list the methods explicitly to avoid importing core/vm (circular risk).
	CreateAccount(common.Address)
	CreateContract(common.Address)
	SubBalance(common.Address, *uint256.Int, tracing.BalanceChangeReason) uint256.Int
	AddBalance(common.Address, *uint256.Int, tracing.BalanceChangeReason) uint256.Int
	GetBalance(common.Address) *uint256.Int
	GetNonce(common.Address) uint64
	SetNonce(common.Address, uint64, tracing.NonceChangeReason)
	GetCodeHash(common.Address) common.Hash
	GetCode(common.Address) []byte
	SetCode(common.Address, []byte) []byte
	GetCodeSize(common.Address) int
	AddRefund(uint64)
	SubRefund(uint64)
	GetRefund() uint64
	GetCommittedState(common.Address, common.Hash) common.Hash
	GetState(common.Address, common.Hash) common.Hash
	SetState(common.Address, common.Hash, common.Hash) common.Hash
	GetStorageRoot(common.Address) common.Hash
	GetTransientState(common.Address, common.Hash) common.Hash
	SetTransientState(common.Address, common.Hash, common.Hash)
	SelfDestruct(common.Address) uint256.Int
	HasSelfDestructed(common.Address) bool
	SelfDestruct6780(common.Address) (uint256.Int, bool)
	Exist(common.Address) bool
	Empty(common.Address) bool
	AddressInAccessList(common.Address) bool
	SlotInAccessList(common.Address, common.Hash) (bool, bool)
	AddAddressToAccessList(common.Address)
	AddSlotToAccessList(common.Address, common.Hash)
	PointCache() *utils.PointCache
	Prepare(ethparams.Rules, common.Address, common.Address, *common.Address, []common.Address, ethtypes.AccessList)
	RevertToSnapshot(int)
	Snapshot() int
	AddLog(*ethtypes.Log)
	AddPreimage(common.Hash, []byte)
	Witness() *stateless.Witness
	AccessEvents() *state.AccessEvents
	Finalise(bool)
} = (*StateDB)(nil)

// emptyCodeHash is the code hash of an empty/non-contract account
var emptyCodeHash = crypto.Keccak256Hash(nil)

// MaxCodeSize is the maximum bytecode to permit for a contract (EIP-170: 24KB)
const MaxCodeSize = 24576

// StateDB implements go-ethereum's vm.StateDB interface using Cosmos SDK KVStore.
type StateDB struct {
	ctx           sdk.Context
	storeService  store.KVStoreService
	accountKeeper authkeeper.AccountKeeper
	bankKeeper    bankkeeper.Keeper
	evmDenom      string

	// Transaction-scoped state
	logs       []*ethtypes.Log
	txHash     common.Hash
	txIndex    int
	logIndex   uint
	refund     uint64
	snapshots  []snapshot
	accessList *accessList
	transient  transientStorage

	// Tracking
	suicided map[common.Address]bool
	created  map[common.Address]bool
	dirty    map[common.Address]bool

	// C1: Sticky error field - checked after EVM execution
	dbErr error

	// M2: Committed state tracking - original values at tx start
	committedState map[common.Address]map[common.Hash]common.Hash
}

type snapshot struct {
	refund     uint64
	logLen     int
	suicided   map[common.Address]bool
	dirty      map[common.Address]bool
	created    map[common.Address]bool
	accessList *accessList
	cacheCtx   sdk.Context   // C2: cached context for KVStore revert
	cacheWrite func()        // C2: write function to commit cache
}

// New creates a new StateDB instance
func New(
	ctx sdk.Context,
	storeService store.KVStoreService,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	evmDenom string,
) *StateDB {
	return &StateDB{
		ctx:            ctx,
		storeService:   storeService,
		accountKeeper:  ak,
		bankKeeper:     bk,
		evmDenom:       evmDenom,
		logs:           make([]*ethtypes.Log, 0),
		suicided:       make(map[common.Address]bool),
		created:        make(map[common.Address]bool),
		dirty:          make(map[common.Address]bool),
		accessList:     newAccessList(),
		transient:      newTransientStorage(),
		committedState: make(map[common.Address]map[common.Hash]common.Hash),
	}
}

// kvStore opens the KVStore from the store service using the current context.
func (s *StateDB) kvStore() store.KVStore {
	return s.storeService.OpenKVStore(s.ctx)
}

// Error returns the first sticky error that occurred during state operations (C1)
func (s *StateDB) Error() error {
	return s.dbErr
}

// setError records the first error encountered (C1)
func (s *StateDB) setError(err error) {
	if s.dbErr == nil && err != nil {
		s.dbErr = err
	}
}

// SetTxContext sets the transaction hash and index for log attribution
func (s *StateDB) SetTxContext(txHash common.Hash, txIndex int) {
	s.txHash = txHash
	s.txIndex = txIndex
}

// --- Account Operations ---

func (s *StateDB) CreateAccount(addr common.Address) {
	kvs := s.kvStore()
	if err := kvs.Set(evmtypes.KeyAccount(addr.Bytes()), []byte{0x01}); err != nil {
		s.setError(err)
	}
	s.created[addr] = true
	s.dirty[addr] = true
}

func (s *StateDB) CreateContract(addr common.Address) {
	// Mark as created contract (same as CreateAccount for our purposes)
	s.CreateAccount(addr)
}

func (s *StateDB) Exist(addr common.Address) bool {
	kvs := s.kvStore()
	has, err := kvs.Has(evmtypes.KeyAccount(addr.Bytes()))
	if err == nil && has {
		return true
	}
	cosmosAddr := sdk.AccAddress(addr.Bytes())
	acct := s.accountKeeper.GetAccount(s.ctx, cosmosAddr)
	return acct != nil
}

func (s *StateDB) Empty(addr common.Address) bool {
	return s.GetNonce(addr) == 0 &&
		s.GetBalance(addr).IsZero() &&
		s.GetCodeHash(addr) == emptyCodeHash
}

// --- Balance Operations ---

func (s *StateDB) GetBalance(addr common.Address) *uint256.Int {
	cosmosAddr := sdk.AccAddress(addr.Bytes())
	coin := s.bankKeeper.GetBalance(s.ctx, cosmosAddr, s.evmDenom)
	bi := coin.Amount.BigInt()
	val, _ := uint256.FromBig(bi)
	if val == nil {
		val = new(uint256.Int)
	}
	return val
}

// GetBalanceBig returns the balance as *big.Int (internal helper for mint/burn).
func (s *StateDB) GetBalanceBig(addr common.Address) *big.Int {
	cosmosAddr := sdk.AccAddress(addr.Bytes())
	coin := s.bankKeeper.GetBalance(s.ctx, cosmosAddr, s.evmDenom)
	return coin.Amount.BigInt()
}

// Transfer moves funds from one account to another using bankKeeper.SendCoins.
// This is the correct method for EVM value transfers (CALL, CALLCODE, etc.)
// because it preserves total supply invariant: no tokens are minted or burned.
//
// WARNING: totalSupply invariant -- always prefer Transfer for account-to-account
// movements. Only use AddBalance/SubBalance for bridging operations where tokens
// genuinely enter/leave the EVM module supply.
func (s *StateDB) Transfer(from, to common.Address, amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	if amount.Sign() < 0 {
		s.setError(ErrNegativeAmount)
		return
	}

	fromCosmos := sdk.AccAddress(from.Bytes())
	toCosmos := sdk.AccAddress(to.Bytes())
	coins := sdk.NewCoins(sdk.NewCoin(s.evmDenom, math.NewIntFromBigInt(amount)))

	// Ensure recipient account exists
	if s.accountKeeper.GetAccount(s.ctx, toCosmos) == nil {
		acct := s.accountKeeper.NewAccountWithAddress(s.ctx, toCosmos)
		s.accountKeeper.SetAccount(s.ctx, acct)
	}

	if err := s.bankKeeper.SendCoins(s.ctx, fromCosmos, toCosmos, coins); err != nil {
		s.setError(err)
		return
	}
	s.dirty[from] = true
	s.dirty[to] = true
}

// AddBalance mints new tokens and credits them to addr.
// WARNING: This increases total supply. Only use for bridging tokens INTO the EVM.
func (s *StateDB) AddBalance(addr common.Address, amount *uint256.Int, reason tracing.BalanceChangeReason) uint256.Int {
	prev := *s.GetBalance(addr)
	if amount.IsZero() {
		return prev
	}
	amountBig := amount.ToBig()
	if amountBig.Sign() < 0 {
		s.setError(ErrNegativeAmount)
		return prev
	}
	cosmosAddr := sdk.AccAddress(addr.Bytes())
	coins := sdk.NewCoins(sdk.NewCoin(s.evmDenom, math.NewIntFromBigInt(amountBig)))

	// Ensure account exists
	if s.accountKeeper.GetAccount(s.ctx, cosmosAddr) == nil {
		acct := s.accountKeeper.NewAccountWithAddress(s.ctx, cosmosAddr)
		s.accountKeeper.SetAccount(s.ctx, acct)
	}

	// C1: Record errors instead of silently swallowing them
	if err := s.bankKeeper.MintCoins(s.ctx, evmtypes.ModuleName, coins); err != nil {
		s.setError(err)
		return prev
	}
	if err := s.bankKeeper.SendCoinsFromModuleToAccount(s.ctx, evmtypes.ModuleName, cosmosAddr, coins); err != nil {
		s.setError(err)
		return prev
	}
	s.dirty[addr] = true
	return prev
}

// SubBalance burns tokens from addr.
// WARNING: This decreases total supply. Only use for bridging tokens OUT of the EVM.
func (s *StateDB) SubBalance(addr common.Address, amount *uint256.Int, reason tracing.BalanceChangeReason) uint256.Int {
	prev := *s.GetBalance(addr)
	if amount.IsZero() {
		return prev
	}
	amountBig := amount.ToBig()
	if amountBig.Sign() < 0 {
		s.setError(ErrNegativeAmount)
		return prev
	}
	cosmosAddr := sdk.AccAddress(addr.Bytes())
	coins := sdk.NewCoins(sdk.NewCoin(s.evmDenom, math.NewIntFromBigInt(amountBig)))

	// C1: Record errors instead of silently swallowing them
	if err := s.bankKeeper.SendCoinsFromAccountToModule(s.ctx, cosmosAddr, evmtypes.ModuleName, coins); err != nil {
		s.setError(err)
		return prev
	}
	if err := s.bankKeeper.BurnCoins(s.ctx, evmtypes.ModuleName, coins); err != nil {
		s.setError(err)
		return prev
	}
	s.dirty[addr] = true
	return prev
}

// --- Nonce Operations ---
//
// DUAL-NONCE MODEL: The EVM module maintains its own nonce in the EVM KVStore
// (KeyNonce), separate from the Cosmos SDK account sequence number managed by
// x/auth. GetNonce checks the EVM store first; if not found, it falls back to
// the Cosmos account sequence. SetNonce only writes to the EVM store. This
// means EVM nonces and Cosmos nonces can diverge if the same account sends
// both Cosmos and EVM transactions. This is by design: EVM nonces follow
// Ethereum semantics, while Cosmos sequences follow SDK semantics.

func (s *StateDB) GetNonce(addr common.Address) uint64 {
	kvs := s.kvStore()
	bz, err := kvs.Get(evmtypes.KeyNonce(addr.Bytes()))
	if err != nil {
		s.setError(err)
	}

	// Always check Cosmos account sequence and return the higher of the two
	// to prevent nonce divergence when the same account sends both EVM and
	// Cosmos SDK transactions.
	var evmNonce uint64
	if bz != nil {
		evmNonce = bytesToUint64(bz)
	}

	cosmosAddr := sdk.AccAddress(addr.Bytes())
	acct := s.accountKeeper.GetAccount(s.ctx, cosmosAddr)
	if acct != nil {
		cosmosNonce := acct.GetSequence()
		if cosmosNonce > evmNonce {
			return cosmosNonce
		}
	}
	return evmNonce
}

func (s *StateDB) SetNonce(addr common.Address, nonce uint64, reason tracing.NonceChangeReason) {
	kvs := s.kvStore()
	if err := kvs.Set(evmtypes.KeyNonce(addr.Bytes()), uint64ToBytes(nonce)); err != nil {
		s.setError(err)
	}
	s.dirty[addr] = true

	// Sync Cosmos account sequence to prevent EVM/Cosmos nonce divergence.
	// If the account exists, update its sequence to match the EVM nonce.
	acct := s.accountKeeper.GetAccount(s.ctx, sdk.AccAddress(addr.Bytes()))
	if acct != nil {
		if err := acct.SetSequence(nonce); err == nil {
			s.accountKeeper.SetAccount(s.ctx, acct)
		}
	}
}

// --- Code Operations ---

func (s *StateDB) GetCodeHash(addr common.Address) common.Hash {
	kvs := s.kvStore()
	bz, err := kvs.Get(evmtypes.KeyCodeHash(addr.Bytes()))
	if err != nil {
		s.setError(err)
	}
	if bz == nil {
		if s.Exist(addr) {
			return emptyCodeHash
		}
		return common.Hash{}
	}
	return common.BytesToHash(bz)
}

func (s *StateDB) GetCode(addr common.Address) []byte {
	kvs := s.kvStore()
	bz, err := kvs.Get(evmtypes.KeyCode(addr.Bytes()))
	if err != nil {
		s.setError(err)
		return nil
	}
	return bz
}

func (s *StateDB) SetCode(addr common.Address, code []byte) []byte {
	prev := s.GetCode(addr)

	// M3: Enforce EIP-170 code size limit (24KB)
	if len(code) > MaxCodeSize {
		s.setError(ErrCodeTooLarge)
		return prev
	}

	kvs := s.kvStore()
	if len(code) == 0 {
		if err := kvs.Delete(evmtypes.KeyCode(addr.Bytes())); err != nil {
			s.setError(err)
		}
		if err := kvs.Delete(evmtypes.KeyCodeHash(addr.Bytes())); err != nil {
			s.setError(err)
		}
	} else {
		if err := kvs.Set(evmtypes.KeyCode(addr.Bytes()), code); err != nil {
			s.setError(err)
		}
		hash := crypto.Keccak256Hash(code)
		if err := kvs.Set(evmtypes.KeyCodeHash(addr.Bytes()), hash.Bytes()); err != nil {
			s.setError(err)
		}
	}
	has, err := kvs.Has(evmtypes.KeyAccount(addr.Bytes()))
	if err != nil {
		s.setError(err)
	}
	if !has {
		if err := kvs.Set(evmtypes.KeyAccount(addr.Bytes()), []byte{0x01}); err != nil {
			s.setError(err)
		}
	}
	s.dirty[addr] = true
	return prev
}

func (s *StateDB) GetCodeSize(addr common.Address) int {
	return len(s.GetCode(addr))
}

// --- Storage Operations ---

func (s *StateDB) GetState(addr common.Address, key common.Hash) common.Hash {
	kvs := s.kvStore()
	bz, err := kvs.Get(evmtypes.KeyStorage(addr.Bytes(), key.Bytes()))
	if err != nil {
		s.setError(err)
		return common.Hash{}
	}
	if bz == nil {
		return common.Hash{}
	}
	return common.BytesToHash(bz)
}

func (s *StateDB) SetState(addr common.Address, key common.Hash, value common.Hash) common.Hash {
	// M2: Track original value for GetCommittedState before overwriting
	if _, ok := s.committedState[addr]; !ok {
		s.committedState[addr] = make(map[common.Hash]common.Hash)
	}
	if _, tracked := s.committedState[addr][key]; !tracked {
		// First write in this tx: record the original value
		s.committedState[addr][key] = s.GetState(addr, key)
	}

	prev := s.GetState(addr, key)

	kvs := s.kvStore()
	if value == (common.Hash{}) {
		if err := kvs.Delete(evmtypes.KeyStorage(addr.Bytes(), key.Bytes())); err != nil {
			s.setError(err)
		}
	} else {
		if err := kvs.Set(evmtypes.KeyStorage(addr.Bytes(), key.Bytes()), value.Bytes()); err != nil {
			s.setError(err)
		}
	}
	s.dirty[addr] = true
	return prev
}

// M2: GetCommittedState returns the value at the start of the transaction, not the current value
func (s *StateDB) GetCommittedState(addr common.Address, key common.Hash) common.Hash {
	if addrState, ok := s.committedState[addr]; ok {
		if val, tracked := addrState[key]; tracked {
			return val
		}
	}
	// Not modified in this tx, so current state IS the committed state
	return s.GetState(addr, key)
}

// GetStorageRoot returns a dummy storage root hash. Our KV-based storage
// does not maintain a Merkle trie per account, so we return a deterministic
// placeholder. go-ethereum only uses this for Verkle witness generation
// which we do not support.
func (s *StateDB) GetStorageRoot(addr common.Address) common.Hash {
	// Return empty hash for non-existent accounts
	if !s.Exist(addr) {
		return common.Hash{}
	}
	// Return a non-empty deterministic hash for existing accounts
	return crypto.Keccak256Hash(addr.Bytes())
}

// --- Transient Storage (EIP-1153) ---

func (s *StateDB) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	return s.transient.Get(addr, key)
}

func (s *StateDB) SetTransientState(addr common.Address, key, value common.Hash) {
	s.transient.Set(addr, key, value)
}

// --- Refund ---

func (s *StateDB) AddRefund(gas uint64) {
	s.refund += gas
}

// M5: SubRefund clamps at zero on underflow instead of panicking.
// go-ethereum panics here, but in our fork we prefer defensive behavior:
// a buggy precompile or opcode metering miscalculation should not be able
// to halt the chain. We log the underflow attempt via the sticky dbErr so
// the EVM caller can still observe the failure, then clamp to zero.
func (s *StateDB) SubRefund(gas uint64) {
	if gas > s.refund {
		s.ctx.Logger().Error(
			"EVM SubRefund underflow -- clamping refund counter to zero",
			"requested_sub", gas,
			"current_refund", s.refund,
		)
		s.refund = 0
		// Mark the transaction as failed so the EVM caller can detect the anomaly
		s.setError(ErrRefundUnderflow)
		return
	}
	s.refund -= gas
}

func (s *StateDB) GetRefund() uint64 {
	return s.refund
}

// --- Self-Destruct (formerly Suicide) ---

func (s *StateDB) SelfDestruct(addr common.Address) uint256.Int {
	prev := *s.GetBalance(addr)
	if !s.Exist(addr) {
		return prev
	}
	s.suicided[addr] = true
	if !prev.IsZero() {
		s.SubBalance(addr, &prev, tracing.BalanceDecreaseSelfdestructBurn)
	}
	return prev
}

func (s *StateDB) HasSelfDestructed(addr common.Address) bool {
	return s.suicided[addr]
}

// SelfDestruct6780 implements post-EIP6780 selfdestruct semantics.
// It only truly destructs if the contract was created in this same transaction.
func (s *StateDB) SelfDestruct6780(addr common.Address) (uint256.Int, bool) {
	bal := *s.GetBalance(addr)
	if s.created[addr] {
		// Created in this tx: full destruct
		s.SelfDestruct(addr)
		return bal, true
	}
	// Not created in this tx: just send balance to beneficiary (handled by caller),
	// do not mark as suicided
	return bal, false
}

// FinalizeDestructs processes suicided accounts, removing code, storage, and nonce.
// Must be called after EVM execution completes (end of transaction).
func (s *StateDB) FinalizeDestructs() {
	// Sort suicided addresses for deterministic cleanup order across validators.
	addrs := make([]common.Address, 0, len(s.suicided))
	for addr := range s.suicided {
		addrs = append(addrs, addr)
	}
	sort.Slice(addrs, func(i, j int) bool {
		return addrs[i].Hex() < addrs[j].Hex()
	})

	for _, addr := range addrs {
		// Delete code
		s.SetCode(addr, nil)

		kvs := s.kvStore()

		// Delete nonce
		if err := kvs.Delete(evmtypes.KeyNonce(addr.Bytes())); err != nil {
			s.setError(err)
		}

		// Delete account marker
		if err := kvs.Delete(evmtypes.KeyAccount(addr.Bytes())); err != nil {
			s.setError(err)
		}

		// Delete all storage for this account
		// Use iterator with prefix range to find all storage keys
		storagePrefix := evmtypes.KeyStorage(addr.Bytes(), nil)
		endPrefix := make([]byte, len(storagePrefix))
		copy(endPrefix, storagePrefix)
		endPrefix = append(endPrefix, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF)
		iter, err := kvs.Iterator(storagePrefix, endPrefix)
		if err != nil {
			s.setError(err)
			continue
		}
		keysToDelete := make([][]byte, 0)
		for ; iter.Valid(); iter.Next() {
			keyCopy := make([]byte, len(iter.Key()))
			copy(keyCopy, iter.Key())
			keysToDelete = append(keysToDelete, keyCopy)
		}
		iter.Close()
		for _, key := range keysToDelete {
			if err := kvs.Delete(key); err != nil {
				s.setError(err)
			}
		}
	}
}

// Finalise implements vm.StateDB. It is called at the end of a transaction.
// The deleteEmptyObjects parameter is ignored in our implementation; we always
// clean up suicided accounts.
func (s *StateDB) Finalise(deleteEmptyObjects bool) {
	s.FinalizeDestructs()
}

// --- Access List (EIP-2929) ---

func (s *StateDB) AddressInAccessList(addr common.Address) bool {
	return s.accessList.ContainsAddress(addr)
}

func (s *StateDB) SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool) {
	return s.accessList.Contains(addr, slot)
}

func (s *StateDB) AddAddressToAccessList(addr common.Address) {
	s.accessList.AddAddress(addr)
}

func (s *StateDB) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	s.accessList.AddSlot(addr, slot)
}

// Prepare implements vm.StateDB. It resets the access list and sets up
// transaction-scoped state per EIP-2929 and EIP-3651.
func (s *StateDB) Prepare(rules ethparams.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses ethtypes.AccessList) {
	s.accessList = newAccessList()
	s.AddAddressToAccessList(sender)
	if dest != nil {
		s.AddAddressToAccessList(*dest)
	}
	// EIP-3651: warm coinbase
	s.AddAddressToAccessList(coinbase)
	for _, addr := range precompiles {
		s.AddAddressToAccessList(addr)
	}
	for _, el := range txAccesses {
		s.AddAddressToAccessList(el.Address)
		for _, key := range el.StorageKeys {
			s.AddSlotToAccessList(el.Address, key)
		}
	}
	// Reset transient storage at the beginning of each transaction
	s.transient = newTransientStorage()
}

// --- Snapshots ---
// C2: Snapshot/RevertToSnapshot now uses CacheContext so KVStore writes are reverted

func (s *StateDB) Snapshot() int {
	id := len(s.snapshots)

	// Copy suicided map
	sd := make(map[common.Address]bool)
	for k, v := range s.suicided {
		sd[k] = v
	}

	// Copy dirty map
	dd := make(map[common.Address]bool)
	for k, v := range s.dirty {
		dd[k] = v
	}

	// Copy created map
	cd := make(map[common.Address]bool)
	for k, v := range s.created {
		cd[k] = v
	}

	// Copy access list
	alCopy := s.accessList.copy()

	// C2: Create a cache context to capture all KVStore writes from this point
	cacheCtx, write := s.ctx.CacheContext()

	s.snapshots = append(s.snapshots, snapshot{
		refund:     s.refund,
		logLen:     len(s.logs),
		suicided:   sd,
		dirty:      dd,
		created:    cd,
		accessList: alCopy,
		cacheCtx:   s.ctx, // save the CURRENT context (parent)
		cacheWrite: write,
	})

	// Switch to the cached context for future operations
	s.ctx = cacheCtx

	return id
}

func (s *StateDB) RevertToSnapshot(id int) {
	if id >= len(s.snapshots) {
		return
	}
	snap := s.snapshots[id]
	s.refund = snap.refund
	s.logs = s.logs[:snap.logLen]
	s.suicided = snap.suicided
	s.dirty = snap.dirty
	s.created = snap.created
	s.accessList = snap.accessList

	// C2: Restore the parent context, discarding all KVStore writes since the snapshot
	// By NOT calling snap.cacheWrite(), the cached writes are discarded
	s.ctx = snap.cacheCtx

	s.snapshots = s.snapshots[:id]
}

// CommitSnapshots writes all pending snapshot caches to the parent store.
// Must be called after successful EVM execution to persist cached KVStore writes.
func (s *StateDB) CommitSnapshots() {
	// Write snapshot caches from bottom to top (oldest first)
	for _, snap := range s.snapshots {
		if snap.cacheWrite != nil {
			snap.cacheWrite()
		}
	}
	s.snapshots = nil
}

// --- Logs ---

func (s *StateDB) AddLog(log *ethtypes.Log) {
	log.TxHash = s.txHash
	log.TxIndex = uint(s.txIndex)
	log.Index = s.logIndex
	s.logIndex++
	s.logs = append(s.logs, log)
}

func (s *StateDB) GetLogs() []*ethtypes.Log {
	return s.logs
}

func (s *StateDB) AddPreimage(hash common.Hash, preimage []byte) {
	// Preimages are SHA3/Keccak256 inputs that map to a 32-byte hash.
	// Bound the data length to 32 bytes to prevent unbounded storage growth.
	if len(preimage) > 32 {
		return
	}
	kvs := s.kvStore()
	key := append([]byte("preimage/"), hash.Bytes()...)
	if err := kvs.Set(key, preimage); err != nil {
		s.setError(err)
	}
}

// --- Verkle / Witness stubs ---

// PointCache returns nil. Verkle trees are not supported.
func (s *StateDB) PointCache() *utils.PointCache {
	return nil
}

// Witness returns nil. Stateless witnesses are not supported.
func (s *StateDB) Witness() *stateless.Witness {
	return nil
}

// AccessEvents returns nil. Verkle access events are not supported.
func (s *StateDB) AccessEvents() *state.AccessEvents {
	return nil
}

// --- ForEachStorage ---

func (s *StateDB) ForEachStorage(addr common.Address, cb func(common.Hash, common.Hash) bool) error {
	kvs := s.kvStore()
	storagePrefix := evmtypes.KeyStorage(addr.Bytes(), nil)
	endPrefix := make([]byte, len(storagePrefix))
	copy(endPrefix, storagePrefix)
	endPrefix = append(endPrefix, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF)

	iter, err := kvs.Iterator(storagePrefix, endPrefix)
	if err != nil {
		return err
	}
	defer iter.Close()

	prefixLen := len(storagePrefix)
	for ; iter.Valid(); iter.Next() {
		rawKey := iter.Key()
		// Strip the prefix to get just the slot bytes
		slotBytes := rawKey[prefixLen:]
		key := common.BytesToHash(slotBytes)
		value := common.BytesToHash(iter.Value())
		if !cb(key, value) {
			break
		}
	}
	return nil
}

// --- Utility ---

func uint64ToBytes(v uint64) []byte {
	bz := make([]byte, 8)
	bz[0] = byte(v >> 56)
	bz[1] = byte(v >> 48)
	bz[2] = byte(v >> 40)
	bz[3] = byte(v >> 32)
	bz[4] = byte(v >> 24)
	bz[5] = byte(v >> 16)
	bz[6] = byte(v >> 8)
	bz[7] = byte(v)
	return bz
}

func bytesToUint64(bz []byte) uint64 {
	if len(bz) < 8 {
		return 0
	}
	return uint64(bz[0])<<56 | uint64(bz[1])<<48 | uint64(bz[2])<<40 |
		uint64(bz[3])<<32 | uint64(bz[4])<<24 | uint64(bz[5])<<16 |
		uint64(bz[6])<<8 | uint64(bz[7])
}
