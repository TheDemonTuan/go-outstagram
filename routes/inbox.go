package routes

import (
	"github.com/gofiber/fiber/v2"
	"outstagram/controllers"
	"outstagram/services"
)

func inboxRouter(r fiber.Router) {
	inboxRoute := r.Group("inbox")

	inboxService := services.NewInboxService()
	inboxController := controllers.NewInboxController(inboxService)

	inboxRoute.Add("POST", ":toUserName", inboxController.InboxSendMessage)
}
