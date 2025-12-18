package handler

import (
	"errors"

	"github.com/AliasgharHeidari/gift-credit/internal/model"
	"github.com/AliasgharHeidari/gift-credit/internal/service"
	"github.com/gofiber/fiber/v2"
)

func UseGiftCode(c *fiber.Ctx) error {

	var req model.Input
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request lmklmklbody",
		})
	}

	err = service.UseGiftCode(req)
	if errors.Is(err, service.InternalErr) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "invalid gift code",
		})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message" : "success",
	})
}
