package handler

import (
	"errors"
	"log"


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

func GiftCodeStatus(c *fiber.Ctx) error {
	GiftCode := c.Params("giftcode")

	result, err := service.GiftCodeStatus(GiftCode)
	if errors.Is(err, service.InternalErr) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unexpected error, please try again later",
		})
	}
	if errors.Is(err, service.ErrNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "GiftCode does not exist, make sure you've entered it correctly",
		})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"Code":        result.Code,
		"Created-At":  result.UpdatedAt.Format("2006-01-02 15:04:05"),
		"Used-counts": result.UsedCount,
		"max-usages":  result.MaxUsage,
		"IsActive":    result.IsActive,
	})

}

func CreateGiftCode(c *fiber.Ctx) error {
	var NewGiftCode model.GiftCode
	err := c.BodyParser(&NewGiftCode)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	err = service.CreateGiftCode(NewGiftCode)
	if errors.Is(err, service.ErrGiftCodeAleadyExist) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "GiftCode already exist",
		})
	}

	if errors.Is(err, service.InternalErr) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal error, please try again later",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"GiftCode Created": nil,
		"Code":             NewGiftCode.Code,
		"Max-usage":        NewGiftCode.MaxUsage,
	})

}
