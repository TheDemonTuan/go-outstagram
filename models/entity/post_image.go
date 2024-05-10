package entity

import (
	"gorm.io/gorm"
	"time"
)

type PostImage struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	PostID string `json:"post_id" gorm:"not null;size:20"`
	Image  string `json:"image" gorm:"not null;size:255"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
