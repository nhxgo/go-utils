package encrypt

import (
	"encoding/base64"
	"testing"
)

// --- helper: generate valid AES key (32 bytes) ---
func validKey() string {
	key := make([]byte, 32) // AES-256
	for i := range key {
		key[i] = byte(i)
	}
	return base64.StdEncoding.EncodeToString(key)
}

// --- Test Encrypt + Decrypt success ---
func TestEncryptDecrypt_Success(t *testing.T) {
	key := validKey()
	plaintext := "hello world"

	encrypted, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("expected %s, got %s", plaintext, decrypted)
	}
}

// --- Test invalid base64 key ---
func TestEncrypt_InvalidBase64Key(t *testing.T) {
	_, err := Encrypt("data", "invalid-base64")

	if err == nil {
		t.Fatal("expected error for invalid base64 key")
	}
}

// --- Test wrong key decryption ---
func TestDecrypt_WrongKey(t *testing.T) {
	key1 := validKey()
	key2 := validKey() + "extra" // different key

	data, err := Encrypt("secret", key1)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	_, err = Decrypt(data, key2)
	if err == nil {
		t.Fatal("expected error for wrong key")
	}
}

// --- Test tampered ciphertext ---
func TestDecrypt_TamperedData(t *testing.T) {
	key := validKey()

	data, err := Encrypt("secret", key)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	// tamper data
	data[len(data)-1] ^= 0xFF

	_, err = Decrypt(data, key)
	if err == nil {
		t.Fatal("expected error for tampered data")
	}
}

// --- Test short ciphertext ---
func TestDecrypt_InvalidCiphertext(t *testing.T) {
	key := validKey()

	_, err := Decrypt([]byte("short"), key)
	if err == nil {
		t.Fatal("expected error for invalid ciphertext")
	}
}

// --- Test IsValidKey ---
func TestIsValidKey(t *testing.T) {
	valid := validKey()

	if !IsValidKey(valid) {
		t.Error("expected valid key")
	}

	if IsValidKey("invalid") {
		t.Error("expected invalid key")
	}

	if IsValidKey("") {
		t.Error("expected empty key to be invalid")
	}
}

// --- Test Hash consistency ---
func TestHash(t *testing.T) {
	secret := "secret"
	data := "data"

	hash1 := Hash(data, secret)
	hash2 := Hash(data, secret)

	if hash1 != hash2 {
		t.Error("hash should be deterministic")
	}
}

// --- Test HashEmail normalization ---
func TestHashEmail(t *testing.T) {
	secret := "secret"

	h1 := HashEmail("Test@Email.com ", secret)
	h2 := HashEmail("test@email.com", secret)

	if h1 != h2 {
		t.Error("email hashing should normalize input")
	}
}

// --- Test HashPassword + CheckPassword ---
func TestPasswordHashing(t *testing.T) {
	password := "mypassword"

	hash, err := HashPassword(password, 10)
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}

	if !CheckPassword(password, hash) {
		t.Error("expected password to match")
	}

	if CheckPassword("wrong", hash) {
		t.Error("expected password mismatch")
	}
}
