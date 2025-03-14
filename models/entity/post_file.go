package entity

import (
	"gorm.io/gorm"
	"time"
)

type PostFile struct {
	ID     string `json:"id" gorm:"primaryKey"`
	PostID string `json:"post_id" gorm:"not null;size:20;index"`
	URL    string `json:"url" gorm:"not null;size:255"`
	Active bool   `json:"active" gorm:"default:true"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
