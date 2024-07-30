package authzserver

import (
	"context"

	"github.com/marmotedu/errors"

	"j-iam/internal/authzserver/analytics"
	"j-iam/internal/authzserver/config"
	"j-iam/internal/authzserver/load"
	"j-iam/internal/authzserver/load/cache"
	"j-iam/internal/authzserver/store/apiserver"
	genericoptions "j-iam/internal/pkg/options"
	genericapiserver "j-iam/internal/pkg/server"
	"j-iam/pkg/log"
	"j-iam/pkg/shutdown"
	"j-iam/pkg/shutdown/shutdownmanagers/posixsignal"
	"j-iam/pkg/storage"
)

// RedisKeyPrefix defines the prefix key in redis for analytics data.
const RedisKeyPrefix = "analytics-"

type authzServer struct {
	gs               *shutdown.GracefulShutdown
	rpcServer        string
	clientCA         string
	redisOptions     *genericoptions.RedisOptions
	genericAPIServer *genericapiserver.GenericAPIServer
	analyticsOptions *analytics.AnalyticsOptions
	redisCancelFunc  context.CancelFunc
}

type preparedAuthzServer struct {
	*authzServer
}

// func createAuthzServer(cfg *config.Config) (*authzServer, error) {.
func createAuthzServer(cfg *config.Config) (*authzServer, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	server := &authzServer{
		gs:               gs,
		redisOptions:     cfg.RedisOptions,
		analyticsOptions: cfg.AnalyticsOptions,
		rpcServer:        cfg.RPCServer,
		clientCA:         cfg.ClientCA,
		genericAPIServer: genericServer,
	}

	return server, nil
}

func (s *authzServer) PrepareRun() preparedAuthzServer {
	_ = s.initialize()

	initRouter(s.genericAPIServer.Engine)

	return preparedAuthzServer{s}
}

// Run start to run AuthzServer.
func (s preparedAuthzServer) Run() error {
	// in order to ensure that the reported data is not lost,
	// please ensure the following graceful shutdown sequence
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		s.genericAPIServer.Close()
		if s.analyticsOptions.Enable {
			analytics.GetAnalytics().Stop()
		}
		s.redisCancelFunc()

		return nil
	}))

	// start shutdown managers
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	return s.genericAPIServer.Run()
}

func buildGenericConfig(cfg *config.Config) (genericConfig *genericapiserver.Config, err error) {
	genericConfig = genericapiserver.NewConfig()
	if err = cfg.GenericServerRunOptions.ApplyTo(genericConfig); err != nil {
		return
	}

	if err = cfg.InsecureServing.ApplyTo(genericConfig); err != nil {
		return
	}

	if err = cfg.SecureServing.ApplyTo(genericConfig); err != nil {
		return
	}

	//if err = cfg.FeatureOptions.ApplyTo(genericConfig); err != nil {
	//	return
	//}
	//

	return
}

func (s *authzServer) buildStorageConfig() *storage.Config {
	return &storage.Config{
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
}

func (s *authzServer) initialize() error {
	ctx, cancel := context.WithCancel(context.Background())
	s.redisCancelFunc = cancel

	// keep redis connected
	go storage.ConnectToRedis(ctx, s.buildStorageConfig())

	cacheIns, err := cache.GetCacheInsOr(apiserver.GetAPIServerFactoryOrDie(s.rpcServer, s.clientCA))
	if err != nil {
		return errors.Wrap(err, "get cache instance failed")
	}

	load.NewLoader(ctx, cacheIns).Start()

	// start analytics service
	if s.analyticsOptions.Enable {
		analyticsStore := storage.RedisCluster{KeyPrefix: RedisKeyPrefix}
		analyticsIns := analytics.NewAnalytics(s.analyticsOptions, &analyticsStore)
		analyticsIns.Start()
	}

	return nil
}
