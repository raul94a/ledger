package repositories

import "fmt"


type ErrNotEnoughFunds struct {
    Message string
}

type ErrEntityNotFound struct {
    Identifier any
}

type ErrTransactionInsertFailed struct {
    Message string
    Reason  error // Wraps original DB error
}

type ErrLedgerEntryInsertFailed struct {
    Message string
    Reason  error // Wraps original DB error
}

type ErrNoRowsAffected struct {
    Message string
    Reason  error // Wraps original DB error
}

func (e *ErrEntityNotFound) Error() string {
    return fmt.Sprintf("Entity not found: %v", e.Identifier)
}

func (e *ErrTransactionInsertFailed) Error() string {
    return fmt.Sprintf("Failed transaction insertion: %s", e.Message)
}

func (e *ErrTransactionInsertFailed) Unwrap() error {
    return e.Reason
}

func (e *ErrLedgerEntryInsertFailed) Error() string {
    return fmt.Sprintf("Failed ledger entry insertion: %s", e.Message)
}

func (e *ErrLedgerEntryInsertFailed) Unwrap() error {
    return e.Reason
}

func (e *ErrNoRowsAffected) Error() string {
    return fmt.Sprintf("No rows affected: %s", e.Message)
}

func (e *ErrNoRowsAffected) Unwrap() error {
    return e.Reason
}

func (e *ErrNotEnoughFunds) Error() string {
    return fmt.Sprintf(e.Message)
}