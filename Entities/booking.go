package entities

import (
	"gorm.io/gorm"
	"time"
)

type Booking struct {
	ID          string         `gorm:"primaryKey;column:id" json:"id"`
	SlotID      string         `gorm:"type:varchar(36);not null;index" json:"slot_id"` // Foreign key to Slot, indexed for performance
	Slot        *Slot          `gorm:"foreignKey:SlotID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"slot"`
	VendorID    string         `gorm:"type:varchar(36);not null;index" json:"vendor_id"` // Foreign key to Vendor, indexed for performance
	Vendor      *Vendor        `gorm:"foreignKey:VendorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"vendor"`
	BookingDate time.Time      `gorm:"type:timestamptz;not null" json:"booking_date"`                                                   // Store in UTC with timezone
	StartDate   time.Time      `gorm:"type:timestamptz;not null" json:"start_date"`                                                     // Store in UTC with timezone
	EndDate     time.Time      `gorm:"type:timestamptz;not null" json:"end_date"`                                                       // Store in UTC with timezone
	Status      string         `gorm:"type:varchar(20);not null;index" json:"status"`                                                   // Indexed for performance
	Payment     *Payment       `gorm:"foreignKey:BookingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"payment"` // Link to Payment
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
