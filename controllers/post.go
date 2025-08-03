package controllers

import (
	"strconv"
	"time"

	"blog/database"
	"blog/middleware"
	"blog/models"
	"blog/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// PostController 文章控制器
type PostController struct{}

// NewPostController 创建文章控制器实例
func NewPostController() *PostController {
	return &PostController{}
}

// CreatePost 创建文章
func (pc *PostController) CreatePost(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "未授权访问")
		return
	}

	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("创建文章参数绑定失败")
		utils.BadRequestResponse(c, "请求参数错误: "+err.Error())
		return
	}

	db := database.GetDB()

	// 创建文章
	post := models.Post{
		Title:        req.Title,
		Content:      req.Content,
		Excerpt:      req.Excerpt,
		UserID:       userID,
		Status:       1, // 默认已发布
		ViewCount:    0,
		CommentCount: 0,
		LikeCount:    0,
		IsTop:        0,
		PublishedAt:  &time.Time{},
	}

	// 如果没有提供摘要，自动生成
	if post.Excerpt == "" && len(post.Content) > 100 {
		post.Excerpt = post.Content[:100] + "..."
	}

	// 设置发布时间
	now := time.Now()
	post.PublishedAt = &now

	if err := db.Create(&post).Error; err != nil {
		logrus.WithError(err).Error("创建文章失败")
		utils.InternalServerErrorResponse(c, "创建文章失败")
		return
	}

	// 预加载用户信息
	if err := db.Preload("User").First(&post, post.ID).Error; err != nil {
		logrus.WithError(err).Error("获取文章详情失败")
		utils.InternalServerErrorResponse(c, "获取文章详情失败")
		return
	}

	logrus.WithFields(logrus.Fields{
		"post_id": post.ID,
		"user_id": userID,
		"title":   post.Title,
	}).Info("文章创建成功")

	utils.SuccessResponse(c, post.ToResponse(), "文章创建成功")
}

// GetPosts 获取文章列表
func (pc *PostController) GetPosts(c *gin.Context) {
	db := database.GetDB()

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 计算偏移量
	offset := (page - 1) * pageSize

	var posts []models.Post
	var total int64

	// 查询总数
	if err := db.Model(&models.Post{}).Where("status = ?", 1).Count(&total).Error; err != nil {
		logrus.WithError(err).Error("查询文章总数失败")
		utils.InternalServerErrorResponse(c, "查询文章列表失败")
		return
	}

	// 查询文章列表（预加载用户信息）
	if err := db.Preload("User").
		Where("status = ?", 1).
		Order("is_top DESC, published_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&posts).Error; err != nil {
		logrus.WithError(err).Error("查询文章列表失败")
		utils.InternalServerErrorResponse(c, "查询文章列表失败")
		return
	}

	// 转换为响应格式
	var postResponses []models.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, post.ToResponse())
	}

	response := gin.H{
		"posts": postResponses,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}

	utils.SuccessResponse(c, response, "获取文章列表成功")
}

// GetPost 获取单篇文章详情
func (pc *PostController) GetPost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的文章ID")
		return
	}

	db := database.GetDB()
	var post models.Post

	// 查询文章（预加载用户信息）
	if err := db.Preload("User").First(&post, uint(postID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "文章不存在")
		} else {
			logrus.WithError(err).Error("查询文章详情失败")
			utils.InternalServerErrorResponse(c, "查询文章详情失败")
		}
		return
	}

	// 检查文章状态
	if post.Status != 1 {
		utils.NotFoundResponse(c, "文章不存在")
		return
	}

	// 增加浏览次数
	if err := db.Model(&post).Update("view_count", gorm.Expr("view_count + ?", 1)).Error; err != nil {
		logrus.WithError(err).Warn("更新文章浏览次数失败")
	}

	utils.SuccessResponse(c, post.ToResponse(), "获取文章详情成功")
}

// UpdatePost 更新文章
func (pc *PostController) UpdatePost(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "未授权访问")
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的文章ID")
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("更新文章参数绑定失败")
		utils.BadRequestResponse(c, "请求参数错误: "+err.Error())
		return
	}

	db := database.GetDB()
	var post models.Post

	// 查询文章
	if err := db.First(&post, uint(postID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "文章不存在")
		} else {
			logrus.WithError(err).Error("查询文章失败")
			utils.InternalServerErrorResponse(c, "查询文章失败")
		}
		return
	}

	// 检查权限（只有作者可以修改）
	if post.UserID != userID {
		logrus.WithFields(logrus.Fields{
			"post_id":      post.ID,
			"post_user":    post.UserID,
			"current_user": userID,
		}).Warn("用户尝试修改非本人文章")
		utils.ForbiddenResponse(c, "无权限修改此文章")
		return
	}

	// 更新文章信息
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Excerpt != "" {
		post.Excerpt = req.Excerpt
	}

	if err := db.Save(&post).Error; err != nil {
		logrus.WithError(err).Error("更新文章失败")
		utils.InternalServerErrorResponse(c, "更新文章失败")
		return
	}

	// 预加载用户信息
	if err := db.Preload("User").First(&post, post.ID).Error; err != nil {
		logrus.WithError(err).Error("获取更新后文章详情失败")
		utils.InternalServerErrorResponse(c, "获取文章详情失败")
		return
	}

	logrus.WithFields(logrus.Fields{
		"post_id": post.ID,
		"user_id": userID,
	}).Info("文章更新成功")

	utils.SuccessResponse(c, post.ToResponse(), "文章更新成功")
}

// DeletePost 删除文章
func (pc *PostController) DeletePost(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "未授权访问")
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的文章ID")
		return
	}

	db := database.GetDB()
	var post models.Post

	// 查询文章
	if err := db.First(&post, uint(postID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "文章不存在")
		} else {
			logrus.WithError(err).Error("查询文章失败")
			utils.InternalServerErrorResponse(c, "查询文章失败")
		}
		return
	}

	// 检查权限（只有作者可以删除）
	if post.UserID != userID {
		logrus.WithFields(logrus.Fields{
			"post_id":      post.ID,
			"post_user":    post.UserID,
			"current_user": userID,
		}).Warn("用户尝试删除非本人文章")
		utils.ForbiddenResponse(c, "无权限删除此文章")
		return
	}

	// 软删除文章
	if err := db.Delete(&post).Error; err != nil {
		logrus.WithError(err).Error("删除文章失败")
		utils.InternalServerErrorResponse(c, "删除文章失败")
		return
	}

	logrus.WithFields(logrus.Fields{
		"post_id": post.ID,
		"user_id": userID,
	}).Info("文章删除成功")

	utils.SuccessResponse(c, nil, "文章删除成功")
}
