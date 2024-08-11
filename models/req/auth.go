package req

import "time"

type AuthLogin struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

type AuthRegister struct {
	Email    string    `json:"email" validate:"required,email,min=5,max=100"`
	Username string    `json:"username" validate:"required,alphanum,min=3,max=50"`
	Password string    `json:"password" validate:"required,min=8"`
	FullName string    `json:"full_name" validate:"required,min=3,max=100"`
	Birthday time.Time `json:"birthday" validate:"required"`
}

type AuthRefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required,jwt"`
}

type AuthLogout struct {
	RefreshToken string `json:"refresh_token" validate:"required,jwt"`
}

type AuthOAuthLogin struct {
	Email    string `json:"email" validate:"required,email,min=5,max=100"`
	Provider int    `json:"provider" validate:"required"`
}

type AuthOAuthRegister struct {
	Email    string    `json:"email" validate:"required,email,min=5,max=100"`
	Username string    `json:"username" validate:"required,alphanum,min=3,max=50"`
	FullName string    `json:"full_name" validate:"required,min=3,max=100"`
	Birthday time.Time `json:"birthday" validate:"required"`
	Provider int       `json:"provider" validate:"required"`
}

type AuthOTPSendEmail struct {
	UserEmail string    `json:"user_email" validate:"required,email,min=5,max=100"`
	Username  string    `json:"username" validate:"required,alphanum,min=3,max=50"`
	Password  string    `json:"password" validate:"required,min=8"`
	FullName  string    `json:"full_name" validate:"required,min=3,max=100"`
	Birthday  time.Time `json:"birthday" validate:"required"`
}

type AuthOTPVerifyEmail struct {
	UserEmail string `json:"user_email" validate:"required,email,min=5,max=100"`
	OtpCode   string `json:"otp_code" validate:"required,min=6,max=6"`
}

type AuthOTPSendEmailResetPassword struct {
	Email string `json:"user_email" validate:"required,email,min=5,max=100"`
}

type AuthResetPassword struct {
	NewPassword string `json:"new_password" validate:"required,min=8"`
	Email       string `json:"email" validate:"required,email,min=5,max=100"`
}
