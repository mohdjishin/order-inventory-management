package models

type Product struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"unique;not null"` // Product name is unique and cannot be null
	Description string  `gorm:""`                // Description is optional
	Price       float64 `gorm:"not null"`        // Price cannot be null
	Category    string  `gorm:"not null"`        // Category cannot be null
	AddedBy     uint    `gorm:"not null"`        // User who added the product
	CreatedAt   int64   `gorm:"autoCreateTime"`  // Automatically set creation time
	UpdatedAt   int64   `gorm:"autoUpdateTime"`  // Automatically set update time
}

// TableName overrides the default table name (products) if necessary
func (Product) TableName() string {
	return "products"
}
