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
	"outstagram/models/req"
	"path/filepath"
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

func (u *UserService) AvatarUploadValidateRequest(body *multipart.Form) (*multipart.FileHeader, error) {
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

	return files[0], nil
}

func (u *UserService) AvatarUploadFile(ctx *fiber.Ctx, file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", errors.New("avatar file is required")
	}

	ext := filepath.Ext(file.Filename)
	randName := common.RandomNString(30)
	newImageName := fmt.Sprintf("%s%s", randName, ext)

	if err := ctx.SaveFile(file, newImageName); err != nil {
		return "", errors.New("error while saving file " + file.Filename)
	}

	isImage := strings.HasPrefix(file.Header.Get("Content-Type"), "image/")

	params := uploader.UploadParams{
		DisplayName: randName,
		Folder:      "avatars",
	}

	if isImage {
		params.Format = "webp"
	}

	result, err := common.CloudinaryUploadFile(newImageName, params)
	if err != nil {
		return "", errors.New("error while uploading file " + file.Filename)
	}

	if err := os.Remove(newImageName); err != nil {
		return "", errors.New("error while deleting local file " + newImageName)
	}

	return result.SecureURL, nil
}

func (u *UserService) UserEditByUserID(userID string, userRecord *req.UserMeUpdate, avatarFile *multipart.FileHeader, ctx *fiber.Ctx) error {
	var user entity.User
	if err := common.DBConn.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusBadRequest, "User not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying user")
	}

	//user.Email = userRecord.Email
	user.Username = userRecord.Username
	user.FullName = userRecord.FullName
	//user.Phone = userRecord.Phone
	user.Birthday = userRecord.Birthday
	user.Bio = userRecord.Bio
	user.Gender = userRecord.Gender

	if avatarFile != nil {
		avatarURL, err := u.AvatarUploadFile(ctx, avatarFile)
		if err != nil {
			return fmt.Errorf("error while uploading avatar: %v", err)
		}

		if user.Avatar != "" && user.Avatar != avatarURL {
			if err := common.CloudinaryDeleteFile(user.Avatar); err != nil {
				return fmt.Errorf("error while deleting old avatar file: %v", err)
			}
		}

		user.Avatar = avatarURL
	} else {
		if userRecord.Avatar != "" {
			user.Avatar = userRecord.Avatar
		}
	}

	if err := common.DBConn.Save(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error while updating user")
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
