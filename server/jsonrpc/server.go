package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"cosmossdk.io/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// Server is the Ethereum JSON-RPC server
type Server struct {
	backend    *Backend
	httpServer *http.Server
	logger     log.Logger
}

// NewServer creates a new JSON-RPC server
func NewServer(backend *Backend, listenAddr string, logger log.Logger) *Server {
	s := &Server{
		backend: backend,
		logger:  logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRequest)

	s.httpServer = &http.Server{
		Addr:         listenAddr,
		Handler:      corsMiddleware(mux),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return s
}

// Start begins serving JSON-RPC requests
func (s *Server) Start() error {
	s.logger.Info("Starting JSON-RPC server", "address", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

// I-08: Restrict CORS to known origins instead of allowing all (*).
func corsMiddleware(next http.Handler) http.Handler {
	allowedOrigins := map[string]bool{
		"https://syreen.net":              true,
		"https://www.syreen.net":          true,
		"https://explorer.syreenlabs.com": true,
		"https://wallet.syreenlabs.com":   true,
		"https://rpc.syreenlabs.com":      true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "https://syreen.net")
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Vary", "Origin")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// Some clients send GET for health checks
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"jsonrpc": "2.0",
			"status":  "Syreen Chain JSON-RPC",
		})
		return
	}

	// Limit request body to 1MB to prevent OOM from oversized payloads.
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, nil, -32700, "Parse error")
		return
	}

	result, rpcErr := s.dispatch(req.Method, req.Params)

	w.Header().Set("Content-Type", "application/json")
	if rpcErr != nil {
		json.NewEncoder(w).Encode(JSONRPCResponse{
			JSONRPC: "2.0",
			Error:   rpcErr,
			ID:      req.ID,
		})
		return
	}

	json.NewEncoder(w).Encode(JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      req.ID,
	})
}

func (s *Server) dispatch(method string, params json.RawMessage) (interface{}, *JSONRPCError) {
	s.logger.Debug("JSON-RPC call", "method", method)

	switch method {
	// --- Chain Identity ---
	case "eth_chainId":
		return hexutil.EncodeBig(s.backend.ChainID()), nil

	case "net_version":
		return fmt.Sprintf("%d", s.backend.ChainID().Int64()), nil

	case "web3_clientVersion":
		return "Syreen/v1.0.0/linux-amd64", nil

	case "web3_sha3":
		return s.web3Sha3(params)

	case "net_listening":
		return true, nil

	case "net_peerCount":
		count, err := s.backend.PeerCount()
		if err != nil {
			return nil, internalError(err)
		}
		return HexInt64(int64(count)), nil

	// --- Block Queries ---
	case "eth_blockNumber":
		height, err := s.backend.BlockNumber()
		if err != nil {
			return nil, internalError(err)
		}
		return HexInt64(height), nil

	case "eth_getBlockByNumber":
		return s.ethGetBlockByNumber(params)

	case "eth_getBlockByHash":
		return s.ethGetBlockByHash(params)

	// --- Account State ---
	case "eth_getBalance":
		return s.ethGetBalance(params)

	case "eth_getTransactionCount":
		return s.ethGetTransactionCount(params)

	case "eth_getCode":
		return s.ethGetCode(params)

	case "eth_getStorageAt":
		return s.ethGetStorageAt(params)

	// --- EVM Execution ---
	case "eth_call":
		return s.ethCall(params)

	case "eth_estimateGas":
		return s.ethEstimateGas(params)

	case "eth_gasPrice":
		return hexutil.EncodeBig(s.backend.GasPrice()), nil

	case "eth_maxPriorityFeePerGas":
		return "0x1", nil

	case "eth_feeHistory":
		return s.ethFeeHistory(params)

	// --- Transaction Queries ---
	case "eth_getTransactionByHash":
		return s.ethGetTransactionByHash(params)

	case "eth_getTransactionReceipt":
		return s.ethGetTransactionReceipt(params)

	case "eth_sendRawTransaction":
		return s.ethSendRawTransaction(params)

	// --- Logs & Filters ---
	case "eth_getLogs":
		return []interface{}{}, nil

	case "eth_newFilter":
		return "0x1", nil

	case "eth_newBlockFilter":
		return "0x1", nil

	case "eth_newPendingTransactionFilter":
		return "0x1", nil

	case "eth_getFilterChanges":
		return []interface{}{}, nil

	case "eth_getFilterLogs":
		return []interface{}{}, nil

	case "eth_uninstallFilter":
		return true, nil

	// --- Misc ---
	case "eth_accounts":
		return []string{}, nil

	case "eth_mining":
		return false, nil

	case "eth_hashrate":
		return "0x0", nil

	case "eth_syncing":
		result, err := s.backend.Syncing()
		if err != nil {
			return nil, internalError(err)
		}
		return result, nil

	case "eth_coinbase":
		return EmptyAddress, nil

	case "eth_getUncleCountByBlockHash", "eth_getUncleCountByBlockNumber":
		return "0x0", nil

	case "eth_getUncleByBlockHashAndIndex", "eth_getUncleByBlockNumberAndIndex":
		return nil, nil

	case "eth_protocolVersion":
		return "0x41", nil

	case "eth_getBlockTransactionCountByHash", "eth_getBlockTransactionCountByNumber":
		return "0x0", nil

	case "eth_sign", "eth_signTransaction":
		return nil, &JSONRPCError{Code: -32601, Message: "not supported: use eth_sendRawTransaction"}

	default:
		return nil, &JSONRPCError{Code: -32601, Message: fmt.Sprintf("method %s not found", method)}
	}
}

// --- Method Implementations ---

func (s *Server) web3Sha3(params json.RawMessage) (interface{}, *JSONRPCError) {
	var args []string
	if err := json.Unmarshal(params, &args); err != nil || len(args) < 1 {
		return nil, invalidParams("expected [data]")
	}
	data, err := hexutil.Decode(args[0])
	if err != nil {
		return nil, invalidParams("invalid hex data")
	}
	return hexutil.Encode(crypto.Keccak256(data)), nil
}

func (s *Server) ethGetBalance(params json.RawMessage) (interface{}, *JSONRPCError) {
	var args []json.RawMessage
	if err := json.Unmarshal(params, &args); err != nil || len(args) < 1 {
		return nil, invalidParams("expected [address, block]")
	}

	var addrStr string
	if err := json.Unmarshal(args[0], &addrStr); err != nil {
		return nil, invalidParams("invalid address")
	}

	if !common.IsHexAddress(addrStr) {
		return nil, invalidParams("invalid hex address")
	}

	addr := common.HexToAddress(addrStr)
	balance, err := s.backend.GetBalance(addr)
	if err != nil {
		return nil, internalError(err)
	}

	return hexutil.EncodeBig(balance), nil
}

func (s *Server) ethGetTransactionCount(params json.RawMessage) (interface{}, *JSONRPCError) {
	var args []json.RawMessage
	if err := json.Unmarshal(params, &args); err != nil || len(args) < 1 {
		return nil, invalidParams("expected [address, block]")
	}

	var addrStr string
	if err := json.Unmarshal(args[0], &addrStr); err != nil {
		return nil, invalidParams("invalid address")
	}

	addr := common.HexToAddress(addrStr)
	nonce, err := s.backend.GetNonce(addr)
	if err != nil {
		return nil, internalError(err)
	}

	return HexUint64(nonce), nil
}

func (s *Server) ethGetCode(params json.RawMessage) (interface{}, *JSONRPCError) {
	var args []json.RawMessage
	if err := json.Unmarshal(params, &args); err != nil || len(args) < 1 {
		return nil, invalidParams("expected [address, block]")
	}

	var addrStr string
	if err := json.Unmarshal(args[0], &addrStr); err != nil {
		return nil, invalidParams("invalid address")
	}

	addr := common.HexToAddress(addrStr)
	code, err := s.backend.GetCode(addr)
	if err != nil {
		return nil, internalError(err)
	}

	if len(code) == 0 {
		return "0x", nil
	}
	return hexutil.Encode(code), nil
}

func (s *Server) ethGetStorageAt(params json.RawMessage) (interface{}, *JSONRPCError) {
	var args []json.RawMessage
	if err := json.Unmarshal(params, &args); err != nil || len(args) < 2 {
		return nil, invalidParams("expected [address, slot, block]")
	}

	var addrStr, slotStr string
	json.Unmarshal(args[0], &addrStr)
	json.Unmarshal(args[1], &slotStr)

	addr := common.HexToAddress(addrStr)
	slot := common.HexToHash(slotStr)
	value, err := s.backend.GetStorageAt(addr, slot)
	if err != nil {
		return nil, internalError(err)
	}

	return value.Hex(), nil
}

func (s *Server) ethCall(params json.RawMessage) (interface{}, *JSONRPCError) {
	var args []json.RawMessage
	if err := json.Unmarshal(params, &args); err != nil || len(args) < 1 {
		return nil, invalidParams("expected [callArgs, block]")
	}

	var callArgs CallArgs
	if err := json.Unmarshal(args[0], &callArgs); err != nil {
		return nil, invalidParams("invalid call args")
	}

	var gas uint64
	if callArgs.Gas != nil {
		gas = uint64(*callArgs.Gas)
	}

	result, err := s.backend.Call(callArgs.From, callArgs.To, callArgs.GetData(), gas)
	if err != nil {
		return nil, &JSONRPCError{Code: 3, Message: err.Error()}
	}

	return hexutil.Encode(result), nil
}

func (s *Server) ethEstimateGas(params json.RawMessage) (interface{}, *JSONRPCError) {
	var args []json.RawMessage
	if err := json.Unmarshal(params, &args); err != nil || len(args) < 1 {
		return nil, invalidParams("expected [callArgs]")
	}

	var callArgs CallArgs
	if err := json.Unmarshal(args[0], &callArgs); err != nil {
		return nil, invalidParams("invalid call args")
	}

	var value *big.Int
	if callArgs.Value != nil {
		value = callArgs.Value.ToInt()
	}

	estimate, err := s.backend.EstimateGas(callArgs.From, callArgs.To, callArgs.GetData(), value)
	if err != nil {
		return nil, &JSONRPCError{Code: 3, Message: err.Error()}
	}

	return HexUint64(estimate), nil
}

func (s *Server) ethGetBlockByNumber(params json.RawMessage) (interface{}, *JSONRPCError) {
	var args []json.RawMessage
	if err := json.Unmarshal(params, &args); err != nil || len(args) < 2 {
		return nil, invalidParams("expected [blockNumber, fullTx]")
	}

	var blockStr string
	var fullTx bool
	json.Unmarshal(args[0], &blockStr)
	json.Unmarshal(args[1], &fullTx)

	blockNum, err := s.backend.ParseBlockNumber(blockStr)
	if err != nil {
		return nil, invalidParams(err.Error())
	}

	block, err := s.backend.GetBlockByNumber(blockNum, fullTx)
	if err != nil {
		return nil, internalError(err)
	}

	return block, nil
}

func (s *Server) ethGetBlockByHash(params json.RawMessage) (interface{}, *JSONRPCError) {
	var args []json.RawMessage
	if err := json.Unmarshal(params, &args); err != nil || len(args) < 2 {
		return nil, invalidParams("expected [blockHash, fullTx]")
	}

	var hashStr string
	var fullTx bool
	json.Unmarshal(args[0], &hashStr)
	json.Unmarshal(args[1], &fullTx)

	hash := common.HexToHash(hashStr)
	block, err := s.backend.GetBlockByHash(hash, fullTx)
	if err != nil {
		return nil, internalError(err)
	}

	return block, nil
}

func (s *Server) ethGetTransactionByHash(params json.RawMessage) (interface{}, *JSONRPCError) {
	var args []string
	if err := json.Unmarshal(params, &args); err != nil || len(args) < 1 {
		return nil, invalidParams("expected [txHash]")
	}

	tx, err := s.backend.GetTransactionByHash(args[0])
	if err != nil {
		return nil, internalError(err)
	}

	return tx, nil
}

func (s *Server) ethGetTransactionReceipt(params json.RawMessage) (interface{}, *JSONRPCError) {
	var args []string
	if err := json.Unmarshal(params, &args); err != nil || len(args) < 1 {
		return nil, invalidParams("expected [txHash]")
	}

	receipt, err := s.backend.GetTransactionReceipt(args[0])
	if err != nil {
		return nil, internalError(err)
	}

	return receipt, nil
}

func (s *Server) ethSendRawTransaction(params json.RawMessage) (interface{}, *JSONRPCError) {
	var args []string
	if err := json.Unmarshal(params, &args); err != nil || len(args) < 1 {
		return nil, invalidParams("expected [rawTxHex]")
	}

	rawTxHex := args[0]
	if rawTxHex == "" {
		return nil, invalidParams("raw transaction data cannot be empty")
	}

	txHash, err := s.backend.BroadcastRawEthTx(rawTxHex)
	if err != nil {
		s.logger.Error("eth_sendRawTransaction failed", "error", err)
		return nil, &JSONRPCError{Code: -32000, Message: err.Error()}
	}

	return txHash, nil
}

func (s *Server) ethFeeHistory(params json.RawMessage) (interface{}, *JSONRPCError) {
	height, err := s.backend.BlockNumber()
	if err != nil {
		return nil, internalError(err)
	}

	// Return a simple fee history
	return map[string]interface{}{
		"oldestBlock":   HexInt64(height),
		"baseFeePerGas": []string{"0x1", "0x1"},
		"gasUsedRatio":  []float64{0.0},
		"reward":        [][]string{{"0x1"}},
	}, nil
}

// --- Error helpers ---

func writeError(w http.ResponseWriter, id json.RawMessage, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JSONRPCResponse{
		JSONRPC: "2.0",
		Error:   &JSONRPCError{Code: code, Message: msg},
		ID:      id,
	})
}

func invalidParams(msg string) *JSONRPCError {
	return &JSONRPCError{Code: -32602, Message: msg}
}

func internalError(err error) *JSONRPCError {
	// Do not leak internal error details to clients.
	_ = err // logged by caller if needed
	return &JSONRPCError{Code: -32603, Message: "internal server error"}
}

// Ensure strings import is used
var _ = strings.TrimPrefix
