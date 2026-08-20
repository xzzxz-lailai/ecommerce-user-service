package router

import (
	"github.com/gin-gonic/gin"
	"user_service/handler/account"
	"user_service/handler/partnerauth"
)

// AccountRoutes 账号路由
func AccountRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register/code", account.SendRegisterEmailCode) // 发送注册邮箱验证码
		authGroup.POST("/register", account.Register)                   // 注册

		authGroup.POST("/login", account.Login) // 登陆

		authGroup.POST("/partner/login", partnerauth.PartnerLogin) // 合作方用户免登录(包含获取合作方用户信息)

		authGroup.POST("/forgot-password/code", account.SendForgetPasswordEmailCode) // 发送忘记密码邮箱验证码
		authGroup.POST("/forgot-password/reset", account.ForgetPassword)             // 忘记密码
	}
}
