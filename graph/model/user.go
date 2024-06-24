package model

type UserProfile struct {
	Username string  `json:"username"`
	User     *User   `json:"user"`
	Posts    []*Post `json:"posts"`
	Friends  []*User `json:"friends"`
}

type UserSuggestion struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
	Role     bool   `json:"role"`
	Active   bool   `json:"active"`

	Posts   []*Post `json:"posts"`
	Friends []*User `json:"friends"`
}
