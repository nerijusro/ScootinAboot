package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthorizationHandler(t *testing.T) {
	authProvider := &mockAuthProvider{}
	handler := NewAuthorizationHandler(authProvider)

	t.Run("When authorizing user returns user api key", func(t *testing.T) {
		router := gin.Default()
		router.GET("/authUser", handler.authorizeUser)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, httptest.NewRequest(http.MethodGet, "/authUser", nil))

		if responseRecoreder.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, responseRecoreder.Code)
		}
	})

	t.Run("When authorizing admin returns admin api key", func(t *testing.T) {
		router := gin.Default()
		router.GET("/authAdmin", handler.authorizeAdmin)

		responseRecoreder := httptest.NewRecorder()
		router.ServeHTTP(responseRecoreder, httptest.NewRequest(http.MethodGet, "/authAdmin", nil))

		if responseRecoreder.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, responseRecoreder.Code)
		}
	})
}

type mockAuthProvider struct{}

func (m *mockAuthProvider) GetAdminApiKey() string {
	return "admin-api-key"
}

func (m *mockAuthProvider) GetUserApiKey() string {
	return "user-api-key"
}
