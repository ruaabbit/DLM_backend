package models

import (
	"time"

	"gorm.io/datatypes"
)

// InspectionRecord 定义挡粮门点检记录模型
type InspectionRecord struct {
	ID                  int            `json:"id" gorm:"primaryKey"`                  // 主键ID
	Unit                string         `json:"unit" gorm:"not null"`                  // 单位
	WarehouseNumber     string         `json:"warehouse_number" gorm:"not null"`      // 仓号
	GrainDoorPosition   string         `json:"grain_door_position" gorm:"not null"`   // 挡粮门位置
	Caretaker           string         `json:"caretaker" gorm:"not null"`             // 保管责任人
	InspectionTime      time.Time      `json:"inspection_time" gorm:"not null"`       // 检查时间
	DeformationCrack    string         `json:"deformation_crack" gorm:"not null"`     // 挡粮门变形和裂痕情况 (无异常/有异常)
	ClosureStatus       string         `json:"closure_status" gorm:"not null"`        // 闭合情况 (闭合良好/闭合不严)
	PinStatus           datatypes.JSON `json:"pin_status" gorm:"type:json"`           // 栓销状况 (正常/松动/变形/缺失)
	MainWallStatus      datatypes.JSON `json:"main_wall_status" gorm:"type:json"`     // 主体墙状况 (正常/破损/有裂缝)
	WarehouseFoundation datatypes.JSON `json:"warehouse_foundation" gorm:"type:json"` // 仓门地基状况 (正常/冻胀/下沉/塌陷/裂痕)
	SafetyRopeInstalled string         `json:"safety_rope_installed" gorm:"not null"` // 安全绳（带）系留装置 (已安装/未安装)
	Remarks             string         `json:"remarks" gorm:"type:text"`              // 补充说明
	Signature           string         `json:"signature" gorm:"not null"`             // 责任人签名
	ContactNumber       string         `json:"contact_number" gorm:"not null"`        // 联系电话
	Images              datatypes.JSON `json:"images"`                                // 图片列表
}
