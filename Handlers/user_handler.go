package Handlers

import (
	"github.com/gofiber/fiber/v2"
	"tln-backend/Entities/dtos"
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

//// CreateUser godoc
//// @Summary Create a user
//// @Description Create a new user with the provided data
//// @Tags users
//// @Accept  json
//// @Produce  json
//// @Param user body entities.RegisterRequest true "User data"
//// @Success 201 {object} dtos.RegisterRequest
//// @Failure 400 {object} string "Invalid input"
//// @Failure 409 {object} string "Username already exists"
//// @Failure 500 {object} string "Internal server error"
//// @Router /users [post]
//func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
//	var user dtos.RegisterRequest
//	if err := c.BodyParser(&user); err != nil {
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
//	}
//
//	if err := h.useCase.CreateUser(&user); err != nil {
//		if err.Error() == "username already exists" {
//			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
//		}
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
//	}
//
//	fmt.Println("User created: ", user)
//
//	return c.Status(fiber.StatusCreated).JSON(user)
//}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user with the provided ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} string "User deleted successfully"
// @Failure 403 {object} string "You are not authorized to delete this user"
// @Failure 500 {object} string "Failed to delete user"
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

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Get a user with the provided ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} dtos.GetUserResponse
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/{id} [get]
// @Security BearerAuth
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	userId := c.Params("id")

	// Fetch user by ID
	user, err := h.useCase.GetUserByID(userId)
	if err != nil {
		// Return a 404 status if the user is not found with a detailed error message
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "User not found",
				"message": "No user found with the provided ID",
				"data":    dtos.GetUserResponse{},
			})
		}

		// On internal server error, return a 500 status with a detailed error message
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
			"data":    dtos.GetUserResponse{},
		})
	}

	response := dtos.GetUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Bookings:  user.Bookings,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
