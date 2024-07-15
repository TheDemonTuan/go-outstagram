package routes

import (
	"github.com/gofiber/fiber/v2"
	"outstagram/controllers"
	"outstagram/services"
)

func adminRouter(r fiber.Router) {
	adminService := services.NewAdminService()
	postService := services.NewPostService()
	userService := services.NewUserService()
	adminController := controllers.NewAdminController(adminService, postService, userService)

	r.Add("DELETE", "posts/:postID/:userID", adminController.AdminDeletePostByPostID)
	r.Add("POST", "ban/users/:userID", adminController.AdminBanUserByUserID)
	r.Add("POST", "block/posts/:postID", adminController.AdminBlockPostByPostID)
}
