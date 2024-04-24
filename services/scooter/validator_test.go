package scooter

import (
	"testing"

	"github.com/nerijusro/scootinAboot/types"
)

func TestScootersValidator(t *testing.T) {
	validator := NewScootersValidator()

	t.Run("When validating create scooter request while given valid request body returns nil", func(t *testing.T) {
		requestBody := types.CreateScooterRequest{
			Location:    types.Location{Latitude: 54.12, Longitude: 25.34},
			IsAvailable: true,
		}

		result := validator.ValidateCreateScooterRequest(&requestBody)
		if result != nil {
			t.Errorf("expected result to be nil, got %s", result.Error())
		}
	})

	t.Run("When validating create scooter request while given invalid latitude returns error", func(t *testing.T) {
		requestBody := types.CreateScooterRequest{
			Location:    types.Location{Latitude: 1234.12, Longitude: 25.34},
			IsAvailable: true,
		}

		result := validator.ValidateCreateScooterRequest(&requestBody)
		if result.Error() != "invalid latitude" {
			t.Errorf("expected result to be invalid latitude, got %s", result.Error())
		}
	})

	t.Run("When validating create scooter request while given invalid longitude returns error", func(t *testing.T) {
		requestBody := types.CreateScooterRequest{
			Location:    types.Location{Latitude: 13.12, Longitude: 3325.34},
			IsAvailable: true,
		}

		result := validator.ValidateCreateScooterRequest(&requestBody)
		if result.Error() != "invalid longitude" {
			t.Errorf("expected result to be invalid longitude, got %s", result.Error())
		}
	})

	t.Run("When validating get scooters query while given valid params returns nil", func(t *testing.T) {
		requestBody := types.GetScootersQueryParameters{
			Availability: "available",
			X1:           25.0,
			X2:           26.0,
			Y1:           54.0,
			Y2:           55.0,
		}

		result := validator.ValidateGetScootersQueryParameters(&requestBody)
		if result != nil {
			t.Errorf("expected result to be nil, got %s", result.Error())
		}
	})

	t.Run("When validating get scooters query while given invalid availablity returns error", func(t *testing.T) {
		requestBody := types.GetScootersQueryParameters{
			Availability: "nothing",
			X1:           25.0,
			X2:           26.0,
			Y1:           54.0,
			Y2:           55.0,
		}

		result := validator.ValidateGetScootersQueryParameters(&requestBody)
		if result.Error() != "invalid availability" {
			t.Errorf("expected result to be invalid availability, got %s", result.Error())
		}
	})

	t.Run("When validating get scooters query while given invalid x1 returns error", func(t *testing.T) {
		requestBody := types.GetScootersQueryParameters{
			Availability: "available",
			X1:           525.0,
			X2:           26.0,
			Y1:           54.0,
			Y2:           55.0,
		}

		result := validator.ValidateGetScootersQueryParameters(&requestBody)
		if result.Error() != "invalid X1" {
			t.Errorf("expected result to be invalid x1, got %s", result.Error())
		}
	})

	t.Run("When validating get scooters query while given invalid x2 returns error", func(t *testing.T) {
		requestBody := types.GetScootersQueryParameters{
			Availability: "available",
			X1:           25.0,
			X2:           526.0,
			Y1:           54.0,
			Y2:           55.0,
		}

		result := validator.ValidateGetScootersQueryParameters(&requestBody)
		if result.Error() != "invalid X2" {
			t.Errorf("expected result to be invalid x2, got %s", result.Error())
		}
	})

	t.Run("When validating get scooters query while given invalid y1 returns error", func(t *testing.T) {
		requestBody := types.GetScootersQueryParameters{
			Availability: "available",
			X1:           25.0,
			X2:           26.0,
			Y1:           554.0,
			Y2:           55.0,
		}

		result := validator.ValidateGetScootersQueryParameters(&requestBody)
		if result.Error() != "invalid Y1" {
			t.Errorf("expected result to be invalid y1, got %s", result.Error())
		}
	})

	t.Run("When validating get scooters query while given invalid y2 returns error", func(t *testing.T) {
		requestBody := types.GetScootersQueryParameters{
			Availability: "available",
			X1:           25.0,
			X2:           26.0,
			Y1:           54.0,
			Y2:           555.0,
		}

		result := validator.ValidateGetScootersQueryParameters(&requestBody)
		if result.Error() != "invalid Y2" {
			t.Errorf("expected result to be invalid y2, got %s", result.Error())
		}
	})
}
