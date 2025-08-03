package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建Gin引擎
	r := gin.Default()

	// 简单的健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "博客系统运行正常",
		})
	})

	// 启动服务器
	port := "8080"
	log.Printf("服务器启动在端口: %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
