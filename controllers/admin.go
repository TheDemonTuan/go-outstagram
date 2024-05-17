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
	rawUserID := ctx.Locals(common.UserIDLocalKey).(string)
	userID := uuid.MustParse(rawUserID)

	if err := a.postService.PostDeleteByPostID(postID, userID); err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Post deleted", nil)
}

func (a *AdminController) AdminBanUserByUserID(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")

	var userRecord entity.User
	if err := a.userService.UserGetByUserID(userID, &userRecord); err != nil {
		return err
	}

	if err := a.userService.UserBanByUserID(userRecord.ID.String()); err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "User banned", nil)
}

func (a *AdminController) AdminUnbanUserByUserID(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")

	var userRecord entity.User
	if err := a.userService.UserGetByUserID(userID, &userRecord); err != nil {
		return err
	}

	if err := a.userService.UserUnbanByUserID(userRecord.ID.String()); err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "User unbanned", nil)
}

//func (a *AdminController) AdminBlockPostByPostID(ctx *fiber.Ctx) error {
//	postID := ctx.Params("postID")
//
//	if err := a.postService.PostBlockByPostID(postID); err != nil {
//		return err
//	}
//
//	return common.CreateResponse(ctx, fiber.StatusOK, "Post blocked", nil)
//}
//
//func (a *AdminController) AdminUnblockPostByPostID(ctx *fiber.Ctx) error {
//	postID := ctx.Params("postID")
//
//	if err := a.postService.PostUnblockByPostID(postID); err != nil {
//		return err
//	}
//
//	return common.CreateResponse(ctx, fiber.StatusOK, "Post unblocked", nil)
//}
