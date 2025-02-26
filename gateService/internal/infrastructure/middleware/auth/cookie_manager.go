package auth

import (
	"gateService/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
)

type CookieManager struct {
	config *config.CookieConfig
}

func NewCookieManager(config *config.CookieConfig) *CookieManager {
	if config == nil {
		config = DefaultCookieConfig()
	}
	return &CookieManager{
		config: config,
	}
}

func DefaultCookieConfig() *config.CookieConfig {
	return &config.CookieConfig{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600,
		Secure:   false,
		HTTPOnly: true,
	}
}

func (m *CookieManager) SetTokenCookie(c *gin.Context, token string) {
	c.SetCookie(
		"token",
		token,
		m.config.MaxAge,
		m.config.Path,
		m.config.Domain,
		m.config.Secure,
		m.config.HTTPOnly,
	)
}

func (m *CookieManager) ClearTokenCookie(c *gin.Context) {
	c.SetCookie(
		"token",
		"",
		-1,
		m.config.Path,
		m.config.Domain,
		m.config.Secure,
		m.config.HTTPOnly,
	)
}

func (m *CookieManager) GetTokenCookie(c *gin.Context) (string, error) {
	token, err := c.Cookie("token")
	if err != nil {
		return "", err
	}
	return token, nil
}
