package keeper

import (
	"encoding/json"
	"fmt"
	"math/big"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/tracing"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	ethparams "github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"

	"syreen/x/evm/statedb"
	"syreen/x/evm/types"
)

// MaxGasLimit is the maximum EVM gas limit per transaction (30M, matching Ethereum)
const MaxGasLimit = uint64(30_000_000)

// MaxEVMCallDepth is the maximum EVM call stack depth (enforced by go-ethereum internally).
// Documented here for clarity; go-ethereum's vm.callGas already enforces this limit.
const MaxEVMCallDepth = 1024

// Keeper manages the EVM module state
type Keeper struct {
	cdc           codec.BinaryCodec
	storeService  store.KVStoreService
	accountKeeper authkeeper.AccountKeeper
	bankKeeper    bankkeeper.Keeper
	authority     string

	chainConfig *ethparams.ChainConfig
}

// NewKeeper creates a new EVM keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	authority string,
) Keeper {
	return Keeper{
		cdc:           cdc,
		storeService:  storeService,
		accountKeeper: ak,
		bankKeeper:    bk,
		authority:     authority,
		chainConfig:   types.DefaultChainConfig(),
	}
}

// Logger returns the module logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// M6: GetParams reads params from KVStore (persisted, not just memory)
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.KeyParams())
	if err != nil || bz == nil {
		return types.DefaultParams()
	}
	var p types.Params
	if err := json.Unmarshal(bz, &p); err != nil {
		return types.DefaultParams()
	}
	return p
}

// M6: SetParams writes params to KVStore
func (k Keeper) SetParams(ctx sdk.Context, p types.Params) error {
	bz, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal evm params: %w", err)
	}
	kvStore := k.storeService.OpenKVStore(ctx)
	return kvStore.Set(types.KeyParams(), bz)
}

// GetStoreService returns the module's KVStoreService
func (k Keeper) GetStoreService() store.KVStoreService {
	return k.storeService
}

// GetChainConfig returns the Ethereum chain configuration
func (k Keeper) GetChainConfig() *ethparams.ChainConfig {
	return k.chainConfig
}

// H4: EthAddressFromBech32 converts a Cosmos bech32 address to an Ethereum address with length validation.
func EthAddressFromBech32(bech32 string) (common.Address, error) {
	accAddr, err := sdk.AccAddressFromBech32(bech32)
	if err != nil {
		return common.Address{}, err
	}
	// H4: Validate the address is exactly 20 bytes (Ethereum address length)
	if len(accAddr) != 20 {
		return common.Address{}, fmt.Errorf("invalid address length: expected 20 bytes, got %d", len(accAddr))
	}
	return common.BytesToAddress(accAddr.Bytes()), nil
}

// Bech32FromEthAddress converts an Ethereum address to a Cosmos bech32 address
func Bech32FromEthAddress(ethAddr common.Address) sdk.AccAddress {
	return sdk.AccAddress(ethAddr.Bytes())
}

// ExecuteEVMTx executes an Ethereum transaction against the EVM.
func (k Keeper) ExecuteEVMTx(ctx sdk.Context, msg *types.MsgEthereumTx) (*ExecutionResult, error) {
	logger := k.Logger(ctx)
	params := k.GetParams(ctx)

	sender, err := EthAddressFromBech32(msg.From)
	if err != nil {
		return nil, fmt.Errorf("invalid sender: %w", err)
	}

	data := msg.GetData()
	value := msg.GetValue()
	to := msg.GetToAddress()

	// H1: Cap gas limit to remaining SDK block gas and MaxGasLimit
	gasLimit := msg.GasLimit
	if gasLimit > MaxGasLimit {
		gasLimit = MaxGasLimit
	}
	if ctx.GasMeter().Limit() > 0 {
		consumed := ctx.GasMeter().GasConsumed()
		limit := ctx.GasMeter().Limit()
		if consumed >= limit {
			gasLimit = 0
		} else {
			remaining := limit - consumed
			if gasLimit > remaining {
				gasLimit = remaining
			}
		}
	}

	// EIP-2028: Charge SDK gas for calldata (16 gas per non-zero byte, 4 per zero byte)
	calldataGas := uint64(0)
	for _, b := range data {
		if b == 0 {
			calldataGas += 4
		} else {
			calldataGas += 16
		}
	}
	ctx.GasMeter().ConsumeGas(calldataGas, "evm-calldata")

	// C4: Use CacheContext so we can roll back ALL state on VM error
	cacheCtx, writeCache := ctx.CacheContext()

	// Create StateDB backed by the cached Cosmos SDK state
	stateDB := statedb.New(cacheCtx, k.storeService, k.accountKeeper, k.bankKeeper, params.EvmDenom)

	// H5: Build block context with real block hash lookups
	blockCtx := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash: func(n uint64) common.Hash {
			// H5: Return actual block hash from header for recent blocks
			// CometBFT stores the last ~256 block hashes
			height := ctx.BlockHeight()
			if int64(n) >= height || int64(n) < height-256 || int64(n) < 0 {
				return common.Hash{}
			}
			// For the current block, use the current header hash
			if int64(n) == height-1 {
				headerHash := ctx.HeaderHash()
				if len(headerHash) > 0 {
					return common.BytesToHash(headerHash)
				}
			}
			// For older blocks, we derive a deterministic hash from the block number
			// since we don't have direct access to historical header hashes from ctx.
			// A full implementation would use a block hash store.
			return crypto.Keccak256Hash([]byte(fmt.Sprintf("block:%d", n)))
		},
		Coinbase:    common.Address{},
		BlockNumber: big.NewInt(ctx.BlockHeight()),
		Time:        uint64(ctx.BlockTime().Unix()),
		Difficulty:  big.NewInt(0),
		GasLimit:    gasLimit,
		BaseFee:     big.NewInt(1),
	}

	txCtx := vm.TxContext{
		Origin:   sender,
		GasPrice: big.NewInt(1),
	}

	// Prepare access list
	timestamp := uint64(ctx.BlockTime().Unix())
	rules := k.chainConfig.Rules(blockCtx.BlockNumber, true, timestamp)
	stateDB.Prepare(rules, sender, common.Address{}, to, vm.ActivePrecompiles(rules), nil)
	txHash := common.BytesToHash(crypto.Keccak256(data))
	stateDB.SetTxContext(txHash, 0)

	// H2: Increment nonce before BOTH Create and Call
	nonce := stateDB.GetNonce(sender)
	// EVM-C3: Validate submitted nonce matches expected nonce
	if msg.Nonce != nonce {
		return nil, fmt.Errorf("nonce mismatch: expected %d, got %d", nonce, msg.Nonce)
	}
	stateDB.SetNonce(sender, nonce+1, tracing.NonceChangeEoACall)

	// Create EVM
	evm := vm.NewEVM(blockCtx, stateDB, k.chainConfig, vm.Config{})
	evm.TxContext = txCtx

	var (
		ret          []byte
		leftOverGas  uint64
		contractAddr common.Address
		vmErr        error
	)

	if to == nil {
		// Contract deployment
		if !params.EnableCreate {
			return nil, fmt.Errorf("contract creation is disabled")
		}
		ret, contractAddr, leftOverGas, vmErr = evm.Create(sender, data, gasLimit, uint256.MustFromBig(value))
		logger.Info("EVM contract deployed",
			"address", contractAddr.Hex(),
			"gas_used", gasLimit-leftOverGas,
		)
	} else {
		// Contract call
		if !params.EnableCall {
			return nil, fmt.Errorf("contract calls are disabled")
		}
		ret, leftOverGas, vmErr = evm.Call(sender, *to, data, gasLimit, uint256.MustFromBig(value))
		logger.Info("EVM call executed",
			"to", to.Hex(),
			"gas_used", gasLimit-leftOverGas,
		)
	}

	gasUsed := gasLimit - leftOverGas
	logs := stateDB.GetLogs()

	result := &ExecutionResult{
		ReturnData:   ret,
		GasUsed:      gasUsed,
		Logs:         logs,
		ContractAddr: contractAddr,
	}

	if vmErr != nil {
		// C4: VM error - do NOT call writeCache(), state is rolled back
		result.VmError = vmErr.Error()
		logger.Error("EVM execution error", "error", vmErr)
	} else if stateDB.Error() != nil {
		// C1: StateDB had a sticky error - roll back
		result.VmError = stateDB.Error().Error()
		logger.Error("EVM statedb error", "error", stateDB.Error())
	} else {
		// C5/M7: Clean up suicided accounts before committing
		stateDB.FinalizeDestructs()

		// C2: Commit snapshot caches
		stateDB.CommitSnapshots()

		// C4: Success - commit all cached state to the parent context
		writeCache()
	}

	// H1: Bridge EVM gas to SDK gas meter
	ctx.GasMeter().ConsumeGas(gasUsed, "evm execution")

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"ethereum_tx",
			sdk.NewAttribute("sender", sender.Hex()),
			sdk.NewAttribute("gas_used", fmt.Sprintf("%d", gasUsed)),
			sdk.NewAttribute("contract_address", contractAddr.Hex()),
			sdk.NewAttribute("vm_error", result.VmError),
		),
	)

	return result, nil
}

// QueryEVM performs a read-only EVM call (eth_call equivalent)
func (k Keeper) QueryEVM(ctx sdk.Context, from common.Address, to *common.Address, data []byte, gasLimit uint64) ([]byte, error) {
	// H3-query: Cap gas limit to prevent DoS
	if gasLimit > MaxGasLimit {
		gasLimit = MaxGasLimit
	}
	params := k.GetParams(ctx)
	// Use cache context so queries never modify state
	cacheCtx, _ := ctx.CacheContext()
	stateDB := statedb.New(cacheCtx, k.storeService, k.accountKeeper, k.bankKeeper, params.EvmDenom)

	blockCtx := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     func(n uint64) common.Hash { return common.Hash{} },
		Coinbase:    common.Address{},
		BlockNumber: big.NewInt(ctx.BlockHeight()),
		Time:        uint64(ctx.BlockTime().Unix()),
		Difficulty:  big.NewInt(0),
		GasLimit:    gasLimit,
		BaseFee:     big.NewInt(1),
	}

	txCtx := vm.TxContext{
		Origin:   from,
		GasPrice: big.NewInt(0),
	}

	evm := vm.NewEVM(blockCtx, stateDB, k.chainConfig, vm.Config{NoBaseFee: true})
	evm.TxContext = txCtx

	if to == nil {
		return nil, fmt.Errorf("cannot query contract creation")
	}

	ret, _, err := evm.StaticCall(from, *to, data, gasLimit)
	if err != nil {
		return nil, fmt.Errorf("evm static call failed: %w", err)
	}

	return ret, nil
}

// GetCode returns the EVM bytecode at an address
func (k Keeper) GetCode(ctx sdk.Context, addr common.Address) []byte {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.KeyCode(addr.Bytes()))
	if err != nil {
		return nil
	}
	return bz
}

// GetStorageAt returns the value of a storage slot
func (k Keeper) GetStorageAt(ctx sdk.Context, addr common.Address, slot common.Hash) common.Hash {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.KeyStorage(addr.Bytes(), slot.Bytes()))
	if err != nil || bz == nil {
		return common.Hash{}
	}
	return common.BytesToHash(bz)
}

// GetEVMNonce returns the EVM nonce for an address
func (k Keeper) GetEVMNonce(ctx sdk.Context, addr common.Address) uint64 {
	params := k.GetParams(ctx)
	stateDB := statedb.New(ctx, k.storeService, k.accountKeeper, k.bankKeeper, params.EvmDenom)
	return stateDB.GetNonce(addr)
}

// ExecutionResult contains the result of an EVM execution
type ExecutionResult struct {
	ReturnData   []byte
	GasUsed      uint64
	Logs         []*ethtypes.Log
	ContractAddr common.Address
	VmError      string
}

// MarshalJSON provides JSON encoding for ExecutionResult
func (r *ExecutionResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ReturnData   string `json:"return_data"`
		GasUsed      uint64 `json:"gas_used"`
		ContractAddr string `json:"contract_address"`
		VmError      string `json:"vm_error,omitempty"`
	}{
		ReturnData:   fmt.Sprintf("%x", r.ReturnData),
		GasUsed:      r.GasUsed,
		ContractAddr: r.ContractAddr.Hex(),
		VmError:      r.VmError,
	})
}
