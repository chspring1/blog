package database

import (
	"log"

	"blog/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	// 数据库连接配置
	dsn := "root:password@tcp(localhost:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库结构
	err = db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Category{},
		&models.Tag{},
		&models.Comment{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected and migrated successfully")
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return db
}
