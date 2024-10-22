package controllers

import (
	"github.com/gofiber/fiber/v2"
	"outstagram/common"
	"outstagram/models/entity"
	"outstagram/services"
	"time"
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
	data := map[string]string{}
	data["type"] = "friend-action"
	data["action"] = "send"
	data["message"] = "send you a friend request!"
	data["fromUserID"] = currentUserInfo.ID.String()
	data["username"] = currentUserInfo.Username
	data["avatar"] = currentUserInfo.Avatar
	data["createdAt"] = time.Now().Format(time.RFC3339)

	if err := common.PusherClient.Trigger(toUserID, "internal-socket", data); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusCreated, "Friend request sent", friendRecord)

}

func (f *FriendController) FriendAcceptRequest(ctx *fiber.Ctx) error {
	currentUserInfo, currenUserInfoIsOk := ctx.Locals(common.UserInfoLocalKey).(entity.User)
	if !currenUserInfoIsOk {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user")
	}

	// Get to user
	toUserID := ctx.Params("toUserID")

	var friendRecord entity.Friend
	if err := f.friendService.AcceptFriendRequest(&friendRecord, currentUserInfo.ID.String(), toUserID); err != nil {
		return err
	}

	// Push notification
	data := map[string]string{}
	data["type"] = "friend-action"
	data["action"] = "accept"
	data["message"] = "accepted your friend request!"
	data["fromUserID"] = currentUserInfo.ID.String()
	data["username"] = currentUserInfo.Username
	data["avatar"] = currentUserInfo.Avatar
	data["createdAt"] = time.Now().Format(time.RFC3339)

	if err := common.PusherClient.Trigger(toUserID, "internal-socket", data); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Friend request accepted", friendRecord)
}

func (f *FriendController) FriendRejectRequest(ctx *fiber.Ctx) error {
	currentUserInfo, currenUserInfoIsOk := ctx.Locals(common.UserInfoLocalKey).(entity.User)
	if !currenUserInfoIsOk {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user")
	}

	// Get to user
	toUserID := ctx.Params("toUserID")

	var friendRecord entity.Friend
	if err := f.friendService.RejectFriendRequest(&friendRecord, currentUserInfo.ID.String(), toUserID); err != nil {
		return err
	}

	// Push notification
	data := map[string]string{}
	data["type"] = "friend-action"
	data["action"] = "reject"
	data["message"] = "rejected your friend request!"
	data["fromUserID"] = currentUserInfo.ID.String()
	data["username"] = currentUserInfo.Username
	data["avatar"] = currentUserInfo.Avatar
	data["createdAt"] = time.Now().Format(time.RFC3339)

	if err := common.PusherClient.Trigger(toUserID, "internal-socket", data); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
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
	if err := f.friendService.GetFriendStatusByUserID(&friendRecord, currentUserID, toUserID); err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Friend record", friendRecord)
}

func (f *FriendController) GetFriendsByUserMe(ctx *fiber.Ctx) error {
	currentUserID, currenUserIdIsOk := ctx.Locals(common.UserIDLocalKey).(string)
	if !currenUserIdIsOk {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user")
	}

	friendsMe, err := f.friendService.GetFriendsUserMe(currentUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Friend list me", friendsMe)

}
