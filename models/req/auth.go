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
