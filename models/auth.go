package models

import (
	"time"
)

// Session 会话模型（可选，用于会话管理）
type Session struct {
	ID        string    `json:"id" gorm:"primaryKey;size:128"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Token     string    `json:"token" gorm:"not null;unique;size:512"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginLog 登录日志模型
type LoginLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	IPAddress string    `json:"ip_address" gorm:"size:45"`
	UserAgent string    `json:"user_agent" gorm:"size:500"`
	LoginAt   time.Time `json:"login_at" gorm:"not null"`
	Success   bool      `json:"success" gorm:"not null;default:true"`
	Message   string    `json:"message" gorm:"size:200"`
}

// TokenBlacklist 令牌黑名单（用于注销）
type TokenBlacklist struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Token     string    `json:"token" gorm:"not null;unique;size:512"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (Session) TableName() string {
	return "sessions"
}

// TableName 指定表名
func (LoginLog) TableName() string {
	return "login_logs"
}

// TableName 指定表名
func (TokenBlacklist) TableName() string {
	return "token_blacklist"
}
