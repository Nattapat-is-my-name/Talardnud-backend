package entities

import (
	"time"
)

type Vendor struct {
	ID        string     `gorm:"primaryKey;column:id" json:"id"`
	Username  string     `gorm:"type:varchar(50);unique;not null" json:"username"`
	Email     string     `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password  string     `gorm:"type:varchar(255);not null" json:"-"`
	Image     string     `gorm:"type:text" json:"image"`
	FirstName string     `gorm:"type:varchar(50);" json:"first_name"`
	LastName  string     `gorm:"type:varchar(50);" json:"last_name"`
	Phone     string     `gorm:"type:varchar(15);unique;not null" json:"phone"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	Bookings  []Booking  `gorm:"foreignKey:VendorID" json:"bookings"`
}
