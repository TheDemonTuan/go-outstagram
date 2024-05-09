package router

import (
	"github.com/gofiber/fiber/v2"
	"outstagram/controllers"
	"outstagram/middleware"
)

func authRouter(r fiber.Router) {
	authRoute := r.Group("auth")

	authRoute.Add("POST", "login", controllers.AuthLogin)
	authRoute.Add("POST", "register", controllers.AuthRegister)
	authRoute.Add("GET", "verify", middleware.Protected(), controllers.AuthVerify)
}
