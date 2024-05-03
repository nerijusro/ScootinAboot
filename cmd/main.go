package main

import (
	"database/sql"
	"log"
	"sync"

	"github.com/go-sql-driver/mysql"
	"github.com/nerijusro/scootinAboot/cmd/api"

	"github.com/nerijusro/scootinAboot/cmd/child"
	"github.com/nerijusro/scootinAboot/config"
	"github.com/nerijusro/scootinAboot/db"
	"github.com/nerijusro/scootinAboot/utils"
)

func main() {
	db := createAndInitializeMySqlStorage()

	serverAddress := utils.NewServerAddress(config.Envs.Protocol, config.Envs.PublicHost, config.Envs.Port)
	server := api.NewAPIServer(serverAddress, db)
	mobileClientDummy := child.NewMobileClientDummy("http://127.0.0.1:8080")

	var wg sync.WaitGroup
	if err := server.Run(&wg); err != nil {
		log.Println("Error starting Gin server:", err)
	}

	go mobileClientDummy.Run(&wg)
	wg.Wait()
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
		log.Fatal(err.Error())
	}
}
