// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type InboxFile struct {
	ID        string  `json:"id"`
	InboxID   string  `json:"inbox_id"`
	Type      *string `json:"type,omitempty"`
	URL       *string `json:"url,omitempty"`
	Active    *bool   `json:"active,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
	UpdatedAt *string `json:"updated_at,omitempty"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

type InboxGetAllBubble struct {
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
	FullName    string `json:"full_name"`
	LastMessage string `json:"last_message"`
	IsRead      bool   `json:"is_read"`
	CreatedAt   string `json:"created_at"`
}

type PostComment struct {
	ID        *string `json:"id,omitempty"`
	PostID    *string `json:"post_id,omitempty"`
	UserID    *string `json:"user_id,omitempty"`
	Content   *string `json:"content,omitempty"`
	Active    *bool   `json:"active,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
	UpdatedAt *string `json:"updated_at,omitempty"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

type PostFile struct {
	ID        *string `json:"id,omitempty"`
	PostID    *string `json:"post_id,omitempty"`
	URL       *string `json:"url,omitempty"`
	Type      *string `json:"type,omitempty"`
	Active    *bool   `json:"active,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
	UpdatedAt *string `json:"updated_at,omitempty"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

type PostLike struct {
	ID        *string `json:"id,omitempty"`
	PostID    *string `json:"post_id,omitempty"`
	UserID    *string `json:"user_id,omitempty"`
	IsLiked   *bool   `json:"is_liked,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
	UpdatedAt *string `json:"updated_at,omitempty"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

type Query struct {
}

type User struct {
	ID        string  `json:"id"`
	Username  *string `json:"username,omitempty"`
	FullName  *string `json:"full_name,omitempty"`
	Email     *string `json:"email,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Avatar    *string `json:"avatar,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	Birthday  *string `json:"birthday,omitempty"`
	Gender    *bool   `json:"gender,omitempty"`
	Role      *bool   `json:"role,omitempty"`
	Active    *bool   `json:"active,omitempty"`
	IsPrivate *bool   `json:"is_private,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
	UpdatedAt *string `json:"updated_at,omitempty"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

type UserSearch struct {
	ID       *string `json:"id,omitempty"`
	Username *string `json:"username,omitempty"`
	FullName *string `json:"full_name,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
	Active   *bool   `json:"active,omitempty"`
}

type UserSuggestion struct {
	ID       *string `json:"id,omitempty"`
	Username *string `json:"username,omitempty"`
	FullName *string `json:"full_name,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
	Active   *bool   `json:"active,omitempty"`
}
