package entities

import (
	"time"
)

type Slot struct {
	ID        string     `gorm:"primaryKey;column:id" json:"id"`
	MarketID  string     `gorm:"type:varchar(36);not null" json:"market_id"`
	Market    Market     `gorm:"foreignKey:MarketID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"market"`
	Name      string     `gorm:"type:varchar(100);not null" json:"name"`
	Price     float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	Size      string     `gorm:"type:varchar(50);not null" json:"size"`
	Status    string     `gorm:"type:varchar(20);not null" json:"status"`
	SlotType  string     `gorm:"type:varchar(50)" json:"slot_type,omitempty"`
	Bookings  []Booking  `gorm:"foreignKey:SlotID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"bookings"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
