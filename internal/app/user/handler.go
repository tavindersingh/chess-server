package user

import "github.com/gofiber/fiber/v2"

type userHandler struct {
	userRepository UserRepository
}

type UserHandler interface {
	CurrentUser(c *fiber.Ctx) error
}

func NewUserHandler(userRepository UserRepository) UserHandler {
	return &userHandler{
		userRepository: userRepository,
	}
}

func (uh *userHandler) CurrentUser(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	user, err := uh.userRepository.GetUser(userId)

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Successfully fetched user",
		"user":    user,
	})
}
