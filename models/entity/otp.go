package entity

import (
	"github.com/google/uuid"
	"time"
)

type Otp struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserEmail string    `json:"user_email" gorm:"unique;index;not null;size:100"`
	OtpCode   string    `json:"otp" gorm:"size:6;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	IsUsed    bool      `json:"is_used" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
