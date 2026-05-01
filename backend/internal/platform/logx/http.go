package logx

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		startedAt := time.Now()
		err := c.Next()
		duration := time.Since(startedAt)
		requestID, _ := c.Locals("requestid").(string)
		log.Printf("req=%s method=%s path=%s status=%d dur=%s", requestID, c.Method(), c.Path(), c.Response().StatusCode(), duration.Round(time.Millisecond))
		return err
	}
}
