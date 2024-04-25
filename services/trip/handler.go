package trip

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nerijusro/scootinAboot/types"
	"github.com/nerijusro/scootinAboot/types/enums"
	"github.com/nerijusro/scootinAboot/types/interfaces"
)

type TripHandler struct {
	validator          interfaces.TripValidator
	tripsReposiotry    interfaces.TripRepository
	scootersRepository interfaces.ScooterRepository
	usersRepository    interfaces.ClientRepository
}

func NewTripHandler(
	validator interfaces.TripValidator,
	tripsRepository interfaces.TripRepository,
	scootersRepository interfaces.ScooterRepository,
	usersRepository interfaces.ClientRepository) *TripHandler {
	return &TripHandler{validator: validator, tripsReposiotry: tripsRepository, scootersRepository: scootersRepository, usersRepository: usersRepository}
}

func (h *TripHandler) RegisterEndpoints(routerGroups map[string]*gin.RouterGroup) {
	userAuthorized := routerGroups["client"]
	userAuthorized.POST("/trips", h.startTrip)
	userAuthorized.PUT("/trips/:id", h.updateTrip)
}

func (h *TripHandler) startTrip(c *gin.Context) {
	clientId := c.GetHeader("Client-Id")
	_, err := uuid.Parse(clientId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Unauthorized request": err.Error()})
		return
	}

	var startTripRequest types.StartTripRequest
	if err := c.BindJSON(&startTripRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request body": err.Error()})
		return
	}

	if err := h.validator.ValidateStartTripRequest(&startTripRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request": err.Error()})
		return
	}

	scooter, scooterOptLockVersion, err := h.scootersRepository.GetScooterById(startTripRequest.ScooterID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error(), "message": "error getting scooter by id"})
		return
	}

	user, userOptLockVersion, err := h.usersRepository.GetUserById(clientId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error(), "message": "error getting user by id"})
		return
	}

	if err := h.validateTripStart(scooter, user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request": err.Error(), "message": "trip cannot be started due to parameter invalidity"})
		return
	}

	trip := types.Trip{
		ID:        uuid.New(),
		ScooterId: scooter.ID,
		ClientId:  user.ID,
	}

	startTripEvent := types.TripEvent{
		TripID:    trip.ID,
		Type:      enums.StartTrip,
		Location:  scooter.Location,
		CreatedAt: startTripRequest.CreatedAt,
		Sequence:  1,
	}

	err = h.tripsReposiotry.StartTrip(trip, scooterOptLockVersion, userOptLockVersion, startTripEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error(), "message": "trip could not be started"})
		return
	}

	c.JSON(http.StatusCreated, startTripEvent)
}

func (h *TripHandler) updateTrip(c *gin.Context) {
	tripId := c.Param("id")
	_, err := uuid.Parse(tripId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request": err.Error()})
		return
	}

	clientId := c.GetHeader("Client-Id")
	_, err = uuid.Parse(clientId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Unauthorized request": err.Error()})
		return
	}

	var request types.TripUpdateRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request body": err.Error()})
		return
	}

	if err := h.validator.ValidateTripUpdateRequest(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request": err.Error()})
		return
	}

	trip, err := h.tripsReposiotry.GetTripById(tripId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error(), "message": "error getting trip by id"})
		return
	}

	if trip.ClientId.String() != clientId {
		c.JSON(http.StatusUnauthorized, gin.H{"Unauthorized request": "user is not allowed to update this trip"})
		return
	}

	if trip.IsFinished {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request": "trip is already finished"})
		return
	}

	_, scooterOptLockVersion, err := h.scootersRepository.GetScooterById(trip.ScooterId.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error(), "message": "error getting scooter by id"})
		return
	}

	tripEvent := types.TripEvent{
		TripID:    trip.ID,
		Location:  request.Location,
		CreatedAt: request.CreatedAt,
		Sequence:  request.Sequence,
	}

	if !request.IsFinishing {
		tripEvent.Type = enums.UpdateTrip
		err = h.tripsReposiotry.UpdateTrip(trip, scooterOptLockVersion, tripEvent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error(), "message": "trip could not be updated"})
			return
		}
	} else {
		_, userOptLockVersion, err := h.usersRepository.GetUserById(clientId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error(), "message": "error getting user by id"})
			return
		}

		tripEvent.Type = enums.EndTrip
		err = h.tripsReposiotry.EndTrip(trip, scooterOptLockVersion, userOptLockVersion, tripEvent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error(), "message": "trip could not be finished"})
			return
		}
	}

	c.JSON(http.StatusOK, tripEvent)
}

func (h *TripHandler) validateTripStart(scooter *types.Scooter, user *types.MobileClient) error {
	if !scooter.IsAvailable {
		return errors.New("scooter is not available")
	}

	if !user.IsEligibleToTravel {
		return errors.New("user is not eligible to travel")
	}

	return nil
}
