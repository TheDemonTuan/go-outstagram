package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"
	"net/http"
	"outstagram/common"
	"outstagram/graph/model"
	"outstagram/models/entity"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

// FromUserInfo is the resolver for the from_user_info field.
func (r *friendResolver) FromUserInfo(ctx context.Context, obj *model.Friend) (*model.User, error) {
	var userRecord *model.User
	if err := r.userService.UserGetByID(obj.FromUserID, &userRecord); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	if userRecord.ID == "" {
		return nil, gqlerror.Errorf("user not found")
	}

	return userRecord, nil
}

// ToUserInfo is the resolver for the to_user_info field.
func (r *friendResolver) ToUserInfo(ctx context.Context, obj *model.Friend) (*model.User, error) {
	var userRecord *model.User
	if err := r.userService.UserGetByID(obj.ToUserID, &userRecord); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	if userRecord.ID == "" {
		return nil, gqlerror.Errorf("user not found")
	}

	return userRecord, nil
}

// Files is the resolver for the files field.
func (r *inboxResolver) Files(ctx context.Context, obj *model.Inbox) ([]*model.InboxFile, error) {
	panic(fmt.Errorf("not implemented: Files - files"))
}

// FromUserInfo is the resolver for the from_user_info field.
func (r *inboxResolver) FromUserInfo(ctx context.Context, obj *model.Inbox) (*model.User, error) {
	var userRecord *model.User
	if err := r.userService.UserGetByID(obj.FromUserID, &userRecord); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	if userRecord.ID == "" {
		return nil, gqlerror.Errorf("user not found")
	}

	return userRecord, nil
}

// ToUserInfo is the resolver for the to_user_info field.
func (r *inboxResolver) ToUserInfo(ctx context.Context, obj *model.Inbox) (*model.User, error) {
	var userRecord *model.User
	if err := r.userService.UserGetByID(obj.ToUserID, &userRecord); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	if userRecord.ID == "" {
		return nil, gqlerror.Errorf("user not found")
	}

	return userRecord, nil
}

// Privacy is the resolver for the privacy field.
func (r *postResolver) Privacy(ctx context.Context, obj *model.Post) (*int, error) {
	var privacy = int(obj.Privacy)
	return &privacy, nil
}

// Type is the resolver for the type field.
func (r *postResolver) Type(ctx context.Context, obj *model.Post) (*int, error) {
	var postType = int(obj.Type)
	return &postType, nil
}

// User is the resolver for the user field.
func (r *postResolver) User(ctx context.Context, obj *model.Post) (*model.User, error) {
	var userRecord *model.User
	if err := r.userService.UserGetByID(obj.UserID, &userRecord); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return userRecord, nil
}

// PostFiles is the resolver for the post_files field.
func (r *postResolver) PostFiles(ctx context.Context, obj *model.Post) ([]*model.PostFile, error) {
	var postFileRecords []*model.PostFile
	if err := r.postFileService.PostFileGetAllByPostID(obj.ID, &postFileRecords); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return postFileRecords, nil
}

// PostLikes is the resolver for the post_likes field.
func (r *postResolver) PostLikes(ctx context.Context, obj *model.Post) ([]*model.PostLike, error) {
	var postLikeRecords []*model.PostLike
	if err := r.postLikeService.PostLikeGetAllByPostID(obj.ID, &postLikeRecords); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return postLikeRecords, nil
}

// PostComments is the resolver for the post_comments field.
func (r *postResolver) PostComments(ctx context.Context, obj *model.Post) ([]*model.PostComment, error) {
	var postCommentRecords []*model.PostComment
	if err := r.postCommentService.PostCommentGetAllByPostID(obj.ID, &postCommentRecords); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return postCommentRecords, nil
}

// PostSaves is the resolver for the post_saves field.
func (r *postResolver) PostSaves(ctx context.Context, obj *model.Post) ([]*model.PostSave, error) {
	panic(fmt.Errorf("not implemented: PostSaves - post_saves"))
}

// PostCommentLikes is the resolver for the post_comment_likes field.
func (r *postResolver) PostCommentLikes(ctx context.Context, obj *model.Post) ([]*model.CommentLike, error) {
	panic(fmt.Errorf("not implemented: PostCommentLikes - post_comment_likes"))
}

// User is the resolver for the user field.
func (r *postCommentResolver) User(ctx context.Context, obj *model.PostComment) (*model.User, error) {
	var userRecord *model.User
	if err := r.userService.UserGetByID(obj.UserID, &userRecord); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return userRecord, nil
}

// Parent is the resolver for the parent field.
func (r *postCommentResolver) Parent(ctx context.Context, obj *model.PostComment) (*model.PostComment, error) {
	var postCommentRecord *model.PostComment
	if (obj.ParentID == "") || (obj.ParentID == "0") {
		return postCommentRecord, nil
	}
	if err := r.postCommentService.PostCommentGetByID(obj.ParentID, &postCommentRecord); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return postCommentRecord, nil
}

// User is the resolver for the user field.
func (r *postLikeResolver) User(ctx context.Context, obj *model.PostLike) (*model.User, error) {
	var userRecord *model.User
	if err := r.userService.UserGetByID(obj.UserID, &userRecord); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return userRecord, nil
}

// UserByUsername is the resolver for the userByUsername field.
func (r *queryResolver) UserByUsername(ctx context.Context, username string) (*model.User, error) {
	var userRecord *model.User
	if err := r.userService.UserGetByUserName(username, &userRecord); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	if userRecord.ID == "" {
		return nil, gqlerror.Errorf("user not found")
	}

	return userRecord, nil
}

// UserProfile is the resolver for the userProfile field.
func (r *queryResolver) UserProfile(ctx context.Context, username string) (*model.UserProfile, error) {
	var userProfile = &model.UserProfile{
		Username: username,
		Posts:    []*model.Post{},
		Reels:    []*model.Post{},
	}

	return userProfile, nil
}

// UserSearch is the resolver for the userSearch field.
func (r *queryResolver) UserSearch(ctx context.Context, keyword string) ([]*model.UserSearch, error) {
	if keyword == "" {
		return nil, nil
	}

	var userSearchRecords []*model.UserSearch
	if err := r.userService.UserSearchByUsernameOrFullName(keyword, &userSearchRecords); err != nil {
		return nil, err
	}

	return userSearchRecords, nil
}

// UserSuggestion the resolver for the userSuggestion field.
func (r *queryResolver) UserSuggestion(ctx context.Context, count int) ([]*model.UserSuggestion, error) {
	currentUserID, isOk := ctx.Value(common.UserIDLocalKey).(string)

	if !isOk {
		return nil, gqlerror.Errorf("user not found")
	}

	var userSuggestions []*model.UserSuggestion
	if err := r.userService.UserSuggestion(currentUserID, count, &userSuggestions); err != nil {
		return nil, err
	}

	return userSuggestions, nil
}

// PostByUsername is the resolver for the postByUsername field.
func (r *queryResolver) PostByUsername(ctx context.Context, username string) ([]*model.Post, error) {
	var posts []*model.Post
	if err := r.postService.PostGetAllByUserName(username, &posts); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return posts, nil
}

// PostByPostID is the resolver for the postByPostId field.
func (r *queryResolver) PostByPostID(ctx context.Context, postID string) (*model.Post, error) {
	currentUserID, isOk := ctx.Value(common.UserIDLocalKey).(string)

	var post *model.Post
	if err := r.postService.PostByPostID(isOk, currentUserID, postID, &post); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return post, nil
}

// PostSuggestions is the resolver for the postSuggestions field.
func (r *queryResolver) PostSuggestions(ctx context.Context, skipPostID string, limit int) ([]*model.Post, error) {
	currentUserID, isOk := ctx.Value(common.UserIDLocalKey).(string)

	var posts []*model.Post
	if err := r.postService.PostGetSuggestions(isOk, currentUserID, skipPostID, limit, &posts); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return posts, nil
}

// PostHomePage is the resolver for the postHomePage field.
func (r *queryResolver) PostHomePage(ctx context.Context, page int) ([]*model.Post, error) {
	currentUserID, isOk := ctx.Value(common.UserIDLocalKey).(string)
	if !isOk {
		return nil, &gqlerror.Error{
			Message: "User not found",
			Extensions: map[string]interface{}{
				"code": http.StatusUnauthorized,
			},
		}
	}

	var posts []*model.Post
	if err := r.postService.PostGetHomePage(page, currentUserID, entity.PostNormal, &posts); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return posts, nil
}

// PostReel is the resolver for the postReel field.
func (r *queryResolver) PostReel(ctx context.Context, page int) ([]*model.Post, error) {
	currentUserID, isOk := ctx.Value(common.UserIDLocalKey).(string)
	if !isOk {
		return nil, gqlerror.Errorf("user not found")
	}

	var posts []*model.Post
	if err := r.postService.PostGetHomePage(page, currentUserID, entity.PostReel, &posts); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return posts, nil
}

// PostExplores is the resolver for the postExplores field.
func (r *queryResolver) PostExplores(ctx context.Context, page int) ([]*model.Post, error) {
	var posts []*model.Post
	if err := r.postService.PostGetExplores(page, &posts); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return posts, nil
}

// InboxGetByUsername is the resolver for the inboxGetByUsername field.
func (r *queryResolver) InboxGetByUsername(ctx context.Context, username string) ([]*model.Inbox, error) {
	currentUserID, isOk := ctx.Value(common.UserIDLocalKey).(string)
	if !isOk {
		return nil, gqlerror.Errorf("user not found")
	}

	var inboxRecords []*model.Inbox
	if err := r.inboxService.InboxGetAllByUserName(currentUserID, username, &inboxRecords); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return inboxRecords, nil
}

// InboxGetAllBubble is the resolver for the inboxGetAllBubble field.
func (r *queryResolver) InboxGetAllBubble(ctx context.Context) ([]*model.InboxGetAllBubble, error) {
	currentUserID, isOk := ctx.Value(common.UserIDLocalKey).(string)
	if !isOk {
		return nil, gqlerror.Errorf("user not found")
	}

	inboxBubbleRecords, err := r.inboxService.InboxGetAllBubble(currentUserID)
	if err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return inboxBubbleRecords, nil
}

// Friends is the resolver for the friends field.
func (r *userResolver) Friends(ctx context.Context, obj *model.User) ([]*model.Friend, error) {
	var friends []*model.Friend
	if err := r.friendService.GetAllFriendsByUserName(obj.Username, &friends); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return friends, nil
}

// User is the resolver for the user field.
func (r *userProfileResolver) User(ctx context.Context, obj *model.UserProfile) (*model.User, error) {
	var userRecord *model.User
	if err := r.userService.UserGetByUserName(obj.Username, &userRecord); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return userRecord, nil
}

// Posts is the resolver for the posts field.
func (r *userProfileResolver) Posts(ctx context.Context, obj *model.UserProfile) ([]*model.Post, error) {
	currentUserID, isOk := ctx.Value(common.UserIDLocalKey).(string)

	var posts []*model.Post
	if err := r.postService.PostProfileGetAllByUserName(isOk, currentUserID, obj.Username, entity.PostNormal, &posts); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return posts, nil
}

// Reels is the resolver for the reels field.
func (r *userProfileResolver) Reels(ctx context.Context, obj *model.UserProfile) ([]*model.Post, error) {
	currentUserID, isOk := ctx.Value(common.UserIDLocalKey).(string)

	var posts []*model.Post
	if err := r.postService.PostProfileGetAllByUserName(isOk, currentUserID, obj.Username, entity.PostReel, &posts); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return posts, nil
}

// Posts is the resolver for the posts field.
func (r *userSuggestionResolver) Posts(ctx context.Context, obj *model.UserSuggestion) ([]*model.Post, error) {
	var posts []*model.Post
	if err := r.postService.PostGetAllByUserName(obj.Username, &posts); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return posts, nil
}

// Friends is the resolver for the friends field.
func (r *userSuggestionResolver) Friends(ctx context.Context, obj *model.UserSuggestion) ([]*model.Friend, error) {
	var friends []*model.Friend
	if err := r.friendService.GetAllFriendsByUserName(obj.Username, &friends); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return friends, nil
}

// Friend returns FriendResolver implementation.
func (r *Resolver) Friend() FriendResolver { return &friendResolver{r} }

// Inbox returns InboxResolver implementation.
func (r *Resolver) Inbox() InboxResolver { return &inboxResolver{r} }

// Post returns PostResolver implementation.
func (r *Resolver) Post() PostResolver { return &postResolver{r} }

// PostComment returns PostCommentResolver implementation.
func (r *Resolver) PostComment() PostCommentResolver { return &postCommentResolver{r} }

// PostLike returns PostLikeResolver implementation.
func (r *Resolver) PostLike() PostLikeResolver { return &postLikeResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

// UserProfile returns UserProfileResolver implementation.
func (r *Resolver) UserProfile() UserProfileResolver { return &userProfileResolver{r} }

// UserSuggestion returns UserSuggestionResolver implementation.
func (r *Resolver) UserSuggestion() UserSuggestionResolver { return &userSuggestionResolver{r} }

type friendResolver struct{ *Resolver }
type inboxResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
type postCommentResolver struct{ *Resolver }
type postLikeResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
type userProfileResolver struct{ *Resolver }
type userSuggestionResolver struct{ *Resolver }
