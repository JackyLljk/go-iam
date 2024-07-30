package cache

import (
	"sync"

	pb "j-iam/internal/pkg/proto/apiserver/v1"

	"github.com/dgraph-io/ristretto"
	"github.com/marmotedu/errors"
	"github.com/ory/ladon"

	"j-iam/internal/authzserver/store"
)

// Cache 存储 secrets 和 policies
type Cache struct {
	lock     *sync.RWMutex
	client   store.Factory
	secrets  *ristretto.Cache
	policies *ristretto.Cache
}

var (
	// ErrSecretNotFound defines secret not found error.
	ErrSecretNotFound = errors.New("secret not found")
	// ErrPolicyNotFound defines policy not found error.
	ErrPolicyNotFound = errors.New("policy not found")
)

var (
	onceCache sync.Once
	cacheIns  *Cache
)

// GetCacheInsOr return store instance.
func GetCacheInsOr(client store.Factory) (*Cache, error) {
	var err error
	if client != nil {
		var (
			secretCache *ristretto.Cache
			policyCache *ristretto.Cache
		)

		onceCache.Do(func() {
			c := &ristretto.Config{
				NumCounters: 1e7,     // number of keys to track frequency of (10M).
				MaxCost:     1 << 30, // maximum cost of cache (1GB).
				BufferItems: 64,      // number of keys per Get buffer.
				Cost:        nil,
			}

			secretCache, err = ristretto.NewCache(c)
			if err != nil {
				return
			}
			policyCache, err = ristretto.NewCache(c)
			if err != nil {
				return
			}

			cacheIns = &Cache{
				client:   client,
				lock:     new(sync.RWMutex),
				secrets:  secretCache,
				policies: policyCache,
			}
		})
	}

	return cacheIns, err
}

// GetSecret return secret detail for the given key.
func (c *Cache) GetSecret(key string) (*pb.SecretInfo, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.secrets.Get(key)
	if !ok {
		return nil, ErrSecretNotFound
	}

	return value.(*pb.SecretInfo), nil
}

// GetPolicy return user's ladon policies for the given user.
func (c *Cache) GetPolicy(key string) ([]*ladon.DefaultPolicy, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.policies.Get(key)
	if !ok {
		return nil, ErrPolicyNotFound
	}

	return value.([]*ladon.DefaultPolicy), nil
}

// Reload 加载 secrets 和 policies (启动 authzserver 和 订阅到更新后重新加载)
// TODO: 看看是不是全量加载，后续应该要改成增量更新，看看怎么实现？
func (c *Cache) Reload() error {
	// TODO：使用 gRPC 加载数据到内存，为什么要加锁？
	c.lock.Lock()
	defer c.lock.Unlock()

	// reload secrets
	secrets, err := c.client.Secrets().List()
	if err != nil {
		return errors.Wrap(err, "list secrets failed")
	}

	c.secrets.Clear()
	for key, val := range secrets {
		c.secrets.Set(key, val, 1)
	}

	// reload policies
	policies, err := c.client.Policies().List()
	if err != nil {
		return errors.Wrap(err, "list policies failed")
	}

	c.policies.Clear()
	for key, val := range policies {
		c.policies.Set(key, val, 1)
	}

	return nil
}
