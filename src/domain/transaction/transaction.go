package transaction_entity // Adjust package name as needed (e.g., internal/models)

import (
    "time"
	"database/sql"
)

// Transaction represents the transactions table in the database.
type TransactionEntity struct {
    ID           int       `json:"id" db:"id"`
    AccountID    int       `json:"account_id" db:"account_id"`
    Type         string    `json:"type" db:"type"`
    Amount       float64   `json:"amount" db:"amount"`
    ToAccountID  sql.NullInt32      `json:"to_account_id" db:"to_account_id"` // Nullable
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}


func ScanTransactionEntity(r *sql.Rows, transaction *TransactionEntity) error {
    return r.Scan(
			&transaction.ID,
			&transaction.Amount,
			&transaction.AccountID,
			&transaction.ToAccountID,
			&transaction.Type,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
}

func FetchTransactionEntities(r *sql.Rows) ([]TransactionEntity,error) {
    var transactions []TransactionEntity
    defer r.Close()
    for r.Next(){
        var transaction TransactionEntity
        er := ScanTransactionEntity(r,&transaction)
		if er != nil {
			return nil,er
		}
		transactions = append(transactions,transaction)
    }

    return transactions,nil
}