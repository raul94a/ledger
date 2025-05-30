package jwt_test

import (
	// "context"
	"fmt"
	api_keycloak "src/api/keycloak"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	// "github.com/stillya/testcontainers-keycloak"
	// "github.com/testcontainers/testcontainers-go"
	// "github.com/testcontainers/testcontainers-go/wait"
)

// Clave JWK de firma (copiada de tu input)

// Token JWT de Keycloak (reemplaza con el token real)
// Test para validar un token JWT de Keycloak
func TestValidateKeycloakToken(t *testing.T) {
	//  ctx := context.Background()

    // // Configure Keycloak container
    // keycloakContainer, err := keycloak.Run(ctx,
    //     "quay.io/keycloak/keycloak:24.0", // Keycloak image and version
    //     testcontainers.WithEnv(map[string]string{
    //         "KEYCLOAK_HTTP_PORT": "7080", // Set Keycloak to listen on port 7080
    //         "KEYCLOAK_FEATURES":  "script",      // Enable the script feature

    //     }),
    //     keycloak.WithAdminUsername("admin"),                // Set admin username
    //     keycloak.WithAdminPassword("admin"),                // Set admin password
    //     keycloak.WithRealmImportFile("../ledger_realm_test.json"), // Import realm from JSON file
    //     testcontainers.WithExposedPorts("7080/tcp"),        // Expose port 7080
    //    testcontainers.WithWaitStrategy(
    //     wait.ForLog("Running the server in development mode. DO NOT use this configuration in production"),
    //    ),
    // )
    // if err != nil {
    //     t.Fatalf("Failed to start Keycloak container: %v", err)
    // }

    // // Ensure container is terminated after use
    // defer func() {
    //     if err := keycloakContainer.Terminate(ctx); err != nil {
    //         t.Fatalf("Failed to terminate Keycloak container: %v", err)
    //     }
    // }()
    godotenv.Load("../../.env")
	client := api_keycloak.BuildKeycloakClientFromEnv()
	// Parsear la clave JWK
	tokenResponse, _ := client.AuthAdminUser()
	tokenString := tokenResponse.AccessToken
	jwks, _ := client.GetJwkCerts()
	jwkKey, _ := jwks.GetSigJwk()
	rsaKey,_ := jwkKey.ComputePublicRsaKey()
   
	// Parsear y validar el token
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar el algoritmo de firma
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Verificar el kid
		kid, ok := token.Header["kid"].(string)
		if !ok || kid != jwkKey.Kid {
			return nil, fmt.Errorf("invalid or missing kid: expected %s, got %s", jwkKey.Kid, kid)
		}

		return &rsaKey, nil
	})

	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if !parsedToken.Valid {
		t.Fatal("Token is not valid")
	}

	// Verificar claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Invalid claims")
	}
	fmt.Println("Claims: ", claims)
	// Validar el issuer
	expectedIssuer := "http://localhost:7080/realms/ledger" // Reemplaza con tu issuer
	if claims["iss"] != expectedIssuer {
		t.Errorf("Expected issuer %s, got %s", expectedIssuer, claims["iss"])
	}

	// Validar la expiración
	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatal("Invalid exp claim")
	}
	if exp < float64(time.Now().Unix()) {
		t.Error("Token is expired")
	}

	// Validar el subject
	if sub, ok := claims["sub"].(string); !ok || sub == "" {
		t.Error("Invalid or missing sub claim")
	}

	// Opcional: Verificar roles (si Keycloak incluye roles en el token)
	if resourceAccess, ok := claims["resource_access"].(map[string]interface{}); ok {
		clientRoles, ok := resourceAccess["<client-id>"].(map[string]interface{})
		if ok {
			roles, ok := clientRoles["roles"].([]interface{})
			if !ok {
				t.Error("No roles found in token")
			} else {
				t.Logf("Roles found: %v", roles)
			}
		}
	}

	fmt.Println("Parsed token: ", parsedToken)

	// Test con un token inválido (por ejemplo, manipulado)
	invalidToken := tokenString + "tampered"
	_, err = jwt.Parse(invalidToken, func(token *jwt.Token) (interface{}, error) {
		return rsaKey, nil
	})
	if err == nil {
		t.Fatal("Expected error for tampered token, got none")
	}
}
