package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/nerijusro/scootinAboot/cmd/api"
	"github.com/nerijusro/scootinAboot/config"
	"github.com/nerijusro/scootinAboot/db"
	"github.com/nerijusro/scootinAboot/types"
)

func main() {
	db := createAndInitializeMySqlStorage()

	serverAddress := types.NewServerAddress(config.Envs.PublicHost, config.Envs.Port)
	server := api.NewAPIServer(serverAddress, db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
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
