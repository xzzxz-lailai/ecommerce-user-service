package auth

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"user_service/httpclient"
	"user_service/model"
	"user_service/pkg"
	"user_service/repo"
)

// SendRegisterEmailCode 发送注册邮箱验证码
func SendRegisterEmailCode(ctx context.Context, req *model.SendRegisterEmailCodeRequest) error {
	// 先查询邮箱是否已经注册，避免已注册邮箱继续发送注册验证码
	emailExist, err := repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if emailExist {
		return errors.New("邮箱已注册")
	}

	// 调用 email-service 发送注册验证码
	return httpclient.SendRegisterEmailCode(ctx, req.Email)
}

// SendForgetPasswordEmailCode 发送忘记密码邮箱验证码
func SendForgetPasswordEmailCode(ctx context.Context, req *model.SendForgetPasswordEmailCodeRequest) error {
	// 先查询邮箱是否存在，只有已注册邮箱才能找回密码
	emailExist, err := repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if !emailExist {
		return errors.New("邮箱未注册")
	}

	// 调用 email-service 发送忘记密码验证码
	return httpclient.SendForgetPasswordEmailCode(ctx, req.Email)
}

// ForgetPassword 验证邮箱验证码后重置密码
func ForgetPassword(ctx context.Context, req *model.ForgetPasswordRequest) error {
	// 先查询邮箱是否存在
	emailExist, err := repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if !emailExist {
		return errors.New("邮箱未注册")
	}

	// 调用 email-service 验证忘记密码验证码
	if err := httpclient.VerifyForgetPasswordEmailCode(ctx, req.Email, req.Code); err != nil {
		return err
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 根据邮箱修改密码
	return repo.UpdatePasswordByEmail(ctx, req.Email, string(hashedPassword))
}

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
	// 查询邮箱是否已经注册
	emailExist, err := repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if emailExist {
		return errors.New("邮箱已注册")
	}
	// 调用email-service验证邮箱函数
	if err := httpclient.VerifyRegisterEmailCode(ctx, req.Email, req.Code); err != nil {
		return err
	}
	// 加密密码
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}
	// 构建用户对象
	user := &model.User{
		Username: req.Username,
		Password: string(hash),
		Email:    req.Email,
	}

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
