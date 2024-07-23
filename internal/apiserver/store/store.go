package store

// 存储层服务单例模式
var client Factory

// Factory 定义 apiserver 存储层服务，抽象工厂方法模式，解耦了仓库层的实现和调用
type Factory interface {
	Users() UserStore
	Secrets() SecretStore
	Policies() PolicyStore
	PolicyAudits() PolicyAuditStore
	Close() error
}

// Client 返回仓库层服务实例
func Client() Factory { return client }

func SetClient(factory Factory) {
	client = factory
}
