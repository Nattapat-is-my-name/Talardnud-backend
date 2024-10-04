package dtos

import entities "tln-backend/Entities"

type BookingRequest struct {
	SlotID      string          `gorm:"uniqueIndex;not null" json:"slot_id"`
	VendorID    string          `json:"vendor_id" validate:"required,uuid"` // Required, selected by the user
	BookingDate string          `json:"booking_date" validate:"required,datetime=2006-01-02"`
	Price       float64         `json:"price" validate:"required,gt=0"`
	Method      entities.Method `json:"method" validate:"required,oneof=PromptPay"`
	MarketID    string          `json:"market_id" validate:"required,uuid"` // Required, selected by the user

}

type CancelBookingRequest struct {
	BookingID string `json:"booking_id" validate:"required"` // The ID of the booking to be canceled.
	VendorID  string `json:"user_id,omitempty"`              // Optional: The ID of the user requesting the cancellation.

}
