package handlers

import (
	"net/http"
	dto "src/api/dto"
	"src/mappers"
	repositories "src/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccountHandler interface {
	FetchAccounts(c *gin.Context)
}

type IAccountHandler struct {
	AccountRepository     repositories.AccountRepository
	TransactionRepository repositories.TransactionRepository
}

// @Summary Returns clients accounts
// @Description Receives the client_id to fetch all the accounts from a Client
// @Accept json
// @Produce json
// @Param client_id
// @Success 200 {object} map[string]interface{} ""
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /accounts/:client_id [get]
func (h *IAccountHandler) FetchAccounts(c *gin.Context) {
	clientIDStr := c.Param("client_id")

	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid identifier"})
		return
	}
	accEntities, err := h.AccountRepository.FetchAccountsByClient(c, clientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "type": "fetching client accounts"})
		return
	}
	var accounts []dto.AccountDto
	for _, accountEntity := range accEntities {
		accountDto, err := mappers.ToAccountDTO(accountEntity, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		balance, err := h.TransactionRepository.FetchAccountBalance(c,nil,accountDto.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		accountDto.Balance = *balance
		accounts = append(accounts, accountDto)
	}
	

	c.JSON(http.StatusOK, gin.H{
		"message": "Cliente creado exitosamente",
		"accounts":  accounts,
	})
}
