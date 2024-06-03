package services

import (
	"errors"
	"outstagram/common"
	"outstagram/models/entity"
)

type PostFileService struct{}

func NewPostFileService() *PostFileService {
	return &PostFileService{}
}

func (pf *PostFileService) PostFileGetAllByPostID(postID string, postFiles interface{}) error {
	if err := common.DBConn.Model(&entity.PostFile{}).Where("post_id = ?", postID).Find(postFiles).Error; err != nil {
		return errors.New("error when getting post files by post id")
	}

	return nil
}
