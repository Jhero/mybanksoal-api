package entity

import (
	"time"

	"gorm.io/gorm"
)

type QuestionStatus string

const (
	StatusDraft     QuestionStatus = "draft"
	StatusPublish   QuestionStatus = "publish"
	StatusUnpublish QuestionStatus = "unpublish"
)

type Question struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"not null" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Answer    string         `gorm:"type:text" json:"answer"`
	Difficulty int            `gorm:"default:1" json:"difficulty"`
	Status    QuestionStatus `gorm:"default:'draft';index" json:"status"`
	CreatedBy uint           `json:"created_by"`
	User      *User          `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
