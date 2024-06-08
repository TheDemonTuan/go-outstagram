package routes

import (
	"github.com/gofiber/fiber/v2"
	"outstagram/middleware"
)

func SetupRouter(app *fiber.App) {
	graphqlRoute := app.Group("graphql", middleware.Protected(), middleware.GraphqlHandler)
	graphqlRouter(graphqlRoute)

	publicAPIRoute := app.Group("api")
	authRouter(publicAPIRoute)

	privateAPIRoute := app.Group("api", middleware.Protected())
	userRouter(privateAPIRoute)
	postRouter(privateAPIRoute)
	friendRouter(privateAPIRoute)

	adminAPIRoute := app.Group("api/admin", middleware.Protected(), middleware.IsAdmin)
	adminRouter(adminAPIRoute)
}
