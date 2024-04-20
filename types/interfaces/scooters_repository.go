package interfaces

import (
	"github.com/google/uuid"
	"github.com/nerijusro/scootinAboot/types"
)

type ScootersRepository interface {
	GetScooterByID(id uuid.UUID) (*types.Scooter, error)
	CreateScooter(scooter types.Scooter) error
}
