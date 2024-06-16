package model

type Post struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	Caption       string `json:"caption"`
	IsHideLike    bool   `json:"is_hide_like"`
	IsHideComment bool   `json:"is_hide_comment"`
	Privacy       int    `json:"privacy"`
	Active        bool   `json:"active"`

	User         *User          `json:"user"`
	PostFiles    []*PostFile    `json:"post_files"`
	PostLikes    []*PostLike    `json:"post_likes"`
	PostComments []*PostComment `json:"post_comments"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}
