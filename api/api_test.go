package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nhxgo/go-utils/errx"
)

// --- helper struct for response decoding ---
type resp struct {
	Message string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// --- Test JSON success ---
func TestJSON_Success(t *testing.T) {
	rec := httptest.NewRecorder()

	data := map[string]string{"hello": "world"}

	JSON(rec, http.StatusOK, data)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var result map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if result["hello"] != "world" {
		t.Errorf("unexpected response body")
	}
}

// --- Test JSON nil body ---
func TestJSON_NilData(t *testing.T) {
	rec := httptest.NewRecorder()

	JSON(rec, http.StatusNoContent, nil)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}

	if rec.Body.Len() != 0 {
		t.Errorf("expected empty body")
	}
}

// --- Test JSON marshal failure ---
func TestJSON_MarshalError(t *testing.T) {
	rec := httptest.NewRecorder()

	// channels cannot be marshaled
	invalid := make(chan int)

	JSON(rec, http.StatusOK, invalid)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

// --- Test ErrorJSON with errx.Error ---
func TestErrorJSON_WithCustomError(t *testing.T) {
	rec := httptest.NewRecorder()

	apiErr := errx.NewError(
		http.StatusBadRequest,
		errx.ValidationErrorCode,
		"invalid input",
		errors.New("bad data"),
	)

	ErrorJSON(rec, nil, apiErr)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	var r resp
	_ = json.Unmarshal(rec.Body.Bytes(), &r)

	if r.Message != "invalid input" {
		t.Errorf("unexpected message: %s", r.Message)
	}

	if r.Code != string(errx.ValidationErrorCode) {
		t.Errorf("unexpected code: %s", r.Code)
	}
}

// --- Test ErrorJSON with generic error ---
func TestErrorJSON_GenericError(t *testing.T) {
	rec := httptest.NewRecorder()

	err := errors.New("something broke")

	ErrorJSON(rec, nil, err)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}

	var r resp
	_ = json.Unmarshal(rec.Body.Bytes(), &r)

	if r.Message != "something broke" {
		t.Errorf("unexpected message")
	}
}

// --- Test ErrorJSON nil error ---
func TestErrorJSON_Nil(t *testing.T) {
	rec := httptest.NewRecorder()

	ErrorJSON(rec, nil, nil)

	if rec.Body.Len() != 0 {
		t.Errorf("expected no response")
	}
}

// --- Test Error (without context) ---
func TestError_WithCustomError(t *testing.T) {
	rec := httptest.NewRecorder()

	apiErr := errx.NewError(
		http.StatusUnauthorized,
		errx.UnauthorizedCode,
		"unauthorized",
		errors.New("no token"),
	)

	Error(rec, apiErr)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}

	var r resp
	_ = json.Unmarshal(rec.Body.Bytes(), &r)

	if r.Code != string(errx.UnauthorizedCode) {
		t.Errorf("unexpected code")
	}
}

// --- Test Error generic ---
func TestError_Generic(t *testing.T) {
	rec := httptest.NewRecorder()

	Error(rec, errors.New("fail"))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500")
	}
}

// --- Test Error nil ---
func TestError_Nil(t *testing.T) {
	rec := httptest.NewRecorder()

	Error(rec, nil)

	if rec.Body.Len() != 0 {
		t.Errorf("expected empty response")
	}
}
