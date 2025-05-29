package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	dto "src/api/dto"
	accountentity "src/domain/account"
	"src/mappers"
	repositories "src/repositories"
	"src/utils"
	"strconv"
)

type AccountHandler interface {
	FetchAccounts(c *gin.Context)
	CreateAccount(c *gin.Context)
}

type IAccountHandler struct {
	AccountRepository     repositories.AccountRepository
	TransactionRepository repositories.TransactionRepository
	RegistryAccountOtpRepository repositories.RegistryAccountOtpRepository
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
	accEntities, error := h.AccountRepository.FetchAccountsByClient(c, clientID)
	if error != nil {
		error.JsonError(c)
		return
	}
	if len(accEntities) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"accounts": make([]dto.AccountDto, 0),
		})
		return
	}
	var accounts []dto.AccountDto
	for _, accountEntity := range accEntities {
		accountDto := mappers.ToAccountDTO(accountEntity, nil)

		balance, error := h.TransactionRepository.FetchAccountBalance(c, nil, accountDto.ID)
		if error != nil {
			error.JsonError(c)
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
	if error := c.ShouldBindJSON(&createAccountReq); error != nil {
		statusCode := http.StatusBadRequest
		reason := http.StatusText(statusCode)
		c.JSON(statusCode, gin.H{"error": error.Error(), "reason": reason})
		return
	}
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

	error := h.AccountRepository.InsertAccount(c, &accountEntity)
	if error != nil {
		error.JsonError(c)
		return
	}
	balance := 0.0
	var balancePtr *float64 = &balance

	c.JSON(http.StatusOK, gin.H{
		"account": mappers.ToAccountDTO(accountEntity, balancePtr),
	})

}

func (h *IAccountHandler) CreateAccountTx(c *gin.Context) {
	var createAccountReq dto.CreateAccountRequest
	const spainCode string = "ES"
	const bankDigits string = "0182"
	const branchDigits string = "0600"
	if error := c.ShouldBindJSON(&createAccountReq); error != nil {
		statusCode := http.StatusBadRequest
		reason := http.StatusText(statusCode)
		c.JSON(statusCode, gin.H{"error": error.Error(), "reason": reason})
		return
	}
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

	error := h.AccountRepository.InsertAccount(c, &accountEntity)
	if err != nil {
		error.JsonError(c)
		return
	}
	balance := 0.0
	var balancePtr *float64 = &balance

	c.JSON(http.StatusOK, gin.H{
		"account": mappers.ToAccountDTO(accountEntity, balancePtr),
	})

}
