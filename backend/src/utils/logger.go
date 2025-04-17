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
	Logger.SetLevel(logrus.TraceLevel)

	Logger.AddHook(&LevelHook{
		Writer:    os.Stdout,
		LogLevels: []logrus.Level{logrus.InfoLevel, logrus.DebugLevel, logrus.TraceLevel},
	})

	Logger.AddHook(&LevelHook{
		Writer: &lumberjack.Logger{
			Filename:   "logs/app.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     30,
			Compress:   true,
		},
		LogLevels: []logrus.Level{logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
	})
}

type LevelHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

func (h *LevelHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = h.Writer.Write([]byte(line))
	return err
}

func (h *LevelHook) Levels() []logrus.Level {
	return h.LogLevels
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
