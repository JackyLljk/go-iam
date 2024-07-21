package apiserver

import (
	"j-iam/internal/apiserver/config"
	genericapiserver "j-iam/internal/pkg/server"
)

type apiServer struct {
	//gs               *shutdown.GracefulShutdown
	//redisOptions     *genericoptions.RedisOptions
	//gRPCAPIServer    *grpcAPIServer
	genericAPIServer *genericapiserver.GenericAPIServer // 通用 API 服务配置（HTTP、HTTPS）
}

// preparedAPIServer
type preparedAPIServer struct {
	*apiServer
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {
	// 优雅关停服务
	//gs := shutdown.New()
	//gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	// 通用服务配置
	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	// 其他服务配置
	//extraConfig, err := buildExtraConfig(cfg)
	//if err != nil {
	//	return nil, err
	//}

	// 根据配置创建实例
	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}
	//extraServer, err := extraConfig.complete().New()
	//if err != nil {
	//	return nil, err
	//}

	// 注册服务
	server := &apiServer{
		//gs:               gs,
		//redisOptions:     cfg.RedisOptions,
		genericAPIServer: genericServer,
		//gRPCAPIServer:    extraServer,
	}

	return server, nil
}

func buildGenericConfig(cfg *config.Config) (genericConfig *genericapiserver.Config, err error) {
	// 初始化通用配置
	genericConfig = genericapiserver.NewConfig()

	// 配置 ServerRunOptions（模式、健康检查、中间件）
	if err = cfg.GenericServerRunOptions.ApplyTo(genericConfig); err != nil {
		return
	}

	// API service feature 配置
	//if err = cfg.FeatureOptions.ApplyTo(genericConfig); lastErr != nil {
	//	return
	//}
	// HTTPS 服务配置
	//if err = cfg.SecureServing.ApplyTo(genericConfig); lastErr != nil {
	//	return
	//}

	if err = cfg.InsecureServing.ApplyTo(genericConfig); err != nil {
		return
	}

	return
}

// PrepareRun 使用配置初始化各种应用
func (s *apiServer) PrepareRun() preparedAPIServer {
	// 初始化 gin 路由（注册中间件、注册 RESTful API）
	initRouter(s.genericAPIServer.Engine)

	// 初始化 redis
	//s.initRedisStore()

	// 优雅关闭回调函数
	//s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
	//	mysqlStore, _ := mysql.GetMySQLFactoryOr(nil)
	//	if mysqlStore != nil {
	//		_ = mysqlStore.Close()
	//	}
	//
	//	s.gRPCAPIServer.Close()
	//	s.genericAPIServer.Close()
	//
	//	return nil
	//}))

	return preparedAPIServer{s}
}

func (s preparedAPIServer) Run() error {
	// 开一个 goroutine，启动 gRPC 服务
	//go s.gRPCAPIServer.Run()

	// start shutdown managers
	// 启动优雅关闭
	//if err := s.gs.Start(); err != nil {
	//	log.Fatalf("start shutdown manager failed: %s", err.Error())
	//}

	return s.genericAPIServer.Run()
}
