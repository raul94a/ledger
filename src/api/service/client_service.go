package services

import (
	"context"
	clientdto "src/api/dto"
	otp_entity "src/domain/registry_accounts_otp"
	app_errors "src/errors"
	"src/mappers"
	"src/repositories"
	"src/utils"
)

type ClientService interface {
	CreateClient(req clientdto.CreateClientRequest) (clientdto.ClientResponse, app_errors.AppError)
}

type clientService struct {
	ClientRepository             repositories.ClientRepository
	RegistryAccountOtpRepository repositories.RegistryAccountOtpRepository
}

func NewClientService(
	clientRepository repositories.ClientRepository,
	registryAccountOtpRepository repositories.RegistryAccountOtpRepository,
) ClientService {
	return &clientService{
		ClientRepository:             clientRepository,
		RegistryAccountOtpRepository: registryAccountOtpRepository,
	}
}

func (s *clientService) CreateClient(req clientdto.CreateClientRequest) (clientdto.ClientResponse, app_errors.AppError) {
	clientEntity, err := mappers.ToClientEntity(req)
	if err != nil {
		return clientdto.ClientResponse{}, &app_errors.ErrBadRequest{Reason: err}
	}
	context := context.Background()
	tx, txError := s.ClientRepository.GetTx()

	if txError != nil {
		return clientdto.ClientResponse{}, txError
	}
	appError := s.ClientRepository.InsertClientTx(context, tx, &clientEntity)
	if appError != nil {
		tx.Rollback()

		return clientdto.ClientResponse{}, appError
	}
	otpCode, err := utils.GenerateRandomOTP(8)
	if err != nil {
		tx.Rollback()
		return clientdto.ClientResponse{}, &app_errors.ErrInternalServer{Reason: err}
	}
	otpEntity := otp_entity.RegisterAccountsOTP{
		ClientID: clientEntity.ID,
		OTP:      otpCode,
	}

	appError = s.RegistryAccountOtpRepository.Insert(context, tx, &otpEntity)
	if appError != nil {
		tx.Rollback()
		return clientdto.ClientResponse{}, appError
	}

	clientResponse, err := mappers.ToClientDTO(clientEntity)
	clientResponse.OTP = otpEntity.OTP

	if err != nil {
		tx.Rollback()
		return clientdto.ClientResponse{}, &app_errors.ErrBadRequest{Reason: err}
	}
	// CreateClientRequest pasa a ClientEntity
	tx.Commit()
	return clientResponse, nil
}
