package models

import (
	"gorm.io/datatypes"
)

// InspectionRecord 定义挡粮门点检记录模型
type InspectionRecord struct {
	ID             int            `json:"id" gorm:"primaryKey"`            // 主键ID
	Unit           string         `json:"unit" gorm:"not null"`            // 单位
	Workshop       string         `json:"workshop" gorm:"not null"`        // 车间名称
	SerialNumber   string         `json:"serial_number" gorm:"not null"`   // 编号
	Date           string         `json:"date" gorm:"not null"`            // 检查日期
	InspectionItem string         `json:"inspection_item" gorm:"not null"` // 检查项目
	Status         string         `json:"status" gorm:"not null"`          // 检查情况
	Images         datatypes.JSON `json:"images"`                          // 图片列表
	Inspector      string         `json:"inspector" gorm:"not null"`       // 记录人员
	Manager        string         `json:"manager" gorm:"not null"`         // 保管责任人
}
