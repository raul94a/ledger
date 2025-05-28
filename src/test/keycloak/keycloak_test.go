package keycloak_test

import (
	"src/api/keycloak"
	"testing"
	"github.com/stretchr/testify/assert"
	    "github.com/joho/godotenv"

)


func TestAdminAuthHandlerTest(t *testing.T){
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file ",err.Error())
	}
	_, err = api_keycloak.AuthAdminUser()
	if err != nil {
		t.Log(err)
	}
	
	assert.True(t,err == nil)
}

func TestKeyCloakJWK(t *testing.T){
	jwk, err := api_keycloak.GetJwkCerts()
	if err != nil {
		t.Log("ERROR ", err)
		
		return
	}
	assert.Equal(t,"sig",jwk.Keys[0].Use)

}