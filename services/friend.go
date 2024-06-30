package services

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"outstagram/common"
	"outstagram/models/entity"
)

type FriendService struct{}

func NewFriendService() *FriendService {
	return &FriendService{}
}

func (f *FriendService) SendFriendRequest(friendRecord *entity.Friend, fromUserID, toUserID string) error {
	if fromUserID == toUserID {
		return fiber.NewError(fiber.StatusBadRequest, "You can't send friend request to yourself")
	}

	uuidFromUserID, err := uuid.Parse(fromUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	uuidToUserID, err := uuid.Parse(toUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	userReceiver := entity.User{}
	if err := common.DBConn.Where("id = ?", toUserID).First(&userReceiver).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusBadRequest, "User receiver not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error when querying user receiver")
	}

	isRequestSent := true
	if err := common.DBConn.Where("from_user_id = ? AND to_user_id = ? OR from_user_id = ? AND to_user_id = ?", fromUserID, toUserID, toUserID, fromUserID).First(&friendRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			isRequestSent = false
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when querying friend record")
		}
	}

	if isRequestSent && (friendRecord.Status == entity.FriendRequested || friendRecord.Status == entity.FriendAccepted) {
		return fiber.NewError(fiber.StatusBadRequest, "Friend request already sent")
	}

	friendRecord.FromUserID = uuidFromUserID
	friendRecord.ToUserID = uuidToUserID

	if isRequestSent && friendRecord.Status == entity.FriendRejected {
		friendRecord.Status = entity.FriendRequested
		if err := common.DBConn.Save(&friendRecord).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return nil
	}

	if err := common.DBConn.Create(&friendRecord).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (f *FriendService) RejectFriendRequest(friendRecord *entity.Friend, fromUserID, toUserID string) error {
	_, err := uuid.Parse(fromUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	_, err = uuid.Parse(toUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := common.DBConn.Where("from_user_id = ? AND to_user_id = ? OR from_user_id = ? AND to_user_id = ?", fromUserID, toUserID, toUserID, fromUserID).First(&friendRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusBadRequest, "Friend request not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error when querying friend record")
	}

	if friendRecord.Status == entity.FriendRejected {
		return fiber.NewError(fiber.StatusBadRequest, "Friend request already rejected")
	}

	friendRecord.Status = entity.FriendRejected

	if err := common.DBConn.Save(&friendRecord).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (f *FriendService) GetFriendList(userID string) ([]entity.Friend, error) {
	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var friends []entity.Friend
	if err := common.DBConn.Where("from_user_id = ? OR to_user_id = ?", userID, userID).Preload("ToUser").Find(&friends).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var friendList []entity.Friend
	for _, friend := range friends {
		if friend.FromUserID.String() == userID {
			friendList = append(friendList, friend)
		} else {
			friendList = append(friendList, friend)
		}
	}

	return friendList, nil
}

func (f *FriendService) GetFriendByUserID(friendRecord *entity.Friend, fromUserID, toUserID string) error {
	_, err := uuid.Parse(fromUserID)
	if err != nil {
		return errors.New("error when parsing fromUserID")
	}

	_, err = uuid.Parse(toUserID)
	if err != nil {
		return errors.New("error when parsing toUserID")
	}

	if err := common.DBConn.Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?) AND status = ?", fromUserID, toUserID, toUserID, fromUserID, entity.FriendAccepted).First(&friendRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("friend not found")
		}
		return errors.New("error when querying friend record")
	}

	return nil
}

func (f *FriendService) GetFriendStatusByUserID(friendRecord *entity.Friend, fromUserID, toUserID string) error {
	_, err := uuid.Parse(fromUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	_, err = uuid.Parse(toUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := common.DBConn.Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)", fromUserID, toUserID, toUserID, fromUserID).First(&friendRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusBadRequest, "Friend record not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error when querying friend record")
	}

	return nil
}

func (f *FriendService) AcceptFriendRequest(friendRecord *entity.Friend, fromUserID, toUserID string) error {
	_, err := uuid.Parse(fromUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	_, err = uuid.Parse(toUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := common.DBConn.Where("from_user_id = ? AND to_user_id = ?", toUserID, fromUserID).First(&friendRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusBadRequest, "Friend request not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error when querying friend record")
	}

	if friendRecord.Status == entity.FriendAccepted {
		return fiber.NewError(fiber.StatusBadRequest, "Friend request already accepted")
	}

	if friendRecord.Status == entity.FriendRejected {
		return fiber.NewError(fiber.StatusBadRequest, "Friend request already rejected")
	}

	friendRecord.Status = entity.FriendAccepted

	if err := common.DBConn.Save(&friendRecord).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (f *FriendService) GetAllFriendsByUserName(userName string, friends interface{}) error {
	var user entity.User
	if err := common.DBConn.Select("id").Where("username = ?", userName).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("error when querying user")
	}

	if err := common.DBConn.Model(&entity.Friend{}).Where("(from_user_id = ? OR to_user_id = ?) AND status = ?", user.ID.String(), user.ID.String(), entity.FriendAccepted).Find(friends).Error; err != nil {
		return errors.New("error when querying friends")
	}

	return nil
}
