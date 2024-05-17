package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"outstagram/common"
	"outstagram/models/entity"
	"outstagram/models/req"
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

func (p *PostController) PostMeGetAll(ctx *fiber.Ctx) error {
	rawUserID := ctx.Locals("currentUserId").(string)
	var postRecords []entity.Post
	if err := common.DBConn.Where("user_id = ?", rawUserID).Preload("PostFiles").Preload("PostLikes").Preload("PostComments").Order("created_at DESC").Find(&postRecords).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "No posts found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying posts")
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Posts found", postRecords)
}

func (p *PostController) PostMeCreate(ctx *fiber.Ctx) error {
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

func (p *PostController) PostMeLikeByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")
	rawUserID := ctx.Locals("currentUserId").(string)
	userID := uuid.MustParse(rawUserID)

	postLike, err := p.postService.PostLikeSaveToDB(postID, userID)
	if err != nil {
		return err
	}

	if !postLike.IsLiked {
		return common.CreateResponse(ctx, fiber.StatusOK, "Post unliked", postLike)
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Post liked", postLike)
}

func (p *PostController) PostMeEditByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")
	rawUserID := ctx.Locals("currentUserId").(string)
	userID := uuid.MustParse(rawUserID)

	bodyData, err := common.Validator[req.PostMeEdit](ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	post, err := p.postService.PostEditSaveToDB(postID, userID, bodyData.Caption)
	if err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Post edited", post)
}

func (p *PostController) PostMeDeleteByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")
	rawUserID := ctx.Locals("currentUserId").(string)
	userID := uuid.MustParse(rawUserID)

	if err := p.postService.PostDeleteSaveToDB(postID, userID); err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Post deleted", nil)
}

func (p *PostController) PostMeCommentByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")
	rawUserID := ctx.Locals("currentUserId").(string)
	userID := uuid.MustParse(rawUserID)

	bodyData, err := common.Validator[req.PostMeComment](ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	postComment, err := p.postService.PostCommentSaveToDB(postID, userID, bodyData.Content)
	if err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusCreated, "Comment created", postComment)
}

func (p *PostController) PostGetAllByUserID(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")

	var postRecords []entity.Post
	if err := common.DBConn.Where("user_id = ?", userID).Preload("PostFiles").Preload("PostLikes").Preload("PostComments").Find(&postRecords).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "No posts found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying posts")
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Posts found", postRecords)
}

func (p *PostController) PostGetByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")

	var post entity.Post
	if err := common.DBConn.Where("id = ?", postID).Preload("PostFiles").Preload("PostLikes").Preload("PostComments").First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Post not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying post")
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Post found", post)
}

func (p *PostController) PostGetAllCommentByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")

	resultComments, err := p.postService.GetAllCommentsByPostID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "No comments found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying comments")
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Comments found", resultComments)
}
