package req

import (
	"github.com/google/uuid"
	"outstagram/models/entity"
	"time"
)

type UserMeUpdate struct {
	Username string    `json:"username" validate:"required,alphanum,min=3,max=50"`
	FullName string    `json:"full_name" validate:"required,min=3,max=100"`
	Birthday time.Time `json:"birthday" validate:"required"`
	Bio      string    `json:"bio" validate:"max=150"`
	Gender   string    `json:"gender" validate:"required,alphanum,min=4,max=6"`
}

type UserMeUpdatePhone struct {
	Phone string `json:"phone" validate:"required,numeric,min=10,max=15"`
}

type UserMeUpdateEmail struct {
	Email string `json:"email" validate:"required,email,min=5,max=100"`
}

type UserMeUpdatePassword struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

type UserReport struct {
	Reason entity.ReportReason `json:"reason" validate:"required"`
	Type   entity.ReportType   `json:"type" validate:"required"`
	Info   string              `json:"info" validate:"required"`
}

type ReportResponse struct {
	ID            uuid.UUID           `json:"id"`
	ByUserID      uuid.UUID           `json:"by_user_id"`
	Type          entity.ReportType   `json:"type"`
	Info          string              `json:"info"`
	Reason        entity.ReportReason `json:"reason"`
	Status        entity.ReportStatus `json:"status"`
	CreatedAt     time.Time           `json:"created_at" `
	UpdatedAt     time.Time           `json:"updated_at" `
	User          *entity.User        `json:"user,omitempty"`
	ReportingUser *entity.User        `json:"reporting_user,omitempty"`
	Post          *entity.Post        `json:"post,omitempty"`
	Comment       *entity.PostComment `json:"comment,omitempty"`
}
