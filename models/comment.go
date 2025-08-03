package models

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Content   string         `json:"content" gorm:"not null;type:text"`
	PostID    uint           `json:"post_id" gorm:"not null;index"`
	Post      Post           `json:"post" gorm:"foreignKey:PostID"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	ParentID  *uint          `json:"parent_id" gorm:"index"`
	Parent    *Comment       `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies   []Comment      `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// CommentResponse 评论响应结构
type CommentResponse struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	PostID    uint      `json:"post_id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	ParentID  *uint     `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse 转换为响应格式
func (c *Comment) ToResponse() CommentResponse {
	return CommentResponse{
		ID:        c.ID,
		Content:   c.Content,
		PostID:    c.PostID,
		UserID:    c.UserID,
		Username:  c.User.Username,
		ParentID:  c.ParentID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}