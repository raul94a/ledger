package handlers

import (
	"net/http"
	api_keycloak "src/api/keycloak"
	dto "src/api/dto"
	"github.com/gin-gonic/gin"
)

type AuthorizationHandler interface {
	Authorization(c *gin.Context)
	
}

type IAuthorizationHandler struct {
	KeycloakClient api_keycloak.KeycloakClient
}


func (h *IAuthorizationHandler) Authorization(c *gin.Context) {

	var req dto.AuthorizationRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	token, err := h.KeycloakClient.AuthUser(req.User,req.Password)
	if err != nil {
		err.JsonError(c)
		return
	}

	c.JSON(http.StatusOK,token)

}
