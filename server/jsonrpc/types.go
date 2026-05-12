package jsonrpc

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  json.RawMessage   `json:"params"`
	ID      json.RawMessage   `json:"id"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
	ID      json.RawMessage `json:"id"`
}

// JSONRPCError represents a JSON-RPC 2.0 error
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// CallArgs represents the arguments to eth_call and eth_estimateGas
type CallArgs struct {
	From     *common.Address `json:"from"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Data     *hexutil.Bytes  `json:"data"`
	Input    *hexutil.Bytes  `json:"input"`
}

// GetData returns the input data, preferring "input" over "data"
func (args *CallArgs) GetData() []byte {
	if args.Input != nil {
		return *args.Input
	}
	if args.Data != nil {
		return *args.Data
	}
	return nil
}

// RPCBlock represents an Ethereum block for JSON-RPC responses
type RPCBlock struct {
	Number           string        `json:"number"`
	Hash             string        `json:"hash"`
	ParentHash       string        `json:"parentHash"`
	Nonce            string        `json:"nonce"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	LogsBloom        string        `json:"logsBloom"`
	TransactionsRoot string        `json:"transactionsRoot"`
	StateRoot        string        `json:"stateRoot"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	Miner            string        `json:"miner"`
	Difficulty       string        `json:"difficulty"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	ExtraData        string        `json:"extraData"`
	Size             string        `json:"size"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Timestamp        string        `json:"timestamp"`
	Transactions     []interface{} `json:"transactions"`
	Uncles           []string      `json:"uncles"`
	MixHash          string        `json:"mixHash"`
	BaseFeePerGas    string        `json:"baseFeePerGas"`
}

// RPCTransaction represents an Ethereum transaction for JSON-RPC responses
type RPCTransaction struct {
	BlockHash        *string `json:"blockHash"`
	BlockNumber      *string `json:"blockNumber"`
	From             string  `json:"from"`
	Gas              string  `json:"gas"`
	GasPrice         string  `json:"gasPrice"`
	Hash             string  `json:"hash"`
	Input            string  `json:"input"`
	Nonce            string  `json:"nonce"`
	To               *string `json:"to"`
	TransactionIndex *string `json:"transactionIndex"`
	Value            string  `json:"value"`
	V                string  `json:"v"`
	R                string  `json:"r"`
	S                string  `json:"s"`
	Type             string  `json:"type"`
}

// RPCReceipt represents a transaction receipt for JSON-RPC responses
type RPCReceipt struct {
	TransactionHash   string      `json:"transactionHash"`
	TransactionIndex  string      `json:"transactionIndex"`
	BlockHash         string      `json:"blockHash"`
	BlockNumber       string      `json:"blockNumber"`
	From              string      `json:"from"`
	To                *string     `json:"to"`
	CumulativeGasUsed string      `json:"cumulativeGasUsed"`
	GasUsed           string      `json:"gasUsed"`
	ContractAddress   *string     `json:"contractAddress"`
	Logs              []RPCLog    `json:"logs"`
	LogsBloom         string      `json:"logsBloom"`
	Status            string      `json:"status"`
	EffectiveGasPrice string      `json:"effectiveGasPrice"`
	Type              string      `json:"type"`
}

// RPCLog represents an EVM log entry
type RPCLog struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	BlockHash        string   `json:"blockHash"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
}

// HexUint64 formats a uint64 as 0x-prefixed hex
func HexUint64(n uint64) string {
	return fmt.Sprintf("0x%x", n)
}

// HexInt64 formats an int64 as 0x-prefixed hex
func HexInt64(n int64) string {
	return fmt.Sprintf("0x%x", n)
}

// HexBig formats a big.Int as 0x-prefixed hex
func HexBig(n *big.Int) string {
	if n == nil {
		return "0x0"
	}
	return fmt.Sprintf("0x%x", n)
}

// EmptyHash is 32 zero bytes as hex
var EmptyHash = "0x0000000000000000000000000000000000000000000000000000000000000000"

// EmptyBloom is 256 zero bytes as hex
var EmptyBloom = "0x" + fmt.Sprintf("%0512x", 0)

// EmptyAddress is 20 zero bytes as hex
var EmptyAddress = "0x0000000000000000000000000000000000000000"
