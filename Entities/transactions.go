package entities

import (
	"time"
)

type Transaction struct {
	ID              string            `gorm:"primaryKey;column:id" json:"id"`
	PaymentID       string            `gorm:"type:varchar(36);index" json:"payment_id"`
	Payment         *Payment          `gorm:"foreignKey:PaymentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"payment"`
	Price           float64           `gorm:"type:decimal(10,2);not null" json:"price"`
	Method          string            `gorm:"type:varchar(50);index" json:"method"`
	Status          TransactionStatus `gorm:"type:varchar(20);not null" json:"status"`
	TransactionDate time.Time         `gorm:"type:timestamptz;index" json:"transaction_date"`
	TransactionID   string            `gorm:"type:varchar(36)" json:"transaction_id,omitempty"`
	Ref1            string            `gorm:"type:varchar(50)" json:"ref1,omitempty"`
	Ref2            string            `gorm:"type:varchar(50)" json:"ref2,omitempty"`
	Ref3            string            `gorm:"type:varchar(50)" json:"ref3,omitempty"`
	Image           string            ` gorm:"type:text" json:"image,omitempty"`
	ExpiresAt       time.Time         `gorm:"type:timestamp;not null" json:"expires_at"`
	CreatedAt       time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       *time.Time        `gorm:"index" json:"deleted_at,omitempty"`
}

type TransactionStatus string

const (
	TransactionPending   TransactionStatus = "pending"
	TransactionCompleted TransactionStatus = "completed"
	TransactionFailed    TransactionStatus = "failed"
	TransactionRefunded  TransactionStatus = "refund"
)
