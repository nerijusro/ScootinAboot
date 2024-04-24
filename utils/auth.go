package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthorizationService struct {
	adminApiKey string
	userApiKey  string
}

func NewAuthService(adminApiKey string, userApiKey string) *AuthorizationService {
	return &AuthorizationService{adminApiKey: adminApiKey, userApiKey: userApiKey}
}

func (s *AuthorizationService) AuthenticateClient(c *gin.Context) {
	apiKey := c.GetHeader("X-API-Key")
	if apiKey != s.adminApiKey && apiKey != s.userApiKey {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.Next()
}

func (s *AuthorizationService) AuthenticateAdmin(c *gin.Context) {
	apiKey := c.GetHeader("X-API-Key")
	if apiKey != s.adminApiKey {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.Next()
}

func (s *AuthorizationService) GetAdminApiKey() string {
	return s.adminApiKey
}

func (s *AuthorizationService) GetUserApiKey() string {
	return s.userApiKey
}
