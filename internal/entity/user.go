package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:255;uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"` // Hide password in JSON
	Role      string         `gorm:"size:50;default:'user'" json:"role"` // e.g., "admin", "editor", "viewer"
	APIKey    string         `gorm:"size:255;uniqueIndex" json:"-"`      // API Key for mobile apps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
