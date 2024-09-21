package entities

import (
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	ID           string         `gorm:"primaryKey;column:id" json:"id"`
	BookingID    string         `gorm:"type:varchar(36);not null;uniqueIndex" json:"booking_id"` // One-to-One with Booking
	Booking      *Booking       `gorm:"foreignKey:BookingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"booking"`
	Amount       float64        `gorm:"type:decimal(10,2);not null;check:amount > 0" json:"amount"`                             // Total amount to be paid
	Method       string         `gorm:"type:varchar(50);not null" json:"method"`                                                // e.g., "Credit Card", "PromptPay"
	Status       string         `gorm:"type:varchar(20);not null" json:"status"`                                                // e.g., "pending", "paid", "failed"
	PaymentDate  time.Time      `gorm:"type:timestamptz;not null" json:"payment_date"`                                          // Date and time of payment initiation
	Transactions []Transaction  `gorm:"foreignKey:PaymentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"transactions"` // List of transactions related to the payment
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
