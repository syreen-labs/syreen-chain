package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
)

// EncryptTx encrypts transaction bytes using AES-256-GCM with the provided key.
// The key must be exactly 32 bytes (AES-256). A random nonce is generated internally.
// Returns the ciphertext (with GCM authentication tag appended) and the nonce used.
func EncryptTx(plaintext []byte, key []byte) (ciphertext []byte, nonce []byte, err error) {
	if len(key) != 32 {
		return nil, nil, fmt.Errorf("key must be 32 bytes for AES-256, got %d", len(key))
	}
	if len(plaintext) == 0 {
		return nil, nil, fmt.Errorf("plaintext cannot be empty")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce = make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext = gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

// DecryptTx decrypts transaction bytes using AES-256-GCM.
// The key must be exactly 32 bytes. The nonce must match the one returned by EncryptTx.
// Returns the original plaintext or an error if decryption/authentication fails.
func DecryptTx(ciphertext []byte, key []byte, nonce []byte) (plaintext []byte, err error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be 32 bytes for AES-256, got %d", len(key))
	}
	if len(ciphertext) == 0 {
		return nil, fmt.Errorf("ciphertext cannot be empty")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	if len(nonce) != gcm.NonceSize() {
		return nil, fmt.Errorf("nonce must be %d bytes, got %d", gcm.NonceSize(), len(nonce))
	}

	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed (authentication error): %w", err)
	}
	return plaintext, nil
}

// DeriveCommitKey derives a deterministic 32-byte AES-256 encryption key from a
// sender's private commitment seed and the target block height.
//
// The derivation uses SHA-256: key = SHA256(seed || blockHeight_LE).
// This ensures each block height produces a unique encryption key for the same seed,
// preventing key reuse across blocks.
func DeriveCommitKey(seed []byte, blockHeight int64) []byte {
	h := sha256.New()
	h.Write(seed)

	heightBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(heightBytes, uint64(blockHeight))
	h.Write(heightBytes)

	return h.Sum(nil) // 32 bytes — suitable for AES-256
}
