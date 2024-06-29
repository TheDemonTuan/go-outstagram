package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type UserOAuth int

const (
	Default UserOAuth = iota
	Facebook
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null;size:50" `
	Password  string    `json:"-" gorm:"not null;size:255"`
	FullName  string    `json:"full_name" gorm:"index;not null;size:100"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null;size:100"`
	Phone     string    `json:"phone" gorm:"uniqueIndex;size:15"`
	Avatar    string    `json:"avatar" gorm:"size:255"`
	Bio       string    `json:"bio" gorm:"size:255"`
	Birthday  time.Time `json:"birthday" gorm:"not null"`
	Gender    bool      `json:"gender" gorm:"default:false"`
	Role      bool      `json:"role" gorm:"default:false"`
	OAuth     UserOAuth `json:"oauth" gorm:"default:0"`
	Active    bool      `json:"active" gorm:"default:true"`
	IsPrivate bool      `json:"is_private" gorm:"default:false"`

	Posts        []Post        `json:"posts" gorm:"foreignKey:UserID;references:ID"`
	PostLikes    []PostLike    `json:"post_likes" gorm:"foreignKey:UserID;references:ID"`
	PostComments []PostComment `json:"post_comments" gorm:"foreignKey:UserID;references:ID"`
	FromFriends  []Friend      `json:"from_friends" gorm:"foreignKey:FromUserID;references:ID"`
	ToFriends    []Friend      `json:"to_friends" gorm:"foreignKey:ToUserID;references:ID"`
	InboxFrom    []Inbox       `json:"inbox_from" gorm:"foreignKey:FromUserID;references:ID"`
	InboxTo      []Inbox       `json:"inbox_to" gorm:"foreignKey:ToUserID;references:ID"`
	Token        []Token       `json:"tokens" gorm:"foreignKey:UserID;references:ID"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
