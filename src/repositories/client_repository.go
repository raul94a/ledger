package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	cliententity "src/domain/client"
	errors "src/errors"
	"gorm.io/gorm"

)

type ClientRepository interface {
	FetchClientById(ctx context.Context, ID int) (cliententity.ClientEntity, errors.AppError)
	FetchClientByIdentification(ctx context.Context, identification string) (cliententity.ClientEntity, errors.AppError)
	FetchClient(ctx context.Context, identification string) (cliententity.ClientEntity, errors.AppError)
	InsertClient(ctx context.Context, client *cliententity.ClientEntity) errors.AppError
	InsertClientTx(ctx context.Context, tx *gorm.DB, client *cliententity.ClientEntity) errors.AppError
	GetTx() (*gorm.DB, errors.AppError)
}

type clientRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewClientRepository(db *gorm.DB, logger *zap.Logger) ClientRepository {
	if db == nil {
		panic("db cannot be nil")
	}
	
	return &clientRepository{db: db, logger: logger}
}

func (r *clientRepository) GetTx() (*gorm.DB, errors.AppError) {
	tx := r.db.Begin(&sql.TxOptions{ReadOnly: false})
	
	return tx, nil
}

func (r *clientRepository) FetchClient(ctx context.Context, identification string) (cliententity.ClientEntity, errors.AppError) {
	
	var client cliententity.ClientEntity = cliententity.ClientEntity{}
	result := r.db.Where("identification = ?",identification).Find(&client)
	err := result.Error
	if err == sql.ErrNoRows {
		r.logger.Error("No client found for " + identification)
		return cliententity.ClientEntity{}, &errors.ErrNotFound{Reason: err, Entity: "Client"}

	}
	if err != nil {
		r.logger.Error("Error occurred: " + err.Error())

		return cliententity.ClientEntity{}, &errors.ErrInternalServer{Reason: err}
	}
	return client, nil
}

func (r *clientRepository) FetchClientById(ctx context.Context, ID int) (cliententity.ClientEntity, errors.AppError) {
	
	var client cliententity.ClientEntity = cliententity.ClientEntity{}
	result := r.db.Where("id = ?",ID).Find(&client)
	err := result.Error
	if err == sql.ErrNoRows {
		r.logger.Error("No client found " + fmt.Sprint(ID))
		return cliententity.ClientEntity{}, &errors.ErrNotFound{Entity: "Client", Reason: err}

	}
	if err != nil {
		r.logger.Error("Error occurred: " + err.Error())
		return cliententity.ClientEntity{}, &errors.ErrInternalServer{Reason: err}
	}
	return client, nil
}

func (r *clientRepository) FetchClientByIdentification(ctx context.Context, identification string) (cliententity.ClientEntity, errors.AppError) {
	return r.FetchClient(ctx,identification)
}

func (r *clientRepository) InsertClient(ctx context.Context, client *cliententity.ClientEntity) errors.AppError {
	// Execute the query and scan the returned values into the client struct
	result := r.db.Create(&client)
	err := result.Error
	
	if err != nil {
		r.logger.Error("Error occurred: " + err.Error())

		return &errors.ErrInternalServer{Reason: err}
	}
	return nil
}

func (r *clientRepository) InsertClientTx(ctx context.Context, tx *gorm.DB, client *cliententity.ClientEntity) errors.AppError {
	result := tx.Create(client)
	// Execute the query and scan the returned values into the client struct
	err := result.Error
	if err != nil {
		r.logger.Error("Error occurred: " + err.Error())

		return &errors.ErrInternalServer{Reason: err}
	}
	return nil
}
