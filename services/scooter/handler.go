package scooter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nerijusro/scootinAboot/types"
)

type ScooterHandler struct {
	repository  types.IScootersRepository
	authService types.IAuthService
	validator   types.IScootersValidator
}

func NewScootersHandler(repository types.IScootersRepository, authService types.IAuthService, validator types.IScootersValidator) *ScooterHandler {
	return &ScooterHandler{repository: repository, authService: authService, validator: validator}
}

func (h *ScooterHandler) RegisterEndpoints(e *gin.Engine) {
	e.POST("/admin/scooters", h.createScooter)
	e.GET("/admin/scooters", h.getAllScooters)
	e.GET("/scooters", h.getScootersByArea)
	e.GET("/scooters/:id", h.getScooter)
}

func (h *ScooterHandler) createScooter(c *gin.Context) {
	if !h.authService.AuthenticateAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var scooterRequest types.CreateScooterRequest
	if err := c.BindJSON(&scooterRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request body": err.Error()})
		return
	}

	if err := h.validator.ValidateCreateScooterRequest(&scooterRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request": err.Error()})
		return
	}

	scooter := types.Scooter{
		ID:          uuid.New(),
		Location:    scooterRequest.Location,
		IsAvailable: scooterRequest.IsAvailable,
	}

	err := h.repository.CreateScooter(scooter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, scooter)
}

func (h *ScooterHandler) getScootersByArea(c *gin.Context) {
	if !h.authService.AuthenticateUser(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var queryParameters types.GetScootersQueryParameters
	if err := c.BindQuery(&queryParameters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request": err.Error()})
		return
	}

	if err := h.validator.ValidateGetScootersQueryParameters(&queryParameters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request": err.Error()})
		return
	}

	scooters, err := h.repository.GetScootersByArea(queryParameters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, scooters)
}

func (h *ScooterHandler) getAllScooters(c *gin.Context) {
	if !h.authService.AuthenticateAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	allScooters, err := h.repository.GetAllScooters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, allScooters)
}

func (h *ScooterHandler) getScooter(c *gin.Context) {
	if !h.authService.AuthenticateAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	idInUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request": err.Error()})
		return
	}

	scooter, _, err := h.repository.GetScooterById(idInUUID.String())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not found": err.Error()})
		return
	}

	c.JSON(http.StatusOK, scooter)
}
