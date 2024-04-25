package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nerijusro/scootinAboot/types"
)

func TestClientHandler(t *testing.T) {
	repository := &mockClientRepository{}
	handler := NewClientsHandler(repository)

	t.Run("When creating user while everything is valid returns ok", func(t *testing.T) {
		requestBody := types.CreateUserRequest{
			FullName: "John Doe",
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPost, "/client/users", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		router := gin.Default()
		router.POST("/client/users", handler.createUser)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusCreated {
			t.Errorf("expected status code %d but got %d", http.StatusCreated, responseRecoreder.Code)
		}

		var response types.MobileClient
		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response.FullName != requestBody.FullName {
			t.Errorf("expected latitude to be %s, got %s", requestBody.FullName, response.FullName)
		}
	})

	t.Run("When creating user while invalid request body returns bad request", func(t *testing.T) {
		requestBody := "I ain't json"

		request, err := http.NewRequest(http.MethodPost, "/client/users", bytes.NewBuffer([]byte(requestBody)))
		if err != nil {
			t.Fatal(err)
		}

		router := gin.Default()
		router.POST("/client/users", handler.createUser)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, responseRecoreder.Code)
		}
	})

	t.Run("When creating user while repository returns error returns internal server error", func(t *testing.T) {
		requestBody := types.CreateUserRequest{
			FullName: "Not John Doe",
		}

		marshalledRequestBody, _ := json.Marshal(requestBody)
		request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		router := gin.Default()
		router.POST("/users", handler.createUser)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, request)

		if responseRecoreder.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, responseRecoreder.Code)
		}
	})
}

type mockClientRepository struct{}

func (m *mockClientRepository) CreateUser(client types.MobileClient) error {
	if client.FullName != "John Doe" {
		return errors.New("internal server error")
	}

	return nil
}

func (m *mockClientRepository) GetUserById(id string) (*types.MobileClient, *int, error) {
	return &types.MobileClient{}, nil, nil
}
