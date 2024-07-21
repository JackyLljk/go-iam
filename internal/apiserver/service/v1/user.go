package v1

import (
	"context"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"j-iam/internal/apiserver/store"
	v1 "j-iam/pkg/model/v1"
)

// UserService 处理用户请求的业务（抽象业务层用户业务对象）
type UserService interface {
	Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error
	List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error)
}

type userService struct {
	store store.Factory
}

// 强制确保 userService 实现了 UserService 接口
var _ UserService = (*userService)(nil)

func (u *userService) Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error {
	return nil
}

func (u *userService) List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
	return nil, nil
}
