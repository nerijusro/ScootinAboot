package scooter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nerijusro/scootinAboot/types"
)

type ScootersHandler struct {
	repository types.IScootersRepository
}

func NewScootersHandler(repository types.IScootersRepository) *ScootersHandler {
	return &ScootersHandler{repository: repository}
}

func (h *ScootersHandler) RegisterEndpoints(e *gin.Engine) {
	e.POST("/scooters", h.createScooter)
	e.GET("/scooters", h.getScooters)
	e.GET("/scooters/:id", h.getScooter)
}

func (h *ScootersHandler) createScooter(c *gin.Context) {
	//Galima turbut perkelt i validacija.go
	var scooterRequest types.CreateScooterRequest
	if err := c.BindJSON(&scooterRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request body": err.Error()})
		return
	}

	scooter := types.Scooter{
		ID:          uuid.New(),
		Location:    scooterRequest.Location,
		IsAvailable: scooterRequest.IsAvailable,
	}

	err := h.repository.CreateScooter(scooter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err})
		return
	}

	c.JSON(http.StatusCreated, scooter)
}

func (h *ScootersHandler) getScooters(c *gin.Context) {
	allScooters, err := h.repository.GetScooters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err})
		return
	}

	c.JSON(http.StatusOK, allScooters)
}

func (h *ScootersHandler) getScooter(c *gin.Context) {
	id := c.Param("id")
	idInUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request": err})
		return
	}

	scooter, err := h.repository.GetScooterById(idInUUID.String())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not found": err.Error()})
		return
	}

	c.JSON(http.StatusOK, scooter)
}
