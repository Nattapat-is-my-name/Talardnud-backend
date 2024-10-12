package dtos

import (
	"time"
	entities "tln-backend/Entities"
)

type BookingResponse struct {
	ID            string                 `json:"id"`
	SlotID        string                 `json:"slotId"`
	VendorID      string                 `json:"vendorId"`
	TransactionID string                 `json:"transactionId"`
	BookingDate   time.Time              `json:"bookingDate"`
	Price         float64                `json:"price"`
	Status        entities.BookingStatus `json:"status"`
	Method        entities.Method        `json:"method"`
	Image         string                 `json:"image,omitempty"`
	ExpiresAt     time.Time              `json:"expiresAt"`
}

type TransactionResponse struct {
	ID              string    `json:"id"`
	PaymentID       string    `json:"paymentId"`
	Amount          float64   `json:"amount"`
	Method          string    `json:"method"`
	Status          string    `json:"status"`
	TransactionDate time.Time `json:"transactionDate"`
	TransactionID   string    `json:"transactionId"`
	Ref1            string    `json:"ref1"`
	Ref2            string    `json:"ref2"`
	Ref3            string    `json:"ref3"`
}
