package models

import (
	"time"

	"gorm.io/gorm"
)

// Category 分类模型
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"unique;not null;size:100"`
	Description string         `json:"description" gorm:"size:500"`
	Posts       []Post         `json:"posts,omitempty" gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Tag 标签模型
type Tag struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"unique;not null;size:50"`
	Posts     []Post         `json:"posts,omitempty" gorm:"many2many:post_tags;"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}
