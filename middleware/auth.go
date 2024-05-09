package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
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
	uuid, ok := claims["uuid"].(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	c.Locals("currentUserId", uuid)
	return c.Next()
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing or malformed JWT")
	}

	return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
}
