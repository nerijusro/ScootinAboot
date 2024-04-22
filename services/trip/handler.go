package trip

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nerijusro/scootinAboot/types"
)

type TripHandler struct {
	authService types.IAuthService
}

func NewTripHandler(authService types.IAuthService) *TripHandler {
	return &TripHandler{authService: authService}
}

func (h *TripHandler) RegisterEndpoints(e *gin.Engine) {
	e.POST("/trips/start", h.startTrip)
	e.PUT("/trips/update", h.updateTrip)
	e.GET("/trips/finish", h.finishTrip)
}

// Tripas turi:
// - id
// - scooterId
// - clientId
// - isFinished
func (h *TripHandler) startTrip(c *gin.Context) {
	if !h.authService.AuthenticateUser(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

}

func (h *TripHandler) updateTrip(c *gin.Context) {
	// implementation of the updateTrip method
}

func (h *TripHandler) finishTrip(c *gin.Context) {
	// implementation of the finishTrip method
}
