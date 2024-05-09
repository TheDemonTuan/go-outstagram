package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func SetupRouter(app *fiber.App) {
	publicAPIRoute := app.Group("api")
	publicAPIRoute.Add("GET", "metrics", monitor.New(monitor.Config{Title: "Social Network Metrics"}))
	authRouter(publicAPIRoute)

	//privateAPIRoute := app.Group("api", middleware.Protected())

}
