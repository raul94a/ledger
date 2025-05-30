package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	otp_entity "src/domain/registry_accounts_otp"
	errors "src/errors"
)

type RegistryAccountOtpRepository interface {
	FetchByClientId(ctx context.Context, clientID int) (otp_entity.RegisterAccountsOTP, errors.AppError)
	Insert(ctx context.Context, tx *sql.Tx, otpEntity *otp_entity.RegisterAccountsOTP) errors.AppError
	Update(ctx context.Context, tx *sql.Tx, otpEntityId int) errors.AppError
}

type registryAccountOtpRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewRegistryAccountOtpRepository(db *sql.DB, logger *slog.Logger) RegistryAccountOtpRepository {
	if db == nil {
		panic("db cannot be nil")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &registryAccountOtpRepository{db: db, logger: logger}
}

func (r *registryAccountOtpRepository) FetchByClientId(ctx context.Context, clientId int) (otp_entity.RegisterAccountsOTP, errors.AppError) {
	query := `
	 SELECT * FROM register_accounts_otp where client_id = $1
	`
	var registerAccountOtpEntity otp_entity.RegisterAccountsOTP = otp_entity.RegisterAccountsOTP{}
	sqlRow := r.db.QueryRowContext(ctx, query, clientId)
	err := sqlRow.Scan(
		&registerAccountOtpEntity.ID,
		&registerAccountOtpEntity.ClientID,
		&registerAccountOtpEntity.OTP,
		&registerAccountOtpEntity.Validated,
		&registerAccountOtpEntity.CreatedAt,
		&registerAccountOtpEntity.UpdatedAt)
	if err == sql.ErrNoRows {
		r.logger.Error("No registry account otp found for " + fmt.Sprint(clientId))
		return otp_entity.RegisterAccountsOTP{}, &errors.ErrNotFound{Reason: err, Entity: "Register Account OTP"}

	}
	if err != nil {
		r.logger.Error("Error occurred: " + err.Error())

		return otp_entity.RegisterAccountsOTP{}, &errors.ErrInternalServer{Reason: err}
	}
	return registerAccountOtpEntity, nil
}

func (r *registryAccountOtpRepository) Insert(ctx context.Context, tx *sql.Tx, otpEntity *otp_entity.RegisterAccountsOTP) errors.AppError {
	query := `
        INSERT INTO register_accounts_otp (
            client_id, otp
        ) VALUES ($1, $2)
        RETURNING id, created_at, updated_at
		`

	err := tx.QueryRowContext(ctx, query,
		otpEntity.ClientID,
		otpEntity.OTP,
	).Scan(
		&otpEntity.ID,
		&otpEntity.CreatedAt,
		&otpEntity.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Error while inserting REGISTER ACCOUNT OTP entity: " + err.Error())
		return &errors.ErrInternalServer{Reason: err}
	}
	return nil
}

func (r *registryAccountOtpRepository) Update(ctx context.Context, tx *sql.Tx, otpEntityId int) errors.AppError {
	query := `
        UPDATE register_accounts_otp 
		SET
			validated = true,
			updated_at = CURRENT_TIMESTAMP -- Changed from CURRENT to CURRENT_TIMESTAMP
		WHERE
		id = $1
		`

	// Execute the query and scan the returned values into the client struct
	sqlResult, err := tx.ExecContext(ctx, query, otpEntityId)
	if err != nil {
		r.logger.Error("Error occurred: " + err.Error())

		return &errors.ErrInternalServer{Reason: err}
	}
	rows, err := sqlResult.RowsAffected()
	if rows == 0 {
		r.logger.Error("No rows found for otpEntity " + fmt.Sprint(otpEntityId))
		return &errors.ErrInternalServer{Reason: err, Message: fmt.Sprintf("no rows found for entity %d", otpEntityId)}
	}
	if err != nil {
		r.logger.Error("Error occurred: " + err.Error())

		return &errors.ErrInternalServer{Reason: err}
	}
	return nil
}
