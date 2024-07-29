package server

import (
	"net"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const jwtKey = "dfVpOK8LZeJLZHYmHdb1VdyRrACKpqoo" // jwt 密钥，临时保存（后面还是从配置中读取）

// generic api service 应用配置

// Config GenericAPIServer（通用应用配置）
type Config struct {
	SecureServing   *SecureServingInfo   // TLS 服务配置
	InsecureServing *InsecureServingInfo // HTTP 服务配置
	Jwt             *JwtInfo             // JWT 中间件配置
	Mode            string               // gin 的运行模式(DebugMode / ReleaseMode / TestMode)
	Middlewares     []string
	Health          bool
	//EnableProfiling bool
	//EnableMetrics   bool
}

// CertKey contains configuration items related to certificate.
type CertKey struct {
	// CertFile is a file containing a PEM-encoded certificate, and possibly the complete certificate chain
	CertFile string
	// KeyFile is a file containing a PEM-encoded private key for the certificate specified by CertFile
	KeyFile string
}

// SecureServingInfo holds configuration of the TLS server.
type SecureServingInfo struct {
	BindAddress string
	BindPort    int
	CertKey     CertKey
}

// Address join host IP address and host port number into a address string, like: 0.0.0.0:8443.
func (s *SecureServingInfo) Address() string {
	return net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort))
}

func NewConfig() *Config {
	return &Config{
		Health:      true,
		Mode:        gin.ReleaseMode,
		Middlewares: []string{},
		//EnableProfiling: true,
		//EnableMetrics:   true,
		Jwt: &JwtInfo{
			Realm: "gin-jwt",
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
		SecureServingInfo:   c.SecureServing,
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
//func LoadConfig(cfg string, defaultName string) {
//	if cfg != "" {
//		viper.SetConfigFile(cfg)
//	} else {
//		viper.AddConfigPath(".")
//		viper.AddConfigPath(filepath.Join(homedir.HomeDir(), RecommendedHomeDir))
//		viper.AddConfigPath("/etc/iam")
//		viper.SetConfigName(defaultName)
//	}
//
//	// Use config file from the flag.
//	viper.SetConfigType("yaml")              // set the type of the configuration to yaml.
//	viper.AutomaticEnv()                     // read in environment variables that match.
//	viper.SetEnvPrefix(RecommendedEnvPrefix) // set ENVIRONMENT variables prefix to IAM.
//	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
//
//	// If a config file is found, read it in.
//	if err := viper.ReadInConfig(); err != nil {
//		log.Warnf("WARNING: viper failed to discover and load the configuration file: %s", err.Error())
//	}
//}
