package controllers

import (
	"github.com/gofiber/fiber/v2"
	"outstagram/common"
	"outstagram/models/entity"
	"outstagram/models/req"
	"outstagram/services"
)

type InboxController struct {
	inboxService *services.InboxService
}

func NewInboxController(inboxService *services.InboxService) *InboxController {
	return &InboxController{
		inboxService: inboxService,
	}
}

func (ib *InboxController) InboxSendMessage(ctx *fiber.Ctx) error {
	// Get current user
	currentUserInfo, currenUserInfoIsOk := ctx.Locals(common.UserInfoLocalKey).(entity.User)
	if !currenUserInfoIsOk {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user")
	}

	bodyData, err := common.RequestBodyValidator[req.InboxSendMessage](ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	toUserName := ctx.Params("toUserName")

	inboxRecord, err := ib.inboxService.SendMessage(currentUserInfo, toUserName, bodyData.Message)
	if err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusCreated, "Message sent", inboxRecord)
}
