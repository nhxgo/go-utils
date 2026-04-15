package request

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- Test Struct ---
type TestRequest struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"gte=1"`
}

func TestGetBody_Success(t *testing.T) {
	body := `{"name":"ajay","age":25}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))

	var result TestRequest
	err := GetBody(req, &result)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Name != "ajay" || result.Age != 25 {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestGetBody_InvalidJSON(t *testing.T) {
	body := `{"name":}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))

	var result TestRequest
	err := GetBody(req, &result)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetBody_UnknownField(t *testing.T) {
	body := `{"name":"ajay","age":25,"extra":"field"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))

	var result TestRequest
	err := GetBody(req, &result)

	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestGetBody_ValidationError(t *testing.T) {
	body := `{"name":"","age":0}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))

	var result TestRequest
	err := GetBody(req, &result)

	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestGetBody_NoBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)

	var result TestRequest
	err := GetBody(req, &result)

	if err == nil {
		t.Fatal("expected error for nil body")
	}
}
func TestGetBodyAndValidate_Success(t *testing.T) {
	body := `{"name":"ajay","age":25}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))

	var result TestRequest
	err := GetBodyAndValidate(req, &result)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWithLoggerAndFromContext(t *testing.T) {
	ctx := context.Background()

	reqCtx := &Context{
		Context: ctx,
		Logger:  nil, // logger.Default() will be used later
	}

	newCtx := WithLogger(ctx, reqCtx)

	result := FromContext(newCtx)

	if result == nil {
		t.Fatal("expected context, got nil")
	}
}

func TestGetContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ctx := GetContext(req)

	if ctx == nil {
		t.Fatal("expected context, got nil")
	}
}

func TestSetAndGetUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	type User struct {
		ID int
	}

	user := User{ID: 1}

	req = SetUser(req, user)

	result, ok := GetUser[User](req)

	if !ok {
		t.Fatal("expected user to be present")
	}

	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, ok := GetUser[string](req)

	if ok {
		t.Fatal("expected no user")
	}
}

func TestValidate_Success(t *testing.T) {
	data := TestRequest{
		Name: "ajay",
		Age:  25,
	}

	err := Validate(data)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_Failure(t *testing.T) {
	data := TestRequest{
		Name: "",
		Age:  0,
	}

	err := Validate(data)

	if err == nil {
		t.Fatal("expected validation error")
	}
}
