package services

import (
	"errors"

	"DLM_backend/database"
	"DLM_backend/models"
	"DLM_backend/utils"
)

// AuthenticateUser 校验用户凭据，并生成 JWT 令牌
func AuthenticateUser(username, password, role string) (string, error) {
	var user models.User

	// 查询指定用户名、密码和角色的用户
	if err := database.DB.Where("username = ? AND password = ? AND role = ?",
		username, password, role).First(&user).Error; err != nil {
		return "", errors.New("invalid credentials")
	}

	// 生成包含用户名和角色的令牌
	token, err := utils.GenerateToken(user.Username, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}
