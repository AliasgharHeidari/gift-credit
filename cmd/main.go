package main

import (
	"log"

	"github.com/AliasgharHeidari/gift-credit/config"
	"github.com/AliasgharHeidari/gift-credit/internal/api/server"
	postgres "github.com/AliasgharHeidari/gift-credit/internal/repository/postgres"
)

func main() {

	cfg, err := config.Load("./config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("running on Port :", cfg.Server.Port)
	log.Println("server Host :", cfg.Server.Host)

	postgres.InitDB(cfg.Database)
	postgres.AutoMigrate()
	server.Start(cfg.Server)
}
