package clientdto

import (api_keycloak "src/api/keycloak")

type AuthorizationResponse = api_keycloak.TokenResponse




type AuthorizationRequest struct {
	User 		string `json:"user" binding:"required"`
	Password	string `json:"password" binding:"required"`
}



