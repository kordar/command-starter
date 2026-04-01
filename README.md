# command-starter

一个用于接入 github.com/kordar/command 的轻量启动器，负责：
- 注册实现了 command.FuncCli 的命令
- 配置并解析命令行参数（选择要执行的命令名称）
- 启动对应的命令逻辑

适配 Go 1.21。

## 安装

```bash
go get github.com/kordar/command-starter
```

## 快速开始

实现一个命令并注册，然后通过命令行选择执行的函数名。

```go
package main

import (
	command_starter "github.com/kordar/command-starter"
	"github.com/kordar/command"
)

type Hello struct{}

func (Hello) Name() string { return "hello" }
func (Hello) Execute(args ...interface{}) { println("hello world") }
func (Hello) GetArgs() []interface{} { return []interface{}{} }

func main() {
	command_starter.AddCli(Hello{})
	m := command_starter.CommandModule{}
	m.Load(map[string]string{
		"name":       "f",
		"usage":      "function name",
		"value":      "",
		"flag-parse": "true",
	})
}
```

运行：

```bash
go run . -f hello
```

## 参数配置

- name：命令行参数名，默认 fname
- usage：参数说明文案，默认 function name
- value：默认值（为空时不执行）
- flag-parse：是否在 Load 调用时执行 flag.Parse，字符串 "true"/"false"

## 注册多个命令

```go
command_starter.AddCli(Foo{}, Bar{}, Baz{})
```

## 注意事项

- 执行完成后会调用 os.Exit(0) 结束进程
- 底层依赖 github.com/kordar/command，请确保版本兼容
