package postgres

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Lapakin/edu-planner/internal/adapter/db"
)

const (
	violatesNotNullConstrainErrCode = "23502"
	foreignKeyErrCode               = "23503"
	uniqueViolationErrCode          = "23505"
	undefinedColumnErrCode          = "42703"
)

var ErrInvalidIsolationLevel = errors.New("invalid isolation level")

func HandleError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		case violatesNotNullConstrainErrCode:
			return &db.ViolatesNotNullConstrainError{
				TableName: pgErr.TableName,
				Message:   pgErr.Message,
				Details:   pgErr.Detail,
			}
		case foreignKeyErrCode:
			return &db.ViolatedForeignKeyError{
				TableName: pgErr.TableName,
				Message:   pgErr.Message,
				Details:   pgErr.Detail,
			}
		case uniqueViolationErrCode:
			return &db.UniqueViolationError{
				TableName: pgErr.TableName,
				Message:   pgErr.Message,
				Details:   pgErr.Detail,
			}
		case undefinedColumnErrCode:
			return &db.UndefinedColumnError{
				TableName: pgErr.TableName,
				Message:   pgErr.Message,
				Details:   pgErr.Detail,
			}
		default:
			return &db.Error{
				TableName: pgErr.TableName,
				Message:   pgErr.Message,
				Details:   pgErr.Detail,
			}
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return db.ErrNoRows
	}

	return nil
}

func IsViolatedForeignKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == foreignKeyErrCode {
			return true
		}
	}

	return false
}

func IsSQLNoRowsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
