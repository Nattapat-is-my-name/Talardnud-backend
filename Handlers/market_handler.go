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

// CreateMarket godoc
// @Summary Create a new market
// @Description Create a new market
// @Tags Market
// @Accept json
// @Produce json
// @Param market body entitiesDtos.MarketRequest true "Market object that needs to be created"
// @Success 201 {object} entitiesDtos.MarketResponse
// @Router /markets/create [post]
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

// GetMarket godoc
// @Summary Get all markets
// @Schemes
// @Tags Market
// @Accept json
// @Produce json
// @Success 200 {object} dtos.GetListMarketResponse
// @Router /markets/get [get]
func (h *MarketHandler) GetMarket(c *fiber.Ctx) error {
	// Call the useCase to get the markets
	markets, errRes := h.useCase.GetMarket()
	if errRes != nil {
		return c.Status(errRes.Code).JSON(entitiesDtos.MarketResponse{
			Status:  "error",
			Message: errRes.Message,
			Data:    nil,
		})
	}

	// Return success response with data
	return c.Status(fiber.StatusOK).JSON(entitiesDtos.MarketResponse{
		Status:  "success",
		Message: "Market fetched successfully",
		Data:    markets, // Could be an empty array or list of markets
	})
}

// EditMarket godoc
// @Summary Edit a market
// @Description Edit a market
// @Tags Market
// @Accept json
// @Produce json
// @Param id path string true "Market ID"
// @Param market body dtos.MarketEditRequest true "Market object that needs to be updated"
// @Success 200 {object} entities.Market
// @Router /markets/edit/{id} [patch]
func (h *MarketHandler) EditMarket(c *fiber.Ctx) error {
	marketID := c.Params("id")

	var req entitiesDtos.MarketEditRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input: unable to parse request body",
		})
	}

	market, errRes := h.useCase.EditMarket(marketID, &req)
	if errRes != nil {
		return c.Status(errRes.Code).JSON(fiber.Map{
			"status":  "error",
			"message": errRes.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Market updated successfully",
		"data":    market,
	})
}

// GetMarketByID godoc
// @Summary Get a market by ID
// @Description Get a market by ID
// @Tags Market
// @Accept json
// @Produce json
// @Param id path string true "Market ID"
// @Success 200 {object} entities.Market
// @Router /markets/get/{id} [get]
func (h *MarketHandler) GetMarketByID(c *fiber.Ctx) error {
	marketID := c.Params("id")

	// Call the useCase to get the market by ID
	market, errRes := h.useCase.GetMarketByID(marketID)
	if errRes != nil {
		return c.Status(errRes.Code).JSON(fiber.Map{
			"status":  "error",
			"message": errRes.Message,
		})
	}

	// Return success response with data
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Market fetched successfully",
		"data":    market,
	})
}

// GetMarketByProviderID godoc
// @Summary Get a market by Provider ID
// @Description Get a market by Provider ID
// @Tags Market
// @Accept json
// @Produce json
// @Param id path string true "Provider ID"
// @Success 200 {object} entitiesDtos.MarketResponse
// @Router /markets/provider/get/{id} [get]
func (h *MarketHandler) GetMarketByProviderID(c *fiber.Ctx) error {
	providerID := c.Params("id")

	// Call the useCase to get the market by ID
	market, errRes := h.useCase.GetMarketByProviderID(providerID)
	if errRes != nil {
		return c.Status(errRes.Code).JSON(fiber.Map{
			"status":  "error",
			"message": errRes.Message,
		})
	}

	// Return success response with data
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Market fetched successfully",
		"data":    market,
	})
}
