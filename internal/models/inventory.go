package models

type Inventory struct {
	ID        uint  `gorm:"primaryKey"`
	ProductID uint  `gorm:"not null;index"` // Foreign key to Product
	Stock     int   `gorm:"not null"`       // Current stock for the product
	AddedBy   uint  `gorm:"not null"`       // User who added/updated the inventory
	UpdatedAt int64 `gorm:"autoUpdateTime"` // Automatically set update time

	// Relationship to Product
	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"` // Ensures Product deletion cascades
}

// TableName overrides the default table name (inventories) if necessary
func (Inventory) TableName() string {
	return "inventories"
}
