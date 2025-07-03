package appRedis

import (
	"context"
	"encoding/json"
	"fmt"
	api_keycloak "src/api/keycloak"

	"github.com/redis/go-redis/v9"
)

type GetTokenResponseType interface{
	GetTypeMessage() string
}
type LoginAgainResponse struct {
	Message string
}
type OkResponse struct {
	Message string
}
type ErrorResponse struct {
	Message string
}
func (c LoginAgainResponse) GetTypeMessage() string {
	return c.Message
}
func (c ErrorResponse) GetTypeMessage() string {
	return c.Message
}
func (c OkResponse) GetTypeMessage() string {
	return c.Message
}

// GenerateKeycloakTokenKey creates a Redis key for a Keycloak token.
func GenerateKeycloakTokenKey(sessionId string) string {
	return fmt.Sprintf("token:%s", sessionId)
}

// SetToken stores a TokenResponse struct as a JSON document in Redis.
func SetToken(ctx context.Context, rdb *redis.Client, sessionId string, token *api_keycloak.TokenResponse) error {

	key := GenerateKeycloakTokenKey(sessionId)
	_, err := rdb.JSONSet(ctx, key, "$", token).Result()
	if err != nil {
		return fmt.Errorf("failed to save document for key %s: %w", key, err)
	}
	fmt.Printf("Guardando session ID en redis token:%s",sessionId)
	return nil
}

// GetToken retrieves a TokenResponse document from Redis.
func GetToken(ctx context.Context, rdb *redis.Client, sessionId string) (*api_keycloak.TokenResponse, GetTokenResponseType) {
	
	key := GenerateKeycloakTokenKey(sessionId)
	res, err := rdb.JSONGet(ctx, key, "$").Result()
	if err != nil || res=="" {
		errString := fmt.Errorf("failed to get Keycloak token document for key %s: %w", key, err)
		if err == redis.Nil {
			return nil, LoginAgainResponse{Message: errString.Error()}
		}
		return nil, LoginAgainResponse{Message: errString.Error()}
	}
	fmt.Println("RESPUESTA DE REDIS " + res)
	var tokens []api_keycloak.TokenResponse // JSONGet with "$" returns an array of the top-level object
	err = json.Unmarshal([]byte(res), &tokens)
	if err != nil {
		errString := fmt.Sprintf("failed to unmarshal Keycloak token document for key %s: %s", key, err.Error())
		return nil, ErrorResponse{Message: errString}
	}
	if len(tokens) == 0 {
		return nil, LoginAgainResponse{Message: ""} // No token found or empty array
	}
	return &tokens[0], OkResponse{Message: ""} // Return the first (and only) token
}

// DeleteToken deletes a Keycloak token document from Redis.
func DeleteToken(ctx context.Context, rdb *redis.Client, sessionId string) error {
	
	key := GenerateKeycloakTokenKey(sessionId)
	_, err := rdb.Del(ctx, key).Result()
	if err != nil {
		fmt.Printf("failed to delete Keycloak token document for key %s: %w", key, err)
		return fmt.Errorf("failed to delete Keycloak token document for key %s: %w", key, err)
	}
	return nil
}
