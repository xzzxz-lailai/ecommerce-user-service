package auth

import (
	"net/http"
	"user_service/model"
	"user_service/pkg"
	auth "user_service/service/auth"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req model.RegisterRequest

	//  接收前端 JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 调用 Service
	err := auth.Register(c.Request.Context(), &req)
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 返回结果
	pkg.Success(c, "注册成功", nil)
}

// Login 用户登录
func Login(c *gin.Context) {
	var req model.LoginRequest

	//  接收前端 JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 调用 Service
	resp, err := auth.Login(c.Request.Context(), &req)
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 返回结果
	pkg.Success(c, "登录成功", resp)
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	var req model.ChangePasswordRequest
	// 接收请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	// 2. 从 JWT 中获取 userID
	userIDValue, exists := c.Get("userID")
	if !exists {
		pkg.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	userID := userIDValue.(int64)
	//  调用 Service
	err := auth.ChangePassword(c.Request.Context(), userID, &req)
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	//  修改成功
	pkg.Success(c, "密码修改成功", nil)
}
