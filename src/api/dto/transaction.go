package clientdto

import (
	"time"
)



type Transaction struct {
    ID          int       `json:"id"`
    AccountID   int       `json:"account_id"`
    Type        string    `json:"type"` // deposit, withdraw, transfer
    Amount      float64   `json:"amount"`
    ToAccountID int      `json:"to_account_id,omitempty"` // For transfers
    CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}