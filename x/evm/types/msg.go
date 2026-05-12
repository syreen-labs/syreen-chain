package types

import (
	"encoding/hex"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// EmptyEthAddress is the zero Ethereum address
var EmptyEthAddress = common.Address{}

// MsgEthereumTx wraps an Ethereum transaction for Cosmos SDK processing
type MsgEthereumTx struct {
	From     string `protobuf:"bytes,1,opt,name=from,proto3" json:"from"`
	To       string `protobuf:"bytes,2,opt,name=to,proto3" json:"to"`
	Value    string `protobuf:"bytes,3,opt,name=value,proto3" json:"value"`
	GasLimit uint64 `protobuf:"varint,4,opt,name=gas_limit,json=gasLimit,proto3" json:"gas_limit"`
	Data     string `protobuf:"bytes,5,opt,name=data,proto3" json:"data"`
	Nonce    uint64 `protobuf:"varint,6,opt,name=nonce,proto3" json:"nonce"`
	RawTxHex string `protobuf:"bytes,7,opt,name=raw_tx_hex,json=rawTxHex,proto3" json:"raw_tx_hex,omitempty"`
}

// MaxTxGasLimit is the maximum gas limit for a single EVM transaction (30M)
const MaxTxGasLimit = uint64(30_000_000)

func (msg *MsgEthereumTx) ValidateBasic() error {
	if msg.From == "" {
		return fmt.Errorf("sender cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(msg.From); err != nil {
		return fmt.Errorf("invalid sender address: %w", err)
	}
	if msg.To != "" {
		if !common.IsHexAddress(msg.To) {
			return fmt.Errorf("invalid to address: must be hex")
		}
	}
	if msg.Data != "" {
		if _, err := hex.DecodeString(stripHexPrefix(msg.Data)); err != nil {
			return fmt.Errorf("invalid data: must be hex encoded: %w", err)
		}
		// Limit calldata to 128KB to prevent mempool abuse.
		if len(msg.Data) > 256*1024 { // hex-encoded = 2x raw bytes
			return fmt.Errorf("data exceeds maximum size of 128KB")
		}
	}
	if msg.GasLimit == 0 {
		return fmt.Errorf("gas limit must be > 0")
	}
	// M5: Validate max gas limit
	if msg.GasLimit > MaxTxGasLimit {
		return fmt.Errorf("gas limit %d exceeds maximum %d", msg.GasLimit, MaxTxGasLimit)
	}
	// M4: Validate value is not negative
	if msg.Value != "" {
		v, ok := new(big.Int).SetString(msg.Value, 10)
		if !ok {
			return fmt.Errorf("invalid value: must be a decimal integer")
		}
		if v.Sign() < 0 {
			return fmt.Errorf("value cannot be negative")
		}
	}
	return nil
}

func (msg *MsgEthereumTx) GetValue() *big.Int {
	v, ok := new(big.Int).SetString(msg.Value, 10)
	if !ok {
		return big.NewInt(0)
	}
	return v
}

func (msg *MsgEthereumTx) GetData() []byte {
	data, _ := hex.DecodeString(stripHexPrefix(msg.Data))
	return data
}

func (msg *MsgEthereumTx) GetToAddress() *common.Address {
	if msg.To == "" {
		return nil
	}
	addr := common.HexToAddress(msg.To)
	return &addr
}

// MsgEthereumTxResponse is the response from executing an EVM transaction (H3)
type MsgEthereumTxResponse struct {
	GasUsed         uint64 `protobuf:"varint,1,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used"`
	VmError         string `protobuf:"bytes,2,opt,name=vm_error,json=vmError,proto3" json:"vm_error,omitempty"`
	ReturnData      string `protobuf:"bytes,3,opt,name=return_data,json=returnData,proto3" json:"return_data,omitempty"`
	ContractAddress string `protobuf:"bytes,4,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty"`
}

// DecodeHexData decodes hex-encoded data, stripping 0x prefix if present
func DecodeHexData(s string) ([]byte, error) {
	return hex.DecodeString(stripHexPrefix(s))
}

func stripHexPrefix(s string) string {
	if len(s) >= 2 && s[:2] == "0x" {
		return s[2:]
	}
	return s
}
