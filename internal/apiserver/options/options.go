package options

import (
	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
	genericoptions "j-iam/internal/pkg/options"
)

// Options 基本参数（是一个配置数据结构，可用来构建应用框架，也作为应用配置的输入）
// Options 实现了 CliOptions，可以在构建 cobra 命令行应用框架时使用
type Options struct {
	GenericServerRunOptions *genericoptions.ServerRunOptions `json:"service"   mapstructure:"service"`
	//GRPCOptions             *genericoptions.GRPCOptions            `json:"grpc"     mapstructure:"grpc"`
	InsecureServing *genericoptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	//SecureServing           *genericoptions.SecureServingOptions   `json:"secure"   mapstructure:"secure"`
	MySQLOptions *genericoptions.MySQLOptions `json:"mysql"    mapstructure:"mysql"`
	//RedisOptions            *genericoptions.RedisOptions           `json:"redis"    mapstructure:"redis"`
	//JwtOptions              *genericoptions.JwtOptions             `json:"jwt"      mapstructure:"jwt"`
	//Log                     *log.Options                           `json:"log"      mapstructure:"log"`
	//FeatureOptions          *genericoptions.FeatureOptions         `json:"feature"  mapstructure:"feature"`
}

func NewOptions() *Options {
	o := Options{
		//GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		//GRPCOptions:             genericoptions.NewGRPCOptions(),
		//InsecureServing:         genericoptions.NewInsecureServingOptions(),
		//SecureServing:           genericoptions.NewSecureServingOptions(),
		//MySQLOptions:            genericoptions.NewMySQLOptions(),
		//RedisOptions:            genericoptions.NewRedisOptions(),
		//JwtOptions:              genericoptions.NewJwtOptions(),
		//Log:                     log.NewOptions(),
		//FeatureOptions:          genericoptions.NewFeatureOptions(),
	}

	return &o
}

func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	//o.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))
	//o.JwtOptions.AddFlags(fss.FlagSet("jwt"))
	//o.GRPCOptions.AddFlags(fss.FlagSet("grpc"))
	//o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	//o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	//o.FeatureOptions.AddFlags(fss.FlagSet("features"))
	//o.InsecureServing.AddFlags(fss.FlagSet("insecure serving"))
	//o.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	//o.Log.AddFlags(fss.FlagSet("logs"))

	return fss
}

func (o *Options) Validate() []error {
	var errs []error

	//errs = append(errs, o.GenericServerRunOptions.Validate()...)
	//errs = append(errs, o.GRPCOptions.Validate()...)
	//errs = append(errs, o.InsecureServing.Validate()...)
	//errs = append(errs, o.SecureServing.Validate()...)
	//errs = append(errs, o.MySQLOptions.Validate()...)
	//errs = append(errs, o.RedisOptions.Validate()...)
	//errs = append(errs, o.JwtOptions.Validate()...)
	//errs = append(errs, o.Log.Validate()...)
	//errs = append(errs, o.FeatureOptions.Validate()...)

	return errs
}
