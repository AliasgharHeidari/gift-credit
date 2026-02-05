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
			"error": "invalid request body",
		})
	}

	NewBalance, err := service.UseGiftCode(req)
	if errors.Is(err, service.InternalErr) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unexpected error, make sure the mobile number is correct and retry",
		})
	} else if errors.Is(err, service.ErrGiftCodeUnavailable) {
		return c.Status(fiber.StatusGone).JSON(fiber.Map{
			"error": "GiftCode is currently unavailable",
		})
	} else if errors.Is(err, service.ErrGiftCodeOutOfUse) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "GiftCode is out of use",
		})
	} else if errors.Is(err, service.ErrGiftCodeAleadyUsed) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "GiftCode is already used for this mobile number",
		})
	} else if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "invalid gift code",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message":         "GiftCode applied successfully",
		"current balance": NewBalance,
	})
}

/* func GiftCodeStatus(c *fiber.Ctx) error {

	var GiftCode model.GiftCodeStatus
	err := c.BodyParser(&GiftCode)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid body request",
		})
	}
	status, err := service.GiftCodeStatus(GiftCode)

}
*/
