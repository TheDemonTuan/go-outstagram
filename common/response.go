package common

import "github.com/gofiber/fiber/v2"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func CreateResponse(c *fiber.Ctx, code int, message string, data interface{}) error {
	return c.Status(code).JSON(Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
