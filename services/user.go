package services

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"outstagram/common"
	"outstagram/models/entity"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) UserGetMe(userID string, userRecord *entity.User) error {
	if err := common.DBConn.Where("id = ?", userID).Find(userRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return fiber.NewError(fiber.StatusBadRequest, "error when fetching user data")
	}
	return nil
}
