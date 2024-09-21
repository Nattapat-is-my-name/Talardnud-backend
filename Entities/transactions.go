package entities

import (
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	ID              string         `gorm:"primaryKey;column:id" json:"id"`
	PaymentID       string         `gorm:"type:varchar(36);index" json:"payment_id"` // Foreign key to Payment
	Payment         *Payment       `gorm:"foreignKey:PaymentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"payment"`
	Amount          float64        `gorm:"type:decimal(10,2);not null;check:amount > 0" json:"amount"` // Amount of the transaction
	Method          string         `gorm:"type:varchar(50);index" json:"method"`                       // Transaction method (e.g., "Credit Card", "PromptPay")
	Status          string         `gorm:"type:varchar(20);not null" json:"status"`                    // Transaction status (e.g., "successful", "failed")
	TransactionDate time.Time      `gorm:"type:timestamptz;index" json:"transaction_date"`
	TransactionID   string         `gorm:"type:varchar(36)" json:"transaction_id,omitempty"` // PromptPay Transaction ID
	Ref1            string         `gorm:"type:varchar(50)" json:"ref1,omitempty"`           // PromptPay Reference 1
	Ref2            string         `gorm:"type:varchar(50)" json:"ref2,omitempty"`           // PromptPay Reference 2
	Ref3            string         `gorm:"type:varchar(50)" json:"ref3,omitempty"`           // PromptPay Reference 3
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete support
}
