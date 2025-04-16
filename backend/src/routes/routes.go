package routes

import (
	"interrupted-export/src/controllers"
	"interrupted-export/src/middleware"

	"github.com/gofiber/fiber/v2"
)

var ApiRoutes = []Route{
	{Method: "POST", Path: "/login", Handler: controllers.Login, Middlewares: []fiber.Handler{middleware.AuthMiddleware()}},
}
