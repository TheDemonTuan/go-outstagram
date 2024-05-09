package services

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

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) AuthenticateUser(usernameOrEmailOrPhone, password string) (*entity.User, error) {
	var userEntity entity.User

	db := common.DBConn.Where("username = ? or email = ? or phone = ?", usernameOrEmailOrPhone, usernameOrEmailOrPhone, usernameOrEmailOrPhone)

	if err := db.First(&userEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid credentials")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Error while querying user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(password)); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid credentials")
	}

	return &userEntity, nil
}

func (s *AuthService) CreateUser(bodyData *req.AuthRegister) (entity.User, error) {
	var existingUser entity.User
	if err := common.DBConn.First(&existingUser, "email = ? OR username = ? OR phone = ?", bodyData.Email, bodyData.Username, bodyData.Phone).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, fiber.NewError(fiber.StatusInternalServerError, "Error while querying user")
		}
	}

	if existingUser.ID != uuid.Nil {
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, "User already exists")
	}

	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(bodyData.Password), bcrypt.DefaultCost)
	if hashedPasswordErr != nil {
		return entity.User{}, fiber.NewError(fiber.StatusInternalServerError, "Error while hashing password")
	}

	newUser := entity.User{
		ID:       uuid.New(),
		Username: bodyData.Username,
		Password: string(hashedPassword),
		FullName: bodyData.FullName,
		Email:    bodyData.Email,
		Phone:    bodyData.Phone,
		Birthday: bodyData.Birthday,
	}
	if err := common.DBConn.Create(&newUser).Error; err != nil {
		return entity.User{}, fiber.NewError(fiber.StatusInternalServerError, "Error while creating user")
	}

	return newUser, nil
}

func (s *AuthService) VerifyUser(currentUserID string) (*entity.User, error) {
	var user entity.User
	if err := common.DBConn.First(&user, "id = ?", currentUserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Error while querying user")
	}
	return &user, nil
}

func (s *AuthService) CreateJWT(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"uuid": userId.String(),
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, tokenSignedErr := token.SignedString([]byte("secret"))
	if tokenSignedErr != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Error while signing token")
	}

	return tokenSigned, nil
}
