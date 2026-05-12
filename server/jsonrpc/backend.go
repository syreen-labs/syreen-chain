package jsonrpc

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	rpcclient "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	evmkeeper "syreen/x/evm/keeper"
	evmtypes "syreen/x/evm/types"
)

// AppStateProvider gives us access to the committed multistore and keepers
// without importing the full app package (avoiding circular deps)
type AppStateProvider interface {
	GetCommitMultiStore() storetypes.CommitMultiStore
	GetBankKeeper() bankkeeper.Keeper
	GetEVMKeeper() evmkeeper.Keeper
	GetLogger() log.Logger
	GetTxConfig() client.TxConfig
	GetAccountKeeper() authkeeper.AccountKeeper
}

// Backend bridges JSON-RPC calls to the Syreen app state
type Backend struct {
	app       AppStateProvider
	cometRPC  *rpcclient.HTTP
	chainID   *big.Int
	logger    log.Logger
}

// NewBackend creates a new JSON-RPC backend
func NewBackend(app AppStateProvider, cometRPCAddr string, logger log.Logger) (*Backend, error) {
	client, err := rpcclient.New(cometRPCAddr, "/websocket")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to CometBFT RPC: %w", err)
	}

	return &Backend{
		app:      app,
		cometRPC: client,
		chainID:  big.NewInt(evmtypes.DefaultChainConfig().ChainID.Int64()),
		logger:   logger,
	}, nil
}

// newSDKContext creates an SDK context from the committed multistore,
// bypassing the IAVL version bug (same pattern as query_handler.go)
func (b *Backend) newSDKContext() sdk.Context {
	cms := b.app.GetCommitMultiStore().CacheMultiStore()
	return sdk.NewContext(cms, cmtproto.Header{}, false, b.app.GetLogger())
}

// ChainID returns the EVM chain ID
func (b *Backend) ChainID() *big.Int {
	return b.chainID
}

// BlockNumber returns the latest block height
func (b *Backend) BlockNumber() (int64, error) {
	status, err := b.cometRPC.Status(context.Background())
	if err != nil {
		return 0, err
	}
	return status.SyncInfo.LatestBlockHeight, nil
}

// GetBalance returns the balance of an address in usyreen
func (b *Backend) GetBalance(addr common.Address) (*big.Int, error) {
	ctx := b.newSDKContext()
	accAddr := evmkeeper.Bech32FromEthAddress(addr)
	coins := b.app.GetBankKeeper().SpendableCoins(ctx, accAddr)
	amount := coins.AmountOf("usyreen")
	return amount.BigInt(), nil
}

// GetNonce returns the EVM nonce for an address
func (b *Backend) GetNonce(addr common.Address) (uint64, error) {
	ctx := b.newSDKContext()
	return b.app.GetEVMKeeper().GetEVMNonce(ctx, addr), nil
}

// GetCode returns the contract bytecode at an address
func (b *Backend) GetCode(addr common.Address) ([]byte, error) {
	ctx := b.newSDKContext()
	return b.app.GetEVMKeeper().GetCode(ctx, addr), nil
}

// GetStorageAt returns the value of a storage slot
func (b *Backend) GetStorageAt(addr common.Address, slot common.Hash) (common.Hash, error) {
	ctx := b.newSDKContext()
	return b.app.GetEVMKeeper().GetStorageAt(ctx, addr, slot), nil
}

// Call performs a read-only EVM call
func (b *Backend) Call(from, to *common.Address, data []byte, gas uint64) ([]byte, error) {
	ctx := b.newSDKContext()

	var fromAddr common.Address
	if from != nil {
		fromAddr = *from
	}

	if gas == 0 {
		gas = 30_000_000
	}

	return b.app.GetEVMKeeper().QueryEVM(ctx, fromAddr, to, data, gas)
}

// EstimateGas estimates the gas needed for a call
func (b *Backend) EstimateGas(from, to *common.Address, data []byte, value *big.Int) (uint64, error) {
	ctx := b.newSDKContext()

	var fromAddr common.Address
	if from != nil {
		fromAddr = *from
	}

	// Run the call and measure gas
	_, err := b.app.GetEVMKeeper().QueryEVM(ctx, fromAddr, to, data, 30_000_000)
	if err != nil {
		return 0, err
	}

	// I-15: Return improved gas estimates with safety margins.
	// The EVM module doesn't precisely track gas consumption, so these are
	// conservative defaults. Contract creation and calls use higher values
	// to avoid out-of-gas failures.
	if to == nil {
		// Contract creation — 2M gas to accommodate large bytecode deployment
		return 2_000_000, nil
	}
	if len(data) == 0 {
		// Simple transfer — standard 21K gas (correct for ETH transfers)
		return 21_000, nil
	}
	// Contract call — 500K gas with safety margin for complex interactions
	return 500_000, nil
}

// GasPrice returns the current gas price
func (b *Backend) GasPrice() *big.Int {
	return big.NewInt(1000) // 0.001 usyreen per gas unit
}

// BroadcastRawEthTx decodes a raw signed Ethereum transaction, wraps it in a
// Cosmos SDK transaction, and broadcasts it via CometBFT.
func (b *Backend) BroadcastRawEthTx(rawTxHex string) (string, error) {
	// 1. Decode raw Ethereum transaction
	rawTxBytes, err := hex.DecodeString(strings.TrimPrefix(rawTxHex, "0x"))
	if err != nil {
		return "", fmt.Errorf("invalid raw tx hex: %w", err)
	}

	ethTx := new(ethtypes.Transaction)
	if err := ethTx.UnmarshalBinary(rawTxBytes); err != nil {
		return "", fmt.Errorf("failed to decode ethereum tx: %w", err)
	}

	// 2. Verify chain ID
	txChainID := ethTx.ChainId()
	if txChainID.Cmp(b.chainID) != 0 {
		return "", fmt.Errorf("chain ID mismatch: tx has %s, expected %s", txChainID.String(), b.chainID.String())
	}

	// 3. Recover sender
	signer := ethtypes.LatestSignerForChainID(b.chainID)
	sender, err := ethtypes.Sender(signer, ethTx)
	if err != nil {
		return "", fmt.Errorf("failed to recover sender: %w", err)
	}

	// 4. Recover the compressed public key
	v, r, s := ethTx.RawSignatureValues()
	txHash := signer.Hash(ethTx)
	sig := make([]byte, 65)
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	copy(sig[32-len(rBytes):32], rBytes)
	copy(sig[64-len(sBytes):64], sBytes)

	vByte := byte(v.Uint64())
	if vByte >= 27 {
		vByte -= 27
	}
	if vByte > 1 {
		vByte = byte(v.Uint64() - b.chainID.Uint64()*2 - 35)
	}
	sig[64] = vByte

	pubKeyBytes, err := crypto.Ecrecover(txHash.Bytes(), sig)
	if err != nil {
		return "", fmt.Errorf("failed to recover public key: %w", err)
	}
	pubKeyUncompressed, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal public key: %w", err)
	}
	compressedPubKey := crypto.CompressPubkey(pubKeyUncompressed)
	cosmosPubKey := &secp256k1.PubKey{Key: compressedPubKey}

	// 5. Build MsgEthereumTx
	senderBech32 := evmkeeper.Bech32FromEthAddress(sender)
	toHex := ""
	if ethTx.To() != nil {
		toHex = ethTx.To().Hex()
	}

	msg := &evmtypes.MsgEthereumTx{
		From:     senderBech32.String(),
		To:       toHex,
		Value:    ethTx.Value().String(),
		GasLimit: ethTx.Gas(),
		Data:     hex.EncodeToString(ethTx.Data()),
		Nonce:    ethTx.Nonce(),
		RawTxHex: rawTxHex,
	}

	// 6. Build Cosmos SDK transaction
	txConfig := b.app.GetTxConfig()
	txBuilder := txConfig.NewTxBuilder()

	if err := txBuilder.SetMsgs(msg); err != nil {
		return "", fmt.Errorf("failed to set messages: %w", err)
	}

	// Set fee: gasLimit * gasPrice (in usyreen)
	gasPrice := ethTx.GasPrice()
	if gasPrice == nil || gasPrice.Sign() == 0 {
		// For EIP-1559 txs, use maxFeePerGas
		gasPrice = ethTx.GasFeeCap()
	}
	if gasPrice == nil || gasPrice.Sign() == 0 {
		gasPrice = big.NewInt(1000) // fallback
	}
	feeAmount := new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(ethTx.Gas()))
	txBuilder.SetGasLimit(ethTx.Gas())
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("usyreen", sdkmath.NewIntFromBigInt(feeAmount))))

	// 7. Get account info for signing
	ctx := b.newSDKContext()
	ak := b.app.GetAccountKeeper()
	acc := ak.GetAccount(ctx, senderBech32)

	var accNum, accSeq uint64
	if acc != nil {
		accNum = acc.GetAccountNumber()
		accSeq = acc.GetSequence()
	}

	// 8. Set signature with the recovered public key
	// We use SIGN_MODE_DIRECT but the actual signature bytes don't matter
	// because EVMSigVerificationDecorator skips Cosmos sig verification
	// and EVMAccountDecorator already verifies the Ethereum signature.
	sigData := signing.SingleSignatureData{
		SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
		Signature: sig[:64], // R || S (without V) — placeholder, not verified by Cosmos
	}
	sigV2 := signing.SignatureV2{
		PubKey:   cosmosPubKey,
		Data:     &sigData,
		Sequence: accSeq,
	}
	if err := txBuilder.SetSignatures(sigV2); err != nil {
		return "", fmt.Errorf("failed to set signatures: %w", err)
	}

	// 9. Set the signer info (account number) via signing
	signerData := authsigning.SignerData{
		ChainID:       "syreen-1",
		AccountNumber: accNum,
		Sequence:      accSeq,
	}
	_ = signerData // Used conceptually; the signature is already set above

	// 10. Encode the transaction
	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return "", fmt.Errorf("failed to encode tx: %w", err)
	}

	// 11. Broadcast via CometBFT
	result, err := b.cometRPC.BroadcastTxSync(context.Background(), txBytes)
	if err != nil {
		return "", fmt.Errorf("broadcast failed: %w", err)
	}

	if result.Code != 0 {
		return "", fmt.Errorf("tx failed (code %d): %s", result.Code, result.Log)
	}

	// Return the Cosmos tx hash as an Ethereum-compatible hash
	ethTxHash := fmt.Sprintf("0x%s", hex.EncodeToString(result.Hash))
	b.logger.Info("Broadcast raw Ethereum tx",
		"eth_sender", sender.Hex(),
		"cosmos_hash", ethTxHash,
		"nonce", ethTx.Nonce(),
		"gas", ethTx.Gas(),
	)

	return ethTxHash, nil
}

// GetBlockByNumber returns block data in Ethereum format
func (b *Backend) GetBlockByNumber(blockNum int64, fullTx bool) (*RPCBlock, error) {
	height := &blockNum

	block, err := b.cometRPC.Block(context.Background(), height)
	if err != nil {
		return nil, err
	}

	txs := make([]interface{}, 0, len(block.Block.Txs))
	for i, tx := range block.Block.Txs {
		txHash := fmt.Sprintf("0x%s", hex.EncodeToString(tx.Hash()))
		if fullTx {
			idx := HexInt64(int64(i))
			txs = append(txs, &RPCTransaction{
				BlockHash:        strPtr(fmt.Sprintf("0x%x", block.BlockID.Hash.Bytes())),
				BlockNumber:      strPtr(HexInt64(block.Block.Height)),
				Hash:             txHash,
				TransactionIndex: &idx,
				From:             EmptyAddress,
				Gas:              "0x5208",
				GasPrice:         "0x3e8",
				Input:            "0x",
				Nonce:            "0x0",
				Value:            "0x0",
				V:                "0x0",
				R:                "0x0",
				S:                "0x0",
				Type:             "0x0",
			})
		} else {
			txs = append(txs, txHash)
		}
	}

	blockHash := fmt.Sprintf("0x%x", block.BlockID.Hash.Bytes())
	parentHash := EmptyHash
	if block.Block.Height > 1 {
		parentHash = fmt.Sprintf("0x%x", block.Block.LastBlockID.Hash.Bytes())
	}

	return &RPCBlock{
		Number:           HexInt64(block.Block.Height),
		Hash:             blockHash,
		ParentHash:       parentHash,
		Nonce:            "0x0000000000000000",
		Sha3Uncles:       EmptyHash,
		LogsBloom:        EmptyBloom,
		TransactionsRoot: fmt.Sprintf("0x%x", block.Block.DataHash.Bytes()),
		StateRoot:        fmt.Sprintf("0x%x", block.Block.AppHash.Bytes()),
		ReceiptsRoot:     EmptyHash,
		Miner:            fmt.Sprintf("0x%x", block.Block.ProposerAddress.Bytes()),
		Difficulty:       "0x0",
		TotalDifficulty:  "0x0",
		ExtraData:        "0x",
		Size:             HexInt64(int64(block.Block.Size())),
		GasLimit:         "0x5f5e100", // 100M
		GasUsed:          HexInt64(int64(len(block.Block.Txs) * 21000)),
		Timestamp:        HexInt64(block.Block.Time.Unix()),
		Transactions:     txs,
		Uncles:           []string{},
		MixHash:          EmptyHash,
		BaseFeePerGas:    "0x1",
	}, nil
}

// GetBlockByHash returns block data by hash
func (b *Backend) GetBlockByHash(hash common.Hash, fullTx bool) (*RPCBlock, error) {
	// Query CometBFT for the block by hash
	hashBytes := hash.Bytes()
	block, err := b.cometRPC.BlockByHash(context.Background(), hashBytes)
	if err != nil {
		return nil, err
	}
	if block == nil || block.Block == nil {
		return nil, nil
	}

	return b.GetBlockByNumber(block.Block.Height, fullTx)
}

// GetTransactionByHash returns transaction data by hash
func (b *Backend) GetTransactionByHash(hash string) (*RPCTransaction, error) {
	// Strip 0x prefix
	hashStr := strings.TrimPrefix(hash, "0x")
	hashBytes, err := hex.DecodeString(hashStr)
	if err != nil {
		return nil, fmt.Errorf("invalid hash: %w", err)
	}

	tx, err := b.cometRPC.Tx(context.Background(), hashBytes, false)
	if err != nil {
		return nil, nil // tx not found
	}

	blockNum := HexInt64(tx.Height)
	idx := HexInt64(int64(tx.Index))
	txHash := fmt.Sprintf("0x%s", hex.EncodeToString(tx.Hash))

	return &RPCTransaction{
		BlockHash:        strPtr(blockNum),
		BlockNumber:      strPtr(blockNum),
		Hash:             txHash,
		TransactionIndex: &idx,
		From:             EmptyAddress,
		Gas:              HexInt64(int64(tx.TxResult.GasWanted)),
		GasPrice:         "0x3e8",
		Input:            "0x",
		Nonce:            "0x0",
		Value:            "0x0",
		V:                "0x0",
		R:                "0x0",
		S:                "0x0",
		Type:             "0x0",
	}, nil
}

// GetTransactionReceipt returns a transaction receipt
func (b *Backend) GetTransactionReceipt(hash string) (*RPCReceipt, error) {
	hashStr := strings.TrimPrefix(hash, "0x")
	hashBytes, err := hex.DecodeString(hashStr)
	if err != nil {
		return nil, fmt.Errorf("invalid hash: %w", err)
	}

	tx, err := b.cometRPC.Tx(context.Background(), hashBytes, false)
	if err != nil {
		return nil, nil // tx not found
	}

	status := "0x1" // success
	if tx.TxResult.Code != 0 {
		status = "0x0" // failure
	}

	blockNum := HexInt64(tx.Height)
	idx := HexInt64(int64(tx.Index))
	txHash := fmt.Sprintf("0x%s", hex.EncodeToString(tx.Hash))

	return &RPCReceipt{
		TransactionHash:   txHash,
		TransactionIndex:  idx,
		BlockHash:         blockNum,
		BlockNumber:       blockNum,
		From:              EmptyAddress,
		To:                nil,
		CumulativeGasUsed: HexInt64(int64(tx.TxResult.GasUsed)),
		GasUsed:           HexInt64(int64(tx.TxResult.GasUsed)),
		ContractAddress:   nil,
		Logs:              []RPCLog{},
		LogsBloom:         EmptyBloom,
		Status:            status,
		EffectiveGasPrice: "0x3e8",
		Type:              "0x0",
	}, nil
}

// PeerCount returns the number of connected peers
func (b *Backend) PeerCount() (int, error) {
	netInfo, err := b.cometRPC.NetInfo(context.Background())
	if err != nil {
		return 0, err
	}
	return netInfo.NPeers, nil
}

// Syncing returns false if the node is synced
func (b *Backend) Syncing() (interface{}, error) {
	status, err := b.cometRPC.Status(context.Background())
	if err != nil {
		return false, err
	}
	if status.SyncInfo.CatchingUp {
		return map[string]string{
			"startingBlock": HexInt64(0),
			"currentBlock":  HexInt64(status.SyncInfo.LatestBlockHeight),
			"highestBlock":  HexInt64(status.SyncInfo.LatestBlockHeight),
		}, nil
	}
	return false, nil
}

// ParseBlockNumber converts block number string ("latest", "earliest", hex) to int64
func (b *Backend) ParseBlockNumber(blockNrOrHash string) (int64, error) {
	switch blockNrOrHash {
	case "latest", "pending", "":
		return b.BlockNumber()
	case "earliest":
		return 1, nil
	default:
		blockNrOrHash = strings.TrimPrefix(blockNrOrHash, "0x")
		n, err := strconv.ParseInt(blockNrOrHash, 16, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid block number: %s", blockNrOrHash)
		}
		return n, nil
	}
}

func strPtr(s string) *string {
	return &s
}
