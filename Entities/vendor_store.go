package entities

import "time"

type Stall struct {
	ID          string    `gorm:"primaryKey;column:id" json:"id"`
	VendorID    string    `gorm:"type:varchar(36);not null" json:"vendor_id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"` // Name of the stall/store
	Image       string    `gorm:"type:text" json:"image"`                 // Image of the stall
	Description string    `gorm:"type:text" json:"description"`           // Description of the stall
	Type        string    `gorm:"type:varchar(50)" json:"type"`           // Type of the stall (food, clothing, etc.)
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
