package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nerijusro/scootinAboot/types"
	"github.com/nerijusro/scootinAboot/types/interfaces"
)

type ClientHandler struct {
	repository interfaces.ClientRepository
}

func NewClientsHandler(repository interfaces.ClientRepository) *ClientHandler {
	return &ClientHandler{repository: repository}
}

func (h *ClientHandler) RegisterEndpoints(routerGroups map[string]*gin.RouterGroup) {
	userAuthorized := routerGroups["client"]
	userAuthorized.POST("/users", h.createUser)
}

func (h *ClientHandler) createUser(c *gin.Context) {
	var userRequest types.CreateUserRequest
	if err := c.BindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request body": err.Error()})
		return
	}

	user := types.MobileClient{
		ID:                 uuid.New(),
		FullName:           userRequest.FullName,
		IsEligibleToTravel: true,
	}

	err := h.repository.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}
