package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nerijusro/scootinAboot/types"
	"github.com/nerijusro/scootinAboot/types/interfaces"
	"github.com/nerijusro/scootinAboot/types/requests"
)

var scooters []types.Scooter

type ScootersHandler struct {
	repository *interfaces.ScootersRepository
}

func NewHandler() *ScootersHandler {
	return &ScootersHandler{}
}

func (h *ScootersHandler) RegisterEndpoints(e *gin.Engine) {
	e.POST("/scooters", createScooter)
	e.GET("/scooters", getScooters)
}

func createScooter(c *gin.Context) {
	//Galima turbut perkelt i validacija.go
	var scooterRequest requests.CreateScooterRequest
	if err := c.BindJSON(&scooterRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request body": err.Error()})
		return
	}

	scooter := types.Scooter{
		ID:          uuid.New(),
		Location:    scooterRequest.Location,
		IsAvailable: scooterRequest.IsAvailable,
		OccupiedBy:  uuid.Nil,
	}

	scooters = append(scooters, scooter)

	c.JSON(http.StatusCreated, scooter)
}

func getScooters(c *gin.Context) {
	c.JSON(http.StatusOK, scooters)
}
