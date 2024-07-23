package user

import (
	srvv1 "j-iam/internal/apiserver/service/v1"
	"j-iam/internal/apiserver/store"
)

// UserController 创建 user handler 控制层接口，处理 user 相关请求
type UserController struct {
	srv srvv1.Service
}

func NewUserController(store store.Factory) *UserController {
	return &UserController{
		srv: srvv1.NewService(store),
	}
}
