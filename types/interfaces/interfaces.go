package interfaces

import (
	"github.com/gin-gonic/gin"
	"github.com/nerijusro/scootinAboot/types"
)

type EndpointHandler interface {
	RegisterEndpoints(e *gin.Engine)
}

type AuthService interface {
	AuthenticateAdmin(c *gin.Context) bool
	AuthenticateUser(c *gin.Context) bool
}

type AuthProvider interface {
	GetAdminApiKey() string
	GetUserApiKey() string
}

type ScootersRepository interface {
	GetScooterById(id string) (*types.Scooter, *int, error)
	GetAllScooters() ([]*types.Scooter, error)
	GetScootersByArea(queryParams types.GetScootersQueryParameters) ([]*types.Scooter, error)
	CreateScooter(scooter types.Scooter) error
}

type ClientsRepository interface {
	CreateUser(client types.MobileClient) error
	GetUserById(id string) (*types.MobileClient, *int, error)
}

type TripsRepository interface {
	GetTripById(id string) (*types.Trip, error)
	StartTrip(trip types.Trip, scooterOptLockVersion *int, userOptLockVersion *int, event types.TripEvent) error
	UpdateTrip(trip *types.Trip, scooterOptLockVersion *int, event types.TripEvent) error
	EndTrip(trip *types.Trip, scooterOptLockVersion *int, userOptLockVersion *int, event types.TripEvent) error
}

type ScootersValidator interface {
	ValidateCreateScooterRequest(request *types.CreateScooterRequest) error
	ValidateGetScootersQueryParameters(queryParams *types.GetScootersQueryParameters) error
}

type TripsValidator interface {
	ValidateStartTripRequest(request *types.StartTripRequest) error
	ValidateTripUpdateRequest(request *types.TripUpdateRequest) error
}
