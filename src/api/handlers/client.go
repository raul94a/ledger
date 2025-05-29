package handlers

import (
	"fmt"
	"net/http"
	clientdto "src/api/dto"
	services "src/api/service"
	repositories "src/repositories"

	"github.com/gin-gonic/gin"
)

type ClientHandler interface {
	CreateClient(c *gin.Context)
}

type IClientHandler struct {
	ClientRepository             repositories.ClientRepository
	RegistryAccountOtpRepository repositories.RegistryAccountOtpRepository
	ClientService				 services.ClientService
}

// @Summary Crea un nuevo cliente
// @Description Recibe los datos de un cliente y lo registra en el sistema.
// @Accept json
// @Produce json
// @Param client body CreateClientRequest true "Datos del cliente para crear"
// @Success 201 {object} map[string]interface{} "Cliente creado exitosamente"
// @Failure 400 {object} map[string]string "Solicitud inválida"
// @Router /clients [post]
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
