package Usecase

import (
	"fmt"
	"github.com/google/uuid"
	"time"
	entities "tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
	"tln-backend/Interfaces"
)

type SlotUseCase struct {
	repo Interfaces.ISlot
}

func NewSlotUseCase(repo Interfaces.ISlot) *SlotUseCase {
	return &SlotUseCase{
		repo: repo,
	}
}
func (su *SlotUseCase) CreateSlots(slotReq *entitiesDtos.SlotGenerationRequest) ([]*entities.Slot, *entitiesDtos.ErrorResponse) {
	// Check if the market exists
	marketExists, err := su.repo.CheckMarketExists(slotReq.MarketID)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to check market existence",
		}
	}
	if !marketExists {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    400,
			Message: fmt.Sprintf("Market with ID %s does not exist", slotReq.MarketID),
		}
	}

	// Parse start and end dates
	startDate, err := time.Parse("2006-01-02", slotReq.DateRange.StartDate)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    400,
			Message: "Invalid start date format",
		}
	}

	endDate, err := time.Parse("2006-01-02", slotReq.DateRange.EndDate)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    400,
			Message: "Invalid end date format",
		}
	}

	// Calculate the number of days in the range
	days := int(endDate.Sub(startDate).Hours() / 24)

	// Calculate total number of slots
	totalSlots := (days + 1) * len(slotReq.Slots)

	// Create a slice to hold all the slots
	slots := make([]*entities.Slot, 0, totalSlots)

	// Generate slots for each day and each SlotData
	for i := 0; i <= days; i++ {
		date := startDate.AddDate(0, 0, i)
		for _, slotData := range slotReq.Slots {
			uniqueSlotID := fmt.Sprintf("%s-%s-%s", slotReq.MarketID, slotData.SlotID, date.Format("2006-01-02"))
			slotEntity := &entities.Slot{
				ID:       uuid.New().String(),
				SlotID:   uniqueSlotID,
				MarketID: slotReq.MarketID,

				Price:     slotData.Price,
				Status:    slotData.Status,
				Category:  slotData.Category,
				Date:      date.Format("2006-01-02"),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			slots = append(slots, slotEntity)
		}
	}

	// Create the slots in the database
	if err := su.repo.CreateSlot(slots); err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to create slots",
		}
	}

	return slots, nil
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
