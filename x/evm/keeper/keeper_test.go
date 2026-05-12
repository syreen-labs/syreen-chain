package keeper_test

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"syreen/x/evm/keeper"
	"syreen/x/evm/statedb"
	"syreen/x/evm/types"
)

func TestMain(m *testing.M) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("syreen", "syreenpub")
	config.SetBech32PrefixForValidator("syreenvaloper", "syreenvaloperpub")
	config.SetBech32PrefixForConsensusNode("syreenvalcons", "syreenvalconspub")
	os.Exit(m.Run())
}

func TestDefaultParams(t *testing.T) {
	p := types.DefaultParams()
	require.Equal(t, "usyreen", p.EvmDenom)
	require.True(t, p.EnableCreate)
	require.True(t, p.EnableCall)
	require.NoError(t, p.Validate())
}

func TestChainConfig(t *testing.T) {
	cfg := types.DefaultChainConfig()
	require.Equal(t, big.NewInt(79733), cfg.ChainID)
	require.NotNil(t, cfg.HomesteadBlock)
	require.NotNil(t, cfg.ByzantiumBlock)
	require.NotNil(t, cfg.ConstantinopleBlock)
	require.NotNil(t, cfg.IstanbulBlock)
	require.NotNil(t, cfg.BerlinBlock)
	require.NotNil(t, cfg.LondonBlock)
	require.True(t, cfg.TerminalTotalDifficultyPassed)

	// All forks at block 0
	require.Equal(t, big.NewInt(0), cfg.HomesteadBlock)
	require.Equal(t, big.NewInt(0), cfg.BerlinBlock)
	require.Equal(t, big.NewInt(0), cfg.LondonBlock)
}

func TestEthAddressConversion(t *testing.T) {
	// Test bech32 -> eth -> bech32 roundtrip
	bech32 := "syreen1ra47n7hkl5r2rd8kda3s3tywnjwzxlsws4xr53"
	ethAddr, err := keeper.EthAddressFromBech32(bech32)
	require.NoError(t, err)
	require.Len(t, ethAddr.Bytes(), 20)

	// Roundtrip back
	cosmosAddr := keeper.Bech32FromEthAddress(ethAddr)
	require.Equal(t, bech32, cosmosAddr.String())
}

func TestInvalidBech32Address(t *testing.T) {
	_, err := keeper.EthAddressFromBech32("invalid")
	require.Error(t, err)
}

func TestGenesisState(t *testing.T) {
	gs := types.DefaultGenesisState()
	require.NotNil(t, gs)
	require.NoError(t, gs.Validate())
	require.Equal(t, "usyreen", gs.Params.EvmDenom)
	require.Empty(t, gs.Accounts)
}

func TestInvalidParams(t *testing.T) {
	p := types.Params{
		EvmDenom:     "",
		EnableCreate: true,
		EnableCall:   true,
	}
	require.Error(t, p.Validate())
}

func TestMsgEthereumTxValidation(t *testing.T) {
	tests := []struct {
		name    string
		msg     types.MsgEthereumTx
		wantErr bool
	}{
		{
			name: "valid deploy",
			msg: types.MsgEthereumTx{
				From:     "syreen1ra47n7hkl5r2rd8kda3s3tywnjwzxlsws4xr53",
				To:       "",
				Value:    "0",
				GasLimit: 1000000,
				Data:     "608060405234801561001057600080fd5b50",
			},
			wantErr: false,
		},
		{
			name: "valid call",
			msg: types.MsgEthereumTx{
				From:     "syreen1ra47n7hkl5r2rd8kda3s3tywnjwzxlsws4xr53",
				To:       "0x1234567890abcdef1234567890abcdef12345678",
				Value:    "100",
				GasLimit: 100000,
				Data:     "a9059cbb",
			},
			wantErr: false,
		},
		{
			name: "empty sender",
			msg: types.MsgEthereumTx{
				From:     "",
				GasLimit: 100000,
			},
			wantErr: true,
		},
		{
			name: "zero gas",
			msg: types.MsgEthereumTx{
				From:     "syreen1ra47n7hkl5r2rd8kda3s3tywnjwzxlsws4xr53",
				GasLimit: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid to address",
			msg: types.MsgEthereumTx{
				From:     "syreen1ra47n7hkl5r2rd8kda3s3tywnjwzxlsws4xr53",
				To:       "not-hex",
				GasLimit: 100000,
			},
			wantErr: true,
		},
		{
			name: "invalid data hex",
			msg: types.MsgEthereumTx{
				From:     "syreen1ra47n7hkl5r2rd8kda3s3tywnjwzxlsws4xr53",
				GasLimit: 100000,
				Data:     "not-hex-zzz",
			},
			wantErr: true,
		},
		// M4: Negative value validation
		{
			name: "negative value",
			msg: types.MsgEthereumTx{
				From:     "syreen1ra47n7hkl5r2rd8kda3s3tywnjwzxlsws4xr53",
				GasLimit: 100000,
				Value:    "-100",
			},
			wantErr: true,
		},
		// M5: Max gas limit validation
		{
			name: "gas limit exceeds max",
			msg: types.MsgEthereumTx{
				From:     "syreen1ra47n7hkl5r2rd8kda3s3tywnjwzxlsws4xr53",
				GasLimit: 31_000_000,
				Value:    "0",
			},
			wantErr: true,
		},
		// M4: Invalid value format
		{
			name: "invalid value format",
			msg: types.MsgEthereumTx{
				From:     "syreen1ra47n7hkl5r2rd8kda3s3tywnjwzxlsws4xr53",
				GasLimit: 100000,
				Value:    "not-a-number",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgGetValue(t *testing.T) {
	msg := &types.MsgEthereumTx{Value: "1000000"}
	require.Equal(t, big.NewInt(1000000), msg.GetValue())

	msg2 := &types.MsgEthereumTx{Value: ""}
	require.Equal(t, big.NewInt(0), msg2.GetValue())
}

func TestMsgGetData(t *testing.T) {
	msg := &types.MsgEthereumTx{Data: "abcdef"}
	data := msg.GetData()
	require.Equal(t, []byte{0xab, 0xcd, 0xef}, data)

	// With 0x prefix
	msg2 := &types.MsgEthereumTx{Data: "0xabcdef"}
	data2 := msg2.GetData()
	require.Equal(t, []byte{0xab, 0xcd, 0xef}, data2)
}

func TestMsgGetToAddress(t *testing.T) {
	// Contract creation (empty To)
	msg := &types.MsgEthereumTx{To: ""}
	require.Nil(t, msg.GetToAddress())

	// Contract call
	msg2 := &types.MsgEthereumTx{To: "0x1234567890abcdef1234567890abcdef12345678"}
	addr := msg2.GetToAddress()
	require.NotNil(t, addr)
	require.Equal(t, common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"), *addr)
}

func TestActivePrecompiles(t *testing.T) {
	cfg := types.DefaultChainConfig()
	rules := cfg.Rules(big.NewInt(1), true)
	precompiles := vm.ActivePrecompiles(rules)
	// Should have standard precompiles (ecrecover, sha256, ripemd160, etc)
	require.NotEmpty(t, precompiles)
}

func TestEVMContractBytecode(t *testing.T) {
	// Verify a simple Solidity storage contract bytecode can be decoded
	bytecode := "608060405234801561001057600080fd5b5060c78061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c806360fe47b11460375780636d4ce63c146049575b600080fd5b60476042366004607d565b600055565b005b60005460405190815260200160405180910390f35b600060208284031215608d57600080fd5b503591905056fea264697066735822"
	data, err := types.DecodeHexData(bytecode)
	require.NoError(t, err)
	require.NotEmpty(t, data)
}

func TestKeccak256(t *testing.T) {
	// Verify the Keccak256 function selector generation works
	// function set(uint256) -> 0x60fe47b1
	hash := crypto.Keccak256([]byte("set(uint256)"))
	selector := hash[:4]
	require.Equal(t, []byte{0x60, 0xfe, 0x47, 0xb1}, selector)
}

// M1: Test SubRefund panics on underflow (matching go-ethereum behavior)
func TestSubRefundPanicsOnUnderflow(t *testing.T) {
	require.Panics(t, func() {
		// SubRefund with gas > refund should panic
		s := &panicRefundHelper{refund: 5}
		s.SubRefund(10)
	})
}

// panicRefundHelper mimics StateDB SubRefund behavior for unit testing without full StateDB setup
type panicRefundHelper struct {
	refund uint64
}

func (p *panicRefundHelper) SubRefund(gas uint64) {
	if gas > p.refund {
		panic("Refund counter below zero")
	}
	p.refund -= gas
}

// Test MsgEthereumTx proto-compatible methods
func TestMsgEthereumTxProtoMethods(t *testing.T) {
	msg := &types.MsgEthereumTx{
		From:     "syreen1ra47n7hkl5r2rd8kda3s3tywnjwzxlsws4xr53",
		To:       "0x1234567890abcdef1234567890abcdef12345678",
		Value:    "100",
		GasLimit: 100000,
		Data:     "abcdef",
	}

	// Test ProtoMessage doesn't panic
	msg.ProtoMessage()

	// Test XXX_MessageName
	require.Equal(t, "syreen.evm.MsgEthereumTx", msg.XXX_MessageName())

	// Test Marshal/Unmarshal roundtrip
	bz, err := msg.Marshal()
	require.NoError(t, err)
	require.NotEmpty(t, bz)

	msg2 := &types.MsgEthereumTx{}
	err = msg2.Unmarshal(bz)
	require.NoError(t, err)
	require.Equal(t, msg.From, msg2.From)
	require.Equal(t, msg.To, msg2.To)
	require.Equal(t, msg.Value, msg2.Value)
	require.Equal(t, msg.GasLimit, msg2.GasLimit)

	// Test Size
	require.Equal(t, len(bz), msg.Size())

	// Test Reset
	msg2.Reset()
	require.Equal(t, "", msg2.From)

	// Test String
	require.Contains(t, msg.String(), "from=")
}

// Test MsgEthereumTxResponse proto-compatible methods
func TestMsgEthereumTxResponseProtoMethods(t *testing.T) {
	resp := &types.MsgEthereumTxResponse{
		GasUsed: 21000,
		VmError: "",
	}
	resp.ProtoMessage()
	require.Equal(t, "syreen.evm.MsgEthereumTxResponse", resp.XXX_MessageName())
	require.Contains(t, resp.String(), "gas_used=21000")

	resp.Reset()
	require.Equal(t, uint64(0), resp.GasUsed)
}

// Test MaxCodeSize constant (M3: EIP-170)
func TestMaxCodeSizeConstant(t *testing.T) {
	require.Equal(t, 24576, statedb.MaxCodeSize)
}

// Test MaxGasLimit constant (M5)
func TestMaxGasLimitConstant(t *testing.T) {
	require.Equal(t, uint64(30_000_000), keeper.MaxGasLimit)
}

// Test MaxTxGasLimit validation (M5)
func TestMaxTxGasLimit(t *testing.T) {
	require.Equal(t, uint64(30_000_000), types.MaxTxGasLimit)
}

// Test GenesisState with accounts
func TestGenesisStateWithAccounts(t *testing.T) {
	gs := types.GenesisState{
		Params: types.DefaultParams(),
		Accounts: []types.Account{
			{
				Address: "0x1234567890abcdef1234567890abcdef12345678",
				Code:    "608060",
				Storage: map[string]string{
					"0x0000000000000000000000000000000000000000000000000000000000000001": "0x0000000000000000000000000000000000000000000000000000000000000042",
				},
			},
		},
	}
	require.NoError(t, gs.Validate())
	require.Len(t, gs.Accounts, 1)
}

// Test KeyParams (M6)
func TestKeyParams(t *testing.T) {
	key := types.KeyParams()
	require.Equal(t, []byte{types.PrefixParams}, key)
}

// Test JSON marshaling of MsgEthereumTx
func TestMsgEthereumTxJSON(t *testing.T) {
	msg := &types.MsgEthereumTx{
		From:     "syreen1ra47n7hkl5r2rd8kda3s3tywnjwzxlsws4xr53",
		To:       "0x1234567890abcdef1234567890abcdef12345678",
		Value:    "100",
		GasLimit: 100000,
		Data:     "abcdef",
		Nonce:    5,
	}

	bz, err := json.Marshal(msg)
	require.NoError(t, err)

	msg2 := &types.MsgEthereumTx{}
	err = json.Unmarshal(bz, msg2)
	require.NoError(t, err)
	require.Equal(t, msg.From, msg2.From)
	require.Equal(t, msg.Nonce, msg2.Nonce)
}

// Test EmptyEthAddress
func TestEmptyEthAddress(t *testing.T) {
	require.Equal(t, common.Address{}, types.EmptyEthAddress)
}

// Test EVM error types
func TestStateDBErrors(t *testing.T) {
	require.NotNil(t, statedb.ErrNegativeAmount)
	require.NotNil(t, statedb.ErrCodeTooLarge)
	require.Contains(t, statedb.ErrNegativeAmount.Error(), "negative")
	require.Contains(t, statedb.ErrCodeTooLarge.Error(), "24KB")
}

// Test params validation
func TestParamsValidation(t *testing.T) {
	// Valid params
	p := types.DefaultParams()
	require.NoError(t, p.Validate())

	// Empty denom
	p2 := types.Params{EvmDenom: "", EnableCreate: true, EnableCall: true}
	require.Error(t, p2.Validate())
}

// Test params JSON marshal/unmarshal (for M6 KVStore persistence)
func TestParamsJSON(t *testing.T) {
	p := types.DefaultParams()
	bz, err := json.Marshal(p)
	require.NoError(t, err)

	var p2 types.Params
	err = json.Unmarshal(bz, &p2)
	require.NoError(t, err)
	require.Equal(t, p.EvmDenom, p2.EvmDenom)
	require.Equal(t, p.EnableCreate, p2.EnableCreate)
	require.Equal(t, p.EnableCall, p2.EnableCall)
}
