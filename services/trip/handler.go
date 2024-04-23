package trip

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nerijusro/scootinAboot/types"
	"github.com/nerijusro/scootinAboot/utils"
)

type TripHandler struct {
	authService        types.IAuthService
	validator          types.ITripsValidator
	tripsReposiotry    types.ITripsRepository
	scootersRepository types.IScootersRepository
	usersRepository    types.IClientsRepository
}

func NewTripHandler(
	authService types.IAuthService,
	validator types.ITripsValidator,
	tripsRepository types.ITripsRepository,
	scootersRepository types.IScootersRepository,
	usersRepository types.IClientsRepository) *TripHandler {
	return &TripHandler{authService: authService, validator: validator, tripsReposiotry: tripsRepository, scootersRepository: scootersRepository, usersRepository: usersRepository}
}

func (h *TripHandler) RegisterEndpoints(e *gin.Engine) {
	e.POST("/trips/start", h.startTrip)
	e.PUT("/trips/status", h.updateTrip)
	e.PUT("/trips/finish", h.finishTrip)
}

func (h *TripHandler) startTrip(c *gin.Context) {
	if !h.authService.AuthenticateUser(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
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

	clientId, err := utils.GetClientIdFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Unauthorized request": err.Error()})
		return
	}

	scooter, scooterOptLockVersion, err := h.scootersRepository.GetScooterById(startTripRequest.ScooterID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error(), "message": "error getting scooter by id"})
		return
	}

	user, userOptLockVersion, err := h.usersRepository.GetUserById(clientId.String())
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
		Type:      "start_trip_event",
		Location:  scooter.Location,
		CreatedAt: startTripRequest.CreatedAt,
		Sequence:  1,
	}

	err = h.tripsReposiotry.StartTrip(trip, scooterOptLockVersion, userOptLockVersion, startTripEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error(), "message": "trip could not be started"})
		return
	}

	c.JSON(http.StatusOK, startTripEvent)
}

func (h *TripHandler) updateTrip(c *gin.Context) {
	if !h.authService.AuthenticateUser(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var updateTripRequest types.TripUpdateRequest
	if err := c.BindJSON(&updateTripRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request body": err.Error()})
		return
	}

	if err := h.validator.ValidateTripUpdateRequest(&updateTripRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request": err.Error()})
		return
	}

	clientId, err := utils.GetClientIdFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Unauthorized request": err.Error()})
		return
	}

	trip, err := h.tripsReposiotry.GetTripById(updateTripRequest.TripID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error(), "message": "error getting trip by id"})
		return
	}

	if trip.ClientId.String() != clientId.String() {
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
		Type:      "update_trip_event",
		Location:  updateTripRequest.Location,
		CreatedAt: updateTripRequest.CreatedAt,
		Sequence:  updateTripRequest.Sequence,
	}

	err = h.tripsReposiotry.UpdateTrip(trip, scooterOptLockVersion, tripEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error(), "message": "trip could not be updated"})
		return
	}

	c.JSON(http.StatusOK, tripEvent)
}

func (h *TripHandler) finishTrip(c *gin.Context) {
	if !h.authService.AuthenticateUser(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var finishTripRequest types.TripUpdateRequest
	if err := c.BindJSON(&finishTripRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request body": err.Error()})
		return
	}

	if err := h.validator.ValidateTripUpdateRequest(&finishTripRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bad request": err.Error()})
		return
	}

	clientId, err := utils.GetClientIdFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Unauthorized request": err.Error()})
		return
	}

	trip, err := h.tripsReposiotry.GetTripById(finishTripRequest.TripID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error(), "message": "error getting trip by id"})
		return
	}

	if trip.ClientId != clientId {
		c.JSON(http.StatusUnauthorized, gin.H{"Unauthorized request": "user is not allowed to finish this trip"})
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

	_, userOptLockVersion, err := h.usersRepository.GetUserById(clientId.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error(), "message": "error getting user by id"})
		return
	}

	tripEvent := types.TripEvent{
		TripID:    trip.ID,
		Type:      "finish_trip_event",
		Location:  finishTripRequest.Location,
		CreatedAt: finishTripRequest.CreatedAt,
		Sequence:  finishTripRequest.Sequence,
	}

	err = h.tripsReposiotry.EndTrip(trip, scooterOptLockVersion, userOptLockVersion, tripEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error(), "message": "trip could not be finished"})
		return
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
