package keycloak_test

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"src/api/keycloak"
	"testing"
)

func TestAdminAuthHandlerTest(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file ", err.Error())
	}
	client := api_keycloak.BuildKeycloakClientFromEnv()

	_, err = client.AuthAdminUser()
	if err != nil {
		t.Log(err)
	}

	assert.True(t, err == nil)
}

func TestKeyCloakJWK(t *testing.T) {
	client := api_keycloak.BuildKeycloakClientFromEnv()

	jwk, err := client.GetJwkCerts()
	if err != nil {
		t.Log("ERROR ", err)

		return
	}
	assert.Equal(t, "sig", jwk.Keys[0].Use)

}
