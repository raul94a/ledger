package handlers

import (
	"net/http"
	dto "src/api/dto"
	accountentity "src/domain/account"
	"src/mappers"
	repositories "src/repositories"
	"src/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccountHandler interface {
	FetchAccounts(c *gin.Context)
	CreateAccount(c *gin.Context)
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
		accountDto:= mappers.ToAccountDTO(accountEntity, nil)
		
		balance, err := h.TransactionRepository.FetchAccountBalance(c, nil, accountDto.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		accountDto.Balance = *balance
		accounts = append(accounts, accountDto)
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
	})
}

func (h *IAccountHandler) CreateAccount(c *gin.Context) {
	var createAccountReq dto.CreateAccountRequest
	const spainCode string = "ES"
	const bankDigits string = "0182"
	const branchDigits string = "0600"
	if err := c.ShouldBindJSON(&createAccountReq); err != nil {
		statusCode := http.StatusBadRequest
		reason := http.StatusText(statusCode)
		c.JSON(statusCode, gin.H{"error": err.Error(), "reason": reason})
		return
	}
	handler := utils.IbanHandler{}

	accNumber := handler.GenerateAccountNumber(10)
	cc,err := handler.DomesticCheckDigits(bankDigits,branchDigits,accNumber)
	if err != nil {
		code := http.StatusInternalServerError
		// LOG
		c.JSON(code, gin.H{"error": "something occurred during IBAN creation", "reason": http.StatusText(code)})
		return
	}
	bban := utils.Bban {
		BankCode: bankDigits,
		BranchCode: branchDigits,
		DomesticCheckDigits: cc,
		AccountNumber: accNumber,
	}
	iban, err := handler.ComputeIban(bban,spainCode)
	if err != nil {
		statusCode := http.StatusBadRequest
		reason := http.StatusText(statusCode)
		// LOG
		c.JSON(statusCode, gin.H{"error": err.Error(), "reason": reason})
		return
	}

	accountEntity := accountentity.AccountEntity{
		ClientID:      createAccountReq.ClientID,
		AccountNumber: iban,
	}

	err = h.AccountRepository.InsertAccount(c, &accountEntity)
	if err != nil {
		statusCode := http.StatusBadRequest
		reason := http.StatusText(statusCode)
		c.JSON(statusCode, gin.H{"error": err.Error(), "reason": reason})
		return
	}
	balance := 0.0
	var balancePtr *float64 = &balance
	
	c.JSON(http.StatusOK, gin.H{
		"account": mappers.ToAccountDTO(accountEntity,balancePtr),
	})

}
