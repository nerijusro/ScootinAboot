package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/nerijusro/scootinAboot/services"
	"github.com/nerijusro/scootinAboot/types"
)

type APIServer struct {
	address *types.ServerAddress
	db      *sql.DB
}

var ginEngine = gin.Default()

func NewAPIServer(address *types.ServerAddress, db *sql.DB) *APIServer {
	return &APIServer{address: address, db: db}
}

func (s *APIServer) StartServer() error {
	//Padaryt su DI
	var handler services.ScootersHandler
	handler.RegisterEndpoints(ginEngine)

	ginEngine.Run(s.address.String())
	return nil
}
