package store

// 存储层服务单例模式
var client Factory

// Factory 定义 apiserver 存储层服务
type Factory interface {
	Users() UserStore
	Close() error
	// Secrets() SecretStore
	// Policies() PolicyStore
	// PolicyAudits() PolicyAuditStore
}

func Client() Factory { return client }

func SetClient(factory Factory) {
	client = factory
}
