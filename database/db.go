// language: go
package database

import (
	"log"

	"DLM_backend/config"
	"DLM_backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB 导出全局数据库实例
var DB *gorm.DB

// InitDB 根据配置选择数据库驱动并初始化 GORM 连接
func InitDB(cfg *config.Config) *gorm.DB {
	var dialector gorm.Dialector
	if cfg.DBDriver == "sqlite3" {
		dialector = sqlite.Open(cfg.DSN())
	} else if cfg.DBDriver == "mysql" {
		dialector = mysql.Open(cfg.DSN())
	} else {
		log.Fatalf("unknown db driver: %s", cfg.DBDriver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移模型
	if err := db.AutoMigrate(&models.User{}, &models.InspectionRecord{}); err != nil {
		log.Fatalf("failed to migrate models: %v", err)
	}

	// 导出数据库实例
	DB = db
	return db
}
