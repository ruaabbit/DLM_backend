package routers

import (
	"DLM_backend/controllers"
	"DLM_backend/utils"

	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化 Gin 路由，并注册接口
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 公共路由
	r.POST("/login", controllers.Login)

	// JWT 鉴权的路由组
	authorized := r.Group("/", utils.JWTAuthMiddleware())
	{
		authorized.POST("/inspection", controllers.CreateInspection)
		authorized.GET("/inspection", controllers.GetInspections)
		authorized.PUT("/inspection", controllers.UpdateInspection)
		authorized.DELETE("/inspection/:id", controllers.DeleteInspection)
	}

	return r
}
