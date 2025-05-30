package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	cliententity "src/domain/client"
	errors "src/errors"
)

type ClientRepository interface {
	FetchClientById(ctx context.Context, ID int) (cliententity.ClientEntity, errors.AppError)
	FetchClientByIdentification(ctx context.Context, identification string) (cliententity.ClientEntity, errors.AppError)
	FetchClient(ctx context.Context, identification string) (cliententity.ClientEntity, errors.AppError)
	InsertClient(ctx context.Context, client *cliententity.ClientEntity) errors.AppError
	InsertClientTx(ctx context.Context, tx *sql.Tx, client *cliententity.ClientEntity) errors.AppError
	GetTx() (*sql.Tx, errors.AppError)
}

type clientRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewClientRepository(db *sql.DB, logger *slog.Logger) ClientRepository {
	if db == nil {
		panic("db cannot be nil")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &clientRepository{db: db, logger: logger}
}

func (r *clientRepository) GetTx() (*sql.Tx, errors.AppError) {
	tx, err := r.db.BeginTx(context.Background(), &sql.TxOptions{ReadOnly: false})
	if err != nil {
		r.logger.Error("Error gettingTransaction " + err.Error())
		return nil, &errors.ErrInternalServer{Reason: err}
	}
	return tx, nil
}

func (r *clientRepository) FetchClient(ctx context.Context, identification string) (cliententity.ClientEntity, errors.AppError) {
	query := `
	 SELECT * FROM clients where identification = $1
	`
	var client cliententity.ClientEntity = cliententity.ClientEntity{}
	sqlRow := r.db.QueryRowContext(ctx, query, identification)
	err := sqlRow.Scan(&client.ID,
		&client.Name,
		&client.Surname1,
		&client.Surname2,
		&client.Email,
		&client.Identification,
		&client.Nationality,
		&client.DateOfBirth,
		&client.Sex,
		&client.Address,
		&client.City,
		&client.Province,
		&client.State,
		&client.ZipCode,
		&client.Telephone,
		&client.CreatedAt,
		&client.UpdatedAt, 
		&client.KcUserId)
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
	query := `
	 SELECT * FROM clients where id = $1
	`
	var client cliententity.ClientEntity = cliententity.ClientEntity{}
	sqlRow := r.db.QueryRowContext(ctx, query, ID)
	err := sqlRow.Scan(&client.ID,
		&client.Name,
		&client.Surname1,
		&client.Surname2,
		&client.Email,
		&client.Identification,
		&client.Nationality,
		&client.DateOfBirth,
		&client.Sex,
		&client.Address,
		&client.City,
		&client.Province,
		&client.State,
		&client.ZipCode,
		&client.Telephone,
		&client.CreatedAt,
		&client.UpdatedAt, 
		&client.KcUserId)
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
	query := `
	 SELECT * FROM clients where identification = $1
	`
	var client cliententity.ClientEntity = cliententity.ClientEntity{}
	sqlRow := r.db.QueryRowContext(ctx, query, identification)
	err := sqlRow.Scan(&client.ID,
		&client.Name,
		&client.Surname1,
		&client.Surname2,
		&client.Email,
		&client.Identification,
		&client.Nationality,
		&client.DateOfBirth,
		&client.Sex,
		&client.Address,
		&client.City,
		&client.Province,
		&client.State,
		&client.ZipCode,
		&client.Telephone,
		&client.CreatedAt,
		&client.UpdatedAt, 
		&client.KcUserId)
	if err == sql.ErrNoRows {
		r.logger.Error("No client found " + fmt.Sprint(identification))
		return cliententity.ClientEntity{}, &errors.ErrNotFound{Entity: "Client", Reason: err}

	}
	if err != nil {
		r.logger.Error("Error occurred: " + err.Error())
		return cliententity.ClientEntity{}, &errors.ErrInternalServer{Reason: err}
	}
	return client, nil
}

func (r *clientRepository) InsertClient(ctx context.Context, client *cliententity.ClientEntity) errors.AppError {
	query := `
        INSERT INTO clients (
            name, surname1, surname2, email, identification, nationality, 
            date_of_birth, sex, address, city, province, state, 
            zip_code, telephone
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
        RETURNING id, created_at, updated_at`

	// Execute the query and scan the returned values into the client struct
	err := r.db.QueryRowContext(ctx, query,
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
		r.logger.Error("Error occurred: " + err.Error())

		return &errors.ErrInternalServer{Reason: err}
	}
	return nil
}

func (r *clientRepository) InsertClientTx(ctx context.Context, tx *sql.Tx, client *cliententity.ClientEntity) errors.AppError {
	query := `
        INSERT INTO clients (
            name, surname1, surname2, email, identification, nationality, 
            date_of_birth, sex, address, city, province, state, 
            zip_code, telephone
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
        RETURNING id, created_at, updated_at`

	// Execute the query and scan the returned values into the client struct
	err := tx.QueryRowContext(ctx, query,
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
		r.logger.Error("Error occurred: " + err.Error())

		return &errors.ErrInternalServer{Reason: err}
	}
	return nil
}
