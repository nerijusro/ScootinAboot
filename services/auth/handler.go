package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nerijusro/scootinAboot/types"
)

type AuthorizationHandler struct {
	authService types.IAuthService
}

func NewAuthorizationHandler(authService types.IAuthService) *AuthorizationHandler {
	return &AuthorizationHandler{authService: authService}
}

func (h *AuthorizationHandler) RegisterEndpoints(e *gin.Engine) {
	e.GET("/authUser", h.authorizeUser)
	e.GET("/authAdmin", h.authorizeAdmin)
}

func (h *AuthorizationHandler) authorizeUser(c *gin.Context) {
	apiKey := types.AuthResponse{StaticApiKey: h.authService.GetUserApiKey()}
	c.JSON(http.StatusOK, apiKey)
}

func (h *AuthorizationHandler) authorizeAdmin(c *gin.Context) {
	apiKey := types.AuthResponse{StaticApiKey: h.authService.GetAdminApiKey()}
	c.JSON(http.StatusOK, apiKey)
}