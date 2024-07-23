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

func (cl *CommentLikeService) CommentLikeGetAllByCommentID(postID string, commentLikes interface{}) error {
	if err := common.DBConn.Model(&entity.CommentLike{}).Joins("JOIN post_comments ON post_comments.id = comment_likes.comment_id").Where("post_comments.post_id = ? AND comment_likes.is_comment_liked = ?", postID, true).Find(commentLikes).Error; err != nil {
		return errors.New("error when getting comment likes by comment id")
	}

	return nil
}
