package trip

import "github.com/gin-gonic/gin"

type TripHandler struct {
}

func NewTripHandler() *TripHandler {
	return &TripHandler{}
}

func (h *TripHandler) RegisterEndpoints(e *gin.Engine) {
	e.POST("/trips/start", h.startTrip)
	e.PUT("/trips/update", h.updateTrip)
	e.GET("/trips/finish", h.finishTrip)
}

func (h *TripHandler) startTrip(c *gin.Context) {
	// implementation of the startTrip method
}

func (h *TripHandler) updateTrip(c *gin.Context) {
	// implementation of the updateTrip method
}

func (h *TripHandler) finishTrip(c *gin.Context) {
	// implementation of the finishTrip method
}
