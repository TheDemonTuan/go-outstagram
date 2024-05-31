package routes

import (
	"github.com/gofiber/fiber/v2"
	"outstagram/graphql"
	"outstagram/middleware"
)

func SetupRouter(app *fiber.App) {
	graphqlRoute := app.Group("graphql")
	graphqlRoute.Add("GET", "", graphql.Query)
	graphqlRoute.Add("POST", "", graphql.Mutation)

	publicAPIRoute := app.Group("api")
	authRouter(publicAPIRoute)

	privateAPIRoute := app.Group("api", middleware.Protected())
	userRouter(privateAPIRoute)
	postRouter(privateAPIRoute)
	friendRouter(privateAPIRoute)

	adminAPIRoute := privateAPIRoute.Use(middleware.IsAdmin)
	adminRouter(adminAPIRoute)
}
