package types

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Scooter struct {
	ID          uuid.UUID `json:"id"`
	Location    Location  `json:"location"`
	IsAvailable bool      `json:"is_available"`
}

type Location struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type CreateScooterRequest struct {
	Location    Location `json:"location" validate:"required"`
	IsAvailable bool     `json:"is_available"`
}

type IScootersRepository interface {
	GetScooterById(id string) (*Scooter, *int, error)
	GetAllScooters() ([]*Scooter, error)
	GetScootersByArea(queryParams GetScootersQueryParameters) ([]*Scooter, error)
	CreateScooter(scooter Scooter) error
}

type AuthResponse struct {
	StaticApiKey string
}

type IAuthService interface {
	GetAdminApiKey() string
	GetUserApiKey() string

	AuthenticateAdmin(c *gin.Context) bool
	AuthenticateUser(c *gin.Context) bool
}

type MobileClient struct {
	ID                 uuid.UUID `json:"id"`
	FullName           string    `json:"full_name"`
	IsEligibleToTravel bool      `json:"is_eligible_to_travel"`
}

type IClientsRepository interface {
	CreateUser(client MobileClient) error
	GetUserById(id string) (*MobileClient, *int, error)
}

type CreateUserRequest struct {
	FullName string `json:"full_name"`
}

type GetScootersQueryParameters struct {
	Availability string  `form:"availability" validate:"required"`
	X1           float64 `form:"x1" validate:"required"`
	X2           float64 `form:"x2" validate:"required"`
	Y1           float64 `form:"y1" validate:"required"`
	Y2           float64 `form:"y2" validate:"required"`
}

type Availability string

const (
	Available   Availability = "available"
	Unavailable Availability = "unavailable"
	All         Availability = "all"
)

type IScootersValidator interface {
	ValidateCreateScooterRequest(request *CreateScooterRequest) error
	ValidateGetScootersQueryParameters(queryParams *GetScootersQueryParameters) error
}

type StartTripRequest struct {
	ScooterID uuid.UUID `json:"scooter_id" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}

type TripUpdateRequest struct {
	TripID    uuid.UUID `json:"trip_id" validate:"required"`
	Location  Location  `json:"location" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	Sequence  int       `json:"sequence" validate:"required"`
}

type ITripsValidator interface {
	ValidateStartTripRequest(request *StartTripRequest) error
	ValidateTripUpdateRequest(request *TripUpdateRequest) error
}

type Trip struct {
	ID         uuid.UUID `json:"id"`
	ScooterId  uuid.UUID `json:"scooter"`
	ClientId   uuid.UUID `json:"client_id"`
	IsFinished bool      `json:"is_finished"`
}

type TripEvent struct {
	TripID    uuid.UUID `json:"trip_id"`
	Type      string    `json:"event_type"`
	Location  Location  `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	Sequence  int       `json:"sequence"`
}

type ITripsRepository interface {
	GetTripById(id string) (*Trip, error)
	StartTrip(trip Trip, scooterOptLockVersion *int, userOptLockVersion *int, event TripEvent) error
	UpdateTrip(trip *Trip, scooterOptLockVersion *int, event TripEvent) error
	EndTrip(trip *Trip, scooterOptLockVersion *int, userOptLockVersion *int, event TripEvent) error
}
