package repositories

import (
	"context"
	"database/sql"
	"fmt"
	otp_entity "src/domain/registry_accounts_otp"
	errors "src/errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RegistryAccountOtpRepository interface {
	FetchByClientId(ctx context.Context, clientID int) (otp_entity.RegisterAccountsOTP, errors.AppError)
	Insert(ctx context.Context, tx *gorm.DB, otpEntity *otp_entity.RegisterAccountsOTP) errors.AppError
	Update(ctx context.Context, tx *gorm.DB, otpEntityId int) errors.AppError
}

type registryAccountOtpRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewRegistryAccountOtpRepository(db *gorm.DB, logger *zap.Logger) RegistryAccountOtpRepository {
	if db == nil {
		panic("db cannot be nil")
	}
	
	return &registryAccountOtpRepository{db: db, logger: logger}
}

func (r *registryAccountOtpRepository) FetchByClientId(ctx context.Context, clientId int) (otp_entity.RegisterAccountsOTP, errors.AppError) {
	// query := `
	//  SELECT * FROM register_accounts_otp where client_id = $1
	// `
	var otpEntity otp_entity.RegisterAccountsOTP
	queryDb := r.db.Where("client_id = ?",clientId).Find(&otpEntity)
	error := queryDb.Error
	
	if error == sql.ErrNoRows {
		r.logger.Error("No registry account otp found for " + fmt.Sprint(clientId))
		return otp_entity.RegisterAccountsOTP{}, &errors.ErrNotFound{Reason: error, Entity: "Register Account OTP"}

	}
	if error != nil {
		r.logger.Error("Error occurred: " + error.Error())

		return otp_entity.RegisterAccountsOTP{}, &errors.ErrInternalServer{Reason: error}
	}
	return otpEntity, nil
}

func (r *registryAccountOtpRepository) Insert(ctx context.Context, tx *gorm.DB, otpEntity *otp_entity.RegisterAccountsOTP) errors.AppError {
	
	result:= tx.Create(otpEntity)
	err := result.Error
	
	if err != nil {
		r.logger.Error("Error while inserting REGISTER ACCOUNT OTP entity: " + err.Error())
		return &errors.ErrInternalServer{Reason: err}
	}
	return nil
}

func (r *registryAccountOtpRepository) Update(ctx context.Context, tx *gorm.DB, otpEntityId int) errors.AppError {
	// query := `
    //     UPDATE register_accounts_otp 
	// 	SET
	// 		validated = true,
	// 		updated_at = CURRENT_TIMESTAMP -- Changed from CURRENT to CURRENT_TIMESTAMP
	// 	WHERE
	// 	id = $1
	// 	`
	var otpEntity = otp_entity.RegisterAccountsOTP {
		ID: otpEntityId,
		Validated: true,
		UpdatedAt: time.Now(),
	}
	result := tx.Save(&otpEntity)
	err := result.Error
	// Execute the query and scan the returned values into the client struct
	if err != nil {
		r.logger.Error("Error occurred: " + err.Error())

		return &errors.ErrInternalServer{Reason: err}
	}
	rows := result.RowsAffected
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
