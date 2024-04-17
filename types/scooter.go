package types

import "github.com/google/uuid"

type ScooterStatus string

const (
	Available = "available"
	Taken     = "taken"
)

type Scooter struct {
	ID         uuid.UUID     `json:"id"`
	Location   Location      `json:"location"`
	Status     ScooterStatus `json:"status"`
	OccupiedBy uuid.UUID     `json:"occupied_by"`
}
