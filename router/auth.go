package router

import (
	"github.com/gin-gonic/gin"
	"user_service/handler/auth"
	"user_service/handler/partnerauth"
)

// AuthRoutes 认证路由
func AuthRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register/code", auth.SendRegisterEmailCode) // 发送注册邮箱验证码
		authGroup.POST("/register", auth.Register)                   // 注册

		authGroup.POST("/login", auth.Login) // 登陆

		authGroup.POST("/partner/login", partnerauth.PartnerLogin) // 合作方用户免登录

		authGroup.POST("/forgot-password/code", auth.SendForgetPasswordEmailCode) // 发送忘记密码邮箱验证码
		authGroup.POST("/forgot-password/reset", auth.ForgetPassword)             // 忘记密码
	}
}
