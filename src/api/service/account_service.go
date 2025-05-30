package services

import (
	"context"
	"database/sql"
	dto "src/api/dto"
	accountentity "src/domain/account"
	cliententity "src/domain/client"
	app_errors "src/errors"
	"src/mappers"
	"src/repositories"
	"src/utils"
)

type AccountService interface {
	CreateAccount(clientId int) (dto.AccountDto, app_errors.AppError)
	CreateAccountTx(context context.Context, tx *sql.Tx, clientId int) (dto.AccountDto, app_errors.AppError)
	CompleteClientRegistrationBankAccount(
		req dto.CompleteClientRegistrationBankAccountRequest,
		clientEntity cliententity.ClientEntity,
	) app_errors.AppError
}

type accountService struct {
	RepositoryWrapper repositories.RepositoryWrapper
}

func NewAccountService(
	wrapper repositories.RepositoryWrapper,
) AccountService {
	return &accountService{
		RepositoryWrapper: wrapper,
	}
}

func (h *accountService) CreateAccount(clientId int) (dto.AccountDto, app_errors.AppError) {

	const spainCode string = "ES"
	const bankDigits string = "0182"
	const branchDigits string = "0600"
	handler := utils.IbanHandler{}
	accNumber := handler.GenerateAccountNumber(10)
	cc := handler.DomesticCheckDigits(bankDigits, branchDigits, accNumber)

	bban := utils.Bban{
		BankCode:            bankDigits,
		BranchCode:          branchDigits,
		DomesticCheckDigits: cc,
		AccountNumber:       accNumber,
	}
	iban, err := handler.ComputeIban(bban, spainCode)
	if err != nil {
		return dto.AccountDto{}, nil

	}

	accountEntity := accountentity.AccountEntity{
		ClientID:      clientId,
		AccountNumber: iban,
	}
	context := context.Background()
	error := h.RepositoryWrapper.AccountRepository.InsertAccount(context, &accountEntity)
	if error != nil {

		return dto.AccountDto{}, nil
	}
	balance := 0.0
	var balancePtr *float64 = &balance

	return mappers.ToAccountDTO(accountEntity, balancePtr), nil

}

func (h *accountService) CreateAccountTx(context context.Context, tx *sql.Tx, clientId int) (dto.AccountDto, app_errors.AppError) {
	const spainCode string = "ES"
	const bankDigits string = "0182"
	const branchDigits string = "0600"
	handler := utils.IbanHandler{}
	accNumber := handler.GenerateAccountNumber(10)
	cc := handler.DomesticCheckDigits(bankDigits, branchDigits, accNumber)

	bban := utils.Bban{
		BankCode:            bankDigits,
		BranchCode:          branchDigits,
		DomesticCheckDigits: cc,
		AccountNumber:       accNumber,
	}
	iban, err := handler.ComputeIban(bban, spainCode)
	if err != nil {
		return dto.AccountDto{}, nil

	}

	accountEntity := accountentity.AccountEntity{
		ClientID:      clientId,
		AccountNumber: iban,
	}

	error := h.RepositoryWrapper.AccountRepository.InsertAccountTx(context, tx, &accountEntity)
	if error != nil {

		return dto.AccountDto{}, nil
	}
	balance := 0.0
	var balancePtr *float64 = &balance

	return mappers.ToAccountDTO(accountEntity, balancePtr), nil

}

func (s *accountService) CompleteClientRegistrationBankAccount(
	req dto.CompleteClientRegistrationBankAccountRequest,
	clientEntity cliententity.ClientEntity,
) app_errors.AppError {

	ctx := context.Background()

	// Fetch Register Account Otp for this entity
	otpEntity, err := s.RepositoryWrapper.RegistryAccountOtpRepository.FetchByClientId(ctx, clientEntity.ID)
	if err != nil {
		return err
	}
	if otpEntity.OTP != req.OTP {
		return &app_errors.ErrBadRequest{Message: "not valid code"}
	}
	if otpEntity.Validated {

		return nil
	}

	// Pre-validation: do we have this user in Keycloak?
	// TODO: 1. Create Keycloak user by KeycloakRepository

	// 2. Update status from Register Accounts Otp
	tx, txErr := s.RepositoryWrapper.ClientRepository.GetTx()
	if txErr != nil {
		return &app_errors.ErrInternalServer{Reason: txErr}
	}

	err = s.RepositoryWrapper.RegistryAccountOtpRepository.Update(ctx, tx, otpEntity.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// 3. Create Bank Account
	_, err = s.CreateAccountTx(ctx, tx, clientEntity.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// 4. Commit Transacction changes
	tx.Commit()

	return nil

}
