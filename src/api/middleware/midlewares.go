package middleware

import (
	"fmt"
	"log"
	"net/http"
	dto "src/api/dto"
	api_keycloak "src/api/keycloak"
	"src/repositories"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("Request - Method: %s | Status: %d | Duration: %v", c.Request.Method, c.Writer.Status(), duration)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	// In a real-world application, you would perform proper authentication here.
	// For the sake of this example, we'll just check if an API key is present.
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}

func KeycloakClientMiddleware(c *gin.Context, client api_keycloak.KeycloakClient) {
	c.Set("keycloak_client", client)
	c.Next()
}

func RepositoryWrapperMiddleware(c *gin.Context, repositoryWrapper *repositories.RepositoryWrapper) {
	c.Set("repository_wrapper", repositoryWrapper)
	c.Next()
}

func AuthorizationMiddleware(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	splittedToken := strings.Split(accessToken, " ")
	if splittedToken[0] != "Bearer" {
		c.AbortWithStatusJSON(400, gin.H{"error": "Bad request"})
		return
	}
	keycloakClient, exists := c.Get("keycloak_client")

	if !exists {
		c.AbortWithStatus(500)
		return
	}
	client, ok := keycloakClient.(api_keycloak.KeycloakClient)
	if !ok {
		c.AbortWithStatus(500)
		return
	}
	token := splittedToken[1]
	parsedToken, err := client.VerifyToken(token)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	if parsedToken == nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	err = api_keycloak.VerifyClaims(parsedToken)
	if err != nil {
		err.JsonError(c)
		return
	}
	// pass the token to the next middleware
	c.Set("token", parsedToken)
	// clientId
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	clientId := claims["client_id"].(float64)
	clientIdInt := int(clientId)
	c.Set("client_id", clientIdInt)
	c.Next()
}

func AuthenticateUserByClientIdMiddleware(c *gin.Context, clientId int) {
	claimsClientId := c.GetInt("client_id")
	if claimsClientId != clientId {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	c.Next()
}

func AuthenticationByClientIdHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIDStr := c.Param("client_id")
		clientID, err := strconv.Atoi(clientIDStr)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "Invalid identifier"})
			return
		}
		AuthenticateUserByClientIdMiddleware(c, clientID)
	}
}

func AuthenticateUserByIdentificationMiddleware(c *gin.Context, identification string) {
	claimsClientId := c.GetInt("client_id")
	repositoryWrapper, exists := c.Get("repository_wrapper")

	if !exists {
		fmt.Println("Repository wrapper does not exists!")
		c.AbortWithStatus(500)
		return
	}
	repositories, ok := repositoryWrapper.(*repositories.RepositoryWrapper)
	if !ok {
		fmt.Println("Repository wrapper parse error!")
		c.AbortWithStatus(500)
		return
	}

	client, err := repositories.ClientRepository.FetchClientByIdentification(c.Request.Context(), identification)
	if err != nil {
		fmt.Println("Error fetching client by identification with request context")

		err.JsonError(c)
		return
	}
	if client.ID != claimsClientId {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	c.Next()
}

func AuthenticateUserByIdentificationHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		identification := c.Param("identification")
		AuthenticateUserByIdentificationMiddleware(c, identification)
	}
}

func AuthenticateUserByAccountIdMiddleware(c *gin.Context, accountId int) {
	claimsClientId := c.GetInt("client_id")
	repositoryWrapper, exists := c.Get("repository_wrapper")

	if !exists {
		c.AbortWithStatus(500)
		return
	}
	repositories, ok := repositoryWrapper.(*repositories.RepositoryWrapper)
	if !ok {
		c.AbortWithStatus(500)
		return
	}

	account, err := repositories.AccountRepository.FetchAccountById(c.Request.Context(), accountId)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	if account.ClientID != claimsClientId {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	c.Next()
}

func AuthenticateByAccountIdHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		accountIdStr := c.Param("account_id")
		fmt.Printf(":account_id param => %s\n", accountIdStr)

		accountID, err := strconv.ParseInt(accountIdStr, 0, 32)
		fmt.Printf(":account_id param parsed => %d\n", accountID)

		if err != nil {
			fmt.Println("Error al parsear account id como parámetro de la url")
			c.JSON(400, gin.H{"error": "Invalid identifier"})
			return
		}
		AuthenticateUserByAccountIdMiddleware(c, int(accountID))
	}
}

func AuthenticatePerformTransactionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var performnTransactionDto dto.PerformTransactionDto
		if error := c.ShouldBindJSON(&performnTransactionDto); error != nil {
			fmt.Printf("ERROR AL OBTENER CUENTA %s", error)

			c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
			return
		}
		// Because Gin consumes the body stream once the body is binded, we have to store in the context
		// for further use in other handlers
		c.Set("perform_transaction_dto", performnTransactionDto)
		claimsClientId := c.GetInt("client_id")
		repositoryWrapper, exists := c.Get("repository_wrapper")

		if !exists {
			c.AbortWithStatus(500)
			return
		}
		repositories, ok := repositoryWrapper.(*repositories.RepositoryWrapper)
		if !ok {
			c.AbortWithStatus(500)
			return
		}
		account, err := repositories.AccountRepository.FetchAccountById(c.Request.Context(), performnTransactionDto.AccountID)
		fmt.Printf("FetchAccountById CALLED")

		if err != nil {
			fmt.Printf("ERROR AL OBTENER CUENTA %s", err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Bad Request"})
			return
		}
		fmt.Printf("ACCOUNT CALLED %v\n", account)
		fmt.Printf("Account Client ID: %d \nClaimed Client Id: %d \n", account.ClientID, claimsClientId)
		if account.ClientID != claimsClientId {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		fmt.Println("All authentication checks passed. Calling c.Next()")

		c.Next()
	}

}

func AppMiddlewares() []func() gin.HandlerFunc {
	var middlewares ([]func() gin.HandlerFunc)
	middlewares = append(middlewares, LoggerMiddleware)
	//middlewares = append(middlewares,AuthMiddleware)
	return middlewares

}
