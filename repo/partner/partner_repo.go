package partner

import (
	"context"
	"user_service/global"
	"user_service/model"
)

// FindPartnerByCode 根据合作方编码查询合作方配置。
func FindPartnerByCode(ctx context.Context, partnerCode string) (*model.Partner, error) {
	var partner model.Partner
	err := global.DB.GetContext(
		ctx,
		&partner,
		`SELECT partner_id, partner_code, partner_name, app_key, app_secret,
			userinfo_url, status, remark, created_at, updated_at
		FROM partners
		WHERE partner_code = ?`,
		partnerCode,
	)
	if err != nil {
		return nil, err
	}

	return &partner, nil
}
