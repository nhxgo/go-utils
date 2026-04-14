package errx

import (
	"errors"
	"net/http"
)

type ApiError string

type Error struct {
	HttpStatus int
	ErrorCode  ApiError
	Message    string
	Err        error
}

func (e *Error) Error() string {
	return "[" + string(e.ErrorCode) + "] " + e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

// ------------------------utils

func NewError(httpStatus int, errorCode ApiError, message string, err error) error {
	return &Error{
		HttpStatus: httpStatus,
		ErrorCode:  errorCode,
		Message:    message,
		Err:        err,
	}
}

// HTTP helpers

func InternalServerError(err error) error {
	return &Error{
		HttpStatus: http.StatusInternalServerError,
		ErrorCode:  InternalServerErrorCode,
		Message:    "internal server error",
		Err:        err,
	}
}

func ConflictError(err error, message string) error {
	return &Error{
		HttpStatus: http.StatusConflict,
		ErrorCode:  ConflictCode,
		Message:    message,
		Err:        err,
	}
}

func NotFoundError(err error, message string) error {
	return &Error{
		HttpStatus: http.StatusNotFound,
		ErrorCode:  NotFoundCode,
		Message:    message,
		Err:        err,
	}
}

func MethodNotAllowedError(err error) error {
	return &Error{
		HttpStatus: http.StatusMethodNotAllowed,
		ErrorCode:  MethodNotAllowedCode,
		Message:    "method not allowed",
		Err:        err,
	}
}

func BadRequestError(err error, message string) error {
	return &Error{
		HttpStatus: http.StatusBadRequest,
		ErrorCode:  BadRequestCode,
		Message:    message,
		Err:        err,
	}
}

func UnauthorizedError(err error, message string) error {
	return &Error{
		HttpStatus: http.StatusUnauthorized,
		ErrorCode:  UnauthorizedCode,
		Message:    message,
		Err:        err,
	}
}

func TooManyRequestError(err error) error {
	return &Error{
		HttpStatus: http.StatusTooManyRequests,
		ErrorCode:  TooManyRequestsCode,
		Message:    "too many requests, please try again later",
		Err:        err,
	}
}

// Crypto / hashing

func NewDecryptError(err error) error {
	return &Error{
		HttpStatus: http.StatusInternalServerError,
		ErrorCode:  InternalServerErrorCode,
		Message:    "failed to decrypt data",
		Err:        errors.Join(ErrDecryptFailed, err),
	}
}

func NewEncryptError(err error) error {
	return &Error{
		HttpStatus: http.StatusInternalServerError,
		ErrorCode:  InternalServerErrorCode,
		Message:    "failed to encrypt data",
		Err:        errors.Join(ErrEncryptFailed, err),
	}
}

func NewHashingError(err error) error {
	return &Error{
		HttpStatus: http.StatusInternalServerError,
		ErrorCode:  InternalServerErrorCode,
		Message:    "failed to hash data",
		Err:        errors.Join(ErrHashingFailed, err),
	}
}
