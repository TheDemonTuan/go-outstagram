package req

import "time"

type UserMeUpdate struct {
	Email    string    `json:"email" validate:"required,email,min=5,max=100"`
	Username string    `json:"username" validate:"required,alphanum,min=3,max=50"`
	FullName string    `json:"full_name" validate:"required,min=3,max=100"`
	Phone    string    `json:"phone" validate:"required,min=10,max=15"`
	Birthday time.Time `json:"birthday" validate:"required"`
	Avatar   string    `json:"avatar"`
	Bio      string    `json:"bio"`
	Gender   bool      `json:"gender"`
}
