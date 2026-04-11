package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nhxgo/go-utils/errx"
)

func PostgresError(message string, err error) errx.Error {
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return errx.NewError(http.StatusConflict, errx.DBUniqueViolation, "Resource already exists", fmt.Errorf("%s: %s", message, pgErr.Detail))

		case "23503": // foreign_key_violation
			return errx.NewError(http.StatusBadRequest, errx.DBForeignKeyViolation, "Invalid reference, foreign key constraint failed", fmt.Errorf("%s: %s", message, pgErr.Message))
		case "23502": // not_null_violation
			return errx.NewError(http.StatusBadRequest, errx.DBNotFound, "Required field missing", fmt.Errorf("%s: %s", message, pgErr.Message))
		case "23514": // check_violation
			return errx.NewError(http.StatusBadRequest, errx.DBError, "Check constraint failed", fmt.Errorf("%s: %s", message, pgErr.Message))

		default: // all other Postgres errors
			return errx.NewError(http.StatusInternalServerError, errx.DBError, "Application error", fmt.Errorf("%s: %s", message, pgErr.Message))
		}
	}

	return errx.InternalServerError(fmt.Errorf("%v :%v", message, err))
}
