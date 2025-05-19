package accountentity

import (
    "time"
)

// AccountBalanceMV represents the account_balances_mv materialized view in the database.
type AccountBalanceMV struct {
    AccountID int       `json:"account_id" db:"account_id"`
    Balance   float64   `json:"balance" db:"balance"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}