package mappers

import (
	"fmt"
	accountdto "src/api/dto"
	accountentity "src/domain/account"
	"time"
)

func ToAccountEntity(account accountdto.AccountDto) (accountentity.AccountEntity, error) {
	var entity accountentity.AccountEntity

	// Parsear CreatedDate
	// Nota: El formato "2006-01-02 15:04:05" es para "YYYY-MM-DD HH:mm:ss"
	createdDate, err := time.Parse("2006-01-02 15:04:05", account.CreatedDate)
	if err != nil {
		return accountentity.AccountEntity{}, fmt.Errorf("error parsing CreatedDate: %w", err)
	}

	// Parsear UpdatedDate
	updatedDate, err := time.Parse("2006-01-02 15:04:05", account.UpdatedDate)
	if err != nil {
		return accountentity.AccountEntity{}, fmt.Errorf("error parsing UpdatedDate: %w", err)
	}

	// Asignar campos directos
	entity.AccountNumber = account.AccountNumber
	entity.ClientID = account.ClientID
	entity.ID = account.ID
	// Asignar fechas parseadas
	entity.CreatedAt = createdDate
	entity.UpdatedAt = updatedDate

	return entity, nil
}

func ToAccountDTO(entity accountentity.AccountEntity, balance *float64) accountdto.AccountDto {
	var dto accountdto.AccountDto

	// Formatear CreatedAt
	// El formato "2006-01-02 15:04:05" es para "YYYY-MM-DD HH:mm:ss"
	dto.CreatedDate = entity.CreatedAt.Format("2006-01-02 15:04:05")

	// Formatear UpdatedAt
	dto.UpdatedDate = entity.UpdatedAt.Format("2006-01-02 15:04:05")

	// Asignar campos directos
	dto.ID = entity.ID
	dto.AccountNumber = entity.AccountNumber
	dto.ClientID = entity.ClientID
	if balance != nil {
		dto.Balance = *balance
	}

	return dto
}
