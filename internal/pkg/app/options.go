package app

import cliflag "github.com/marmotedu/component-base/pkg/cli/flag"

type CliOptions interface {
	Flags() (fss cliflag.NamedFlagSets)
	Validate() []error
}
