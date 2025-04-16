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

type VerifyInput struct {
	Code string `json:"code" validate:"required,len=6"`
}

func NewMail(c *fiber.Ctx) error {
	var input EmailInput
	if err := c.BodyParser(&input); err != nil {
		return utils.HandleError(c, err, fiber.StatusBadRequest, "invalid request body")
	}

	if err := rules.Validate.Struct(input); err != nil {
		return utils.HandleError(c, err, fiber.StatusBadRequest, err.Error())
	}

	user := c.Locals("user").(*models.User)

	otp := models.OTP{
		UserID:    user.ID,
		Email:     input.Email,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := otp.Create(c.Context()); err != nil {
		return utils.HandleError(c, err, fiber.StatusInternalServerError, "internal server error")
	}

	return c.JSON(fiber.Map{
		"message": "Please check your email for the OTP code",
	})
}

func VerifyMail(c *fiber.Ctx) error {
	var input VerifyInput
	if err := c.BodyParser(&input); err != nil {
		return utils.HandleError(c, err, fiber.StatusBadRequest, "invalid request body")
	}

	if err := rules.Validate.Struct(input); err != nil {
		return utils.HandleError(c, err, fiber.StatusBadRequest, err.Error())
	}

	user := c.Locals("user").(*models.User)
	otp := models.OTP{
		UserID: user.ID,
		Code:   input.Code,
	}

	err := otp.Verify(c.Context())
	if err != nil {
		if err.Error() == "invalid OTP code" {
			return utils.HandleError(c, err, fiber.StatusBadRequest, "invalid OTP code")
		}

		return utils.HandleError(c, err, fiber.StatusInternalServerError, "internal server error")
	}

	return c.JSON(fiber.Map{
		"message": "Email verified successfully",
	})
}
