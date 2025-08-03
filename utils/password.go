package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 对密码进行加密
func HashPassword(password string) (string, error) {
	// 使用bcrypt对密码进行哈希加密
	// cost 参数控制加密强度，推荐值为12-14
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword 验证密码
func CheckPassword(password, hashedPassword string) bool {
	// 比较密码和哈希值
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
