package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"
	"outstagram/common"
	"outstagram/graph/model"

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
		User:     &model.User{},
		Posts:    []*model.Post{},
		Friends:  []*model.User{},
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
	_, isOk := ctx.Value(common.UserIDLocalKey).(string)

	var post *model.Post
	if err := r.postService.PostGetAllByPostID(isOk, postID, &post); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return post, nil
}

// PostHomePage is the resolver for the postHomePage field.
func (r *queryResolver) PostHomePage(ctx context.Context, page int) ([]*model.Post, error) {
	currentUserID, isOk := ctx.Value(common.UserIDLocalKey).(string)
	if !isOk {
		return nil, gqlerror.Errorf("user not found")
	}

	var posts []*model.Post
	if err := r.postService.PostGetHomePage(page, currentUserID, &posts); err != nil {
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

// User is the resolver for the user field.
func (r *userProfileResolver) User(ctx context.Context, obj *model.UserProfile) (*model.User, error) {
	var userRecord *model.User
	if err := r.userService.UserGetByUserName(obj.Username, &userRecord); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	if userRecord.ID == "" {
		return nil, gqlerror.Errorf("user not found")
	}

	return userRecord, nil
}

// Posts is the resolver for the posts field.
func (r *userProfileResolver) Posts(ctx context.Context, obj *model.UserProfile) ([]*model.Post, error) {
	var posts []*model.Post
	if err := r.postService.PostGetAllByUserName(obj.Username, &posts); err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return posts, nil
}

// Friends is the resolver for the friends field.
func (r *userProfileResolver) Friends(ctx context.Context, obj *model.UserProfile) ([]*model.Friend, error) {
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

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// UserProfile returns UserProfileResolver implementation.
func (r *Resolver) UserProfile() UserProfileResolver { return &userProfileResolver{r} }

type friendResolver struct{ *Resolver }
type inboxResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
type postCommentResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userProfileResolver struct{ *Resolver }
