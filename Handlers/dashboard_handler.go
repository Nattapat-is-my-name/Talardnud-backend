package Handlers

import (
	"github.com/gofiber/fiber/v2"
	"tln-backend/Usecase"
)

type DashboardHandler struct {
	useCase *Usecase.DashboardUseCase
}

func NewDashboardHandler(useCase *Usecase.DashboardUseCase) *DashboardHandler {
	return &DashboardHandler{useCase: useCase}
}

func (h *DashboardHandler) GetDashboardData(c *fiber.Ctx) error {
	data, err := h.useCase.GetDashboardData("") // Empty marketID to get all markets
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

// Add a new handler for single market if needed
func (h *DashboardHandler) GetSingleMarketStats(c *fiber.Ctx) error {
	marketID := c.Params("id")
	if marketID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Market ID is required",
		})
	}

	stats, err := h.useCase.GetSingleMarketStats(marketID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   stats,
	})
}
