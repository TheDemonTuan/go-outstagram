package controllers

import (
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
	bodyData, err := common.Validator[req.AuthLogin](ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userEntity, err := c.authService.AuthenticateUser(bodyData.Username, bodyData.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	token, err := c.authService.CreateJWT(userEntity.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ctx.Set("TDT-Auth-Token", token)

	return common.CreateResponse(ctx, fiber.StatusOK, "Login successfully", fiber.Map{
		"token": token,
	})
}

func (c *AuthController) AuthRegister(ctx *fiber.Ctx) error {
	bodyData, err := common.Validator[req.AuthRegister](ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	newUser, err := c.authService.CreateUser(bodyData)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	token, err := c.authService.CreateJWT(newUser.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ctx.Set("TDT-Auth-Token", token)

	return common.CreateResponse(ctx, fiber.StatusCreated, "Register successfully", fiber.Map{
		"token": token,
	})
}

func (c *AuthController) AuthVerify(ctx *fiber.Ctx) error {
	currentUserID, currenUserIdIsOk := ctx.Locals("currentUserId").(string)
	if !currenUserIdIsOk {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	user, err := c.authService.VerifyUser(currentUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "User is verified", fiber.Map{
		"user": user,
	})
}
