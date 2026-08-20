package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"user_service/pkg"
)

// Authorization 鉴权中间件
func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取 Token
		token := c.GetHeader("Authorization")
		if token == "" {
			pkg.Error(c, 401, "请先登录")
			c.Abort()
			return
		}

		// 检查格式是否是 "Bearer xxx"。
		parts := strings.SplitN(token, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			pkg.Error(c, 401, "Token 格式错误")
			c.Abort()
			return
		}

		// 2. 解析 Token
		claims, err := pkg.ParseToken(parts[1])
		if err != nil {
			pkg.Error(c, 401, "请先登录")
			c.Abort()
			return
		}

		// 3. 把用户信息存入上下文，后续 handler 直接取用
		c.Set("subjectType", claims.SubjectType)
		c.Set("userID", claims.UserID)
		c.Set("partnerID", claims.PartnerID)
		c.Set("partnerUserID", claims.PartnerUserID)

		c.Next()
	}
}

// RequireUser 只允许公司内部账号访问。
func RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		Authorization()(c)
		if c.IsAborted() {
			return
		}

		if c.GetString("subjectType") != "user" {
			pkg.Error(c, 403, "无权限访问")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePartnerUser 只允许合作方用户访问。
func RequirePartnerUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		Authorization()(c)
		if c.IsAborted() {
			return
		}

		if c.GetString("subjectType") != "partner_user" {
			pkg.Error(c, 403, "无权限访问")
			c.Abort()
			return
		}

		c.Next()
	}
}
