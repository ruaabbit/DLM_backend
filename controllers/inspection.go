package controllers

import (
	"net/http"
	"strconv"

	"DLM_backend/models"
	"DLM_backend/services"

	"github.com/gin-gonic/gin"
)

// CreateInspection 处理新增点检记录请求
func CreateInspection(c *gin.Context) {
	var record models.InspectionRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := services.CreateInspectionRecord(&record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create record"})
		return
	}
	c.JSON(http.StatusOK, created)
}

// GetInspections 处理查询点检记录请求
func GetInspections(c *gin.Context) {
	records, err := services.GetInspectionRecords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get records"})
		return
	}
	c.JSON(http.StatusOK, records)
}

// UpdateInspection 处理更新点检记录请求
func UpdateInspection(c *gin.Context) {
	var record models.InspectionRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := services.UpdateInspectionRecord(&record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update record"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteInspection 处理删除点检记录请求
func DeleteInspection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := services.DeleteInspectionRecord(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete record"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "record deleted"})
}
