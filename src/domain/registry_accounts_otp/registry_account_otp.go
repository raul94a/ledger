package registry_account_otp_entity


import (
	"time"
)

// RegisterAccountsOTP represents the public.register_accounts_otp table
type RegisterAccountsOTP struct {
	ID        int        `json:"id" db:"id"`
	ClientID  int       `json:"client_id" db:"client_id"` // Using *int for nullable integer
	OTP       string    `json:"otp" db:"otp"`             // Using *string for nullable varchar
	Validated bool       `json:"validated" db:"validated"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}