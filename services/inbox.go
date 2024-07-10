package services

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"outstagram/common"
	"outstagram/graph/model"
	"outstagram/models/entity"
	"sync"
	"time"
)

type InboxService struct{}

func NewInboxService() *InboxService {
	return &InboxService{}
}

func (ib *InboxService) SendMessage(fromUserInfo entity.User, toUserName, message string) (entity.Inbox, error) {
	fromUserUUID, err := uuid.Parse(fromUserInfo.ID.String())
	if err != nil {
		return entity.Inbox{}, fiber.NewError(fiber.StatusBadRequest, "Invalid from user ID")
	}

	userRecord := entity.User{}
	if err := common.DBConn.Where("username = ?", toUserName).First(&userRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Inbox{}, fiber.NewError(fiber.StatusBadRequest, "User not found")
		}
		return entity.Inbox{}, fiber.NewError(fiber.StatusInternalServerError, "Error while getting user record")
	}

	if fromUserInfo.ID.String() == userRecord.ID.String() {
		return entity.Inbox{}, fiber.NewError(fiber.StatusBadRequest, "You can't send message to yourself")
	}

	inboxRecord := entity.Inbox{}
	inboxRecord.FromUserID = fromUserUUID
	inboxRecord.ToUserID = userRecord.ID
	inboxRecord.Message = message

	if err := common.DBConn.Save(&inboxRecord).Error; err != nil {
		return entity.Inbox{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Push notification
	data := map[string]string{}
	data["type"] = "inbox-action"
	data["action"] = "send"
	data["message"] = "accepted your friend request!"
	data["fromUserID"] = fromUserInfo.ID.String()
	data["username"] = fromUserInfo.Username
	data["toUserName"] = userRecord.Username
	//data["avatar"] = currentUserInfo.Avatar
	data["createdAt"] = time.Now().Format(time.RFC3339)

	if err := common.PusherClient.Trigger(userRecord.ID.String(), "internal-socket", data); err != nil {
		return entity.Inbox{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return inboxRecord, nil
}

func (ib *InboxService) InboxGetAllByUserName(fromUserID, username string, inboxRecords interface{}) error {
	var userRecord entity.User
	if err := common.DBConn.Where("username = ?", username).First(&userRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user to send not found")
		}
		return errors.New("error while getting user record")
	}

	if fromUserID == userRecord.ID.String() {
		return errors.New("you can't send message to yourself")
	}

	if err := common.DBConn.Model(entity.Inbox{}).Where("from_user_id = ? OR to_user_id = ?", userRecord.ID, userRecord.ID).Order("created_at ASC").Find(inboxRecords).Error; err != nil {
		return errors.New("error while getting inbox records")
	}

	return nil
}

func fetchUserAndInboxRecord(userID string, wg *sync.WaitGroup, ch chan<- *model.InboxGetAllBubble, errCh chan<- error) {
	defer wg.Done()

	var userRecord entity.User
	if err := common.DBConn.Select("username", "full_name", "avatar").Where("id = ?", userID).First(&userRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errCh <- errors.New("user not found")
			return
		}
		errCh <- errors.New("error while getting user record")
		return
	}

	var inboxRecord entity.Inbox
	if err := common.DBConn.Select("message", "created_at").Where("from_user_id = ? OR to_user_id = ?", userID, userID).Order("created_at DESC").First(&inboxRecord).Error; err != nil {
		errCh <- errors.New("error while getting inbox record")
		return
	}

	ch <- &model.InboxGetAllBubble{
		Username:    userRecord.Username,
		FullName:    userRecord.FullName,
		Avatar:      userRecord.Avatar,
		IsRead:      inboxRecord.IsRead,
		LastMessage: inboxRecord.Message,
		CreatedAt:   inboxRecord.CreatedAt.String(),
	}
}

func (ib *InboxService) InboxGetAllBubble(userID string) ([]*model.InboxGetAllBubble, error) {
	var inboxRecords []entity.Inbox
	if err := common.DBConn.Select("from_user_id", "to_user_id").Where("from_user_id = ? OR to_user_id = ?", userID, userID).Find(&inboxRecords).Error; err != nil {
		return nil, errors.New("error while getting inbox records")
	}

	// Get all unique user ID
	var listInboxID []string
	for _, inboxRecord := range inboxRecords {
		currentID := inboxRecord.FromUserID.String()
		if inboxRecord.FromUserID.String() == userID {
			currentID = inboxRecord.ToUserID.String()
		}
		listInboxID = append(listInboxID, currentID)
	}

	// Remove duplicate user ID
	uniqueStrings := make(map[string]bool)
	var filteredStrings []string
	for _, str := range listInboxID {
		if _, ok := uniqueStrings[str]; !ok {
			uniqueStrings[str] = true
			filteredStrings = append(filteredStrings, str)
		}
	}

	// Fetch user and inbox record concurrently
	var wg sync.WaitGroup
	ch := make(chan *model.InboxGetAllBubble, len(filteredStrings))
	errCh := make(chan error, len(filteredStrings))

	for _, userID := range filteredStrings {
		wg.Add(1)
		go fetchUserAndInboxRecord(userID, &wg, ch, errCh)
	}

	wg.Wait()
	close(ch)
	close(errCh)

	var inboxBubbleRecords []*model.InboxGetAllBubble
	for record := range ch {
		inboxBubbleRecords = append(inboxBubbleRecords, record)
	}

	if len(errCh) > 0 {
		for err := range errCh {
			return nil, errors.New(fmt.Sprintf("error: %s", err.Error()))
		}
	}

	return inboxBubbleRecords, nil
}
