package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/nerijusro/scootinAboot/services/auth"
	"github.com/nerijusro/scootinAboot/services/scooter"
	"github.com/nerijusro/scootinAboot/types"
)

type APIServer struct {
	address *types.ServerAddress
	db      *sql.DB
}

func NewAPIServer(address *types.ServerAddress, db *sql.DB) *APIServer {
	return &APIServer{address: address, db: db}
}

func (s *APIServer) Run() error {
	ginEngine := gin.Default()

	staticUserApiKey := "my_static_user_api_key"
	staticAdminApiKey := "my_static_admin_api_key"

	authService := auth.NewAuthService(staticAdminApiKey, staticUserApiKey)

	authHandler := auth.NewAuthorizationHandler(authService)
	authHandler.RegisterEndpoints(ginEngine)

	//Padaryt su DI
	scootersRepository := scooter.NewScootersRepository(s.db)
	scootersHandler := scooter.NewScootersHandler(scootersRepository, authService)
	scootersHandler.RegisterEndpoints(ginEngine)

	ginEngine.Run(s.address.String())
	return nil
}
