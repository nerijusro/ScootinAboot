package client

// import (
// 	"bytes"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/nerijusro/scootinAboot/types"
// )

// func TestClientHandler(t *testing.T) {
// 	authService := &mockAuthService{}
// 	repository := &mockClientsRepository{}
// 	handler := NewClientsHandler(repository, authService)

// 	t.Run("When creating user while everything is valid returns ok", func(t *testing.T) {
// 		requestBody := types.CreateUserRequest{
// 			FullName: "John Doe",
// 		}

// 		marshalledRequestBody, _ := json.Marshal(requestBody)

// 		request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(marshalledRequestBody))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		request.Header.Set("api-key", "some-api-key")

// 		router := gin.Default()
// 		router.POST("/users", handler.createUser)

// 		responseRecoreder := httptest.NewRecorder()
// 		router.ServeHTTP(responseRecoreder, request)

// 		if responseRecoreder.Code != http.StatusCreated {
// 			t.Errorf("expected status code %d but got %d", http.StatusCreated, responseRecoreder.Code)
// 		}

// 		var response types.MobileClient
// 		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
// 			t.Fatal(err)
// 		}

// 		if response.FullName != requestBody.FullName {
// 			t.Errorf("expected latitude to be %s, got %s", requestBody.FullName, response.FullName)
// 		}
// 	})

// 	t.Run("When creating user while static api key is not given returns unauthorized", func(t *testing.T) {
// 		requestBody := types.CreateUserRequest{
// 			FullName: "John Doe",
// 		}

// 		marshalledRequestBody, _ := json.Marshal(requestBody)

// 		request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(marshalledRequestBody))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		router := gin.Default()
// 		router.POST("/users", handler.createUser)

// 		responseRecoreder := httptest.NewRecorder()
// 		router.ServeHTTP(responseRecoreder, request)

// 		if responseRecoreder.Code != http.StatusUnauthorized {
// 			t.Errorf("expected status code %d but got %d", http.StatusUnauthorized, responseRecoreder.Code)
// 		}

// 		var response map[string]interface{}
// 		if err := json.NewDecoder(responseRecoreder.Body).Decode(&response); err != nil {
// 			t.Fatal(err)
// 		}

// 		if response["error"] != "Unauthorized" {
// 			t.Errorf("expected message to be Unauthorized, got %f", response["error"])
// 		}
// 	})

// 	t.Run("When creating user while invalid request body returns bad request", func(t *testing.T) {
// 		requestBody := "I ain't json"

// 		request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte(requestBody)))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		request.Header.Set("api-key", "some-api-key")

// 		router := gin.Default()
// 		router.POST("/users", handler.createUser)

// 		responseRecoreder := httptest.NewRecorder()
// 		router.ServeHTTP(responseRecoreder, request)

// 		if responseRecoreder.Code != http.StatusBadRequest {
// 			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, responseRecoreder.Code)
// 		}
// 	})

// 	t.Run("When creating user while repository returns error returns internal server error", func(t *testing.T) {
// 		requestBody := types.CreateUserRequest{
// 			FullName: "Not John Doe",
// 		}

// 		marshalledRequestBody, _ := json.Marshal(requestBody)

// 		request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(marshalledRequestBody))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		request.Header.Set("api-key", "some-api-key")

// 		router := gin.Default()
// 		router.POST("/users", handler.createUser)

// 		responseRecoreder := httptest.NewRecorder()
// 		router.ServeHTTP(responseRecoreder, request)

// 		if responseRecoreder.Code != http.StatusInternalServerError {
// 			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, responseRecoreder.Code)
// 		}
// 	})
// }

// type mockAuthService struct{}

// func (m *mockAuthService) AuthenticateAdmin(c *gin.Context) bool {
// 	return c.GetHeader("api-key") == "some-api-key"
// }

// func (m *mockAuthService) AuthenticateUser(c *gin.Context) bool {
// 	return c.GetHeader("api-key") == "some-api-key"
// }

// type mockClientsRepository struct{}

// func (m *mockClientsRepository) CreateUser(client types.MobileClient) error {
// 	if client.FullName != "John Doe" {
// 		return errors.New("internal server error")
// 	}

// 	return nil
// }

// func (m *mockClientsRepository) GetUserById(id string) (*types.MobileClient, *int, error) {
// 	return &types.MobileClient{}, nil, nil
// }
