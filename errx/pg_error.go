package errx

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

func MapPgError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique violation
			return &Error{
				HttpStatus: http.StatusConflict,
				ErrorCode:  ConflictCode,
				Message:    "duplicate record",
				Err:        err,
			}

		case "23503": // foreign key violation
			return &Error{
				HttpStatus: http.StatusBadRequest,
				ErrorCode:  BadRequestCode,
				Message:    "invalid reference (foreign key violation)",
				Err:        err,
			}

		case "23502": // not null violation
			return &Error{
				HttpStatus: http.StatusBadRequest,
				ErrorCode:  ValidationErrorCode,
				Message:    "missing required field",
				Err:        err,
			}

		case "42P01": // undefined table
			return &Error{
				HttpStatus: http.StatusInternalServerError,
				ErrorCode:  InternalServerErrorCode,
				Message:    "database schema error",
				Err:        err,
			}

		case "40001": // serialization failure (retryable)
			return &Error{
				HttpStatus: http.StatusInternalServerError,
				ErrorCode:  DBTimeoutCode,
				Message:    "transaction conflict, please retry",
				Err:        err,
			}

		default:
			return &Error{
				HttpStatus: http.StatusInternalServerError,
				ErrorCode:  DBErrorCode,
				Message:    "database error",
				Err:        err,
			}
		}
	}

	var apiErr *Error
	if errors.As(err, &apiErr) {
		return NewError(apiErr.HttpStatus, apiErr.ErrorCode, apiErr.Message, apiErr.Err)
	} else {
		return InternalServerError(fmt.Errorf("%v :%v", "internal server error", err))
	}
}
