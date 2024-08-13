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
)

type AdminController struct {
	adminService *services.AdminService
	postService  *services.PostService
	userService  *services.UserService
}

func NewAdminController(adminService *services.AdminService, postService *services.PostService, userService *services.UserService) *AdminController {
	return &AdminController{
		adminService: adminService,
		postService:  postService,
		userService:  userService,
	}
}

func (a *AdminController) AdminDeletePostByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")
	rawUserID := ctx.Params("userID")
	userID := uuid.MustParse(rawUserID)

	var post entity.Post
	if err := common.DBConn.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusBadRequest, "Post not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying post")
	}

	if err := a.postService.PostDeleteByPostID(postID, userID); err != nil {
		return err
	}

	post.Active = false

	if err := common.DBConn.Save(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while updating post to inactive")
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Post deleted and deactivated", nil)
}

func (a *AdminController) AdminBanUserByUserID(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")

	var userRecord entity.User
	if err := a.userService.UserGetByID(userID, &userRecord); err != nil {
		return err
	}

	userActive, err := a.userService.UserBanByUserID(userRecord.ID.String())
	if err != nil {
		return err
	}

	if userActive == true {
		return common.CreateResponse(ctx, fiber.StatusOK, "User unbanned", userActive)

	}

	return common.CreateResponse(ctx, fiber.StatusOK, "User banned", userActive)
}

func (a *AdminController) AdminBlockPostByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")

	postActive, err := a.postService.PostBlockByPostID(postID)
	if err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Post blocked", postActive)
}

func (a *AdminController) AdminDeleteCommentOnPostByCommentID(ctx *fiber.Ctx) error {
	postID := ctx.Params("postID")

	rawUserID := ctx.Params("userID")
	userID := uuid.MustParse(rawUserID)

	rawCommentID := ctx.Params("commentID")
	commentID := uuid.MustParse(rawCommentID)

	if err := a.postService.PostDeleteCommentOnPostByCommentID(postID, userID, commentID); err != nil {
		return err
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Comment deleted", nil)
}

func (a *AdminController) AdminGetReports(ctx *fiber.Ctx) error {
	var reportRecords []entity.Report
	var reportResponses []req.ReportResponse

	if err := common.DBConn.Order("created_at DESC").Find(&reportRecords).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "No report found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying reports")
	}

	for _, report := range reportRecords {
		response := req.ReportResponse{
			ID:        report.ID,
			ByUserID:  report.ByUserID,
			Type:      report.Type,
			Info:      report.Info,
			Reason:    report.Reason,
			Status:    report.Status,
			CreatedAt: report.CreatedAt,
		}

		var reportingUser entity.User
		if err := common.DBConn.Where("id = ?", report.ByUserID).First(&reportingUser).Error; err == nil {
			response.ReportingUser = &reportingUser
		}

		switch report.Type {
		case 4:
			var user entity.User
			if err := common.DBConn.Where("id = ?", report.Info).First(&user).Error; err == nil {
				response.User = &user
			}
		case 1, 2:
			var post entity.Post

			if err := common.DBConn.
				Joins("JOIN users ON users.id = posts.user_id").
				Where("posts.id = ?", report.Info).
				Where("posts.active = ?", true).
				Where("users.active = ?", true).
				Preload("PostFiles").
				First(&post).Error; err == nil {
				response.Post = &post

				var user entity.User
				if err := common.DBConn.Where("id = ?", post.UserID).First(&user).Error; err == nil {
					response.User = &user
				}
			}
		case 3:
			var comment entity.PostComment
			if err := common.DBConn.Where("id = ?", report.Info).First(&comment).Error; err == nil {
				response.Comment = &comment

				var user entity.User
				if err := common.DBConn.Where("id = ?", comment.UserID).First(&user).Error; err == nil {
					response.User = &user
				}
			}
		}

		reportResponses = append(reportResponses, response)
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Reports found", reportResponses)
}

func (a *AdminController) AdminUpdateStatusReport(ctx *fiber.Ctx) error {
	var reportRecord entity.Report
	rawReportID := ctx.Params("reportID")
	reportID, err := uuid.Parse(rawReportID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid report ID format")
	}

	if err := common.DBConn.Where("id = ?", reportID).First(&reportRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Report not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while querying report")
	}

	reportRecord.Status = 1

	if err := common.DBConn.Save(&reportRecord).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while updating report status")
	}

	return common.CreateResponse(ctx, fiber.StatusOK, "Report status updated", reportRecord.Status)
}
