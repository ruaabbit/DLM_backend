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
		// 点检记录相关接口
		authorized.POST("/inspection", controllers.CreateInspection)
		authorized.GET("/inspection", controllers.GetInspections)
		authorized.PUT("/inspection", controllers.UpdateInspection)
		authorized.DELETE("/inspection/:id", controllers.DeleteInspection)

		// 获取当前登录用户的点检记录
		authorized.GET("/user/inspections", controllers.GetUserInspections)

		// 用户个人信息相关接口
		authorized.GET("/profile", controllers.GetUserProfile)
		authorized.PUT("/profile", controllers.UpdateUserProfile)

		// 图片上传接口
		authorized.POST("/upload/image", controllers.UploadImage)

	}

	r.Static("/images", "./uploads/images")

	return r
}
