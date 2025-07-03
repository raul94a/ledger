package appRedis

import (
	"fmt"
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)



// IndexName constant for the tokens index
const (
	TokensIndex = "idx:tokens"
	// ... other index names
)

// CreateAllIndexes should be called on application startup
func CreateAllIndexes(ctx context.Context, rdb *redis.Client, zLogger *zap.Logger) error {
	// ... existing index creation calls (e.g., CreateUsersIndex)

	err := CreateTokensIndex(ctx, rdb,zLogger)
	if err != nil {
		return fmt.Errorf("failed to create tokens index: %w", err)
	}
	zLogger.Sugar().Info(fmt.Sprintf("RedisSerch index %s ensured",TokensIndex))

	return nil
}

// CreateTokensIndex defines and creates the 'idx:keycloak_tokens' index for TokenResponse.
func CreateTokensIndex(ctx context.Context, rdb *redis.Client,zLogger *zap.Logger) error {
	// Check if index already exists to prevent errors on re-creation
	_, err := rdb.FTInfo(ctx, TokensIndex).Result()
	if err == nil {
		//log.Printf("RediSearch index '%s' already exists, skipping creation.", TokensIndex)
		return nil
	}
	// Check for the specific "Unknown Index Name" error, otherwise it's a real error
	if err.Error() != "Unknown Index Name" {
		//return fmt.Errorf("failed to check existing index '%s': %w", TokensIndex, err)
	}

	// Define the schema fields for the TokenResponse struct
	schema := []*redis.FieldSchema{
		{
			FieldName: "$.access_token",
			As:        "accessToken",
			FieldType: redis.SearchFieldTypeText,
			// You might use Tag if you want exact matches and case-insensitivity on this,
			// but for long tokens, Text is fine for basic indexing, or no indexing if not searchable.
			// Options: &redis.FieldOptions{NoIndex: true}, // If you don't need to search by access_token
		},
		{
			FieldName: "$.expires_in",
			As:        "expiresIn",
			FieldType: redis.SearchFieldTypeNumeric,
		},
		{
			FieldName: "$.refresh_expires_in",
			As:        "refreshExpiresIn",
			FieldType: redis.SearchFieldTypeNumeric,
		},
		{
			NoIndex: true,
			FieldName: "$.refresh_token",
			As:        "refreshToken",
			FieldType: redis.SearchFieldTypeText,
			// Options: &redis.FieldOptions{NoIndex: true}, // If you don't need to search by refresh_token
		},
		{
			FieldName: "$.token_type",
			As:        "tokenType",
			FieldType: redis.SearchFieldTypeTag, // "Bearer" can be a good tag field
		},
		{
			FieldName: "$.not-before-policy", // Note the dash in JSON tag
			As:        "notBeforePolicy",
			FieldType: redis.SearchFieldTypeNumeric,
		},
		{
			FieldName: "$.session_state",
			As:        "sessionState",
			FieldType: redis.SearchFieldTypeTag, // Session state can be good for exact matches
		},
		{
			FieldName: "$.scope",
			As:        "scope",
			FieldType: redis.SearchFieldTypeTag, // Scopes are often space-separated strings, Tag allows finding exact tokens (e.g., "openid profile")
		},
	}

	// Options for the index (OnJSON: true for JSON documents, Prefix for keys)
	options := &redis.FTCreateOptions{
		OnJSON: true,
		// Choose a prefix that identifies your Keycloak tokens.
		// For example, if you store them as "token:session_id".
		// Let's assume you'll store them like "token:<user_id_or_session_id>"
		Prefix: []interface{}{"token:"},
	}

	// Create the index
	_, err = rdb.FTCreate(ctx, TokensIndex, options, schema...).Result()
	if err != nil {
		errorStr := fmt.Sprintf("failed to create RediSearch index '%s': %s", TokensIndex, err.Error())
		zLogger.Error(errorStr)
		return err
	}
	return nil
}