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

// check name of slot
func (repo *SlotRepository) CheckSlotName(name string) bool {
	var slot entities.Slot
	repo.db.Where("name = ?", name).First(&slot)
	if slot.ID != "" {
		return true
	}
	return false
}

//func (repo *SlotRepository) GetSlotWithMarketAndProviderByID(slotID string) (*entities.Slot, error) {
//	var slot entities.Slot
//
//	// Use nested Preload to load Market and its associated Provider
//	err := repo.db.Preload("Market").Preload("Market.Provider").Where("id = ?", slotID).First(&slot).Error
//	if err != nil {
//		return nil, err
//	}
//
//	return &slot, nil
//}

func (repo *SlotRepository) GetSlots(slotID string) (*entities.Slot, error) {
	var slot entities.Slot

	result := repo.db.Where("slot_id = ?", slotID).First(&slot)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("slot not found")
		}
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	return &slot, nil
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

func (repo *SlotRepository) GetSlotsByDate(marketID, date string) ([]*entities.Slot, error) {
	var slots []*entities.Slot

	err := repo.db.Where("market_id = ? AND date = ?", marketID, date).Find(&slots).Error
	if err != nil {
		return nil, err
	}

	return slots, nil
}

func (repo *SlotRepository) UpdateSlotStatus(slotID string, status entities.SlotStatus) error {
	var slot entities.Slot
	result := repo.db.Model(&slot).Where("ID = ?", slotID).Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
