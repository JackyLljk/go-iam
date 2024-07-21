package server

import (
	"github.com/gin-gonic/gin"
)

// generic api service 应用配置

// Config GenericAPIServer（通用应用配置）
type Config struct {
	//SecureServing   *SecureServingInfo   // TLS 服务配置
	InsecureServing *InsecureServingInfo // HTTP 服务配置
	//Jwt             *JwtInfo             // JWT 中间件配置
	Mode        string // gin 的运行模式(DebugMode / ReleaseMode / TestMode)
	Middlewares []string
	Health      bool
	//EnableProfiling bool
	//EnableMetrics   bool
}

func NewConfig() *Config {
	return &Config{
		Health:      true,
		Mode:        gin.ReleaseMode,
		Middlewares: []string{},
		//EnableProfiling: true,
		//EnableMetrics:   true,
		//Jwt: &JwtInfo{
		//	Realm:      "iam jwt",
		//	Timeout:    1 * time.Hour,
		//	MaxRefresh: 1 * time.Hour,
		//},
	}
}

// InsecureServingInfo 完整的 HTTP 服务地址 URL:Port
type InsecureServingInfo struct {
	Address string
}

// CompletedConfig 完成的 GenericAPIServer 配置
type CompletedConfig struct {
	*Config
}

func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{c}
}

func (c CompletedConfig) New() (*GenericAPIServer, error) {
	// setMode before gin.New()
	gin.SetMode(c.Mode)

	s := &GenericAPIServer{
		//SecureServingInfo:   c.SecureServing,
		InsecureServingInfo: c.InsecureServing,
		health:              c.Health,
		middlewares:         c.Middlewares,
		//enableMetrics:       c.EnableMetrics,
		//enableProfiling:     c.EnableProfiling,
		Engine: gin.New(),
	}

	initGenericAPIServer(s)

	return s, nil
}
