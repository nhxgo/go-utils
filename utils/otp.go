package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP(length int) (string, error) {
	const digits = "0123456789"

	otp := make([]byte, length)

	for i := range otp {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", fmt.Errorf("failed to generate otp: %w", err)
		}
		otp[i] = digits[n.Int64()]
	}

	return string(otp), nil
}
