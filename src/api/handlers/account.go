package handlers

import (
	"context"
	"net/http"
	dto "src/api/dto"
	api_keycloak "src/api/keycloak"
	services "src/api/service"
	"src/mappers"
	repositories "src/repositories"
	"strconv"
	"github.com/gin-gonic/gin"
	"fmt"
)

type AccountHandler interface {
	FetchAccounts(c *gin.Context)
	CreateAccount(c *gin.Context)
	// UpdateRegisterAccountOtpStatus(c *gin.Context)
	// CreateUser(c *gin.Context)
	CompleteNewUserRegistration(c *gin.Context)
}

type IAccountHandler struct {
	KeycloakClient               api_keycloak.KeycloakClient
	AccountService               services.AccountService
	ClientRepository             repositories.ClientRepository
	AccountRepository            repositories.AccountRepository
	TransactionRepository        repositories.TransactionRepository
	RegistryAccountOtpRepository repositories.RegistryAccountOtpRepository
}

// @Summary Returns client accounts
// @Description Receives the client_id to fetch all the accounts from a Client
// @Accept json
// @Produce json
// @Param client_id path int true "The ID of the client whose accounts want to be fetched"
// @Success 200 {object} []dto.AccountDto "Successfully retrieved client accounts"
// @Failure 400 {object} map[string]string "Bad Request - Invalid client ID"
// @Router /accounts/{client_id} [get] // Changed :client_id to {client_id} for Swagger consistency
// @security BearerAuth
// @tags Accounts
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

// @Summary Creates a new account
// @Description Creates a new bank account for a given client ID
// @Accept json
// @Produce json
// @Param request body dto.CreateAccountRequest true "Client ID for the new account"
// @Success 200 {object} object{account=dto.AccountDto} "Created account successfully"
// @Failure 400 {object} map[string]string "Bad Request - Invalid input or client ID"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /accounts [post]
// @security BearerAuth
// @tags Accounts
func (h *IAccountHandler) CreateAccount(c *gin.Context) {
	var createAccountReq dto.CreateAccountRequest
	if error := c.ShouldBindJSON(&createAccountReq); error != nil {
		statusCode := http.StatusBadRequest
		reason := http.StatusText(statusCode)
		c.JSON(statusCode, gin.H{"error": error.Error(), "reason": reason})
		return
	}

	account, err := h.AccountService.CreateAccount(createAccountReq.ClientID)

	if err != nil {
		err.JsonError(c)
		return
	}

	c.JSON(http.StatusOK, account)

}


// @Summary Completes new user registration
// @Description Finalizes the registration process for a new user, including Keycloak user creation and account service completion.
// @Accept json
// @Produce json
// @Param request body dto.CompleteClientRegistrationBankAccountRequest true "Registration details for the new user"
// @Success 201 {object} map[string]interface{} "User registration completed successfully" // 201 for creation
// @Failure 400 {object} map[string]string "Bad Request - Invalid input"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /accounts/completeNewUserRegistration [post] // Assuming this is the intended route
// @tags Accounts
func (h *IAccountHandler) CompleteNewUserRegistration(c *gin.Context) {
	var completeClientRegistration dto.CompleteClientRegistrationBankAccountRequest
	if err := c.Bind(&completeClientRegistration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 1. Fetch the client
	ctx := context.Background()
	clientEntity, err := h.ClientRepository.FetchClientByIdentification(ctx, completeClientRegistration.Identification)
	if err != nil {
		err.JsonError(c)
		return
	}
	// 2. Create Keycloak User
	credential := api_keycloak.KcCredentials{
		Type:  "password",
		Value: completeClientRegistration.Pin,
	}
	credentials := make([]api_keycloak.KcCredentials, 0)
	credentials = append(credentials, credential)
	kcCreateUserRequest := api_keycloak.KcCreateUserRequest{
		Username:    completeClientRegistration.Identification,
		FirstName:   clientEntity.Name,
		LastName:    clientEntity.Surname1,
		Email:       clientEntity.Email,
		Enabled:     true,
		Credentials: credentials,
		Attributes: map[string][]string{},
	}
	kcCreateUserRequest.AddAttribute("client_id", fmt.Sprint(clientEntity.ID))
	err = h.KeycloakClient.CreateUser(kcCreateUserRequest)
	if err != nil {
		err.JsonError(c)
		return
	}
	// 3. Call Service for completion of new client registration!

	err = h.AccountService.CompleteClientRegistrationBankAccount(completeClientRegistration, clientEntity)
	if err != nil {
		err.JsonError(c)
		return
	}
	c.JSON(http.StatusCreated, gin.H{})

	// 3. Return OK

}
