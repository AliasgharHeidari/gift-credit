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
	app.Post("/gift/use", handler.UseGiftCode)

	// Get giftcode status
/* 	app.("/gift/status", handler.GiftCodeStatus) */

	
	app.Listen(cfg.Port)
}
