package services

import (
	"outstagram/common"
	"outstagram/models/entity"
)

type PostFileService struct{}

func NewPostFileService() *PostFileService {
	return &PostFileService{}
}

func (pf *PostFileService) PostFileGetAllByPostID(postID string, postFiles *[]entity.PostFile) error {
	if err := common.DBConn.Where("post_id = ?", postID).Find(&postFiles).Error; err != nil {
		return err
	}

	return nil
}
