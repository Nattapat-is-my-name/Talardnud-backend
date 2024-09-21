package entities

type MarketProvider struct {
	ID      string   `gorm:"primaryKey;column:id" json:"id"`
	Name    string   `gorm:"type:varchar(100);not null" json:"name"`
	Phone   string   `gorm:"type:varchar(20)" json:"phone"`
	Email   string   `gorm:"type:varchar(100)" json:"email"`
	Address string   `gorm:"type:varchar(255)" json:"address"`
	Markets []Market `gorm:"foreignKey:ProviderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"markets"`
}
