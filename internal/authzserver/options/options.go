// Package options contains flags and options for initializing an apiserver
package options

import (
	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
	"github.com/marmotedu/component-base/pkg/json"

	"j-iam/internal/authzserver/analytics"
	genericoptions "j-iam/internal/pkg/options"
	"j-iam/internal/pkg/server"
	"j-iam/pkg/log"
)

const clientCa = "cert/server.pem"

// Options authzserver Options 配置
type Options struct {
	RPCServer               string                                 `json:"rpcserver"      mapstructure:"rpcserver"`
	ClientCA                string                                 `json:"client-ca-file" mapstructure:"client-ca-file"`
	GenericServerRunOptions *genericoptions.ServerRunOptions       `json:"server"         mapstructure:"server"`
	InsecureServing         *genericoptions.InsecureServingOptions `json:"insecure"       mapstructure:"insecure"`
	RedisOptions            *genericoptions.RedisOptions           `json:"redis"          mapstructure:"redis"`
	Log                     *log.Options                           `json:"log"            mapstructure:"log"`
	AnalyticsOptions        *analytics.AnalyticsOptions            `json:"analytics"      mapstructure:"analytics"`
	SecureServing           *genericoptions.SecureServingOptions   `json:"secure"         mapstructure:"secure"`
	//FeatureOptions   *genericoptions.FeatureOptions `json:"feature"        mapstructure:"feature"`
}

// NewOptions creates a new Options object with default parameters.
func NewOptions() *Options {
	o := Options{
		RPCServer: "127.0.0.1:8081",
		//ClientCA:                "",
		ClientCA:                clientCa, // 设置为临时值
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		InsecureServing:         genericoptions.NewInsecureServingOptions(),
		RedisOptions:            genericoptions.NewRedisOptions(),
		Log:                     log.NewOptions(),
		AnalyticsOptions:        analytics.NewAnalyticsOptions(),
		SecureServing:           genericoptions.NewSecureServingOptions(),
		//FeatureOptions:          genericoptions.NewFeatureOptions(),
	}

	o.GenericServerRunOptions.Middlewares = []string{"recovery", "logger", "nocache", "cors", "dump"}
	o.InsecureServing.BindPort = 8082

	return &o
}

// ApplyTo applies the run options to the method receiver and returns self.
func (o *Options) ApplyTo(c *server.Config) error {
	return nil
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.AnalyticsOptions.AddFlags(fss.FlagSet("analytics"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure serving"))
	o.Log.AddFlags(fss.FlagSet("logs"))
	o.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	//o.FeatureOptions.AddFlags(fss.FlagSet("features"))

	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := fss.FlagSet("misc")
	fs.StringVar(&o.RPCServer, "rpcserver", o.RPCServer, "The address of iam rpc server. "+
		"The rpc server can provide all the secrets and policies to use.")
	fs.StringVar(&o.ClientCA, "client-ca-file", o.ClientCA, ""+
		"If set, any request presenting a client certificate signed by one of "+
		"the authorities in the client-ca-file is authenticated with an identity "+
		"corresponding to the CommonName of the client certificate.")

	return fss
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

// Complete set default Options.
func (o *Options) Complete() error {
	//return o.SecureServing.Complete()
	return nil
}

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.GenericServerRunOptions.Validate()...)
	errs = append(errs, o.InsecureServing.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.AnalyticsOptions.Validate()...)
	errs = append(errs, o.SecureServing.Validate()...)
	//errs = append(errs, o.FeatureOptions.Validate()...)

	return errs
}
