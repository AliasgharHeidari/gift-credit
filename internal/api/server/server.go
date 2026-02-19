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
	
	// Use GiftCode
	app.Post("/gift", handler.UseGiftCode)

	// Get GiftCode status
 	app.Get("/gift/status/:giftcode", handler.GiftCodeStatus)

	// Create GiftCode
	app.Post("/gift/Create", handler.CreateGiftCode)

	// Get GiftCode list
	app.Get("/gift/", handler.GetGiftCodeList)
	
	// Delete GiftCode
	app.Delete("/gift/", handler.DeleteGiftCode)


	
	app.Listen(cfg.Port)
}
