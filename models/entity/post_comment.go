package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type PostComment struct {
	ID      string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	PostID  string    `json:"post_id" gorm:"not null;size:20"`
	UserID  uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Content string    `json:"content" gorm:"not null;size:255"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
