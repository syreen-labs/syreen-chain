package jsonrpc

import (
	"fmt"

	"cosmossdk.io/log"
)

// StartJSONRPCServer creates and starts the JSON-RPC server as a goroutine.
// It connects to the CometBFT RPC to query blocks/txs, and reads chain state
// directly from the app's committed multistore (bypassing the IAVL version bug).
func StartJSONRPCServer(app AppStateProvider, listenAddr string, cometRPCAddr string, logger log.Logger) (*Server, error) {
	rpcLogger := logger.With("module", "json-rpc")

	backend, err := NewBackend(app, cometRPCAddr, rpcLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON-RPC backend: %w", err)
	}

	srv := NewServer(backend, listenAddr, rpcLogger)

	go func() {
		if err := srv.Start(); err != nil {
			rpcLogger.Error("JSON-RPC server stopped", "error", err)
		}
	}()

	rpcLogger.Info("JSON-RPC server started",
		"listen", listenAddr,
		"comet_rpc", cometRPCAddr,
		"chain_id", backend.ChainID().Int64(),
	)

	return srv, nil
}
