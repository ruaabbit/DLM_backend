package services

import (
	"DLM_backend/database"
	"DLM_backend/models"
	"fmt"
	"strings"
	"time"
)

// CreateInspectionRecord 新建点检记录
func CreateInspectionRecord(record *models.InspectionRecord) (*models.InspectionRecord, error) {
	if err := database.DB.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

// GetInspectionRecords 获取所有点检记录（兼容原有API）
func GetInspectionRecords() ([]models.InspectionRecord, error) {
	// 使用一个足够大的值获取所有数据
	_, records, err := GetInspectionRecordsWithPagination(1, 1000)
	return records, err
}

// GetInspectionRecordsWithPagination 获取分页的点检记录
func GetInspectionRecordsWithPagination(page, pageSize int) (int64, []models.InspectionRecord, error) {
	var records []models.InspectionRecord
	var total int64

	// 获取总记录数
	if err := database.DB.Model(&models.InspectionRecord{}).Count(&total).Error; err != nil {
		return 0, nil, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询分页数据
	if err := database.DB.Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		return 0, nil, err
	}

	return total, records, nil
}

// GetInspectionRecordsWithFilters 获取带条件过滤的分页点检记录
func GetInspectionRecordsWithFilters(page, pageSize int, filters map[string]interface{}) (int64, []models.InspectionRecord, error) {
	var records []models.InspectionRecord
	var total int64
	query := database.DB.Model(&models.InspectionRecord{})

	// 提前复制filters，避免修改原始filters
	workingFilters := make(map[string]interface{})
	for k, v := range filters {
		workingFilters[k] = v
	}

	// 处理关键字搜索
	if keyword, ok := workingFilters["keyword"].(string); ok && keyword != "" {
		query = query.Where(
			"unit LIKE ? OR warehouse_number LIKE ? OR grain_door_position LIKE ? OR "+
				"caretaker LIKE ? OR remarks LIKE ? OR signature LIKE ? OR "+
				"deformation_crack_description LIKE ? OR closure_description LIKE ? OR "+
				"pin_description LIKE ? OR main_wall_description LIKE ? OR "+
				"warehouse_foundation_description LIKE ? OR safety_rope_description LIKE ? OR "+
				"contact_number LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
			"%"+keyword+"%",
		)
		delete(workingFilters, "keyword")
	}

	// 处理日期范围过滤
	if startDateStr, ok := workingFilters["start_date"].(string); ok && startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			query = query.Where("inspection_time >= ?", startDate)
		}
		delete(workingFilters, "start_date")
	} else if startDate, ok := workingFilters["start_date"].(time.Time); ok {
		query = query.Where("inspection_time >= ?", startDate)
		delete(workingFilters, "start_date")
	}

	if endDateStr, ok := workingFilters["end_date"].(string); ok && endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			// 设置到当天结束
			endDate = endDate.Add(24*time.Hour - time.Second)
			query = query.Where("inspection_time <= ?", endDate)
		}
		delete(workingFilters, "end_date")
	} else if endDate, ok := workingFilters["end_date"].(time.Time); ok {
		query = query.Where("inspection_time <= ?", endDate)
		delete(workingFilters, "end_date")
	}

	// 处理JSON字段的过滤（支持多值查询，用逗号分隔）
	jsonFields := []string{"pin_status", "main_wall_status", "warehouse_foundation"}
	for _, fieldName := range jsonFields {
		if value, ok := workingFilters[fieldName].(string); ok && value != "" {
			values := strings.Split(value, ",")
			for _, v := range values {
				v = strings.TrimSpace(v)
				// 对每个值添加一个EXISTS条件
				query = query.Where(fmt.Sprintf("EXISTS (SELECT 1 FROM json_each(%s) WHERE value = ?)", fieldName), v)
			}
			delete(workingFilters, fieldName)
		}
	}

	// 应用其余的过滤条件
	for key, value := range workingFilters {
		// 如果值是字符串并且不为空，则添加模糊查询
		if strVal, ok := value.(string); ok && strVal != "" {
			query = query.Where(key+" LIKE ?", "%"+strVal+"%")
		} else if value != nil {
			// 对于其他非空值，使用精确匹配
			query = query.Where(key+" = ?", value)
		}
	}

	// 获取总记录数
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询分页数据
	if err := query.Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		return 0, nil, err
	}

	return total, records, nil
}

// UpdateInspectionRecord 更新点检记录
func UpdateInspectionRecord(record *models.InspectionRecord) (*models.InspectionRecord, error) {
	if err := database.DB.Save(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

// DeleteInspectionRecord 删除点检记录
func DeleteInspectionRecord(id int) error {
	if err := database.DB.Delete(&models.InspectionRecord{}, id).Error; err != nil {
		return err
	}
	return nil
}
