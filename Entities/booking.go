package entities

import (
	"time"
)

type Booking struct {
	ID          string     `gorm:"primaryKey;column:id" json:"id"`
	SlotID      string     `gorm:"type:varchar(36);not null" json:"slot_id"`
	Slot        *Slot      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"slot"`
	VendorID    string     `gorm:"type:varchar(36);not null;index" json:"vendor_id"`
	Vendor      *Vendor    `gorm:"foreignKey:VendorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"vendor"`
	BookingDate time.Time  `gorm:"type:timestamptz;not null" json:"booking_date"`
	StartDate   time.Time  `gorm:"type:timestamptz;not null" json:"start_date"`
	EndDate     time.Time  `gorm:"type:timestamptz;not null" json:"end_date"`
	Status      string     `gorm:"type:varchar(20);not null;index" json:"status"`
	Payment     *Payment   `gorm:"foreignKey:BookingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"payment"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
