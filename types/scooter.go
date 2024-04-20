package types

import "github.com/google/uuid"

type Scooter struct {
	ID          uuid.UUID `json:"id"`
	Location    Location  `json:"location"`
	IsAvailable bool      `json:"is_available"`
	OccupiedBy  uuid.UUID `json:"occupied_by"`
}
