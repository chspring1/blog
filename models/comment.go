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
	Status    int            `json:"status" gorm:"default:1;comment:1-正常 0-隐藏"`
	LikeCount int            `json:"like_count" gorm:"default:0"`
	IPAddress string         `json:"ip_address" gorm:"size:45"`
	UserAgent string         `json:"user_agent" gorm:"size:500"`
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
	Status    int       `json:"status"`
	LikeCount int       `json:"like_count"`
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
		Status:    c.Status,
		LikeCount: c.LikeCount,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}
