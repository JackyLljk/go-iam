package v1

import "j-iam/internal/apiserver/store"

// Service defines functions used to return resource interface. 定义用于返回资源接口的函数
type Service interface {
	Users() UserService
	Secrets() SecretService
	Policies() PolicyService
}

type service struct {
	store store.Factory
}

// NewService returns Service interface.
func NewService(store store.Factory) Service {
	return &service{
		store: store,
	}
}

func (s *service) Users() UserService {
	return newUsers(s)
}

func (s *service) Secrets() SecretService {
	return newSecrets(s)
}

func (s *service) Policies() PolicyService {
	return newPolicies(s)
}
