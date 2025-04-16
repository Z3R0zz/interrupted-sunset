package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JwtSecret = []byte("")

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	JwtSecret = []byte(os.Getenv("JWT_SECRET"))
}
