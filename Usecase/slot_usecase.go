package Usecase

import (
	"fmt"
	"strings"
	"time"
	entities "tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
	"tln-backend/Interfaces"
	"tln-backend/contact"
)

type SlotUseCase struct {
	repo Interfaces.ISlot
}

var _ contact.ISlotUseCase = (*SlotUseCase)(nil)

func NewSlotUseCase(repo Interfaces.ISlot) *SlotUseCase {
	return &SlotUseCase{
		repo: repo,
	}
}
func (su *SlotUseCase) CreateOrUpdateLayout(marketID string, layout []entitiesDtos.ZoneLayout) ([]*entities.Slot, *entitiesDtos.ErrorResponse) {
	// Validate the marketID
	marketExists, err := su.repo.CheckMarketExists(marketID)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to check market existence: " + err.Error(),
		}
	}
	if !marketExists {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    400,
			Message: fmt.Sprintf("Market with ID %s does not exist", marketID),
		}
	}

	// Validate the layout
	//if err := su.validateLayout(layout); err != nil {
	//	return nil, &entitiesDtos.ErrorResponse{
	//		Code:    400,
	//		Message: "Invalid layout: " + err.Error(),
	//	}
	//}

	// Get existing slots for the market
	existingSlots, err := su.repo.GetSlotsByMarketID(marketID)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch existing slots: " + err.Error(),
		}
	}

	// Create a map for quick lookup of existing slots
	existingSlotMap := make(map[string]*entities.Slot)
	for _, slot := range existingSlots {
		key := fmt.Sprintf("%s-%s-%s-%s", slot.MarketID, slot.Zone, slot.Name, slot.Date)
		existingSlotMap[key] = slot
	}

	slotsToUpsert := make([]*entities.Slot, 0)

	// Process layout and prepare slots for upsertion
	for _, zoneLayout := range layout {
		for _, stall := range zoneLayout.Stalls {
			uniqueSlotID := fmt.Sprintf("%s-%s-%s-%s", marketID, zoneLayout.Zone, stall.Name, zoneLayout.Date.Format("2006-01-02"))

			category, err := parseCategory(stall.StallType)
			if err != nil {
				return nil, &entitiesDtos.ErrorResponse{
					Code:    400,
					Message: fmt.Sprintf("Invalid category for stall %s: %s", stall.Name, err.Error()),
				}
			}

			slot := &entities.Slot{
				ID:        uniqueSlotID,
				MarketID:  marketID,
				Zone:      zoneLayout.Zone,
				Name:      stall.Name,
				Width:     stall.Width,
				Height:    stall.Height,
				Price:     stall.Price,
				Status:    entities.StatusAvailable,
				Category:  category,
				Date:      zoneLayout.Date.Format("2006-01-02"),
				UpdatedAt: time.Now(),
			}

			if existingSlot, exists := existingSlotMap[uniqueSlotID]; exists {
				// Update existing slot
				slot.CreatedAt = existingSlot.CreatedAt
			} else {
				// Set CreatedAt for new slots
				slot.CreatedAt = time.Now()
			}

			slotsToUpsert = append(slotsToUpsert, slot)
		}
	}

	// Upsert slots
	updatedSlots, err := su.repo.UpsertSlots(slotsToUpsert)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to upsert slots: " + err.Error(),
		}
	}

	return updatedSlots, nil
}

func (su *SlotUseCase) EditSlot(slotID string, updates *entitiesDtos.SlotUpdateDTO) (*entities.Slot, *entitiesDtos.ErrorResponse) {
	// Retrieve the existing slot
	existingSlot, err := su.repo.GetSlots(slotID)
	if err != nil {
		if err.Error() == "slot not found" {
			return nil, &entitiesDtos.ErrorResponse{
				Code:    404,
				Message: "Slot not found",
			}
		}
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to retrieve slot: " + err.Error(),
		}
	}

	// Apply updates
	if updates.Name != nil {
		existingSlot.Name = *updates.Name
	}
	if updates.Width != nil {
		existingSlot.Width = int(*updates.Width)
	}
	if updates.Height != nil {
		existingSlot.Height = int(*updates.Height)
	}
	if updates.Price != 0 {
		existingSlot.Price = updates.Price
	}
	if updates.Category != "" {
		existingSlot.Category = updates.Category
	}
	if updates.Status != "" {
		existingSlot.Status = updates.Status
	}

	// Update the slot in the repository
	updatedSlot, err := su.repo.UpdateSlot(existingSlot)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to update slot: " + err.Error(),
		}
	}

	return updatedSlot, nil
}

func (su *SlotUseCase) DeleteSlot(slotID string) *entitiesDtos.ErrorResponse {
	// Check if the slot exists
	_, err := su.repo.GetSlots(slotID)
	if err != nil {
		if err.Error() == "slot not found" {
			return &entitiesDtos.ErrorResponse{
				Code:    404,
				Message: "Slot not found",
			}
		}
		return &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to retrieve slot: " + err.Error(),
		}
	}

	// Delete the slot
	if err := su.repo.DeleteSlot(slotID); err != nil {
		return &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to delete slot: " + err.Error(),
		}
	}

	return nil
}

func (su *SlotUseCase) DeleteSlotByDateAndZone(slotID string, zoneID string, date string) *entitiesDtos.ErrorResponse {

	date = strings.TrimSpace(date)
	fmt.Printf("Received parameters: slotID=%s, zoneID=%s, date=%s\n", slotID, zoneID, date)
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return &entitiesDtos.ErrorResponse{
			Code:    400,
			Message: "Invalid date format. Please use YYYY-MM-DD",
		}
	}

	// Delete the slots
	if err := su.repo.DeleteSlotByDateAndZone(slotID, zoneID, date); err != nil {
		return &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to delete slots: " + err.Error(),
		}
	}

	return nil
}
func parseCategory(category string) (entities.Category, error) {
	switch strings.ToLower(category) {
	case "clothes", "clothing":
		return entities.CategoryClothes, nil
	case "food":
		return entities.CategoryFood, nil
	case "crafts", "craft":
		return entities.CategoryCrafts, nil
	case "produce":
		return entities.CategoryProduce, nil
	case "electronics", "electronic":
		return entities.CategoryElectronics, nil
	case "services", "service":
		return entities.CategoryServices, nil
	case "other":
		return entities.CategoryOther, nil
	default:
		return "", fmt.Errorf("unknown category: %s", category)
	}
}
func (su *SlotUseCase) GetSlots(slotID string) (*entities.Slot, *entitiesDtos.ErrorResponse) {
	slot, err := su.repo.GetSlots(slotID)
	if err != nil {
		if err.Error() == "slot not found" {
			return nil, &entitiesDtos.ErrorResponse{
				Code:    404,
				Message: "Slot not found",
			}
		}
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to retrieve slot: " + err.Error(),
		}
	}

	return slot, nil
}
func (su *SlotUseCase) GetSlotsByDate(marketID string, date string) ([]*entities.Slot, *entitiesDtos.ErrorResponse) {
	// Validate the date format
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    400,
			Message: "Invalid date format. Please use YYYY-MM-DD",
		}
	}

	// Call the repository method to get slots
	slots, err := su.repo.GetSlotsByDate(marketID, date)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to retrieve slots: " + err.Error(),
		}
	}

	// Check if any slots were found
	if len(slots) == 0 {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    404,
			Message: "No slots found for the given date and market",
		}
	}

	return slots, nil
}

func (su *SlotUseCase) UpdateSlotStatus(slotID string, status entities.SlotStatus) (*entities.Slot, *entitiesDtos.ErrorResponse) {
	// Call the repository method to update the slot status
	if err := su.repo.UpdateSlotStatus(slotID, status); err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to update slot status: " + err.Error(),
		}
	}

	// Retrieve the updated slot
	slot, err := su.GetSlots(slotID)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to retrieve updated slot: " + err.Message,
		}
	}

	return slot, nil
}

func (su *SlotUseCase) GetProviderSlots(marketID string) ([]*entities.Slot, *entitiesDtos.ErrorResponse) {
	// Call the repository method to get slots
	slots, err := su.repo.GetProviderSlots(marketID)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to retrieve slots: " + err.Error(),
		}
	}
	// Check if any slots were found
	if len(slots) == 0 {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    404,
			Message: "No slots found for the provider",
		}
	}

	return slots, nil
}
