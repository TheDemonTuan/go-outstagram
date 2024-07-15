package services

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	if err := common.DBConn.Model(&entity.PostComment{}).Where("post_id = ?", postID).Order("created_at DESC").Find(postComments).Error; err != nil {
		return errors.New("error when getting post comments by post id")
	}

	return nil
}
