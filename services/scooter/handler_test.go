package scooter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nerijusro/scootinAboot/types"
)

func TestScootersHandler(t *testing.T) {
	scootersRepository := &mockScootersRepository{}
	authService := &mockAuthService{}
	validator := &mockScooterRequestValidator{} // Add this line

	scootersHandler := NewScootersHandler(scootersRepository, authService, validator) // Update this line

	t.Run("Should pass given valid request body while creating scooter", func(t *testing.T) {
		requestBody := types.CreateScooterRequest{
			Location:    types.Location{Latitude: 54.12, Longitude: 25.34},
			IsAvailable: true,
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)

		request, err := http.NewRequest(http.MethodPost, "/scooters", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		router := gin.Default()
		router.POST("/scooters", scootersHandler.createScooter)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusCreated {
			t.Errorf("expected status code %d but got %d", http.StatusCreated, responseRecoreder.Code)
		}
	})
}

type mockScooterRequestValidator struct{}

// ValidateCreateScooterRequest implements types.IScootersValidator.
func (m *mockScooterRequestValidator) ValidateCreateScooterRequest(request *types.CreateScooterRequest) error {
	panic("unimplemented")
}

// ValidateGetScootersQueryParameters implements types.IScootersValidator.
func (m *mockScooterRequestValidator) ValidateGetScootersQueryParameters(queryParams *types.GetScootersQueryParameters) error {
	panic("unimplemented")
}

// ValidateCreateScooterRequest implements types.IScooterRequestValidator.
func (m *mockScooterRequestValidator) GetAndValidateCreateScooterRequest(request types.CreateScooterRequest) error {
	panic("unimplemented")
}

// ValidateGetScootersQueryParameters implements types.IScooterRequestValidator.
func (m *mockScooterRequestValidator) GetAndValidateGetScootersQueryParameters(queryParams types.GetScootersQueryParameters) error {
	panic("unimplemented")
}

type mockScootersRepository struct{}

// GetAllScooters implements types.IScootersRepository.
func (m *mockScootersRepository) GetAllScooters() ([]*types.Scooter, error) {
	panic("unimplemented")
}

// GetScootersByArea implements types.IScootersRepository.
func (m *mockScootersRepository) GetScootersByArea(queryParams types.GetScootersQueryParameters) ([]*types.Scooter, error) {
	panic("unimplemented")
}

func (m *mockScootersRepository) GetScooters() ([]*types.Scooter, error) {
	return nil, nil
}

func (m *mockScootersRepository) CreateScooter(types.Scooter) error {
	return nil
}

func (m *mockScootersRepository) GetScooterById(string) (*types.Scooter, *int, error) {
	return nil, nil, nil
}

type mockAuthService struct{}

// AuthenticateAdmin implements types.IAuthService.
func (m *mockAuthService) AuthenticateAdmin(c *gin.Context) bool {
	panic("unimplemented")
}

// AuthenticateUser implements types.IAuthService.
func (m *mockAuthService) AuthenticateUser(c *gin.Context) bool {
	panic("unimplemented")
}

// GetAdminApiKey implements types.IAuthService.
func (m *mockAuthService) GetAdminApiKey() string {
	panic("unimplemented")
}

// GetUserApiKey implements types.IAuthService.
func (m *mockAuthService) GetUserApiKey() string {
	panic("unimplemented")
}
