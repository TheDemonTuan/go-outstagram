package req

import (
	"outstagram/models/entity"
)

type PostMeEdit struct {
	Caption string             `json:"caption" validate:"required,min=1,max=2200"`
	Privacy entity.PostPrivacy `json:"privacy" validate:"omitempty"`
}

type PostMeComment struct {
	Content string `json:"content" validate:"required,min=1,max=255"`
}
