package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/nerijusro/scootinAboot/utils"
)

type APIServer struct {
	address        *utils.ServerAddress
	db             *sql.DB
	serviceLocator *utils.ServiceLocator
}

func NewAPIServer(address *utils.ServerAddress, db *sql.DB, serviceLocator *utils.ServiceLocator) *APIServer {
	return &APIServer{address: address, db: db, serviceLocator: serviceLocator}
}

// TODO
// Testai
// Child procesas
// Dockerfile
// Dokumentacija
func (s *APIServer) Run() error {
	ginEngine := gin.Default()
	for _, handler := range s.serviceLocator.EndpointHandlers {
		handler.RegisterEndpoints(ginEngine)
	}

	ginEngine.Run(s.address.String())
	return nil
}
