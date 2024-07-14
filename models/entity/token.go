package entity

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID       string    `json:"user_id" gorm:"not null;size:20"`
	RefreshToken string    `json:"refresh_token" gorm:"unique;index;not null;size:255"`
	ExpiredAt    time.Time `json:"expired_at" gorm:"not null"`
	Active       bool      `json:"active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
