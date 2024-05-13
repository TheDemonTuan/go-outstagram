package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"outstagram/common"
	"outstagram/models/entity"
	"outstagram/services"
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
	if err := common.DBConn.Where("user_id = ?", rawUserID).Preload("PostFiles").Find(&postRecords).Error; err != nil {
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

	caption, files, err := p.postService.PostCreateValidateRequest(form)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	localPaths, cloudinaryPaths, err := p.postService.PostCreateUploadFiles(ctx, files)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	rawUserID := ctx.Locals("currentUserId").(string)
	userID := uuid.MustParse(rawUserID)

	newPost, err := p.postService.PostCreateSaveToDB(userID, caption, localPaths, cloudinaryPaths)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusCreated, "Post created", newPost)
}
