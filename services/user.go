package services

import (
	"errors"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"mime/multipart"
	"os"
	"outstagram/common"
	"outstagram/models/entity"
	"outstagram/models/req"
	"regexp"
	"strings"
	"time"
)

type UserService struct {
	friendService *FriendService
}

func NewUserService() *UserService {
	return &UserService{
		friendService: NewFriendService(),
	}
}

func (u *UserService) UserGetByID(userID string, userRecord interface{}) error {
	if err := common.DBConn.Model(&entity.User{}).Where("id = ?", userID).Find(userRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("error while querying user")
	}
	return nil

}

func (u *UserService) UserSearchByUsernameOrFullName(keyword string, users interface{}) error {
	if err := common.DBConn.Model(&entity.User{}).Where("username LIKE ? OR full_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("error while querying user")
	}

	return nil
}

func (u *UserService) UserGetByUserName(username string, userRecord interface{}) error {
	if err := common.DBConn.Model(&entity.User{}).Where("username = ?", username).Find(userRecord).Error; err != nil {
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

	if err := common.DBConn.Omit("phone").Save(&user).Error; err != nil {
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

func (u *UserService) UserBanByUserID(userID string) (bool, error) {
	var user entity.User
	if err := common.DBConn.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fiber.NewError(fiber.StatusBadRequest, "User not found")
		}
		return false, fiber.NewError(fiber.StatusInternalServerError, "Error while querying user")
	}

	user.Active = !user.Active

	if err := common.DBConn.Omit("phone").Save(&user).Error; err != nil {
		return false, fiber.NewError(fiber.StatusInternalServerError, "error while updating user")
	}

	return user.Active, nil
}

func (u *UserService) UserMeEditProfileValidateRequest(userRecord *req.UserMeUpdate) error {
	match, err := regexp.Match(`^[a-zA-ZÀ-ỹ\s]+$`, []byte(userRecord.FullName))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid full name")
	}

	if !match {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid full name")
	}

	if userRecord.Gender != "male" && userRecord.Gender != "female" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid gender")
	}

	return nil
}

func (u *UserService) IsUserAtLeast13(birthday time.Time) bool {
	thirteenYearsAgo := time.Now().AddDate(-13, 0, 0)
	return birthday.Before(thirteenYearsAgo) || birthday.Equal(thirteenYearsAgo)
}

func (u *UserService) UserMeEditProfileSaveToDB(ctx *fiber.Ctx, userRecord *req.UserMeUpdate) (entity.User, error) {
	userInfo := ctx.Locals(common.UserInfoLocalKey).(entity.User)

	genderConvert := true
	if userRecord.Gender == "male" {
		genderConvert = false
	}

	if !u.IsUserAtLeast13(userRecord.Birthday) {
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, "User must be at least 13 years old")
	}

	if userInfo.Username == userRecord.Username && userInfo.FullName == userRecord.FullName && userInfo.Birthday.Format("2006-01-02") == userRecord.Birthday.Format("2006-01-02") && userInfo.Bio == userRecord.Bio && userInfo.Gender == genderConvert {
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, "No change")
	}

	var existingUser entity.User
	if err := common.DBConn.Where("username = ?", userRecord.Username).First(&existingUser).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, fiber.NewError(fiber.StatusInternalServerError, "Error while querying username")
		}
	}

	if existingUser.ID != userInfo.ID {
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, "Username already exists")
	}

	userInfo.Username = userRecord.Username
	userInfo.FullName = userRecord.FullName
	userInfo.Birthday = userRecord.Birthday
	userInfo.Bio = userRecord.Bio
	userInfo.Gender = genderConvert

	if err := common.DBConn.Omit("phone").Save(&userInfo).Error; err != nil {
		return entity.User{}, fiber.NewError(fiber.StatusInternalServerError, "error while updating user")
	}

	return existingUser, nil
}

func (u *UserService) UserSuggestion(userID string, count int, users interface{}) error {
	friends, err := u.friendService.GetListFriendID(userID)
	if err != nil {
		return err
	}

	if len(friends) == 0 {
		friends = append(friends, userID)
	}

	if err := common.DBConn.Model(&entity.User{}).Not("id = ?", userID).Where("id NOT IN ?", friends).Where("active = ?", true).Order("RANDOM()").Limit(count).Find(users).Error; err != nil {
		return errors.New("error while querying user")
	}
	return nil
}

func (u *UserService) UserMeEditPrivateSaveToDB(ctx *fiber.Ctx) error {
	userInfo := ctx.Locals(common.UserInfoLocalKey).(entity.User)

	userInfo.IsPrivate = !userInfo.IsPrivate

	if err := common.DBConn.Omit("phone").Save(&userInfo).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error while updating user")
	}

	ctx.Locals(common.UserInfoLocalKey, userInfo)

	return nil
}

func (u *UserService) UserMeEditPhoneSaveToDB(ctx *fiber.Ctx, phone string) error {
	userInfo := ctx.Locals(common.UserInfoLocalKey).(entity.User)

	if userInfo.Phone == phone {
		return fiber.NewError(fiber.StatusBadRequest, "New phone is the same as the current phone")
	}

	var existingUser entity.User

	if err := common.DBConn.Where("phone = ? and id != ?", phone, userInfo.ID).Select("phone").First(&existingUser).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusInternalServerError, "Error while querying phone")
		}
	}

	if existingUser.Phone != "" {
		return fiber.NewError(fiber.StatusBadRequest, "Phone number already exists")
	}

	userInfo.Phone = phone

	if err := common.DBConn.Save(&userInfo).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error while updating user")
	}

	ctx.Locals(common.UserInfoLocalKey, userInfo)

	return nil
}

func (u *UserService) UserMeEditEmailSaveToDB(ctx *fiber.Ctx, email string) error {
	userInfo := ctx.Locals(common.UserInfoLocalKey).(entity.User)

	if userInfo.Email == email {
		return fiber.NewError(fiber.StatusBadRequest, "New email is the same as the current email")
	}

	var existingUser entity.User
	if err := common.DBConn.Where("email = ? and id != ?", email, userInfo.ID).First(&existingUser).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying email")
	}

	if existingUser.Email != "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email already exists")
	}

	userInfo.Email = email

	if err := common.DBConn.Save(&userInfo).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error while updating user")
	}

	ctx.Locals(common.UserInfoLocalKey, userInfo)

	return nil
}

func (u *UserService) UserMeDeleteAvatarSaveToDB(ctx *fiber.Ctx) error {
	userInfo := ctx.Locals(common.UserInfoLocalKey).(entity.User)

	if userInfo.Avatar == "" {
		return fiber.NewError(fiber.StatusBadRequest, "No avatar to delete")
	}

	publicID, err := common.GetPublicIDFromURL("users/avatar", userInfo.Avatar)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error while getting public ID")
	}

	if err := u.UserMeDeleteAvatar(publicID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error while deleting file")
	}

	userInfo.Avatar = ""

	if err := common.DBConn.Save(&userInfo).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error while updating user")
	}

	ctx.Locals(common.UserInfoLocalKey, userInfo)

	return nil
}

func (u *UserService) UserEditPasswordSaveToDB(ctx *fiber.Ctx, bodyData *req.UserMeUpdatePassword) error {
	userInfo := ctx.Locals(common.UserInfoLocalKey).(entity.User)

	if bodyData.CurrentPassword == bodyData.NewPassword {
		return fiber.NewError(fiber.StatusBadRequest, "new password cannot be the same as the current password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(bodyData.CurrentPassword))
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "current password is incorrect")
	}

	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(bodyData.NewPassword), bcrypt.DefaultCost)
	if hashedPasswordErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while hashing password")
	}

	userInfo.Password = string(hashedPassword)

	if err := common.DBConn.Omit("phone").Save(&userInfo).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error while updating user")
	}

	ctx.Locals(common.UserInfoLocalKey, userInfo)

	return nil

}
