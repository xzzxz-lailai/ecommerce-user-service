package partnerauth

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"user_service/model"
	"user_service/pkg"
	partnerRepo "user_service/repo/partner"
)

type partnerUserInfo struct {
	ExternalUserID string `json:"external_user_id"`
	Nickname       string `json:"nickname"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Avatar         string `json:"avatar"`
}

type partnerUserInfoRequest struct {
	AppKey   string `json:"app_key"`
	AuthCode string `json:"auth_code"`
}

type partnerUserInfoResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    partnerUserInfo `json:"data"`
}

// PartnerLogin 合作方用户免登录
func PartnerLogin(ctx context.Context, req *model.PartnerLoginRequest) (*model.PartnerLoginResponse, error) {
	// 查询合作方配置。
	partner, err := partnerRepo.FindPartnerByCode(ctx, req.PartnerCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("合作方不存在")
		}
		return nil, err
	}

	if partner.Status != 1 {
		return nil, errors.New("合作方已禁用")
	}

	// 通过 auth_code 换取合作方用户信息。
	userInfo, err := fetchPartnerUserInfo(ctx, partner, req.AuthCode)
	if err != nil {
		return nil, err
	}
	if userInfo.ExternalUserID == "" {
		return nil, errors.New("合作方用户ID为空")
	}

	// 查询或创建平台用户映射。
	partnerUser, err := partnerRepo.FindPartnerUser(ctx, partner.PartnerID, userInfo.ExternalUserID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		partnerUser = &model.PartnerUser{
			PartnerID:      partner.PartnerID,
			ExternalUserID: userInfo.ExternalUserID,
			Nickname:       userInfo.Nickname,
			Phone:          userInfo.Phone,
			Email:          userInfo.Email,
			Avatar:         userInfo.Avatar,
			Status:         1,
		}
		if err := partnerRepo.CreatePartnerUser(ctx, partnerUser); err != nil {
			return nil, err
		}
	} else {
		if err := partnerRepo.UpdatePartnerUserLoginTime(ctx, partnerUser.PartnerUserID); err != nil {
			return nil, err
		}
	}

	// 生成我们系统自己的登录 token。
	token, err := pkg.GeneratePartnerUserToken(partner.PartnerID, partnerUser.PartnerUserID)
	if err != nil {
		return nil, errors.New("登录失败，请稍后重试")
	}

	return &model.PartnerLoginResponse{
		Token:         token,
		PartnerUserID: partnerUser.PartnerUserID,
		Nickname:      partnerUser.Nickname,
	}, nil
}

// fetchPartnerUserInfo 通过 auth_code 获取合作方用户信息。
func fetchPartnerUserInfo(ctx context.Context, partner *model.Partner, authCode string) (*partnerUserInfo, error) {
	reqBody := partnerUserInfoRequest{
		AppKey:   partner.AppKey,
		AuthCode: authCode,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	// 调用合作方用户信息接口。
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, partner.UserinfoURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, errors.New("调用合作方用户信息接口失败")
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("合作方用户信息接口状态异常：%d", httpResp.StatusCode)
	}

	var resp partnerUserInfoResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, errors.New("解析合作方用户信息失败")
	}
	if resp.Code != 0 {
		return nil, errors.New(resp.Message)
	}

	return &resp.Data, nil
}
