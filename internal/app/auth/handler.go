package auth

import "github.com/gofiber/fiber/v2"

type authHandler struct {
	authService AuthServiceImpl
}

type AuthHandlerImpl interface {
	AnonymousLogin(c *fiber.Ctx) error
}

func NewAuthHandler(authService AuthServiceImpl) AuthHandlerImpl {
	return &authHandler{
		authService: authService,
	}
}

func (ah *authHandler) AnonymousLogin(c *fiber.Ctx) error {
	token, err := ah.authService.AnonymousLogin()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Successfully logged in",
		"tokens": &map[string]string{
			"accessToken": token,
		},
	})
}
