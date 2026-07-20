package telemetry

import (
	"strconv"

	appmetrics "github.com/PainQlers/backend/pkg/telemetry/metrics"
	"github.com/gofiber/fiber/v2"
)

func MetricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		err := c.Next()

		appmetrics.HTTPRequestsTotal.
			WithLabelValues(
				c.Method(),
				c.Route().Path,
				strconv.Itoa(c.Response().StatusCode()),
			).
			Inc()

		return err
	}
}
