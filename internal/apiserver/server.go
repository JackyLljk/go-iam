package apiserver

import (
	"context"
	"fmt"
	"j-iam/internal/apiserver/config"
	cachev1 "j-iam/internal/apiserver/controller/v1/cache"
	"j-iam/internal/apiserver/store"
	"j-iam/internal/apiserver/store/mysql"
	genericoptions "j-iam/internal/pkg/options"
	pb "j-iam/internal/pkg/proto/apiserver/v1"
	genericapiserver "j-iam/internal/pkg/server"
	"j-iam/pkg/log"
	"j-iam/pkg/shutdown"
	"j-iam/pkg/shutdown/shutdownmanagers/posixsignal"
	"j-iam/pkg/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type apiServer struct {
	genericAPIServer *genericapiserver.GenericAPIServer // 通用 API 服务配置（HTTP、HTTPS）
	gs               *shutdown.GracefulShutdown         // 优雅关停服务
	redisOptions     *genericoptions.RedisOptions
	gRPCAPIServer    *grpcAPIServer
}

// preparedAPIServer
type preparedAPIServer struct {
	*apiServer
}

// ExtraConfig 定义 apiserver 的额外配置
type ExtraConfig struct {
	Addr         string
	MaxMsgSize   int
	mysqlOptions *genericoptions.MySQLOptions
	ServerCert   genericoptions.GeneratableKeyCert
	// etcdOptions      *genericoptions.EtcdOptions
}

type completedExtraConfig struct {
	*ExtraConfig
}

// 填充字段
func (c *ExtraConfig) complete() *completedExtraConfig {
	if c.Addr == "" {
		c.Addr = "127.0.0.1:8081"
	}

	return &completedExtraConfig{c}
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {
	// 优雅关停服务
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	// 通用服务配置
	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	// 额外服务配置
	extraConfig, err := buildExtraConfig(cfg)
	if err != nil {
		return nil, err
	}

	// 根据配置创建实例
	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}
	extraServer, err := extraConfig.complete().New()
	if err != nil {
		return nil, err
	}

	// 注册服务
	server := &apiServer{
		gs:               gs,
		genericAPIServer: genericServer,
		redisOptions:     cfg.RedisOptions,
		gRPCAPIServer:    extraServer,
	}

	return server, nil
}

// PrepareRun 使用配置初始化各种应用
func (s *apiServer) PrepareRun() preparedAPIServer {
	// 初始化 gin 路由（注册中间件、注册 RESTful API）
	initRouter(s.genericAPIServer.Engine)

	// 初始化 redis
	s.initRedisStore()

	// 优雅关闭回调函数
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		mysqlStore, _ := mysql.GetMySQLFactory(nil)
		if mysqlStore != nil {
			_ = mysqlStore.Close()
		}

		s.gRPCAPIServer.Close()
		s.genericAPIServer.Close()

		return nil
	}))

	return preparedAPIServer{s}
}

func (s preparedAPIServer) Run() error {
	// 开一个 goroutine，启动 gRPC 服务
	go s.gRPCAPIServer.Run()

	// 启动优雅关闭 managers
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	return s.genericAPIServer.Run()
}

// New 创建 gRPC 实例
func (c *completedExtraConfig) New() (*grpcAPIServer, error) { // 加上 grpc 服务
	//creds, err := credentials.NewServerTLSFromFile(c.ServerCert.CertKey.CertFile, c.ServerCert.CertKey.KeyFile)
	creds, err := credentials.NewServerTLSFromFile("cert/server.pem", "cert/server.key") // 临时使用一下
	if err != nil {
		log.Fatalf("Failed to generate credentials %s", err.Error())
	}
	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(c.MaxMsgSize), grpc.Creds(creds)}
	grpcServer := grpc.NewServer(opts...)

	storeIns, _ := mysql.GetMySQLFactory(c.mysqlOptions)

	//storeIns, _ := etcd.GetEtcdFactoryOr(c.etcdOptions, nil) 	// 可选 etcd ？

	store.SetClient(storeIns)

	cacheIns, err := cachev1.GetCacheInsOr(storeIns)
	if err != nil {
		log.Fatalf("Failed to get cache instance: %s", err.Error())
	}

	pb.RegisterCacheServer(grpcServer, cacheIns)

	reflection.Register(grpcServer)

	return &grpcAPIServer{grpcServer, c.Addr}, nil
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
	//if err = cfg.SecureServing.ApplyTo(genericConfig); err != nil {
	//	return
	//}

	if err = cfg.InsecureServing.ApplyTo(genericConfig); err != nil {
		return
	}

	return
}

func buildExtraConfig(cfg *config.Config) (*ExtraConfig, error) {
	// 1. grpc 服务配置
	// 2. 证书相关配置
	// 3. MySQL 相关配置
	return &ExtraConfig{
		mysqlOptions: cfg.MySQLOptions,
		Addr:         fmt.Sprintf("%s:%d", cfg.GRPCOptions.BindAddress, cfg.GRPCOptions.BindPort),
		MaxMsgSize:   cfg.GRPCOptions.MaxMsgSize,
		//ServerCert:   cfg.SecureServing.ServerCert,
		// etcdOptions:      cfg.EtcdOptions,
	}, nil
}

// initRedisStore 初始化 Redis
func (s *apiServer) initRedisStore() {
	ctx, cancel := context.WithCancel(context.Background())
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		cancel()

		return nil
	}))

	redisConfig := &storage.Config{
		Host:                  s.redisOptions.Host,
		Port:                  s.redisOptions.Port,
		Addrs:                 s.redisOptions.Addrs,
		MasterName:            s.redisOptions.MasterName,
		Username:              s.redisOptions.Username,
		Password:              s.redisOptions.Password,
		Database:              s.redisOptions.Database,
		MaxIdle:               s.redisOptions.MaxIdle,
		MaxActive:             s.redisOptions.MaxActive,
		Timeout:               s.redisOptions.Timeout,
		EnableCluster:         s.redisOptions.EnableCluster,
		UseSSL:                s.redisOptions.UseSSL,
		SSLInsecureSkipVerify: s.redisOptions.SSLInsecureSkipVerify,
	}

	// try to connect to redis
	go storage.ConnectToRedis(ctx, redisConfig)
}
