package models

type Role uint

// String method for role to print as "admin", "customer", "supplier"
func (r Role) String() string {
	return [...]string{"admin", "customer", "supplier"}[r]
}

const (
	AdminRole Role = iota + 1
	CustomerRole
	SupplierRole
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	FirstName string `gorm:"not null"`        // First name cannot be null
	LastName  string `gorm:"not null"`        // Last name cannot be null
	Email     string `gorm:"unique;not null"` // Email is unique and cannot be null
	Password  string `gorm:"not null"`        // Password is required
	Phone     string `gorm:"not null"`        // Phone number is required
	Role      Role   `gorm:"not null"`        // Role is required
	CreatedAt int64  `gorm:"autoCreateTime"`  // Automatically set creation time
	UpdatedAt int64  `gorm:"autoUpdateTime"`  // Automatically set update time
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
