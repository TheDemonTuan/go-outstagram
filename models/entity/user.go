package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID       uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username string    `json:"username" gorm:"unique;not null;size:50" `
	Password string    `json:"-" gorm:"not null;size:255"`
	FullName string    `json:"full_name" gorm:"not null;size:100"`
	Email    string    `json:"email" gorm:"unique;not null;size:100"`
	Phone    string    `json:"phone" gorm:"unique;size:15"`
	Avatar   string    `json:"avatar" gorm:"size:255"`
	Bio      string    `json:"bio" gorm:"size:255"`
	Birthday time.Time `json:"birthday" gorm:"not null"`
	Gender   bool      `json:"gender"`
	Role     bool      `json:"role" gorm:"default:false"`
	Active   bool      `json:"active" gorm:"default:true"`

	Posts        []Post        `json:"posts" gorm:"foreignKey:UserID;references:ID"`
	PostLikes    []PostLike    `json:"post_likes" gorm:"foreignKey:UserID;references:ID"`
	PostComments []PostComment `json:"post_comments" gorm:"foreignKey:UserID;references:ID"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
