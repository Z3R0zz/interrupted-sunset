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

func User(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	queue := models.Queue{UserID: user.ID}
	ctx := c.Context()

	exists, err := queue.ExistsInQueue(ctx)
	if err != nil {
		return utils.HandleError(c, err, fiber.StatusInternalServerError, "failed to check queue existence")
	}

	status := ""
	if exists {
		if status, err = queue.GetStatus(ctx); err != nil {
			return utils.HandleError(c, err, fiber.StatusInternalServerError, "failed to get queue status")
		}
	}

	return c.JSON(fiber.Map{
		"ID":       user.ID,
		"Username": user.Username,
		"Email":    user.Email,
		"Verified": user.EmailVerifiedAt != nil,
		"Queue":    exists,
		"Status":   status,
	})
}
