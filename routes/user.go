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
	userRoute.Add("PATCH", "me/avatar", userController.UserMeUploadAvatar)
	userRoute.Add("PATCH", "me/profile", userController.UserMeEditProfile)
	userRoute.Add("PATCH", "me/private", userController.UserMeEditPrivate)
	userRoute.Add("PATCH", "me/phone", userController.UserMeEditPhone)
	userRoute.Add("PATCH", "me/email", userController.UserMeEditEmail)
	userRoute.Add("GET", ":userID", userController.UserGetByUserID)
	//userRoute.Add("POST", "delete", userController.UserDelete)
}
