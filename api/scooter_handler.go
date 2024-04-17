package api

import (
	"ScootinAboot/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var scooters []types.Scooter

func AddScooterHandlerEndpoints(e *gin.Engine) {
	e.POST("/scooters", createScooter)
	e.GET("/scooters", getScooters)
}

func createScooter(c *gin.Context) {
	var scooter types.Scooter
	scooter.ID = uuid.New()

	scooters = append(scooters, scooter)

	c.JSON(http.StatusCreated, scooter)
}

func getScooters(c *gin.Context) {
	c.JSON(http.StatusOK, scooters)
}
