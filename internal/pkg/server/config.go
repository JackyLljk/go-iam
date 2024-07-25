package server

import (
	"path/filepath"
	"strings"
	"time"

	"j-iam/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/util/homedir"
	"github.com/spf13/viper"
)

const (
	// RecommendedHomeDir defines the default directory used to place all iam service configurations.
	RecommendedHomeDir = ".iam"
	// RecommendedEnvPrefix defines the ENV prefix used by all iam service.
	RecommendedEnvPrefix = "IAM"

	jwtKey = "dfVpOK8LZeJLZHYmHdb1VdyRrACKpqoo" // jwt 密钥，临时保存（后面还是从配置中读取）
)

// generic api service 应用配置

// Config GenericAPIServer（通用应用配置）
type Config struct {
	//SecureServing   *SecureServingInfo   // TLS 服务配置
	InsecureServing *InsecureServingInfo // HTTP 服务配置
	Jwt             *JwtInfo             // JWT 中间件配置
	Mode            string               // gin 的运行模式(DebugMode / ReleaseMode / TestMode)
	Middlewares     []string
	Health          bool
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
		Jwt: &JwtInfo{
			Realm: "iam jwt",
			Key:   jwtKey,
			//Timeout:    1 * time.Hour,
			//MaxRefresh: 1 * time.Hour,
			Timeout:    24 * time.Hour,
			MaxRefresh: 3 * 24 * time.Hour,
		},
	}
}

// InsecureServingInfo 完整的 HTTP 服务地址 URL:Port
type InsecureServingInfo struct {
	Address string
}

// JwtInfo 定义用于创建 JWT 认证中间件的 JWT 字段
type JwtInfo struct {
	// defaults to "iam jwt"
	Realm string
	// defaults to empty
	Key string
	// defaults to one hour
	Timeout time.Duration
	// defaults to zero
	MaxRefresh time.Duration
}

// CompletedConfig 完成的 GenericAPIServer 配置
type CompletedConfig struct {
	*Config
}

// Complete 填充字段
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

// LoadConfig 从配置文件中读取
func LoadConfig(cfg string, defaultName string) {
	if cfg != "" {
		viper.SetConfigFile(cfg)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(homedir.HomeDir(), RecommendedHomeDir))
		viper.AddConfigPath("/etc/iam")
		viper.SetConfigName(defaultName)
	}

	// Use config file from the flag.
	viper.SetConfigType("yaml")              // set the type of the configuration to yaml.
	viper.AutomaticEnv()                     // read in environment variables that match.
	viper.SetEnvPrefix(RecommendedEnvPrefix) // set ENVIRONMENT variables prefix to IAM.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("WARNING: viper failed to discover and load the configuration file: %s", err.Error())
	}
}
