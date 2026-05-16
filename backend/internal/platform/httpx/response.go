package httpx

import "github.com/gofiber/fiber/v2"

func Success(c *fiber.Ctx, status int, message string, data any) error {
	payload := fiber.Map{
		"success":    true,
		"message":    message,
		"request_id": requestID(c),
	}
	if data != nil {
		payload["data"] = data
	}
	return c.Status(status).JSON(payload)
}

func requestID(c *fiber.Ctx) string {
	if c == nil {
		return ""
	}
	if value, ok := c.Locals("requestid").(string); ok {
		return value
	}
	return ""
}
