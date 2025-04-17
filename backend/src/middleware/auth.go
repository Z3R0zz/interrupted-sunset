package middleware

import (
	"interrupted-export/src/config"
	"interrupted-export/src/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		token := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		claims := &jwt.MapClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
			return config.JwtSecret, nil
		})
		if err != nil || !parsedToken.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"middleware": true,
				"error":      "Unauthorized",
			})
		}

		expAt, ok := (*claims)["exp"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"middleware": true,
				"error":      "Unauthorized",
			})
		}

		if int64(expAt) < jwt.TimeFunc().Unix() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"middleware": true,
				"error":      "Unauthorized",
			})
		}

		userIDFloat, ok := (*claims)["user_id"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"middleware": true,
				"error":      "Unauthorized",
			})
		}

		user := &models.User{ID: uint(userIDFloat)}
		if err := user.Get(c.Context()); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"middleware": true,
				"error":      "Unauthorized",
			})
		}

		banned, err := user.IsBanned()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"middleware": true,
				"error":      "Internal server error",
			})
		}

		if banned {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"middleware": true,
				"error":      "User is banned",
			})
		}

		c.Locals("user", user)

		return c.Next()
	}
}
