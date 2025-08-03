package models

import (
	"time"

	"gorm.io/gorm"
)

// Post 博客文章模型
type Post struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Title        string         `json:"title" gorm:"not null;size:200"`
	Content      string         `json:"content" gorm:"type:text"`
	Summary      string         `json:"summary" gorm:"size:500"`
	Excerpt      string         `json:"excerpt" gorm:"size:500"`
	Status       int            `json:"status" gorm:"default:1;comment:1-已发布 0-草稿"`
	ViewCount    uint           `json:"view_count" gorm:"default:0"`
	CommentCount int            `json:"comment_count" gorm:"default:0"`
	LikeCount    int            `json:"like_count" gorm:"default:0"`
	IsTop        int            `json:"is_top" gorm:"default:0;comment:1-置顶 0-普通"`
	UserID       uint           `json:"user_id" gorm:"not null;index"`
	User         User           `json:"user" gorm:"foreignKey:UserID"`
	CategoryID   *uint          `json:"category_id" gorm:"index"`
	Category     *Category      `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Tags         []Tag          `json:"tags,omitempty" gorm:"many2many:post_tags;"`
	Comments     []Comment      `json:"comments,omitempty" gorm:"foreignKey:PostID"`
	PublishedAt  *time.Time     `json:"published_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// PostResponse 博客文章响应结构
type PostResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Summary     string    `json:"summary"`
	Status      int       `json:"status"`
	ViewCount   uint      `json:"view_count"`
	UserID      uint      `json:"user_id"`
	Username    string    `json:"username"`
	CategoryID  *uint     `json:"category_id"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse 转换为响应格式
func (p *Post) ToResponse() PostResponse {
	var category string
	if p.Category != nil {
		category = p.Category.Name
	}

	var tags []string
	for _, tag := range p.Tags {
		tags = append(tags, tag.Name)
	}

	return PostResponse{
		ID:         p.ID,
		Title:      p.Title,
		Content:    p.Content,
		Summary:    p.Summary,
		Status:     p.Status,
		ViewCount:  p.ViewCount,
		UserID:     p.UserID,
		Username:   p.User.Username,
		CategoryID: p.CategoryID,
		Category:   category,
		Tags:       tags,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

// TableName 指定表名
func (Post) TableName() string {
	return "posts"
}