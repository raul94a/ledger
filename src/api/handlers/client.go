package handlers

import(
	mappers "src/mappers"
	"net/http"
	clientdto "src/api/dto"
	repositories "src/repositories"
	"github.com/gin-gonic/gin"
)




type ClientHandler interface {
	CreateClient(c *gin.Context)
}


type IClientHandler struct {
	ClientRepository repositories.ClientRepository
}

// @Summary Crea un nuevo cliente
// @Description Recibe los datos de un cliente y lo registra en el sistema.
// @Accept json
// @Produce json
// @Param client body CreateClientRequest true "Datos del cliente para crear"
// @Success 201 {object} map[string]interface{} "Cliente creado exitosamente"
// @Failure 400 {object} map[string]string "Solicitud inválida"
// @Router /clients [post]
func (h *IClientHandler) CreateClient(c *gin.Context){
	var client clientdto.CreateClientRequest
	if error := c.ShouldBindJSON(&client); error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}

	clientEntity, error := mappers.ToClientEntity(client)
	if(error != nil){
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}
	context := c.Request.Context()
	err := h.ClientRepository.InsertClient(context, &clientEntity)
	if err != nil {
		err.JsonError(c)
		return
	}
	// Si el parsing fue exitoso, 'client' ahora contiene los datos del JSON
	clientResponse, error := mappers.ToClientDTO(clientEntity)
	// Aquí podrías guardar el usuario en una base de datos, realizar validaciones adicionales, etc.
	// Por simplicidad, solo devolvemos un mensaje de éxito.
	if error != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	// CreateClientRequest pasa a ClientEntity

	c.JSON(http.StatusOK, gin.H{
		"message": "Cliente creado exitosamente",
		"client":  clientResponse,
	})
}