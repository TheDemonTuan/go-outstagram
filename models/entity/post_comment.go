package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type PostComment struct {
	ID       uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	PostID   string    `json:"post_id" gorm:"not null;size:20;index"`
	UserID   uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	ParentID uuid.UUID `json:"parent_id" gorm:"type:uuid;default:null;index"`
	Content  string    `json:"content" gorm:"not null;size:255"`
	Active   bool      `json:"active" gorm:"default:true"`

	Children []PostComment `json:"children" gorm:"foreignKey:ParentID"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
