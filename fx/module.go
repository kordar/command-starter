package command_starter

import (
	"context"
	"flag"
	"log/slog"

	"github.com/kordar/command"
	gocfgmodulefx "github.com/kordar/gocfg-load-module/fx/v2"
	"github.com/spf13/cast"
	"go.uber.org/fx"
)

// ModuleConfig command 模块的配置
type ModuleConfig struct {
	Name      string // 命令参数名，默认 "fname"
	Usage     string // 参数使用说明
	Value     string // 默认值
	FlagParse bool   // 是否自动调用 flag.Parse()（注意：和 cobra 一起用时建议设为 false）
}

// cfgModule 实现 gocfg-load-module/fx 的 GoCfgModule 接口
type cfgModule struct {
	name string
}

var _ gocfgmodulefx.GoCfgModule = cfgModule{}

// StarterModule 创建一个可被 gocfg-load-module/fx 加载的 starter
func StarterModule(name string) gocfgmodulefx.GoCfgModule {
	return cfgModule{name: name}
}

// Name 实现 GoCfgModule 接口
func (m cfgModule) Name() string {
	return m.name
}

// Load 实现 GoCfgModule 接口，解析配置并返回 []fx.Option
func (m cfgModule) Load(data any) []fx.Option {
	slog.Info("Module Load Complete", "module", "command-starter(fx)")
	cfg := buildModuleConfig(data)
	return []fx.Option{
		Module(cfg),
	}
}

// buildModuleConfig 从配置数据构造 ModuleConfig
func buildModuleConfig(data any) ModuleConfig {
	if data == nil {
		return ModuleConfig{}
	}

	cfg := ModuleConfig{
		Name:  "fname",
		Usage: "function name",
	}

	item := cast.ToStringMapString(data)
	if item["name"] != "" {
		cfg.Name = item["name"]
	}
	if item["usage"] != "" {
		cfg.Usage = item["usage"]
	}
	if item["value"] != "" {
		cfg.Value = item["value"]
	}
	if item["flag-parse"] != "" {
		cfg.FlagParse = cast.ToBool(item["flag-parse"])
	}

	return cfg
}

// Module 构造 command starter 的 fx 模块
func Module(cfg ModuleConfig) fx.Option {
	options := []fx.Option{
		fx.Supply(cfg),
		fx.Provide(
			provideCommandRegistry,
		),
	}

	// required 模式下允许无配置接入，此时只提供注册表，不绑定标准 flag。
	if shouldInitFlag(cfg) {
		options = append(options, fx.Invoke(initCommandStarter))
	}

	return fx.Module("command-starter", options...)
}

func shouldInitFlag(cfg ModuleConfig) bool {
	return cfg.FlagParse || cfg.Name != "" || cfg.Usage != "" || cfg.Value != ""
}

// CommandRegistry 用于聚合注入的 command.FuncCli
type CommandRegistry struct {
	Cli []command.FuncCli `group:"command-funccli"`
}

// provideCommandRegistry 把 grouped slice 聚合为一个 registry
func provideCommandRegistry(
	in struct {
		fx.In
		Cli []command.FuncCli `group:"command-funccli"`
	},
) CommandRegistry {
	return CommandRegistry{
		Cli: in.Cli,
	}
}

// CmdArgs 从 CLI (cobra/root) 传入的命令参数
type CmdArgs struct {
	Name    string   // 要执行的命令名称
	Args    []string // 命令附加参数
	RawArgs []string // root 收到的完整参数
}

// ExecuteCommand 在 fx 启动后执行指定的命令。
// 由 starter.CmdServerStarter 通过 fx.Invoke 触发。
func ExecuteCommand(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, args CmdArgs, reg CommandRegistry) {
	registerCommands(reg)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				exitCode := runCLICommand(args, reg)
				if err := shutdowner.Shutdown(fx.ExitCode(exitCode)); err != nil {
					slog.Error("关闭 Fx 失败", "name", args.Name, "error", err)
				}
			}()
			return nil
		},
	})
}

// StartFunc 暴露给外部调用，直接启动命令（避免和 cobra 等冲突）
func StartFunc(cmd string, reg CommandRegistry) {
	registerCommands(reg)
	command.StartCmd(cmd)
}

// initCommandStarter 按配置初始化 command starter（只用于 flag 模式）
func initCommandStarter(cfg ModuleConfig, reg CommandRegistry) {
	cfg = normalizeModuleConfig(cfg)

	registerCommands(reg)

	// 绑定 flag 参数
	var f string
	flag.StringVar(&f, cfg.Name, cfg.Value, cfg.Usage)

	// 仅在配置中明确开启时才 flag.Parse()
	if cfg.FlagParse {
		flag.Parse()
		// 启动命令
		command.StartCmd(f)
	}
}

func normalizeModuleConfig(cfg ModuleConfig) ModuleConfig {
	if cfg.Name == "" {
		cfg.Name = "fname"
	}
	if cfg.Usage == "" {
		cfg.Usage = "function name"
	}
	return cfg
}

func registerCommands(reg CommandRegistry) {
	for _, c := range reg.Cli {
		command.SetCmd(c)
	}
}

func runCLICommand(args CmdArgs, reg CommandRegistry) (exitCode int) {
	defer func() {
		if r := recover(); r != nil {
			exitCode = 1
			slog.Error("CLI 命令执行发生 panic", "name", args.Name, "panic", r)
		}
	}()

	if args.Name == "" {
		slog.Error("未提供待执行的命令名称")
		return 1
	}

	cmd, ok := findCommand(args.Name, reg)
	if !ok {
		slog.Error("待执行的程序不存在", "name", args.Name)
		return 1
	}

	if err := command.ExecuteFunc(cmd, args.RawArgs); err != nil {
		return 1
	}
	return 0
}

func findCommand(name string, reg CommandRegistry) (command.FuncCli, bool) {
	for _, c := range reg.Cli {
		if c.Name() == name {
			return c, true
		}
	}
	return nil, false
}
