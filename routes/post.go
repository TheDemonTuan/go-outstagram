package routes

import (
	"github.com/gofiber/fiber/v2"
	"outstagram/controllers"
	"outstagram/services"
)

func postRouter(r fiber.Router) {
	postRoute := r.Group("posts")

	postService := services.NewPostService()
	friendService := services.NewFriendService()
	postController := controllers.NewPostController(postService, friendService)

	//Me API
	postRoute.Add("GET", "me", postController.PostMeGetAll)
	postRoute.Add("POST", "me", postController.PostMeCreate)
	postRoute.Add("POST", "me/like/:postID", postController.PostMeLikeByPostID)
	postRoute.Add("POST", "me/comment/:postID", postController.PostMeCommentByPostID)
	postRoute.Add("PUT", "me/:postID", postController.PostMeEditByPostID)
	postRoute.Add("DELETE", "me/:postID", postController.PostMeDeleteByPostID)
	//User API
	postRoute.Add("GET", "user/:userID", postController.PostGetAllByUserID)
	postRoute.Add("GET", ":postID", postController.PostGetByPostID)
}
