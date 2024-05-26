package common

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func RequestBodyValidator[T any](c *fiber.Ctx) (*T, error) {
	body := new(T)

	if err := c.BodyParser(body); err != nil {
		return nil, errors.New("invalid request")
	}

	var validate = validator.New(validator.WithRequiredStructEnabled())
	if errs := validate.Struct(body); errs != nil {
		err := errs.(validator.ValidationErrors)[0]
		switch err.Tag() {
		case "required":
			return nil, errors.New(fmt.Sprintf("%s không được để trống", err.Field()))
		case "email":
			return nil, errors.New(fmt.Sprintf("%s không phải định dạng email ", err.Field()))
		case "len":
			return nil, errors.New(fmt.Sprintf("%s phải dài chính xác %v ký tự", err.Field(), err.Param()))
		case "min":
			return nil, errors.New(fmt.Sprintf("%s phải có ít nhất %v ký tự", err.Field(), err.Param()))
		case "max":
			return nil, errors.New(fmt.Sprintf("%s không được dài hơn %v ký tự", err.Field(), err.Param()))
		case "alphanum":
			return nil, errors.New(fmt.Sprintf("%s chỉ được phép chứa ký tự hoặc là số", err.Field()))
		case "alphanumunicode":
			return nil, errors.New(fmt.Sprintf("%s chỉ được phép chứa ký tự hoặc là số", err.Field()))
		case "alphaunicode":
			return nil, errors.New(fmt.Sprintf("%s chỉ được phép là ký tự", err.Field()))
		case "number":
			return nil, errors.New(fmt.Sprintf("%s chỉ được phép là số", err.Field()))
		case "gte":
			return nil, errors.New(fmt.Sprintf("%s phải nhiều hơn hoặc bằng %v", err.Field(), err.Param()))
		case "lte":
			return nil, errors.New(fmt.Sprintf("%s phải ít hơn hoặc bằng %v", err.Field(), err.Param()))
		default:
			return nil, errors.New(fmt.Sprintf("%s: %v must satisfy %s %v criteria", err.Field(), err.Value(), err.Tag(), err.Param()))
		}
	}

	return body, nil
}
