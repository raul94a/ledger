package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	api_keycloak "src/api/keycloak"
	"strings"
	"time"
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

func AuthenticationMiddleware(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	splittedToken := strings.Split(accessToken, " ")
	if splittedToken[0] != "Bearer" {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	token := splittedToken[1]
	parsedToken, err := api_keycloak.VerifyToken(token)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	if parsedToken == nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	err = api_keycloak.VerifyClaims(parsedToken)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	// pass the token to the next middleware
	c.Set("token",parsedToken)
	c.Next()

}

func AppMiddlewares() []func() gin.HandlerFunc {
	var middlewares ([]func() gin.HandlerFunc)
	middlewares = append(middlewares, LoggerMiddleware)
	//middlewares = append(middlewares,AuthMiddleware)

	return middlewares

}
