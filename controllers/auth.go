package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"outstagram/common"
	"outstagram/models/entity"
	"outstagram/models/req"
	"time"
)

func createJWT(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"uuid": userId.String(),
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, tokenSignedErr := token.SignedString([]byte("secret"))
	if tokenSignedErr != nil {
		return "", errors.New("error while signing token")
	}

	return tokenSigned, nil
}

func AuthLogin(c *fiber.Ctx) error {
	bodyData, err := common.Validator[req.AuthLogin](c)

	if err != nil || bodyData == nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var userEntity entity.User

	db := common.DBConn.Where(" username = ?", bodyData.Username)

	if err := db.First(&userEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid credentials")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(bodyData.Password)); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid password")
	}

	token, tokenIsErr := createJWT(userEntity.UserId)
	if tokenIsErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tokenIsErr.Error())
	}

	c.Set("TDT-Auth-Token", token)

	return c.Status(fiber.StatusOK).JSON(common.NewResponse(fiber.StatusOK, "Login successfully", nil))
}

func AuthRegister(c *fiber.Ctx) error {
	bodyData, err := common.Validator[req.AuthRegister](c)

	if err != nil || bodyData == nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var userEntity entity.User

	if err := common.DBConn.First(&userEntity, "email = ? OR username = ? OR phone = ?", bodyData.Email, bodyData.Username, bodyData.Phone).Error; err == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Email, username or phone number already exists")
	}

	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(bodyData.Password), bcrypt.DefaultCost)

	if hashedPasswordErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while hashing password")
	}

	newUser := entity.User{
		UserId:   uuid.New(),
		Username: bodyData.Username,
		Password: string(hashedPassword),
		Name:     bodyData.Name,
		Email:    bodyData.Email,
		Phone:    bodyData.Phone,
		Birthday: bodyData.Birthday,
	}

	if err := common.DBConn.Create(&newUser).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while creating user")
	}

	token, tokenIsErr := createJWT(newUser.UserId)

	if tokenIsErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while creating token")
	}

	c.Set("TDT-Auth-Token", token)
	return c.Status(fiber.StatusCreated).JSON(common.NewResponse(fiber.StatusOK, "Register successfully", newUser))
}

func AuthVerify(c *fiber.Ctx) error {
	currentUserID, currenUserIdIsOk := c.Locals("currentUserId").(string)
	if !currenUserIdIsOk {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	userEntity := entity.User{}
	if err := common.DBConn.First(&userEntity, "user_id = ?", currentUserID).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}

		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	return c.JSON(common.NewResponse(
		fiber.StatusOK,
		"Success",
		userEntity),
	)
}
