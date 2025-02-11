package authzserver

import (
	"github.com/marmotedu/errors"

	"j-iam/internal/authzserver/load/cache"
	"j-iam/internal/pkg/middleware"
	"j-iam/internal/pkg/middleware/auth"
)

// 创建 cache 认证策略
func newCacheAuth() middleware.AuthStrategy {
	return auth.NewCacheStrategy(getSecretFunc())
}

// 返回密钥信息
func getSecretFunc() func(string) (auth.Secret, error) {
	return func(kid string) (auth.Secret, error) {
		cli, err := cache.GetCacheInsOr(nil)
		if err != nil || cli == nil {
			return auth.Secret{}, errors.Wrap(err, "get cache instance failed")
		}

		secret, err := cli.GetSecret(kid)
		if err != nil {
			return auth.Secret{}, err
		}

		// 返回的密钥包括以下字段
		return auth.Secret{
			Username: secret.Username,
			ID:       secret.SecretId,
			Key:      secret.SecretKey,
			Expires:  secret.Expires,
		}, nil
	}
}
