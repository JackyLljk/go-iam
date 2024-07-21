package options

import (
	"j-iam/internal/pkg/server"
	"net"
	"strconv"
)

// InsecureServingOptions HTTP 服务配置，服务地址和端口
type InsecureServingOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port"    mapstructure:"bind-port"`
}

func NewInsecureServingOptions() *InsecureServingOptions {
	return &InsecureServingOptions{
		BindAddress: "127.0.0.1",
		BindPort:    8080,
	}
}

func (s *InsecureServingOptions) ApplyTo(c *server.Config) error {
	c.InsecureServing = &server.InsecureServingInfo{
		Address: net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort)),
	}

	return nil
}
