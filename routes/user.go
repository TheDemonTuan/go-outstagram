package routes

import (
	"github.com/gofiber/fiber/v2"
	"outstagram/controllers"
	"outstagram/services"
)

func userRouter(r fiber.Router) {
	userRoute := r.Group("users")

	userService := services.NewUserService()
	userController := controllers.NewUserController(userService)

	userRoute.Add("GET", "me", userController.UserGetMe)
	userRoute.Add("PUT", "me/update", userController.UserEditMe)
	userRoute.Add("GET", ":userID", userController.UserGetByUserID)
	//userRoute.Add("POST", "delete", userController.UserDelete)
}
