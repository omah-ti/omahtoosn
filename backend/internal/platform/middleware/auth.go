package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/security"
)

const ClaimsKey = "claims"

// mengembalikan fiber middleware yang memvalidasi JWT
func NewAuthMiddleware(tokens security.TokenProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string

		authHeader := c.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				tokenString = parts[1]
			}
		}

		if tokenString == "" {
			tokenString = c.Cookies("access_token")
		}

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "missing authorization token",
			})
		}

		claims, err := tokens.ParseAccessToken(tokenString)
		if err != nil {
			switch err {
			case security.ErrTokenExpired:
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"message": "token has expired",
				})
			case security.ErrTokenMalformed:
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"message": "token is malformed",
				})
			default:
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"message": "token is invalid",
				})
			}
		}

		c.Locals(ClaimsKey, claims)
		return c.Next()
	}
}

// ambil JWT claims dari fiber context
func ClaimsFromContext(c *fiber.Ctx) (*security.Claims, bool) {
	claims, ok := c.Locals(ClaimsKey).(*security.Claims)
	return claims, ok
}
