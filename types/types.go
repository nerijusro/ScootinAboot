package types

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Scooter struct {
	ID          uuid.UUID `json:"id"`
	Location    Location  `json:"location"`
	IsAvailable bool      `json:"is_available"`
	OccupiedBy  uuid.UUID `json:"occupied_by"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CreateScooterRequest struct {
	Location    Location `json:"location"`
	IsAvailable bool     `json:"is_available"`
}

type IScootersRepository interface {
	GetScooterById(id string) (*Scooter, error)
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
}

type CreateUserRequest struct {
	FullName string `json:"full_name"`
}

type GetScootersQueryParameters struct {
	Availability string  `form:"availability"`
	X1           float64 `form:"x1"`
	X2           float64 `form:"x2"`
	Y1           float64 `form:"y1"`
	Y2           float64 `form:"y2"`
}

type Availability string

const (
	Available   Availability = "available"
	Unavailable Availability = "unavailable"
	All         Availability = "all"
)
