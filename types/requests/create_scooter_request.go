package requests

import "github.com/nerijusro/scootinAboot/types"

type CreateScooterRequest struct {
	Location    types.Location `json:"location"`
	IsAvailable bool           `json:"is_available"`
}
