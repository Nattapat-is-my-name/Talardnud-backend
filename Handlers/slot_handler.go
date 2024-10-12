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

// CreateSlot godoc
// @Summary Create a slot
// @Description Create a new slot with the provided data
// @Tags slots
// @Accept  json
// @Produce  json
// @Param slot body dtos.SlotGenerationRequest true "Slot data"
// @Success 201 {object} []entities.Slot
// @Failure 400 {object} string "Invalid input"
// @Failure 409 {object} string "Slot already exists"
// @Failure 500 {object} string "Internal server error"
// @Router /slots/create [post]
func (h *SlotHandler) CreateSlot(c *fiber.Ctx) error {
	var slotReq entitiesDtos.SlotGenerationRequest

	if err := c.BodyParser(&slotReq); err != nil {
		return c.Status(400).JSON(&entitiesDtos.ErrorResponse{
			Code:    400,
			Message: "Invalid request body",
		})
	}

	slot, errRes := h.useCase.CreateSlots(&slotReq)
	if errRes != nil {
		return c.Status(errRes.Code).JSON(errRes)
	}

	return c.Status(201).JSON(slot)
}

// GetSlot godoc
// @Summary Get all slots
// @Description Get all slots
// @Tags slots
// @Accept json
// @Produce json
// @Param id path string true "Market ID"
// @Success 200 {object} entities.Slot
// @Router /slots/get/{id} [get]
func (h *SlotHandler) GetSlot(c *fiber.Ctx) error {
	marketID := c.Params("id")
	slots, errRes := h.useCase.GetSlots(marketID)
	if errRes != nil {
		return c.Status(errRes.Code).JSON(errRes)
	}

	return c.Status(200).JSON(slots)
}

// GetSlotByDate godoc
// @Summary Get slots by date
// @Description Get slots by date
// @Tags slots
// @Accept json
// @Produce json
// @Param marketID path string true "MarketID"
// @Param date path string true "Date"
// @Success 200 {object} []entities.Slot
// @Router /slots/markets/{marketID}/date/{date} [get]
func (h *SlotHandler) GetSlotByDate(c *fiber.Ctx) error {
	marketID := c.Params("marketID")
	selectDate := c.Params("date")
	slots, errRes := h.useCase.GetSlotsByDate(marketID, selectDate)
	if errRes != nil {
		return c.Status(errRes.Code).JSON(
			&entitiesDtos.ErrorResponse{
				Code:    errRes.Code,
				Message: errRes.Message,
			})
	}

	return c.Status(200).JSON(slots)
}
