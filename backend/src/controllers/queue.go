package controllers

import (
	"interrupted-export/src/models"
	"interrupted-export/src/utils"

	"github.com/gofiber/fiber/v2"
)

func Queue(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	queue := models.Queue{UserID: user.ID}

	if user.EmailVerifiedAt == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Email verification is required to create a queue entry.",
		})
	}

	exists, err := queue.ExistsInQueue(c.Context())
	if err != nil {
		return utils.HandleError(c, err, fiber.StatusInternalServerError, "failed to check queue existence")
	}

	if exists {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You already have a queue entry.",
		})
	}

	if err := queue.Insert(c.Context()); err != nil {
		return utils.HandleError(c, err, fiber.StatusInternalServerError, "failed to insert queue entry")
	}

	return c.JSON(fiber.Map{
		"message": "Queue entry created successfully.",
	})
}
