package services

import (
	"errors"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"mime/multipart"
	"os"
	"outstagram/common"
	"outstagram/models/entity"
	"strings"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) UserGetByUserID(userID string, userRecord *entity.User) error {
	if err := common.DBConn.Where("id = ?", userID).Find(userRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("error while querying user")
	}
	return nil

}

func (u *UserService) UserSearchByUsernameOrFullName(keyword string, users *[]entity.User) error {
	if err := common.DBConn.Where("username LIKE ? OR full_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("error while querying user")
	}

	return nil
}

func (u *UserService) UserGetByUserName(username string, userRecord *entity.User) error {
	if err := common.DBConn.Where("username = ?", username).Find(userRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("error while querying user")
	}
	return nil
}

func (u *UserService) UserMeUploadAvatarValidateRequest(body *multipart.Form) (*multipart.FileHeader, error) {
	if body == nil {
		return nil, errors.New("request body is required")
	}

	if body.File == nil {
		return nil, errors.New("request body file is required")
	}

	files := body.File["avatar"]

	if len(files) == 0 {
		return nil, errors.New("avatar file is required")
	}

	if len(files) > 1 {
		return nil, errors.New("only one avatar file is allowed")
	}

	acceptType := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/jpg":  true,
		"image/webp": true,
	}

	for _, file := range files {
		if !acceptType[file.Header["Content-Type"][0]] {
			return nil, errors.New(file.Filename + "file type is not supported")
		}
	}

	return files[0], nil
}

func (u *UserService) UserMeUploadAvatar(file *multipart.FileHeader, ctx *fiber.Ctx) (*uploader.UploadResult, error) {
	ext := strings.Split(file.Header["Content-Type"][0], "/")[1]
	name := ctx.Locals(common.UserIDLocalKey).(string)
	newAvatarName := fmt.Sprintf("%s.%s", name, ext)

	//Save to local
	if err := ctx.SaveFile(file, newAvatarName); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "error while saving file")
	}

	result, err := common.CloudinaryUploadFile(newAvatarName, uploader.UploadParams{
		Folder:      "users/avatar",
		DisplayName: name,
		Format:      "webp",
	})

	if err != nil {
		if err := os.Remove(newAvatarName); err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, "error while deleting file")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, "error while uploading file")
	}

	if err := os.Remove(newAvatarName); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "error while deleting file")
	}

	return result, nil
}

func (u *UserService) UserMeDeleteAvatar(publicID string) error {
	if err := common.CloudinaryDeleteFile(publicID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error while deleting file")
	}

	return nil
}

func (u *UserService) UserMeSaveAvatarToDB(userID string, secureURL string) error {
	var user entity.User
	if err := common.DBConn.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusBadRequest, "User not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying user")
	}

	oldAvatar := user.Avatar
	user.Avatar = secureURL

	if err := common.DBConn.Save(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error while updating user")
	}

	//Delete old avatar
	if oldAvatar != "" {
		publicID, err := common.GetPublicIDFromURL("users/avatar", oldAvatar)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "error while getting public ID")
		}
		if err := u.UserMeDeleteAvatar(publicID); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "error while deleting file")
		}
	}

	return nil
}

func (u *UserService) UserBanByUserID(userID string) error {
	var user entity.User
	if err := common.DBConn.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusBadRequest, "User not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying user")
	}

	user.Active = false

	if err := common.DBConn.Save(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error while updating user")
	}

	return nil
}

func (u *UserService) UserUnbanByUserID(userID string) error {
	var user entity.User
	if err := common.DBConn.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusBadRequest, "User not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying user")
	}

	user.Active = true

	if err := common.DBConn.Save(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error while updating user")
	}

	return nil
}
