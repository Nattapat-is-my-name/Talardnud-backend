package dtos

import "time"

type BookingRequest struct {
	SlotID    string    `json:"slot_id" validate:"required,uuid"`               // Required, selected by the user
	VendorID  string    `json:"vendor_id" validate:"required,uuid"`             // Required, selected by the user
	StartDate time.Time `json:"start_date" validate:"required"`                 // Required, start of booking
	EndDate   time.Time `json:"end_date" validate:"required,gtfield=StartDate"` // Required, end of booking
	Amount    float64   `json:"amount" validate:"required,gt=0"`                // Required, amount to be paid
	Method    string    `json:"method" validate:"required"`                     // Required, payment method
	MarketID  string    `json:"market_id" validate:"required,uuid"`             // Required, selected by the user

}

type CancelBookingRequest struct {
	BookingID string `json:"booking_id" validate:"required"` // The ID of the booking to be canceled.
	VendorID  string `json:"user_id,omitempty"`              // Optional: The ID of the user requesting the cancellation.

}
