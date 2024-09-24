package dtos

import "time"

type SlotRequest struct {
	MarketID    string    `json:"market_id" validate:"required,uuid"`
	Name        string    `json:"name" validate:"required"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	Size        string    `json:"size" validate:"required"`
	Status      string    `json:"status" validate:"required"`
	SlotType    string    `json:"slot_type,omitempty"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	Description string    `json:"description,omitempty"`
}
