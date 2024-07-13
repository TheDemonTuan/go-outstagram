package entity

import (
	"github.com/google/uuid"
	"time"
)

type PostSave struct {
	ID     uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	PostID string    `json:"post_id" gorm:"not null;size:20;index"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`

	Post Post `json:"-" gorm:"foreignKey:PostID;references:ID"`
	User User `json:"-" gorm:"foreignKey:UserID;references:ID"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
