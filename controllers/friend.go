package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"outstagram/common"
	"outstagram/models/entity"
	"outstagram/services"
)

type FriendController struct {
	friendService *services.FriendService
}

func NewFriendController(friendService *services.FriendService) *FriendController {
	return &FriendController{
		friendService: friendService,
	}
}

func (f *FriendController) FriendSendRequest(ctx *fiber.Ctx) error {
	// Get current user
	currentUserInfo, currenUserInfoIsOk := ctx.Locals(common.UserInfoLocalKey).(entity.User)
	if !currenUserInfoIsOk {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user")
	}

	// Get to user
	toUserID := ctx.Params("toUserID")

	var friendRecord entity.Friend
	if err := f.friendService.SendFriendRequest(&friendRecord, currentUserInfo.ID.String(), toUserID); err != nil {
		return err
	}

	// Push notification
	data := map[string]string{"message": "" + currentUserInfo.Username + " sent you a friend request!", "fromUserID": currentUserInfo.ID.String()}

	if err := common.PusherClient.Trigger(toUserID, "friend-notification", data); err != nil {
		fmt.Println(err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusCreated, "Friend request sent", friendRecord)

}

func (f *FriendController) FriendAcceptRequest(ctx *fiber.Ctx) error {
	return nil
}

func (f *FriendController) FriendRejectRequest(ctx *fiber.Ctx) error {
	currentUserID, currenUserIdIsOk := ctx.Locals(common.UserIDLocalKey).(string)
	if !currenUserIdIsOk {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user")
	}

	// Get to user
	toUserID := ctx.Params("toUserID")

	var friendRecord entity.Friend
	if err := f.friendService.RejectFriendRequest(&friendRecord, currentUserID, toUserID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Friend request rejected", friendRecord)
}

func (f *FriendController) FriendList(ctx *fiber.Ctx) error {
	currentUserID, currenUserIdIsOk := ctx.Locals(common.UserIDLocalKey).(string)
	if !currenUserIdIsOk {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user")
	}

	friendList, err := f.friendService.GetFriendList(currentUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Friend list", friendList)
}

func (f *FriendController) GetFriendByUserID(ctx *fiber.Ctx) error {
	currentUserID, currenUserIdIsOk := ctx.Locals(common.UserIDLocalKey).(string)
	if !currenUserIdIsOk {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user")
	}

	toUserID := ctx.Params("toUserID")

	var friendRecord entity.Friend
	if err := f.friendService.GetFriendByUserID(&friendRecord, currentUserID, toUserID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Friend record", friendRecord)
}
