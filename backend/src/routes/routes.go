package routes

import (
	"interrupted-export/src/controllers"
	"interrupted-export/src/middleware"

	"github.com/gofiber/fiber/v2"
)

var ApiRoutes = []Route{
	{Method: "POST", Path: "/login", Handler: controllers.Login},

	{Method: "POST", Path: "/mail/new", Handler: controllers.NewMail, Middlewares: []fiber.Handler{middleware.AuthMiddleware()}},
	{Method: "POST", Path: "/mail/verify", Handler: controllers.VerifyMail, Middlewares: []fiber.Handler{middleware.AuthMiddleware()}},
}
