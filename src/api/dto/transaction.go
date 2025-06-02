package clientdto

import (
	"time"
)



type PerformTransactionDto struct {
    AccountID   int       `json:"account_id"`
    Type        string    `json:"type"` // ADD, WITHDRAWAL, TRANSFER
    Amount      float64   `json:"amount"`
    ToAccountNumber *string       `json:"to_account_number,omitempty"` // For transfers
}

type TransactionDto struct {
    ID          int       `json:"id"`
    AccountID   int       `json:"account_id"`
    Type        string    `json:"type"` // ADD, WITHDRAWAL, TRANSFER
    Amount      float64   `json:"amount"`
    ToAccountNumber *string `json:"to_account_number"` // For transfers
    CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}