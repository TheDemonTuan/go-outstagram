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
	userService := services.NewUserService()
	tokenService := services.NewTokenService()
	authController := controllers.NewAuthController(authService, userService, tokenService)

	authRoute.Add("POST", "login", authController.AuthLogin)
	authRoute.Add("POST", "register", authController.AuthRegister)
	authRoute.Add("DELETE", "logout", authController.AuthLogout)
	authRoute.Add("POST", "oauth/login", authController.AuthOAuthLogin)
	authRoute.Add("POST", "oauth/register", authController.AuthOAuthRegister)
	authRoute.Add("GET", "verify", middleware.Protected(), authController.AuthVerify)
	authRoute.Add("POST", "refresh", authController.AuthRefreshToken)
}
