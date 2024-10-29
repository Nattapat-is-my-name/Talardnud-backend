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
	res, errResponse := h.useCase.Register(req.Username, req.Password, req.Email, req.PhoneNumber, req.Firstname, req.Lastname)
	if errResponse != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to register user",
			"details": errResponse,
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// ProviderLogin godoc
// @Summary Provider Login
// @Description Login for market providers with the provided credentials
// @Tags auth
// @Accept  json
// @Produce  json
// @Param login body dtos.ProviderLoginRequest true "Provider Login data"
// @Success 200 {object} dtos.ProviderLoginResponse
// @Failure 400 {object} string "Invalid input"
// @Failure 401 {object} string "Invalid credentials"
// @Router /auth/provider/login [post]
func (h *AuthHandler) ProviderLogin(c *fiber.Ctx) error {
	var req dtos.ProviderLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid input",
			"details": "Failed to parse login request",
		})
	}

	// Authenticate the provider
	response, err := h.useCase.ProviderLogin(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Login failed",
			"details": "Invalid email or password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":       "success",
		"message":      "Provider login successful",
		"access_token": response.AccessToken,
		"provider_id":  response.ProviderID,
	})

}

// RegisterProvider godoc
// @Summary Register Provider
// @Description Register a new market provider with the provided data
// @Tags auth
// @Accept  json
// @Produce  json
// @Param register body dtos.MarketProviderRequest true "Register provider request"
// @Success 200 {object} entities.MarketProvider
// @Failure 400 {object} string "Failed to register provider"
// @Failure 409 {object} string "Provider email already exists"
// @Router /auth/provider/register [post]
func (h *AuthHandler) RegisterProvider(c *fiber.Ctx) error {
	var req dtos.MarketProviderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid input",
			"details": "Failed to parse registration request",
		})
	}

	res, errResponse := h.useCase.RegisterProvider(req.Username, req.Phone, req.Email, req.Password)
	if errResponse != nil {
		return c.Status(errResponse.Code).JSON(fiber.Map{
			"error":   "Failed to register provider",
			"details": errResponse.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
