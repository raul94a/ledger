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
	FetchAccountById(ctx context.Context, ID int) (accountentity.AccountEntity, error)
	FetchAccountIdByAccountNumber(ctx context.Context, iban string) (*int, error)
	FetchAccountsByClient(ctx context.Context, clientID int) ([]accountentity.AccountEntity, error)
	InsertAccount(ctx context.Context, account *accountentity.AccountEntity) error
	createAccountBalance(ctx context.Context, tx *sql.Tx, account *accountentity.AccountEntity) error
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

func (r *accountRepository) FetchAccountIdByAccountNumber(ctx context.Context, iban string)(*int, error){
	fmt.Println("ACCOUNT NUMBER ", iban)
	query := `
		SELECT id from accounts where account_number = $1
	`
	fmt.Println("ACCOUNT NUMBER ", iban)

	var id *int
	err := r.db.QueryRowContext(ctx,query,iban).Scan(&id)
	if err == sql.ErrNoRows {
		r.logger.Error("Error occurred: " + err.Error())
		return nil, err
	}
	if err != nil {
		r.logger.Error("Error occurred: " + err.Error())
		return nil, err
	}
	r.logger.Info("ID ES ", id)
	return id,nil
	
}

func (r *accountRepository) FetchAccountsByClient(ctx context.Context, clientID int) ([]accountentity.AccountEntity, error) {
	query := `
	 SELECT * FROM accounts where client_id = $1
	`

	sqlRows, err := r.db.QueryContext(ctx, query, clientID)

	if err != nil {
		r.logger.Error("Error occurred: " + err.Error())
		return nil, err
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
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, account)
	}
	if err := sqlRows.Err(); err != nil {
		r.logger.Error("Error occurred after scanning rows: " + err.Error())
		return nil, fmt.Errorf("failed to fetch accounts: %w", err)
	}

	if len(accounts) == 0 {
		r.logger.Warn("No accounts found for client: " + fmt.Sprint(clientID))
		return nil, &errors.ErrEntityNotFound{Identifier: fmt.Sprint(clientID)}
	}

	return accounts, nil
}

func (r *accountRepository) FetchAccountById(ctx context.Context, ID int) (accountentity.AccountEntity, error) {
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
		return accountentity.AccountEntity{}, &errors.ErrEntityNotFound{Identifier: fmt.Sprint(ID)}

	}
	if err != nil {
		r.logger.Error("Error occurred: " + err.Error())
		return accountentity.AccountEntity{}, err
	}
	return account, nil
}

func (r *accountRepository) createAccountBalance(ctx context.Context, tx *sql.Tx, account *accountentity.AccountEntity) error {
	query := `
	INSERT INTO account_balances (
            account_id, balance
        ) VALUES ($1, $2)`
	initBalance := 0.0
	_, err := tx.ExecContext(ctx, query, account.ID, initBalance)
	if err != nil {
		errString := fmt.Sprintf("Error inserting new account_balance (ACCOUNT_ID: %d). %s",account.ID,err.Error())
		r.logger.Error(errString)

		return err
	}
	
	return nil
}

func (r *accountRepository) InsertAccount(ctx context.Context, account *accountentity.AccountEntity) error {
	tx, txErr := r.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false, Isolation: 0})

	if txErr != nil {
		errorStr := fmt.Sprint("Error creating Tx for Account Insertion. AccountID: %d", account.ID)
		r.logger.Error(errorStr)
		return txErr
	}

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
		return err
	}
	err = r.createAccountBalance(ctx,tx,account)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
