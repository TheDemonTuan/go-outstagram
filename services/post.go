package services

import (
	"errors"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
	"os"
	"outstagram/common"
	"outstagram/models/entity"
	"strings"
)

type PostService struct{}

func NewPostService() *PostService {
	return &PostService{}
}

func (p *PostService) PostCreateValidateRequest(body *multipart.Form) (string, []*multipart.FileHeader, error) {
	if body == nil {
		return "", nil, errors.New("request body is required")
	}

	if body.Value == nil {
		return "", nil, errors.New("request body value is required")
	}

	if body.File == nil {
		return "", nil, errors.New("request body file is required")
	}

	caption := body.Value["caption"]
	files := body.File["files"]

	if len(caption) == 0 || len(caption[0]) == 0 {
		return "", nil, errors.New("caption is required")
	}

	if len(caption[0]) > 2200 {
		return "", nil, errors.New("caption is too long")
	}

	if len(files) == 0 {
		return "", nil, errors.New("files are required")
	}

	if len(files) > 10 {
		return "", nil, errors.New("files are too many")
	}

	acceptType := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/jpg":  true,
		"image/webp": true,
		"video/mp4":  true,
	}

	for _, file := range files {
		if !acceptType[file.Header["Content-Type"][0]] {
			return "", nil, errors.New(file.Filename + "file type is not supported")
		}
	}

	return caption[0], files, nil
}

func (p *PostService) PostCreateUploadFiles(ctx *fiber.Ctx, files []*multipart.FileHeader) ([]string, []*uploader.UploadResult, error) {
	var localPaths []string
	var cloudinaryPaths []*uploader.UploadResult
	for _, file := range files {
		ext := strings.Split(file.Header["Content-Type"][0], "/")[1]
		randName := common.RandomNString(30)
		newImageName := fmt.Sprintf("%s.%s", randName, ext)

		//Save to local
		if err := ctx.SaveFile(file, newImageName); err != nil {
			if err := p.PostCreateDeleteFiles(localPaths, cloudinaryPaths, true); err != nil {
				return nil, nil, err
			}
			return nil, nil, errors.New("error while saving file " + file.Filename)
		}

		isImage := file.Header["Content-Type"][0][:5] == "image"

		params := uploader.UploadParams{
			DisplayName: randName,
			Folder:      "posts",
		}

		if isImage {
			params = uploader.UploadParams{
				Format:      "webp",
				DisplayName: randName,
				Folder:      "posts",
			}
		}

		//Upload to cloudinary
		result, err := common.CloudinaryUploadFile(newImageName, params)
		if err != nil {
			if err := p.PostCreateDeleteFiles(localPaths, cloudinaryPaths, true); err != nil {
				return nil, nil, err
			}
			return nil, nil, errors.New("error while uploading file " + file.Filename)
		}

		cloudinaryPaths = append(cloudinaryPaths, result)
		localPaths = append(localPaths, newImageName)
	}

	if err := p.PostCreateDeleteFiles(localPaths, cloudinaryPaths, false); err != nil {
		return nil, nil, err
	}

	return localPaths, cloudinaryPaths, nil
}

func (p *PostService) PostCreateDeleteFiles(localPaths []string, cloudinaryPaths []*uploader.UploadResult, isDeleteCloud bool) error {
	isErr := false
	for _, filePath := range localPaths {
		if err := os.Remove(filePath); err != nil {
			isErr = true
		}
	}

	if isDeleteCloud || isErr {
		for _, cloudinaryPath := range cloudinaryPaths {
			if err := common.CloudinaryDeleteFile(cloudinaryPath.PublicID); err != nil {
				return err
			}
		}
	}

	if isErr {
		return errors.New("error while deleting files")
	}

	return nil
}

func (p *PostService) PostCreateSaveToDB(userID uuid.UUID, caption string, localPaths []string, cloudinaryPaths []*uploader.UploadResult) (entity.Post, error) {
	newPostID := "TD" + common.RandomNString(18)
	newPost := entity.Post{
		ID:      newPostID,
		UserID:  userID,
		Caption: caption,
	}

	for _, filePath := range cloudinaryPaths {
		newPost.PostFiles = append(newPost.PostFiles, entity.PostFile{
			ID:     filePath.PublicID,
			PostID: newPostID,
			URL:    filePath.SecureURL,
		})
	}

	if err := common.DBConn.Create(&newPost).Error; err != nil {
		if err := p.PostCreateDeleteFiles(localPaths, cloudinaryPaths, true); err != nil {
			return entity.Post{}, err
		}
		return entity.Post{}, errors.New("error while creating post")
	}

	return newPost, nil
}

func (p *PostService) PostLikeSaveToDB(postID string, userID uuid.UUID) error {
	var post entity.Post
	if err := common.DBConn.Where("id = ?", postID).First(&post).Error; err != nil {
		return errors.New("post not found")
	}

	var postLike entity.PostLike
	if err := common.DBConn.Where("post_id = ? AND user_id = ?", postID, userID).First(&postLike).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			postLike = entity.PostLike{
				PostID: postID,
				UserID: userID,
			}

			if err := common.DBConn.Create(&postLike).Error; err != nil {
				return errors.New("error while creating like")
			}
			return nil
		}
		return errors.New("error while querying like")
	}

	if err := common.DBConn.Delete(&postLike).Error; err != nil {
		return errors.New("error while deleting like")
	}

	return nil
}
