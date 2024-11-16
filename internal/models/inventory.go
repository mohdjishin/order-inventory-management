package models

import "time"

type Inventory struct {
	ID        uint      `gorm:"primaryKey"`
	ProductID uint      `gorm:"not null;index"`
	Stock     int       `gorm:"not null"`
	AddedBy   uint      `gorm:"not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	BasePrice float64   `gorm:"not null"`
	//Initial stock
	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}

// TableName overrides the default table name (inventories) if necessary
func (Inventory) TableName() string {
	return "inventories"
}
