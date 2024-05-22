package graphql

import (
	"errors"
	"github.com/graphql-go/graphql"
	"outstagram/models/entity"
	"outstagram/services"
)

var userTypes = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"full_name": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"phone": &graphql.Field{
				Type: graphql.String,
			},
			"avatar": &graphql.Field{
				Type: graphql.String,
			},
			"bio": &graphql.Field{
				Type: graphql.String,
			},
			"birthday": &graphql.Field{
				Type: graphql.DateTime,
			},
			"gender": &graphql.Field{
				Type: graphql.Boolean,
			},
			"role": &graphql.Field{
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
		},
	})

var userGetByUserID = &graphql.Field{
	Name:        "GetUserByUserID",
	Type:        userTypes,
	Description: "Get user by user id",
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

		var user entity.User
		userService := services.NewUserService()
		if err := userService.UserGetByUserID(id, &user); err != nil {
			return nil, err
		}

		return user, nil
	},
}

var userSearchByUsernameOrFullName = &graphql.Field{
	Name:        "SearchUserByUsernameOrFullName",
	Type:        graphql.NewList(userTypes),
	Description: "Search user by username or full name",
	Args: graphql.FieldConfigArgument{
		"keyword": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		keyword, ok := params.Args["keyword"].(string)
		if !ok {
			return nil, nil
		}

		var users []entity.User
		userService := services.NewUserService()
		if err := userService.UserSearchByUsernameOrFullName(keyword, &users); err != nil {
			return nil, err
		}

		return users, nil
	},
}

var userGetByUserName = &graphql.Field{
	Name:        "GetUserByUserName",
	Type:        userTypes,
	Description: "Get user by username",
	Args: graphql.FieldConfigArgument{
		"username": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		username, ok := params.Args["username"].(string)
		if !ok {
			return nil, nil
		}

		var user entity.User
		userService := services.NewUserService()
		if err := userService.UserGetByUserName(username, &user); err != nil {
			return nil, err
		}

		if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
			return nil, errors.New("user not found")
		}

		return user, nil
	},
}
