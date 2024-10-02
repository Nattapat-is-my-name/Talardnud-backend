package Handlers

import (
	"github.com/gofiber/fiber/v2"
	entities "tln-backend/Entities"
	"tln-backend/Entities/dtos"
	"tln-backend/Usecase"
)

type AuthHandler struct {
	useCase *Usecase.AuthUseCase
}

func NewAuthHandler(uc *Usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		useCase: uc,
	}
}

// Login godoc
// @Summary Login
// @Description Login with the provided credentials
// @Tags auth
// @Accept  json
// @Produce  json
// @Param login body entities.LoginRequest true "Login data"
// @Success 200 {object} entities.LoginResponse
// @Failure 400 {object} string "Invalid input"
// @Failure 401 {object} string "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req entities.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid input",
			"details": "Failed to parse login request",
		})
	}

	// Authenticate the user with either username or email
	response, err := h.useCase.Login(req.UsernameOrEmail, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Login failed",
			"details": "Invalid username/email or password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":       "success",
		"message":      "Login successful",
		"access_token": response.AccessToken,
		"vendor_id":    response.VendorID,
	})
}

// Register godoc
// @Summary Register
// @Description Register a new user with the provided data
// @Tags auth
// @Accept  json
// @Produce  json
// @Param register body dtos.RegisterRequest true "Register request"
// @Success 200 {object} dtos.RegisterResponse
// @Failure 400 {object} string "Failed to register user"
// @Failure 409 {object} string "Email already exists"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dtos.RegisterRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid input",
			"details": "Failed to parse registration request",
		})
	}

	// Register the user and capture the detailed error response
	res, errResponse := h.useCase.Register(req.Username, req.Password, req.Email, req.PhoneNumber)
	if errResponse != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to register user",
			"details": errResponse,
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
