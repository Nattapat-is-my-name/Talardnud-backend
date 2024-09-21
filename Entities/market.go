package entities

type Market struct {
	ID          string         `gorm:"primaryKey;column:id" json:"id"`
	ProviderID  string         `gorm:"type:varchar(36);not null" json:"provider_id"`
	Provider    MarketProvider `gorm:"foreignKey:ProviderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"provider"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Address     string         `gorm:"type:varchar(255)" json:"address"`
	Description string         `gorm:"type:text" json:"description"`
	Image       string         `gorm:"type:varchar(255)" json:"image"`
	Slots       []Slot         `gorm:"foreignKey:MarketID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"slots"`
	OpenTime    string         `gorm:"type:varchar(10)" json:"open_time"`
	CloseTime   string         `gorm:"type:varchar(10)" json:"close_time"`
	Latitude    string         `gorm:"type:varchar(20)" json:"latitude"`
	Longitude   string         `gorm:"type:varchar(20)" json:"longitude"`
}
