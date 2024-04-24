package repository

import (
	"time"

	prodrepo "github.com/swarit-pandey/e-commerce/product/pkg/repository"
	userrepo "github.com/swarit-pandey/e-commerce/user/pkg/repository"
)

type Order struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	UserID    uint          `json:"-"`
	User      userrepo.User `gorm:"foreignKey:UserID" json:"user"`
	Status    string        `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Items     []OrderItem   `gorm:"foreignKey:OrderID" json:"items"`
	Address   OrderAddress  `gorm:"foreignKey:OrderID" json:"address"`
}

type OrderItem struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	OrderID   uint             `json:"-"`
	ProductID uint             `json:"-"`
	Product   prodrepo.Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  uint             `json:"quantity"`
	Price     float64          `json:"price"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type OrderAddress struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OrderID      uint      `json:"-"`
	AddressLine1 string    `json:"address_line1"`
	AddressLine2 string    `json:"address_line2"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	Country      string    `json:"country"`
	PostalCode   string    `json:"postal_code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
