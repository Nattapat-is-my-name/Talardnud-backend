package Repository

import (
	"gorm.io/gorm"
	entities "tln-backend/Entities"
)

type SlotRepository struct {
	db *gorm.DB
}

func NewSlotRepository(db *gorm.DB) *SlotRepository {
	return &SlotRepository{db: db}
}

func (repo *SlotRepository) CreateSlot(slot *entities.Slot) error {

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

func (repo *SlotRepository) GetSlotWithMarketAndProviderByID(slotID string) (*entities.Slot, error) {
	var slot entities.Slot

	// Use nested Preload to load Market and its associated Provider
	err := repo.db.Preload("Market").Preload("Market.Provider").Where("id = ?", slotID).First(&slot).Error
	if err != nil {
		return nil, err
	}

	return &slot, nil
}

func (repo *SlotRepository) GetSlots(marketID string) ([]*entities.Slot, error) {
	var slots []*entities.Slot

	err := repo.db.Where("market_id = ?", marketID).Find(&slots).Error
	if err != nil {
		return nil, err
	}

	return slots, nil
}
