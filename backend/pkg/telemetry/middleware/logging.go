package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggingMiddleware() fiber.Handler {

	return func(c *fiber.Ctx) error {

		start := time.Now()

		err := c.Next()

		slog.Info(
			"HTTP Request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"latency", time.Since(start),
			"ip", c.IP(),
		)

		return err
	}
}
