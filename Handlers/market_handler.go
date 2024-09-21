package Handlers

import (
	"github.com/gofiber/fiber/v2"
	entitiesDtos "tln-backend/Entities/dtos"
	"tln-backend/Usecase"
)

type MarketHandler struct {
	useCase *Usecase.MarketUseCase
}

func NewMarketHandler(useCase *Usecase.MarketUseCase) *MarketHandler {
	return &MarketHandler{useCase: useCase}
}

func (h *MarketHandler) CreateMarket(c *fiber.Ctx) error {
	var req entitiesDtos.MarketRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input: unable to parse request body",
		})
	}

	market, errRes := h.useCase.CreateMarket(&req)
	if errRes != nil {
		return c.Status(errRes.Code).JSON(fiber.Map{
			"status":  "error",
			"message": errRes.Message,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Market created successfully",
		"data":    market,
	})
}
