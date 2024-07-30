package authzserver

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"

	"j-iam/internal/authzserver/controller/v1/authorize"
	"j-iam/internal/authzserver/load/cache"
	"j-iam/internal/pkg/code"
	"j-iam/pkg/log"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
}

func installController(g *gin.Engine) *gin.Engine {
	// 基于 secretID、secretKey 生成 cache jwt token 进行认证
	// TODO: 为什么需要这个认证？是根据 secrets 进行的 jwt 认证？
	// TODO 流程：登录 apiserver -> 查询 secrets 和自己的 policies -> 根据授权策略，选择资源访问 -> 到 authzserver 先认证，后鉴权
	//auth := newCacheAuth()
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "page not found."), nil)
	})

	cacheIns, _ := cache.GetCacheInsOr(nil)
	if cacheIns == nil {
		log.Panicf("get nil cache instance")
	}

	apiv1 := g.Group("/v1")
	{
		authzController := authorize.NewAuthzController(cacheIns)

		// Router for authorization

		// apiserver 端管理员创建用户权限策略
		// authzserver 端根据拉取的策略和密钥对用户进行鉴权
		apiv1.POST("/authz", authzController.Authorize)

		// TODO：怎么操作成功鉴权的资源？
	}

	return g
}
