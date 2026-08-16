package router

import (
	"github.com/gin-gonic/gin"
	"user_service/handler/auth"
)

// AuthRoutes 认证路由
func AuthRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", auth.Register) // 注册
		authGroup.POST("/login", auth.Login)       // 登陆
	}
}
