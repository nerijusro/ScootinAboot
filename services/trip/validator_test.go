package trip

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nerijusro/scootinAboot/types"
)

func TestTripsValidator(t *testing.T) {
	validator := NewTripsValidator()

	t.Run("When validating start trip request while given valid request body returns nil", func(t *testing.T) {
		requestBody := types.StartTripRequest{
			CreatedAt: time.Now(),
			ScooterID: uuid.New(),
		}

		result := validator.ValidateStartTripRequest(&requestBody)
		if result != nil {
			t.Errorf("expected result to be nil, got %s", result.Error())
		}
	})

	t.Run("When validating start trip request while given invalid scooter id returns error", func(t *testing.T) {
		requestBody := types.StartTripRequest{
			CreatedAt: time.Now(),
			ScooterID: uuid.Nil,
		}

		result := validator.ValidateStartTripRequest(&requestBody)
		if result.Error() != "Key: 'StartTripRequest.ScooterID' Error:Field validation for 'ScooterID' failed on the 'required' tag" {
			t.Errorf("expected result to be: Key: 'StartTripRequest.ScooterID' Error:Field validation for 'ScooterID' failed on the 'required' tag, got %s", result.Error())
		}
	})

	t.Run("When validating start trip request while given invalid created at returns error", func(t *testing.T) {
		requestBody := types.StartTripRequest{
			CreatedAt: time.Now().Add(time.Hour),
			ScooterID: uuid.New(),
		}

		result := validator.ValidateStartTripRequest(&requestBody)
		if result.Error() != "invalid created_at" {
			t.Errorf("expected result to be invalid created_at, got %s", result.Error())
		}
	})

	t.Run("When validating trip update request while given valid request body returns nil", func(t *testing.T) {
		requestBody := types.TripUpdateRequest{
			Location:  types.Location{Latitude: 54.12, Longitude: 25.34},
			CreatedAt: time.Now(),
			TripID:    uuid.New(),
			Sequence:  2,
		}

		result := validator.ValidateTripUpdateRequest(&requestBody)
		if result != nil {
			t.Errorf("expected result to be nil, got %s", result.Error())
		}
	})

	t.Run("When validating trip update request while given invalid location returns error", func(t *testing.T) {
		requestBody := types.TripUpdateRequest{
			Location:  types.Location{Latitude: 1234.12, Longitude: 25.34},
			CreatedAt: time.Now(),
			TripID:    uuid.New(),
			Sequence:  2,
		}

		result := validator.ValidateTripUpdateRequest(&requestBody)
		if result.Error() != "invalid latitude" {
			t.Errorf("expected result to be invalid latitude, got %s", result.Error())
		}
	})

	t.Run("When validating trip update request while given invalid created at returns error", func(t *testing.T) {
		requestBody := types.TripUpdateRequest{
			Location:  types.Location{Latitude: 54.12, Longitude: 25.34},
			CreatedAt: time.Now().Add(time.Hour),
			TripID:    uuid.New(),
			Sequence:  2,
		}

		result := validator.ValidateTripUpdateRequest(&requestBody)
		if result.Error() != "invalid created_at" {
			t.Errorf("expected result to be invalid created_at, got %s", result.Error())
		}
	})

	t.Run("When validating trip update request while given invalid trip id returns error", func(t *testing.T) {
		requestBody := types.TripUpdateRequest{
			Location:  types.Location{Latitude: 54.12, Longitude: 25.34},
			CreatedAt: time.Now(),
			TripID:    uuid.Nil,
			Sequence:  2,
		}

		result := validator.ValidateTripUpdateRequest(&requestBody)
		if result.Error() != "Key: 'TripUpdateRequest.TripID' Error:Field validation for 'TripID' failed on the 'required' tag" {
			t.Errorf("expected result to be: Key: 'TripUpdateRequest.TripID' Error:Field validation for 'TripID' failed on the 'required' tag, got %s", result.Error())
		}
	})

	t.Run("When validating trip update request while given invalid sequence returns error", func(t *testing.T) {
		requestBody := types.TripUpdateRequest{
			Location:  types.Location{Latitude: 54.12, Longitude: 25.34},
			CreatedAt: time.Now(),
			TripID:    uuid.New(),
			Sequence:  0,
		}

		result := validator.ValidateTripUpdateRequest(&requestBody)
		if result.Error() != "Key: 'TripUpdateRequest.Sequence' Error:Field validation for 'Sequence' failed on the 'required' tag" {
			t.Errorf("expected result to be: Key: 'TripUpdateRequest.Sequence' Error:Field validation for 'Sequence' failed on the 'required' tag, got %s", result.Error())
		}
	})
}
