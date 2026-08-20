package partnerauth

import (
	"net/http"
	"user_service/model"
	"user_service/pkg"
	"user_service/service/partnerauth"

	"github.com/gin-gonic/gin"
)

// PartnerLogin 合作方用户免登录。
func PartnerLogin(c *gin.Context) {
	var req model.PartnerLoginRequest

	// 接收前端 JSON。
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 调用 Service。
	resp, err := partnerauth.PartnerLogin(c.Request.Context(), &req)
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	pkg.Success(c, "登录成功", resp)
}
