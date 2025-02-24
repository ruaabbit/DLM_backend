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
