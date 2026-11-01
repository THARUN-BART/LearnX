package tests

import (
	"testing"

	"github.com/THARUN-BART/LearnX/backend/internal/crypto"
)

func TestCrypto_EncryptDecrypt(t *testing.T) {
	key := crypto.DeriveKey("super-secret-key-12345")
	originalText := "MySecretPassword123!"

	encrypted, err := crypto.Encrypt([]byte(originalText), key)
	if err != nil {
		t.Fatalf("unexpected encryption error: %v", err)
	}

	if encrypted == originalText {
		t.Fatalf("encrypted text must not match plaintext")
	}

	decrypted, err := crypto.Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("unexpected decryption error: %v", err)
	}

	if string(decrypted) != originalText {
		t.Fatalf("expected %q, got %q", originalText, string(decrypted))
	}
}

func TestCrypto_DecryptWithInvalidKey(t *testing.T) {
	key1 := crypto.DeriveKey("secret-1")
	key2 := crypto.DeriveKey("secret-2")

	encrypted, err := crypto.Encrypt([]byte("confidential"), key1)
	if err != nil {
		t.Fatalf("unexpected encryption error: %v", err)
	}

	_, err = crypto.Decrypt(encrypted, key2)
	if err == nil {
		t.Fatalf("expected decryption error with wrong key, got nil")
	}
}

func TestCrypto_GenerateRandomToken(t *testing.T) {
	tok1, err := crypto.GenerateRandomToken(32)
	if err != nil {
		t.Fatalf("unexpected token error: %v", err)
	}
	tok2, err := crypto.GenerateRandomToken(32)
	if err != nil {
		t.Fatalf("unexpected token error: %v", err)
	}

	if tok1 == tok2 {
		t.Fatalf("random tokens must not match")
	}
	if len(tok1) != 64 {
		t.Fatalf("expected 64 hex characters, got %d", len(tok1))
	}
}

func TestCrypto_HashString(t *testing.T) {
	h1 := crypto.HashString("test-string")
	h2 := crypto.HashString("test-string")
	h3 := crypto.HashString("other-string")

	if h1 != h2 {
		t.Fatalf("hashes of identical string should match")
	}
	if h1 == h3 {
		t.Fatalf("hashes of different strings must not match")
	}
}
