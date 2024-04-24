package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nerijusro/scootinAboot/types"
	"github.com/nerijusro/scootinAboot/types/interfaces"
)

type AuthorizationHandler struct {
	authProvider interfaces.AuthProvider
}

func NewAuthorizationHandler(authProvider interfaces.AuthProvider) *AuthorizationHandler {
	return &AuthorizationHandler{authProvider: authProvider}
}

func (h *AuthorizationHandler) RegisterEndpoints(routerGroups map[string]*gin.RouterGroup) {
	rootGroup := routerGroups["root"]
	rootGroup.GET("/authUser", h.authorizeUser)
	rootGroup.GET("/authAdmin", h.authorizeAdmin)
}

func (h *AuthorizationHandler) authorizeUser(c *gin.Context) {
	apiKey := types.AuthResponse{StaticApiKey: h.authProvider.GetUserApiKey()}
	c.JSON(http.StatusOK, apiKey)
}

func (h *AuthorizationHandler) authorizeAdmin(c *gin.Context) {
	apiKey := types.AuthResponse{StaticApiKey: h.authProvider.GetAdminApiKey()}
	c.JSON(http.StatusOK, apiKey)
}
