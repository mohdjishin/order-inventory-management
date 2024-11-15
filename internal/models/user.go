package models

type Role uint

func (r Role) String() string {
	return [...]string{"admin", "customer", "supplier"}[r-1]
}
func GetRoleID(s string) uint {

	if s == "admin" {
		return 0
	} else if s == "customer" {
		return 1
	} else if s == "supplier" {
		return 2
	}
	return 0
}

const (
	SuperAdmin Role = iota + 1
	CustomerRole
	SupplierRole
)

type User struct {
	Id        uint   `gorm:"primaryKey" json:"id,omitempty"`             // User ID
	FirstName string `gorm:"not null" json:"firstName,omitempty"`        // First name cannot be null
	LastName  string `gorm:"not null" json:"lastName,omitempty"`         // Last name cannot be null
	Email     string `gorm:"unique;not null" json:"email,omitempty"`     // Email is unique and cannot be null
	Password  string `gorm:"not null" json:"-"`                          // Password is required, omitted in JSON response
	Phone     string `gorm:"not null" json:"phone,omitempty"`            // Phone number is required
	Role      Role   `gorm:"not null" json:"-"`                          // Role is required
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at,omitempty"` // Automatically set creation time
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at,omitempty"` // Automatically set update time

	// This field is used to approve the user by the admin only for suppliers
	Approved bool `gorm:"default:false" json:"approved"` // Approved is false by default
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
