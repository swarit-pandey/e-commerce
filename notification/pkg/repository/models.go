package repository

import (
	"time"

	userrepo "github.com/swarit-pandey/e-commerce/user/pkg/repository"
)

type Notification struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	UserID    uint          `json:"-"`
	User      userrepo.User `gorm:"foreignKey:UserID" json:"user"`
	Message   string        `json:"message"`
	IsRead    bool          `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type NotificationType struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserNotificationPreference struct {
	UserID             uint      `gorm:"primaryKey" json:"-"`
	NotificationTypeID uint      `gorm:"primaryKey" json:"-"`
	IsEnabled          bool      `gorm:"default:true" json:"is_enabled"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
