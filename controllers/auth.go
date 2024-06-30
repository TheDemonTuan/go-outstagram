package controllers

import (
	"github.com/gofiber/fiber/v2"
	"outstagram/common"
	"outstagram/models/entity"
	"outstagram/models/req"
	"outstagram/services"
)

type AuthController struct {
	authService  *services.AuthService
	userService  *services.UserService
	tokenService *services.TokenService
}

func NewAuthController(authService *services.AuthService, userService *services.UserService, tokenService *services.TokenService) *AuthController {
	return &AuthController{
		authService:  authService,
		userService:  userService,
		tokenService: tokenService,
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

	accessToken, err := c.authService.GenerateAccessToken(userRecord.ID.String())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	refreshToken, err := c.authService.GenerateRefreshToken(userRecord.ID.String())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Login successfully", fiber.Map{
		"user":          userRecord,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
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

	accessToken, err := c.authService.GenerateAccessToken(newUser.ID.String())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	refreshToken, err := c.authService.GenerateRefreshToken(newUser.ID.String())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusCreated, "Register successfully", fiber.Map{
		"user":          newUser,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (c *AuthController) AuthVerify(ctx *fiber.Ctx) error {
	user, isOK := ctx.Locals(common.UserInfoLocalKey).(entity.User)
	if !isOK {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "User is verified", fiber.Map{
		"user": user,
	})
}

func (c *AuthController) AuthRefreshToken(ctx *fiber.Ctx) error {
	bodyData, err := common.RequestBodyValidator[req.AuthRefreshToken](ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userId, err := c.authService.ValidateRefreshToken(bodyData.RefreshToken, true)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var userRecord entity.User
	if err := c.userService.UserGetByID(userId, &userRecord); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	accessToken, err := c.authService.GenerateAccessToken(userId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Refresh token successfully", fiber.Map{
		"user":         userRecord,
		"access_token": accessToken,
	})
}

func (c *AuthController) AuthLogout(ctx *fiber.Ctx) error {
	bodyData, err := common.RequestBodyValidator[req.AuthLogout](ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userID, err := c.authService.ValidateRefreshToken(bodyData.RefreshToken, true)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := c.tokenService.DeleteRefreshTokenByToken(userID, bodyData.RefreshToken); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Logout successfully", nil)
}

func (c *AuthController) AuthOAuthLogin(ctx *fiber.Ctx) error {
	bodyData, err := common.RequestBodyValidator[req.AuthOAuthLogin](ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userRecord, err := c.authService.AuthOAuthLogin(bodyData)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	accessToken, err := c.authService.GenerateAccessToken(userRecord.ID.String())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	refreshToken, err := c.authService.GenerateRefreshToken(userRecord.ID.String())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "OAuth successfully", fiber.Map{
		"user":          userRecord,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (c *AuthController) AuthOAuthRegister(ctx *fiber.Ctx) error {
	bodyData, err := common.RequestBodyValidator[req.AuthOAuthRegister](ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userRecord, err := c.authService.AuthOAuthRegister(bodyData)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	accessToken, err := c.authService.GenerateAccessToken(userRecord.ID.String())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	refreshToken, err := c.authService.GenerateRefreshToken(userRecord.ID.String())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusCreated, "OAuth successfully", fiber.Map{
		"user":          userRecord,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
