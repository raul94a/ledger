package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	accountentity "src/domain/account"
	errors "src/errors"

)

type AccountRepository interface {
	FetchAccountById(ctx context.Context, ID int) (accountentity.AccountEntity, errors.AppError)
	FetchAccountIdByAccountNumber(ctx context.Context, iban string) (*int, errors.AppError)
	FetchAccountsByClient(ctx context.Context, clientID int) ([]accountentity.AccountEntity, errors.AppError)
	InsertAccount(ctx context.Context, account *accountentity.AccountEntity) errors.AppError
	InsertAccountTx(ctx context.Context, tx *sql.Tx, account *accountentity.AccountEntity) errors.AppError
	createAccountBalance(ctx context.Context, tx *sql.Tx, account *accountentity.AccountEntity) errors.AppError
}

type accountRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewAccountRepository(db *sql.DB, logger *slog.Logger) AccountRepository {
	if db == nil {
		panic("db cannot be nil")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &accountRepository{db: db, logger: logger}
}

func (r *accountRepository) FetchAccountIdByAccountNumber(ctx context.Context, iban string)(*int, errors.AppError ){
	fmt.Println("ACCOUNT NUMBER ", iban)
	query := `
		SELECT id from accounts where account_number = $1
	`
	

	var id *int
	err := r.db.QueryRowContext(ctx,query,iban).Scan(&id)
	if err == sql.ErrNoRows {
		r.logger.Error("Error occurred: " + err.Error())
		return nil, &errors.ErrNotFound{Entity: "Account", Reason: err}
	}
	if err != nil {
		r.logger.Error("Error occurred: " + err.Error())
		return nil, &errors.ErrInternalServer{Reason: err}
	}
	return id,nil
	
}

func (r *accountRepository) FetchAccountsByClient(ctx context.Context, clientID int) ([]accountentity.AccountEntity, errors.AppError) {
	query := `
	 SELECT * FROM accounts where client_id = $1
	`

	sqlRows, err := r.db.QueryContext(ctx, query, clientID)

	if err != nil {
		r.logger.Error("Error occurred: " + err.Error())
		return nil, &errors.ErrInternalServer{Reason: err}
	}
	defer sqlRows.Close()
	var accounts []accountentity.AccountEntity = make([]accountentity.AccountEntity, 0)

	for sqlRows.Next() {
		var account accountentity.AccountEntity = accountentity.AccountEntity{}
		scanError := sqlRows.Scan(
			&account.ID,
			&account.ClientID,
			&account.AccountNumber,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if scanError != nil {
			r.logger.Error("Error occurred while scanning account: " + err.Error())
			return nil, &errors.ErrInternalServer{Reason: err}
		}
		accounts = append(accounts, account)
	}
	if err := sqlRows.Err(); err != nil {
		r.logger.Error("Error occurred after scanning rows: " + err.Error())
		return nil, &errors.ErrInternalServer{Reason: err}
	}
	// // no error
	if len(accounts) == 0 {
		r.logger.Warn("No accounts found for client: " + fmt.Sprint(clientID))
		return make([]accountentity.AccountEntity, 0) ,nil
	}

	return accounts, nil
}

func (r *accountRepository) FetchAccountById(ctx context.Context, ID int) (accountentity.AccountEntity, errors.AppError) {
	query := `
	 SELECT * FROM accounts where id = $1
	`
	var account accountentity.AccountEntity = accountentity.AccountEntity{}
	sqlRow := r.db.QueryRowContext(ctx, query, ID)
	err := sqlRow.Scan(
		&account.ID,
		&account.ClientID,
		&account.AccountNumber,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		r.logger.Error("No account found " + fmt.Sprint(ID))
		return accountentity.AccountEntity{}, &errors.ErrNotFound{Entity: "Account"}

	}
	if err != nil {
		r.logger.Error("Error occurred: " + err.Error())
		return accountentity.AccountEntity{}, &errors.ErrInternalServer{Reason: err}
	}
	return account, nil
}

func (r *accountRepository) createAccountBalance(ctx context.Context, tx *sql.Tx, account *accountentity.AccountEntity) errors.AppError  {
	query := `
	INSERT INTO account_balances (
            account_id, balance
        ) VALUES ($1, $2)`
	initBalance := 0.0
	_, err := tx.ExecContext(ctx, query, account.ID, initBalance)
	if err != nil {
		errString := fmt.Sprintf("Error inserting new account_balance (ACCOUNT_ID: %d). %s",account.ID,err.Error())
		r.logger.Error(errString)

		return &errors.ErrInternalServer{Reason: err}
	}
	
	return nil
}

func (r *accountRepository) InsertAccountTx(ctx context.Context, tx *sql.Tx,account *accountentity.AccountEntity) errors.AppError {
	
	query := `
	INSERT INTO accounts (
            client_id, account_number
        ) VALUES ($1, $2) 
		RETURNING id,created_at, updated_at`

	// Execute the query and scan the returned values into the client struct
	err := tx.QueryRowContext(ctx, query,
		account.ClientID,
		account.AccountNumber,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		r.logger.Error("Error occurred inserting account: " + err.Error() + " .ClientID: " + fmt.Sprint(account.ClientID))
		tx.Rollback()
		return &errors.ErrInternalServer{Reason: err}
	}
	err = r.createAccountBalance(ctx,tx,account)
	if err != nil {
		tx.Rollback()
		return &errors.ErrInternalServer{Reason: err}
	}
	return nil
}


func (r *accountRepository) InsertAccount(ctx context.Context, account *accountentity.AccountEntity) errors.AppError {
	
	query := `
	INSERT INTO accounts (
            client_id, account_number
        ) VALUES ($1, $2) 
		RETURNING id,created_at, updated_at`

	// Execute the query and scan the returned values into the client struct
	tx, txError := r.db.BeginTx(ctx,&sql.TxOptions{ReadOnly: false})
	if txError != nil {
		return &errors.ErrInternalServer{Reason: txError}
	}
	err := tx.QueryRowContext(ctx, query,
		account.ClientID,
		account.AccountNumber,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		r.logger.Error("Error occurred inserting account: " + err.Error() + " .ClientID: " + fmt.Sprint(account.ClientID))
		tx.Rollback()
		return &errors.ErrInternalServer{Reason: err}
	}
	err = r.createAccountBalance(ctx,tx,account)
	if err != nil {
		tx.Rollback()
		return &errors.ErrInternalServer{Reason: err}
	}
	return nil
}
