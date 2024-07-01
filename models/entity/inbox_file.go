package entity

import (
	"github.com/google/uuid"
	"time"
)

type InboxFileType int

const (
	InboxFileVideo InboxFileType = iota
	InboxFileImage
)

type InboxFile struct {
	ID      uint          `json:"id" gorm:"primaryKey"`
	InboxID uuid.UUID     `json:"inbox_id" gorm:"type:uuid;not null;index"`
	Type    InboxFileType `json:"type" gorm:"not null;default:1"`
	URL     string        `json:"url" gorm:"not null;size:255"`
	Active  bool          `json:"active" gorm:"default:true"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt time.Time `json:"deleted_at" gorm:"index"`
}
