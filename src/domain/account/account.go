package accountentity

import (
    "time"
	"database/sql"
)

// Account represents the accounts table in the database.
type AccountEntity struct {
    ID           int       `json:"id" db:"id"`
    ClientID     int       `json:"client_id" db:"client_id"`
    AccountNumber string    `json:"account_number" db:"account_number"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}


func ScanAccountEntity(r *sql.Rows, account *AccountEntity) error {
    return r.Scan(
			&account.ID,
			&account.AccountNumber,
			&account.ClientID,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
}

func FetchAccountEntities(r *sql.Rows) ([]AccountEntity,error) {
    var accounts []AccountEntity
    defer r.Close()
    for r.Next(){
        var account AccountEntity
        er := ScanAccountEntity(r,&account)
		if er != nil {
			return nil,er
		}
		accounts = append(accounts,account)
    }

    return accounts,nil
}

