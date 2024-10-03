package entities

import "time"

type Market struct {
	ID          string     `gorm:"primaryKey;column:id" json:"id"`
	ProviderID  string     `gorm:"type:varchar(36);not null" json:"provider_id"`
	Name        string     `gorm:"type:varchar(100);not null" json:"name"`
	Phone       string     `gorm:"type:varchar(20)" json:"phone"`
	Address     string     `gorm:"type:varchar(255)" json:"address"`
	Description string     `gorm:"type:text" json:"description"`
	Image       string     `gorm:"type:varchar(255)" json:"image"`
	OpenTime    string     `gorm:"type:varchar(10)" json:"open_time"`
	CloseTime   string     `gorm:"type:varchar(10)" json:"close_time"`
	Latitude    string     `gorm:"type:varchar(20)" json:"latitude"`
	Longitude   string     `gorm:"type:varchar(20)" json:"longitude"`
	Slots       []Slot     `gorm:"foreignKey:MarketID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"slots"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
