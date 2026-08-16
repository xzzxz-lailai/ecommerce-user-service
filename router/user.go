package router

import (
	"github.com/gin-gonic/gin"
	"user_service/handler/auth"
)

// UserRoutes 用户路由
func UserRoutes(r *gin.RouterGroup) {
	userGroup := r.Group("/user")
	{
		userGroup.PUT("/password", auth.ChangePassword)
	}
}
