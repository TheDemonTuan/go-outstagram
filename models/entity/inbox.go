package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Inbox struct {
	ID         uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	FromUserID uuid.UUID `json:"from_user_id" gorm:"not null"`
	ToUserID   uuid.UUID `json:"to_user_id" gorm:"not null"`
	Message    string    `json:"message" gorm:"not null;type:text"`
	IsRead     bool      `json:"is_read" gorm:"default:false"`

	Files []InboxFile `json:"files" gorm:"foreignKey:InboxID;references:ID"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
