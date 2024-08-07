package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"outstagram/common"
	"outstagram/models/entity"
	"outstagram/models/req"
	"outstagram/services"
	"strconv"
	"strings"
	"time"
)

type PostController struct {
	postService   *services.PostService
	friendService *services.FriendService
	userService   *services.UserService
}

func NewPostController(postService *services.PostService, friendService *services.FriendService, userService *services.UserService) *PostController {
	return &PostController{
		postService:   postService,
		friendService: friendService,
		userService:   userService,
	}
}

func (p *PostController) PostMeGetAll(ctx *fiber.Ctx) error {
	rawUserID := ctx.Locals(common.UserIDLocalKey).(string)
	var postRecords []entity.Post
	if err := common.DBConn.Where("user_id = ?", rawUserID).Preload("PostFiles").Preload("PostLikes").Preload("PostComments").Order("created_at DESC").Find(&postRecords).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "No posts found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying posts")
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Posts found", postRecords)
}

func (p *PostController) PostMeCreate(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	caption, privacy, files, postType, isHideLike, isHideComment, err := p.postService.PostCreateValidateRequest(form)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	localPaths, cloudinaryPaths, err := p.postService.PostCreateUploadFiles(ctx, files)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	rawUserID := ctx.Locals(common.UserIDLocalKey).(string)
	userID := uuid.MustParse(rawUserID)

	privacyInt, err := strconv.Atoi(privacy)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	isHiddenLike, err := strconv.ParseBool(isHideLike)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	isHiddenComment, err := strconv.ParseBool(isHideComment)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	newPost, err := p.postService.PostCreateByUserID(userID, caption, entity.PostPrivacy(privacyInt), isHiddenLike, isHiddenComment, postType, localPaths, cloudinaryPaths)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.CreateResponse(ctx, fiber.StatusCreated, "Post created", newPost)
}

func (p *PostController) PostMeLikeByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")
	currentUserInfo, userInfoIsOk := ctx.Locals(common.UserInfoLocalKey).(entity.User)
	if !userInfoIsOk {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user")
	}

	var postRecord entity.Post
	postLike, postUserID, err := p.postService.PostLikeByPostID(postID, currentUserInfo.ID, &postRecord)
	if err != nil {
		return err
	}

	if !postLike.IsLiked {
		return common.CreateResponse(ctx, fiber.StatusOK, "Post unliked", postLike)
	}

	// Push notification
	if currentUserInfo.ID.String() != postUserID {
		data := map[string]interface{}{}
		data["type"] = "post-like"
		data["message"] = "liked your post!"
		data["userID"] = currentUserInfo.ID
		data["username"] = currentUserInfo.Username
		data["avatar"] = currentUserInfo.Avatar
		data["postID"] = postRecord.ID
		data["postType"] = entity.PostType(postRecord.Type).PostTypeString()
		data["postLike"] = postLike
		data["createdAt"] = time.Now().Format(time.RFC3339)

		if err := common.PusherClient.Trigger(postUserID, "internal-socket", data); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Post liked", postLike)
}

func (p *PostController) PostMeEditByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")
	rawUserID := ctx.Locals(common.UserIDLocalKey).(string)
	userID := uuid.MustParse(rawUserID)

	bodyData, err := common.RequestBodyValidator[req.PostMeEdit](ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	post, err := p.postService.PostEditByPostID(postID, userID, bodyData.Caption, bodyData.Privacy)
	if err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Post edited", post)
}

func (p *PostController) PostHiddenCommentByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")
	rawUserID := ctx.Locals(common.UserIDLocalKey).(string)
	userID := uuid.MustParse(rawUserID)

	isHiddenComment, err := p.postService.PostHiddenCommentByPostID(postID, userID)
	if err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Post edited", isHiddenComment)

}

func (p *PostController) PostHiddenLikeByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")
	rawUserID := ctx.Locals(common.UserIDLocalKey).(string)
	userID := uuid.MustParse(rawUserID)

	isHiddenLike, err := p.postService.PostHiddenLikeByPostID(postID, userID)
	if err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Post edited", isHiddenLike)

}

func (p *PostController) PostMeDeleteByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")
	rawUserID := ctx.Locals(common.UserIDLocalKey).(string)
	userID := uuid.MustParse(rawUserID)

	if err := p.postService.PostDeleteByPostID(postID, userID); err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusNoContent, "Post deleted", nil)
}

func (p *PostController) PostMeCommentByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")
	currentUserInfo, userInfoIsOk := ctx.Locals(common.UserInfoLocalKey).(entity.User)
	if !userInfoIsOk {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user")
	}

	parentID := ctx.Query("parentID")

	bodyData, err := common.RequestBodyValidator[req.PostMeComment](ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	content := strings.TrimSpace(bodyData.Content)

	if content == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Comment content is empty")
	}

	var postRecord entity.Post
	var userParentRecord entity.User
	postComment, postUserID, err := p.postService.PostCommentByPostID(postID, currentUserInfo.ID, content, parentID, &postRecord, &userParentRecord)
	if err != nil {
		return err
	}

	// Push notification
	data := map[string]interface{}{}
	data["type"] = "post-comment"
	data["message"] = "post commented: " + content
	data["userID"] = currentUserInfo.ID
	data["username"] = currentUserInfo.Username
	data["avatar"] = currentUserInfo.Avatar
	data["postID"] = postRecord.ID
	data["postType"] = entity.PostType(postRecord.Type).PostTypeString()
	data["postComment"] = postComment
	data["createdAt"] = time.Now().Format(time.RFC3339)

	if currentUserInfo.ID.String() != postUserID && userParentRecord.ID.String() != postUserID {
		if err := common.PusherClient.Trigger(postUserID, "internal-socket", data); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	if userParentRecord.ID.String() != "" && userParentRecord.ID.String() != currentUserInfo.ID.String() {
		data["message"] = "post replied: " + content
		data["parentID"] = userParentRecord.ID
		data["parentUsername"] = userParentRecord.Username

		if err := common.PusherClient.Trigger(userParentRecord.ID.String(), "internal-socket", data); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return common.CreateResponse(ctx, fiber.StatusCreated, "Comment created", postComment)
}

func (p *PostController) PostGetAllByUserID(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")

	var postRecords []entity.Post
	if err := common.DBConn.Where("user_id = ?", userID).Preload("PostFiles").Preload("PostLikes").Preload("PostComments").Find(&postRecords).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "No posts found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying posts")
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Posts found", postRecords)
}

func (p *PostController) PostGetByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")

	var post entity.Post
	if err := common.DBConn.Where("id = ?", postID).Preload("PostFiles").Preload("PostLikes").Preload("PostComments").First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Post not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying post")
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Post found", post)
}

func (p *PostController) PostMeSaveByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")
	currentUserInfo, userInfoIsOk := ctx.Locals(common.UserInfoLocalKey).(entity.User)
	if !userInfoIsOk {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user")
	}

	postSave, err := p.postService.PostSaveByPostID(postID, currentUserInfo.ID.String())
	if err != nil {
		return err
	}

	if postSave.ID == uuid.Nil {
		return common.CreateResponse(ctx, fiber.StatusOK, "Post remove save", nil)
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Post saved", postSave)
}

func (p *PostController) PostMeGetAllSaved(ctx *fiber.Ctx) error {
	rawUserID := ctx.Locals(common.UserIDLocalKey).(string)
	var postRecords []entity.Post
	if err := common.DBConn.Model(&entity.Post{}).Joins("JOIN post_saves ON post_saves.post_id = posts.id").Joins("JOIN users ON users.id = posts.user_id").Where("post_saves.user_id = ?", rawUserID).Where("posts.active = ?", true).
		Where("users.active = ?", true).Preload("PostFiles").Preload("PostLikes").Preload("PostComments").Find(&postRecords).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "No posts found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying posts")
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Posts found", postRecords)
}

func (p *PostController) PostDeleteCommentOnPostByCommentID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")

	commentID, err := uuid.Parse(ctx.Params("commentID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid comment ID")
	}

	var userID uuid.UUID
	if rawUserID := ctx.Query("userID"); rawUserID != "" {
		if userID, err = uuid.Parse(rawUserID); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
		}
	}

	currentUserInfo, userInfoIsOk := ctx.Locals(common.UserInfoLocalKey).(entity.User)
	if !userInfoIsOk {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user")
	}

	postRecord, err := p.postService.PostGetByPostID(postID)
	if err != nil {
		return err
	}

	deleteUserID := currentUserInfo.ID
	if postRecord.UserID == currentUserInfo.ID && userID != uuid.Nil {
		deleteUserID = userID
	}

	if err := p.postService.PostDeleteCommentOnPostByCommentID(postID, deleteUserID, commentID); err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Comment deleted", nil)
}

func (p *PostController) PostMeLikeCommentByCommentID(ctx *fiber.Ctx) error {
	rawCommentID := ctx.Params("commentID")
	commentID := uuid.MustParse(rawCommentID)

	currentUserInfo, userInfoIsOk := ctx.Locals(common.UserInfoLocalKey).(entity.User)
	if !userInfoIsOk {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user")
	}

	var commentRecord entity.PostComment
	postCommentLike, commentUserID, err := p.postService.PostLikeCommentByCommentID(commentID, currentUserInfo.ID, &commentRecord)
	if err != nil {
		return err
	}

	if !postCommentLike.IsCommentLiked {
		return common.CreateResponse(ctx, fiber.StatusOK, "Comment unliked", postCommentLike)
	}

	// Push notification
	if currentUserInfo.ID.String() != commentUserID {
		data := map[string]interface{}{}
		data["type"] = "comment-like"
		data["message"] = "liked your comment!"
		data["userID"] = currentUserInfo.ID
		data["username"] = currentUserInfo.Username
		data["avatar"] = currentUserInfo.Avatar
		data["commentID"] = commentRecord.ID
		//data["postType"] = entity.PostType(commentRecord.Type).PostTypeString()
		data["CommentLike"] = postCommentLike
		data["createdAt"] = time.Now().Format(time.RFC3339)

		if err := common.PusherClient.Trigger(commentUserID, "internal-socket", data); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Comment liked", postCommentLike)

}

func (p *PostController) PostMeGetAllDeletedByUserID(ctx *fiber.Ctx) error {
	rawUserID := ctx.Locals(common.UserIDLocalKey).(string)
	var postRecords []entity.Post

	if err := common.DBConn.Unscoped().Model(&entity.Post{}).Where("user_id = ? AND deleted_at IS NOT NULL", rawUserID).Preload("PostFiles").Preload("PostLikes").Preload("PostComments").Find(&postRecords).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "No deleted posts found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying posts")
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Deleted posts found", postRecords)
}

func (p *PostController) PostMeRestoreByPostID(ctx *fiber.Ctx) error {

	rawUserID := ctx.Locals(common.UserIDLocalKey).(string)
	bodyData, err := common.RequestBodyValidator[req.PostMeRestore](ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	postRestore, err := p.postService.PostMeRestoreByPostID(rawUserID, bodyData.PostIDs)
	if err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Post restored", postRestore)
}

func (p *PostController) PostMeGetAllLiked(ctx *fiber.Ctx) error {
	rawUserID := ctx.Locals(common.UserIDLocalKey).(string)
	var postRecords []entity.Post
	if err := common.DBConn.Model(&entity.Post{}).Joins("JOIN post_likes ON post_likes.post_id = posts.id").Joins("JOIN users ON users.id = posts.user_id").Where("post_likes.user_id = ?", rawUserID).Where("posts.active = ?", true).
		Where("users.active = ?", true).Preload("PostFiles").Preload("PostLikes").Preload("PostComments").Find(&postRecords).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "No posts found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying posts")
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Posts found", postRecords)

}

func (p *PostController) PostMeGetAllCommented(ctx *fiber.Ctx) error {
	rawUserID := ctx.Locals(common.UserIDLocalKey).(string)
	var postRecords []entity.Post

	// Fetch posts where the user has commented, ensuring unique posts
	if err := common.DBConn.Model(&entity.Post{}).
		Joins("JOIN post_comments ON post_comments.post_id = posts.id").
		Joins("JOIN users ON users.id = posts.user_id").
		Where("post_comments.user_id = ?", rawUserID).
		Where("posts.active = ?", true).
		Where("users.active = ?", true).
		Preload("PostFiles").
		Preload("PostLikes").
		Preload("PostComments").
		Find(&postRecords).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "No posts found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying posts")
	}

	// Use a map to ensure unique posts
	uniquePostMap := make(map[string]entity.Post)
	for _, post := range postRecords {
		uniquePostMap[post.ID] = post
	}

	uniquePosts := make([]entity.Post, 0, len(uniquePostMap))
	userIDs := make([]uuid.UUID, 0, len(uniquePostMap))
	for _, post := range uniquePostMap {
		uniquePosts = append(uniquePosts, post)
		userIDs = append(userIDs, post.UserID)
	}

	// Fetch users associated with the posts
	var users []entity.User
	if err := common.DBConn.Where("id IN (?)", userIDs).Find(&users).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying users")
	}

	// Create a map of users for easy lookup
	userMap := make(map[uuid.UUID]entity.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// Create response objects
	postResponses := make([]req.PostResponse, len(uniquePosts))
	for i, post := range uniquePosts {
		postResponses[i] = req.PostResponse{
			Post: post,
			User: userMap[post.UserID],
		}
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Posts found", postResponses)
}
