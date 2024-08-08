package req

import (
	"github.com/google/uuid"
	"outstagram/models/entity"
)

type PostMeEdit struct {
	Caption string             `json:"caption" validate:"required,min=1,max=2200"`
	Privacy entity.PostPrivacy `json:"privacy" validate:"omitempty"`
}

type PostMeComment struct {
	Content string `json:"content" validate:"required,min=1,max=255"`
}

type PostMeRestore struct {
	PostIDs []string `json:"post_ids"`
}

type PostMeDelete struct {
	PostIDs []string `json:"post_ids"`
}

type PostResponse struct {
	Post entity.Post `json:"post"`
	User entity.User `json:"user"`
}

type DeleteCommentsRequest struct {
	CommentIDs []struct {
		PostID    string    `json:"post_id"`
		CommentID uuid.UUID `json:"comment_id"`
	} `json:"comment_ids"`
}

//type DeleteCommentsRequest struct {
//	CommentIDs []uuid.UUID `json:"comment_ids"`
//}
