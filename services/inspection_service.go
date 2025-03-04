package services

import (
	"DLM_backend/database"
	"DLM_backend/models"
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

	// 应用过滤条件
	for key, value := range filters {
		// 如果值是字符串并且不为空，则添加模糊查询
		if strVal, ok := value.(string); ok && strVal != "" {
			query = query.Where(key+" LIKE ?", "%"+strVal+"%")
		} else if value != nil {
			// 对于其他非空值，使用精确匹配
			query = query.Where(key+" = ?", value)
		}
	}

	// 如果有日期范围过滤
	if startDate, ok := filters["start_date"].(time.Time); ok {
		query = query.Where("inspection_time >= ?", startDate)
		delete(filters, "start_date") // 删除已处理的特殊条件
	}
	if endDate, ok := filters["end_date"].(time.Time); ok {
		query = query.Where("inspection_time <= ?", endDate)
		delete(filters, "end_date") // 删除已处理的特殊条件
	}

	// 完整的关键字搜索功能
	if keyword, ok := filters["keyword"].(string); ok && keyword != "" {
		query = query.Where(
			"unit LIKE ? OR warehouse_number LIKE ? OR grain_door_position LIKE ? OR "+
				"caretaker LIKE ? OR remarks LIKE ? OR signature LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
		)
		delete(filters, "keyword") // 从过滤条件中移除已处理的关键字
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
