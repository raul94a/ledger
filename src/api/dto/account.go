package clientdto

import (
	"time"
)


type AccountDto struct {
    ID            int       `json:"id"`
    UserID        string    `json:"user_id"` // From Keycloak
    AccountNumber string    `json:"account_number"`
    CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Balance		  float64   `json:"balance"`
}