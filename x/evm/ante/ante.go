package ante

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	evmkeeper "syreen/x/evm/keeper"
	evmtypes "syreen/x/evm/types"
)

// contextKey is a private type for context keys in this package
type contextKey string

const evmRawTxVerifiedKey contextKey = "evm_raw_tx_verified"

// EVMAccountKeeper defines the account keeper methods needed by EVM ante decorators
type EVMAccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(ctx context.Context, acc sdk.AccountI)
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

// IsEVMRawTx checks if the tx contains a single MsgEthereumTx with RawTxHex set
func IsEVMRawTx(tx sdk.Tx) (*evmtypes.MsgEthereumTx, bool) {
	msgs := tx.GetMsgs()
	if len(msgs) != 1 {
		return nil, false
	}
	ethMsg, ok := msgs[0].(*evmtypes.MsgEthereumTx)
	if !ok || ethMsg.RawTxHex == "" {
		return nil, false
	}
	return ethMsg, true
}

// IsEVMVerified checks if the context has the EVM raw tx verified flag
func IsEVMVerified(ctx sdk.Context) bool {
	val := ctx.Value(evmRawTxVerifiedKey)
	if val == nil {
		return false
	}
	verified, ok := val.(bool)
	return ok && verified
}

// ---------------------------------------------------------------------------
// EVMAccountDecorator — ensures sender account exists for raw Ethereum txs,
// verifies the Ethereum signature, and sets the verified flag in context.
// Must be placed BEFORE fee deduction and signature decorators.
// ---------------------------------------------------------------------------

type EVMAccountDecorator struct {
	ak        EVMAccountKeeper
	evmKeeper evmkeeper.Keeper
}

func NewEVMAccountDecorator(ak EVMAccountKeeper, ek evmkeeper.Keeper) EVMAccountDecorator {
	return EVMAccountDecorator{ak: ak, evmKeeper: ek}
}

func (d EVMAccountDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	ethMsg, isEVM := IsEVMRawTx(tx)
	if !isEVM {
		return next(ctx, tx, simulate)
	}

	// Decode raw Ethereum transaction
	rawTxBytes, err := hex.DecodeString(strings.TrimPrefix(ethMsg.RawTxHex, "0x"))
	if err != nil {
		return ctx, fmt.Errorf("invalid raw tx hex: %w", err)
	}

	ethTx := new(ethtypes.Transaction)
	if err := ethTx.UnmarshalBinary(rawTxBytes); err != nil {
		return ctx, fmt.Errorf("failed to decode raw ethereum tx: %w", err)
	}

	// Verify chain ID
	chainID := d.evmKeeper.GetChainConfig().ChainID
	signer := ethtypes.LatestSignerForChainID(chainID)

	// Recover sender from Ethereum signature
	sender, err := ethtypes.Sender(signer, ethTx)
	if err != nil {
		return ctx, fmt.Errorf("failed to recover sender from ethereum signature: %w", err)
	}

	// Verify sender matches msg.From
	expectedFrom := evmkeeper.Bech32FromEthAddress(sender)
	if ethMsg.From != expectedFrom.String() {
		return ctx, fmt.Errorf("sender mismatch: recovered %s but msg has %s", expectedFrom.String(), ethMsg.From)
	}

	// Ensure account exists for the sender
	acc := d.ak.GetAccount(ctx, expectedFrom)
	if acc == nil {
		acc = d.ak.NewAccountWithAddress(ctx, expectedFrom)
		d.ak.SetAccount(ctx, acc)
	}

	// Set the public key on the account if not already set
	if acc.GetPubKey() == nil {
		// Recover the full public key from the Ethereum signature
		v, r, s := ethTx.RawSignatureValues()
		txHash := signer.Hash(ethTx)

		// Reconstruct the 65-byte signature [R || S || V]
		sig := make([]byte, 65)
		rBytes := r.Bytes()
		sBytes := s.Bytes()
		copy(sig[32-len(rBytes):32], rBytes)
		copy(sig[64-len(sBytes):64], sBytes)

		// V needs to be 0 or 1 for crypto.Ecrecover
		vByte := byte(v.Uint64())
		if vByte >= 27 {
			vByte -= 27
		}
		// For EIP-155, V = chainID * 2 + 35 + recovery_id
		if vByte > 1 {
			vByte = byte(v.Uint64() - chainID.Uint64()*2 - 35)
		}
		sig[64] = vByte

		pubKeyBytes, err := crypto.Ecrecover(txHash.Bytes(), sig)
		if err != nil {
			return ctx, fmt.Errorf("failed to recover public key: %w", err)
		}

		// Convert to compressed secp256k1 public key (33 bytes)
		pubKeyUncompressed, err := crypto.UnmarshalPubkey(pubKeyBytes)
		if err != nil {
			return ctx, fmt.Errorf("failed to unmarshal public key: %w", err)
		}
		compressedPubKey := crypto.CompressPubkey(pubKeyUncompressed)

		cosmosPubKey := &secp256k1.PubKey{Key: compressedPubKey}

		err = acc.SetPubKey(cosmosPubKey)
		if err != nil {
			return ctx, fmt.Errorf("failed to set public key: %w", err)
		}
		d.ak.SetAccount(ctx, acc)
	}

	// Mark context as EVM-verified so sig verification decorators can skip
	ctx = ctx.WithValue(evmRawTxVerifiedKey, true)

	return next(ctx, tx, simulate)
}

// ---------------------------------------------------------------------------
// EVMSigVerificationDecorator — replaces the standard SigVerificationDecorator.
// For EVM raw txs (flagged by EVMAccountDecorator), skips Cosmos signature
// verification. For all other txs, delegates to the original decorator.
// ---------------------------------------------------------------------------

type EVMSigVerificationDecorator struct {
	original sdk.AnteDecorator
}

func NewEVMSigVerificationDecorator(original sdk.AnteDecorator) EVMSigVerificationDecorator {
	return EVMSigVerificationDecorator{original: original}
}

func (d EVMSigVerificationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if IsEVMVerified(ctx) {
		// EVM signature already verified by EVMAccountDecorator — skip Cosmos sig check
		return next(ctx, tx, simulate)
	}
	return d.original.AnteHandle(ctx, tx, simulate, next)
}

// ---------------------------------------------------------------------------
// EVMSkipDecorator — wraps any ante decorator and skips it for EVM raw txs.
// Used for decorators that don't apply to Ethereum-signed transactions
// (e.g., SetPubKey, ValidateSigCount, SigGasConsume).
//
// SECURITY WARNING: EVMSkipDecorator must NOT be used to wrap fee deduction
// decorators (e.g., DeductFeeDecorator). Wrapping a fee decorator would allow
// EVM transactions to bypass fee payment entirely, enabling fee-free spam.
// Only use this for signature-related decorators that are redundant for EVM
// transactions (whose signatures are already verified by EVMAccountDecorator).
// ---------------------------------------------------------------------------

type EVMSkipDecorator struct {
	original sdk.AnteDecorator
}

func NewEVMSkipDecorator(original sdk.AnteDecorator) EVMSkipDecorator {
	return EVMSkipDecorator{original: original}
}

// AnteHandle skips the wrapped decorator for EVM-verified transactions.
// This is safe ONLY for non-fee decorators. See the security warning above.
func (d EVMSkipDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if IsEVMVerified(ctx) {
		return next(ctx, tx, simulate)
	}
	return d.original.AnteHandle(ctx, tx, simulate, next)
}

// Ensure all decorators implement the interface
var _ sdk.AnteDecorator = EVMAccountDecorator{}
var _ sdk.AnteDecorator = EVMSigVerificationDecorator{}
var _ sdk.AnteDecorator = EVMSkipDecorator{}

// Ensure cryptotypes.PubKey is used
var _ cryptotypes.PubKey = &secp256k1.PubKey{}
