package types

import "github.com/google/uuid"

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
	GetScooters() ([]*Scooter, error)
	CreateScooter(scooter Scooter) error
}