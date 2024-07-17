package apiserver

import (
	"honnef.co/go/tools/config"
	"j-iam/internal/apiserver/options"
	"j-iam/internal/pkg/app"
)

// NewApp 基于选项模式建立应用配置
func NewApp(basename string) *app.App {
	opts := options.NewOptions() // 创建带有默认值的 Options

	// 根据回调函数配置可选参数
	application := app.NewApp("j-iam API Server",
		basename,
		app.WithOptions(opts))

	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		//log.Init(opts.Log)
		//defer log.Flush()
		//
		//cfg, err := config.CreateConfigFromOptions(opts)
		//if err != nil {
		//	return err
		//}
		//
		return nil
	}
}

func Run(cfg *config.Config) error {
	//server, err := createAPIServer(cfg)
	//if err != nil {
	//	return err
	//}
	//
	//return server.PrepareRun().Run()
	return nil
}
