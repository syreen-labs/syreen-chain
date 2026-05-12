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
	From     string `json:"from"`      // cosmos bech32 sender
	To       string `json:"to"`        // hex destination (empty for contract creation)
	Value    string `json:"value"`     // amount in usyreen
	GasLimit uint64 `json:"gas_limit"` // gas limit
	Data     string `json:"data"`      // hex-encoded calldata or contract bytecode
	Nonce    uint64 `json:"nonce"`     // EVM nonce
	RawTxHex string `json:"raw_tx_hex,omitempty"` // RLP-encoded signed Ethereum tx (set by eth_sendRawTransaction)
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
	GasUsed         uint64 `json:"gas_used"`
	VmError         string `json:"vm_error,omitempty"`
	ReturnData      string `json:"return_data,omitempty"`
	ContractAddress string `json:"contract_address,omitempty"`
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
