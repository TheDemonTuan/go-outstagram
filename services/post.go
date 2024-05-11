package services

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func (p *PostService) GetStaticPath() string {
	return os.Getenv("STATIC_PATH") + "/posts/"
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
		"video/mp4":  true,
	}

	for _, file := range files {
		if !acceptType[file.Header["Content-Type"][0]] {
			return "", nil, errors.New(file.Filename + "file type is not supported")
		}
	}

	return caption[0], files, nil
}

func (p *PostService) PostCreateUploadFiles(ctx *fiber.Ctx, files []*multipart.FileHeader) ([]string, error) {
	var filePaths []string
	for _, file := range files {
		ext := strings.Split(file.Header["Content-Type"][0], "/")[1]
		randName := common.RandomNString(30)
		newImageName := fmt.Sprintf("%s.%s", randName, ext)

		//Check for errors
		if err := ctx.SaveFile(file, p.GetStaticPath()+newImageName); err != nil {
			if err := p.PostCreateDeleteFiles(filePaths); err != nil {
				return nil, err
			}
			return nil, errors.New("error while saving file " + file.Filename)
		}

		filePaths = append(filePaths, newImageName)
	}

	return filePaths, nil
}

func (p *PostService) PostCreateDeleteFiles(filePaths []string) error {
	for _, imagePath := range filePaths {
		if err := os.Remove(p.GetStaticPath() + imagePath); err != nil {
			return errors.New("error while deleting file " + imagePath)
		}
	}

	return nil
}

func (p *PostService) PostCreateSaveToDB(userID uuid.UUID, caption string, filePaths []string) (entity.Post, error) {
	newPost := entity.Post{
		ID:      "TD" + common.RandomNString(18),
		UserID:  userID,
		Caption: caption,
		Files:   p.GetStaticPath() + "," + strings.Join(filePaths, ","),
	}

	if err := common.DBConn.Create(&newPost).Error; err != nil {
		if err := p.PostCreateDeleteFiles(filePaths); err != nil {
			return entity.Post{}, err
		}
		return entity.Post{}, errors.New("error while creating post")
	}

	return newPost, nil
}
