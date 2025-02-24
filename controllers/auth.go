package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "DLM_backend/services"
)

// Login 用于处理用户登录请求
func Login(c *gin.Context) {
    var loginData struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    token, err := services.AuthenticateUser(loginData.Username, loginData.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}