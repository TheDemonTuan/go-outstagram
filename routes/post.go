package routes

import (
	"github.com/gofiber/fiber/v2"
	"outstagram/controllers"
	"outstagram/services"
)

func postRouter(r fiber.Router) {
	postRoute := r.Group("posts")

	postService := services.NewPostService()
	postController := controllers.NewPostController(postService)

	postRoute.Add("GET", "", postController.PostGetAll)
	postRoute.Add("POST", "", postController.PostCreate)
}
