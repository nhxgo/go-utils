package jwtx

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// --- helper to generate RSA keys ---
func generateKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate keys: %v", err)
	}
	return priv, &priv.PublicKey
}

// --- Test Generate + Validate success ---
func TestGenerateAndValidate_Success(t *testing.T) {
	priv, pub := generateKeys(t)

	payload := map[string]any{
		"user_id": 123,
		"role":    "admin",
	}

	token, err := Generate(priv, time.Hour, payload)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	claims := jwt.MapClaims{}
	result, err := Validate(token, pub, claims)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}

	if result["user_id"] != float64(123) {
		t.Errorf("expected user_id 123, got %v", result["user_id"])
	}

	if result["role"] != "admin" {
		t.Errorf("expected role admin, got %v", result["role"])
	}
}

// --- Test invalid token ---
func TestValidate_InvalidToken(t *testing.T) {
	_, pub := generateKeys(t)

	claims := jwt.MapClaims{}
	_, err := Validate("invalid.token.here", pub, claims)

	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

// --- Test wrong public key ---
func TestValidate_WrongKey(t *testing.T) {
	priv1, _ := generateKeys(t)
	_, pub2 := generateKeys(t)

	token, err := Generate(priv1, time.Hour, map[string]any{"id": 1})
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	claims := jwt.MapClaims{}
	_, err = Validate(token, pub2, claims)

	if err == nil {
		t.Fatal("expected error with wrong public key")
	}
}

// --- Test issuer mismatch ---
func TestValidate_IssuerMismatch(t *testing.T) {
	priv, pub := generateKeys(t)

	// manually create token with wrong issuer
	claims := jwt.MapClaims{
		"iss": "wrong",
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := token.SignedString(priv)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = Validate(tokenStr, pub, jwt.MapClaims{})
	if err == nil {
		t.Fatal("expected issuer mismatch error")
	}
}

// --- Test expired token ---
func TestValidate_ExpiredToken(t *testing.T) {
	priv, pub := generateKeys(t)

	token, err := Generate(priv, -time.Hour, map[string]any{"id": 1})
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	_, err = Validate(token, pub, jwt.MapClaims{})
	if err == nil {
		t.Fatal("expected expired token error")
	}
}

// --- Test reserved claims override protection ---
func TestGenerate_ReservedClaimsNotOverridden(t *testing.T) {
	priv, pub := generateKeys(t)

	payload := map[string]any{
		"iss": "hacker",
		"exp": 9999999999,
		"iat": 0,
		"id":  1,
	}

	token, err := Generate(priv, time.Hour, payload)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	claims := jwt.MapClaims{}
	result, err := Validate(token, pub, claims)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}

	if result["iss"] != Issuer {
		t.Errorf("issuer should not be overridden")
	}
}

// --- Test LoadPrivateKey / LoadPublicKey ---
func TestLoadKeys(t *testing.T) {
	priv, _ := generateKeys(t)

	privFile := "test_private.pem"
	pubFile := "test_public.pem"

	// --- encode private key ---
	privBytes := x509.MarshalPKCS1PrivateKey(priv)
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})

	err := os.WriteFile(privFile, privPEM, 0644)
	if err != nil {
		t.Fatalf("failed to write private key: %v", err)
	}
	defer os.Remove(privFile)

	// --- encode public key ---
	pubBytes := x509.MarshalPKCS1PublicKey(&priv.PublicKey)
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubBytes,
	})

	err = os.WriteFile(pubFile, pubPEM, 0644)
	if err != nil {
		t.Fatalf("failed to write public key: %v", err)
	}
	defer os.Remove(pubFile)

	// --- test loading ---
	_, err = LoadPrivateKey(privFile)
	if err != nil {
		t.Errorf("failed to load private key: %v", err)
	}

	_, err = LoadPublicKey(pubFile)
	if err != nil {
		t.Errorf("failed to load public key: %v", err)
	}
}
