package scooter

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nerijusro/scootinAboot/types"
)

func TestScootersHandler(t *testing.T) {
	authService := &mockAuthService{}
	scootersRepository := &mockScootersRepository{}
	validator := &mockScooterRequestValidator{}
	scootersHandler := NewScootersHandler(scootersRepository, authService, validator)

	t.Run("When creating scooter while given valid request body returns status created", func(t *testing.T) {
		requestBody := types.CreateScooterRequest{
			Location:    types.Location{Latitude: 54.12, Longitude: 25.34},
			IsAvailable: true,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)

		request, err := http.NewRequest(http.MethodPost, "/admin/scooters", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("api-key", "some-api-key")

		router := gin.Default()
		router.POST("/admin/scooters", scootersHandler.createScooter)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusCreated {
			t.Errorf("expected status code %d but got %d", http.StatusCreated, responseRecoreder.Code)
		}

		var response types.Scooter
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response.Location.Latitude != requestBody.Location.Latitude {
			t.Errorf("expected latitude to be %f, got %f", requestBody.Location.Latitude, response.Location.Latitude)
		}

		if response.Location.Longitude != requestBody.Location.Longitude {
			t.Errorf("expected longitude to be %f, got %f", requestBody.Location.Longitude, response.Location.Longitude)
		}

		if response.IsAvailable != requestBody.IsAvailable {
			t.Errorf("expected availability to be %t, got %t", requestBody.IsAvailable, response.IsAvailable)
		}
	})

	t.Run("When creating scooter while required parameters are missing returns bad request", func(t *testing.T) {
		requestBody := map[string]string{
			"is_available": "true",
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)

		request, err := http.NewRequest(http.MethodPost, "/admin/scooters", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("api-key", "some-api-key")

		router := gin.Default()
		router.POST("/admin/scooters", scootersHandler.createScooter)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, responseRecoreder.Code)
		}
	})

	t.Run("When creating scooter while static api key is missing returns unauthorized", func(t *testing.T) {
		requestBody := types.CreateScooterRequest{
			Location:    types.Location{Latitude: 54.12, Longitude: 25.34},
			IsAvailable: true,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)

		request, err := http.NewRequest(http.MethodPost, "/admin/scooters", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		router := gin.Default()
		router.POST("/admin/scooters", scootersHandler.createScooter)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %d but got %d", http.StatusUnauthorized, responseRecoreder.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response["error"] != "Unauthorized" {
			t.Errorf("expected message to be Unauthorized, got %f", response["error"])
		}
	})

	t.Run("When creating scooter while request body is not JSON returns bad request", func(t *testing.T) {
		requestBody := "not-json"

		request, err := http.NewRequest(http.MethodPost, "/admin/scooters", bytes.NewBuffer([]byte(requestBody)))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("api-key", "some-api-key")

		router := gin.Default()
		router.POST("/admin/scooters", scootersHandler.createScooter)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, responseRecoreder.Code)
		}
	})

	t.Run("When creating scooter while request body content is not valid returns bad request", func(t *testing.T) {
		requestBody := types.CreateScooterRequest{
			Location:    types.Location{Latitude: 999.0, Longitude: 25.34},
			IsAvailable: true,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPost, "/admin/scooters", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("api-key", "some-api-key")

		router := gin.Default()
		router.POST("/admin/scooters", scootersHandler.createScooter)

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

	t.Run("When creating scooter while repository fails returns internal server error", func(t *testing.T) {
		requestBody := types.CreateScooterRequest{
			Location:    types.Location{Latitude: 0, Longitude: 0},
			IsAvailable: true,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPost, "/admin/scooters", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("api-key", "some-api-key")

		router := gin.Default()
		router.POST("/admin/scooters", scootersHandler.createScooter)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, responseRecoreder.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response["Internal Server Error"] != "issue while creating scooter" {
			t.Errorf("expected message to be: issue while creating scooter, got %f", response["Internal Server Error"])
		}
	})

	t.Run("When getting scooters by area while everything is valid returns ok", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/scooters?availability=available&x1=25.0&x2=26.0&y1=54.0&y2=55.0", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("api-key", "some-api-key")

		router := gin.Default()
		router.GET("/scooters", scootersHandler.getScootersByArea)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, responseRecoreder.Code)
		}

		var response []*types.Scooter
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		responseElement := response[0]
		if responseElement.Location.Latitude != 12.12 {
			t.Errorf("expected latitude to be: 12.12, got %f", responseElement.Location.Latitude)
		}

		if responseElement.Location.Longitude != 44.34 {
			t.Errorf("expected longitude to be 44.34, got %f", responseElement.Location.Longitude)
		}

		if responseElement.IsAvailable != true {
			t.Errorf("expected availability to be true, got %t", responseElement.IsAvailable)
		}
	})

	t.Run("When getting scooters by area while query parameters are not valid returns bad request", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/scooters?invalidParameter=some-invalid-param", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("api-key", "some-api-key")

		router := gin.Default()
		router.GET("/scooters", scootersHandler.getScootersByArea)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadGateway, responseRecoreder.Code)
		}
	})

	t.Run("When getting scooters by area while query parameter content is not valid returns bad request", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/scooters?availability=noSuchThingAsAvailable&x1=25.0&x2=26.0&y1=54.0&y2=55.0", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("api-key", "some-api-key")

		router := gin.Default()
		router.GET("/scooters", scootersHandler.getScootersByArea)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, responseRecoreder.Code)
		}
	})

	t.Run("When getting scooters by area while repository fails returns not found", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/scooters?availability=available&x1=0&x2=0&y1=0&y2=0", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("api-key", "some-api-key")

		router := gin.Default()
		router.GET("/scooters", scootersHandler.getScootersByArea)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusNotFound {
			t.Errorf("expected status code %d but got %d", http.StatusNotFound, responseRecoreder.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response["Not Found"] != "issue while getting scooters by area" {
			t.Errorf("expected message to be: issue while getting scooters by area, got %f", response["Not Found"])
		}
	})

	t.Run("When getting all scooters while everything is valid returns ok", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/admin/scooters", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("api-key", "some-api-key")

		router := gin.Default()
		router.GET("/admin/scooters", scootersHandler.getAllScooters)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, responseRecoreder.Code)
		}

		var response []*types.Scooter
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		responseElement := response[0]
		if responseElement.Location.Latitude != 12.12 {
			t.Errorf("expected latitude to be: 12.12, got %f", responseElement.Location.Latitude)
		}

		if responseElement.Location.Longitude != 44.34 {
			t.Errorf("expected longitude to be 44.34, got %f", responseElement.Location.Longitude)
		}

		if responseElement.IsAvailable != true {
			t.Errorf("expected availability to be true, got %t", responseElement.IsAvailable)
		}
	})

	t.Run("When getting all scooters while static api key is missing returns unauthorized", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/admin/scooters", nil)
		if err != nil {
			t.Fatal(err)
		}

		router := gin.Default()
		router.GET("/admin/scooters", scootersHandler.getAllScooters)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %d but got %d", http.StatusUnauthorized, responseRecoreder.Code)
		}
	})

	t.Run("When getting scooter by id while everything is valid returns ok", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/scooters/e3344268-d649-4c19-a20c-a0c64a5a6623", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("api-key", "some-api-key")

		router := gin.Default()
		router.GET("/scooters/:id", scootersHandler.getScooter)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, responseRecoreder.Code)
		}

		var response types.Scooter
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response.ID.String() != "e3344268-d649-4c19-a20c-a0c64a5a6623" {
			t.Errorf("expected id to be: e3344268-d649-4c19-a20c-a0c64a5a6623, got %s", response.ID.String())
		}

		if response.Location.Latitude != 12.12 {
			t.Errorf("expected latitude to be: 12.12, got %f", response.Location.Latitude)
		}

		if response.Location.Longitude != 44.34 {
			t.Errorf("expected longitude to be 44.34, got %f", response.Location.Longitude)
		}

		if response.IsAvailable != true {
			t.Errorf("expected availability to be true, got %t", response.IsAvailable)
		}
	})

	t.Run("When getting scooter by id while static api key is missing returns unauthorized", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/scooters/e3344268-d649-4c19-a20c-a0c64a5a6623", nil)
		if err != nil {
			t.Fatal(err)
		}

		router := gin.Default()
		router.GET("/scooters/:id", scootersHandler.getScooter)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %d but got %d", http.StatusUnauthorized, responseRecoreder.Code)
		}
	})

	t.Run("When getting scooter by id while scooter is not existant returns not found", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/scooters/e3344268-1234-4c19-a20c-a0c64a5a6623", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("api-key", "some-api-key")

		router := gin.Default()
		router.GET("/scooters/:id", scootersHandler.getScooter)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusNotFound {
			t.Errorf("expected status code %d but got %d", http.StatusNotFound, responseRecoreder.Code)
		}
	})
}

type mockScooterRequestValidator struct{}

func (m *mockScooterRequestValidator) ValidateCreateScooterRequest(request *types.CreateScooterRequest) error {
	if request.Location.Latitude < -90 || request.Location.Latitude > 90 {
		return errors.New("invalid latitude")
	}

	return nil
}

func (m *mockScooterRequestValidator) ValidateGetScootersQueryParameters(queryParams *types.GetScootersQueryParameters) error {
	if queryParams.Availability != "available" && queryParams.Availability != "unavailable" && queryParams.Availability != "all" {
		return errors.New("invalid availability")
	}

	return nil
}

type mockScootersRepository struct{}

func (m *mockScootersRepository) CreateScooter(scooter types.Scooter) error {
	if scooter.Location.Longitude == 0 || scooter.Location.Latitude == 0 {
		return errors.New("issue while creating scooter")
	}

	return nil
}

func (m *mockScootersRepository) GetAllScooters() ([]*types.Scooter, error) {
	id, _ := uuid.Parse("e3344268-d649-4c19-a20c-a0c64a5a6623")
	var scooters []*types.Scooter
	scooters = append(scooters, &types.Scooter{
		ID:          id,
		Location:    types.Location{Latitude: 12.12, Longitude: 44.34},
		IsAvailable: true,
	})

	return scooters, nil
}

func (m *mockScootersRepository) GetScooterById(id string) (*types.Scooter, *int, error) {
	if id == "e3344268-d649-4c19-a20c-a0c64a5a6623" {
		id, _ := uuid.Parse("e3344268-d649-4c19-a20c-a0c64a5a6623")

		return &types.Scooter{
			ID:          id,
			Location:    types.Location{Latitude: 12.12, Longitude: 44.34},
			IsAvailable: true,
		}, nil, nil
	}

	return nil, nil, errors.New("issue while getting scooter by id")
}

func (m *mockScootersRepository) GetScootersByArea(queryParams types.GetScootersQueryParameters) ([]*types.Scooter, error) {
	if queryParams.X1 == 0 || queryParams.X2 == 0 || queryParams.Y1 == 0 || queryParams.Y2 == 0 {
		return nil, errors.New("issue while getting scooters by area")
	}

	id, _ := uuid.Parse("e3344268-d649-4c19-a20c-a0c64a5a6623")
	return []*types.Scooter{
		{
			ID:          id,
			Location:    types.Location{Latitude: 12.12, Longitude: 44.34},
			IsAvailable: true,
		},
	}, nil
}

type mockAuthService struct{}

func (m *mockAuthService) AuthenticateAdmin(c *gin.Context) bool {
	return c.GetHeader("api-key") == "some-api-key"
}

func (m *mockAuthService) AuthenticateUser(c *gin.Context) bool {
	return c.GetHeader("api-key") == "some-api-key"
}
