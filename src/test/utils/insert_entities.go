package utils

import (
	"context"
	"database/sql"
	accountentity "src/domain/account"
	cliententity "src/domain/client"
	transaction_entity "src/domain/transaction"
	"fmt"

)


func InsertClient(ctx context.Context, db *sql.DB, client *cliententity.ClientEntity) error {
	query := `
        INSERT INTO clients (
            name, surname1, surname2, email, identification, nationality, 
            date_of_birth, sex, address, city, province, state, 
            zip_code, telephone
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
        RETURNING id, created_at, updated_at`

	// Execute the query and scan the returned values into the client struct
	err := db.QueryRowContext(ctx, query,
		client.Name,
		client.Surname1,
		client.Surname2,
		client.Email,
		client.Identification,
		client.Nationality,
		client.DateOfBirth,
		client.Sex,
		client.Address,
		client.City,
		client.Province,
		client.State,
		client.ZipCode,
		client.Telephone,
	).Scan(&client.ID, &client.CreatedAt, &client.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert client: %w", err)
	}
	return nil
}

func InsertAccount(ctx context.Context, db *sql.DB, account *accountentity.AccountEntity, clientID int) error {
	query := `
        INSERT INTO accounts (
            client_id, account_number
        ) VALUES ($1, $2)
        RETURNING id, created_at, updated_at`

	// Execute the query and scan the returned values into the client struct
	err := db.QueryRowContext(ctx, query,
		clientID,
		account.AccountNumber,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert account: %w", err)
	}
	return nil
}

func InsertTransaction(ctx context.Context, tx *sql.Tx, transaction *transaction_entity.TransactionEntity, clientID int) error {
	query := `
        INSERT INTO transactions (
            client_id, account_id, to_account_id, amount, type
        ) VALUES ($1, $2)
        RETURNING id, created_at, updated_at`

	// Execute the query and scan the returned values into the client struct
	err := tx.QueryRowContext(ctx, query,
		clientID,
		transaction.AccountID,
		transaction.ToAccountID,
		transaction.Amount,
		transaction.Type,
	).Scan(&transaction.ID, &transaction.CreatedAt, &transaction.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}
	return nil
}

func InsertLedgerEntry(ctx context.Context, tx *sql.Tx, transaction *transaction_entity.TransactionEntity, accountID int, typeTransaction string)error{
	query := `
        INSERT INTO ledger_entries (
            transaction_id, account_id, type, amount
        ) VALUES ($1, $2)
        RETURNING id, created_at, updated_at`

	// Execute the query and scan the returned values into the client struct
	err := tx.QueryRowContext(ctx, query,
		transaction.ID,
		accountID,
		typeTransaction,
		transaction.Amount,
	).Scan(&transaction.ID, &transaction.CreatedAt, &transaction.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert ledger_entry: %w", err)
	}
	return nil
}

func AccountTransactionTx(ctx context.Context, db *sql.DB, transaction *transaction_entity.TransactionEntity, clientID int, typeFrom string, typeTo string) error {
	options := sql.TxOptions{
		Isolation: 0,
		ReadOnly: false,
	}
	tx,error := db.BeginTx(ctx, &options)

	if error != nil {
		fmt.Println("Error beginning transaction: " + error.Error())
		return error
	}

	insertTransactionError := InsertTransaction(ctx,tx,transaction, clientID)
	if insertTransactionError != nil{
		error = tx.Rollback()
		fmt.Println("Error inserting account transaction: " + error.Error())
		return nil
	}

	insertLedgerEntryError := InsertLedgerEntry(ctx,tx,transaction, transaction.AccountID, typeFrom)
	if insertLedgerEntryError != nil{
		error = tx.Rollback()
		fmt.Println("Error inserting ledger Entry transaction (from): " + error.Error())
		return nil
	}

	insertLedgerEntryError = InsertLedgerEntry(ctx,tx,transaction, *transaction.ToAccountID, typeTo)
	if insertLedgerEntryError != nil{
		error = tx.Rollback()
		fmt.Println("Error inserting ledger Entry transaction (to): " + error.Error())
		return nil
	}


	tx.Commit()
	return nil
	
}