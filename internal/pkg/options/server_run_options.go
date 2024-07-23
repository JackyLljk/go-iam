package options

import (
	"j-iam/internal/pkg/server"

	"github.com/spf13/pflag"
)

// 定义运行通用应用（generic api service）的配置

type ServerRunOptions struct {
	// 运行模式(gin）
	Mode string `json:"mode"        mapstructure:"mode"`
	// 健康检查
	Health bool `json:"health"     mapstructure:"health"`
	// 中间件使用
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
}

// NewServerRunOptions 配置默认参数
func NewServerRunOptions() *ServerRunOptions {
	defaults := server.NewConfig()

	return &ServerRunOptions{
		Mode:        defaults.Mode,
		Health:      defaults.Health,
		Middlewares: defaults.Middlewares,
	}
}

// ApplyTo 将通用应用配置添加到配置中
func (s *ServerRunOptions) ApplyTo(c *server.Config) error {
	c.Mode = s.Mode
	c.Health = s.Health
	c.Middlewares = s.Middlewares

	return nil
}

func (s *ServerRunOptions) Validate() []error {
	return []error{}
}

func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs.StringVar(&s.Mode, "service.mode", s.Mode, ""+
		"Start the service in a specified service mode. Supported service mode: debug, test, release.")

	fs.BoolVar(&s.Health, "service.health", s.Health, ""+
		"Add self readiness check and install /health router.")

	fs.StringSliceVar(&s.Middlewares, "service.middlewares", s.Middlewares, ""+
		"List of allowed middlewares for service, comma separated. If this list is empty default middlewares will be used.")
}
