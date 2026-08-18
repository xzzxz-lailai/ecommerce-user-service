package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"user_service/etcd"
)

// 从 etcd 查询 email-service 的服务地址
func GetEmailService() (string, error) {
	// 从 etcd 查询服务名为 "email-service" 的所有地址
	addrs, err := etcd.GetService("email-service")
	if err != nil {
		return "", err
	}
	if len(addrs) == 0 {
		return "", errors.New("email-service 服务不可用")
	}
	return addrs[0], nil
}

// 注册场景的验证码类型
// email-service 里会根据 code_type 区分：注册验证码、找回密码验证码等
const CodeTypeRegister = "register"

// 发送邮箱验证码时，请求 email-service 需要传的 JSON 参数
// 对应 email-service 的接口：POST /email/code/send
type SendEmailCodeRequest struct {
	Email    string `json:"email" binding:"required,email"`
	CodeType string `json:"code_type"` // 验证码类型
}
type VerifyEmailCodeRequest struct {
	Email    string `json:"email" binding:"required,email"`
	CodeType string `json:"code_type"` // 验证码类型
	Code     string `json:"code"`      // 验证码
}

// SendRegisterEmailCode 调用 email-service服务的,发送邮箱函数
func SendRegisterEmailCode(ctx context.Context, email string) error {
	// 1. 从 etcd 查询 email-service 的地址
	// 比如拿到：127.0.0.1:8081
	addr, err := GetEmailService()
	if err != nil {
		return err
	}
	// 2. 组装请求参数
	// 这个 JSON 会发送给 email-service
	// 最终内容类似：
	// {
	//   "email": "test@qq.com",
	//   "code_type": "register"
	// }
	body, err := json.Marshal(SendEmailCodeRequest{
		Email:    email,
		CodeType: CodeTypeRegister,
	})
	if err != nil {
		return err
	}
	// 3. 拼接 email-service 的发送验证码接口地址
	// 如果 addr = "127.0.0.1:8081"
	// 那 url = "http://127.0.0.1:8081/email/code/send"
	url := fmt.Sprintf("http://%s/api/v1/email/code/send", addr)

	// 4. 创建一个 HTTP POST 请求
	// ctx 用来控制请求生命周期
	// bytes.NewReader(body) 是把 JSON 数据放到请求体里
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	// 5. 告诉 email-service：我发送的是 JSON
	req.Header.Set("Content-Type", "application/json")

	// 6. 创建 HTTP 客户端
	// Timeout 表示最多等 5 秒，超过就算失败
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 7. 发送请求给 email-service
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 8. 判断 email-service 返回的 HTTP 状态码
	// 200 表示成功
	// 不是 200，就认为发送验证码失败
	if resp.StatusCode != http.StatusOK {
		return errors.New("发送邮箱验证码失败")
	}

	// 9. 没有错误，说明验证码发送成功
	return nil
}

// VerifyRegisterEmailCode 通过 HTTP 调用 email-service 验证注册邮箱验证码
func VerifyRegisterEmailCode(ctx context.Context, email, code string) error {
	// 1. 从 etcd 查询 email-service 的地址
	// 比如拿到：127.0.0.1:8081
	addr, err := GetEmailService()
	if err != nil {
		return err
	}
	// 2. 组装请求参数
	// 这个 JSON 会发送给 email-service
	// 最终内容类似：
	// {
	//   "email": "test@qq.com",
	//   "code_type": "register",
	//   "code": "123456"
	// }
	body, err := json.Marshal(VerifyEmailCodeRequest{
		Email:    email,
		CodeType: CodeTypeRegister,
		Code:     code,
	})
	if err != nil {
		return err
	}
	// 3. 拼接 email-service 的验证码验证接口地址
	// 如果 addr = "127.0.0.1:8081"
	// 那 url = "http://127.0.0.1:8081/api/v1/email/code/verify"
	url := fmt.Sprintf("http://%s/api/v1/email/code/verify", addr)
	// 4. 创建一个 HTTP POST 请求
	// ctx 用来控制请求生命周期
	// bytes.NewReader(body) 是把 JSON 数据放到请求体里
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	// 5. 告诉 email-service：我发送的是 JSON
	req.Header.Set("Content-Type", "application/json")
	// 6. 创建 HTTP 客户端
	// Timeout 表示最多等 5 秒，超过就算失败
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	// 7. 发送请求给 email-service
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 8. 判断 email-service 返回的 HTTP 状态码
	// 200 表示验证码正确
	// 不是 200，就认为验证码错误或已过期
	if resp.StatusCode != http.StatusOK {
		return errors.New("邮箱验证码错误")
	}

	// 9. 没有错误，说明验证码验证成功
	return nil
}
