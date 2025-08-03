package controllers

import (
	"strconv"

	"blog/database"
	"blog/middleware"
	"blog/models"
	"blog/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CommentController 评论控制器
type CommentController struct{}

// NewCommentController 创建评论控制器实例
func NewCommentController() *CommentController {
	return &CommentController{}
}

// CreateComment 创建评论
func (cc *CommentController) CreateComment(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "未授权访问")
		return
	}

	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的文章ID")
		return
	}

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("创建评论参数绑定失败")
		utils.BadRequestResponse(c, "请求参数错误: "+err.Error())
		return
	}

	db := database.GetDB()

	// 检查文章是否存在
	var post models.Post
	if err := db.First(&post, uint(postID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "文章不存在")
		} else {
			logrus.WithError(err).Error("查询文章失败")
			utils.InternalServerErrorResponse(c, "查询文章失败")
		}
		return
	}

	// 检查文章状态
	if post.Status != 1 {
		utils.BadRequestResponse(c, "文章已被删除或未发布")
		return
	}

	// 创建评论
	comment := models.Comment{
		Content:   req.Content,
		UserID:    userID,
		PostID:    uint(postID),
		Status:    1, // 默认正常状态
		LikeCount: 0,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}

	if err := db.Create(&comment).Error; err != nil {
		logrus.WithError(err).Error("创建评论失败")
		utils.InternalServerErrorResponse(c, "创建评论失败")
		return
	}

	// 更新文章评论数
	if err := db.Model(&post).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
		logrus.WithError(err).Warn("更新文章评论数失败")
	}

	// 预加载用户信息
	if err := db.Preload("User").First(&comment, comment.ID).Error; err != nil {
		logrus.WithError(err).Error("获取评论详情失败")
		utils.InternalServerErrorResponse(c, "获取评论详情失败")
		return
	}

	logrus.WithFields(logrus.Fields{
		"comment_id": comment.ID,
		"post_id":    postID,
		"user_id":    userID,
	}).Info("评论创建成功")

	utils.SuccessResponse(c, comment.ToResponse(), "评论创建成功")
}

// GetComments 获取文章评论列表
func (cc *CommentController) GetComments(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的文章ID")
		return
	}

	db := database.GetDB()

	// 检查文章是否存在
	var post models.Post
	if err := db.First(&post, uint(postID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "文章不存在")
		} else {
			logrus.WithError(err).Error("查询文章失败")
			utils.InternalServerErrorResponse(c, "查询文章失败")
		}
		return
	}

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 计算偏移量
	offset := (page - 1) * pageSize

	var comments []models.Comment
	var total int64

	// 查询总数
	if err := db.Model(&models.Comment{}).
		Where("post_id = ? AND status = ?", uint(postID), 1).
		Count(&total).Error; err != nil {
		logrus.WithError(err).Error("查询评论总数失败")
		utils.InternalServerErrorResponse(c, "查询评论列表失败")
		return
	}

	// 查询评论列表（预加载用户信息）
	if err := db.Preload("User").
		Where("post_id = ? AND status = ?", uint(postID), 1).
		Order("created_at ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&comments).Error; err != nil {
		logrus.WithError(err).Error("查询评论列表失败")
		utils.InternalServerErrorResponse(c, "查询评论列表失败")
		return
	}

	// 转换为响应格式
	var commentResponses []models.CommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, comment.ToResponse())
	}

	response := gin.H{
		"comments": commentResponses,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}

	utils.SuccessResponse(c, response, "获取评论列表成功")
}

// DeleteComment 删除评论
func (cc *CommentController) DeleteComment(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "未授权访问")
		return
	}

	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的评论ID")
		return
	}

	db := database.GetDB()
	var comment models.Comment

	// 查询评论
	if err := db.First(&comment, uint(commentID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "评论不存在")
		} else {
			logrus.WithError(err).Error("查询评论失败")
			utils.InternalServerErrorResponse(c, "查询评论失败")
		}
		return
	}

	// 检查权限（只有评论作者可以删除）
	if comment.UserID != userID {
		logrus.WithFields(logrus.Fields{
			"comment_id":   comment.ID,
			"comment_user": comment.UserID,
			"current_user": userID,
		}).Warn("用户尝试删除非本人评论")
		utils.ForbiddenResponse(c, "无权限删除此评论")
		return
	}

	// 软删除评论
	if err := db.Delete(&comment).Error; err != nil {
		logrus.WithError(err).Error("删除评论失败")
		utils.InternalServerErrorResponse(c, "删除评论失败")
		return
	}

	// 更新文章评论数
	var post models.Post
	if err := db.First(&post, comment.PostID).Error; err == nil {
		if err := db.Model(&post).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
			logrus.WithError(err).Warn("更新文章评论数失败")
		}
	}

	logrus.WithFields(logrus.Fields{
		"comment_id": comment.ID,
		"user_id":    userID,
	}).Info("评论删除成功")

	utils.SuccessResponse(c, nil, "评论删除成功")
}
