package authzserver

// 创建 authzserver 应用

import (
	"j-iam/internal/authzserver/config"
	"j-iam/internal/authzserver/options"
	"j-iam/internal/pkg/app"
	"j-iam/pkg/log"
)

const commandDesc = `Authorization(Authz) 服务使用 ladon 作为访问策略库，可以基于此进行资源保护

Find more ladon information at:
    https://github.com/ory/ladon`

// NewApp 基于默认参数构建 App
func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("IAM Authorization Server",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithNoVersion(),
		app.WithNoConfig(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.Log)
		defer log.Flush()

		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		return Run(cfg)
	}
}

// Run runs the specified AuthzServer. This should never exit.
func Run(cfg *config.Config) error {
	server, err := createAuthzServer(cfg)
	if err != nil {
		return err
	}

	return server.PrepareRun().Run()
}
