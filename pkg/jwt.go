package pkg

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"user_service/config"
)

type Claims struct {
	UserID               int64 `json:"user_id"` // 用户 ID
	jwt.RegisteredClaims       // 内嵌标准 JWT 声明字段
}

// GenerateToken 生成 JWT Token
func GenerateToken(UserID int64) (string, error) {
	// 从配置文件读取密钥和过期时间
	secret := config.Cfg.JWT.Secret
	expire := config.Cfg.JWT.Expire

	claims := Claims{
		UserID: UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析 JWT Token
func ParseToken(tokenStr string) (*Claims, error) {
	secret := config.Cfg.JWT.Secret

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, errors.New("token 无效或已过期")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("token 解析失败")
	}

	return claims, nil
}
