package req

import (
	"outstagram/models/entity"
	"time"
)

type PostMeEdit struct {
	Caption string             `json:"caption" validate:"required,min=1,max=2200"`
	Privacy entity.PostPrivacy `json:"privacy" validate:"omitempty"`
}

type PostMeComment struct {
	Content string `json:"content" validate:"required,min=1,max=255"`
}

type PostComment struct {
	ID        string      `json:"id"`
	Content   string      `json:"content"`
	CreatedAt time.Time   `json:"created_at"`
	User      UserComment `json:"user"`
}

type UserComment struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
}
