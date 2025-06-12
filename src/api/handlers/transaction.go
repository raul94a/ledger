package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	dto "src/api/dto"
	trasnactionentity "src/domain/transaction"
	mappers "src/mappers"
	repositories "src/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionHandler interface {
	PerformTransaction(c *gin.Context)
	GetTransactions(c *gin.Context)
}

type ITransactionHandler struct {
	TransactionRepository repositories.TransactionRepository
	AccountRepository     repositories.AccountRepository
}

// POST
func (h *ITransactionHandler) PerformTransaction(c *gin.Context) {
	var performnTransactionDto dto.PerformTransactionDto
	if error := c.ShouldBindJSON(&performnTransactionDto); error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
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

	for _, element := range validTransactionTypes {
		if element == performnTransactionDto.Type {
			isValidType = true
		}
	}
	if !isValidType {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction type is not valid"})
		return
	}
	// 2. Transfer with Non empty ToAccountId
	if performnTransactionDto.ToAccountNumber == nil && performnTransactionDto.Type == "TRANSFER" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction type is not valid", "reason": "TRANSFER type needs to set to_account_id"})
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
		if err != nil {
			err.JsonError(c)
			return
		}
		if id == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
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
	transactionEntity := trasnactionentity.TransactionEntity{
		AccountID:   performnTransactionDto.AccountID,
		ToAccountID: toAccountIdSql,
		Type:        performnTransactionDto.Type,
		Amount:      performnTransactionDto.Amount,
	}

	err := h.TransactionRepository.InsertTransactionLedgerTx(c, &transactionEntity)
	if err != nil {

		err.JsonError(c)
		return
	}
	transactionDto := dto.TransactionDto{
		ID:              transactionEntity.ID,
		AccountID:       transactionEntity.AccountID,
		Type:            transactionEntity.Type,
		Amount:          transactionEntity.Amount,
		ToAccountNumber: performnTransactionDto.ToAccountNumber,
		CreatedAt:       transactionEntity.CreatedAt,
		UpdatedAt:       transactionEntity.UpdatedAt,
	}
	c.JSON(http.StatusOK, gin.H{"transaction": transactionDto})
}

func (h *ITransactionHandler) GetTransactions(c *gin.Context) {
	accountId := c.Param("account_id")
	countStr := c.Query("count")
	pageStr := c.Query("page")

	// token validation, account id validation
	accountIdInt, err := strconv.ParseInt(accountId, 0, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	countInt, err := strconv.ParseInt(countStr, 0, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	pageInt, err := strconv.ParseInt(pageStr, 0, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	// Pre-validation
	// account_id belongs to clientId ? 

	pagination, error := h.TransactionRepository.GetTransactions(context.Background(), int(accountIdInt), int(pageInt), int(countInt))
	if error != nil {
		error.JsonError(c)
		return
	}
	paginationDto, err := mappers.ToPaginationTransactionDto(pagination)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(http.StatusOK,paginationDto)

	// fetchFromRepository
}
