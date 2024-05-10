package controllers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"os"
	"outstagram/common"
	"outstagram/models/entity"
	"outstagram/services"
	"strings"
)

type PostController struct {
	postService *services.PostService
}

func NewPostController(postService *services.PostService) *PostController {
	return &PostController{
		postService: postService,
	}
}

func (p *PostController) PostGetAll(ctx *fiber.Ctx) error {
	rawUserID := ctx.Locals("currentUserId").(string)
	var postRecords []entity.Post
	if err := common.DBConn.Where("user_id = ?", rawUserID).Find(&postRecords).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "No posts found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying posts")
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Posts found", postRecords)
}

func (p *PostController) PostCreate(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}

	caption := form.Value["caption"]
	files := form.File["files"]

	acceptTypes := []string{"image/jpeg", "image/png", "image/jpg", "video/mp4"}
	for _, file := range files {
		//Validate if the file is an image or video
		isAcceptType := false
		for _, acceptType := range acceptTypes {
			if file.Header["Content-Type"][0] == acceptType {
				isAcceptType = true
				break
			}
		}
		if !isAcceptType {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid file type")
		}
	}

	var filePaths []string
	for _, file := range files {
		ext := strings.Split(file.Header["Content-Type"][0], "/")[1]
		randName := common.RandomString(25)
		newFileName := fmt.Sprintf("%s.%s", randName, ext)

		err := ctx.SaveFile(file, fmt.Sprintf("./%s/posts/%s", os.Getenv("STATIC_PATH"), newFileName))

		//Check for errors
		if err != nil {
			for _, imagePath := range filePaths {
				_ = os.Remove(fmt.Sprintf("./%s/posts/%s", os.Getenv("STATIC_PATH"), imagePath))
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		filePaths = append(filePaths, newFileName)
	}

	rawUserID := ctx.Locals("currentUserId").(string)

	userID := uuid.MustParse(rawUserID)

	newPost := entity.Post{
		ID:      common.RandomString(15),
		UserID:  userID,
		Caption: caption[0],
		Files:   strings.Join(filePaths, ","), //Convert slice to string
	}

	if err := common.DBConn.Create(&newPost).Error; err != nil {
		for _, imagePath := range filePaths {
			_ = os.Remove(fmt.Sprintf("./%s/posts/%s", os.Getenv("STATIC_PATH"), imagePath))
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while creating post")
	}

	return common.CreateResponse(ctx, fiber.StatusCreated, "Post created", newPost)
}
