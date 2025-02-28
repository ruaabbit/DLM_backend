package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"DLM_backend/models"
	"DLM_backend/services"
	"DLM_backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// InspectionRequest 用于处理前端传来的点检记录请求
type InspectionRequest struct {
	Unit                           string         `json:"unit" binding:"required"`                  // 单位
	WarehouseNumber                string         `json:"warehouse_number" binding:"required"`      // 仓号
	GrainDoorPosition              string         `json:"grain_door_position" binding:"required"`   // 挡粮门位置
	Caretaker                      string         `json:"caretaker" binding:"required"`             // 保管责任人
	InspectionTime                 time.Time      `json:"inspection_time" binding:"required"`       // 检查时间
	DeformationCrack               string         `json:"deformation_crack" binding:"required"`     // 挡粮门变形和裂痕情况
	DeformationCrackDescription    string         `json:"deformation_crack_description"`            // 挡粮门变形和裂痕情况说明
	ClosureStatus                  string         `json:"closure_status" binding:"required"`        // 闭合情况
	ClosureDescription             string         `json:"closure_description"`                      // 闭合情况说明
	PinStatus                      datatypes.JSON `json:"pin_status" binding:"required"`            // 栓销状况
	PinDescription                 string         `json:"pin_description"`                          // 栓销状况说明
	MainWallStatus                 datatypes.JSON `json:"main_wall_status" binding:"required"`      // 主体墙状况
	MainWallDescription            string         `json:"main_wall_description"`                    // 主体墙状况说明
	WarehouseFoundation            datatypes.JSON `json:"warehouse_foundation" binding:"required"`  // 仓门地基状况
	WarehouseFoundationDescription string         `json:"warehouse_foundation_description"`         // 仓门地基状况说明
	SafetyRopeInstalled            string         `json:"safety_rope_installed" binding:"required"` // 安全绳（带）系留装置
	SafetyRopeDescription          string         `json:"safety_rope_description"`                  // 安全绳（带）系留装置说明
	Remarks                        string         `json:"remarks"`                                  // 补充说明
	Signature                      string         `json:"signature" binding:"required"`             // 责任人签名
	ContactNumber                  string         `json:"contact_number" binding:"required"`        // 联系电话
	Images                         datatypes.JSON `json:"images"`                                   // 图片列表
}

// CreateInspection 处理新增点检记录请求
func CreateInspection(c *gin.Context) {
	var requestData InspectionRequest
	if err := c.ShouldBindJSON(&requestData); err != nil {
		utils.ErrorResponse(c, err.Error())
		return
	}

	// 处理数组字段转换为JSON
	pinStatusJSON, _ := json.Marshal(requestData.PinStatus)
	mainWallStatusJSON, _ := json.Marshal(requestData.MainWallStatus)
	warehouseFoundationJSON, _ := json.Marshal(requestData.WarehouseFoundation)

	// 转换为模型
	record := models.InspectionRecord{
		Unit:                           requestData.Unit,
		WarehouseNumber:                requestData.WarehouseNumber,
		GrainDoorPosition:              requestData.GrainDoorPosition,
		Caretaker:                      requestData.Caretaker,
		InspectionTime:                 requestData.InspectionTime,
		DeformationCrack:               requestData.DeformationCrack,
		DeformationCrackDescription:    requestData.DeformationCrackDescription,
		ClosureStatus:                  requestData.ClosureStatus,
		ClosureDescription:             requestData.ClosureDescription,
		PinStatus:                      pinStatusJSON,
		PinDescription:                 requestData.PinDescription,
		MainWallStatus:                 mainWallStatusJSON,
		MainWallDescription:            requestData.MainWallDescription,
		WarehouseFoundation:            warehouseFoundationJSON,
		WarehouseFoundationDescription: requestData.WarehouseFoundationDescription,
		SafetyRopeInstalled:            requestData.SafetyRopeInstalled,
		SafetyRopeDescription:          requestData.SafetyRopeDescription,
		Remarks:                        requestData.Remarks,
		Signature:                      requestData.Signature,
		ContactNumber:                  requestData.ContactNumber,
	}

	// 处理图片数据
	if len(requestData.Images) > 0 {
		// 转换图片数组为JSON
		imagesJSON, _ := json.Marshal(requestData.Images)
		record.Images = imagesJSON
	}

	created, err := services.CreateInspectionRecord(&record)
	if err != nil {
		utils.ErrorResponse(c, "failed to create record")
		return
	}
	utils.SuccessResponse(c, created)
}

// GetInspections 处理查询点检记录请求
func GetInspections(c *gin.Context) {
	records, err := services.GetInspectionRecords()
	if err != nil {
		utils.ErrorResponse(c, "failed to get records")
		return
	}
	utils.SuccessResponse(c, records)
}

// UpdateInspection 处理更新点检记录请求
func UpdateInspection(c *gin.Context) {
	var requestData InspectionRequest
	if err := c.ShouldBindJSON(&requestData); err != nil {
		utils.ErrorResponse(c, err.Error())
		return
	}

	// 获取记录ID
	id, exists := c.GetQuery("id")
	if !exists {
		utils.ErrorResponse(c, "missing record id")
		return
	}

	recordID, err := strconv.Atoi(id)
	if err != nil {
		utils.ErrorResponse(c, "invalid record id")
		return
	}

	// 处理数组字段转换为JSON
	pinStatusJSON, _ := json.Marshal(requestData.PinStatus)
	mainWallStatusJSON, _ := json.Marshal(requestData.MainWallStatus)
	warehouseFoundationJSON, _ := json.Marshal(requestData.WarehouseFoundation)

	// 转换为模型
	record := models.InspectionRecord{
		ID:                             recordID,
		Unit:                           requestData.Unit,
		WarehouseNumber:                requestData.WarehouseNumber,
		GrainDoorPosition:              requestData.GrainDoorPosition,
		Caretaker:                      requestData.Caretaker,
		InspectionTime:                 requestData.InspectionTime,
		DeformationCrack:               requestData.DeformationCrack,
		DeformationCrackDescription:    requestData.DeformationCrackDescription,
		ClosureStatus:                  requestData.ClosureStatus,
		ClosureDescription:             requestData.ClosureDescription,
		PinStatus:                      pinStatusJSON,
		PinDescription:                 requestData.PinDescription,
		MainWallStatus:                 mainWallStatusJSON,
		MainWallDescription:            requestData.MainWallDescription,
		WarehouseFoundation:            warehouseFoundationJSON,
		WarehouseFoundationDescription: requestData.WarehouseFoundationDescription,
		SafetyRopeInstalled:            requestData.SafetyRopeInstalled,
		SafetyRopeDescription:          requestData.SafetyRopeDescription,
		Remarks:                        requestData.Remarks,
		Signature:                      requestData.Signature,
		ContactNumber:                  requestData.ContactNumber,
	}

	// 处理图片数据
	if len(requestData.Images) > 0 {
		// 转换图片数组为JSON
		imagesJSON, _ := json.Marshal(requestData.Images)
		record.Images = imagesJSON
	}

	updated, err := services.UpdateInspectionRecord(&record)
	if err != nil {
		utils.ErrorResponse(c, "failed to update record")
		return
	}
	utils.SuccessResponse(c, updated)
}

// DeleteInspection 处理删除点检记录请求
func DeleteInspection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, "invalid id")
		return
	}
	if err := services.DeleteInspectionRecord(id); err != nil {
		utils.ErrorResponse(c, "failed to delete record")
		return
	}
	utils.SuccessResponse(c, gin.H{"message": "record deleted"})
}
