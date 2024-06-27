package model

type Friend struct {
	ID         string `json:"id"`
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
	Status     int    `json:"status"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`

	FromUserInfo *User `json:"from_user_info" gorm:"-"`
	ToUserInfo   *User `json:"to_user_info" gorm:"-"`
}
