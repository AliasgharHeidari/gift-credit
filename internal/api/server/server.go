package server

import (
	"github.com/AliasgharHeidari/gift-credit/config"
	"github.com/AliasgharHeidari/gift-credit/internal/api/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

)

func Start(cfg config.ServerConfig) {
	app := fiber.New()

	app.Use(logger.New())
	
	//Use giftcode
	app.Post("/gift", handler.UseGiftCode)

	// Get giftcode status
 	app.Get("/gift/status/:giftcode", handler.GiftCodeStatus)

	// Create GiftCode
	app.Post("/gift/Create", handler.CreateGiftCode)

	
	app.Listen(cfg.Port)
}
