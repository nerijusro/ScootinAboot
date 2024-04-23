package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/nerijusro/scootinAboot/cmd/api"
	"github.com/nerijusro/scootinAboot/config"
	"github.com/nerijusro/scootinAboot/db"
	"github.com/nerijusro/scootinAboot/services/auth"
	"github.com/nerijusro/scootinAboot/services/client"
	"github.com/nerijusro/scootinAboot/services/scooter"
	"github.com/nerijusro/scootinAboot/services/trip"
	"github.com/nerijusro/scootinAboot/types/interfaces"
	"github.com/nerijusro/scootinAboot/utils"
)

func main() {
	db := createAndInitializeMySqlStorage()
	serviceLocator := createAndInitializeDependencies(db)

	serverAddress := utils.NewServerAddress(config.Envs.PublicHost, config.Envs.Port)
	server := api.NewAPIServer(serverAddress, db, serviceLocator)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func createAndInitializeDependencies(db *sql.DB) *utils.ServiceLocator {
	serviceLocator := &utils.ServiceLocator{
		EndpointHandlers: make(map[string]interfaces.EndpointHandler),
	}

	authService := auth.NewAuthService(config.Envs.StaticAdminApiKey, config.Envs.StaticUserApiKey)

	authHandler := auth.NewAuthorizationHandler(authService)

	clientsRepository := client.NewClientsRepository(db)
	clientHandler := client.NewClientsHandler(clientsRepository, authService)

	scootersRepository := scooter.NewScootersRepository(db)
	scootersRequestValidator := scooter.NewScootersValidator()
	scootersHandler := scooter.NewScootersHandler(scootersRepository, authService, scootersRequestValidator)

	tripsRepository := trip.NewTripsRepository(db)
	tripsValidator := trip.NewTripsValidator()
	tripHandler := trip.NewTripHandler(authService, tripsValidator, tripsRepository, scootersRepository, clientsRepository)

	serviceLocator.RegisterService("authHandler", authHandler)
	serviceLocator.RegisterService("clientHandler", clientHandler)
	serviceLocator.RegisterService("scootersHandler", scootersHandler)
	serviceLocator.RegisterService("tripHandler", tripHandler)

	return serviceLocator
}

func createAndInitializeMySqlStorage() *sql.DB {
	db, err := db.NewMySqlStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  config.Envs.Net,
		AllowNativePasswords: config.Envs.AllowNativePasswords,
		ParseTime:            config.Envs.ParseTime,
	})

	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)
	log.Println("Successfully connected to database")
	return db
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}
