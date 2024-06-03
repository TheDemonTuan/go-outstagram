package graph

//go:generate go run github.com/99designs/gqlgen generate

import "outstagram/services"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userService     services.UserService
	postService     services.PostService
	postFileService services.PostFileService
	//postCommentService services.PostCommentService
	//postLikeService services.PostLikeService
}
