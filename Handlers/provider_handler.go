package Handlers

import (
	"github.com/gofiber/fiber/v2"
	entities "tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
	"tln-backend/Usecase"
)

type MarketProvider struct {
	useCase *Usecase.ProviderUseCase
}

func NewMarketProvider(useCase *Usecase.ProviderUseCase) *MarketProvider {
	return &MarketProvider{useCase: useCase}
}

// CreateProvider creates a new market provider and handles errors appropriately.
func (uc *MarketProvider) CreateProvider(c *fiber.Ctx) error {
	var marketProvider entitiesDtos.MarketProviderRequest

	if err := c.BodyParser(&marketProvider); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&entitiesDtos.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid input: unable to parse request body",
		})
	}

	createdProvider, errRes := uc.useCase.CreateProvider(&marketProvider)
	if errRes != nil {
		return c.Status(errRes.Code).JSON(errRes)
	}

	return c.Status(fiber.StatusCreated).JSON(createdProvider)
}

// UpdateProvider updates an existing market provider and handles errors appropriately.
func (uc *MarketProvider) UpdateProvider(c *fiber.Ctx) error {
	var marketProvider entities.MarketProvider

	if err := c.BodyParser(&marketProvider); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&entitiesDtos.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid input: unable to parse request body",
		})
	}

	updatedProvider, errRes := uc.useCase.UpdateProvider(&marketProvider)
	if errRes != nil {
		return c.Status(errRes.Code).JSON(errRes)
	}

	return c.Status(fiber.StatusOK).JSON(updatedProvider)
}
