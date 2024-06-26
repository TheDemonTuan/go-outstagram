package model

import "outstagram/models/entity"

type Post struct {
	ID            string             `json:"id"`
	UserID        string             `json:"user_id"`
	Caption       string             `json:"caption"`
	IsHideLike    bool               `json:"is_hide_like"`
	IsHideComment bool               `json:"is_hide_comment"`
	Privacy       entity.PostPrivacy `json:"privacy"`
	Active        bool               `json:"active"`

	User         *User          `json:"user"`
	PostFiles    []*PostFile    `json:"post_files"`
	PostLikes    []*PostLike    `json:"post_likes"`
	PostComments []*PostComment `json:"post_comments"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type PostComment struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	ParentID  string `json:"parent_id"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`

	User   *User        `json:"user"`
	Parent *PostComment `json:"parent"`
}

type PostLike struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	UserID    string `json:"user_id"`
	IsLiked   bool   `json:"is_liked"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`

	User *User `json:"user"`
}
