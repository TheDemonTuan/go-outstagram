package controllers

import (
	"github.com/gofiber/fiber/v2"
	"outstagram/common"
	"outstagram/models/entity"
	"outstagram/services"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (u *UserController) UserGetMe(ctx *fiber.Ctx) error {
	rawUserID := ctx.Locals("currentUserId").(string)
	var userRecord entity.User

	if err := u.userService.UserGetMe(rawUserID, &userRecord); err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Get user successfully", userRecord)
}
