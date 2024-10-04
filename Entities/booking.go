package entities

import (
	"time"
)

type Booking struct {
	ID          string        `gorm:"primaryKey;column:id" json:"id"`
	SlotID      string        `gorm:"not null;index" json:"slot_id"`
	Slot        *Slot         `gorm:"foreignKey:SlotID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"slot"`
	VendorID    string        `gorm:"type:varchar(36);not null;index" json:"vendor_id"`
	MarketID    string        `gorm:"type:varchar(36);not null" json:"market_id"`
	Vendor      *Vendor       `gorm:"foreignKey:VendorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"vendor"`
	BookingDate time.Time     `gorm:"type:date;not null" json:"booking_date"` // Changed from Date to BookingDate
	Status      BookingStatus `gorm:"type:varchar(20);not null;index" json:"status"`
	Method      Method        `gorm:"type:varchar(20);not null" json:"method"`
	Price       float64       `gorm:"type:decimal(10,2);not null" json:"price"`
	Payment     *Payment      `gorm:"foreignKey:BookingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"payment"`
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	ExpiresAt   time.Time     `gorm:"type:timestamp;not null" json:"expires_at"`
}
type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusCancelled BookingStatus = "cancelled"
	StatusCompleted BookingStatus = "completed"
)

type Method string

const (
	MethodPromptPay Method = "PromptPay"
)
