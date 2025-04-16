package utils

import (
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger = logrus.New()

func init() {
	Logger.SetFormatter(&logrus.JSONFormatter{})

	logFile := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	Logger.SetOutput(multiWriter)
	Logger.SetLevel(logrus.InfoLevel)
}

func HandleError(c *fiber.Ctx, err error, status int, message string) error {
	Logger.WithError(err).WithFields(logrus.Fields{
		"status": status,
		"path":   c.Path(),
		"method": c.Method(),
	}).Error(message)

	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}
