package models

import (
	"time"

	"blog/utils"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string         `json:"username" gorm:"type:varchar(50);uniqueIndex;not null" validate:"required,min=3,max=50"`
	Password  string         `json:"-" gorm:"type:varchar(255);not null" validate:"required,min=6"`
	Email     string         `json:"email" gorm:"type:varchar(100);uniqueIndex;not null" validate:"required,email"`
	Nickname  string         `json:"nickname" gorm:"type:varchar(50)"`
	Avatar    string         `json:"avatar" gorm:"type:varchar(255)"`
	Bio       string         `json:"bio" gorm:"type:text"`
	Status    int8           `json:"status" gorm:"type:tinyint;default:1;index"` // 1-正常，0-禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Posts    []Post    `json:"posts,omitempty" gorm:"foreignKey:UserID"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子 - 密码加密
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	return utils.CheckPassword(password, u.Password)
}

// UserResponse 用户响应结构（不包含敏感信息）
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
