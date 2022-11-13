package middleware

import "github.com/gofiber/fiber/v2"

func IsHeader(c *fiber.Ctx) error {
	bearer := c.Get("secret-key")
	if bearer == "" {
		return c.Status(401).JSON(fiber.Map{"error": "header not set"})
	}
	return c.Next()
}

func SecretKey(c *fiber.Ctx) error {
	userKey := c.Get("secret-key")
	if userKey == "" {
		return c.Status(401).JSON(fiber.Map{"error": "secret-key not set"})
	}
	vesselKey := c.Get("vessel-secret-key")
	if vesselKey == "" {
		return c.Status(401).JSON(fiber.Map{"error": "vessel-secret-key not set"})
	}
	return c.Next()
}
