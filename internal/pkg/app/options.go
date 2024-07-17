package app

import cliflag "github.com/marmotedu/component-base/pkg/cli/flag"

// 定义应用框架的各种配置的抽象接口

// CliOptions 抽象接口，用于从命令行读取参数的配置
type CliOptions interface {
	Flags() (fss cliflag.NamedFlagSets)
	Validate() []error
}

// ConfigureOptions 抽象接口，用于从配置文件中读取参数的配置
type ConfigureOptions interface {
	ApplyFlags() []error
}

// CompleteableOptions 抽象出可以完成的配置？？？
type CompleteableOptions interface {
	Complete() error
}

// PrintableOptions 抽象接口，可打印的配置
type PrintableOptions interface {
	String() string
}
