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

// GetInspectionRecords 获取所有点检记录
func GetInspectionRecords() ([]models.InspectionRecord, error) {
	var records []models.InspectionRecord
	if err := database.DB.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
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
