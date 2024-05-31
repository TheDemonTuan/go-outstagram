package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type FriendStatus int

const (
	FriendRequested FriendStatus = iota
	FriendAccepted
	FriendRejected
)

type Friend struct {
	ID         uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	FromUserID uuid.UUID    `json:"from_user_id" gorm:"not null"`
	ToUserID   uuid.UUID    `json:"to_user_id" gorm:"not null"`
	Status     FriendStatus `json:"status" gorm:"default:0"`

	FromUser User `gorm:"foreignKey:FromUserID"`
	ToUser   User `gorm:"foreignKey:ToUserID"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
