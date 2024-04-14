package token

import (
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

func NewTokenHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("token")

		if token == "" {
			slog.Error("user has no token!!!")
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized access")
		}

		return c.Next()
	}
}
