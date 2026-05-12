package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt_Roundtrip(t *testing.T) {
	seed := []byte("test-seed-for-key-derivation!!")
	key := DeriveCommitKey(seed, 100)

	plaintext := []byte("this is a secret transaction body")
	ciphertext, nonce, err := EncryptTx(plaintext, key)
	if err != nil {
		t.Fatalf("EncryptTx failed: %v", err)
	}

	if len(ciphertext) == 0 {
		t.Fatal("ciphertext is empty")
	}
	if len(nonce) == 0 {
		t.Fatal("nonce is empty")
	}
	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("ciphertext should not equal plaintext")
	}

	decrypted, err := DecryptTx(ciphertext, key, nonce)
	if err != nil {
		t.Fatalf("DecryptTx failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("decrypted does not match plaintext: got %q, want %q", decrypted, plaintext)
	}
}

func TestDecryptTx_WrongKey(t *testing.T) {
	key1 := DeriveCommitKey([]byte("seed-1"), 100)
	key2 := DeriveCommitKey([]byte("seed-2"), 100)

	plaintext := []byte("secret data")
	ciphertext, nonce, err := EncryptTx(plaintext, key1)
	if err != nil {
		t.Fatalf("EncryptTx failed: %v", err)
	}

	_, err = DecryptTx(ciphertext, key2, nonce)
	if err == nil {
		t.Fatal("expected decryption to fail with wrong key")
	}
}

func TestDecryptTx_TamperedCiphertext(t *testing.T) {
	key := DeriveCommitKey([]byte("seed"), 50)
	plaintext := []byte("important tx data")
	ciphertext, nonce, err := EncryptTx(plaintext, key)
	if err != nil {
		t.Fatalf("EncryptTx failed: %v", err)
	}

	// Tamper with ciphertext
	ciphertext[0] ^= 0xFF

	_, err = DecryptTx(ciphertext, key, nonce)
	if err == nil {
		t.Fatal("expected decryption to fail with tampered ciphertext")
	}
}

func TestEncryptTx_InvalidKeyLength(t *testing.T) {
	_, _, err := EncryptTx([]byte("data"), []byte("short-key"))
	if err == nil {
		t.Fatal("expected error for invalid key length")
	}
}

func TestEncryptTx_EmptyPlaintext(t *testing.T) {
	key := DeriveCommitKey([]byte("seed"), 1)
	_, _, err := EncryptTx([]byte{}, key)
	if err == nil {
		t.Fatal("expected error for empty plaintext")
	}
}

func TestDecryptTx_EmptyCiphertext(t *testing.T) {
	key := DeriveCommitKey([]byte("seed"), 1)
	_, err := DecryptTx([]byte{}, key, make([]byte, 12))
	if err == nil {
		t.Fatal("expected error for empty ciphertext")
	}
}

func TestDeriveCommitKey_DifferentHeights(t *testing.T) {
	seed := []byte("deterministic-seed")
	key1 := DeriveCommitKey(seed, 100)
	key2 := DeriveCommitKey(seed, 101)
	key3 := DeriveCommitKey(seed, 100)

	if bytes.Equal(key1, key2) {
		t.Fatal("keys for different heights should differ")
	}
	if !bytes.Equal(key1, key3) {
		t.Fatal("keys for the same seed and height should be identical")
	}
	if len(key1) != 32 {
		t.Fatalf("key length should be 32 bytes, got %d", len(key1))
	}
}

func TestDeriveCommitKey_DifferentSeeds(t *testing.T) {
	key1 := DeriveCommitKey([]byte("seed-a"), 100)
	key2 := DeriveCommitKey([]byte("seed-b"), 100)

	if bytes.Equal(key1, key2) {
		t.Fatal("keys for different seeds should differ")
	}
}
