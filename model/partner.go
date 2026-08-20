package model

import "time"

// Partner 合作方配置，对应 partners 表。
type Partner struct {
	PartnerID   int64     `db:"partner_id"`   // 合作方ID，系统内部唯一标识
	PartnerCode string    `db:"partner_code"` // 合作方编码，用于接口对接识别，例如 moushi_app
	PartnerName string    `db:"partner_name"` // 合作方名称，例如某市App
	AppKey      string    `db:"app_key"`      // 分配给合作方的AppKey，用于接口身份识别
	AppSecret   string    `db:"app_secret"`   // 分配给合作方的密钥，用于签名校验或调用对方接口
	UserinfoURL string    `db:"userinfo_url"` // 合作方用户信息接口地址，用于通过auth_code换取用户信息
	Status      int8      `db:"status"`       // 状态：1启用，0禁用
	Remark      string    `db:"remark"`       // 备注信息
	CreatedAt   time.Time `db:"created_at"`   // 创建时间
	UpdatedAt   time.Time `db:"updated_at"`   // 更新时间
}

// PartnerUser 合作方用户映射，对应 partner_users 表。
type PartnerUser struct {
	PartnerUserID  int64     `db:"partner_user_id"`  // 平台用户映射ID，系统内部唯一标识
	PartnerID      int64     `db:"partner_id"`       // 合作方ID，对应 partners.partner_id
	ExternalUserID string    `db:"external_user_id"` // 合作方系统里的用户唯一ID
	Nickname       string    `db:"nickname"`         // 用户昵称，由合作方返回，可为空
	Phone          string    `db:"phone"`            // 用户手机号，由合作方返回，可为空
	Email          string    `db:"email"`            // 用户邮箱，由合作方返回，可为空
	Avatar         string    `db:"avatar"`           // 用户头像地址，由合作方返回，可为空
	Status         int8      `db:"status"`           // 状态：1正常，0禁用
	LastLoginAt    time.Time `db:"last_login_at"`    // 最近一次免登录进入商城的时间
	CreatedAt      time.Time `db:"created_at"`       // 创建时间
	UpdatedAt      time.Time `db:"updated_at"`       // 更新时间
}

// PartnerLoginRequest 合作方用户免登录请求。
type PartnerLoginRequest struct {
	PartnerCode string `json:"partner_code" binding:"required"` // 合作方编码
	AuthCode    string `json:"auth_code" binding:"required"`    // 合作方生成的一次性登录code
}

// PartnerLoginResponse 合作方用户免登录响应。
type PartnerLoginResponse struct {
	Token         string `json:"token"`           // 我们系统签发的JWT
	PartnerUserID int64  `json:"partner_user_id"` // 我们系统内的平台用户ID
	Nickname      string `json:"nickname"`        // 用户昵称
}
