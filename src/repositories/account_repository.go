package repositories

import (
	"context"
	"database/sql"
	"fmt"
	accountentity "src/domain/account"
	errors "src/errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AccountRepository interface {
	FetchAccountById(ctx context.Context, ID int) (accountentity.AccountEntity, errors.AppError)
	FetchAccountIdByAccountNumber(ctx context.Context, iban string) (*int, errors.AppError)
	FetchAccountsByClient(ctx context.Context, clientID int) ([]accountentity.AccountEntity, errors.AppError)
	InsertAccount(ctx context.Context, account *accountentity.AccountEntity) errors.AppError
	InsertAccountTx(ctx context.Context, tx *gorm.DB, account *accountentity.AccountEntity) errors.AppError
	createAccountBalance(ctx context.Context, tx *gorm.DB, account *accountentity.AccountEntity) errors.AppError
}

type accountRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAccountRepository(db *gorm.DB, logger *zap.Logger) AccountRepository {
	if db == nil {
		panic("db cannot be nil")
	}
	return &accountRepository{db: db, logger: logger}
}

func (r *accountRepository) FetchAccountIdByAccountNumber(ctx context.Context, iban string)(*int, errors.AppError ){
	result := r.db.Where("account_number = ?", iban)
	err := result.Error
	var id *int
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

	var accounts []accountentity.AccountEntity
	result := r.db.Where("client_id = ?",clientID).Find(&accounts)
	err := result.Error

	if err != nil {
		r.logger.Error("Error occurred: " + err.Error())
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
	
	var account accountentity.AccountEntity = accountentity.AccountEntity{}
	result := r.db.Where("id = ?",ID).Find(&account)
	err := result.Error
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

func (r *accountRepository) createAccountBalance(ctx context.Context, tx *gorm.DB, account *accountentity.AccountEntity) errors.AppError  {
	query := `
	INSERT INTO account_balances (
            account_id, balance
        ) VALUES ($1, $2)`
	initBalance := 0.0

	tx = tx.Raw(query, account.ID, initBalance)
	err := tx.Error
	if err != nil {
		errString := fmt.Sprintf("Error inserting new account_balance (ACCOUNT_ID: %d). %s",account.ID,err.Error())
		r.logger.Error(errString)

		return &errors.ErrInternalServer{Reason: err}
	}
	
	return nil
}

func (r *accountRepository) InsertAccountTx(ctx context.Context, tx *gorm.DB,account *accountentity.AccountEntity) errors.AppError {
	
	// Execute the query and scan the returned values into the client struct
	result := r.db.Create(account)
	err := result.Error


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
	
	// Execute the query and scan the returned values into the client struct
	tx  := r.db.Begin(&sql.TxOptions{ReadOnly: false})
	if tx.Error != nil {
		return &errors.ErrInternalServer{Reason: tx.Error}
	}
	result := tx.Create(account)
	err := result.Error
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
