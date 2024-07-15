package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"outstagram/common"
	"outstagram/models/entity"
	"outstagram/services"
)

type AdminController struct {
	adminService *services.AdminService
	postService  *services.PostService
	userService  *services.UserService
}

func NewAdminController(adminService *services.AdminService, postService *services.PostService, userService *services.UserService) *AdminController {
	return &AdminController{
		adminService: adminService,
		postService:  postService,
		userService:  userService,
	}
}

func (a *AdminController) AdminDeletePostByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")
	rawUserID := ctx.Params("userID")
	userID := uuid.MustParse(rawUserID)

	if err := a.postService.PostDeleteByPostID(postID, userID); err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Post deleted", nil)
}

func (a *AdminController) AdminBanUserByUserID(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")

	var userRecord entity.User
	if err := a.userService.UserGetByID(userID, &userRecord); err != nil {
		return err
	}

	userActive, err := a.userService.UserBanByUserID(userRecord.ID.String())
	if err != nil {
		return err
	}

	if userActive == true {
		return common.CreateResponse(ctx, fiber.StatusOK, "User unbanned", userActive)

	}

	return common.CreateResponse(ctx, fiber.StatusOK, "User banned", userActive)
}

func (a *AdminController) AdminBlockPostByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")

	postActive, err := a.postService.PostBlockByPostID(postID)
	if err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Post blocked", postActive)
}
