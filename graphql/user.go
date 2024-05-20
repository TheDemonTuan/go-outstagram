package graphql

import (
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
			"password": &graphql.Field{
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
