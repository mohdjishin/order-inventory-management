package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey"`
	Name        string         `gorm:"unique;not null"`
	Description string         `gorm:""`
	Price       float64        `gorm:"not null"`
	Category    string         `gorm:"not null"`
	AddedBy     uint           `gorm:"not null"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete field

	InventoryID uint `gorm:"not null;index"`
	BlackListed bool `gorm:"default:false" json:"blacklisted,omitempty"` // Blacklisted flag for suppliers
}

func (Product) TableName() string {
	return "products"
}
