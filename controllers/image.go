package controllers

import (
	"os"
	"strconv"
	"time"

	"DLM_backend/utils"

	"github.com/gin-gonic/gin"
)

// UploadImage 处理图片上传请求
func UploadImage(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, "未能获取上传文件: "+err.Error())
		return
	}

	// 创建存储目录（如果不存在）
	uploadDir := "uploads/images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.ServerErrorResponse(c, "创建目录失败: "+err.Error())
		return
	}

	// 生成唯一文件名（防止文件名冲突）
	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + file.Filename
	filepath := uploadDir + "/" + filename

	// 保存文件
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		utils.ServerErrorResponse(c, "保存文件失败: "+err.Error())
		return
	}

	// 返回文件访问路径
	fileURL := "/images/" + filename
	utils.SuccessResponse(c, gin.H{
		"url":      fileURL,
		"filename": filename,
	})
}
