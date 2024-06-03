package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"outstagram/common"
	"outstagram/models/entity"
	"outstagram/services"
)

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		SuccessHandler: jwtSuccess,
		ErrorHandler:   jwtError,
	})
}

func jwtSuccess(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, ok := claims["uuid"].(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	userService := services.NewUserService()
	var userRecord entity.User
	if err := userService.UserGetByID(userID, &userRecord); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	if !userRecord.Active {
		return fiber.NewError(fiber.StatusUnauthorized, "User is not active")
	}

	c.Locals(common.UserIDLocalKey, userID)
	c.Locals(common.UserInfoLocalKey, userRecord)

	return c.Next()
}

func jwtError(c *fiber.Ctx, err error) error {
	if c.Path() == "/graphql" || c.Path() == "/graphql/playground" {
		return c.Next()
	}

	if err.Error() == "Missing or malformed JWT" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing or malformed JWT")
	}

	return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
}

func GraphqlHandler(c *fiber.Ctx) error {
	currentUserID, isOK := c.Locals(common.UserIDLocalKey).(string)

	if !isOK {
		c.Context().SetUserValue(common.UserInfoLocalKey, nil)
		return c.Next()
	}
	c.Context().SetUserValue(common.UserInfoLocalKey, currentUserID)

	return c.Next()
}
