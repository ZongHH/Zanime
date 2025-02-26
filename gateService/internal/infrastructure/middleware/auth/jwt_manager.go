package auth

import (
	"errors"
	"fmt"
	"gateService/internal/domain/entity"
	"gateService/internal/infrastructure/config"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken     = errors.New("无效的token")
	ErrTokenExpired     = errors.New("token已过期")
	ErrInvalidSignature = errors.New("无效的签名")
	ErrMalformedToken   = errors.New("格式错误的token")
)

// JWTManager JWT管理器
type JWTManager struct {
	config *config.JWTConfig
}

// CustomClaims 自定义Claims结构体
type CustomClaims struct {
	UserInfo  *entity.UserInfo `json:"user_info"`
	TokenType string           `json:"token_type"`
	jwt.RegisteredClaims
}

// NewJWTManager 创建JWT管理器实例
func NewJWTManager(config *config.JWTConfig) *JWTManager {
	if config == nil {
		config = DefaultJWTConfig()
	}
	return &JWTManager{config: config}
}

func DefaultJWTConfig() *config.JWTConfig {
	return &config.JWTConfig{
		SecretKey:    "your-secret-key",
		Issuer:       "your-issuer",
		AccessToken:  config.JWTTokenConfig{},
		RefreshToken: config.JWTTokenConfig{},
		TokenType:    "Bearer",
	}
}

// GenerateToken 生成JWT令牌
func (m *JWTManager) GenerateToken(userInfo *entity.UserInfo) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserInfo:  userInfo,           // 用户信息
		TokenType: m.config.TokenType, // 令牌类型
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.config.Issuer,                                              // 令牌签发者
			Subject:   fmt.Sprintf("user_%d", userInfo.UserID),                      // 令牌主题,包含用户ID
			IssuedAt:  jwt.NewNumericDate(now),                                      // 令牌签发时间
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.AccessToken.ExpireTime)), // 令牌过期时间
			NotBefore: jwt.NewNumericDate(now),                                      // 令牌生效时间
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.SecretKey))
}

// GenerateRefreshToken 生成刷新令牌
func (m *JWTManager) GenerateRefreshToken(userInfo *entity.UserInfo) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserInfo:  userInfo,  // 用户信息
		TokenType: "Refresh", // 令牌类型为刷新令牌
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.config.Issuer,                                               // 令牌签发者
			Subject:   fmt.Sprintf("refresh_user_%d", userInfo.UserID),               // 令牌主题,包含用户ID
			IssuedAt:  jwt.NewNumericDate(now),                                       // 令牌签发时间
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.RefreshToken.ExpireTime)), // 令牌过期时间
			NotBefore: jwt.NewNumericDate(now),                                       // 令牌生效时间
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.SecretKey))
}

// ParseToken 解析并验证token
func (m *JWTManager) ParseToken(tokenString string) (*CustomClaims, error) {
	// 移除Bearer前缀
	tokenString = strings.TrimPrefix(tokenString, m.config.TokenType+" ")

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignature
		}
		return []byte(m.config.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrMalformedToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// RefreshToken 刷新token
func (m *JWTManager) RefreshToken(refreshToken string) (string, error) {
	claims, err := m.ParseToken(refreshToken)
	if err != nil {
		return "", err
	}

	if claims.TokenType != "Refresh" {
		return "", errors.New("不是有效的刷新令牌")
	}

	return m.GenerateToken(claims.UserInfo)
}

// ValidateToken 验证token的有效性
func (m *JWTManager) ValidateToken(tokenString string) bool {
	_, err := m.ParseToken(tokenString)
	return err == nil
}

// GetUserFromToken 从token中获取用户信息
func (m *JWTManager) GetUserFromToken(tokenString string) (*entity.UserInfo, error) {
	claims, err := m.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	return claims.UserInfo, nil
}
