package models

import (
	"gorm.io/gorm"
)

type Order struct {
	ID         uint    `gorm:"primaryKey"`
	UserID     uint    `gorm:"not null;index"` // Foreign key for User
	ProductID  uint    `gorm:"not null;index"` // Foreign key for Product
	Quantity   int     `gorm:"not null"`       // Quantity of product ordered
	TotalPrice float64 `gorm:"not null"`       // Total price for the order
	Status     string  `gorm:"default:'PENDING'"`
	CreatedAt  int64   `gorm:"autoCreateTime"`
	UpdatedAt  int64   `gorm:"autoUpdateTime"`

	User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"` // Changed from Customer to User
	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}

// TableName specifies the table name for the Order model.
func (Order) TableName() string {
	return "orders"
}

// Migration function to add constraints and indexes
func AddOrderConstraints(db *gorm.DB) {
	// Add CHECK constraint for Status field (only allows 'PENDING', 'COMPLETED', 'CANCELLED')
	err := db.Exec(`
		ALTER TABLE orders
		ADD CONSTRAINT check_status
		CHECK (status IN ('PENDING', 'COMPLETED', 'CANCELLED'));
	`).Error
	if err != nil {
		panic("Failed to add check constraint for status")
	}

	// Optionally, if required, add unique index on UserID + ProductID for order uniqueness
	err = db.Exec(`
		CREATE UNIQUE INDEX idx_user_product ON orders (user_id, product_id);
	`).Error
	if err != nil {
		panic("Failed to add unique index for user and product")
	}
}
