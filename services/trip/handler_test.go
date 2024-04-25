package trip

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nerijusro/scootinAboot/types"
	"github.com/nerijusro/scootinAboot/types/enums"
)

func TestTripHandler(t *testing.T) {
	validator := &mockTripValidator{}
	tripRepository := &mockTripRepository{}
	scooterRepository := &mockScooterRepository{}
	userRepository := &mockClientRepository{}

	handler := NewTripHandler(validator, tripRepository, scooterRepository, userRepository)

	t.Run("When starting trip while everything is valid returns ok", func(t *testing.T) {
		requestBody := types.StartTripRequest{
			CreatedAt: time.Now(),
			ScooterID: uuid.New(),
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPost, "/client/trips", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		clientId := uuid.New().String()
		request.Header.Set("client-id", clientId)

		router := gin.Default()
		router.POST("/client/trips", handler.startTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusCreated {
			t.Errorf("expected status code %d but got %d", http.StatusCreated, responseRecoreder.Code)
		}

		var response types.TripEvent
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response.TripID == uuid.Nil {
			t.Errorf("expected newly generated trip id")
		}

		if response.Type != enums.StartTrip {
			t.Errorf("expected trip type to be %s, got %s", enums.StartTrip, response.Type)
		}
	})

	t.Run("When starting trip while client id is not set returns unauthorized", func(t *testing.T) {
		requestBody := types.StartTripRequest{
			CreatedAt: time.Now(),
			ScooterID: uuid.New(),
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPost, "/client/trips", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		router := gin.Default()
		router.POST("/client/trips", handler.startTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %d but got %d", http.StatusUnauthorized, responseRecoreder.Code)
		}
	})

	t.Run("When starting trip while invalid request body returns bad request", func(t *testing.T) {
		requestBody := "I ain't json"

		request, err := http.NewRequest(http.MethodPost, "/client/trips", bytes.NewBuffer([]byte(requestBody)))
		if err != nil {
			t.Fatal(err)
		}

		clientId := uuid.New().String()
		request.Header.Set("client-id", clientId)

		router := gin.Default()
		router.POST("/client/trips", handler.startTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, responseRecoreder.Code)
		}
	})

	t.Run("When starting trip while request body content is not valid returns bad request", func(t *testing.T) {
		requestBody := types.StartTripRequest{
			CreatedAt: time.Now(),
			ScooterID: uuid.Nil,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPost, "/client/trips", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		clientId := uuid.New().String()
		request.Header.Set("client-id", clientId)

		router := gin.Default()
		router.POST("/client/trips", handler.startTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, responseRecoreder.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response["Bad request"] != "invalid scooter id" {
			t.Errorf("expected message to be: invalid scooter id, got %f", response["error"])
		}
	})

	t.Run("When starting trip while scooter id is not found returns internal server error", func(t *testing.T) {
		id := "03de5edd-e9d7-4c4e-1111-0ff9c07b6a37"
		idUuid, err := uuid.Parse(id)
		if err != nil {
			t.Fatal(err)
		}

		requestBody := types.StartTripRequest{
			CreatedAt: time.Now(),
			ScooterID: idUuid,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPost, "/client/trips", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		clientId := uuid.New().String()
		request.Header.Set("client-id", clientId)

		router := gin.Default()
		router.POST("/client/trips", handler.startTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, responseRecoreder.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response["Internal Server Error"] != "scooter with id 03de5edd-e9d7-4c4e-1111-0ff9c07b6a37 not found" {
			t.Errorf("expected message to be: scooter with id 03de5edd-e9d7-4c4e-1111-0ff9c07b6a37 not found, got %f", response["error"])
		}
	})

	t.Run("When starting trip while user id is not found returns internal server error", func(t *testing.T) {
		requestBody := types.StartTripRequest{
			CreatedAt: time.Now(),
			ScooterID: uuid.New(),
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPost, "/client/trips", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		clientId := "bec6a2fb-896f-473e-2222-a4208d033498"
		request.Header.Set("client-id", clientId)

		router := gin.Default()
		router.POST("/client/trips", handler.startTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, responseRecoreder.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response["Internal Server Error"] != "error getting user by id" {
			t.Errorf("expected message to be: error getting user by id, got %f", response["error"])
		}
	})

	t.Run("When starting trip while scooter or user is not eligible to travel returns bad request", func(t *testing.T) {
		id := "03de5edd-e9d7-4c4e-2222-0ff9c07b6a37"
		idUuid, err := uuid.Parse(id)
		if err != nil {
			t.Fatal(err)
		}

		requestBody := types.StartTripRequest{
			CreatedAt: time.Now(),
			ScooterID: idUuid,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPost, "/client/trips", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		clientId := uuid.New().String()
		request.Header.Set("client-id", clientId)

		router := gin.Default()
		router.POST("/client/trips", handler.startTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, responseRecoreder.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response["message"] != "trip cannot be started due to parameter invalidity" {
			t.Errorf("expected message to be: trip cannot be started due to parameter invalidity, got %f", response["message"])
		}
	})

	t.Run("When starting trip while repository returns error returns internal server error", func(t *testing.T) {
		requestBody := types.StartTripRequest{
			CreatedAt: time.Now(),
			ScooterID: uuid.New(),
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPost, "/client/trips", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		clientId := "bec6a2fb-896f-473e-a3d5-a4208d033498"
		request.Header.Set("client-id", clientId)

		router := gin.Default()
		router.POST("/client/trips", handler.startTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, responseRecoreder.Code)
		}
	})

	t.Run("When updating trip while everything is valid returns ok", func(t *testing.T) {
		requestBody := types.TripUpdateRequest{
			Location:  types.Location{Latitude: 54.12, Longitude: 25.34},
			CreatedAt: time.Now(),
			Sequence:  2,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPut, "/client/trips/5266c8a2-7a04-45ab-a26c-2a6c9e73bb30", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		clientId := "bec6a2fb-896f-473e-1111-a4208d033498"
		request.Header.Set("client-id", clientId)

		router := gin.Default()
		router.PUT("/client/trips/:id", handler.updateTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, responseRecoreder.Code)
		}

		var response types.TripEvent
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response.Type != enums.UpdateTrip {
			t.Errorf("expected trip type to be %s, got %s", enums.UpdateTrip, response.Type)
		}

		if response.Sequence != 2 {
			t.Errorf("expected sequence to be 2, got %d", response.Sequence)
		}

		if response.Location.Latitude != 54.12 {
			t.Errorf("expected latitude to be 54.12, got %f", response.Location.Latitude)
		}

		if response.Location.Longitude != 25.34 {
			t.Errorf("expected longitude to be 25.34, got %f", response.Location.Longitude)
		}

		if response.TripID.String() != "5266c8a2-7a04-45ab-a26c-2a6c9e73bb30" {
			t.Errorf("expected trip id to be 5266c8a2-7a04-45ab-a26c-2a6c9e73bb30, got %s", response.TripID.String())
		}
	})

	t.Run("When updating trip while client id is not set returns unauthorized", func(t *testing.T) {
		requestBody := types.TripUpdateRequest{
			Location:  types.Location{Latitude: 54.12, Longitude: 25.34},
			CreatedAt: time.Now(),
			Sequence:  2,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPut, "/client/trips/5266c8a2-7a04-45ab-a26c-2a6c9e73bb30", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		router := gin.Default()
		router.PUT("/client/trips/:id", handler.updateTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %d but got %d", http.StatusUnauthorized, responseRecoreder.Code)
		}
	})

	t.Run("When updating trip while invalid request body returns bad request", func(t *testing.T) {
		requestBody := "I ain't json"

		request, err := http.NewRequest(http.MethodPut, "/client/trips/5266c8a2-7a04-45ab-a26c-2a6c9e73bb30", bytes.NewBuffer([]byte(requestBody)))
		if err != nil {
			t.Fatal(err)
		}

		clientId := "5266c8a2-1234-45ab-a26c-a4208d033498"
		request.Header.Set("client-id", clientId)

		router := gin.Default()
		router.PUT("/client/trips/:id", handler.updateTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, responseRecoreder.Code)
		}
	})

	t.Run("When updating trip while request body content is not valid returns bad request", func(t *testing.T) {
		requestBody := types.TripUpdateRequest{
			Location:  types.Location{Latitude: 1234.12, Longitude: 25.34},
			CreatedAt: time.Now(),
			Sequence:  2,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPut, "/client/trips/5266c8a2-7a04-45ab-a26c-2a6c9e73bb30", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		clientId := "5266c8a2-1234-45ab-a26c-a4208d033498"
		request.Header.Set("client-id", clientId)

		router := gin.Default()
		router.PUT("/client/trips/:id", handler.updateTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, responseRecoreder.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response["Bad request"] != "invalid latitude" {
			t.Errorf("expected message to be: invalid latitude, got %f", response["error"])
		}
	})

	t.Run("When updating trip while the trip is not found returns internal server error", func(t *testing.T) {
		requestBody := types.TripUpdateRequest{
			Location:  types.Location{Latitude: 54.12, Longitude: 25.34},
			CreatedAt: time.Now(),
			Sequence:  2,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPut, "/client/trips/5266c8a2-7a04-45ab-1111-2a6c9e73bb30", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		clientId := "5266c8a2-1234-45ab-a26c-a4208d033498"
		request.Header.Set("client-id", clientId)

		router := gin.Default()
		router.PUT("/client/trips/:id", handler.updateTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, responseRecoreder.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response["Internal Server Error"] != "error getting trip by id" {
			t.Errorf("expected message to be: error getting trip by id, got %f", response["error"])
		}
	})

	t.Run("When updating trip while trip is not clients returns unauthorized", func(t *testing.T) {
		requestBody := types.TripUpdateRequest{
			Location:  types.Location{Latitude: 54.12, Longitude: 25.34},
			CreatedAt: time.Now(),
			Sequence:  2,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPut, "/client/trips/5266c8a2-7a04-45ab-3333-2a6c9e73bb30", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		clientId := "5266c8a2-1234-45ab-3333-a4208d033498"
		request.Header.Set("client-id", clientId)

		router := gin.Default()
		router.PUT("/client/trips/:id", handler.updateTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %d but got %d", http.StatusUnauthorized, responseRecoreder.Code)
		}
	})

	t.Run("When updating trip while it is already finished returns bad request", func(t *testing.T) {
		requestBody := types.TripUpdateRequest{
			Location:  types.Location{Latitude: 54.12, Longitude: 25.34},
			CreatedAt: time.Now(),
			Sequence:  2,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPut, "/client/trips/5266c8a2-7a04-45ab-3333-2a6c9e73bb30", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		clientId := "bec6a2fb-896f-473e-1111-a4208d033498"
		request.Header.Set("client-id", clientId)

		router := gin.Default()
		router.PUT("/client/trips/:id", handler.updateTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, responseRecoreder.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response["Bad request"] != "trip is already finished" {
			t.Errorf("expected message to be: trip is already finished, got %f", response["error"])
		}
	})

	t.Run("When updating trip while scooter could not be found returns internal server error", func(t *testing.T) {
		requestBody := types.TripUpdateRequest{
			Location:  types.Location{Latitude: 54.12, Longitude: 25.34},
			CreatedAt: time.Now(),
			Sequence:  2,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPut, "/client/trips/5266c8a2-7a04-45ab-4444-2a6c9e73bb30", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		clientId := "bec6a2fb-896f-473e-1111-a4208d033498"
		request.Header.Set("client-id", clientId)

		router := gin.Default()
		router.PUT("/client/trips/:id", handler.updateTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, responseRecoreder.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response["message"] != "error getting scooter by id" {
			t.Errorf("expected message to be: error getting scooter by id, got %f", response["error"])
		}
	})

	t.Run("When updating trip while repository fails error returns internal server error", func(t *testing.T) {
		requestBody := types.TripUpdateRequest{
			Location:  types.Location{Latitude: 54.12, Longitude: 25.34},
			CreatedAt: time.Now(),
			Sequence:  2,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPut, "/client/trips/5266c8a2-7a04-45ab-7777-2a6c9e73bb30", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		clientId := "bec6a2fb-896f-473e-1111-a4208d033498"
		request.Header.Set("client-id", clientId)

		router := gin.Default()
		router.PUT("/client/trips/:id", handler.updateTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, responseRecoreder.Code)
		}
	})

	t.Run("When ending trip while repository can not find user by id return internal server error", func(t *testing.T) {
		requestBody := types.TripUpdateRequest{
			Location:    types.Location{Latitude: 54.12, Longitude: 25.34},
			CreatedAt:   time.Now(),
			IsFinishing: true,
			Sequence:    2,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPut, "/client/trips/5266c8a2-7a04-45ab-5555-2a6c9e73bb30", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		clientId := "bec6a2fb-896f-473e-2222-a4208d033498"
		request.Header.Set("client-id", clientId)

		router := gin.Default()
		router.PUT("/client/trips/:id", handler.updateTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, responseRecoreder.Code)
		}
	})

	t.Run("When ending trip while repository fails returns internal server error", func(t *testing.T) {
		requestBody := types.TripUpdateRequest{
			Location:    types.Location{Latitude: 54.12, Longitude: 25.34},
			CreatedAt:   time.Now(),
			IsFinishing: true,
			Sequence:    2,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPut, "/client/trips/5266c8a2-7a04-45ab-1111-2a6c9e73bb30", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		clientId := "bec6a2fb-896f-473e-1111-a4208d033498"
		request.Header.Set("client-id", clientId)

		router := gin.Default()
		router.PUT("/client/trips/:id", handler.updateTrip)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, responseRecoreder.Code)
		}
	})
}

type mockTripRepository struct{}

func (m *mockTripRepository) GetTripById(id string) (*types.Trip, error) {
	idUuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	validClientIdUuid, err := uuid.Parse("bec6a2fb-896f-473e-1111-a4208d033498")
	if err != nil {
		return nil, err
	}

	nonExistantClientIdUuid, err := uuid.Parse("bec6a2fb-896f-473e-2222-a4208d033498")
	if err != nil {
		return nil, err
	}

	validScooterIdUuid, err := uuid.Parse("03de5edd-e9d7-4c4e-2222-0ff9c07b6a37")
	if err != nil {
		return nil, err
	}

	nonExistantScooterIdUuid, err := uuid.Parse("03de5edd-e9d7-4c4e-1111-0ff9c07b6a37")
	if err != nil {
		return nil, err
	}

	if id == "5266c8a2-7a04-45ab-1111-2a6c9e73bb30" {
		return nil, errors.New("error getting trip by id")
	}

	if id == "5266c8a2-7a04-45ab-2222-2a6c9e73bb30" {
		return &types.Trip{ClientId: uuid.New()}, nil
	}

	if id == "5266c8a2-7a04-45ab-3333-2a6c9e73bb30" {
		return &types.Trip{ClientId: validClientIdUuid, IsFinished: true}, nil
	}

	if id == "5266c8a2-7a04-45ab-4444-2a6c9e73bb30" {
		return &types.Trip{ID: idUuid, ClientId: validClientIdUuid, ScooterId: nonExistantScooterIdUuid}, nil
	}

	if id == "5266c8a2-7a04-45ab-5555-2a6c9e73bb30" {
		return &types.Trip{ID: idUuid, ClientId: nonExistantClientIdUuid, ScooterId: validScooterIdUuid}, nil
	}

	return &types.Trip{ID: idUuid, ClientId: validClientIdUuid, ScooterId: validScooterIdUuid}, nil
}

func (m *mockTripRepository) StartTrip(trip types.Trip, scooterOptLockVersion *int, userOptLockVersion *int, event types.TripEvent) error {
	if trip.ClientId.String() == "bec6a2fb-896f-473e-a3d5-a4208d033498" {
		return errors.New("trip can not be started")
	}

	return nil
}

func (m *mockTripRepository) UpdateTrip(trip *types.Trip, scooterOptLockVersion *int, event types.TripEvent) error {
	if trip.ID.String() == "5266c8a2-7a04-45ab-7777-2a6c9e73bb30" {
		return errors.New("error getting trip by id")
	}

	return nil
}

func (m *mockTripRepository) EndTrip(trip *types.Trip, scooterOptLockVersion *int, userOptLockVersion *int, event types.TripEvent) error {
	if trip.ClientId.String() == "5266c8a2-7a04-45ab-7777-2a6c9e73bb30" {
		return errors.New("error getting user by id")
	}

	return nil
}

type mockScooterRepository struct{}

func (m *mockScooterRepository) GetScooterById(id string) (*types.Scooter, *int, error) {
	if id == "03de5edd-e9d7-4c4e-1111-0ff9c07b6a37" {
		return nil, nil, errors.New("scooter with id 03de5edd-e9d7-4c4e-1111-0ff9c07b6a37 not found")
	}

	if id == "03de5edd-e9d7-4c4e-2222-0ff9c07b6a37" {
		return &types.Scooter{IsAvailable: false}, new(int), nil
	}

	return &types.Scooter{IsAvailable: true}, new(int), nil
}

// GetScootersByArea implements interfaces.ScooterRepository.
func (m *mockScooterRepository) GetScootersByArea(queryParams types.GetScootersQueryParameters) ([]*types.Scooter, error) {
	panic("unimplemented")
}

// CreateScooter implements interfaces.ScooterRepository.
func (m *mockScooterRepository) CreateScooter(scooter types.Scooter) error {
	panic("unimplemented")
}

// GetAllScooters implements interfaces.ScooterRepository.
func (m *mockScooterRepository) GetAllScooters() ([]*types.Scooter, error) {
	panic("unimplemented")
}

type mockClientRepository struct{}

// CreateUser implements interfaces.ClientRepository.
func (m *mockClientRepository) CreateUser(client types.MobileClient) error {
	panic("unimplemented")
}

func (m *mockClientRepository) GetUserById(id string) (*types.MobileClient, *int, error) {
	idUuid, err := uuid.Parse(id)
	if err != nil {
		return nil, nil, err
	}

	if id == "bec6a2fb-896f-473e-2222-a4208d033498" {
		return nil, nil, errors.New("error getting user by id")
	}

	return &types.MobileClient{ID: idUuid, IsEligibleToTravel: true}, new(int), nil
}

type mockTripValidator struct{}

func (m *mockTripValidator) ValidateStartTripRequest(request *types.StartTripRequest) error {
	if request.ScooterID == uuid.Nil {
		return errors.New("invalid scooter id")
	}

	return nil
}

func (m *mockTripValidator) ValidateTripUpdateRequest(request *types.TripUpdateRequest) error {
	if request.Location.Latitude > 90 || request.Location.Latitude < -90 {
		return errors.New("invalid latitude")
	}

	return nil
}
