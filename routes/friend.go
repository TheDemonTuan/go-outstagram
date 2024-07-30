package routes

import (
	"github.com/gofiber/fiber/v2"
	"outstagram/controllers"
	"outstagram/services"
)

func friendRouter(r fiber.Router) {
	friendRoute := r.Group("friends")

	friendService := services.NewFriendService()
	friendController := controllers.NewFriendController(friendService)

	friendRoute.Get("", friendController.FriendList)
	friendRoute.Get("me", friendController.GetFriendsByUserMe)
	friendRoute.Add("POST", "/:toUserID/request", friendController.FriendSendRequest)
	friendRoute.Add("PATCH", "/:toUserID/accept", friendController.FriendAcceptRequest)
	friendRoute.Add("DELETE", "/:toUserID/reject", friendController.FriendRejectRequest)
	friendRoute.Get("/:toUserID", friendController.GetFriendByUserID)
}
