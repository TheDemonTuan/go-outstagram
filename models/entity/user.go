package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UserId    uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"userid"`
	Username  string    `json:"username" gorm:"unique;not null;size:50" `
	Password  string    `json:"password" gorm:"not null;size:255"`
	Name      string    `json:"name" gorm:"not null;size:100"`
	Email     string    `json:"email" gorm:"size:100"`
	Phone     string    `json:"phone" gorm:"size:15"`
	Avatar    string    `json:"avatar" gorm:"size:255"`
	Birthday  time.Time `json:"birthday" gorm:"not null"`
	Bio       string    `json:"bio" gorm:"size:255"`
	Pronouns  string    `json:"pronouns" gorm:"size:50"`
	Gender    bool      `json:"gender" gorm:"size:10"`
	Active    bool      `json:"active" gorm:"default:true"`
	Role      bool      `json:"role" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
