package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse 返回成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

// ErrorResponse 返回错误响应
func ErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": message})
}

// UnauthorizedResponse 返回未授权的错误响应
func UnauthorizedResponse(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": message})
}

// ServerErrorResponse 返回服务器内部错误响应
func ServerErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": message})
}

// NotFoundResponse 返回资源未找到的响应
func NotFoundResponse(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{"success": false, "error": message})
}
