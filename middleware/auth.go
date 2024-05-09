package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte("secret")},
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error",
			"message": "Missing or malformed JWT", "data": nil})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error",
		"message": "Invalid or expired JWT",
		"data":    nil})
}
