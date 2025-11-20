package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("无效的Token")
	ErrTokenExpired = errors.New("Token已过期")
)

// JWTClaims JWT 声明
type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Manager JWT 管理器
type Manager struct {
	secretKey     string
	tokenDuration time.Duration
}

// NewManager 创建 JWT 管理器
func NewManager(secretKey string, tokenDuration time.Duration) *Manager {
	return &Manager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

// ValidateToken 验证 JWT Token
func (m *Manager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
