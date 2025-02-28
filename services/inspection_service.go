package services

import (
	"DLM_backend/database"
	"DLM_backend/models"
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
