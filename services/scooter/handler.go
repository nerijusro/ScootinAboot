package scooter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nerijusro/scootinAboot/types"
	"github.com/nerijusro/scootinAboot/types/interfaces"
)

type ScooterHandler struct {
	repository interfaces.ScootersRepository
	validator  interfaces.ScootersValidator
}

func NewScootersHandler(repository interfaces.ScootersRepository, validator interfaces.ScootersValidator) *ScooterHandler {
	return &ScooterHandler{repository: repository, validator: validator}
}

func (h *ScooterHandler) RegisterEndpoints(routerGroups map[string]*gin.RouterGroup) {
	adminAuthorized := routerGroups["admin"]
	adminAuthorized.POST("/scooters", h.createScooter)
	adminAuthorized.GET("/scooters", h.getAllScooters)

	userAuthorized := routerGroups["client"]
	userAuthorized.GET("/scooters", h.getScootersByArea)
	userAuthorized.GET("/scooters/:id", h.getScooter)
}

func (h *ScooterHandler) createScooter(c *gin.Context) {
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
		c.JSON(http.StatusNotFound, gin.H{"Not Found": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"scooters": scooters})
}

func (h *ScooterHandler) getAllScooters(c *gin.Context) {
	allScooters, err := h.repository.GetAllScooters()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not Found": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"scooters": allScooters})
}

func (h *ScooterHandler) getScooter(c *gin.Context) {
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
