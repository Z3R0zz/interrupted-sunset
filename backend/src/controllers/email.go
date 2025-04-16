package controllers

import (
	"interrupted-export/src/models"
	"interrupted-export/src/rules"
	"interrupted-export/src/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type EmailInput struct {
	Email string `json:"email" validate:"email,max=200,required"`
}

func NewEmail(c *fiber.Ctx) error {
	var input EmailInput
	if err := c.BodyParser(&input); err != nil {
		return utils.HandleError(c, err, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := rules.Validate.Struct(input); err != nil {
		return utils.HandleError(c, err, fiber.StatusBadRequest, err.Error())
	}

	otp := models.OTP{
		UserID:    c.Locals("user_id").(uint),
		Email:     input.Email,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := otp.Create(c.Context()); err != nil {
		return utils.HandleError(c, err, fiber.StatusInternalServerError, "Failed to create OTP")
	}

	return c.JSON(fiber.Map{
		"message": "OTP sent to email",
	})
}
