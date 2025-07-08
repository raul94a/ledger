package handlers

import (
	"context"
	"fmt"
	"net/http"
	clientdto "src/api/dto"
	services "src/api/service"
	repositories "src/repositories"

	"github.com/gin-gonic/gin"
)

type ClientHandler interface {
	CreateClient(c *gin.Context)
	GetClientByIdentification(c *gin.Context)
}

type IClientHandler struct {
	ClientRepository             repositories.ClientRepository
	RegistryAccountOtpRepository repositories.RegistryAccountOtpRepository
	ClientService                services.ClientService
}

// @Summary Creates a new client
// @Description Recibe los datos de un cliente y lo registra en el sistema.
// @Accept json
// @Produce json
// @Param client body clientdto.CreateClientRequest true "Datos del cliente para crear"
// @Success 201 {object} clientdto.ClientResponse ""
// @Failure 400 {object} app_errors.ErrorMessageJsonType ""
// @Failure 404 {object} app_errors.ErrorJsonType ""
// @Router /clients [post]
// @tags Clients
func (h *IClientHandler) CreateClient(c *gin.Context) {
	var client clientdto.CreateClientRequest
	if error := c.ShouldBindJSON(&client); error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}

	clientResponse, err := h.ClientService.CreateClient(client)
	if err != nil {
		fmt.Println("Error ", err.Error())
		err.JsonError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Cliente creado exitosamente",
		"client":  clientResponse,
	})
}

// @Summary Gets a new client by its identification -Spanish National Document Identity-
// @Description Returns the bank client.
// @Accept json
// @Produce json
// @Param identification path string true "Identification of the client with format 00000000X"
// @Success 201 {object} clientdto.ClientResponse ""
// @Failure 400 {object} app_errors.ErrorMessageJsonType ""
// @Failure 404 {object} app_errors.ErrorJsonType ""
// @Router /clients/{identification} [get]
// @security BearerAuth
// @tags Clients
func (h *IClientHandler) GetClientByIdentification(c *gin.Context) {
	identification := c.Param("identification")
	clientResponse, err := h.ClientRepository.FetchClientByIdentification(context.Background(), identification)
	if err != nil {
		fmt.Println("Error ", err.Error())
		err.JsonError(c)
		return
	}
	c.JSON(http.StatusOK, clientResponse)
}
