# 个人博客系统 API 文档

## 项目介绍

这是一个使用 Go + Gin + GORM 开发的个人博客系统后端，提供完整的用户认证、文章管理和评论功能。

## 技术栈

- **Go 1.24.5** - 编程语言
- **Gin** - Web框架
- **GORM** - ORM库
- **MySQL** - 数据库
- **JWT** - 用户认证
- **bcrypt** - 密码加密
- **Logrus** - 日志记录

## 环境变量配置

```bash
# 数据库配置
DB_HOST=localhost
DB_USERNAME=root
DB_PASSWORD=your_password
DB_NAME=blog_db

# 服务器配置
PORT=8080
```

## API 接口

### 基础URL
```
http://localhost:8080/api/v1
```

### 1. 用户认证

#### 用户注册
```http
POST /auth/register
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com",
  "nickname": "测试用户"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "nickname": "测试用户",
      "status": 1,
      "created_at": "2025-08-03T10:00:00Z",
      "updated_at": "2025-08-03T10:00:00Z"
    }
  }
}
```

#### 用户登录
```http
POST /auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

### 2. 用户管理 (需要认证)

**认证头:**
```http
Authorization: Bearer <your_jwt_token>
```

#### 获取个人信息
```http
GET /user/profile
```

#### 更新个人信息
```http
PUT /user/profile
Content-Type: application/json

{
  "nickname": "新昵称",
  "avatar": "http://example.com/avatar.jpg",
  "bio": "个人简介"
}
```

#### 修改密码
```http
PUT /user/password
Content-Type: application/json

{
  "old_password": "oldpassword",
  "new_password": "newpassword123"
}
```

### 3. 文章管理

#### 获取文章列表 (公开)
```http
GET /posts?page=1&page_size=10
```

**响应:**
```json
{
  "code": 200,
  "message": "获取文章列表成功",
  "data": {
    "posts": [
      {
        "id": 1,
        "title": "文章标题",
        "content": "文章内容...",
        "excerpt": "文章摘要",
        "user_id": 1,
        "status": 1,
        "view_count": 100,
        "comment_count": 5,
        "published_at": "2025-08-03T10:00:00Z",
        "user": {
          "id": 1,
          "username": "testuser",
          "nickname": "测试用户"
        }
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 10,
      "total": 25,
      "total_page": 3
    }
  }
}
```

#### 获取文章详情 (公开)
```http
GET /posts/1
```

#### 创建文章 (需要认证)
```http
POST /posts
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "title": "我的第一篇博客",
  "content": "这是文章的详细内容...",
  "excerpt": "这是文章摘要"
}
```

#### 更新文章 (需要认证，仅作者)
```http
PUT /posts/1
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "title": "更新后的标题",
  "content": "更新后的内容..."
}
```

#### 删除文章 (需要认证，仅作者)
```http
DELETE /posts/1
Authorization: Bearer <your_jwt_token>
```

### 4. 评论管理

#### 获取文章评论 (公开)
```http
GET /posts/1/comments?page=1&page_size=20
```

#### 创建评论 (需要认证)
```http
POST /posts/1/comments
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "content": "这是一条评论内容"
}
```

#### 删除评论 (需要认证，仅作者)
```http
DELETE /comments/1
Authorization: Bearer <your_jwt_token>
```

### 5. 系统健康检查

#### 健康检查 (公开)
```http
GET /health
```

**响应:**
```json
{
  "status": "ok",
  "message": "博客系统运行正常"
}
```

## 错误响应格式

```json
{
  "code": 400,
  "message": "错误信息描述"
}
```

### 常见错误码

- **200** - 成功
- **400** - 请求参数错误
- **401** - 未授权访问
- **403** - 权限不足
- **404** - 资源不存在
- **500** - 服务器内部错误

## 部署说明

### 1. 环境准备

```bash
# 安装 Go 1.24.5+
# 安装 MySQL 8.0+
```

### 2. 数据库初始化

```sql
CREATE DATABASE blog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. 运行项目

```bash
# 下载依赖
go mod tidy

# 设置环境变量
export DB_PASSWORD="your_mysql_password"
export DB_NAME="blog_db"

# 运行项目
go run main.go
```

### 4. 使用Docker部署

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o blog-server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/blog-server .
EXPOSE 8080
CMD ["./blog-server"]
```

## 数据库表结构

详细的数据库表结构请参考 `database/schema.sql` 文件。

主要包含：
- **users** - 用户表
- **posts** - 文章表
- **comments** - 评论表
- **categories** - 分类表
- **tags** - 标签表
- **post_tags** - 文章标签关联表

## 日志记录

系统使用 Logrus 进行日志记录，日志格式为 JSON，包含以下信息：
- 请求方法和路径
- 用户身份信息
- 操作结果和错误信息
- 时间戳

## 安全特性

1. **密码加密** - 使用 bcrypt 对密码进行哈希加密
2. **JWT认证** - 使用 JWT 进行用户身份验证
3. **权限控制** - 用户只能操作自己的数据
4. **参数验证** - 对所有输入参数进行严格验证
5. **SQL注入防护** - 使用 GORM 的参数化查询防止 SQL 注入
6. **日志审计** - 记录所有重要操作的日志

## 性能优化

1. **数据库索引** - 为常用查询字段添加索引
2. **分页查询** - 对列表查询实现分页
3. **预加载** - 使用 GORM 的 Preload 减少 N+1 查询
4. **软删除** - 使用软删除提高数据安全性
5. **统计计数** - 维护文章浏览数、评论数等统计信息

## 扩展功能建议

1. **文件上传** - 实现图片上传功能
2. **标签系统** - 完善标签管理功能
3. **搜索功能** - 实现文章搜索
4. **缓存机制** - 使用 Redis 缓存热点数据
5. **邮件通知** - 实现评论通知功能
6. **管理后台** - 开发管理员后台界面
