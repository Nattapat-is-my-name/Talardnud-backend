package Repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	entities "tln-backend/Entities"
)

type SlotRepository struct {
	db *gorm.DB
}

func NewSlotRepository(db *gorm.DB) *SlotRepository {
	return &SlotRepository{db: db}
}

func (repo *SlotRepository) CreateSlot(slot []*entities.Slot) error {

	if err := repo.db.Create(slot).Error; err != nil {
		return err
	}
	return nil

}

func (r *SlotRepository) UpsertSlots(slots []*entities.Slot) ([]*entities.Slot, error) {
	updatedSlots := make([]*entities.Slot, 0, len(slots))

	for _, slot := range slots {
		result := r.db.Save(slot)
		if result.Error != nil {
			return nil, result.Error
		}
		updatedSlots = append(updatedSlots, slot)
	}

	return updatedSlots, nil
}

func (repo *SlotRepository) UpdateSlot(slot *entities.Slot) (*entities.Slot, error) {
	result := repo.db.Save(slot)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update slot: %w", result.Error)
	}

	// Retrieve the updated slot to ensure we return the most up-to-date data
	updatedSlot, err := repo.GetSlots(slot.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated slot: %w", err)
	}

	return updatedSlot, nil
}

// check name of slot
func (repo *SlotRepository) CheckSlotName(name string) bool {
	var slot entities.Slot
	repo.db.Where("name = ?", name).First(&slot)
	if slot.ID != "" {
		return true
	}
	return false
}

func (repo *SlotRepository) GetProviderSlots(marketID string) ([]*entities.Slot, error) {
	var slots []*entities.Slot
	err := repo.db.Where("market_id = ?", marketID).Find(&slots).Error
	if err != nil {
		return nil, err
	}

	return slots, nil
}

func (repo *SlotRepository) GetSlots(slotID string) (*entities.Slot, error) {
	var slot entities.Slot

	result := repo.db.Where("ID = ?", slotID).First(&slot)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("slot not found")
		}
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	return &slot, nil
}

func (repo *SlotRepository) DeleteSlot(slotID string) error {
	var slot entities.Slot
	result := repo.db.Where("ID = ?", slotID).Delete(&slot)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CheckMarketExists
func (repo *SlotRepository) CheckMarketExists(marketID string) (bool, error) {
	var market entities.Market
	repo.db.Where("id = ?", marketID).First(&market)
	if market.ID != "" {
		return true, nil
	}
	return false, nil
}

func (repo *SlotRepository) GetSlotsByMarketID(marketID string) ([]*entities.Slot, error) {
	var slots []*entities.Slot
	err := repo.db.Where("market_id = ?", marketID).Find(&slots).Error
	if err != nil {
		return nil, err
	}

	return slots, nil
}

func (repo *SlotRepository) UpdateSlots(slots []*entities.Slot) error {
	for _, slot := range slots {
		result := repo.db.Model(&slot).Updates(slot)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func (repo *SlotRepository) GetSlotsByDate(marketID, date string) ([]*entities.Slot, error) {
	var slots []*entities.Slot

	err := repo.db.Where("market_id = ? AND date = ?", marketID, date).Find(&slots).Error
	if err != nil {
		return nil, err
	}

	return slots, nil
}

func (repo *SlotRepository) DeleteSlotByDateAndZone(markeID, zoneID, date string) error {
	var slot entities.Slot
	result := repo.db.Where("market_id = ? AND zone = ? AND date = ?", markeID, zoneID, date).Delete(&slot)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *SlotRepository) UpdateSlotStatus(slotID string, status entities.SlotStatus) error {
	var slot entities.Slot
	result := repo.db.Model(&slot).Where("ID = ?", slotID).Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
