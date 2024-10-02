package Usecase

import (
	"github.com/google/uuid"
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
func (su *SlotUseCase) CreateSlot(slotReq *entitiesDtos.SlotRequest) (*entities.Slot, *entitiesDtos.ErrorResponse) {
	// Check if slot name already exists
	if su.repo.CheckSlotName(slotReq.Name) {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    400,
			Message: "Slot name already exists",
		}
	}

	// Create slot entity from request
	slotEntity := entities.Slot{
		ID:       uuid.New().String(),
		MarketID: slotReq.MarketID,
		Name:     slotReq.Name,
		Price:    slotReq.Price,
		Size:     slotReq.Size,
		Status:   "available",
	}

	// Create the slot in the database
	if err := su.repo.CreateSlot(&slotEntity); err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to create slot",
		}
	}

	createdSlot, err := su.repo.GetSlotWithMarketAndProviderByID(slotEntity.ID)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to retrieve slot details: " + err.Error(),
		}
	}

	return createdSlot, nil
}

func (su *SlotUseCase) GetSlots(markerID string) ([]*entities.Slot, *entitiesDtos.ErrorResponse) {
	slots, err := su.repo.GetSlots(markerID)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to retrieve slots: " + err.Error(),
		}
	}

	return slots, nil
}
