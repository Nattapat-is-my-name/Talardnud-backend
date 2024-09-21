package Handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	entities "tln-backend/Entities"
	"tln-backend/Usecase"
)

// UserHandler handles user-related requests.
type UserHandler struct {
	useCase *Usecase.UserUseCase
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(uc *Usecase.UserUseCase) *UserHandler {
	return &UserHandler{useCase: uc}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body entities.User true "User data"
// @Success 201 {object} entities.User
// @Failure 400 {object} string "Invalid input"
// @Failure 500 {object} string "Internal server error"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user entities.RegisterRequest
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.useCase.CreateUser(&user); err != nil {
		if err.Error() == "username already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	fmt.Println("User created: ", user)

	return c.Status(fiber.StatusCreated).JSON(user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user with the provided ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} string "User deleted successfully"
// @Failure 403 {object} string "You are not authorized to delete this user"
// @Failure 500 {object} string "Internal server error"
// @Router /users/{id} [delete]

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {

	userIdFromToken := c.Locals("userID").(string)

	userIdToDelete := c.Params("id")

	if userIdFromToken != userIdToDelete {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to delete this user",
		})
	}

	err := h.useCase.DeleteUser(userIdToDelete)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
		"user_id": userIdToDelete,
	})
}
