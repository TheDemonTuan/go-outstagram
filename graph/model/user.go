package model

type UserProfile struct {
	Username string  `json:"username"`
	User     *User   `json:"user"`
	Posts    []*Post `json:"posts"`
	Friends  []*User `json:"friends"`
}
