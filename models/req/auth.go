package req

import "time"

type AuthLogin struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=8"`
	Token    string `json:"token"`
}

type AuthRegister struct {
	Username string    `json:"username" validate:"required,min=3,max=30"`
	Password string    `json:"password" validate:"required,min=8"`
	Email    string    `json:"email" validate:"required,email,min=5,max=100"`
	Phone    string    `json:"phone" validate:"required,min=10,max=15"`
	Name     string    `json:"name" validate:"required,min=3,max=100"`
	Birthday time.Time `json:"birthday" validate:"required"`
}
