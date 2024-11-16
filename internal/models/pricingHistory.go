package models

import "time"

type PricingHistory struct {
	ID        uint      `gorm:"primaryKey"`
	ProductID uint      `gorm:"not null;index"` // Foreign key to Product
	OldPrice  float64   `gorm:"not null"`
	NewPrice  float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}

func (PricingHistory) TableName() string {
	return "pricing_histories"
}
