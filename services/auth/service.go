package auth

import "github.com/gin-gonic/gin"

type AuthorizationService struct {
	AdminApiKey string
	UserApiKey  string
}

func NewAuthService(adminApiKey string, userApiKey string) *AuthorizationService {
	return &AuthorizationService{AdminApiKey: adminApiKey, UserApiKey: userApiKey}
}

func (s *AuthorizationService) AuthenticateUser(c *gin.Context) bool {
	apiKey := c.GetHeader("X-API-Key")
	return apiKey == s.AdminApiKey || apiKey == s.UserApiKey
}

func (s *AuthorizationService) AuthenticateAdmin(c *gin.Context) bool {
	apiKey := c.GetHeader("X-API-Key")
	return apiKey == s.AdminApiKey
}

func (s *AuthorizationService) GetAdminApiKey() string {
	return s.AdminApiKey
}

func (s *AuthorizationService) GetUserApiKey() string {
	return s.UserApiKey
}
