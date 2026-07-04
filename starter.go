package command_starter

import (
	"flag"

	"github.com/kordar/command"
)

// Config 描述 command starter 的行为。
type Config struct {
	Name      string // 命令行参数名，默认 "fname"
	Usage     string // 参数使用说明，默认 "function name"
	Value     string // 默认值
	FlagParse bool   // 是否自动调用 flag.Parse()
}

// DefaultConfig 返回默认配置。
func DefaultConfig() Config {
	return Config{
		Name:  "fname",
		Usage: "function name",
	}
}

// AddCli 注册一组 command.FuncCli 实现。
func AddCli(cli ...command.FuncCli) {
	for _, c := range cli {
		command.SetCmd(c)
	}
}

// Start 绑定 flag 并启动命令。
// 若 cfg.FlagParse 为 true，则自动调用 flag.Parse() 解析命令行参数。
// 否则只注册 flag 定义，由调用方决定何时 Parse。
func Start(cfg Config) {
	if cfg.Name == "" {
		cfg.Name = "fname"
	}
	if cfg.Usage == "" {
		cfg.Usage = "function name"
	}

	var f string
	flag.StringVar(&f, cfg.Name, cfg.Value, cfg.Usage)

	if cfg.FlagParse {
		flag.Parse()
		command.StartCmd(f)
	}
}

// StartWithName 是 Start 的便捷包装，仅指定参数名即启动。
func StartWithName(name string) {
	Start(Config{Name: name, FlagParse: true})
}
