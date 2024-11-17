package models

import "time"

type PricingHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id,omitempty"`
	ProductID uint      `gorm:"not null;index" json:"productId,omitempty"`
	OldPrice  float64   `gorm:"not null" json:"oldPrice,omitempty"`
	NewPrice  float64   `gorm:"not null" json:"newPrice,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;" json:"-"`
}

func (PricingHistory) TableName() string {
	return "pricing_histories"
}
