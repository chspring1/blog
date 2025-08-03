package main

import (
	"os"

	"blog/database"
	"blog/routes"
	"github.com/sirupsen/logrus"
)

func init() {
	// 配置日志
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)
}

func main() {
	// 数据库配置
	config := database.DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     3306,
		Username: getEnv("DB_USERNAME", "root"),
		Password: getEnv("DB_PASSWORD", "password"), // 请修改为实际密码
		Database: getEnv("DB_NAME", "blog_db"),      // 请修改为实际数据库名
		Charset:  "utf8mb4",
	}

	// 初始化数据库连接
	err := database.InitDatabase(config)
	if err != nil {
		logrus.WithError(err).Fatal("数据库连接失败")
	}
	defer database.CloseDatabase()

	// 自动迁移数据库表
	err = database.AutoMigrate()
	if err != nil {
		logrus.WithError(err).Fatal("数据库迁移失败")
	}

	// 设置路由
	r := routes.SetupRoutes()

	// 启动服务器
	port := getEnv("PORT", "8080")
	logrus.WithField("port", port).Info("博客系统启动中...")

	if err := r.Run(":" + port); err != nil {
		logrus.WithError(err).Fatal("服务器启动失败")
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
