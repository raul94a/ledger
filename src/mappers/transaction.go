package mappers

import (
	dto "src/api/dto"
	transaction_entity "src/domain/transaction"
	pagination "src/domain/pagination"
)

func ToPaginationTransactionDto(transactionPagination pagination.Pagination[transaction_entity.TransactionEntity]) (pagination.Pagination[dto.TransactionDto],error){

	items := transactionPagination.Items
	itemsDto := make([]dto.TransactionDto,0)
	for _, el := range items {
		transactionDto,err := ToTransactionDto(el)
		if err != nil {
			return pagination.Pagination[dto.TransactionDto]{}, err
		}
		itemsDto = append(itemsDto, transactionDto)
	}


	var transactionPaginationDto = pagination.Pagination[dto.TransactionDto]{
		Page: transactionPagination.Page,
		LastPage: transactionPagination.LastPage,
		Count: transactionPagination.Count,
		Items: itemsDto,
	}

	return transactionPaginationDto, nil
	
}

func ToTransactionDto(entity transaction_entity.TransactionEntity) (dto.TransactionDto,error){


	var transaction dto.TransactionDto
	// Formatear DateOfBirth
	// El formato "2006-01-02" es para "YYYY-MM-DD"


	// Formatear CreatedAt
	// El formato "2006-01-02 15:04:05" es para "YYYY-MM-DD HH:mm:ss"
	transaction.CreatedAt = entity.CreatedAt

	// Formatear UpdatedAt
	transaction.UpdatedAt = entity.UpdatedAt

	// Asignar campos directos
	transaction.ID = entity.ID
	transaction.AccountID = entity.AccountID
	transaction.Amount = entity.Amount
	transaction.Type = entity.Type

	transaction.ToAccountNumber = nil
	if entity.ToAccountNumber.Valid {
	
		transaction.ToAccountNumber = &entity.ToAccountNumber.String
	}

	
	return transaction, nil
}