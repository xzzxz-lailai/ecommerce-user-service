package repo

import (
	"context"
	"user_service/global"
	"user_service/model"
)

// ExistsByUsername 按用户名查询是否存在
func ExistsByUsername(ctx context.Context, username string) (bool, error) {

	var count int

	err := global.DB.GetContext(
		ctx,
		&count,
		"SELECT COUNT(*) FROM users WHERE username = ?", username,
	)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// FindByUsername 根据用户名查询用户
func FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := global.DB.GetContext(
		ctx,
		&user,
		"SELECT user_id, username, password FROM users WHERE username = ?",
		username,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ExistsByEmail 按邮箱查询是否存在
func ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int
	err := global.DB.GetContext(ctx, &count, "SELECT COUNT(*) FROM users WHERE email = ?", email)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 根据 userID 查询用户
func FindByID(ctx context.Context, usserID int64) (*model.User, error) {
	var user model.User
	err := global.DB.GetContext(ctx,
		&user,
		"SELECT user_id,username, password FROM users WHERE user_id = ?",
		usserID,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func Create(ctx context.Context, user *model.User) error {
	_, err := global.DB.ExecContext(
		ctx,
		"INSERT INTO users (username, password, email) VALUES (?, ?, ?)",
		user.Username,
		user.Password,
		user.Email,
	)

	return err
}
func UpdatePassword(ctx context.Context, userID int64, password string) error {
	_, err := global.DB.ExecContext(
		ctx,
		"UPDATE users SET password = ? WHERE user_id = ?",
		password,
		userID,
	)

	return err
}
