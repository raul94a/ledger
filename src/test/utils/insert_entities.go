package utils

import (
	"context"
	"database/sql"
	"fmt"
	accountentity "src/domain/account"
	cliententity "src/domain/client"
	transaction_entity "src/domain/transaction"
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
		RETURNING id`

	// Execute the query and scan the returned values into the client struct
	err := db.QueryRowContext(ctx, query,
		clientID,
		account.AccountNumber,
	).Scan(&account.ID)
	if err != nil {
		return fmt.Errorf("failed to insert account: %w", err)
	}


	return nil
}

func InsertTransaction(ctx context.Context, tx *sql.Tx, transaction *transaction_entity.TransactionEntity) error {
	query := `
        INSERT INTO transactions (
            account_id, to_account_id, amount, type
        ) VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at`

	// Execute the query and scan the returned values into the client struct
	
	err := tx.QueryRowContext(ctx, query,
		transaction.AccountID,
		transaction.ToAccountID,
		transaction.Amount,
		transaction.Type,
	).Scan(&transaction.ID, &transaction.CreatedAt, &transaction.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err.Error())
	}
	return nil
}


func InsertLedgerEntry(ctx context.Context, tx *sql.Tx, transaction *transaction_entity.TransactionEntity, accountID int, typeTransaction string) error {
	query := `
        INSERT INTO ledger_entries (
            transaction_id, account_id, type, amount
        ) VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at`

	fmt.Println("Insertando ledger entry con el siguiente obj de transacción")
	fmt.Println("TRANSACTION ID " + fmt.Sprint(transaction.ID))
	fmt.Println("ACCOUNT ID " + fmt.Sprint(transaction.AccountID))
	fmt.Println("TYPE " + transaction.Type)
	fmt.Println("AMOUNT " + fmt.Sprint(transaction.Amount))

	// Execute the query and scan the returned values into the client struct
	_,err := tx.ExecContext(ctx, query,
		transaction.ID,
		accountID,
		typeTransaction,
		transaction.Amount,
	)

	if err != nil {
		fmt.Println("failed to insert ledger_entry ", transaction)
		return fmt.Errorf("");
	}
	return nil
}

func AccountTransactionTx(ctx context.Context, db *sql.DB, transaction *transaction_entity.TransactionEntity, clientID int, typeFrom string, typeTo string) error {
	options := sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	}
	tx, error := db.BeginTx(ctx, &options)
	
	if error != nil {
		fmt.Println("Error beginning transaction: " + error.Error())
		return error
	}

	insertTransactionError := InsertTransaction(ctx, tx, transaction)
	fmt.Println(insertTransactionError)
	if insertTransactionError != nil {
	    tx.Rollback()
		fmt.Println("Error inserting account transaction: " + insertTransactionError.Error())
		return insertTransactionError
	}

	insertLedgerEntryError := InsertLedgerEntry(ctx, tx, transaction, transaction.AccountID, typeFrom)
	
	if insertLedgerEntryError != nil {
		tx.Rollback()
		fmt.Println("Error inserting ledger Entry transaction (from): " + insertLedgerEntryError.Error())
		fmt.Println(transaction)
		return insertLedgerEntryError
	}
	if transaction.ToAccountID.Valid {
		value := transaction.ToAccountID.Int32
		insertLedgerEntryError = InsertLedgerEntry(ctx, tx, transaction, int(value), typeTo)
		if insertLedgerEntryError != nil {
			tx.Rollback()
			fmt.Println("Error inserting ledger Entry transaction (to): " + insertLedgerEntryError.Error())
			return insertLedgerEntryError
		}
	}

	tx.Commit()
	return nil

}
