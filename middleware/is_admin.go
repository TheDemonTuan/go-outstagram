package middleware

import (
	"github.com/gofiber/fiber/v2"
	"outstagram/common"
	"outstagram/models/entity"
)

func IsAdmin(ctx *fiber.Ctx) error {
	userInfo := ctx.Locals(common.UserInfoLocalKey).(entity.User)

	if !userInfo.Role {
		return fiber.NewError(fiber.StatusForbidden, "You are not an admin")
	}

	return ctx.Next()
}
