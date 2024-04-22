package scooter

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/nerijusro/scootinAboot/types"
)

type ScootersValidator struct{}

func NewScootersValidator() *ScootersValidator {
	return &ScootersValidator{}
}

func (s *ScootersValidator) ValidateCreateScooterRequest(request *types.CreateScooterRequest) error {
	validator := validator.New()
	if err := validator.Struct(request); err != nil {
		return err
	}

	if request.Location.Latitude < -90 || request.Location.Latitude > 90 {
		return errors.New("invalid latitude")
	}

	if request.Location.Longitude < -180 || request.Location.Longitude > 180 {
		return errors.New("invalid longitude")
	}

	return nil
}

func (s *ScootersValidator) ValidateGetScootersQueryParameters(queryParams *types.GetScootersQueryParameters) error {
	validator := validator.New()
	if err := validator.Struct(queryParams); err != nil {
		return err
	}

	if !isValidAvailability(queryParams.Availability) {
		return errors.New("invalid availability")
	}

	if queryParams.X1 < -180 || queryParams.X1 > 180 {
		return errors.New("invalid X1")
	}

	if queryParams.X2 < -180 || queryParams.X2 > 180 {
		return errors.New("invalid X2")
	}

	if queryParams.Y1 < -90 || queryParams.Y1 > 90 {
		return errors.New("invalid Y1")
	}

	if queryParams.Y2 < -90 || queryParams.Y2 > 90 {
		return errors.New("invalid Y2")
	}

	return nil
}

func isValidAvailability(availability string) bool {
	validAvailabilities := map[string]bool{
		"available":   true,
		"unavailable": true,
		"all":         true,
	}
	_, ok := validAvailabilities[availability]
	return ok
}