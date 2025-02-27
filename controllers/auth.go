package controllers

import (
	"DLM_backend/database"
	"DLM_backend/models"
	"DLM_backend/services"
	"DLM_backend/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// LoginRequest 用户登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"` // 角色: "keeper"(保管员) 或 "admin"(管理员)
}

// UserProfileRequest 用户个人信息修改请求结构体
type UserProfileRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// Login 用于处理用户登录请求
func Login(c *gin.Context) {
	var loginData LoginRequest
	if err := c.ShouldBindJSON(&loginData); err != nil {
		utils.ErrorResponse(c, err.Error())
		return
	}

	// 验证角色是否有效
	if loginData.Role != "keeper" && loginData.Role != "admin" {
		utils.ErrorResponse(c, "invalid role, must be 'keeper' or 'admin'")
		return
	}

	// 根据用户名、密码和角色进行认证
	token, err := services.AuthenticateUser(loginData.Username, loginData.Password, loginData.Role)
	if err != nil {
		utils.ErrorResponse(c, "invalid credentials")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"token": token,
		"user": gin.H{
			"username": loginData.Username,
			"role":     loginData.Role,
		},
	})
}

// GetUserProfile 获取用户个人信息
func GetUserProfile(c *gin.Context) {
	// 从JWT中获取用户名
	claims, exists := c.Get("claims")
	if !exists {
		utils.UnauthorizedResponse(c, "token claims not found")
		return
	}

	username := claims.(jwt.MapClaims)["username"].(string)

	// 从数据库获取用户信息
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		utils.NotFoundResponse(c, "user not found")
		return
	}

	// 返回用户信息，排除敏感信息如密码
	utils.SuccessResponse(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"name":     user.Name,
		"phone":    user.Phone,
	})
}

// UpdateUserProfile 更新用户个人信息
func UpdateUserProfile(c *gin.Context) {
	// 从JWT中获取用户名
	claims, exists := c.Get("claims")
	if !exists {
		utils.UnauthorizedResponse(c, "token claims not found")
		return
	}

	username := claims.(jwt.MapClaims)["username"].(string)

	// 解析请求数据
	var profileData UserProfileRequest
	if err := c.ShouldBindJSON(&profileData); err != nil {
		utils.ErrorResponse(c, err.Error())
		return
	}

	// 更新用户信息
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		utils.NotFoundResponse(c, "user not found")
		return
	}

	// 只更新提供的字段
	if profileData.Name != "" {
		user.Name = profileData.Name
	}
	if profileData.Phone != "" {
		user.Phone = profileData.Phone
	}

	if err := database.DB.Save(&user).Error; err != nil {
		utils.ServerErrorResponse(c, "failed to update user profile")
		return
	}

	// 返回更新后的用户信息
	utils.SuccessResponse(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"name":     user.Name,
		"phone":    user.Phone,
	})
}
