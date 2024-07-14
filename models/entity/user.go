package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type UserOAuth int

const (
	OAuthDefault UserOAuth = iota
	OAuthFacebook
	OAuthGoogle
	OAuthGithub
)

func (u UserOAuth) EnumIndex() int {
	return int(u)
}

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username  string    `json:"username" gorm:"unique;index;not null;size:50" `
	Password  string    `json:"-" gorm:"not null;size:255"`
	FullName  string    `json:"full_name" gorm:"index;not null;size:100"`
	Email     string    `json:"email" gorm:"unique;index;not null;size:100"`
	Phone     string    `json:"phone" gorm:"unique;index;size:15"`
	Avatar    string    `json:"avatar" gorm:"size:255"`
	Bio       string    `json:"bio" gorm:"size:255"`
	Birthday  time.Time `json:"birthday" gorm:"not null"`
	Gender    bool      `json:"gender" gorm:"default:false"`
	Role      bool      `json:"role" gorm:"default:false"`
	OAuth     UserOAuth `json:"oauth" gorm:"column:oauth;default:0"`
	Active    bool      `json:"active" gorm:"default:true"`
	IsPrivate bool      `json:"is_private" gorm:"default:false"`

	Posts        []Post        `json:"-" gorm:"foreignKey:UserID;references:ID"`
	PostLikes    []PostLike    `json:"-" gorm:"foreignKey:UserID;references:ID"`
	PostComments []PostComment `json:"-" gorm:"foreignKey:UserID;references:ID"`
	FromFriends  []Friend      `json:"-" gorm:"foreignKey:FromUserID;references:ID"`
	ToFriends    []Friend      `json:"-" gorm:"foreignKey:ToUserID;references:ID"`
	InboxFrom    []Inbox       `json:"-" gorm:"foreignKey:FromUserID;references:ID"`
	InboxTo      []Inbox       `json:"-" gorm:"foreignKey:ToUserID;references:ID"`
	Token        []Token       `json:"-" gorm:"foreignKey:UserID;references:ID"`
	PostSaves    []PostSave    `json:"-" gorm:"foreignKey:UserID;references:ID"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
