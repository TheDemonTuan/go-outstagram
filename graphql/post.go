package graphql

import (
	"github.com/graphql-go/graphql"
	"outstagram/models/entity"
	"outstagram/services"
)

var postTypes = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Post",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"user_id": &graphql.Field{
				Type: graphql.String,
			},
			"caption": &graphql.Field{
				Type: graphql.String,
			},
			"is_hide_like": &graphql.Field{
				Type: graphql.Boolean,
			},
			"is_hide_comment": &graphql.Field{
				Type: graphql.Boolean,
			},
			"active": &graphql.Field{
				Type: graphql.Boolean,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"deleted_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"user": &graphql.Field{
				Type: userTypes,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					post := params.Source.(entity.Post)
					var user entity.User
					userService := services.NewUserService()
					if err := userService.UserGetByUserID(post.UserID.String(), &user); err != nil {
						return nil, err
					}

					return user, nil
				},
			},
			"files": &graphql.Field{
				Type: graphql.NewList(postFileTypes),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					post := params.Source.(entity.Post)
					var postFiles []entity.PostFile
					postFileService := services.NewPostFileService()
					if err := postFileService.PostFileGetAllByPostID(post.ID, &postFiles); err != nil {
						return nil, err
					}

					return postFiles, nil
				},
			},
		},
	},
)

var postGetAll = &graphql.Field{
	Name:        "GetAllPosts",
	Type:        graphql.NewList(postTypes),
	Description: "Get all post list",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		var posts []entity.Post
		postService := services.NewPostService()
		if err := postService.PostGetAll(&posts); err != nil {
			return nil, err
		}
		return posts, nil
	},
}

var postGetAllByUserID = &graphql.Field{
	Name:        "GetAllPostByUserID",
	Type:        graphql.NewList(postTypes),
	Description: "Get all post by user id",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		id, ok := params.Args["id"].(string)
		if !ok {
			return nil, nil
		}

		var posts []entity.Post
		postService := services.NewPostService()
		if err := postService.PostGetAllByUserID(id, &posts); err != nil {
			return nil, err
		}

		return posts, nil
	},
}
