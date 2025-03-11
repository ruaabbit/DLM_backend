package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"DLM_backend/database"
	"DLM_backend/models"
	"DLM_backend/services"
	"DLM_backend/utils"

	"github.com/dgrijalva/jwt-go"
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

	// 从JWT获取用户信息
	claims, exists := c.Get("claims")
	if !exists {
		utils.UnauthorizedResponse(c, "token claims not found")
		return
	}
	username := claims.(jwt.MapClaims)["username"].(string)

	// 查询用户ID
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		utils.NotFoundResponse(c, "user not found")
		return
	}

	// 处理数组字段转换为JSON
	pinStatusJSON, _ := json.Marshal(requestData.PinStatus)
	mainWallStatusJSON, _ := json.Marshal(requestData.MainWallStatus)
	warehouseFoundationJSON, _ := json.Marshal(requestData.WarehouseFoundation)

	// 转换为模型
	record := models.InspectionRecord{
		UserID:                         user.ID, // 添加用户ID
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

// GetInspections 处理查询点检记录请求，支持分页和过滤
func GetInspections(c *gin.Context) {
	// 获取分页参数，默认第1页，每页10条
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 确保参数有效
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // 限制pageSize范围，防止请求过大数据
	}

	// 构建过滤条件
	filters := make(map[string]interface{})

	// 添加常规字段过滤
	if unit := c.Query("unit"); unit != "" {
		filters["unit"] = unit
	}
	if warehouse := c.Query("warehouse_number"); warehouse != "" {
		filters["warehouse_number"] = warehouse
	}
	if position := c.Query("grain_door_position"); position != "" {
		filters["grain_door_position"] = position
	}
	if caretaker := c.Query("caretaker"); caretaker != "" {
		filters["caretaker"] = caretaker
	}
	if signature := c.Query("signature"); signature != "" {
		filters["signature"] = signature
	}

	// 添加日期范围过滤
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			filters["start_date"] = startDate
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			// 设置为当天结束时间
			endDate = endDate.Add(24*time.Hour - time.Second)
			filters["end_date"] = endDate
		}
	}

	// 添加状况类型过滤
	if deformation := c.Query("deformation_crack"); deformation != "" {
		filters["deformation_crack"] = deformation
	}
	if closure := c.Query("closure_status"); closure != "" {
		filters["closure_status"] = closure
	}
	if safety := c.Query("safety_rope_installed"); safety != "" {
		filters["safety_rope_installed"] = safety
	}

	// 添加JSON字段过滤
	if pinStatus := c.Query("pin_status"); pinStatus != "" {
		filters["pin_status"] = pinStatus
	}
	if mainWallStatus := c.Query("main_wall_status"); mainWallStatus != "" {
		filters["main_wall_status"] = mainWallStatus
	}
	if warehouseFoundation := c.Query("warehouse_foundation"); warehouseFoundation != "" {
		filters["warehouse_foundation"] = warehouseFoundation
	}

	// 关键字搜索 (在多个字段中查找)
	if keyword := c.Query("keyword"); keyword != "" {
		// 这里需要修改服务层函数来支持关键字搜索
		// 本例中暂时使用备注字段进行模糊查询
		filters["keyword"] = keyword
	}

	// 调用服务层获取带过滤的分页数据
	total, records, err := services.GetInspectionRecordsWithFilters(page, pageSize, filters)
	if err != nil {
		utils.ErrorResponse(c, "failed to get records")
		return
	}

	// 计算总页数
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	utils.SuccessResponse(c, gin.H{
		"records": records,
		"pagination": gin.H{
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": totalPages,
		},
	})
}

// GetUserInspections 获取当前登录用户的点检记录
func GetUserInspections(c *gin.Context) {
	// 从JWT中获取用户名
	claims, exists := c.Get("claims")
	if !exists {
		utils.UnauthorizedResponse(c, "token claims not found")
		return
	}
	username := claims.(jwt.MapClaims)["username"].(string)

	// 查询用户ID
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		utils.NotFoundResponse(c, "user not found")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 确保参数有效
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 构建过滤条件，加入用户ID
	filters := make(map[string]interface{})
	filters["user_id"] = user.ID

	// 添加日期范围过滤
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			filters["start_date"] = startDate
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			// 设置为当天结束时间
			endDate = endDate.Add(24*time.Hour - time.Second)
			filters["end_date"] = endDate
		}
	}

	// 调用服务层获取带过滤的分页数据
	total, records, err := services.GetInspectionRecordsWithFilters(page, pageSize, filters)
	if err != nil {
		utils.ErrorResponse(c, "failed to get records")
		return
	}

	// 计算总页数
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	// 获取用户累计打卡次数
	var totalCount int64
	if err := database.DB.Model(&models.InspectionRecord{}).Where("user_id = ?", user.ID).Count(&totalCount).Error; err != nil {
		utils.ServerErrorResponse(c, "failed to get total count")
		return
	}

	// 获取本月打卡次数
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)
	var monthlyCount int64
	if err := database.DB.Model(&models.InspectionRecord{}).
		Where("user_id = ? AND inspection_time BETWEEN ? AND ?",
			user.ID, startOfMonth, endOfMonth).
		Count(&monthlyCount).Error; err != nil {
		utils.ServerErrorResponse(c, "failed to get monthly count")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"records": records,
		"pagination": gin.H{
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": totalPages,
		},
		"stats": gin.H{
			"monthlyCount": monthlyCount, // 本月打卡次数
			"totalCount":   totalCount,   // 累计打卡次数
		},
	})
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
