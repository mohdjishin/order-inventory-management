package models

type Inventory struct {
	ID        uint    `gorm:"primaryKey"`
	ProductID uint    `gorm:"not null;index"` // Foreign key to Product
	Stock     int     `gorm:"not null"`
	UpdatedAt int64   `gorm:"autoUpdateTime"`
	Product   Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}

func (Inventory) TableName() string {
	return "inventories"
}
