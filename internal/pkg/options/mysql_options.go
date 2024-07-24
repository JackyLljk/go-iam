package options

import (
	"time"

	"github.com/spf13/pflag"
)

type MySQLOptions struct {
	Host                  string        `json:"host,omitempty"                     mapstructure:"host"`
	Username              string        `json:"username,omitempty"                 mapstructure:"username"`
	Password              string        `json:"-"                                  mapstructure:"password"`
	Database              string        `json:"database"                           mapstructure:"database"`
	MaxIdleConnections    int           `json:"max-idle-connections,omitempty"     mapstructure:"max-idle-connections"`
	MaxOpenConnections    int           `json:"max-open-connections,omitempty"     mapstructure:"max-open-connections"`
	MaxConnectionLifeTime time.Duration `json:"max-connection-life-time,omitempty" mapstructure:"max-connection-life-time"`
	LogLevel              int           `json:"log-level"                          mapstructure:"log-level"`
}

// NewMySQLOptions 创建一个 'zero' value 实例
func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{
		// 临时修改一下
		Host:     "127.0.0.1:13306",
		Username: "root",
		Password: "iam123",
		Database: "iam",
		//Host:                  "127.0.0.1:3306",
		//Username:              "",
		//Password:              "",
		//Database:              "",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
		LogLevel:              1, // Silent
	}
}

// Validate 检验传递给 MySQLOptions 的标志 (flags) (实际上是用于实现 CliOptions 接口)
func (m *MySQLOptions) Validate() []error {
	return []error{}
}

// AddFlags 将与特定 APIServer 的 mysql 存储相关的标志添加到指定的 FlagSet 中
func (m *MySQLOptions) AddFlags(fs *pflag.FlagSet) {

}
