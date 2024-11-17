package models

import "time"

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
	Id        uint      `gorm:"primaryKey" json:"id,omitempty"`
	FirstName string    `gorm:"not null" json:"firstName,omitempty"`
	LastName  string    `gorm:"not null" json:"lastName,omitempty"`
	Email     string    `gorm:"unique;not null" json:"email,omitempty"`
	Password  string    `gorm:"not null" json:"-"`
	Phone     string    `gorm:"not null" json:"phone,omitempty"`
	Role      Role      `gorm:"not null" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	// This field is used to approve the supplier by the admin only for suppliers
	Approved    bool `gorm:"default:false" json:"approved"`              // Approved is false by default
	BlackListed bool `gorm:"default:false" json:"blacklisted,omitempty"` // Blacklisted flag for suppliers

}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
