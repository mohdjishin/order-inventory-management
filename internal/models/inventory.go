package models

import (
	"time"

	"gorm.io/gorm"
)

type Inventory struct {
	ID          uint           `gorm:"primaryKey"`
	ProductID   uint           `gorm:"not null;index"`
	Stock       int            `gorm:"not null"`
	AddedBy     uint           `gorm:"not null"`
	BasePrice   float64        `gorm:"not null"`
	BlackListed bool           `gorm:"default:false" json:"blacklisted,omitempty"` // Blacklisted flag for suppliers
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                             // Soft delete field
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	//Initial stock
	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}

// TableName overrides the default table name (inventories) if necessary
func (Inventory) TableName() string {
	return "inventories"
}
