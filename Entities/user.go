package entities

import (
	"gorm.io/gorm"
	"time"
)

type Vendor struct {
	ID        string         `gorm:"primaryKey;column:id" json:"id"`                   // UUID primary key
	Username  string         `gorm:"type:varchar(50);unique;not null" json:"username"` // Unique and non-null
	Email     string         `gorm:"type:varchar(100);unique;not null" json:"email"`   // Unique and non-null
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`              // Exclude from JSON responses
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`                 // Auto timestamp for creation
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`                 // Auto timestamp for updates
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`                // Soft delete support
	Bookings  []Booking      `gorm:"foreignKey:VendorID" json:"bookings"`              // Relationship to bookings
}
