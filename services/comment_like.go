package services

import (
	"errors"
	"outstagram/common"
	"outstagram/models/entity"
)

type CommentLikeService struct{}

func NewCommentLikeService() *CommentLikeService {
	return &CommentLikeService{}
}

func (cl *CommentLikeService) CommentLikeGetAllByCommentID(commentID string, commentLikes interface{}) error {
	if err := common.DBConn.Model(&entity.PostLike{}).Where("comment_id = ? AND is_comment_liked = ?", commentID, true).Find(commentLikes).Error; err != nil {
		return errors.New("error when getting comment likes by comment id")
	}

	return nil
}
