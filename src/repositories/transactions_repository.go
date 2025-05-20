package repositories

import (
	"context"
	"database/sql"
	"log/slog"
	ledgerentity "src/domain/ledger"
	transaction_entity "src/domain/transaction"
)

type TransactionRepository interface {
	// FetchTransactionById(ctx context.Context, ID int) (transaction_entity.TransactionEntity, error)
	// FetchTransactionsByAccount(ctx context.Context, accountID int) ([]transaction_entity.TransactionEntity, error)
	InsertTransaction(ctx context.Context, tx *sql.Tx, transaction *transaction_entity.TransactionEntity) error
	InsertLedgerEntry(ctx context.Context, tx *sql.Tx, ledgerTransaction *ledgerentity.LedgerTransaction) error
	InsertTransactionLedgerTx(ctx context.Context, transaction *transaction_entity.TransactionEntity) error
	// UpdateAccountBalance(ctx context.Context, tx *sql.Tx,transaction *transaction_entity.TransactionEntity) error
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

func (r *transactionRepository) InsertTransaction(ctx context.Context, tx *sql.Tx, transaction *transaction_entity.TransactionEntity) error {
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
		r.logger.Error("Account id: %s, amount: %v, type: %s", transaction.AccountID, transaction.Amount,transaction.Type)
		return &ErrTransactionInsertFailed{Message: "", Reason: err}
	}
	return nil
}

func (r *transactionRepository) InsertLedgerEntry(ctx context.Context, tx *sql.Tx, ledgerTransaction *ledgerentity.LedgerTransaction) error {
	query := `
        INSERT INTO ledger_entries (
            transaction_id, account_id, type, amount
        ) VALUES ($1, $2, $3, $4)`

	transaction := ledgerTransaction.Transaction
	// Execute the query and scan the returned values into the client struct
	result, err := tx.ExecContext(ctx, query,
		transaction.ID,
		ledgerTransaction.AccountID,
		ledgerTransaction.TransactionType,
		transaction.Amount,
	)
	
	if err != nil {
		r.logger.Error("Error occurred while inserting ledger_entry: %s ", err.Error())
		isFromTransaction := ledgerTransaction.Transaction.AccountID == ledgerTransaction.AccountID
		origin := ""
		if isFromTransaction {
			origin = "(from)"
		} else {
			origin = "(to)"
		}
		return &ErrLedgerEntryInsertFailed{Message: origin, Reason: err}
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		

		return &ErrNoRowsAffected{Message: "Problem during inserting ledgerEntry"}
	}
	return nil
}

func (r *transactionRepository) InsertTransactionLedgerTx(ctx context.Context, transaction *transaction_entity.TransactionEntity) error {
	tx, txErr := r.db.BeginTx(ctx, &sql.TxOptions{
		ReadOnly:  false,
		Isolation: 0,
	})

	if txErr != nil {
		r.logger.Error("Error occurred while beginning transaction ledger: %s ", txErr.Error())
		return txErr
	}

	err := r.InsertTransaction(ctx, tx, transaction)

	if err != nil {
		r.logger.Error("Error occurred in Txn while inserting transaction: %s ", err.Error())
		tx.Rollback()
		return err
	}

	// we always starts with the account triggering the transaction
	transactionType := ""
	if transaction.Type != "ADD" {
		transactionType = "DEBIT"
	} else {
		transactionType = "CREDIT"
	}
	transactionLedger := ledgerentity.LedgerTransaction{
		Transaction:     *transaction,
		TransactionType: transactionType,
		AccountID:       transaction.AccountID,
	}
	err = r.InsertLedgerEntry(ctx, tx, &transactionLedger)
	if err != nil {
		r.logger.Error("Error occurred in Txn %d while inserting transaction_ledger (from): %s ",transaction.ID, err.Error())
		r.logger.Error("Account id: %s, amount: %v, type: %s", transaction.AccountID, transaction.Amount,transaction.Type)

		tx.Rollback()
		return err
	}

	if !transaction.ToAccountID.Valid {
		r.logger.Warn("No attached ToAccountId in Transaction %d. Tx finishing without errors.", transaction.ID)
		tx.Commit()
		return nil
	}
	// If there's a ToAccountId, another ledger entry has to be inserted
	transactionLedger.Transaction.AccountID = int(transaction.ToAccountID.Int32)

	if transactionType == "CREDIT" {
		transactionType = "DEBIT"
	} else {
		transactionType = "CREDIT"
	}

	err = r.InsertLedgerEntry(ctx, tx, &transactionLedger)
	if err != nil {
		r.logger.Error("Error occurred in Txn %d while inserting transaction_ledger (to): %s ",transaction.ID, err.Error())
		r.logger.Error("Account id: %s, amount: %v, type: %s", transaction.AccountID, transaction.Amount,transaction.Type)
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil

}
