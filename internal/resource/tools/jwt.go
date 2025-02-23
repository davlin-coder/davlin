package tools

import (
	"fmt"
	"time"

	"github.com/davlin-coder/davlin/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims 定义JWT的payload结构
type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTManager JWT管理器
type JWTManager struct {
	config *config.JWTConfig
}

// NewJWTManager 创建JWT管理器实例
func NewJWTManager(cfg *config.Config) *JWTManager {
	return &JWTManager{config: &cfg.JWT}
}

// GenerateToken 生成JWT令牌
func (m *JWTManager) GenerateToken(userID uint, username string) (string, error) {
	claims := JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			Issuer:    "davlin-auth",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.config.Expire * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.config.SecretKey))
	if err != nil {
		return "", fmt.Errorf("生成token失败: %v", err)
	}

	return tokenString, nil
}

// ParseToken 解析JWT令牌
func (m *JWTManager) ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.config.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析token失败: %v", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("无效的token")
}
