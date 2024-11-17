package models

import "time"

type ShippingState uint

const (
	SPending   ShippingState = iota + 1
	SShipping                // 2
	SDelivered               // 3
	SCancelled               // 4
	SReturned                // 5
)

func (s ShippingState) String() string {
	return []string{"PENDING", "SHIPPING", "DELIVERED", "CANCELLED", "RETURNED"}[s-1]
}

var ShippingStates = map[string]ShippingState{
	"pending":   SPending,
	"shipping":  SShipping,
	"delivered": SDelivered,
	"cancelled": SCancelled,
	"returned":  SReturned,
}

type Order struct {
	ID             uint             `gorm:"primaryKey" json:"id,omitempty"`
	UserID         uint             `gorm:"not null;index" json:"userId,omitempty"`
	ProductID      uint             `gorm:"not null;index" json:"productId,omitempty"`
	Quantity       int              `gorm:"not null" json:"quantity,omitempty"`
	TotalPrice     float64          `gorm:"not null" json:"totalPrice,omitempty"`
	Status         string           `gorm:"default:'PENDING'" json:"status,omitempty"`
	CreatedAt      time.Time        `gorm:"autoCreateTime" json:"-"`
	UpdatedAt      time.Time        `gorm:"autoUpdateTime" json:"-" `
	SupplierID     uint             `gorm:"not null;index" json:"supplierId,omitempty"`
	User           User             `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
	Product        Product          `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;" json:"-"`
	Supplier       User             `gorm:"foreignKey:SupplierID;constraint:OnDelete:CASCADE;" json:"-"`
	ReturnStatus   string           `gorm:"default:''" json:"returnStatus,omitempty"`
	ShippingStates []ShipmentStatus `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;" json:"shippingStates,omitempty"`
	ShippingState  ShippingState    `gorm:"default:1" json:"shippingState,omitempty" `
}

func (Order) TableName() string {
	return "orders"
}

type ShipmentStatus struct {
	ID             uint          `gorm:"primaryKey"`
	OrderID        uint          `gorm:"not null;index"`
	Status         ShippingState `gorm:"not null"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime"`
	AdditionalInfo string        `gorm:"type:text"`
}

func (ShipmentStatus) TableName() string {
	return "shipment_statuses"
}
