package routes

import (
	"blog/controllers"
	"blog/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SetupRoutes 设置路由
func SetupRoutes() *gin.Engine {
	// 设置gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建gin引擎
	r := gin.New()

	// 使用中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 创建控制器实例
	userController := controllers.NewUserController()
	postController := controllers.NewPostController()
	commentController := controllers.NewCommentController()

	// API版本分组
	v1 := r.Group("/api/v1")

	// 认证相关路由（无需认证）
	auth := v1.Group("/auth")
	{
		auth.POST("/register", userController.Register) // 用户注册
		auth.POST("/login", userController.Login)       // 用户登录
	}

	// 用户相关路由（需要认证）
	user := v1.Group("/user")
	user.Use(middleware.AuthMiddleware()) // 应用认证中间件
	{
		user.GET("/profile", userController.GetProfile)      // 获取个人信息
		user.PUT("/profile", userController.UpdateProfile)   // 更新个人信息
		user.PUT("/password", userController.ChangePassword) // 修改密码
	}

	// 文章相关路由
	posts := v1.Group("/posts")
	{
		// 公共接口（无需认证）
		posts.GET("", postController.GetPosts)    // 获取文章列表
		posts.GET("/:id", postController.GetPost) // 获取文章详情

		// 需要认证的接口
		posts.Use(middleware.AuthMiddleware())
		posts.POST("", postController.CreatePost)       // 创建文章
		posts.PUT("/:id", postController.UpdatePost)    // 更新文章
		posts.DELETE("/:id", postController.DeletePost) // 删除文章
	}

	// 评论相关路由
	comments := v1.Group("/posts/:post_id/comments")
	{
		// 公共接口（无需认证）
		comments.GET("", commentController.GetComments) // 获取评论列表

		// 需要认证的接口
		comments.Use(middleware.AuthMiddleware())
		comments.POST("", commentController.CreateComment) // 创建评论
	}

	// 评论管理路由（需要认证）
	commentManage := v1.Group("/comments")
	commentManage.Use(middleware.AuthMiddleware())
	{
		commentManage.DELETE("/:id", commentController.DeleteComment) // 删除评论
	}

	// 健康检查接口
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "博客系统运行正常",
		})
	})

	logrus.Info("路由配置完成")
	return r
}
