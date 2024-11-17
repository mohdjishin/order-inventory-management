package models

import (
	"time"

	"gorm.io/gorm"
)

type Inventory struct {
	ID          uint           `gorm:"primaryKey" json:"id,omitempty"`
	ProductID   uint           `gorm:"not null;index" json:"productId,omitempty"`
	Stock       int            `gorm:"not null" json:"stock,omitempty"`
	AddedBy     uint           `gorm:"not null" json:"addedBy,omitempty"`
	BasePrice   float64        `gorm:"not null" json:"basePrice,omitempty"`
	BlackListed bool           `gorm:"default:false" json:"blacklisted,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"-"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"-"`

	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;" json:"product,omitempty"`
}

func (Inventory) TableName() string {
	return "inventories"
}
