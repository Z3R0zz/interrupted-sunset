package main

import (
	"fmt"
	"interrupted-export/src/config"
	"interrupted-export/src/database"
	"interrupted-export/src/routes"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.Load()

	app := fiber.New(fiber.Config{
		Prefork:               true,
		CaseSensitive:         true,
		StrictRouting:         true,
		DisableStartupMessage: false,
		ServerHeader:          "INTERRUPTED",
		AppName:               "INTERRUPTED API",
	})

	database.Connect(os.Getenv("DATABASE_URL"))

	routes.SetupAppRoutes(app)

	err := app.Listen(os.Getenv("PORT"))
	if err != nil {
		fmt.Println("Error starting web server:", err)
	}
}
