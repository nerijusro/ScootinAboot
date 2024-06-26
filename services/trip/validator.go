package trip

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/nerijusro/scootinAboot/types"
)

type TripValidator struct{}

var Validator = validator.New()

func NewTripValidator() *TripValidator {
	return &TripValidator{}
}

func (s *TripValidator) ValidateStartTripRequest(request *types.StartTripRequest) error {
	if err := Validator.Struct(request); err != nil {
		return err
	}

	if request.CreatedAt.IsZero() || request.CreatedAt.After(time.Now()) {
		return errors.New("invalid created_at")
	}

	return nil
}

func (s *TripValidator) ValidateTripUpdateRequest(request *types.TripUpdateRequest) error {
	if err := Validator.Struct(request); err != nil {
		return err
	}

	if request.Location.Latitude < -90 || request.Location.Latitude > 90 {
		return errors.New("invalid latitude")
	}

	if request.Location.Longitude < -180 || request.Location.Longitude > 180 {
		return errors.New("invalid longitude")
	}

	if request.CreatedAt.IsZero() || request.CreatedAt.After(time.Now()) {
		return errors.New("invalid created_at")
	}

	if request.Sequence < 1 {
		return errors.New("invalid sequence")
	}

	return nil
}
