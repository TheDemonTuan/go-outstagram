package services

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
	"outstagram/common"
	"outstagram/models/entity"
	"outstagram/models/req"
	"regexp"
	"time"
)

type AuthService struct {
	userService  *UserService
	tokenService *TokenService
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) AuthenticateUser(usernameOrEmailOrPhone, password string) (*entity.User, error) {
	var userEntity entity.User

	db := common.DBConn.Where("username = ? OR email = ? OR phone = ? AND oauth = ?", usernameOrEmailOrPhone, usernameOrEmailOrPhone, usernameOrEmailOrPhone, entity.OAuthDefault)

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

func (s *AuthService) ValidateFullName(fullName string) bool {
	match, err := regexp.Match(`^[a-zA-ZÀ-ỹ\s]+$`, []byte(fullName))
	if err != nil {
		return false
	}
	return match
}

func (s *AuthService) IsUserAtLeast13(birthday time.Time) bool {
	thirteenYearsAgo := time.Now().AddDate(-13, 0, 0)
	return birthday.Before(thirteenYearsAgo) || birthday.Equal(thirteenYearsAgo)
}

func (s *AuthService) CreateUser(bodyData *req.AuthRegister) (entity.User, error) {
	var existingUser entity.User
	if err := common.DBConn.First(&existingUser, "email = ? OR username = ?", bodyData.Email, bodyData.Username).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, fiber.NewError(fiber.StatusInternalServerError, "Error while querying user")
		}
	}

	if existingUser.ID != uuid.Nil {
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, "User already exists")
	}

	if !s.IsUserAtLeast13(bodyData.Birthday) {
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, "User must be at least 13 years old")
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
		Birthday: bodyData.Birthday,
	}
	if err := common.DBConn.Omit("phone").Create(&newUser).Error; err != nil {
		return entity.User{}, fiber.NewError(fiber.StatusInternalServerError, "Error while creating user")
	}

	return newUser, nil
}

func (s *AuthService) GenerateAccessToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"uuid": userId,
		"exp":  time.Now().Add(time.Minute * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenSigned, tokenSignedErr := token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if tokenSignedErr != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Error while signing access token")
	}

	return tokenSigned, nil
}

func (s *AuthService) GenerateRefreshToken(userId string) (string, error) {
	exp := time.Now().Add(time.Hour * 24 * 15)

	claims := jwt.MapClaims{
		"uuid": userId,
		"exp":  exp.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenSigned, tokenSignedErr := token.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
	if tokenSignedErr != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Error while signing refresh token")
	}

	if err := s.tokenService.SaveRefreshToken(userId, tokenSigned, exp); err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Error while saving refresh token")
	}

	return tokenSigned, nil
}

func (s *AuthService) ValidateRefreshToken(refreshToken string, isCheckDB bool) (string, error) {
	claims := jwt.MapClaims{}
	token, tokenErr := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("REFRESH_TOKEN_SECRET")), nil
	})

	if tokenErr != nil {
		return "", fiber.NewError(fiber.StatusBadRequest, "Invalid refresh token")
	}

	if !token.Valid {
		return "", fiber.NewError(fiber.StatusBadRequest, "Invalid refresh token")
	}

	userId, isOK := claims["uuid"].(string)
	if !isOK {
		return "", fiber.NewError(fiber.StatusBadRequest, "Invalid refresh token")
	}

	if isCheckDB {
		if _, err := s.tokenService.GetRefreshTokenByToken(refreshToken); err != nil {
			return "", fiber.NewError(fiber.StatusBadRequest, "Invalid refresh token")
		}
	}

	return userId, nil
}

func (s *AuthService) AuthOAuthLogin(bodyData *req.AuthOAuthLogin) (entity.User, error) {
	if bodyData.Provider == entity.OAuthDefault.EnumIndex() {
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, "Invalid provider")
	}

	var userEntity entity.User
	if err := common.DBConn.First(&userEntity, "email = ?", bodyData.Email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, fiber.NewError(fiber.StatusBadRequest, "User not found")
		}
		return entity.User{}, fiber.NewError(fiber.StatusInternalServerError, "Error while querying user")
	}

	if userEntity.OAuth != entity.UserOAuth(bodyData.Provider) {
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, "Invalid provider")
	}

	return userEntity, nil
}

func (s *AuthService) AuthOAuthRegister(bodyData *req.AuthOAuthRegister) (entity.User, error) {
	if bodyData.Provider == entity.OAuthDefault.EnumIndex() {
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, "Invalid provider")
	}

	var existingUser entity.User
	if err := common.DBConn.First(&existingUser, "email = ? or username = ?", bodyData.Email, bodyData.Username).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, fiber.NewError(fiber.StatusInternalServerError, "Error while querying user")
		}
	}

	if existingUser.ID != uuid.Nil {
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, "User already exists")
	}

	if !s.IsUserAtLeast13(bodyData.Birthday) {
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, "User must be at least 13 years old")
	}

	newUser := entity.User{
		ID:       uuid.New(),
		Username: bodyData.Username,
		FullName: bodyData.FullName,
		Email:    bodyData.Email,
		Birthday: bodyData.Birthday,
		OAuth:    entity.UserOAuth(bodyData.Provider),
	}

	if err := common.DBConn.Omit("phone").Create(&newUser).Error; err != nil {
		return entity.User{}, fiber.NewError(fiber.StatusInternalServerError, "Error while creating user")
	}

	return newUser, nil
}
