package evm

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"syreen/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"syreen/x/evm/client/cli"
	"syreen/x/evm/keeper"
	"syreen/x/evm/types"
)

var (
	_ module.HasGenesis  = AppModule{}
	_ module.AppModule   = AppModule{}
	_ module.HasServices = AppModule{}
)

// AppModule implements the AppModule interface
type AppModule struct {
	keeper keeper.Keeper
}

func NewAppModule(k keeper.Keeper) AppModule {
	return AppModule{
		keeper: k,
	}
}

func (AppModule) Name() string { return types.ModuleName }

func (AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

func (AppModule) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (AppModule) DefaultGenesis(_ codec.JSONCodec) json.RawMessage {
	gs := types.DefaultGenesisState()
	bz, _ := json.Marshal(gs)
	return bz
}

func (AppModule) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var gs types.GenesisState
	if err := json.Unmarshal(bz, &gs); err != nil {
		panic(fmt.Errorf("failed to unmarshal evm genesis state: %w", err))
	}
	return gs.Validate()
}

func (AppModule) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {}

func (AppModule) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

func (AppModule) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// H3: RegisterServices registers the MsgServer so MsgEthereumTx can be delivered
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(am.keeper))
}

// M8: InitGenesis imports genesis accounts (code + storage) into the EVM KVStore
func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONCodec, data json.RawMessage) {
	var gs types.GenesisState
	if err := json.Unmarshal(data, &gs); err != nil {
		ctx.Logger().Error(
			"FATAL: malformed evm genesis JSON — halting startup",
			"module", types.ModuleName,
			"error", err,
		)
		panic(fmt.Errorf("failed to unmarshal %s genesis: %w", types.ModuleName, err))
	}

	// Persist params to KVStore (M6)
	if err := am.keeper.SetParams(ctx, gs.Params); err != nil {
		ctx.Logger().Error(
			"FATAL: failed to persist evm params during genesis init",
			"module", types.ModuleName,
			"error", err,
		)
		panic(fmt.Errorf("failed to set %s params: %w", types.ModuleName, err))
	}

	// M8: Import genesis accounts
	kvStore := am.keeper.GetStoreService().OpenKVStore(ctx)
	for _, acct := range gs.Accounts {
		if acct.Address == "" {
			continue
		}
		addr := common.HexToAddress(acct.Address)

		// Mark account as existing
		_ = kvStore.Set(types.KeyAccount(addr.Bytes()), []byte{0x01})

		// Import code
		if acct.Code != "" {
			code, err := types.DecodeHexData(acct.Code)
			if err != nil {
				ctx.Logger().Error(
					"FATAL: invalid hex bytecode in evm genesis account",
					"module", types.ModuleName,
					"address", acct.Address,
					"error", err,
				)
				panic(fmt.Errorf("invalid code hex for %s: %w", acct.Address, err))
			}
			if len(code) > 0 {
				_ = kvStore.Set(types.KeyCode(addr.Bytes()), code)
				codeHash := crypto.Keccak256Hash(code)
				_ = kvStore.Set(types.KeyCodeHash(addr.Bytes()), codeHash.Bytes())
			}
		}

		// Import storage — sort keys for deterministic state writes
		storageKeys := make([]string, 0, len(acct.Storage))
		for hexKey := range acct.Storage {
			storageKeys = append(storageKeys, hexKey)
		}
		sort.Strings(storageKeys)

		for _, hexKey := range storageKeys {
			hexVal := acct.Storage[hexKey]
			key := common.HexToHash(hexKey)
			val := common.HexToHash(hexVal)
			if val != (common.Hash{}) {
				_ = kvStore.Set(types.KeyStorage(addr.Bytes(), key.Bytes()), val.Bytes())
			}
		}

		// Import nonce (NEW-8 fix)
		if acct.Nonce > 0 {
			_ = kvStore.Set(types.KeyNonce(addr.Bytes()), uint64ToBytes(acct.Nonce))
		}
	}
}

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

// M9: ExportGenesis exports EVM state (params + all accounts with code/storage/nonces)
func (am AppModule) ExportGenesis(ctx sdk.Context, _ codec.JSONCodec) json.RawMessage {
	params := am.keeper.GetParams(ctx)

	kvStore := am.keeper.GetStoreService().OpenKVStore(ctx)
	accountMap := make(map[common.Address]*types.Account)

	// Export accounts with code
	codePrefix := []byte{types.PrefixCode}
	codePrefixEnd := append([]byte{types.PrefixCode}, 0xFF)
	codeIter, err := kvStore.Iterator(codePrefix, codePrefixEnd)
	if err == nil {
		for ; codeIter.Valid(); codeIter.Next() {
			addrBytes := codeIter.Key()[1:] // strip prefix byte
			if len(addrBytes) < 20 {
				continue
			}
			addr := common.BytesToAddress(addrBytes[:20])
			code := codeIter.Value()

			// Gather storage
			storageMap := make(map[string]string)
			storagePrefix := types.KeyStorage(addr.Bytes(), nil)
			storageEnd := make([]byte, len(storagePrefix))
			copy(storageEnd, storagePrefix)
			storageEnd = append(storageEnd, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF)
			storageIter, sErr := kvStore.Iterator(storagePrefix, storageEnd)
			if sErr == nil {
				prefixLen := len(storagePrefix)
				for ; storageIter.Valid(); storageIter.Next() {
					slotBytes := storageIter.Key()[prefixLen:]
					storageMap[common.BytesToHash(slotBytes).Hex()] = common.BytesToHash(storageIter.Value()).Hex()
				}
				storageIter.Close()
			}

			accountMap[addr] = &types.Account{
				Address: addr.Hex(),
				Code:    fmt.Sprintf("%x", code),
				Storage: storageMap,
			}
		}
		codeIter.Close()
	}

	// NEW-8 fix: Also export EOA nonces (accounts with nonces but no code)
	noncePrefix := []byte{types.PrefixNonce}
	noncePrefixEnd := append([]byte{types.PrefixNonce}, 0xFF)
	nonceIter, err := kvStore.Iterator(noncePrefix, noncePrefixEnd)
	if err == nil {
		for ; nonceIter.Valid(); nonceIter.Next() {
			addrBytes := nonceIter.Key()[1:] // strip prefix byte
			if len(addrBytes) < 20 {
				continue
			}
			addr := common.BytesToAddress(addrBytes[:20])
			nonce := bytesToUint64(nonceIter.Value())
			if nonce == 0 {
				continue
			}
			if acct, exists := accountMap[addr]; exists {
				acct.Nonce = nonce
			} else {
				accountMap[addr] = &types.Account{
					Address: addr.Hex(),
					Nonce:   nonce,
				}
			}
		}
		nonceIter.Close()
	}

	// Sort account addresses for deterministic genesis export
	sortedAddrs := make([]common.Address, 0, len(accountMap))
	for addr := range accountMap {
		sortedAddrs = append(sortedAddrs, addr)
	}
	sort.Slice(sortedAddrs, func(i, j int) bool {
		return sortedAddrs[i].Hex() < sortedAddrs[j].Hex()
	})

	accounts := make([]types.Account, 0, len(accountMap))
	for _, addr := range sortedAddrs {
		accounts = append(accounts, *accountMap[addr])
	}

	gs := types.GenesisState{
		Params:   params,
		Accounts: accounts,
	}
	bz, _ := json.Marshal(gs)
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

func (am AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) BeginBlock(goCtx context.Context) error {
	if !config.IsModuleEnabled("evm") {
		return nil
	}
	return nil
}

func (am AppModule) EndBlock(_ context.Context) error {
	if !config.IsModuleEnabled("evm") {
		return nil
	}
	return nil
}

func (AppModule) IsOnePerModuleType() {}
func (AppModule) IsAppModule()        {}
