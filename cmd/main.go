package main

import (
	"log"

	"DLM_backend/config"
	"DLM_backend/database"
	"DLM_backend/routers"
	"DLM_backend/utils"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 设置 JWT 密钥
	utils.SetJWTSecret(cfg.JWTSecret)

	// 初始化数据库连接
	db := database.InitDB(cfg)
	_ = db // 后续可使用 db 进行数据库操作

	// 初始化路由
	r := routers.SetupRouter()

	// 启动服务
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
