package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID            string    `json:"id" gorm:"primaryKey;not null;size:20"`
	UserID        uuid.UUID `json:"user_id" gorm:"not null;type:uuid"`
	Caption       string    `json:"caption" gorm:"not null;size:2200"`
	Files         string    `json:"files" gorm:"not null;"`
	IsHideLike    bool      `json:"is_hide_like" gorm:"default:false"`
	IsHideComment bool      `json:"is_hide_comment" gorm:"default:false"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
