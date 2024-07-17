package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/marmotedu/component-base/pkg/cli/globalflag"
	"github.com/marmotedu/component-base/pkg/term"
	"github.com/marmotedu/component-base/pkg/version"
	"github.com/marmotedu/component-base/pkg/version/verflag"
	"github.com/marmotedu/log"
	"os"

	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 基于选项模式(Options Pattern)，构建应用配置

// App cli 应用架构
type App struct {
	basename    string               // 命令名称 app name (此处即 "iam-apiserver")
	name        string               // 简短描述
	description string               // 详细描述
	options     CliOptions           // 命令行参数接口
	runFunc     RunFunc              // 启动应用的回调函数
	silence     bool                 // 静默模式(true 则程序启动信息、配置信息和版本信息不会在控制台中打印)
	noVersion   bool                 // 设置应用程序是否提供版本 flag
	noConfig    bool                 // 设置应用程序是否提供配置 flag
	commands    []*Command           // cli 应用的子命令
	args        cobra.PositionalArgs // 函数类型 func(cmd *Command, args []string) error
	cmd         *cobra.Command       // Cobra 库的命令参数
}

// Option 定义用于初始化应用程序架构的可选字段（基于选项模式配置可选字段）
type Option func(*App)

// 如果需要附加配置字段时的逻辑，可以使用下面的写法
//type Option interface {
//	apply(*App)
//}
//
//type optionFunc func(*App)
//
//func (f optionFunc) apply(a *App) {
//	//添加配置逻辑
//	f(a)
//}

// RunFunc 定义启动应用的
type RunFunc func(basename string) error

// WithXXX 函数，在创建 App 时为配置的每个字段设置值

func WithRunFunc(run RunFunc) Option {
	return func(a *App) {
		a.runFunc = run
	}
}

// WithOptions 从配置或命令行中，读取并设置应用框架配置
func WithOptions(opt CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

func WithDescription(desc string) Option {
	return func(a *App) {
		a.description = desc
	}
}

// WithDefaultValidArgs 默认配置为不需要传入命令行参数，只使用根命令 run
// 当你运行命令并提供任何参数时，程序将输出错误消息并退出，因为 WithDefaultValidArgs 配置的验证逻辑不允许任何参数
func WithDefaultValidArgs() Option {
	return func(a *App) {
		a.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		}
	}
}

func WithSilence() Option {
	return func(a *App) {
		a.silence = true
	}
}

// WithNoConfig 设置应用框架不提供 Config
func WithNoConfig() Option {
	return func(a *App) {
		a.noConfig = true
	}
}

// WithValidArgs 配置有效命令行参数(如果命令参数无效，则不会执行)
func WithValidArgs(args cobra.PositionalArgs) Option {
	return func(a *App) {
		a.args = args
	}
}

func WithNoVersion() Option {
	return func(a *App) {
		a.noVersion = true
	}
}

// NewApp 基于选项模式配置应用架构可选字段
func NewApp(name string, basename string, opts ...Option) *App {
	// 初始化默认应用
	a := &App{
		name:     name,
		basename: basename,
	}

	// 选项模式：根据 opts 参数设置其他选项（不设置即为默认值）
	// 这一步获取绑定了可选字段的闭包（回调函数），用于配置 App 的可选字段
	for _, o := range opts {
		o(a)
	}

	// 根据配置选项后的应用，构建命令行命令
	a.buildCommand()

	return a
}

func (a *App) buildCommand() {
	// 创建命令行命令
	cmd := cobra.Command{
		Use:   FormatBaseName(a.basename),
		Short: a.name,
		Long:  a.description,
		// stop printing usage when the command errors(设置发生错误时不打印使用信息）
		SilenceUsage:  true,
		SilenceErrors: true,

		// 将 App 中的命令参数设置在命令中
		Args: a.args,
	}

	// 可以自定义无效命令的 usage
	// cmd.SetUsageTemplate(usageTemplate)

	// 设置标准输入输出
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)

	// 初始化flags，且标志是排序的
	cmd.Flags().SortFlags = true
	cliflag.InitFlags(cmd.Flags())

	// 检查是否有子命令，有则添加到父命令 cmd
	// TODO: 子命令哪来的？
	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		// 设置自定义的 help 命令
		cmd.SetHelpCommand(helpCommand(FormatBaseName(a.basename)))
	}

	// 注册了应用启动回调函数，
	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	// 将 options 的 flags 添加到 cobra 实例的 FlagSet 中
	var namedFlagSets cliflag.NamedFlagSets
	if a.options != nil {
		// 这里的 Options 实现了 Flags() (fss cliflag.NamedFlagSets) 方法
		namedFlagSets = a.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}
	}

	// 添加版本相关 flag 到 global flagset 中
	if !a.noVersion {
		verflag.AddFlags(namedFlagSets.FlagSet("global"))
	}

	// 添加配置相关 flag 到 global flagset 中
	//if !a.noConfig {
	//	addConfigFlag(a.basename, namedFlagSets.FlagSet("global"))
	//}

	globalflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())
	// add new global flagset to cmd FlagSet
	cmd.Flags().AddFlagSet(namedFlagSets.FlagSet("global"))

	// 添加命令模板到 flagset 中
	addCmdTemplate(&cmd, namedFlagSets)

	// 注册 cobra.Command 到 App
	a.cmd = &cmd
}

func (a *App) Run() {
	// cmd.Execute() 调用命令的执行入口
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
	fmt.Println("link start!")
}

func (a *App) Command() *cobra.Command {
	return a.cmd
}

func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	printWorkingDir()
	cliflag.PrintFlags(cmd.Flags())
	if !a.noVersion {
		// display application version information
		verflag.PrintAndExitIfRequested()
	}

	if !a.noConfig {
		// 绑定 key 到 Flag (viper 支持 pflag 包)
		// 绑定后，viper 可以使用标志的 name 访问其值 viper.GetXXX(name)
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}

		// 将 options 反序列化
		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	if !a.silence {
		log.Infof("%v Starting %s ...", progressMessage, a.name)
		if !a.noVersion {
			log.Infof("%v Version: `%s`", progressMessage, version.Get().ToJSON())
		}
		if !a.noConfig {
			log.Infof("%v Config file used: `%s`", progressMessage, viper.ConfigFileUsed())
		}
	}

	if a.options != nil {
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}

	// run application
	// 启动应用（internal/apiserver/app.go run()）
	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}

	return nil
}

func (a *App) applyOptionRules() error {
	//if completeableOptions, ok := a.options.(CompleteableOptions); ok {
	//	if err := completeableOptions.Complete(); err != nil {
	//		return err
	//	}
	//}
	//
	//if errs := a.options.Validate(); len(errs) != 0 {
	//	return errors.NewAggregate(errs)
	//}
	//
	//if printableOptions, ok := a.options.(PrintableOptions); ok && !a.silence {
	//	log.Infof("%v Config: `%s`", progressMessage, printableOptions.String())
	//}

	return nil
}

// printWorkingDir 获取当前工作目录并打印
func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Infof("%v WorkingDir: %s", progressMessage, wd)
}

// addCmdTemplate 添加命令模板到 cmd
func addCmdTemplate(cmd *cobra.Command, namedFlagSets cliflag.NamedFlagSets) {
	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		_, _ = fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)

		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
}
