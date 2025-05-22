package clientdto

import (
	
)

type AccountDto struct {
	ID            int       `json:"id"`
	ClientID      int    `json:"client_id"` // From Keycloak
	AccountNumber string    `json:"account_number"`
	Balance       float64   `json:"balance"`
	CreatedDate    string `json:"created_date" binding:"required,datetime=2006-01-02 15:04:05"` // ISO 8601 date (YYYY-MM-DD HH:mm:ss)
    UpdatedDate    string `json:"updated_date" binding:"required,datetime=2006-01-02 15:04:05"` // ISO 8601 date (YYYY-MM-DD HH:mm:ss)
}


type CreateAccountRequest struct {
	ClientID int `json:"client_id"`
}


