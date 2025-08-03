package models

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" validate:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6" validate:"required,min=6"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Nickname string `json:"nickname" binding:"max=50" validate:"max=50"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// ChangePasswordRequest 修改密码请求结构
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UpdateProfileRequest 更新个人信息请求结构
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Avatar   string `json:"avatar" binding:"max=255"`
	Bio      string `json:"bio" binding:"max=500"`
}

// CreatePostRequest 创建文章请求结构
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,max=255"`
	Content string `json:"content" binding:"required"`
	Excerpt string `json:"excerpt" binding:"max=500"`
}

// UpdatePostRequest 更新文章请求结构
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"max=255"`
	Content string `json:"content"`
	Excerpt string `json:"excerpt" binding:"max=500"`
}

// CreateCommentRequest 创建评论请求结构
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,max=1000"`
}
