package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/log"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

// GenericAPIServer 保存 api api service 的服务状态(理解为 gin.Engine 的装饰器)
type GenericAPIServer struct {
	middlewares []string

	//SecureServingInfo   *SecureServingInfo   // HTTPS（TLS）服务配置
	InsecureServingInfo *InsecureServingInfo // HTTP 服务配置
	*gin.Engine

	health                       bool
	enableMetrics                bool
	enableProfiling              bool
	insecureServer, secureServer *http.Server
	ShutdownTimeout              time.Duration
}

// initGenericAPIServer 配置 gin 服务
func initGenericAPIServer(s *GenericAPIServer) {
	// do some setup
	// s.GET(path, ginSwagger.WrapHandler(swaggerFiles.Handler))

	//s.Setup()
	//s.InstallMiddlewares()	// 调用中间件
	//s.InstallAPIs()	// 加载 API
}

// Run 会生成 http 服务器。仅当端口最初无法侦听时，它才会返回。
func (s *GenericAPIServer) Run() error {
	// 为了可扩展性，使用自定义的 HTTP 服务
	s.insecureServer = &http.Server{
		Addr:           s.InsecureServingInfo.Address,
		Handler:        s,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// For scalability, use custom HTTP configuration mode here
	//s.secureServer = &http.Server{
	//	Addr:    s.SecureServingInfo.Address(),
	//	Handler: s,
	//	// ReadTimeout:    10 * time.Second,
	//	// WriteTimeout:   10 * time.Second,
	//	// MaxHeaderBytes: 1 << 20,
	//}

	// 使用 “errgroup.Group” 包，实现web核心功能" 一进程、多服务（内外有别，安全性）"
	var eg errgroup.Group

	// 在 goroutine 中初始化 HTTP 服务，以便它不会阻塞 main goroutine 的正常关机处理
	// 一进程多服务：监听 HTTP、HTTPS 端口
	eg.Go(func() error {
		log.Infof("Start to listening the incoming requests on http address: %s", s.InsecureServingInfo.Address)
		if err := s.insecureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())

			return err
		}

		log.Infof("Server on %s stopped", s.InsecureServingInfo.Address)

		return nil
	})

	// 启动 HTTPS 服务
	//eg.Go(func() error {
	//	key, cert := s.SecureServingInfo.CertKey.KeyFile, s.SecureServingInfo.CertKey.CertFile
	//	if cert == "" || key == "" || s.SecureServingInfo.BindPort == 0 {
	//		return nil
	//	}
	//
	//	log.Infof("Start to listening the incoming requests on https address: %s", s.SecureServingInfo.Address())
	//
	//	if err := s.secureServer.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
	//		log.Fatal(err.Error())
	//
	//		return err
	//	}
	//
	//	log.Infof("Server on %s stopped", s.SecureServingInfo.Address())
	//
	//	return nil
	//})

	// Ping the service to make sure the router is working.
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//if s.health {
	//	if err := s.ping(ctx); err != nil {
	//		return err
	//	}
	//}

	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
	}

	return nil
}
