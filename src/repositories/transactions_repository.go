package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	ledgerentity "src/domain/ledger"
	transaction_entity "src/domain/transaction"
	errors "src/errors"
	validators "src/validators"
	"strings"
)

type TransactionRepository interface {
	// FetchTransactionById(ctx context.Context, ID int) (transaction_entity.TransactionEntity, error)
	// FetchTransactionsByAccount(ctx context.Context, accountID int) ([]transaction_entity.TransactionEntity, error)
	FetchAccountBalance(ctx context.Context, tx *sql.Tx, accountID int) (*float64, errors.AppError)
	updateAccountBalance(ctx context.Context, tx *sql.Tx, ledgerTransaction ledgerentity.LedgerTransaction) errors.AppError
	InsertTransaction(ctx context.Context, tx *sql.Tx, transaction *transaction_entity.TransactionEntity) errors.AppError
	InsertLedgerEntry(ctx context.Context, tx *sql.Tx, ledgerTransaction *ledgerentity.LedgerTransaction) errors.AppError
	InsertTransactionLedgerTx(ctx context.Context, transaction *transaction_entity.TransactionEntity) errors.AppError
}

type transactionRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewTransactionRepository(db *sql.DB, logger *slog.Logger) TransactionRepository {
	if db == nil {
		panic("db cannot be nil")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &transactionRepository{db: db, logger: logger}
}

func (r *transactionRepository) InsertTransaction(ctx context.Context, tx *sql.Tx, transaction *transaction_entity.TransactionEntity) errors.AppError {
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
		r.logger.Error("Error occurred while inserting transaction: %s", err.Error())
		r.logger.Error("Account id: %s, amount: %v, type: %s", transaction.AccountID, transaction.Amount, transaction.Type)
		return &errors.ErrInternalServer{Reason: err}
	}
	return nil
}

func (r *transactionRepository) InsertLedgerEntry(ctx context.Context, tx *sql.Tx, ledgerTransaction *ledgerentity.LedgerTransaction) errors.AppError {
	query := `
        INSERT INTO ledger_entries (
            transaction_id, account_id, type, amount
        ) VALUES ($1, $2, $3, $4)`

	transaction := ledgerTransaction.Transaction
	// Execute the query and scan the returned values into the client struct
	result, err := tx.ExecContext(ctx, query,
		transaction.ID,
		ledgerTransaction.AccountID,
		ledgerTransaction.LedgerType,
		transaction.Amount,
	)

	if err != nil {
		isFromTransaction := ledgerTransaction.Transaction.AccountID == ledgerTransaction.AccountID
		origin := ""
		if isFromTransaction {
			origin = "(from)"
		} else {
			origin = "(to)"
		}
		errString := fmt.Sprintf("Error %s occurred in Txn while inserting ledger entry (%s)", err.Error(), origin)
		r.logger.Error(errString)
		errString = fmt.Sprintf("Account id: %d, amount: %v, type: %s", transaction.AccountID, transaction.Amount, transaction.Type)
		r.logger.Error(errString)
		return &errors.ErrInternalServer{Reason: err}
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		errString := fmt.Sprintf("No rows affected for transaction %v", ledgerTransaction)
		r.logger.Error(errString)
		return &errors.ErrInternalServer{Reason: err}
	}
	err = r.updateAccountBalance(ctx, tx, *ledgerTransaction)
	if err != nil {
		return &errors.ErrInternalServer{Reason: err}
	}
	return nil
}

/**
* Database transaction to move funds: ADD, WITHDRAWAL OR TRANSFER
* 1. Initialize database transaction (Tx)
* 2. Check balances if TransactionType is WITHDRAWAL OR TRANSFER
* 3. Insert the transaction —Money exchange— into the database
* 4. Insert the LedgerEntry and update the balance for the source account
* 5. Insert the LedgerEntry and update the balance of the destination account
*
*
 */

func (r *transactionRepository) InsertTransactionLedgerTx(ctx context.Context, transaction *transaction_entity.TransactionEntity) errors.AppError {
	tx, txErr := r.db.BeginTx(ctx, &sql.TxOptions{
		ReadOnly:  false,
		Isolation: 0,
	})

	if txErr != nil {
		r.logger.Error("Error occurred while beginning transaction ledger: %s ", txErr.Error())
		return &errors.ErrInternalServer{Reason: txErr}
	}

	// we always starts with the account triggering the transaction
	ledgerType := ""

	balance, err := r.FetchAccountBalance(ctx, tx, transaction.AccountID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = validators.ValidateTransactionBalance(*transaction, *balance, r.logger)

	if err != nil {
		tx.Rollback()
		return err
	}

	err = r.InsertTransaction(ctx, tx, transaction)

	if err != nil {
		tx.Rollback()
		return err
	}

	if strings.ToUpper(transaction.Type) != "ADD" {
		ledgerType = "DEBIT"
	} else {
		ledgerType = "CREDIT"
	}

	transactionLedger := ledgerentity.LedgerTransaction{
		Transaction: *transaction,
		LedgerType:  ledgerType,
		AccountID:   transaction.AccountID,
	}

	err = r.InsertLedgerEntry(ctx, tx, &transactionLedger)
	if err != nil {
		tx.Rollback()
		return err
	}

	if !transaction.ToAccountID.Valid {
		warning := fmt.Sprintf(
			"No attached ToAccountId in transaction %d. Tx finishing without errors for account_id %d",
			transaction.ID,
			transaction.AccountID,
		)
		r.logger.Warn(warning)
		tx.Commit()
		return nil
	}
	// If there's a valid ToAccountId, another ledger entry has to be inserted

	if ledgerType == "CREDIT" {
		ledgerType = "DEBIT"
	} else {
		ledgerType = "CREDIT"
	}
	transactionLedger.AccountID = int(transaction.ToAccountID.Int32)
	transactionLedger.LedgerType = ledgerType

	err = r.InsertLedgerEntry(ctx, tx, &transactionLedger)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil

}

func (r *transactionRepository) updateAccountBalance(
	ctx context.Context,
	tx *sql.Tx,
	ledgerTransaction ledgerentity.LedgerTransaction,
) errors.AppError {
	query := ""
	mustAdd := ledgerTransaction.LedgerType == "CREDIT"
	infoStr := fmt.Sprint("Must add credit to ledger: %v", mustAdd)
	r.logger.Info(infoStr)
	if mustAdd {
		query = `UPDATE account_balances SET balance = balance + $1 where account_id = $2`

	} else {
		query = `UPDATE account_balances SET balance = balance - $1 where account_id = $2`
	}

	_, err := tx.ExecContext(ctx, query, ledgerTransaction.Transaction.Amount, ledgerTransaction.AccountID)
	if err != nil {
		errStr := fmt.Sprintf(
			"Error occurred in Txn while updating account balance. ACTION: %s, AMOUNT: %f, ACCOUNT_ID: %d",
			ledgerTransaction.LedgerType,
			ledgerTransaction.Transaction.Amount,
			ledgerTransaction.AccountID,
		)
		r.logger.Error(errStr)
		return &errors.ErrInternalServer{Reason: err}
	}

	return nil

}

func (r *transactionRepository) FetchAccountBalance(ctx context.Context, tx *sql.Tx, accountID int) (*float64, errors.AppError) {
	query := `SELECT balance from account_balances where account_id = $1`
	var balance *float64
	if tx == nil {
		err := r.db.QueryRowContext(ctx, query, accountID).Scan(&balance)
		if err != nil {
			errStr := fmt.Sprintf(
				"Error occurred while fetching account balance. ACCOUNT_ID: %v",
				accountID,
			)
			r.logger.Error(errStr)
			return nil, &errors.ErrInternalServer{Reason: err}

		}
	} else {
		err := tx.QueryRowContext(ctx, query, accountID).Scan(&balance)
		if err != nil {
			errStr := fmt.Sprintf(
				"Error occurred while fetching account balance. ACCOUNT_ID: %v",
				accountID,
			)
			r.logger.Error(errStr)
			return nil, &errors.ErrInternalServer{Reason: err}
		}
	}
	return balance, nil

}
