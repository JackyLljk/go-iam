package apiserver

import (
	"j-iam/internal/apiserver/config"
	"j-iam/internal/apiserver/options"
	"j-iam/internal/pkg/app"
)

const commandDesc = `j-apiserver 验证和配置 API 对象的数据，
包括用户、策略、密钥等。API 服务器为 REST 操作提供服务，以执行 API 对象管理。
在以下位置查找更多 j-apiserver 信息：`

// NewApp 基于选项模式建立应用配置
func NewApp(basename string) *app.App {
	opts := options.NewOptions() // 创建带有默认值的 Options

	// 根据回调函数配置可选参数
	application := app.NewApp("j-iam API Server",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		//app.WithNoVersion(),	// 不打印应用版本信息
		//app.WithNoConfig(),	// 不打印应用配置信息
		//app.WithSilence(),	// 静默模式，不打印版本和配置信息
		app.WithRunFunc(run(opts)))

	return application
}

// run 注册在 App.runFunc 的应用启动回调函数
func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		//log.Init(opts.Log)
		//defer log.Flush()

		// 根据应用配置，构建 apiserver 配置
		cfg := config.CreateConfigFromOptions(opts)
		return Run(cfg)
	}
}

func Run(cfg *config.Config) error {
	server, err := createAPIServer(cfg)
	if err != nil {
		return err
	}

	return server.PrepareRun().Run()
}
