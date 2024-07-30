package apiserver

import (
	"j-iam/internal/apiserver/controller/v1/policy"
	"j-iam/internal/apiserver/controller/v1/secret"
	"j-iam/internal/apiserver/controller/v1/user"
	"j-iam/internal/apiserver/store/mysql"
	"j-iam/internal/pkg/code"
	"j-iam/internal/pkg/middleware"
	"j-iam/internal/pkg/middleware/auth"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

// installMiddleware 注册通用（所有API都用的）中间件到 gin 路由引擎
func installMiddleware(g *gin.Engine) {
	//gin.SetMode(gin.DebugMode)
}

// installController 注册 RESTful API 到 gin 路由引擎
func installController(g *gin.Engine) *gin.Engine {
	// Middlewares.
	// 创建 jwt 认证中间件，并加载到 gin
	jwtStrategy, _ := newJWTAuth().(auth.JWTStrategy)

	// 登录 /login 需要 Basic 认证和 Bearer 认证
	g.POST("/login", jwtStrategy.LoginHandler)
	g.POST("/logout", jwtStrategy.LogoutHandler)

	// 令牌刷新，刷新令牌的有效时间要长于认证
	g.POST("/refresh", jwtStrategy.RefreshHandler)

	// 基于策略模式，创建可切换的认证策略集
	auto := newAutoAuth()

	g.NoRoute(auto.AuthFunc(), func(c *gin.Context) {
		core.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
	})

	// v1 handlers, requiring authentication
	storeIns, _ := mysql.GetMySQLFactory(nil)
	v1 := g.Group("/v1")
	{
		// user RESTful resource
		userv1 := v1.Group("/users")
		{
			userController := user.NewUserController(storeIns)

			userv1.POST("", userController.Create)

			// 卡这里了！认证不通过
			userv1.Use(auto.AuthFunc(), middleware.Validation()) // 为什么不让访问？？？
			//v1.PUT("/find_password", userController.FindPassword)
			userv1.DELETE("", userController.DeleteCollection) // admin api
			userv1.DELETE(":name", userController.Delete)      // admin api
			userv1.PUT(":name/change-password", userController.ChangePassword)
			userv1.PUT(":name", userController.Update)
			userv1.GET("", userController.List)
			userv1.GET(":name", userController.Get) // admin api
		}

		// 密钥和策略的访问需要认证通过才行
		v1.Use(auto.AuthFunc())

		// policy RESTful resource
		policyv1 := v1.Group("/policies", middleware.Publish())
		{
			policyController := policy.NewPolicyController(storeIns)

			policyv1.POST("", policyController.Create)
			policyv1.DELETE("", policyController.DeleteCollection)
			policyv1.DELETE(":name", policyController.Delete)
			policyv1.PUT(":name", policyController.Update)
			policyv1.GET("", policyController.List)
			policyv1.GET(":name", policyController.Get)
		}

		// secret RESTful resource
		secretv1 := v1.Group("/secrets", middleware.Publish())
		{
			secretController := secret.NewSecretController(storeIns)

			secretv1.POST("", secretController.Create)
			secretv1.DELETE(":name", secretController.Delete)
			secretv1.PUT(":name", secretController.Update)
			secretv1.GET("", secretController.List)
			secretv1.GET(":name", secretController.Get)
		}
	}

	return g
}
