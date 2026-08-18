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
const (
	CodeTypeRegister       = "register"
	CodeTypeForgetPassword = "forget_password"
)

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

// SendRegisterEmailCode 调用 email-service，发送注册场景下的邮箱验证码
func SendRegisterEmailCode(ctx context.Context, email string) error {
	// 1. 从 etcd 查询 email-service 的地址
	// 比如拿到：127.0.0.1:8081
	addr, err := GetEmailService()
	if err != nil {
		return err
	}
	// 2. 组装发送验证码需要的请求参数
	// 这个 JSON 会发送给 email-service 的验证码发送接口
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
	// 那 url = "http://127.0.0.1:8081/api/v1/email/code/send"
	url := fmt.Sprintf("http://%s/api/v1/email/code/send", addr)

	// 4. 通过 HTTP POST 调用 email-service 发送验证码
	// doEmailPost 内部会处理单次请求超时和有限次数重试
	if err := doEmailPost(url, body); err != nil {
		return errors.New("发送邮箱验证码失败")
	}

	// 5. 没有错误，说明验证码发送成功
	return nil
}

// VerifyRegisterEmailCode 调用 email-service，验证注册场景下的邮箱验证码
func VerifyRegisterEmailCode(ctx context.Context, email, code string) error {
	// 1. 从 etcd 查询 email-service 的地址
	// 比如拿到：127.0.0.1:8081
	addr, err := GetEmailService()
	if err != nil {
		return err
	}
	// 2. 组装验证验证码需要的请求参数
	// 这个 JSON 会发送给 email-service 的验证码验证接口
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

	// 4. 通过 HTTP POST 调用 email-service 验证验证码
	// doEmailPost 内部会处理单次请求超时和有限次数重试
	if err := doEmailPost(url, body); err != nil {
		return errors.New("邮箱验证码错误")
	}

	// 5. 没有错误，说明验证码验证成功
	return nil
}

// SendForgetPasswordEmailCode 调用 email-service，发送忘记密码场景下的邮箱验证码
func SendForgetPasswordEmailCode(ctx context.Context, email string) error {
	// 1. 从 etcd 查询 email-service 的地址
	// 比如拿到：127.0.0.1:8081
	addr, err := GetEmailService()
	if err != nil {
		return err
	}
	// 2. 组装发送验证码需要的请求参数
	// 这个 JSON 会发送给 email-service 的验证码发送接口
	// 最终内容类似：
	// {
	//   "email": "test@qq.com",
	//   "code_type": "forget_password"
	// }
	body, err := json.Marshal(SendEmailCodeRequest{
		Email:    email,
		CodeType: CodeTypeForgetPassword,
	})
	if err != nil {
		return err
	}
	// 3. 拼接 email-service 的发送验证码接口地址
	// 如果 addr = "127.0.0.1:8081"
	// 那 url = "http://127.0.0.1:8081/api/v1/email/code/send"
	url := fmt.Sprintf("http://%s/api/v1/email/code/send", addr)

	// 4. 通过 HTTP POST 调用 email-service 发送验证码
	// doEmailPost 内部会处理单次请求超时和有限次数重试
	if err := doEmailPost(url, body); err != nil {
		return errors.New("发送找回密码验证码失败")
	}

	// 5. 没有错误，说明验证码发送成功
	return nil
}

// VerifyForgetPasswordEmailCode 调用 email-service，验证忘记密码场景下的邮箱验证码
func VerifyForgetPasswordEmailCode(ctx context.Context, email, code string) error {
	// 1. 从 etcd 查询 email-service 的地址
	// 比如拿到：127.0.0.1:8081
	addr, err := GetEmailService()
	if err != nil {
		return err
	}
	// 2. 组装验证验证码需要的请求参数
	// 这个 JSON 会发送给 email-service 的验证码验证接口
	// 最终内容类似：
	// {
	//   "email": "test@qq.com",
	//   "code_type": "forget_password",
	//   "code": "123456"
	// }
	body, err := json.Marshal(VerifyEmailCodeRequest{
		Email:    email,
		CodeType: CodeTypeForgetPassword,
		Code:     code,
	})
	if err != nil {
		return err
	}
	// 3. 拼接 email-service 的验证码验证接口地址
	// 如果 addr = "127.0.0.1:8081"
	// 那 url = "http://127.0.0.1:8081/api/v1/email/code/verify"
	url := fmt.Sprintf("http://%s/api/v1/email/code/verify", addr)

	// 4. 通过 HTTP POST 调用 email-service 验证验证码
	// doEmailPost 内部会处理单次请求超时和有限次数重试
	if err := doEmailPost(url, body); err != nil {
		return errors.New("找回密码验证码错误")
	}

	// 5. 没有错误，说明验证码验证成功
	return nil
}

const (
	emailRequestTimeout = 2 * time.Second        // 单次请求最多等待 2 秒，超过就算超时失败
	emailRetryCount     = 3                      // 最多请求 3 次：第 1 次失败后，还可以重试 2 次
	emailRetryInterval  = 500 * time.Millisecond // 每次失败后，等待 500 毫秒再重试
)

// doEmailPost 发送 POST 请求到 email-service，失败时会有限重试
func doEmailPost(url string, body []byte) error {
	// 创建 HTTP 客户端，并设置单次请求超时时间
	client := &http.Client{
		Timeout: emailRequestTimeout,
	}
	var lastErr error // 保存最后一次失败的错误，最终 3 次都失败时返回它

	// 最多请求 emailRetryCount 次
	for i := 0; i < emailRetryCount; i++ {
		// 发起 POST 请求，请求体是 JSON
		resp, err := client.Post(url, "application/json", bytes.NewReader(body))
		// 如果请求本身失败，比如网络错误、超时
		if err != nil {
			lastErr = err
		} else {
			// 请求成功拿到响应后，要关闭响应体，释放连接资源
			resp.Body.Close()
			// 状态码是 200，说明调用成功，直接返回 nil
			if resp.StatusCode == http.StatusOK {
				return nil
			}
			// 请求发出去了，但 email-service 返回的状态码不是 200
			lastErr = fmt.Errorf("email-service 返回状态码: %d", resp.StatusCode)
		}
		// 如果还没到最后一次，就等一会儿再重试
		if i < emailRetryCount-1 {
			time.Sleep(emailRetryInterval)
		}
	}
	// 所有次数都失败，返回最后一次错误
	return lastErr
}
