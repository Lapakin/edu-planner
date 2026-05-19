package db

import (
	"errors"
	"fmt"
)

type ViolatesNotNullConstrainError struct {
	TableName string `json:"table_name"`
	Message   string `json:"message"`
	Details   string `json:"details"`
}

func (e *ViolatesNotNullConstrainError) Error() string {
	return fmt.Sprintf("ViolatesNotNullConstrainError: Table '%s' has a NOT NULL constraint violation. Message: %s. Details: %s", e.TableName, e.Message, e.Details)
}

type ViolatedForeignKeyError struct {
	TableName string `json:"table_name"`
	Message   string `json:"message"`
	Details   string `json:"details"`
}

func (e *ViolatedForeignKeyError) Error() string {
	return fmt.Sprintf("ViolatedForeignKeyError: Table '%s' has a foreign key violation. Message: %s. Details: %s", e.TableName, e.Message, e.Details)
}

type UniqueViolationError struct {
	TableName string `json:"table_name"`
	Message   string `json:"message"`
	Details   string `json:"details"`
}

func (e *UniqueViolationError) Error() string {
	return fmt.Sprintf("UniqueViolationError: Table '%s' has a unique constraint violation. Message: %s. Details: %s", e.TableName, e.Message, e.Details)
}

type UndefinedColumnError struct {
	TableName string `json:"table_name"`
	Message   string `json:"message"`
	Details   string `json:"details"`
}

func (e *UndefinedColumnError) Error() string {
	return fmt.Sprintf("UndefinedColumnError: Table '%s' has an undefined column. Message: %s. Details: %s", e.TableName, e.Message, e.Details)
}

type Error struct {
	TableName string `json:"table_name"`
	Message   string `json:"message"`
	Details   string `json:"details"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Database error in table %s: %s - %s", e.TableName, e.Message, e.Details)
}

var ErrNoRows = errors.New("error no rows")
