package middleware

import (
	"strings"

	"blog/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logrus.Warn("请求缺少认证令牌")
			utils.UnauthorizedResponse(c, "请提供认证令牌")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			logrus.Warn("令牌格式错误")
			utils.UnauthorizedResponse(c, "令牌格式错误")
			c.Abort()
			return
		}

		// 提取token
		tokenString := authHeader[len(bearerPrefix):]
		if tokenString == "" {
			logrus.Warn("令牌为空")
			utils.UnauthorizedResponse(c, "令牌不能为空")
			c.Abort()
			return
		}

		// 验证token
		valid, claims := utils.ValidateToken(tokenString)
		if !valid {
			logrus.WithField("token", tokenString).Warn("令牌验证失败")
			utils.UnauthorizedResponse(c, "令牌无效或已过期")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		logrus.WithFields(logrus.Fields{
			"user_id":  claims.UserID,
			"username": claims.Username,
		}).Info("用户认证成功")

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（不强制要求认证）
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			const bearerPrefix = "Bearer "
			if strings.HasPrefix(authHeader, bearerPrefix) {
				tokenString := authHeader[len(bearerPrefix):]
				if tokenString != "" {
					valid, claims := utils.ValidateToken(tokenString)
					if valid {
						c.Set("user_id", claims.UserID)
						c.Set("username", claims.Username)
					}
				}
			}
		}
		c.Next()
	}
}

// GetCurrentUserID 从上下文获取当前用户ID
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint)
	return id, ok
}

// GetCurrentUsername 从上下文获取当前用户名
func GetCurrentUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}

	name, ok := username.(string)
	return name, ok
}
