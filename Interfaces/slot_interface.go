package Interfaces

import entities "tln-backend/Entities"

type ISlot interface {
	CreateSlot(slot []*entities.Slot) error
	CheckSlotName(name string) bool
	GetSlots(slotID string) (*entities.Slot, error)
	CheckMarketExists(marketID string) (bool, error)
	GetSlotsByDate(marketID, date string) ([]*entities.Slot, error)
	UpdateSlotStatus(slotID string, status entities.SlotStatus) error
}
