package entity

import (
	"gorm.io/gorm"
	"time"
)

type PostFileType int

const (
	PostFileVideo PostFileType = iota
	PostFileImage
)

type PostFile struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	PostID    string         `json:"post_id" gorm:"not null;size:20"`
	URL       string         `json:"url" gorm:"not null;size:255"`
	Type      PostFileType   `json:"type" gorm:"not null;default:0"`
	Active    bool           `json:"active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
