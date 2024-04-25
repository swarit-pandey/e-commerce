package repository

import "time"

type User struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	Username     string        `gorm:"uniqueIndex" json:"username"`
	Email        string        `gorm:"uniqueIndex" json:"email"`
	PasswordHash string        `json:"-"`
	Name         string        `json:"name"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Addresses    []UserAddress `gorm:"foreignKey:UserID" json:"addresses"`
}

type UserAddress struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `json:"-"`
	AddressLine1 string    `json:"address_line1"`
	AddressLine2 string    `json:"address_line2"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	Country      string    `json:"country"`
	PostalCode   string    `json:"postal_code"`
	IsPrimary    bool      `gorm:"default:false" json:"is_primary"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
