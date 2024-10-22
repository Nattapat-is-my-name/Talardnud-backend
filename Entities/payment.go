package entities

import (
	"time"
)

type Payment struct {
	ID          string        `gorm:"primaryKey;column:id" json:"id"`
	BookingID   string        `gorm:"type:varchar(36);not null;uniqueIndex" json:"booking_id"`
	Booking     *Booking      `gorm:"foreignKey:BookingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"booking"`
	Price       float64       `gorm:"type:decimal(10,2);not null" json:"price"`
	Method      Method        `gorm:"type:varchar(50);not null" json:"method"`
	Status      PaymentStatus `gorm:"type:varchar(20);not null" json:"status"`
	PaymentDate time.Time     `gorm:"type:timestamptz;not null" json:"payment_date"`

	Transactions []Transaction `gorm:"foreignKey:PaymentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"transactions"`

	ExpiresAt time.Time  `gorm:"type:timestamp;not null" json:"expires_at"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentCompleted PaymentStatus = "completed"
	PaymentFailed    PaymentStatus = "failed"
)
