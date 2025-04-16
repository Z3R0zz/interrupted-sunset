package routes

import (
	"interrupted-export/src/controllers"
	"interrupted-export/src/middleware"

	"github.com/gofiber/fiber/v2"
)

var ApiRoutes = []Route{
	{Method: "POST", Path: "/login", Handler: controllers.Login},

	{Method: "POST", Path: "/email/new", Handler: controllers.NewEmail, Middlewares: []fiber.Handler{middleware.AuthMiddleware()}},
}
