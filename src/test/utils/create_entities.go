package utils

import (
	"database/sql"
	accountentity "src/domain/account"
	cliententity "src/domain/client"
	ledgerentity "src/domain/ledger"
	transaction_entity "src/domain/transaction"
	"time"

	"github.com/google/uuid"
)

func CreateClientTest(ID int, name string, email string) cliententity.ClientEntity {
	return cliententity.ClientEntity{
		ID:             ID,
		Name:           name,
		Surname1:       "Doe",
		Surname2:       sql.NullString{String: "Smith", Valid: true}, // Optional, can be {Valid: false} for NULL
		Email:          email,
		Identification: "ABC123456" + string(ID),
		Nationality:    "US",
		DateOfBirth:    time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
		Sex:            "M",
		Address:        "123 Main St",
		City:           "Springfield",
		Province:       "Illinois",
		State:          sql.NullString{String: "IL", Valid: true}, // Optional, can be {Valid: false} for NULL
		ZipCode:        "62701",
		Telephone:      "+1-555-123-4567",
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
}

func CreateAccount(clientId int) accountentity.AccountEntity {
	return accountentity.AccountEntity{
		ClientID:      clientId,
		AccountNumber: uuid.New().String(), // Random UUID
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func CreateTransaction(accountID int, toAccountID int, amount float64, typ string) transaction_entity.TransactionEntity {
	toAccountIDPtr := &toAccountID // Pointer for nullable field
	return transaction_entity.TransactionEntity{
		AccountID:   accountID,
		Type:        typ,
		Amount:      amount,
		ToAccountID: toAccountIDPtr,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// InitializeLedgerEntry creates a new LedgerEntry instance.
func InitializeLedgerEntry(transactionID, accountID int, entryType string, amount float64) ledgerentity.LedgerEntryEntity {
	return ledgerentity.LedgerEntryEntity{
		TransactionID: transactionID,
		AccountID:     accountID,
		Type:          entryType, // "credit" or "debit"
		Amount:        amount,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
