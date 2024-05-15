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
	IsHideLike    bool      `json:"is_hide_like" gorm:"default:false"`
	IsHideComment bool      `json:"is_hide_comment" gorm:"default:false"`
	Active        bool      `json:"active" gorm:"default:true"`

	PostFiles    []PostFile    `json:"post_images" gorm:"foreignKey:PostID;references:ID"`
	PostLikes    []PostLike    `json:"post_likes" gorm:"foreignKey:PostID;references:ID"`
	PostComments []PostComment `json:"post_comments" gorm:"foreignKey:PostID;references:ID"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
