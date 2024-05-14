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

	postRoute.Add("GET", "me", postController.PostGetAll)
	postRoute.Add("GET", "user/:userID", postController.PostGetAllByUserID)
	postRoute.Add("GET", ":postID", postController.PostGetByPostID)
	postRoute.Add("POST", "", postController.PostCreate)
	postRoute.Add("POST", "like/:postID", postController.PostLikeByPostID)
}
