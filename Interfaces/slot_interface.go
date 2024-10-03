package Interfaces

import entities "tln-backend/Entities"

type ISlot interface {
	CreateSlot(slot []*entities.Slot) error
	CheckSlotName(name string) bool
	//GetSlotWithMarketAndProviderByID(slotID string) (*entities.Slot, error)
	GetSlots(slotID string) (*entities.Slot, error)
	CheckMarketExists(marketID string) (bool, error)
	GetSlotsByDate(marketID, date string) ([]*entities.Slot, error)
}
