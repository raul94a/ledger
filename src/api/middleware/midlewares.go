package middleware

import (
	"log"
	api_keycloak "src/api/keycloak"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

func AuthenticationMiddleware(c *gin.Context) {
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
		err.JsonError(c)
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
	c.Next()

}

func AppMiddlewares() []func() gin.HandlerFunc {
	var middlewares ([]func() gin.HandlerFunc)
	middlewares = append(middlewares, LoggerMiddleware)
	//middlewares = append(middlewares,AuthMiddleware)

	return middlewares

}
