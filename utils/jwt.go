package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// jwtSecret 默认的密钥，可以通过 SetJWTSecret 覆盖
var jwtSecret = []byte("secret")

// SetJWTSecret 设置 JWT 密钥
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

// GenerateToken 根据用户名和角色生成 JWT 令牌
func GenerateToken(username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析 JWT 令牌
func ParseToken(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
