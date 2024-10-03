package entities

import (
	"time"
)

type Slot struct {
	ID        string     `gorm:"primaryKey;column:id" json:"id"`
	SlotID    string     `gorm:"uniqueIndex;not null" json:"slot_id"`
	MarketID  string     `gorm:"type:varchar(36);not null" json:"market_id"`
	Name      string     `gorm:"type:varchar(100);" json:"name"`
	Price     float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	Status    SlotStatus `gorm:"type:varchar(20);not null" json:"status"`
	Category  Category   `gorm:"type:varchar(50);not null" json:"category"`
	Booking   []Booking  `gorm:"foreignKey:SlotID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"booking,omitempty"`
	Date      string     `gorm:"type:date;not null" json:"date" validate:"required,datetime=2006-01-02"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
type SlotStatus string
type Category string

const (
	StatusAvailable   SlotStatus = "available"
	StatusBooked      SlotStatus = "booked"
	StatusMaintenance SlotStatus = "maintenance"

	CategoryClothes     Category = "clothes"
	CategoryFood        Category = "food"
	CategoryCrafts      Category = "crafts"
	CategoryProduce     Category = "produce"
	CategoryElectronics Category = "electronics"
	CategoryServices    Category = "services"
	CategoryOther       Category = "other"
)
