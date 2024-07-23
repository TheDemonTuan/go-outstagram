package services

import (
	"errors"
	"outstagram/common"
	"outstagram/models/entity"
)

type PostSaveService struct{}

func NewPostSaveService() *PostSaveService {
	return &PostSaveService{}
}

func (pl *PostSaveService) PostSaveGetAllByPostID(postID string, postSaves interface{}) error {
	if err := common.DBConn.Model(&entity.PostSave{}).Where("post_id = ?", postID).Find(postSaves).Error; err != nil {
		return errors.New("error when getting post saves by post id")
	}

	return nil
}
