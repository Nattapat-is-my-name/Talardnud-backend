package dtos

import "time"

type SlotResponse struct {
	ID          string    `json:"id"`
	MarketID    string    `json:"market_id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Size        string    `json:"size"`
	Status      string    `json:"status"`
	SlotType    string    `json:"slot_type,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Description string    `json:"description,omitempty"`
}
