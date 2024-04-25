package interfaces

import (
	"github.com/gin-gonic/gin"
	"github.com/nerijusro/scootinAboot/types"
)

type EndpointHandler interface {
	RegisterEndpoints(routerGroups map[string]*gin.RouterGroup)
}

type AuthService interface {
	AuthenticateAdmin(c *gin.Context)
	AuthenticateClient(c *gin.Context)
}

type AuthProvider interface {
	GetAdminApiKey() string
	GetUserApiKey() string
}

type ScooterRepository interface {
	GetScooterById(id string) (*types.Scooter, *int, error)
	GetAllScooters() ([]*types.Scooter, error)
	GetScootersByArea(queryParams types.GetScootersQueryParameters) ([]*types.Scooter, error)
	CreateScooter(scooter types.Scooter) error
}

type ClientRepository interface {
	CreateUser(client types.MobileClient) error
	GetUserById(id string) (*types.MobileClient, *int, error)
}

type TripRepository interface {
	GetTripById(id string) (*types.Trip, error)
	StartTrip(trip types.Trip, scooterOptLockVersion *int, userOptLockVersion *int, event types.TripEvent) error
	UpdateTrip(trip *types.Trip, scooterOptLockVersion *int, event types.TripEvent) error
	EndTrip(trip *types.Trip, scooterOptLockVersion *int, userOptLockVersion *int, event types.TripEvent) error
}

type ScooterValidator interface {
	ValidateCreateScooterRequest(request *types.CreateScooterRequest) error
	ValidateGetScootersQueryParameters(queryParams *types.GetScootersQueryParameters) error
}

type TripValidator interface {
	ValidateStartTripRequest(request *types.StartTripRequest) error
	ValidateTripUpdateRequest(request *types.TripUpdateRequest) error
}
