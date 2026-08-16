package auth

import (
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"user_service/model"
	"user_service/pkg"
	"user_service/repo"
)

// Register 用户注册
func Register(ctx context.Context, req *model.RegisterRequest) error {

	// 查询用户名是否存在
	exist, err := repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return err
	}

	if exist {
		return errors.New("用户名已存在")
	}

	// 加密密码
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: string(hash),
	}

	// 保存到数据库
	return repo.Create(ctx, user)
}

func Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	// 根据用户名查询用户
	user, err := repo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("用户没有注册")
		}
		return nil, err
	}
	// 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	// 生成token
	token, err := pkg.GenerateToken(user.UserID)
	if err != nil {
		return nil, errors.New("登录失败，请稍后重试")
	}
	// 返回登录结结果
	return &model.LoginResponse{
		Username: user.Username,
		Token:    token,
	}, nil
}

// 修改密码
func ChangePassword(ctx context.Context, userID int64, req *model.ChangePasswordRequest) error {
	// 根据 userID 查询用户
	user, err := repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("用户不存在")
		}
		return err
	}
	// 校验旧密码
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.OldPassword),
	)
	if err != nil {
		return errors.New("旧密码错误")
	}
	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)
	// 修改密码
	err = repo.UpdatePassword(
		ctx,
		userID,
		string(hashedPassword),
	)
	if err != nil {
		return err
	}

	return nil
}
