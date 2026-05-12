package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"cosmossdk.io/core/store"
	"cosmossdk.io/math"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/compute/types"
	"syreen/x/compute/vm"
)

type Keeper struct {
	cdc           codec.Codec
	storeService  store.KVStoreService
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	authority     string
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	authority string,
) Keeper {
	return Keeper{
		cdc:           cdc,
		storeService:  storeService,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		authority:     authority,
	}
}

func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Create stores code bytecode and returns a code ID
func (k Keeper) Create(ctx context.Context, creator string, wasmCode []byte, instantiateAccess *types.AccessConfig) (uint64, error) {
	params := k.GetParams(ctx)

	// Check upload permission
	if !k.isUploadAllowed(params.CodeUploadAccess, creator) {
		return 0, types.ErrUploadDenied
	}

	// Validate code size
	if uint64(len(wasmCode)) > params.MaxWasmCodeSize {
		return 0, types.ErrCodeTooLarge
	}

	// Get next code ID
	codeID := k.getNextCodeID(ctx)
	k.setNextCodeID(ctx, codeID+1)

	// Compute code hash
	hash := sha256.Sum256(wasmCode)

	// Determine instantiate permission
	permission := types.AccessConfig{Permission: params.InstantiateDefaultPermission}
	if instantiateAccess != nil {
		permission = *instantiateAccess
	}

	// Store code info
	codeInfo := types.CodeInfo{
		CodeID:                codeID,
		Creator:               creator,
		CodeHash:              hash[:],
		InstantiatePermission: permission,
	}
	k.setCodeInfo(ctx, codeID, codeInfo)

	// Store bytecode
	k.setCodeBytecode(ctx, codeID, wasmCode)

	k.Logger(ctx).Info("stored code", "code_id", codeID, "creator", creator, "size", len(wasmCode))
	return codeID, nil
}

// MaxWasmCodeSize is the hard maximum bytecode size for WASM execution (500KB).
// This prevents DoS via oversized contract code regardless of module params.
const MaxWasmCodeSize = 500 * 1024

// Instantiate creates a new contract instance from stored code
func (k Keeper) Instantiate(ctx context.Context, codeID uint64, creator, admin, label string, initMsg json.RawMessage, funds sdk.Coins) (string, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Verify code exists
	codeInfo, found := k.GetCodeInfo(ctx, codeID)
	if !found {
		return "", types.ErrCodeNotFound
	}

	// Check instantiate permission
	if !k.isInstantiateAllowed(codeInfo.InstantiatePermission, creator) {
		return "", types.ErrInstantiateDenied
	}

	// Generate contract address deterministically
	contractAddr := k.generateContractAddress(ctx, codeID, creator, label)

	// Ensure the account exists
	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return "", types.ErrInvalidCreator
	}

	contractAccAddr, err := sdk.AccAddressFromBech32(contractAddr)
	if err != nil {
		return "", fmt.Errorf("invalid generated contract address: %w", err)
	}

	// Create account for contract if it doesn't exist
	if k.accountKeeper.GetAccount(ctx, contractAccAddr) == nil {
		acc := k.accountKeeper.NewAccountWithAddress(ctx, contractAccAddr)
		k.accountKeeper.SetAccount(ctx, acc)
	}

	// Transfer funds if any
	if funds.IsAllPositive() {
		if err := k.bankKeeper.SendCoins(ctx, creatorAddr, contractAccAddr, funds); err != nil {
			return "", err
		}
	}

	// Store contract info
	contractInfo := types.ContractInfo{
		Address:   contractAddr,
		CodeID:    codeID,
		Creator:   creator,
		Admin:     admin,
		Label:     label,
		CreatedAt: sdkCtx.BlockTime().Unix(),
	}
	k.setContractInfo(ctx, contractAddr, contractInfo)

	// Load the contract code (JSON contract definition)
	code := k.getCodeBytecode(ctx, codeID)
	if code == nil {
		// Legacy/non-VM code: just store the init message
		k.SetContractState(ctx, contractAddr, []byte("_init_msg"), initMsg)
		k.Logger(ctx).Info("instantiated contract (legacy)", "address", contractAddr, "code_id", codeID, "creator", creator)
		return contractAddr, nil
	}

	// Enforce hard bytecode size limit
	if len(code) > MaxWasmCodeSize {
		return "", fmt.Errorf("wasm code size %d exceeds maximum %d", len(code), MaxWasmCodeSize)
	}

	// Charge gas proportional to bytecode size before execution
	sdkCtx.GasMeter().ConsumeGas(uint64(len(code))+10000, "wasm-instantiate")

	// Try to parse as a VM contract definition
	def, parseErr := vm.ParseContractDefinition(code)
	if parseErr != nil {
		// Not a VM contract -- fall back to legacy behavior
		k.SetContractState(ctx, contractAddr, []byte("_init_msg"), initMsg)
		k.Logger(ctx).Info("instantiated contract (legacy)", "address", contractAddr, "code_id", codeID, "creator", creator)
		return contractAddr, nil
	}

	// Parse the init message
	var msgMap map[string]interface{}
	if err := json.Unmarshal(initMsg, &msgMap); err != nil {
		return "", types.ErrInvalidMsg
	}

	// Run the VM instantiate handler
	params := k.GetParams(ctx)
	blockHeight := strconv.FormatInt(sdkCtx.BlockHeight(), 10)
	execCtx := vm.NewExecutionContext(contractAddr, creator, blockHeight, nil, msgMap, params.MaxContractGas)

	engine := vm.NewEngine()
	if err := engine.ExecuteInstantiate(def, execCtx); err != nil {
		return "", fmt.Errorf("contract instantiate failed: %w", err)
	}

	// Wrap state save and bank messages in a single CacheContext for atomicity.
	{
		cacheCtx, writeCache := sdkCtx.CacheContext()
		if err := k.saveVMState(cacheCtx, contractAddr, execCtx.State); err != nil {
			return "", fmt.Errorf("contract instantiate state save failed: %w", err)
		}
		if err := k.executeBankMsgs(cacheCtx, execCtx.BankMsgs); err != nil {
			return "", fmt.Errorf("contract bank msg failed: %w", err)
		}
		writeCache() // commit both state save and bank msgs atomically
	}

	k.Logger(ctx).Info("instantiated contract", "address", contractAddr, "code_id", codeID, "creator", creator)
	return contractAddr, nil
}

// Execute calls a smart contract with the given message, running the JSON VM engine.
func (k Keeper) Execute(ctx context.Context, contractAddr, sender string, msg json.RawMessage, funds sdk.Coins) ([]byte, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Verify contract exists
	contractInfo, found := k.GetContractInfo(ctx, contractAddr)
	if !found {
		return nil, types.ErrContractNotFound
	}

	// Verify sender is valid
	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return nil, types.ErrInvalidCreator
	}

	// Transfer funds if any
	if funds.IsAllPositive() {
		contractAccAddr, err := sdk.AccAddressFromBech32(contractAddr)
		if err != nil {
			return nil, types.ErrContractNotFound
		}
		if err := k.bankKeeper.SendCoins(ctx, senderAddr, contractAccAddr, funds); err != nil {
			return nil, err
		}
	}

	// Load contract code
	code := k.getCodeBytecode(ctx, contractInfo.CodeID)
	if code == nil {
		return nil, fmt.Errorf("contract code not found for code_id %d", contractInfo.CodeID)
	}

	// Enforce hard bytecode size limit
	if len(code) > MaxWasmCodeSize {
		return nil, fmt.Errorf("wasm code size %d exceeds maximum %d", len(code), MaxWasmCodeSize)
	}

	// Charge gas proportional to bytecode size before execution
	sdkCtx.GasMeter().ConsumeGas(uint64(len(code))+10000, "wasm-execution")

	// Parse the contract definition
	def, parseErr := vm.ParseContractDefinition(code)
	if parseErr != nil {
		// Fall back to legacy placeholder behavior for non-VM contracts
		return k.executeLegacy(ctx, contractAddr, sender, msg)
	}

	// Determine which handler to call from the message
	handlerName, handlerMsg, err := vm.DetermineHandler(msg)
	if err != nil {
		return nil, types.ErrInvalidMsg
	}

	// Load current contract state from KV store
	state := k.loadVMState(ctx, contractAddr)

	// Create execution context
	params := k.GetParams(ctx)
	blockHeight := strconv.FormatInt(sdkCtx.BlockHeight(), 10)
	execCtx := vm.NewExecutionContext(contractAddr, sender, blockHeight, state, handlerMsg, params.MaxContractGas)

	// Run the VM
	engine := vm.NewEngine()
	if err := engine.Execute(def, handlerName, execCtx); err != nil {
		return nil, fmt.Errorf("contract execution failed: %w", err)
	}

	// Bridge compute VM gas to Cosmos SDK gas meter.
	// Conversion factor: 1000 internal VM gas units = 1 SDK gas unit.
	sdkCtx.GasMeter().ConsumeGas(execCtx.GasUsed/1000, "compute execution")

	// Wrap state save and bank messages in a single CacheContext for atomicity.
	// If either fails, all changes are discarded — no partial state is committed.
	{
		cacheCtx, writeCache := sdkCtx.CacheContext()
		if execCtx.StateChanged {
			if err := k.saveVMState(cacheCtx, contractAddr, execCtx.State); err != nil {
				return nil, fmt.Errorf("contract execution state save failed: %w", err)
			}
		}
		if err := k.executeBankMsgs(cacheCtx, execCtx.BankMsgs); err != nil {
			return nil, fmt.Errorf("contract bank msg failed: %w", err)
		}
		writeCache() // commit both state save and bank msgs atomically
	}

	// Emit events
	k.emitVMEvents(ctx, execCtx.Events)

	k.Logger(ctx).Info("executed contract", "address", contractAddr, "sender", sender, "handler", handlerName)

	if execCtx.Response != nil {
		return execCtx.Response, nil
	}

	resp, _ := json.Marshal(map[string]string{"status": "ok"})
	return resp, nil
}

// executeLegacy is the old placeholder execution path for non-VM contracts
func (k Keeper) executeLegacy(ctx context.Context, contractAddr, sender string, msg json.RawMessage) ([]byte, error) {
	var msgMap map[string]json.RawMessage
	if err := json.Unmarshal(msg, &msgMap); err != nil {
		return nil, types.ErrInvalidMsg
	}

	const maxKeys = 50
	if len(msgMap) > maxKeys {
		return nil, fmt.Errorf("too many state keys in execution message: %d (max %d)", len(msgMap), maxKeys)
	}

	// Sort map keys for deterministic state writes across all validators.
	msgKeys := make([]string, 0, len(msgMap))
	for key := range msgMap {
		msgKeys = append(msgKeys, key)
	}
	sort.Strings(msgKeys)

	for _, key := range msgKeys {
		value := msgMap[key]
		stateKey := fmt.Sprintf("exec/%s/%s", sender, key)
		k.SetContractState(ctx, contractAddr, []byte(stateKey), value)
	}

	resp, _ := json.Marshal(map[string]string{"status": "ok"})
	return resp, nil
}

// Migrate runs a code upgrade for a contract
func (k Keeper) Migrate(ctx context.Context, contractAddr, sender string, newCodeID uint64, msg json.RawMessage) ([]byte, error) {
	contractInfo, found := k.GetContractInfo(ctx, contractAddr)
	if !found {
		return nil, types.ErrContractNotFound
	}

	// Only admin can migrate
	if contractInfo.Admin == "" || contractInfo.Admin != sender {
		return nil, types.ErrUnauthorized
	}

	// Verify new code exists
	_, found = k.GetCodeInfo(ctx, newCodeID)
	if !found {
		return nil, types.ErrCodeNotFound
	}

	// Update contract code ID
	contractInfo.CodeID = newCodeID
	k.setContractInfo(ctx, contractAddr, contractInfo)

	// Store migration message
	k.SetContractState(ctx, contractAddr, []byte("_migrate_msg"), msg)

	k.Logger(ctx).Info("migrated contract", "address", contractAddr, "new_code_id", newCodeID)

	resp, _ := json.Marshal(map[string]string{"status": "migrated"})
	return resp, nil
}

// UpdateAdmin sets a new admin for a contract
func (k Keeper) UpdateAdmin(ctx context.Context, contractAddr, sender, newAdmin string) error {
	contractInfo, found := k.GetContractInfo(ctx, contractAddr)
	if !found {
		return types.ErrContractNotFound
	}

	if contractInfo.Admin != sender {
		return types.ErrUnauthorized
	}

	contractInfo.Admin = newAdmin
	k.setContractInfo(ctx, contractAddr, contractInfo)

	k.Logger(ctx).Info("updated contract admin", "address", contractAddr, "new_admin", newAdmin)
	return nil
}

// QuerySmart runs a readonly query against a contract using the VM engine
func (k Keeper) QuerySmart(ctx context.Context, contractAddr string, queryMsg json.RawMessage) ([]byte, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	contractInfo, found := k.GetContractInfo(ctx, contractAddr)
	if !found {
		return nil, types.ErrContractNotFound
	}

	// Load contract code
	code := k.getCodeBytecode(ctx, contractInfo.CodeID)
	if code == nil {
		return nil, fmt.Errorf("contract code not found for code_id %d", contractInfo.CodeID)
	}

	// Parse the contract definition
	def, parseErr := vm.ParseContractDefinition(code)
	if parseErr != nil {
		// Fall back to legacy query behavior
		return k.querySmartLegacy(ctx, contractAddr, queryMsg)
	}

	// Determine which query handler to call
	handlerName, handlerMsg, err := vm.DetermineHandler(queryMsg)
	if err != nil {
		return nil, types.ErrInvalidMsg
	}

	// Load current contract state
	state := k.loadVMState(ctx, contractAddr)

	// Create execution context (read-only)
	params := k.GetParams(ctx)
	blockHeight := strconv.FormatInt(sdkCtx.BlockHeight(), 10)
	execCtx := vm.NewExecutionContext(contractAddr, "", blockHeight, state, handlerMsg, params.MaxContractGas)

	// Run the VM query
	engine := vm.NewEngine()
	if err := engine.ExecuteQuery(def, handlerName, execCtx); err != nil {
		return nil, fmt.Errorf("contract query failed: %w", err)
	}

	if execCtx.Response != nil {
		return execCtx.Response, nil
	}

	return json.RawMessage(`{}`), nil
}

// querySmartLegacy is the old placeholder query path for non-VM contracts
func (k Keeper) querySmartLegacy(ctx context.Context, contractAddr string, queryMsg json.RawMessage) ([]byte, error) {
	var queryMap map[string]json.RawMessage
	if err := json.Unmarshal(queryMsg, &queryMap); err != nil {
		return nil, types.ErrInvalidMsg
	}

	// Sort query keys for deterministic response ordering.
	queryKeys := make([]string, 0, len(queryMap))
	for key := range queryMap {
		queryKeys = append(queryKeys, key)
	}
	sort.Strings(queryKeys)

	result := make(map[string]json.RawMessage)
	for _, key := range queryKeys {
		state := k.GetContractState(ctx, contractAddr, []byte("exec/"+key))
		if state != nil {
			result[key] = state
		}
	}

	resp, _ := json.Marshal(result)
	return resp, nil
}

// QueryRaw returns the raw contract state for a given key
func (k Keeper) QueryRaw(ctx context.Context, contractAddr string, key []byte) []byte {
	_, found := k.GetContractInfo(ctx, contractAddr)
	if !found {
		return nil
	}
	return k.GetContractState(ctx, contractAddr, key)
}

// GetCodeInfo returns the code info for a given code ID
func (k Keeper) GetCodeInfo(ctx context.Context, codeID uint64) (types.CodeInfo, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.CodeInfoKey(codeID))
	if err != nil || bz == nil {
		return types.CodeInfo{}, false
	}
	var codeInfo types.CodeInfo
	if err := json.Unmarshal(bz, &codeInfo); err != nil {
		return types.CodeInfo{}, false
	}
	return codeInfo, true
}

// GetContractInfo returns the contract info for a given address
func (k Keeper) GetContractInfo(ctx context.Context, contractAddr string) (types.ContractInfo, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.ContractInfoKey(contractAddr))
	if err != nil || bz == nil {
		return types.ContractInfo{}, false
	}
	var contractInfo types.ContractInfo
	if err := json.Unmarshal(bz, &contractInfo); err != nil {
		return types.ContractInfo{}, false
	}
	return contractInfo, true
}

// GetContractState returns a contract state value for a given key
func (k Keeper) GetContractState(ctx context.Context, contractAddr string, key []byte) []byte {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.ContractStateKey(contractAddr, key))
	if err != nil || bz == nil {
		return nil
	}
	return bz
}

// SetContractState stores a contract state entry
func (k Keeper) SetContractState(ctx context.Context, contractAddr string, key []byte, value []byte) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Set(types.ContractStateKey(contractAddr, key), value)
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

// Internal helpers

func (k Keeper) setCodeInfo(ctx context.Context, codeID uint64, info types.CodeInfo) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(info)
	kvStore.Set(types.CodeInfoKey(codeID), bz)
}

func (k Keeper) setCodeBytecode(ctx context.Context, codeID uint64, code []byte) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Set(types.CodeBytecodeKey(codeID), code)
}

func (k Keeper) getCodeBytecode(ctx context.Context, codeID uint64) []byte {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.CodeBytecodeKey(codeID))
	if err != nil {
		return nil
	}
	return bz
}

func (k Keeper) setContractInfo(ctx context.Context, contractAddr string, info types.ContractInfo) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(info)
	kvStore.Set(types.ContractInfoKey(contractAddr), bz)
}

func (k Keeper) getNextCodeID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.NextCodeIDKey))
	if err != nil || bz == nil {
		return 1
	}
	return types.BytesToUint64(bz)
}

func (k Keeper) setNextCodeID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Set([]byte(types.NextCodeIDKey), types.Uint64ToBytes(id))
}

func (k Keeper) generateContractAddress(ctx context.Context, codeID uint64, creator, label string) string {
	// Include an instance counter to guarantee uniqueness even with same codeID+creator+label
	instanceID := k.nextInstanceID(ctx)
	data := fmt.Sprintf("%d/%s/%s/%d", codeID, creator, label, instanceID)
	hash := sha256.Sum256([]byte(data))
	addr := sdk.AccAddress(hash[:20])
	return addr.String()
}

func (k Keeper) nextInstanceID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte("instance_counter"))
	var id uint64
	if err == nil && bz != nil {
		id = types.BytesToUint64(bz)
	}
	id++
	kvStore.Set([]byte("instance_counter"), types.Uint64ToBytes(id))
	return id
}

func (k Keeper) isUploadAllowed(access types.AccessConfig, sender string) bool {
	switch access.Permission {
	case types.AccessTypeEverybody:
		return true
	case types.AccessTypeOnlyAddress:
		return access.Address == sender
	case types.AccessTypeNobody:
		return false
	default:
		return false // Fail-closed: unknown permission types are denied
	}
}

func (k Keeper) isInstantiateAllowed(access types.AccessConfig, sender string) bool {
	switch access.Permission {
	case types.AccessTypeEverybody:
		return true
	case types.AccessTypeOnlyAddress:
		return access.Address == sender
	case types.AccessTypeNobody:
		return false
	default:
		return false // Fail-closed: unknown permission types are denied
	}
}

// loadVMState loads the VM state map for a contract from the KV store.
// All VM state is stored as a single JSON blob under the "_vm_state" key.
func (k Keeper) loadVMState(ctx context.Context, contractAddr string) map[string]string {
	bz := k.GetContractState(ctx, contractAddr, []byte("_vm_state"))
	if bz == nil {
		return make(map[string]string)
	}
	var state map[string]string
	if err := json.Unmarshal(bz, &state); err != nil {
		return make(map[string]string)
	}
	return state
}

// saveVMState persists the VM state map to the KV store as a single JSON blob.
// C-05: Returns error instead of silently ignoring marshal failures.
func (k Keeper) saveVMState(ctx context.Context, contractAddr string, state map[string]string) error {
	bz, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal VM state: %w", err)
	}
	k.SetContractState(ctx, contractAddr, []byte("_vm_state"), bz)
	return nil
}

// DeleteContractState removes a contract state entry
func (k Keeper) DeleteContractState(ctx context.Context, contractAddr string, key []byte) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Delete(types.ContractStateKey(contractAddr, key))
}

// executeBankMsgs processes bank send messages queued by the VM.
// I-06: All sends are wrapped in a CacheContext for atomicity — if any send
// fails, all previous sends within this batch are discarded.
func (k Keeper) executeBankMsgs(ctx context.Context, msgs []vm.BankMsg) error {
	if len(msgs) == 0 {
		return nil
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cacheCtx, write := sdkCtx.CacheContext()

	for _, bm := range msgs {
		fromAddr, err := sdk.AccAddressFromBech32(bm.FromAddress)
		if err != nil {
			return fmt.Errorf("invalid from address %s: %w", bm.FromAddress, err)
		}
		toAddr, err := sdk.AccAddressFromBech32(bm.ToAddress)
		if err != nil {
			return fmt.Errorf("invalid to address %s: %w", bm.ToAddress, err)
		}
		amount, ok := math.NewIntFromString(bm.Amount)
		if !ok {
			return fmt.Errorf("invalid amount %s: not a valid integer", bm.Amount)
		}
		coins := sdk.NewCoins(sdk.NewCoin(bm.Denom, amount))
		if err := k.bankKeeper.SendCoins(cacheCtx, fromAddr, toAddr, coins); err != nil {
			return err // all previous sends are discarded
		}
	}

	write() // commit all sends atomically
	return nil
}

// emitVMEvents translates VM events into SDK events.
// I-10: Attribute keys are sorted before iterating to ensure deterministic
// event ordering across all validators, regardless of Go map iteration order.
func (k Keeper) emitVMEvents(ctx context.Context, events []vm.Event) {
	if len(events) == 0 {
		return
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	for _, evt := range events {
		// Sort attribute keys for deterministic iteration
		keys := make([]string, 0, len(evt.Attributes))
		for k := range evt.Attributes {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		attrs := make([]sdk.Attribute, 0, len(evt.Attributes))
		for _, k := range keys {
			v := evt.Attributes[k]
			attrs = append(attrs, sdk.NewAttribute(k, v))
		}
		sdkCtx.EventManager().EmitEvent(sdk.NewEvent("wasm-"+evt.Type, attrs...))
	}
}
