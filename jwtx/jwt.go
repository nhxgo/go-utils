package jwtx

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const Issuer = "nhxgo"

func Generate(
	privateKey *rsa.PrivateKey,
	expiry time.Duration,
	payload interface{},
) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(expiry).Unix(),
		"iss": Issuer,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	var payloadMap map[string]any
	if err := json.Unmarshal(payloadBytes, &payloadMap); err != nil {
		return "", err
	}

	for k, val := range payloadMap {
		if k == "exp" || k == "iat" || k == "iss" {
			continue
		}
		claims[k] = val
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(privateKey)
}
func Validate[T jwt.Claims](
	tokenStr string,
	publicKey *rsa.PublicKey,
	claims T,
) (T, error) {

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodRS256 {
			return claims, fmt.Errorf("invalid signing method")
		}
		return publicKey, nil
	})
	if err != nil {
		return claims, err
	}

	c, ok := token.Claims.(T)
	if !ok || !token.Valid {
		return claims, fmt.Errorf("invalid token")
	}
	cIssuer, err := c.GetIssuer()
	if err != nil {
		return claims, err
	}
	if cIssuer != Issuer {
		return claims, fmt.Errorf("Issure mismatch")
	}

	return c, nil
}
