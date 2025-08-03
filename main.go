package main

import (
	"log"

	"blog/database"
	"blog/routes"
)

func main() {
	// 初始化数据库
	database.InitDB()

	// 设置路由
	r := routes.SetupRoutes()

	// 启动服务器
	port := "8080"
	log.Printf("服务器启动在端口: %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
