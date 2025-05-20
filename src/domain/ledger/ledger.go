package ledgerentity

import (
	transaction_entity "src/domain/transaction"
	"database/sql"
	"time"
)

// LedgerEntryEntity represents the ledger_entries table in the database.
type LedgerEntryEntity struct {
	ID            int       `json:"id" db:"id"`
	TransactionID int       `json:"transaction_id" db:"transaction_id"`
	AccountID     int       `json:"account_id" db:"account_id"`
	Type          string    `json:"type" db:"type"`
	Amount        float64   `json:"amount" db:"amount"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type LedgerTransaction struct {
	Transaction transaction_entity.TransactionEntity
	AccountID int
	TransactionType string // CREDIT / DEBIT

}

func ScanLedgerEntryEntity(r *sql.Rows, ledgerEntry *LedgerEntryEntity) error {
	return r.Scan(
		&ledgerEntry.ID,
		&ledgerEntry.Amount,
		&ledgerEntry.TransactionID,
		&ledgerEntry.AccountID,
		&ledgerEntry.Type,
		&ledgerEntry.CreatedAt,
		&ledgerEntry.UpdatedAt,
	)
}

func FetchLedgerEntities(r *sql.Rows) ([]LedgerEntryEntity, error) {
	var ledgerEntries []LedgerEntryEntity
	defer r.Close()
	for r.Next() {
		var ledgerEntry LedgerEntryEntity
		er := ScanLedgerEntryEntity(r, &ledgerEntry)
		if er != nil {
			return nil, er
		}
		ledgerEntries = append(ledgerEntries, ledgerEntry)
	}

	return ledgerEntries, nil
}
