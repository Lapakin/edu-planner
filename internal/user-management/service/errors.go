package service

import (
	"errors"
	"fmt"

	"github.com/Lapakin/edu-planner/internal/adapter/db"

	pg "github.com/Lapakin/edu-planner/internal/adapter/db/postgres"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrUnmarshal = errors.New("unmarshal error")
	ErrInternal  = errors.New("internal server error")
)

func handleDBError(err error) error {
	if dbErr := pg.HandleError(err); dbErr != nil {
		if uniqueErr, ok := errors.AsType[*db.UniqueViolationError](dbErr); ok {
			return NewDuplicateUniqueValueError(uniqueErr)
		}
		if fkErr, ok := errors.AsType[*db.ViolatedForeignKeyError](dbErr); ok {
			return NewViolatedForeignKeyValueError(fkErr)
		}
		if errors.Is(dbErr, db.ErrNoRows) {
			return ErrNotFound
		}
	}
	return ErrInternal
}

type DuplicateUniqueValueError struct {
	Table   string `json:"table"`
	Details string `json:"details"`
}

func NewDuplicateUniqueValueError(err *db.UniqueViolationError) *DuplicateUniqueValueError {
	return &DuplicateUniqueValueError{
		Table:   err.TableName,
		Details: err.Details,
	}
}

type ViolatedForeignKeyValueError struct {
	Table   string `json:"table"`
	Details string `json:"details"`
}

func NewViolatedForeignKeyValueError(err *db.ViolatedForeignKeyError) *ViolatedForeignKeyValueError {
	return &ViolatedForeignKeyValueError{
		Table:   err.TableName,
		Details: err.Details,
	}
}

func (e *DuplicateUniqueValueError) Error() string {
	return fmt.Sprintf("Duplicate unique value error in table %s: %s", e.Table, e.Details)
}

func (e *ViolatedForeignKeyValueError) Error() string {
	return fmt.Sprintf("Foreign key violation value error in table %s: %s", e.Table, e.Details)
}
