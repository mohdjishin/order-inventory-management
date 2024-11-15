package models

type PricingHistory struct {
	ID        uint    `gorm:"primaryKey"`
	ProductID uint    `gorm:"not null;index"` // Foreign key to Product
	Price     float64 `gorm:"not null"`
	CreatedAt int64   `gorm:"autoCreateTime"`
	Product   Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}

func (PricingHistory) TableName() string {
	return "pricing_histories"
}
