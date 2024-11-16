package models

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"unique;not null"`
	Description string    `gorm:""`
	Price       float64   `gorm:"not null"`
	Category    string    `gorm:"not null"`
	AddedBy     uint      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	InventoryID uint      `gorm:"not null;index"`
}

func (Product) TableName() string {
	return "products"
}
