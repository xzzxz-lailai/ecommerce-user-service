package middleware

import (
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
			c.Abort() // 终止请求，不再往下执行
			return
		}

		////检查格式是否是 "Bearer xxx"
		//parts := strings.SplitN(token, " ", 2)
		//if len(parts) != 2 || parts[0] != "Bearer" {
		//	c.JSON(http.StatusUnauthorized, handler.Response{
		//		Code: 401,
		//		Msg:  "Token 格式错误",
		//		Data: nil,
		//	})
		//	c.Abort()
		//	return
		//}

		// 2. 解析 Token
		claims, err := pkg.ParseToken(token)
		if err != nil {
			pkg.Error(c, 401, "请先登录")
			c.Abort()
			return
		}

		// 3. 把用户信息存入上下文，后续 handler 直接取用
		c.Set("userID", claims.UserID)

		c.Next()
	}
}
