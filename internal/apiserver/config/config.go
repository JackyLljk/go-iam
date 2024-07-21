package config

import "j-iam/internal/apiserver/options"

// 定义控制流 j-apiserver 的配置

// Config 是运行 pump 服务的配置（
type Config struct {
	*options.Options
}

// CreateConfigFromOptions 根据给定的 pump 命令行或配置文件选项创建正在运行的配置实例
func CreateConfigFromOptions(opts *options.Options) *Config {
	return &Config{opts}
}
