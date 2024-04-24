package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/nerijusro/scootinAboot/config"
	"github.com/nerijusro/scootinAboot/services/auth"
	"github.com/nerijusro/scootinAboot/services/client"
	"github.com/nerijusro/scootinAboot/services/scooter"
	"github.com/nerijusro/scootinAboot/services/trip"
	"github.com/nerijusro/scootinAboot/types/interfaces"
	"github.com/nerijusro/scootinAboot/utils"
)

type APIServer struct {
	address *utils.ServerAddress
	db      *sql.DB
}

func NewAPIServer(address *utils.ServerAddress, db *sql.DB) *APIServer {
	return &APIServer{address: address, db: db}
}

// TODO
// trip handler testai
// Child procesas
// Dockerfile
// Dokumentacija
// Gal iseitu autentifikacija i middleware
// Data dump
func (s *APIServer) Run() error {
	ginEngine := gin.Default()

	serviceLocator := buildServiceLocator(s.db)
	authService := serviceLocator.AuthMiddlewares["authMiddleware"]
	routerGroups := enableAuthenticationAndGetRouterGroups(ginEngine, authService)

	for _, handler := range serviceLocator.EndpointHandlers {
		handler.RegisterEndpoints(routerGroups)
	}

	ginEngine.Run(s.address.String())
	return nil
}

func buildServiceLocator(db *sql.DB) *utils.ServiceLocator {
	authService := utils.NewAuthService(config.Envs.StaticAdminApiKey, config.Envs.StaticUserApiKey)

	authHandler := auth.NewAuthorizationHandler(authService)

	clientsRepository := client.NewRepository(db)
	clientHandler := client.NewClientsHandler(clientsRepository)

	scootersRepository := scooter.NewRepository(db)
	scootersRequestValidator := scooter.NewScootersValidator()
	scootersHandler := scooter.NewScootersHandler(scootersRepository, scootersRequestValidator)

	tripsRepository := trip.NewRepository(db)
	tripsValidator := trip.NewTripsValidator()
	tripHandler := trip.NewTripHandler(tripsValidator, tripsRepository, scootersRepository, clientsRepository)

	serviceLocator := &utils.ServiceLocator{
		EndpointHandlers: make(map[string]interfaces.EndpointHandler),
		AuthMiddlewares:  make(map[string]interfaces.AuthService),
	}

	serviceLocator.RegisterAuthMiddleware("authMiddleware", authService)

	serviceLocator.RegisterEndpointHandler("authHandler", authHandler)
	serviceLocator.RegisterEndpointHandler("clientHandler", clientHandler)
	serviceLocator.RegisterEndpointHandler("scootersHandler", scootersHandler)
	serviceLocator.RegisterEndpointHandler("tripHandler", tripHandler)

	return serviceLocator
}

func enableAuthenticationAndGetRouterGroups(e *gin.Engine, authService interfaces.AuthService) map[string]*gin.RouterGroup {
	adminAuthorized := e.Group("/admin", authService.AuthenticateAdmin)
	userAuthorized := e.Group("/client", authService.AuthenticateClient)

	routerGroups := map[string]*gin.RouterGroup{
		"root":   &e.RouterGroup,
		"admin":  adminAuthorized,
		"client": userAuthorized,
	}

	return routerGroups
}
