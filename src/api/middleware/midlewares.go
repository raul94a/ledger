package middleware

import (
	"fmt"
	"log"
	api_keycloak "src/api/keycloak"
	"src/repositories"
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
	fmt.Printf("ClientID CLAIMS AUTHENTICATION %d \n", claims["client_id"])
	clientId := claims["client_id"].(float64)
	clientIdInt := int(clientId) 
	c.Set("client_id", clientIdInt)
	c.Next()
}

func AuthenticateUserByClientIdMiddleware(c *gin.Context, clientId int) {
	claimsClientId, exists := c.Get("client_id")
    if !exists {
        c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized: Client ID not found in claims"})
        return
    }

    fmt.Printf("Client Id: %d\n", clientId)
    fmt.Printf("Claims Client Id: %v\n", claimsClientId)

    // Type assertion to int
    
	if claimsClientId != clientId{
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}

    c.Next()
}

func AuthenticateUserByIdentificationMiddleware(c *gin.Context, identification string) {
	claimsClientId := c.GetInt("client_id")
	repositoryWrapper, exists := c.Get("repository_wrapper")

	if !exists {
		c.AbortWithStatus(500)
		return
	}
	repositories, ok := repositoryWrapper.(repositories.RepositoryWrapper)
	if !ok {
		c.AbortWithStatus(500)
		return
	}

	client, err := repositories.ClientRepository.FetchClientByIdentification(c.Request.Context(), identification)
	if err != nil {
		err.JsonError(c)
		return
	}
	if client.ID != claimsClientId {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	c.Next()
}

func AuthenticateUserByAccountIdMiddleware(c *gin.Context, accountId int) {
	claimsClientId := c.GetInt("client_id")
	repositoryWrapper, exists := c.Get("repository_wrapper")

	if !exists {
		c.AbortWithStatus(500)
		return
	}
	repositories, ok := repositoryWrapper.(repositories.RepositoryWrapper)
	if !ok {
		c.AbortWithStatus(500)
		return
	}

	account, err := repositories.AccountRepository.FetchAccountById(c.Request.Context(), accountId)
	if err != nil {
		err.JsonError(c)
		return
	}
	if account.ClientID != claimsClientId {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	c.Next()
}

func AppMiddlewares() []func() gin.HandlerFunc {
	var middlewares ([]func() gin.HandlerFunc)
	middlewares = append(middlewares, LoggerMiddleware)
	//middlewares = append(middlewares,AuthMiddleware)

	return middlewares

}
