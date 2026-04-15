package utils

import (
	"strings"
	"testing"
)

func TestGenerateOTP_Length(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"length 4", 4},
		{"length 6", 6},
		{"length 8", 8},
		{"length 0", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otp, err := GenerateOTP(tt.length)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(otp) != tt.length {
				t.Errorf("expected length %d, got %d", tt.length, len(otp))
			}
		})
	}
}

func TestGenerateOTP_DigitsOnly(t *testing.T) {
	otp, err := GenerateOTP(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, ch := range otp {
		if !strings.ContainsRune("0123456789", ch) {
			t.Errorf("otp contains non-digit character: %c", ch)
		}
	}
}

func TestGenerateOTP_Uniqueness(t *testing.T) {
	otp1, err1 := GenerateOTP(6)
	otp2, err2 := GenerateOTP(6)

	if err1 != nil || err2 != nil {
		t.Fatalf("unexpected error: %v %v", err1, err2)
	}

	// Not guaranteed, but highly unlikely to match
	if otp1 == otp2 {
		t.Logf("warning: two OTPs matched (rare but possible): %s", otp1)
	}
}

func TestGenerateOTP_MultipleCalls(t *testing.T) {
	for i := 0; i < 100; i++ {
		otp, err := GenerateOTP(6)
		if err != nil {
			t.Fatalf("unexpected error on iteration %d: %v", i, err)
		}

		if len(otp) != 6 {
			t.Errorf("expected length 6, got %d", len(otp))
		}
	}
}
