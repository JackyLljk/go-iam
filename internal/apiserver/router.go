package apiserver

import (
	"github.com/gin-gonic/gin"
	"j-iam/internal/apiserver/controller/v1/user"
	"j-iam/internal/apiserver/store/mysql"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

// installMiddleware 注册通用（所有API都用的）中间件到 gin 路由引擎
func installMiddleware(g *gin.Engine) {}

// installController 注册 RESTful API 到 gin 路由引擎
func installController(g *gin.Engine) {

	storsIns, _ := mysql.GetMysqlFactory(nil)
	v1 := g.Group("/v1")
	{
		//	// user RESTful resource
		userV1 := v1.Group("/users")
		{
			userController := user.NewUserController(storeIns)

			userV1.POST("", userController.Create)
			userV1.GET("", userController.List)
			//		userv1.GET(":name", userController.Get) // admin api
			//	}
			//
			//	v1.Use(auto.AuthFunc())
		}
	}
}
