package dtos

import (
	"time"
)

type BookingResponse struct {
	ID          string               `json:"id"`
	SlotID      string               `json:"slotId"`
	VendorID    string               `json:"vendorId"`
	BookingDate time.Time            `json:"bookingDate"`
	StartDate   time.Time            `json:"startDate"`
	EndDate     time.Time            `json:"endDate"`
	Status      string               `json:"status"`
	Amount      float64              `json:"amount"`
	Method      string               `json:"method"`
	PromptPay   *PromptPayResponse   `json:"promptPay,omitempty"`
	Transaction *TransactionResponse `json:"transaction,omitempty"`
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
