package dtos

import "time"

type SlotRequest struct {
	MarketID    string    `json:"market_id" validate:"required,uuid"`             // Required, the market associated with the slot
	Name        string    `json:"name" validate:"required"`                       // Required, the name of the slot
	Price       float64   `json:"price" validate:"required,gt=0"`                 // Required, the price of the slot
	Size        string    `json:"size" validate:"required"`                       // Required, the size of the slot
	Status      string    `json:"status" validate:"required"`                     // Required, the status of the slot
	SlotType    string    `json:"slot_type,omitempty"`                            // Optional, the type of the slot
	StartTime   time.Time `json:"start_time" validate:"required"`                 // Required, start time of the slot
	EndTime     time.Time `json:"end_time" validate:"required,gtfield=StartTime"` // Required, end time of the slot, must be after StartTime
	Description string    `json:"description,omitempty"`                          // Optional, a description for the slot (e.g., purpose, etc.)
}
