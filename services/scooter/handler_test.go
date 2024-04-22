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

	scootersHandler := NewScootersHandler(scootersRepository, authService)

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

type mockScootersRepository struct{}

func (m *mockScootersRepository) GetScooters() ([]*types.Scooter, error) {
	return nil, nil
}

func (m *mockScootersRepository) CreateScooter(types.Scooter) error {
	return nil
}

func (m *mockScootersRepository) GetScooterById(string) (*types.Scooter, error) {
	return nil, nil
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
