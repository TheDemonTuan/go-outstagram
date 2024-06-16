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
	"outstagram/models/entity"
	"outstagram/models/req"
	"strings"
)

type PostService struct{}

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

func (p *PostService) PostGetAllByPostID(isOk bool, postID string, postRecords interface{}) error {
	if !isOk {
		if err := common.DBConn.Model(&entity.Post{}).Where("id = ? AND privacy IN ?", postID, []entity.PostPrivacy{entity.PostPublic}).Find(postRecords).Error; err != nil {
			return errors.New("error while querying post")
		}
		return nil
	}

	if err := common.DBConn.Where("id = ? AND privacy IN ?", postID, []entity.PostPrivacy{entity.PostOnlyFriend, entity.PostPublic}).Find(postRecords).Error; err != nil {
		return err
	}

	return nil
}

func (p *PostService) PostGetAll(posts *[]entity.Post) error {
	if err := common.DBConn.Find(&posts).Error; err != nil {
		return err
	}

	return nil
}

func (p *PostService) PostCreateValidateRequest(body *multipart.Form) (string, []*multipart.FileHeader, error) {
	if body == nil {
		return "", nil, errors.New("request body is required")
	}

	if body.Value == nil {
		return "", nil, errors.New("request body value is required")
	}

	if body.File == nil {
		return "", nil, errors.New("request body file is required")
	}

	caption := body.Value["caption"]
	files := body.File["files"]

	if len(caption) == 0 || len(caption[0]) == 0 {
		return "", nil, errors.New("caption is required")
	}

	if len(caption[0]) > 2200 {
		return "", nil, errors.New("caption is too long")
	}

	if len(files) == 0 {
		return "", nil, errors.New("files are required")
	}

	if len(files) > 10 {
		return "", nil, errors.New("files are too many")
	}

	acceptType := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/jpg":  true,
		"image/webp": true,
		"video/mp4":  true,
	}

	for _, file := range files {
		if !acceptType[file.Header["Content-Type"][0]] {
			return "", nil, errors.New(file.Filename + "file type is not supported")
		}
	}

	return caption[0], files, nil
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

func (p *PostService) PostCreateByUserID(userID uuid.UUID, caption string, localPaths []string, cloudinaryPaths []*uploader.UploadResult) (entity.Post, error) {
	newPostID := "TD" + common.RandomNString(18)
	newPost := entity.Post{
		ID:      newPostID,
		UserID:  userID,
		Caption: strings.Trim(caption, " "),
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

func (p *PostService) PostLikeByPostID(postID string, userID uuid.UUID) (entity.PostLike, error) {
	var postRecord entity.Post
	if err := common.DBConn.Where("id = ?", postID).First(&postRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.PostLike{}, fiber.NewError(fiber.StatusBadRequest, "Post not found")
		}
		return entity.PostLike{}, fiber.NewError(fiber.StatusInternalServerError, "Error while querying post")
	}

	var userRecord entity.User
	if err := common.DBConn.Where("id = ?", userID).First(&userRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.PostLike{}, fiber.NewError(fiber.StatusBadRequest, "User not found")
		}
		return entity.PostLike{}, fiber.NewError(fiber.StatusInternalServerError, "Error while querying user")
	}

	var postLike entity.PostLike
	if err := common.DBConn.Where("post_id = ? AND user_id = ?", postID, userID).First(&postLike).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			postLike = entity.PostLike{
				PostID: postID,
				UserID: userID,
			}

			if err := common.DBConn.Create(&postLike).Error; err != nil {
				return entity.PostLike{}, fiber.NewError(fiber.StatusInternalServerError, "Error while creating post like")
			}

			return postLike, nil
		}
		return entity.PostLike{}, fiber.NewError(fiber.StatusInternalServerError, "Error while querying post like")
	}

	if postLike.IsLiked {
		postLike.IsLiked = false
	} else {
		postLike.IsLiked = true
	}

	if err := common.DBConn.Save(&postLike).Error; err != nil {
		return entity.PostLike{}, fiber.NewError(fiber.StatusInternalServerError, "Error while updating post like")
	}

	return postLike, nil
}

func (p *PostService) PostEditByPostID(postID string, userID uuid.UUID, caption string) (entity.Post, error) {
	var postRecord entity.Post
	if err := common.DBConn.Where("id = ? and user_id = ?", postID, userID).First(&postRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Post{}, fiber.NewError(fiber.StatusBadRequest, "Post not found")
		}
		return entity.Post{}, fiber.NewError(fiber.StatusInternalServerError, "Error while querying post")
	}

	if caption == postRecord.Caption {
		return postRecord, nil
	}

	postRecord.Caption = strings.Trim(caption, " ")

	if err := common.DBConn.Save(&postRecord).Error; err != nil {
		return entity.Post{}, fiber.NewError(fiber.StatusInternalServerError, "Error while updating post")
	}

	return postRecord, nil
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

func (p *PostService) PostCommentByPostID(postID string, userID uuid.UUID, content string) (entity.PostComment, error) {
	var postRecord entity.Post
	if err := common.DBConn.Where("id = ?", postID).First(&postRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.PostComment{}, fiber.NewError(fiber.StatusBadRequest, "Post not found")
		}
		return entity.PostComment{}, fiber.NewError(fiber.StatusInternalServerError, "Error while querying post")
	}

	var userRecord entity.User
	if err := common.DBConn.Where("id = ?", userID).First(&userRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.PostComment{}, fiber.NewError(fiber.StatusBadRequest, "User not found")
		}
		return entity.PostComment{}, fiber.NewError(fiber.StatusInternalServerError, "Error while querying user")
	}

	newPostComment := entity.PostComment{
		PostID:  postID,
		UserID:  userID,
		Content: strings.Trim(content, " "),
	}

	if err := common.DBConn.Create(&newPostComment).Error; err != nil {
		return entity.PostComment{}, fiber.NewError(fiber.StatusInternalServerError, "Error while creating post comment")
	}

	return newPostComment, nil
}

func (p *PostService) PostGetAllCommentByPostID(postID string) ([]req.PostComment, error) {
	var commentsWithUsers []struct {
		entity.PostComment
		entity.User
	}

	err := common.DBConn.Table("post_comments").
		Select("post_comments.*, users.*").
		Joins("JOIN users ON post_comments.user_id = users.id").
		Where("post_comments.post_id = ?", postID).
		Find(&commentsWithUsers).Error
	if err != nil {
		return nil, err
	}

	var resultComments []req.PostComment
	for _, c := range commentsWithUsers {
		resultComments = append(resultComments, req.PostComment{
			ID:        c.PostComment.ID.String(),
			Content:   c.PostComment.Content,
			CreatedAt: c.PostComment.CreatedAt,
			User: req.UserComment{
				ID:       c.User.ID.String(),
				Username: c.User.Username,
				FullName: c.User.FullName,
				Avatar:   c.User.Avatar,
			},
		})
	}

	return resultComments, nil
}

func (p *PostService) PostGetHomePage(page int, currentUserID string, posts interface{}) error {
	var friendRecords []entity.Friend
	if err := common.DBConn.Where("(from_user_id = ? OR to_user_id = ?) AND status = ?", currentUserID, currentUserID, entity.FriendAccepted).Select("from_user_id", "to_user_id").Find(&friendRecords).Error; err != nil {
		return errors.New("error while querying followings")
	}

	var friends []string
	for _, f := range friendRecords {
		if f.FromUserID.String() == currentUserID {
			friends = append(friends, f.ToUserID.String())
		} else {
			friends = append(friends, f.FromUserID.String())
		}
	}

	friends = append(friends, currentUserID)

	postsPerPage := 2
	offset := (page - 1) * postsPerPage

	if err := common.DBConn.Model(&entity.Post{}).Where("(user_id IN ? AND privacy IN ?) OR user_id = ?", friends, []entity.PostPrivacy{entity.PostOnlyFriend, entity.PostPublic}, currentUserID).Order("created_at desc").Offset(offset).Limit(postsPerPage).Find(posts).Error; err != nil {
		return errors.New("error while querying posts")
	}

	return nil
}
