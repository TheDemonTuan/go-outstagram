package routes

import (
	"github.com/gofiber/fiber/v2"
	"outstagram/controllers"
	"outstagram/services"
)

func adminRouter(r fiber.Router) {
	adminRoute := r.Group("admin")

	adminService := services.NewAdminService()
	postService := services.NewPostService()
	userService := services.NewUserService()
	adminController := controllers.NewAdminController(adminService, postService, userService)

	adminRoute.Add("DELETE", "posts/:postID", adminController.AdminDeletePostByPostID)
	adminRoute.Add("POST", "ban/users/:userID", adminController.AdminBanUserByUserID)
}
