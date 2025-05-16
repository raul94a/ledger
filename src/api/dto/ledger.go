package clientdto

import (
	"time"
)


type LedgerEntryDto struct {
    ID           int       `json:"id"`
    TransactionID int      `json:"transaction_id"`
    AccountID    int       `json:"account_id"`
    Type         string    `json:"type"` // credit, debit
    Amount       float64   `json:"amount"`
    CreatedAt    time.Time `json:"created_at"`
}