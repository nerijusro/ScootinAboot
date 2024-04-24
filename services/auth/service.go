package auth

import "github.com/gin-gonic/gin"

type AuthorizationService struct {
	adminApiKey string
	userApiKey  string
}

func NewAuthService(adminApiKey string, userApiKey string) *AuthorizationService {
	return &AuthorizationService{adminApiKey: adminApiKey, userApiKey: userApiKey}
}

func (s *AuthorizationService) AuthenticateUser(c *gin.Context) bool {
	apiKey := c.GetHeader("X-API-Key")
	return apiKey == s.adminApiKey || apiKey == s.userApiKey
}

func (s *AuthorizationService) AuthenticateAdmin(c *gin.Context) bool {
	apiKey := c.GetHeader("X-API-Key")
	return apiKey == s.adminApiKey
}

func (s *AuthorizationService) GetAdminApiKey() string {
	return s.adminApiKey
}

func (s *AuthorizationService) GetUserApiKey() string {
	return s.userApiKey
}
