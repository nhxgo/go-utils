package errx

import (
	"errors"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

// --- Test Error() method ---
func TestError_String(t *testing.T) {
	err := &Error{
		ErrorCode: BadRequestCode,
		Message:   "invalid input",
	}

	expected := "[BAD_REQUEST] invalid input"

	if err.Error() != expected {
		t.Errorf("expected %s, got %s", expected, err.Error())
	}
}

// --- Test Unwrap ---
func TestError_Unwrap(t *testing.T) {
	baseErr := errors.New("root error")

	err := &Error{
		Err: baseErr,
	}

	if !errors.Is(err, baseErr) {
		t.Error("unwrap failed")
	}
}

// --- Test NewError ---
func TestNewError(t *testing.T) {
	baseErr := errors.New("test")

	err := NewError(400, BadRequestCode, "bad request", baseErr)

	e, ok := err.(*Error)
	if !ok {
		t.Fatal("expected *Error type")
	}

	if e.HttpStatus != 400 || e.ErrorCode != BadRequestCode || e.Message != "bad request" {
		t.Error("unexpected error fields")
	}
}

// --- Test HTTP helpers ---
func TestHTTPHelpers(t *testing.T) {
	tests := []struct {
		name       string
		errFunc    func() error
		statusCode int
		code       ApiError
	}{
		{"internal", func() error { return InternalServerError(nil) }, http.StatusInternalServerError, InternalServerErrorCode},
		{"bad request", func() error { return BadRequestError(nil, "bad") }, http.StatusBadRequest, BadRequestCode},
		{"unauthorized", func() error { return UnauthorizedError(nil, "unauth") }, http.StatusUnauthorized, UnauthorizedCode},
		{"conflict", func() error { return ConflictError(nil, "conflict") }, http.StatusConflict, ConflictCode},
		{"not found", func() error { return NotFoundError(nil, "not found") }, http.StatusNotFound, NotFoundCode},
		{"method not allowed", func() error { return MethodNotAllowedError(nil) }, http.StatusMethodNotAllowed, MethodNotAllowedCode},
		{"too many requests", func() error { return TooManyRequestError(nil) }, http.StatusTooManyRequests, TooManyRequestsCode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errFunc()

			e, ok := err.(*Error)
			if !ok {
				t.Fatal("expected *Error")
			}

			if e.HttpStatus != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, e.HttpStatus)
			}

			if e.ErrorCode != tt.code {
				t.Errorf("expected code %s, got %s", tt.code, e.ErrorCode)
			}
		})
	}
}

// --- Test Crypto errors ---
func TestCryptoErrors(t *testing.T) {
	baseErr := errors.New("low-level")

	tests := []struct {
		name string
		fn   func(error) error
	}{
		{"decrypt", NewDecryptError},
		{"encrypt", NewEncryptError},
		{"hashing", NewHashingError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn(baseErr)

			e, ok := err.(*Error)
			if !ok {
				t.Fatal("expected *Error")
			}

			if !errors.Is(e.Err, baseErr) {
				t.Error("expected wrapped error")
			}
		})
	}
}

// --- Test MapPgError: nil ---
func TestMapPgError_Nil(t *testing.T) {
	if MapPgError(nil) != nil {
		t.Error("expected nil")
	}
}

// --- Test MapPgError: unique violation ---
func TestMapPgError_UniqueViolation(t *testing.T) {
	pgErr := &pgconn.PgError{Code: "23505"}

	err := MapPgError(pgErr)

	e := err.(*Error)

	if e.HttpStatus != http.StatusConflict {
		t.Error("expected conflict status")
	}
}

// --- Test MapPgError: foreign key ---
func TestMapPgError_ForeignKey(t *testing.T) {
	pgErr := &pgconn.PgError{Code: "23503"}

	err := MapPgError(pgErr)

	e := err.(*Error)

	if e.HttpStatus != http.StatusBadRequest {
		t.Error("expected bad request")
	}
}

// --- Test MapPgError: not null ---
func TestMapPgError_NotNull(t *testing.T) {
	pgErr := &pgconn.PgError{Code: "23502"}

	err := MapPgError(pgErr)

	e := err.(*Error)

	if e.ErrorCode != ValidationErrorCode {
		t.Error("expected validation error")
	}
}

// --- Test MapPgError: unknown pg error ---
func TestMapPgError_UnknownPgError(t *testing.T) {
	pgErr := &pgconn.PgError{Code: "99999"}

	err := MapPgError(pgErr)

	e := err.(*Error)

	if e.ErrorCode != DBErrorCode {
		t.Error("expected DB error code")
	}
}

// --- Test MapPgError: already ApiError ---
func TestMapPgError_AlreadyApiError(t *testing.T) {
	original := BadRequestError(errors.New("bad"), "bad request")

	err := MapPgError(original)

	e := err.(*Error)

	if e.HttpStatus != http.StatusBadRequest {
		t.Error("expected same status")
	}
}

// --- Test MapPgError: generic error ---
func TestMapPgError_Generic(t *testing.T) {
	err := MapPgError(errors.New("random"))

	e := err.(*Error)

	if e.HttpStatus != http.StatusInternalServerError {
		t.Error("expected internal server error")
	}
}
