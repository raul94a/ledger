package handlers

import (
	"database/sql"
	"net/http"
	"fmt"
	dto "src/api/dto"
	trasnactionentity "src/domain/transaction"
	repositories "src/repositories"

	"github.com/gin-gonic/gin"
)

type TransactionHandler interface {
	PerformTransaction(c *gin.Context)
}

type ITransactionHandler struct {
	TransactionRepository repositories.TransactionRepository
	AccountRepository     repositories.AccountRepository
}

// POST
func (h *ITransactionHandler) PerformTransaction(c *gin.Context) {
	var performnTransactionDto dto.PerformTransactionDto
	if err := c.ShouldBindJSON(&performnTransactionDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// pre-validation
	// 
	// Is AccountID from a trusted source?
	// The API must protect accounts from impersonation
	// so, check if the account belongs to the user calling this webservice is necessary

	// Validate fields
	// 1. Transaction type
	var validTransactionTypes = []string{"ADD", "TRANSFER", "WITHDRAWAL"}
	isValidType := false

	for _,element := range validTransactionTypes {
		if element == performnTransactionDto.Type {
			isValidType = true
		}
	}
	if !isValidType {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction type is not valid"})
		return
	} 
	// 2. Transfer with Non empty ToAccountId
	if performnTransactionDto.ToAccountNumber == nil && performnTransactionDto.Type == "TRANSFER"  {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction type is not valid", "reason":"TRANSFER type needs to set to_account_id"})
		return
	}
	// 3. Negative money
	if performnTransactionDto.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount cannot be negative."})
		return
	}
	toAccountIdSql := sql.NullInt32{}
	if performnTransactionDto.ToAccountNumber != nil {
		// Fetch AccountID
		accNr := *performnTransactionDto.ToAccountNumber
		fmt.Println("ACCOUNT NUMBER ", accNr)
		id, err := h.AccountRepository.FetchAccountIdByAccountNumber(c, accNr)
		if err != nil || id == nil{
			if err == sql.ErrNoRows{
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		toAccountIdSql.Valid = true
		toAccountIdSql.Int32 = int32(*id)

	} else {
		toAccountIdSql = sql.NullInt32{
			Int32: int32(-1),
			Valid: false,
		}

	}
	transactionEntity := trasnactionentity.TransactionEntity {
		AccountID: performnTransactionDto.AccountID,
		ToAccountID: toAccountIdSql,
		Type: performnTransactionDto.Type,
		Amount: performnTransactionDto.Amount,
	}

	err := h.TransactionRepository.InsertTransactionLedgerTx(c,&transactionEntity)
	if err != nil {
		
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})		
		return
	}
	transactionDto := dto.TransactionDto {
		ID: transactionEntity.ID,
		AccountID: transactionEntity.AccountID,
		Type: transactionEntity.Type,
		Amount: transactionEntity.Amount,
		ToAccountNumber: performnTransactionDto.ToAccountNumber,
		CreatedAt: transactionEntity.CreatedAt,
		UpdatedAt: transactionEntity.UpdatedAt,
	}
	c.JSON(http.StatusOK, gin.H{"transaction": transactionDto})
}



