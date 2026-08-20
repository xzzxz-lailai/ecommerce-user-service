package partner

import (
	"context"
	"user_service/global"
	"user_service/model"
)

// FindPartnerUser 根据合作方和外部用户ID查询用户映射。
func FindPartnerUser(ctx context.Context, partnerID int64, externalUserID string) (*model.PartnerUser, error) {
	var partnerUser model.PartnerUser
	err := global.DB.GetContext(
		ctx,
		&partnerUser,
		`SELECT partner_user_id, partner_id, external_user_id, nickname, phone,
			email, avatar, status, last_login_at, created_at, updated_at
		FROM partner_users
		WHERE partner_id = ? AND external_user_id = ?`,
		partnerID,
		externalUserID,
	)
	if err != nil {
		return nil, err
	}

	return &partnerUser, nil
}

// CreatePartnerUser 创建合作方用户映射。
func CreatePartnerUser(ctx context.Context, partnerUser *model.PartnerUser) error {
	result, err := global.DB.ExecContext(
		ctx,
		`INSERT INTO partner_users (partner_id, external_user_id, nickname, phone, email, avatar, status, last_login_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())`,
		partnerUser.PartnerID,
		partnerUser.ExternalUserID,
		partnerUser.Nickname,
		partnerUser.Phone,
		partnerUser.Email,
		partnerUser.Avatar,
		partnerUser.Status,
	)
	if err != nil {
		return err
	}

	partnerUserID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	partnerUser.PartnerUserID = partnerUserID

	return nil
}

// UpdatePartnerUserLoginTime 更新合作方用户最近登录时间。
func UpdatePartnerUserLoginTime(ctx context.Context, partnerUserID int64) error {
	_, err := global.DB.ExecContext(
		ctx,
		"UPDATE partner_users SET last_login_at = NOW() WHERE partner_user_id = ?",
		partnerUserID,
	)

	return err
}
