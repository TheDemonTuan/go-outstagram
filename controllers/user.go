package controllers

import (
	"github.com/gofiber/fiber/v2"
	"outstagram/common"
	"outstagram/models/entity"
	"outstagram/models/req"
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
	var userRecord entity.User

	userRecord = ctx.Locals(common.UserInfoLocalKey).(entity.User)

	return common.CreateResponse(ctx, fiber.StatusOK, "Get user successfully", userRecord)
}

func (u *UserController) UserGetByUserID(ctx *fiber.Ctx) error {
	rawUserID := ctx.Params("userID")
	var userRecord entity.User

	if err := u.userService.UserGetByID(rawUserID, &userRecord); err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Get user successfully", userRecord)
}

func (u *UserController) UserMeUploadAvatar(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	fileHeader, err := u.userService.UserMeUploadAvatarValidateRequest(form)
	if err != nil {
		return err
	}

	uploadResult, err := u.userService.UserMeUploadAvatar(fileHeader, ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	userID := ctx.Locals(common.UserIDLocalKey).(string)
	if err := u.userService.UserMeSaveAvatarToDB(userID, uploadResult.SecureURL); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Avatar uploaded successfully", uploadResult.SecureURL)
}

func (u *UserController) UserMeEditProfile(ctx *fiber.Ctx) error {
	bodyData, err := common.RequestBodyValidator[req.UserMeUpdate](ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := u.userService.UserMeEditProfileValidateRequest(bodyData); err != nil {
		return err
	}

	user, err := u.userService.UserMeEditProfileSaveToDB(ctx, bodyData)
	if err != nil {
		return err
	}
	return common.CreateResponse(ctx, fiber.StatusOK, "Profile updated", user)

}

func (u *UserController) UserMeEditPrivate(ctx *fiber.Ctx) error {

	if err := u.userService.UserMeEditPrivateSaveToDB(ctx); err != nil {
		return err
	}

	userInfo := ctx.Locals(common.UserInfoLocalKey).(entity.User)

	return common.CreateResponse(ctx, fiber.StatusOK, "Private updated", userInfo.IsPrivate)

}

