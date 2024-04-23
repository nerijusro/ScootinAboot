package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/nerijusro/scootinAboot/config"
	"github.com/nerijusro/scootinAboot/services/auth"
	"github.com/nerijusro/scootinAboot/services/client"
	"github.com/nerijusro/scootinAboot/services/scooter"
	"github.com/nerijusro/scootinAboot/services/trip"
	"github.com/nerijusro/scootinAboot/types"
)

type APIServer struct {
	address *types.ServerAddress
	db      *sql.DB
}

func NewAPIServer(address *types.ServerAddress, db *sql.DB) *APIServer {
	return &APIServer{address: address, db: db}
}

// TODO
// Padaryt service locatoriu
// Testai
// Child procesas
// Dockerfile
// Dokumentacija
// Kaip ir types grupavima butu gerai sutvarkyt
// Auth service gal irgi ne vietoj?
func (s *APIServer) Run() error {
	ginEngine := gin.Default()

	authService := auth.NewAuthService(config.Envs.StaticAdminApiKey, config.Envs.StaticUserApiKey)

	authHandler := auth.NewAuthorizationHandler(authService)
	authHandler.RegisterEndpoints(ginEngine)

	clientsRepository := client.NewClientsRepository(s.db)
	clientHandler := client.NewClientsHandler(clientsRepository, authService)
	clientHandler.RegisterEndpoints(ginEngine)

	scootersRepository := scooter.NewScootersRepository(s.db)
	scootersRequestValidator := scooter.NewScootersValidator()
	scootersHandler := scooter.NewScootersHandler(scootersRepository, authService, scootersRequestValidator)
	scootersHandler.RegisterEndpoints(ginEngine)

	tripsRepository := trip.NewTripsRepository(s.db)
	tripsValidator := trip.NewTripsValidator()
	tripHandler := trip.NewTripHandler(authService, tripsValidator, tripsRepository, scootersRepository, clientsRepository)
	tripHandler.RegisterEndpoints(ginEngine)

	ginEngine.Run(s.address.String())
	return nil
}
