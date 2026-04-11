package errx

import (
	"errors"
	"net/http"
)

type ApiError string

type errorImpl struct {
	HttpStatus  int
	ErrorCode   ApiError
	Description error
	Message     string
}
type Error *errorImpl

var (
	DecryptError       = errors.New("Failed to Decrypt")
	EncryptError       = errors.New("Failed to Encrypt")
	WrongPasswordError = errors.New("Wrong password")
)

const (
	UnknownError ApiError = "UNKNOWN_ERROR"
	Empty        ApiError = ""

	internalServer ApiError = "INTERNAL_SERVER_ERROR"
	notFound       ApiError = "NOT_FOUND"
	unauthorized   ApiError = "UNAUTHORIZED"
	badRequest     ApiError = "BAD_REQUEST"
	conflict       ApiError = "CONFLICT"

	UnprocessableEntity ApiError = "UNPROCESSABLE_ENTITY"
	ValidationError     ApiError = "VALIDATION_ERROR"
	methodNotAllowed    ApiError = "METHOD_NOT_ALLOWED"

	// Database
	DBError               ApiError = "DATABASE_ERROR"
	DBUniqueViolation     ApiError = "DB_UNIQUE_VIOLATION"
	DBForeignKeyViolation ApiError = "DB_FOREIGN_KEY_VIOLATION"
	DBNotFound            ApiError = "DB_NOT_FOUND"
	DBTimeout             ApiError = "DB_TIMEOUT"
	DBConnectionFailed    ApiError = "DB_CONNECTION_FAILED"

	UserNotFound ApiError = "USER_NOT_FOUND"

	// Business Logic
	InvalidCredentials ApiError = "INVALID_CREDENTIALS"
	UserAlreadyExists  ApiError = "USER_ALREADY_EXISTS"
	PermissionDenied   ApiError = "PERMISSION_DENIED"
)

func NewError(HttpStatus int, errorCode ApiError, message string, description error) Error {
	return &errorImpl{
		HttpStatus:  HttpStatus,
		ErrorCode:   errorCode,
		Description: description,
		Message:     message,
	}
}

func InternalServerError(description error) Error {
	return &errorImpl{
		HttpStatus:  http.StatusInternalServerError,
		ErrorCode:   internalServer,
		Description: description,
		Message:     "Internal Server Error",
	}
}

func ConflictError(description error, message string) Error {
	return &errorImpl{
		HttpStatus:  http.StatusConflict,
		ErrorCode:   conflict,
		Description: description,
		Message:     message,
	}
}

func NotFoundError(description error, message string) Error {
	return &errorImpl{
		HttpStatus:  http.StatusNotFound,
		ErrorCode:   notFound,
		Description: description,
		Message:     message,
	}
}

func MethodNotAllowedError(description error) Error {
	return &errorImpl{
		HttpStatus:  http.StatusMethodNotAllowed,
		ErrorCode:   methodNotAllowed,
		Description: description,
		Message:     "Method Not Allowed",
	}
}

func BadRequestError(description error, message string) Error {
	return &errorImpl{
		HttpStatus:  http.StatusBadRequest,
		ErrorCode:   badRequest,
		Description: description,
		Message:     message,
	}
}

func UnauthorizedError(description error, message string) Error {
	return &errorImpl{
		HttpStatus:  http.StatusUnauthorized,
		ErrorCode:   unauthorized,
		Description: description,
		Message:     message,
	}
}
