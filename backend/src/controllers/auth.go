package controllers

import (
	"interrupted-export/src/models"
	"interrupted-export/src/rules"
	"interrupted-export/src/utils"

	"github.com/gofiber/fiber/v2"
)

type LoginInput struct {
	Username string `json:"username" validate:"max=45,min=3,required"`
	Password string `json:"password" validate:"max=255,min=8,required"`
}

func Login(c *fiber.Ctx) error {
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return utils.HandleError(c, err, fiber.StatusBadRequest, "invalid request body")
	}

	if err := rules.Validate.Struct(input); err != nil {
		return utils.HandleError(c, err, fiber.StatusBadRequest, err.Error())
	}

	user := models.User{
		Username: input.Username,
		Password: input.Password,
	}

	token, err := user.AttemptLogin(c.Context())
	if err != nil {
		return utils.HandleError(c, err, fiber.StatusUnauthorized, err.Error())
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}
