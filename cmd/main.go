package main

import (
	"github.com/AliasgharHeidari/gift-credit/internal/api/server"
	postgres "github.com/AliasgharHeidari/gift-credit/internal/repository/postgres"
)

func main() {
	postgres.InitDB()
	postgres.AutoMigrate()
	server.Start()

}
