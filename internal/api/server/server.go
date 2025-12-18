package server

import (
	"os"

	"github.com/AliasgharHeidari/gift-credit/internal/api/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func Start() {
	app := fiber.New()

	app.Use(logger.New())
	
	app.Post("/gift/use", handler.UseGiftCode)

	err := godotenv.Load("./.env")
	if err != nil {
		panic(err)
	}

	port := os.Getenv("SERVER_PORT")

	app.Listen(port)

}
