package services

import (
	"errors"
	"outstagram/common"
	"outstagram/models/entity"
)

type PostCommentService struct{}

func NewPostCommentService() *PostCommentService {
	return &PostCommentService{}
}

func (pc *PostCommentService) PostCommentGetByID(postCommentID string, postComment interface{}) error {
	if err := common.DBConn.Where("id = ?", postCommentID).First(postComment).Error; err != nil {
		return errors.New("error when getting post comment by id")
	}

	return nil
}

func (pc *PostCommentService) PostCommentGetAllByPostID(postID string, postComments interface{}) error {
	if err := common.DBConn.Model(&entity.PostComment{}).Where("post_id = ?", postID).Find(postComments).Error; err != nil {
		return errors.New("error when getting post comments by post id")
	}

	return nil
}
