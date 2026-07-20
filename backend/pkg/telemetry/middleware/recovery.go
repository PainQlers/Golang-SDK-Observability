package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func RecoveryMiddleware() fiber.Handler {

	return func(c *fiber.Ctx) error {

		defer func() {

			if r := recover(); r != nil {

				slog.Error(
					"panic recovered",
					"panic", r,
				)

				_ = c.Status(fiber.StatusInternalServerError).
					JSON(fiber.Map{
						"error": "internal server error",
					})
			}

		}()

		return c.Next()
	}
}
