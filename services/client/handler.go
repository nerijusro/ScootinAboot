package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nerijusro/scootinAboot/types"
	"github.com/nerijusro/scootinAboot/types/interfaces"
)

type ClientHandler struct {
	repository  interfaces.ClientsRepository
	authService interfaces.AuthService
}

func NewClientsHandler(repository interfaces.ClientsRepository, authService interfaces.AuthService) *ClientHandler {
	return &ClientHandler{repository: repository, authService: authService}
}

func (h *ClientHandler) RegisterEndpoints(e *gin.Engine) {
	e.POST("/users", h.createUser)
}

func (h *ClientHandler) createUser(c *gin.Context) {
	if !h.authService.AuthenticateUser(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

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
