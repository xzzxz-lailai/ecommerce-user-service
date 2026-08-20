package router

import (
	"github.com/gin-gonic/gin"
	"user_service/middleware"
)

func Router() *gin.Engine {
	r := gin.Default()

	// CORS 中间件
	r.Use(middleware.CORS())

	// API v1 路由组
	api := r.Group("/api/v1")
	{
		// 公开路由(无需鉴权)
		AccountRoutes(api)

		// 公司内部账号路由(需要鉴权)
		private := api.Group("", middleware.RequireUser())
		{
			UserRoutes(private)
		}
	}

	return r
}
