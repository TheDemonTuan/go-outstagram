package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type PostLike struct {
	ID     string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	PostID string    `json:"post_id" gorm:"not null;size:20"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`

	Post Post `json:"-" gorm:"foreignKey:PostID;references:ID"`
	User User `json:"-" gorm:"foreignKey:UserID;references:ID"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
