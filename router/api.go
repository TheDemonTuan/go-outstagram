package router

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {
	publicAPIRoute := app.Group("api")
	authRouter(publicAPIRoute)

	//privateAPIRoute := app.Group("api", middleware.Protected())
}
