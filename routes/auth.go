package routes

import (
	"github.com/gofiber/fiber/v2"
	"outstagram/controllers"
	"outstagram/middleware"
	"outstagram/services"
)

func authRouter(r fiber.Router) {
	authRoute := r.Group("auth")

	authService := services.NewAuthService()
	authController := controllers.NewAuthController(authService)

	authRoute.Add("POST", "login", authController.AuthLogin)
	authRoute.Add("POST", "register", authController.AuthRegister)
	authRoute.Add("GET", "verify", middleware.Protected(), authController.AuthVerify)
}
