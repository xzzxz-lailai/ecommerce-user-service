package model

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Code     string `json:"code" binding:"required,len=6"`
}

type SendRegisterEmailCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type SendForgetPasswordEmailCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Code        string `json:"code" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"gi`
}
type LoginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
