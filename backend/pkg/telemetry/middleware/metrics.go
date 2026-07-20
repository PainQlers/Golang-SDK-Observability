package middleware

import (
	"strconv"
	"time"

	instrument "github.com/PainQlers/backend/pkg/telemetry/instrument"
	"github.com/gofiber/fiber/v2"
)

func MetricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		start := time.Now()

		err := c.Next()

		status := c.Response().StatusCode()

		route := "unknown"
		if c.Route() != nil {
			route = c.Route().Path
		}

		instrument.RecordHTTPRequest(
			c.Method(),
			route,
			strconv.Itoa(status),
			time.Since(start).Seconds(),
		)

		return err
	}
}
