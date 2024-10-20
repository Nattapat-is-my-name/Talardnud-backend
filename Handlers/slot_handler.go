package Handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/url"
	entitiesDtos "tln-backend/Entities/dtos"
	"tln-backend/Usecase"
)

type SlotHandler struct {
	useCase *Usecase.SlotUseCase
}

func NewSlotHandler(useCase *Usecase.SlotUseCase) *SlotHandler {
	return &SlotHandler{useCase: useCase}
}

// CreateOrUpdateLayout godoc
// @Summary Create or update layout
// @Description Create or update layout
// @Tags slots
// @Accept json
// @Produce json
// @Param marketId path string true "Market ID"
// @Param layout body dtos.LayoutRequest true "Layout data"
// @Success 200 {object} string
// @Router /slots/{marketId}/create [post]
func (h *SlotHandler) CreateOrUpdateLayout(c *fiber.Ctx) error {
	marketID := c.Params("marketId")
	fmt.Println("Market ID: ", marketID)
	var layoutRequest entitiesDtos.LayoutRequest
	if err := c.BodyParser(&layoutRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON: " + err.Error(),
		})
	}

	updatedSlots, errResponse := h.useCase.CreateOrUpdateLayout(marketID, layoutRequest.Layout)
	if errResponse != nil {
		return c.Status(errResponse.Code).JSON(fiber.Map{
			"error": errResponse.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "Layout processed successfully",
		"marketId":     marketID,
		"updatedSlots": updatedSlots,
	})
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

// EditSlot godoc
// @Summary Edit slot
// @Description Edit slot
// @Tags slots
// @Accept json
// @Produce json
// @Param id path string true "Slot ID"
// @Param updateDTO body dtos.SlotUpdateDTO true "Slot update data"
// @Success 200 {object} entities.Slot
// @Router /slots/edit/{id} [patch]
func (h *SlotHandler) EditSlot(c *fiber.Ctx) error {
	slotID := c.Params("id")

	// Parse the update data from the request body
	var updateDTO entitiesDtos.SlotUpdateDTO
	if err := c.BodyParser(&updateDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&entitiesDtos.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body: " + err.Error(),
		})
	}

	// Call the usecase layer with both slotID and updateDTO
	updatedSlot, errRes := h.useCase.EditSlot(slotID, &updateDTO)
	if errRes != nil {
		return c.Status(errRes.Code).JSON(&entitiesDtos.ErrorResponse{
			Code:    errRes.Code,
			Message: errRes.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedSlot)
}

// DeleteSlot godoc
// @Summary Delete slot
// @Description Delete slot
// @Tags slots
// @Accept json
// @Produce json
// @Param id path string true "Slot ID"
// @Success 200 {object} string
// @Router /slots/delete/{id} [delete]
func (h *SlotHandler) DeleteSlot(c *fiber.Ctx) error {
	slotID := c.Params("id")
	errRes := h.useCase.DeleteSlot(slotID)
	if errRes != nil {
		return c.Status(errRes.Code).JSON(&entitiesDtos.ErrorResponse{
			Code:    errRes.Code,
			Message: errRes.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Slot deleted successfully",
	})
}

// DeleteSlotByDateAndZone godoc
// @Summary Delete slot by date and zone
// @Description Delete slot by date and zone
// @Tags slots
// @Accept json
// @Produce json
// @Param id path string true "Slot ID"
// @Param zoneID path string true "Zone ID"
// @Param date path string true "Date"
// @Success 200 {object} string
// @Router /slots/delete/{id}/zone/{zoneID}/date/{date} [delete]
func (h *SlotHandler) DeleteSlotByDateAndZone(c *fiber.Ctx) error {
	slotID := c.Params("id")
	zoneID := c.Params("zoneID")
	date, err := url.QueryUnescape(c.Params("date"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&entitiesDtos.ErrorResponse{
			Code:    400,
			Message: "Invalid date format in URL",
		})
	}

	errRes := h.useCase.DeleteSlotByDateAndZone(slotID, zoneID, date)
	if errRes != nil {
		return c.Status(errRes.Code).JSON(&entitiesDtos.ErrorResponse{
			Code:    errRes.Code,
			Message: errRes.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Slot deleted successfully",
	})
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

// GetProviderSlots godoc
// @Summary Get provider slots
// @Description Get provider slots
// @Tags slots
// @Accept json
// @Produce json
// @Param id path string true "Provider ID"
// @Success 200 {object} []entities.Slot
// @Router /slots/provider/get/{id} [get]
func (h *SlotHandler) GetProviderSlots(c *fiber.Ctx) error {
	providerID := c.Params("id")
	slots, errRes := h.useCase.GetProviderSlots(providerID)
	if errRes != nil {
		return c.Status(errRes.Code).JSON(
			&entitiesDtos.ErrorResponse{
				Code:    errRes.Code,
				Message: errRes.Message,
			})
	}

	return c.Status(200).JSON(slots)
}
