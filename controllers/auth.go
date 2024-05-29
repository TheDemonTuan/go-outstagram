package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"outstagram/common"
	"outstagram/models/req"
	"outstagram/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) AuthLogin(ctx *fiber.Ctx) error {
	bodyData, err := common.RequestBodyValidator[req.AuthLogin](ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userRecord, err := c.authService.AuthenticateUser(bodyData.Username, bodyData.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	token, err := c.authService.CreateJWT(userRecord.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Login successfully", fiber.Map{
		"user":  userRecord,
		"token": token,
	})
}

func (c *AuthController) AuthRegister(ctx *fiber.Ctx) error {
	bodyData, err := common.RequestBodyValidator[req.AuthRegister](ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !c.authService.ValidateFullName(bodyData.FullName) {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid full name")
	}

	newUser, err := c.authService.CreateUser(bodyData)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	token, err := c.authService.CreateJWT(newUser.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusCreated, "Register successfully", fiber.Map{
		"user":  newUser,
		"token": token,
	})
}

func (c *AuthController) AuthVerify(ctx *fiber.Ctx) error {
	currentUserID, currenUserIdIsOk := ctx.Locals(common.UserIDLocalKey).(string)
	if !currenUserIdIsOk {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	user, err := c.authService.VerifyUser(currentUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	data := map[string]string{"message": "hello " + user.FullName}

	if err := common.PusherClient.Trigger("my-channel", "my-event", data); err != nil {
		fmt.Println(err.Error())
	}
	return common.CreateResponse(ctx, fiber.StatusOK, "User is verified", fiber.Map{
		"user": user,
	})
}
