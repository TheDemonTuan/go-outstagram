package controllers

import (
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"outstagram/common"
	"outstagram/models/entity"
	"outstagram/models/req"
	"outstagram/services"
	"time"
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
	rawUserID := ctx.Locals(common.UserIDLocalKey).(string)
	var userRecord entity.User

	if err := u.userService.UserGetByUserID(rawUserID, &userRecord); err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Get user successfully", userRecord)
}

func (u *UserController) UserEditMe(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var avatarFile *multipart.FileHeader
	if form != nil {
		avatarFile, err = u.userService.AvatarUploadValidateRequest(form)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	rawUserID := ctx.Locals(common.UserIDLocalKey).(string)

	username := form.Value["username"][0]
	fullName := form.Value["full_name"][0]
	birthdayStr := form.Value["birthday"][0]
	bio := form.Value["bio"][0]
	genderStr := form.Value["gender"][0]

	var avatar string
	if len(form.Value["avatar"]) > 0 {
		avatar = form.Value["avatar"][0]
	}

	birthday, err := time.Parse(time.RFC3339, birthdayStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid birthday format")
	}

	gender := genderStr == "true"

	userRecord := req.UserMeUpdate{
		Username: username,
		FullName: fullName,
		Birthday: birthday,
		Avatar:   avatar,
		Bio:      bio,
		Gender:   gender,
	}
	//var userRecord req.UserMeUpdate
	//if err := ctx.BodyParser(&userRecord); err != nil {
	//	return fiber.NewError(fiber.StatusBadRequest, err.Error())
	//}

	if err := u.userService.UserEditByUserID(rawUserID, &userRecord, avatarFile, ctx); err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Update user successfully", userRecord)
}

func (u *UserController) UserGetByUserID(ctx *fiber.Ctx) error {
	rawUserID := ctx.Params("userID")
	var userRecord entity.User

	if err := u.userService.UserGetByUserID(rawUserID, &userRecord); err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Get user successfully", userRecord)
}
