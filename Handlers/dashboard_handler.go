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

// GetWeeklyStats  godoc
// @Summary Get weekly stats for a market
// @Description Get weekly stats for a market with the market ID
// @Tags dashboard
// @Accept  json
// @Produce  json
// @Param id path string true "Market ID"
// @Success 200 {object} entities.DashboardResponse
// @Failure 400 {object} string "Invalid input"
// @Failure 500 {object} string "Internal server error"
// @Router /dashboard/weekly/{id} [get]
func (h *DashboardHandler) GetWeeklyStats(c *fiber.Ctx) error {
	marketID := c.Params("id")
	if marketID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Market ID is required",
		})
	}

	stats, err := h.useCase.GetWeeklyData(marketID)
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
