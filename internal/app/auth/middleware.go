package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	jwtManager *JwtManager
}

func NewAuthMiddleware(jwtManager *JwtManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (am *AuthMiddleware) RequireAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	userId, err := am.jwtManager.ValidateToken(token)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	
	c.Locals("userId", userId)
	return c.Next()
}
