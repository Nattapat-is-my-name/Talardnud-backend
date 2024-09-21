package Handlers

import (
	"github.com/gofiber/fiber/v2"
	entitiesDtos "tln-backend/Entities/dtos"
	"tln-backend/Usecase"
)

type SlotHandler struct {
	useCase *Usecase.SlotUseCase
}

func NewSlotHandler(useCase *Usecase.SlotUseCase) *SlotHandler {
	return &SlotHandler{useCase: useCase}
}

func (h *SlotHandler) CreateSlot(c *fiber.Ctx) error {
	var slotReq entitiesDtos.SlotRequest

	if err := c.BodyParser(&slotReq); err != nil {
		return c.Status(400).JSON(&entitiesDtos.ErrorResponse{
			Code:    400,
			Message: "Invalid request body",
		})
	}

	slot, errRes := h.useCase.CreateSlot(&slotReq)
	if errRes != nil {
		return c.Status(errRes.Code).JSON(errRes)
	}

	return c.Status(201).JSON(slot)
}
