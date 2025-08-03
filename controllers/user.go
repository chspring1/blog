package controllers

import (
	"blog/database"
	"blog/middleware"
	"blog/models"
	"blog/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// UserController 用户控制器
type UserController struct{}

// NewUserController 创建用户控制器实例
func NewUserController() *UserController {
	return &UserController{}
}

// Register 用户注册
func (uc *UserController) Register(c *gin.Context) {
	var req models.RegisterRequest

	// 绑定JSON数据到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("用户注册参数绑定失败")
		utils.BadRequestResponse(c, "请求参数错误: "+err.Error())
		return
	}

	db := database.GetDB()

	// 检查用户名是否已存在
	var existUser models.User
	if err := db.Where("username = ?", req.Username).First(&existUser).Error; err == nil {
		logrus.WithField("username", req.Username).Warn("用户名已存在")
		utils.BadRequestResponse(c, "用户名已存在")
		return
	}

	// 检查邮箱是否已存在
	if err := db.Where("email = ?", req.Email).First(&existUser).Error; err == nil {
		logrus.WithField("email", req.Email).Warn("邮箱已被注册")
		utils.BadRequestResponse(c, "邮箱已被注册")
		return
	}

	// 创建新用户
	user := models.User{
		Username: req.Username,
		Password: req.Password, // 密码会在BeforeCreate钩子中自动加密
		Email:    req.Email,
		Nickname: req.Nickname,
		Status:   1, // 默认激活状态
	}

	// 保存用户到数据库
	if err := db.Create(&user).Error; err != nil {
		logrus.WithError(err).Error("用户注册失败")
		utils.InternalServerErrorResponse(c, "注册失败")
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		logrus.WithError(err).Error("生成JWT令牌失败")
		utils.InternalServerErrorResponse(c, "生成令牌失败")
		return
	}

	// 返回注册成功响应
	response := models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}

	logrus.WithFields(logrus.Fields{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("用户注册成功")

	utils.SuccessResponse(c, response, "注册成功")
}

// Login 用户登录
func (uc *UserController) Login(c *gin.Context) {
	var req models.LoginRequest

	// 绑定JSON数据到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("用户登录参数绑定失败")
		utils.BadRequestResponse(c, "请求参数错误: "+err.Error())
		return
	}

	db := database.GetDB()

	// 查找用户（支持用户名或邮箱登录）
	var user models.User
	if err := db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.WithField("username", req.Username).Warn("用户登录失败：用户不存在")
			utils.UnauthorizedResponse(c, "用户名或密码错误")
		} else {
			logrus.WithError(err).Error("查询用户失败")
			utils.InternalServerErrorResponse(c, "登录失败")
		}
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		logrus.WithField("user_id", user.ID).Warn("尝试登录被禁用的账户")
		utils.UnauthorizedResponse(c, "账户已被禁用")
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		logrus.WithField("user_id", user.ID).Warn("用户登录失败：密码错误")
		utils.UnauthorizedResponse(c, "用户名或密码错误")
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		logrus.WithError(err).Error("生成JWT令牌失败")
		utils.InternalServerErrorResponse(c, "生成令牌失败")
		return
	}

	// 返回登录成功响应
	response := models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}

	logrus.WithFields(logrus.Fields{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("用户登录成功")

	utils.SuccessResponse(c, response, "登录成功")
}

// GetProfile 获取用户个人信息
func (uc *UserController) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "未授权访问")
		return
	}

	db := database.GetDB()
	var user models.User

	if err := db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "用户不存在")
		} else {
			logrus.WithError(err).Error("获取用户信息失败")
			utils.InternalServerErrorResponse(c, "获取用户信息失败")
		}
		return
	}

	utils.SuccessResponse(c, user.ToResponse(), "获取成功")
}

// UpdateProfile 更新用户个人信息
func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "未授权访问")
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("更新用户信息参数绑定失败")
		utils.BadRequestResponse(c, "请求参数错误: "+err.Error())
		return
	}

	db := database.GetDB()
	var user models.User

	if err := db.First(&user, userID).Error; err != nil {
		utils.NotFoundResponse(c, "用户不存在")
		return
	}

	// 更新用户信息
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	if err := db.Save(&user).Error; err != nil {
		logrus.WithError(err).Error("更新用户信息失败")
		utils.InternalServerErrorResponse(c, "更新失败")
		return
	}

	logrus.WithField("user_id", userID).Info("用户信息更新成功")
	utils.SuccessResponse(c, user.ToResponse(), "更新成功")
}

// ChangePassword 修改密码
func (uc *UserController) ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "未授权访问")
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("修改密码参数绑定失败")
		utils.BadRequestResponse(c, "请求参数错误: "+err.Error())
		return
	}

	db := database.GetDB()
	var user models.User

	if err := db.First(&user, userID).Error; err != nil {
		utils.NotFoundResponse(c, "用户不存在")
		return
	}

	// 验证旧密码
	if !user.CheckPassword(req.OldPassword) {
		logrus.WithField("user_id", userID).Warn("修改密码失败：原密码错误")
		utils.BadRequestResponse(c, "原密码错误")
		return
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		logrus.WithError(err).Error("密码加密失败")
		utils.InternalServerErrorResponse(c, "密码加密失败")
		return
	}

	// 更新密码
	user.Password = hashedPassword
	if err := db.Save(&user).Error; err != nil {
		logrus.WithError(err).Error("密码修改失败")
		utils.InternalServerErrorResponse(c, "密码修改失败")
		return
	}

	logrus.WithField("user_id", userID).Info("密码修改成功")
	utils.SuccessResponse(c, nil, "密码修改成功")
}
