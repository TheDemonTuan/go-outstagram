package services

import (
	"errors"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
	"os"
	"outstagram/common"
	"outstagram/graph/model"
	"outstagram/models/entity"
	"strings"
)

type PostService struct {
	friendService *FriendService
}

func NewPostService() *PostService {
	return &PostService{}
}

func (p *PostService) PostGetByPostID(postID string) (entity.Post, error) {
	var post entity.Post
	if err := common.DBConn.Where("id = ?", postID).First(&post).Error; err != nil {
		return entity.Post{}, err
	}

	return post, nil
}

func (p *PostService) PostGetAllByUserID(userID string, posts *[]entity.Post) error {
	if err := common.DBConn.Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return err
	}

	return nil
}

func (p *PostService) PostGetAllByUserName(username string, posts interface{}) error {
	var user entity.User
	if err := common.DBConn.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("error while querying user")
	}

	if err := common.DBConn.Model(&entity.Post{}).Where("user_id = ?", user.ID).Find(posts).Error; err != nil {
		return errors.New("error while querying post")
	}

	return nil
}

func (p *PostService) PostProfileGetAllByUserName(isOK bool, currentUserID, username string, postType entity.PostType, posts interface{}) error {
	var user entity.User
	if err := common.DBConn.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("error while querying user")
	}

	if isOK {
		if currentUserID == user.ID.String() {
			if err := common.DBConn.Model(&entity.Post{}).Where("user_id = ? AND type = ?", user.ID, postType).Order("created_at desc").Find(posts).Error; err != nil {
				return errors.New("error while querying post")
			}
			return nil
		}

		isFriend := true
		var friendRecord entity.Friend
		if err := p.friendService.GetFriendByUserID(&friendRecord, currentUserID, user.ID.String()); err != nil {
			isFriend = false
		}

		if isFriend {
			if err := common.DBConn.Model(&entity.Post{}).Where("user_id = ? AND privacy IN ? AND type = ?", user.ID.String(), []entity.PostPrivacy{entity.PostOnlyFriend, entity.PostPublic}, postType).Order("created_at desc").Find(posts).Error; err != nil {
				return errors.New("error while querying posts")
			}
			return nil
		}
	}

	if err := common.DBConn.Model(&entity.Post{}).Where("user_id = ? AND privacy IN ? AND type = ?", user.ID.String(), []entity.PostPrivacy{entity.PostPublic}, postType).Order("created_at desc").Find(posts).Error; err != nil {
		return errors.New("error while querying post")
	}

	return nil
}

func (p *PostService) PostByPostID(isOk bool, currentUserID, postID string, postRecords interface{}) error {
	if err := common.DBConn.Model(entity.Post{}).Where("id = ?", postID).First(postRecords).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return errors.New("error while querying post")
	}

	post := postRecords.(**model.Post)

	if isOk {
		if currentUserID == (*post).UserID {
			return nil
		}

		isFriend := true
		var friendRecord entity.Friend
		if err := p.friendService.GetFriendByUserID(&friendRecord, currentUserID, (*post).UserID); err != nil {
			isFriend = false
		}

		if isFriend {
			if (*post).Privacy == entity.PostOnlyFriend || (*post).Privacy == entity.PostPublic {
				return nil
			}
		}
	}

	if (*post).Privacy == entity.PostPublic {
		return nil
	}

	return errors.New("post is private")
}

func (p *PostService) PostGetSuggestions(isOK bool, currentUserID, skipPostID string, limit int, postRecords interface{}) error {
	var post entity.Post
	if err := common.DBConn.Where("id = ?", skipPostID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("skip post not found")
		}
		return errors.New("skip error while querying post")

	}
	if isOK {
		if currentUserID == post.UserID.String() {
			if err := common.DBConn.Model(&entity.Post{}).Not("id = ?", skipPostID).Where("user_id = ?", post.UserID).Order("created_at desc").Limit(limit).Find(postRecords).Error; err != nil {
				return errors.New("error while querying post")
			}
			return nil
		}

		isFriend := true
		var friendRecord entity.Friend
		if err := p.friendService.GetFriendByUserID(&friendRecord, currentUserID, post.UserID.String()); err != nil {
			isFriend = false
		}

		if isFriend {
			if err := common.DBConn.Model(&entity.Post{}).Not("id = ?", skipPostID).Where("user_id = ? AND privacy IN ?", post.UserID.String(), []entity.PostPrivacy{entity.PostOnlyFriend, entity.PostPublic}).Order("created_at desc").Limit(limit).Find(postRecords).Error; err != nil {
				return errors.New("error while querying posts")
			}
			return nil
		}
	}

	if err := common.DBConn.Model(&entity.Post{}).Not("id = ?", skipPostID).Where("user_id = ? AND privacy IN ?", post.UserID.String(), []entity.PostPrivacy{entity.PostPublic}).Order("created_at desc").Limit(limit).Find(postRecords).Error; err != nil {
		return errors.New("error while querying post")
	}

	return nil
}

func (p *PostService) PostGetAll(posts *[]entity.Post) error {
	if err := common.DBConn.Find(&posts).Error; err != nil {
		return err
	}

	return nil
}

func (p *PostService) PostCreateValidateRequest(body *multipart.Form) (string, string, []*multipart.FileHeader, entity.PostType, string, string, error) {
	if body == nil {
		return "", "", nil, entity.PostNormal, "", "", errors.New("request body is required")
	}

	if body.Value == nil {
		return "", "", nil, entity.PostNormal, "", "", errors.New("request body value is required")
	}

	if body.File == nil {
		return "", "", nil, entity.PostNormal, "", "", errors.New("request body file is required")
	}

	caption := body.Value["caption"]
	privacy := body.Value["privacy"]
	files := body.File["files"]
	isHideLike := body.Value["is_hide_like"]
	isHideComment := body.Value["is_hide_comment"]

	if len(caption) == 0 || len(caption[0]) == 0 {
		return "", "", nil, entity.PostNormal, "", "", errors.New("caption is required")
	}

	if len(caption[0]) > 2200 {
		return "", "", nil, entity.PostNormal, "", "", errors.New("caption is too long")
	}

	if len(privacy) == 0 || len(privacy[0]) == 0 {
		return "", "", nil, entity.PostNormal, "", "", errors.New("privacy is required")

	}

	if privacy[0] != "0" && privacy[0] != "1" && privacy[0] != "2" {
		return "", "", nil, entity.PostNormal, "", "", errors.New("privacy is invalid")

	}

	if len(isHideLike) == 0 || len(isHideLike[0]) == 0 {
		return "", "", nil, entity.PostNormal, "", "", errors.New("is_hide_like is required")

	}

	if isHideLike[0] != "true" && isHideLike[0] != "false" {
		return "", "", nil, entity.PostNormal, "", "", errors.New("is_hide_like is invalid")

	}

	if len(isHideComment) == 0 || len(isHideComment[0]) == 0 {
		return "", "", nil, entity.PostNormal, "", "", errors.New("is_hide_comment is required")

	}

	if len(files) == 0 {
		return "", "", nil, entity.PostNormal, "", "", errors.New("files are required")
	}

	if len(files) > 10 {
		return "", "", nil, entity.PostNormal, "", "", errors.New("files are too many")
	}

	validImageContentTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/jpg":  true,
		"image/webp": true,
	}

	validVideoContentTypes := map[string]bool{
		"video/mp4": true,
	}

	totalImage := 0
	totalVideo := 0

	for _, file := range files {
		if validImageContentTypes[file.Header["Content-Type"][0]] {
			totalImage++
		}

		if validVideoContentTypes[file.Header["Content-Type"][0]] {
			totalVideo++
		}

		if !validImageContentTypes[file.Header["Content-Type"][0]] && !validVideoContentTypes[file.Header["Content-Type"][0]] {
			return "", "", nil, entity.PostNormal, "", "", errors.New("file content type is invalid")
		}

		if totalImage > 0 && totalVideo > 0 {
			return "", "", nil, entity.PostNormal, "", "", errors.New("files must be all images or all videos")
		}
	}

	if totalImage == 0 && totalVideo == 0 {
		return "", "", nil, entity.PostNormal, "", "", errors.New("files must be all images or all videos")
	}

	var postType entity.PostType
	if totalImage > 0 {
		postType = entity.PostNormal
	}

	if totalVideo > 0 {
		postType = entity.PostReel
	}

	return caption[0], privacy[0], files, postType, isHideLike[0], isHideComment[0], nil
}

func (p *PostService) PostCreateUploadFiles(ctx *fiber.Ctx, files []*multipart.FileHeader) ([]string, []*uploader.UploadResult, error) {
	var localPaths []string
	var cloudinaryPaths []*uploader.UploadResult
	for _, file := range files {
		ext := strings.Split(file.Header["Content-Type"][0], "/")[1]
		randName := common.RandomNString(30)
		newImageName := fmt.Sprintf("%s.%s", randName, ext)

		//Save to local
		if err := ctx.SaveFile(file, newImageName); err != nil {
			if err := p.PostCreateDeleteFiles(localPaths, cloudinaryPaths, true); err != nil {
				return nil, nil, err
			}
			return nil, nil, errors.New("error while saving file " + file.Filename)
		}

		isImage := file.Header["Content-Type"][0][:5] == "image"

		params := uploader.UploadParams{
			DisplayName: randName,
			Folder:      "posts",
		}

		if isImage {
			params = uploader.UploadParams{
				Format:      "webp",
				DisplayName: randName,
				Folder:      "posts",
			}
		}

		//Upload to cloudinary
		result, err := common.CloudinaryUploadFile(newImageName, params)
		if err != nil {
			if err := p.PostCreateDeleteFiles(localPaths, cloudinaryPaths, true); err != nil {
				return nil, nil, err
			}
			return nil, nil, errors.New("error while uploading file " + file.Filename)
		}

		cloudinaryPaths = append(cloudinaryPaths, result)
		localPaths = append(localPaths, newImageName)
	}

	if err := p.PostCreateDeleteFiles(localPaths, cloudinaryPaths, false); err != nil {
		return nil, nil, err
	}

	return localPaths, cloudinaryPaths, nil
}

func (p *PostService) PostCreateDeleteFiles(localPaths []string, cloudinaryPaths []*uploader.UploadResult, isDeleteCloud bool) error {
	isErr := false
	for _, filePath := range localPaths {
		if err := os.Remove(filePath); err != nil {
			isErr = true
		}
	}

	if isDeleteCloud || isErr {
		for _, cloudinaryPath := range cloudinaryPaths {
			if err := common.CloudinaryDeleteFile(cloudinaryPath.PublicID); err != nil {
				return err
			}
		}
	}

	if isErr {
		return errors.New("error while deleting files")
	}

	return nil
}

func (p *PostService) PostCreateByUserID(userID uuid.UUID, caption string, privacy entity.PostPrivacy, isHideLike bool, isHideComment bool, postType entity.PostType, localPaths []string, cloudinaryPaths []*uploader.UploadResult) (entity.Post, error) {

	newPostID := "TD" + common.RandomNString(18)
	newPost := entity.Post{
		ID:            newPostID,
		UserID:        userID,
		Caption:       strings.Trim(caption, " "),
		Privacy:       privacy,
		IsHideComment: isHideComment,
		IsHideLike:    isHideLike,
		Type:          postType,
	}

	for _, filePath := range cloudinaryPaths {
		newPost.PostFiles = append(newPost.PostFiles, entity.PostFile{
			ID:     filePath.PublicID,
			PostID: newPostID,
			URL:    filePath.SecureURL,
		})
	}

	if err := common.DBConn.Create(&newPost).Error; err != nil {
		if err := p.PostCreateDeleteFiles(localPaths, cloudinaryPaths, true); err != nil {
			return entity.Post{}, err
		}
		return entity.Post{}, errors.New("error while creating post")
	}

	return newPost, nil
}

func (p *PostService) PostLikeByPostID(postID string, userID uuid.UUID, postRecord *entity.Post) (entity.PostLike, string, error) {
	if err := common.DBConn.Where("id = ?", postID).First(&postRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.PostLike{}, "", fiber.NewError(fiber.StatusBadRequest, "Post not found")
		}
		return entity.PostLike{}, "", fiber.NewError(fiber.StatusInternalServerError, "Error while querying post")
	}

	var userRecord entity.User
	if err := common.DBConn.Where("id = ?", userID).First(&userRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.PostLike{}, "", fiber.NewError(fiber.StatusBadRequest, "User not found")
		}
		return entity.PostLike{}, "", fiber.NewError(fiber.StatusInternalServerError, "Error while querying user")
	}

	var postLike entity.PostLike
	if err := common.DBConn.Where("post_id = ? AND user_id = ?", postID, userID).First(&postLike).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			postLike = entity.PostLike{
				PostID: postID,
				UserID: userID,
			}

			if err := common.DBConn.Create(&postLike).Error; err != nil {
				return entity.PostLike{}, "", fiber.NewError(fiber.StatusInternalServerError, "Error while creating post like")
			}

			return postLike, postRecord.UserID.String(), nil
		}
		return entity.PostLike{}, "", fiber.NewError(fiber.StatusInternalServerError, "Error while querying post like")
	}

	if postLike.IsLiked {
		postLike.IsLiked = false
	} else {
		postLike.IsLiked = true
	}

	if err := common.DBConn.Save(&postLike).Error; err != nil {
		return entity.PostLike{}, "", fiber.NewError(fiber.StatusInternalServerError, "Error while updating post like")
	}

	return postLike, postRecord.UserID.String(), nil
}

func (p *PostService) PostEditByPostID(postID string, userID uuid.UUID, caption string, privacy entity.PostPrivacy) (entity.Post, error) {
	var postRecord entity.Post
	if err := common.DBConn.Where("id = ? and user_id = ?", postID, userID).First(&postRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Post{}, fiber.NewError(fiber.StatusBadRequest, "Post not found")
		}
		return entity.Post{}, fiber.NewError(fiber.StatusInternalServerError, "Error while querying post")
	}
	if postRecord.Caption == caption && postRecord.Privacy == privacy {
		return entity.Post{}, fiber.NewError(fiber.StatusBadRequest, "Nothing to update")
	}

	postRecord.Caption = strings.Trim(caption, " ")
	postRecord.Privacy = privacy

	if err := common.DBConn.Save(&postRecord).Error; err != nil {
		return entity.Post{}, fiber.NewError(fiber.StatusInternalServerError, "Error while updating post")
	}

	return postRecord, nil
}

func (p *PostService) PostHiddenCommentByPostID(postID string, userID uuid.UUID) (bool, error) {
	var postRecord entity.Post
	if err := common.DBConn.Where("id = ? and user_id = ?", postID, userID).First(&postRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fiber.NewError(fiber.StatusBadRequest, "Post not found")
		}
		return false, fiber.NewError(fiber.StatusInternalServerError, "Error while querying post")

	}
	postRecord.IsHideComment = !postRecord.IsHideComment

	if err := common.DBConn.Save(&postRecord).Error; err != nil {
		return false, fiber.NewError(fiber.StatusInternalServerError, "Error while updating post")
	}

	return postRecord.IsHideComment, nil
}

func (p *PostService) PostHiddenLikeByPostID(postID string, userID uuid.UUID) (bool, error) {
	var postRecord entity.Post
	if err := common.DBConn.Where("id = ? and user_id = ?", postID, userID).First(&postRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fiber.NewError(fiber.StatusBadRequest, "Post not found")
		}
		return false, fiber.NewError(fiber.StatusInternalServerError, "Error while querying post")

	}
	postRecord.IsHideLike = !postRecord.IsHideLike

	if err := common.DBConn.Save(&postRecord).Error; err != nil {
		return false, fiber.NewError(fiber.StatusInternalServerError, "Error while updating post")
	}

	return postRecord.IsHideLike, nil
}

func (p *PostService) PostDeleteByPostID(postID string, userID uuid.UUID) error {
	var postRecord entity.Post
	if err := common.DBConn.Where("id = ? and user_id = ?", postID, userID).First(&postRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusBadRequest, "Post not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying post")
	}

	if err := common.DBConn.Delete(&postRecord).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while deleting post")
	}

	return nil
}

func (p *PostService) PostCommentByPostID(postID string, userID uuid.UUID, content string, parentID string, postRecord *entity.Post, userParentRecord *entity.User) (entity.PostComment, string, error) {
	if err := common.DBConn.Where("id = ?", postID).First(&postRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.PostComment{}, "", fiber.NewError(fiber.StatusBadRequest, "Post not found")
		}
		return entity.PostComment{}, "", fiber.NewError(fiber.StatusInternalServerError, "Error while querying post")
	}

	newPostComment := entity.PostComment{
		PostID:  postID,
		UserID:  userID,
		Content: strings.Trim(content, " "),
	}

	if parentID != "" {
		parentIDUUID, err := uuid.Parse(parentID)
		if err != nil {
			return entity.PostComment{}, "", fiber.NewError(fiber.StatusBadRequest, "Parent comment ID is invalid")
		}

		var postCommentRecord entity.PostComment
		if err := common.DBConn.Where("id = ?", parentID).First(&postCommentRecord).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return entity.PostComment{}, "", fiber.NewError(fiber.StatusBadRequest, "Parent comment not found")
			}
			return entity.PostComment{}, "", fiber.NewError(fiber.StatusInternalServerError, "Error while querying parent comment")
		}

		if err := common.DBConn.Select("id", "username").Where("id = ?", postCommentRecord.UserID).First(&userParentRecord).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return entity.PostComment{}, "", fiber.NewError(fiber.StatusBadRequest, "Parent comment user not found")
			}
			return entity.PostComment{}, "", fiber.NewError(fiber.StatusInternalServerError, "Error while querying parent comment user")
		}

		newPostComment.ParentID = parentIDUUID
	}

	if err := common.DBConn.Create(&newPostComment).Error; err != nil {
		return entity.PostComment{}, "", fiber.NewError(fiber.StatusInternalServerError, "Error while creating post comment")
	}

	return newPostComment, postRecord.UserID.String(), nil
}

func (p *PostService) PostGetHomePage(page int, currentUserID string, postType entity.PostType, posts interface{}) error {
	friends, err := p.friendService.GetListFriendID(currentUserID)
	if err != nil {
		return err
	}

	postsPerPage := 2
	offset := (page - 1) * postsPerPage

	if err := common.DBConn.Model(&entity.Post{}).Where("((user_id IN ? AND privacy IN ?) OR user_id = ?) AND type = ?", friends, []entity.PostPrivacy{entity.PostOnlyFriend, entity.PostPublic}, currentUserID, postType).Order("created_at desc").Offset(offset).Limit(postsPerPage).Find(posts).Error; err != nil {
		return errors.New("error while querying posts")
	}

	return nil
}

func (p *PostService) PostGetExplores(page int, posts interface{}) error {
	if err := common.DBConn.Model(&entity.Post{}).Where("privacy IN ?", []entity.PostPrivacy{entity.PostOnlyFriend, entity.PostPublic}).Order("created_at desc").Find(posts).Error; err != nil {
		return errors.New("error while querying posts")
	}

	return nil
}
