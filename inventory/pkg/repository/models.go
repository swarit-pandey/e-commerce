package repository

import (
	"time"

	prodrepo "github.com/swarit-pandey/e-commerce/product/pkg/repository"
)

type Inventory struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	ProductID uint             `json:"-"`
	Product   prodrepo.Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  uint             `json:"quantity"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type InventoryLocation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type InventoryLocationProduct struct {
	InventoryLocationID uint      `gorm:"primaryKey" json:"-"`
	InventoryID         uint      `gorm:"primaryKey" json:"-"`
	Quantity            uint      `json:"quantity"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
