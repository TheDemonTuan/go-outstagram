package services

import (
	"errors"
	"outstagram/common"
	"outstagram/models/entity"
)

type PostLikeService struct{}

func NewPostLikeService() *PostLikeService {
	return &PostLikeService{}
}

func (pl *PostLikeService) PostLikeGetAllByPostID(postID string, postLikes interface{}) error {
	if err := common.DBConn.Model(&entity.PostLike{}).Where("post_id = ?", postID).Find(postLikes).Error; err != nil {
		return errors.New("error when getting post likes by post id")
	}

	return nil
}
