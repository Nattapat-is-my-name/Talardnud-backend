package dtos

import (
	entities "tln-backend/Entities"
)

type SlotData struct {
	SlotID   string              `json:"slot_id" validate:"required,regexp=^[A-G]-[1-8]$"`
	Price    float64             `json:"price" validate:"required,gt=0"`
	Status   entities.SlotStatus `json:"status" validate:"required,oneof=available booked maintenance"`
	Category entities.Category   `json:"category" validate:"required,oneof=clothes food crafts produce electronics services other"`
}

type DateRange struct {
	StartDate string `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `json:"end_date" validate:"required,datetime=2006-01-02,gtfield=StartDate"`
}

type SlotGenerationRequest struct {
	MarketID  string     `json:"market_id" validate:"required,uuid"`
	DateRange DateRange  `json:"date_range" validate:"required"`
	Slots     []SlotData `json:"slots" validate:"required,dive,min=1"`
}
