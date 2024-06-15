package model

type Inbox struct {
	ID         string `json:"id"`
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
	Message    string `json:"message"`
	IsRead     bool   `json:"is_read"`

	Files []InboxFile `json:"files"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}
