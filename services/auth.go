package services

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"os"
	"outstagram/common"
	"outstagram/models/entity"
	"outstagram/models/req"
	"regexp"
	"time"
)

const (
	fromEmail = "nguyenninhdan123456@gmail.com"
)

func sendEmail(fromEmail, toEmail, subject, body string) error {
	sg := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	from := mail.NewEmail("Outstagram", fromEmail)
	to := mail.NewEmail("Recipient", toEmail)
	message := mail.NewSingleEmail(from, subject, to, body, body)

	// Gửi email và lấy phản hồi
	response, err := sg.Send(message)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	log.Printf("SendGrid Response Status Code: %d", response.StatusCode)

	return nil
}

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

	var otpRecord entity.Otp
	if err := common.DBConn.Where("user_email = ? AND is_used = ?", bodyData.Email, true).First(&otpRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, fiber.NewError(fiber.StatusBadRequest, "Email has not been verified")
		}
		return entity.User{}, fiber.NewError(fiber.StatusInternalServerError, "Error while checking OTP verification")
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

func (s *AuthService) sendOtpEmail(email, otp string) error {
	subject := otp + " is your Outstagram code"
	body := "Hi,\n\nSomeone tried to sign up / reset password for an Outstagram account with " + " " + email + ". If it was you, enter this confirmation code in the app:\n\n " + otp
	err := sendEmail(fromEmail, email, subject, body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (s *AuthService) sendOtpEmailResetPassword(email, otp string) error {
	subject := email + ", we've made it easy to get back on Outstagram"
	body := "Hi " + email + ",\n\nSorry to hear you’re having trouble logging into Instagram. We got a message that you forgot your password. If it was you, enter this confirmation code in the app:\n\n " + otp
	err := sendEmail(fromEmail, email, subject, body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (s *AuthService) generateAndSaveNewOTP(existingOtp *entity.Otp) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Outstagram",
		AccountName: existingOtp.UserEmail,
	})
	if err != nil {
		return "", err
	}

	otp, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		return "", err
	}

	existingOtp.OtpCode = otp
	existingOtp.ExpiresAt = time.Now().Add(10 * time.Minute)

	if err := common.DBConn.Save(existingOtp).Error; err != nil {
		return "", err
	}

	return otp, nil
}

func (s *AuthService) generateAndSaveNewOTPResetPassword(existingOtp *entity.Otp) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Outstagram",
		AccountName: existingOtp.UserEmail,
	})
	if err != nil {
		return "", err
	}

	otp, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		return "", err
	}

	existingOtp.OtpCode = otp
	existingOtp.IsUsed = false
	existingOtp.ExpiresAt = time.Now().Add(10 * time.Minute)

	if err := common.DBConn.Save(existingOtp).Error; err != nil {
		return "", err
	}

	return otp, nil
}

func (s *AuthService) GenerateAndSaveOTP(bodyData *req.AuthOTPSendEmail) error {
	var existingUser entity.User

	if err := common.DBConn.First(&existingUser, "email = ? OR username = ?", bodyData.UserEmail, bodyData.Username).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusInternalServerError, "Error while querying user")
		}
	}

	if existingUser.ID != uuid.Nil {
		return fiber.NewError(fiber.StatusBadRequest, "User already exists")
	}

	if !s.IsUserAtLeast13(bodyData.Birthday) {
		return fiber.NewError(fiber.StatusBadRequest, "User must be at least 13 years old")
	}

	var existingOtp entity.Otp

	if err := common.DBConn.First(&existingOtp, "user_email = ? AND is_used = ?", bodyData.UserEmail, false).Error; err == nil {
		if existingOtp.ExpiresAt.Before(time.Now()) {
			otp, err := s.generateAndSaveNewOTP(&existingOtp)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			return s.sendOtpEmail(bodyData.UserEmail, otp)
		} else {
			return s.sendOtpEmail(bodyData.UserEmail, existingOtp.OtpCode)
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying otp")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Outstagram",
		AccountName: bodyData.UserEmail,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	otp, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	otpRecord := entity.Otp{
		ID:        uuid.New(),
		UserEmail: bodyData.UserEmail,
		OtpCode:   otp,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := common.DBConn.Create(&otpRecord).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while creating otp")
	}

	return s.sendOtpEmail(bodyData.UserEmail, otp)

}
func (s *AuthService) VerifyOTP(bodyData *req.AuthOTPVerifyEmail) error {
	var otpRecord entity.Otp

	if err := common.DBConn.Where("user_email = ? and otp_code = ?", bodyData.UserEmail, bodyData.OtpCode).First(&otpRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "OTP not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying OTP")
	}

	if otpRecord.IsUsed {
		return fiber.NewError(fiber.StatusBadRequest, "OTP has already been used")
	}

	if time.Now().After(otpRecord.ExpiresAt) {
		return fiber.NewError(fiber.StatusBadRequest, "OTP has expired")
	}

	otpRecord.IsUsed = true
	if err := common.DBConn.Save(&otpRecord).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update OTP status")
	}

	return nil
}

func (s *AuthService) GenerateAndSendPasswordResetOTP(bodyData *req.AuthOTPSendEmailResetPassword) error {
	var existingUser entity.User

	if err := common.DBConn.First(&existingUser, "email = ?", bodyData.Email).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusInternalServerError, "Error while querying user")
		}
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	var existingOtp entity.Otp

	if err := common.DBConn.First(&existingOtp, "user_email = ?", bodyData.Email).Error; err == nil {

		if existingOtp.IsUsed || existingOtp.ExpiresAt.Before(time.Now()) {
			otp, err := s.generateAndSaveNewOTPResetPassword(&existingOtp)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Error while querying OTP")
			}
			return s.sendOtpEmailResetPassword(bodyData.Email, otp)
		}

		return s.sendOtpEmailResetPassword(bodyData.Email, existingOtp.OtpCode)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying OTP")
	}

	return nil
}

func (s *AuthService) AuthResetPasswordSaveToDB(ctx *fiber.Ctx, bodyData *req.AuthResetPassword) error {

	var userRecord entity.User

	if err := common.DBConn.First(&userRecord, "email = ?", bodyData.Email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying user")
	}

	if bcrypt.CompareHashAndPassword([]byte(userRecord.Password), []byte(bodyData.NewPassword)) == nil {
		return fiber.NewError(fiber.StatusBadRequest, "New password cannot be the same as the old password")
	}

	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(bodyData.NewPassword), bcrypt.DefaultCost)
	if hashedPasswordErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while hashing password")
	}

	userRecord.Password = string(hashedPassword)

	if err := common.DBConn.Omit("phone").Save(&userRecord).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error while updating user")
	}

	return nil

}

func (s *AuthService) GenerateAccessToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"uuid": userId,
		"exp":  time.Now().Add(time.Minute * 5).Unix(),
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
