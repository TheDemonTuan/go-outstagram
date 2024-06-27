package model

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	FullName  *string   `json:"full_name,omitempty"`
	Email     *string   `json:"email,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	Avatar    *string   `json:"avatar,omitempty"`
	Bio       *string   `json:"bio,omitempty"`
	Birthday  *string   `json:"birthday,omitempty"`
	Gender    *bool     `json:"gender,omitempty"`
	Role      *bool     `json:"role,omitempty"`
	Active    *bool     `json:"active,omitempty"`
	IsPrivate *bool     `json:"is_private,omitempty"`
	Friends   []*Friend `json:"friends" gorm:"-"`
	CreatedAt *string   `json:"created_at,omitempty"`
	UpdatedAt *string   `json:"updated_at,omitempty"`
	DeletedAt *string   `json:"deleted_at,omitempty"`
}

type UserProfile struct {
	Username string  `json:"username"`
	User     *User   `json:"user"`
	Posts    []*Post `json:"posts"`
}

type UserSuggestion struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
	Role     bool   `json:"role"`
	Active   bool   `json:"active"`

	Posts   []*Post `json:"posts" gorm:"-"`
	Friends []*User `json:"friends" gorm:"-"`
}
