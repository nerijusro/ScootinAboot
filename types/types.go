package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/nerijusro/scootinAboot/types/enums"
)

// Entities
type Scooter struct {
	ID          uuid.UUID `json:"id"`
	Location    Location  `json:"location"`
	IsAvailable bool      `json:"is_available"`
}

type Location struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type MobileClient struct {
	ID                 uuid.UUID `json:"id"`
	FullName           string    `json:"full_name"`
	IsEligibleToTravel bool      `json:"is_eligible_to_travel"`
}

type Trip struct {
	ID         uuid.UUID `json:"id"`
	ScooterId  uuid.UUID `json:"scooter"`
	ClientId   uuid.UUID `json:"client_id"`
	IsFinished bool      `json:"is_finished"`
}

type TripEvent struct {
	TripID    uuid.UUID           `json:"trip_id"`
	Type      enums.TripEventType `json:"event_type"`
	Location  Location            `json:"location"`
	CreatedAt time.Time           `json:"created_at"`
	Sequence  int                 `json:"sequence"`
}

// Requests
type CreateScooterRequest struct {
	Location    Location `json:"location" validate:"required"`
	IsAvailable bool     `json:"is_available"`
}

type CreateUserRequest struct {
	FullName string `json:"full_name"`
}

type StartTripRequest struct {
	ScooterID uuid.UUID `json:"scooter_id" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}

type TripUpdateRequest struct {
	Location    Location  `json:"location" validate:"required"`
	CreatedAt   time.Time `json:"created_at" validate:"required"`
	IsFinishing bool      `json:"is_finishing"`
	Sequence    int       `json:"sequence" validate:"required"`
}

// Query parameters
type GetScootersQueryParameters struct {
	Availability string  `form:"availability" validate:"required"`
	X1           float64 `form:"x1" validate:"required"`
	X2           float64 `form:"x2" validate:"required"`
	Y1           float64 `form:"y1" validate:"required"`
	Y2           float64 `form:"y2" validate:"required"`
}

// Responses
type AuthResponse struct {
	StaticApiKey string
}

type GetScootersResponse struct {
	Scooters []*Scooter `json:"scooters"`
}
