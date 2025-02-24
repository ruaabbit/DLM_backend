package services

import (
	"errors"

	"DLM_backend/database"
	"DLM_backend/models"
	"DLM_backend/utils"
)

// AuthenticateUser 校验用户凭据，并生成 JWT 令牌
func AuthenticateUser(username, password string) (string, error) {
	var user models.User
	if err := database.DB.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.Username, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}
